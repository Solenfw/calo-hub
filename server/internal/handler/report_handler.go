// Package handler contains the request handlers for the report-oriented routes.
//
// These endpoints bridge the database-backed Martin report metadata and the PDF
// generation service. The goal is to keep report lookup, validation, and output
// creation in one place so the route layer remains simple and the business
// behavior stays easy to test.
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/solenfw/calo-hub/internal/dto"
	"github.com/solenfw/calo-hub/internal/models"
	"github.com/solenfw/calo-hub/internal/repository"
	"github.com/solenfw/calo-hub/internal/service/pdf"
)

// ReportStore defines the report operations required by the HTTP layer.
//
// Keeping the interface in the consumer package allows the handlers to be tested
// with lightweight fakes instead of a real repository object.
type ReportStore interface {
	GetAllMartinReportLists(ctx context.Context) ([]models.MartinReportList, error)
	GetMartinReportListByName(ctx context.Context, name string) (models.MartinReportList, error)
	GetMartinReportProductsByReportID(ctx context.Context, reportID int) ([]models.MartinReportProduct, error)
	CreateMartinReportList(ctx context.Context, name string) (models.MartinReportList, error)
	DeleteMartinReportList(ctx context.Context, reportID int) error
	SetMartinReportProducts(ctx context.Context, reportID int, products []models.MartinReportProduct) error
}

var _ ReportStore = (*repository.ReportRepository)(nil)

// ReportHandler serves the Martin report endpoints that resolve report metadata,
// product snapshots, and generated PDF output.
type ReportHandler struct {
	store   ReportStore
	catalog CatalogStore
}

// NewReportHandler wires a report-capable store into the HTTP handler.
func NewReportHandler(store ReportStore, catalogStores ...CatalogStore) *ReportHandler {
	var catalog CatalogStore
	if len(catalogStores) > 0 {
		catalog = catalogStores[0]
	}

	return &ReportHandler{
		store:   store,
		catalog: catalog,
	}
}

func (h *ReportHandler) resolveReportProductImages(ctx context.Context, products []models.MartinReportProduct) ([]models.MartinReportProduct, error) {
	if h.catalog == nil || len(products) == 0 {
		return products, nil
	}

	resolved := make([]models.MartinReportProduct, len(products))
	copy(resolved, products)

	for i, product := range resolved {
		if product.Image != nil && strings.TrimSpace(*product.Image) != "" {
			continue
		}
		if strings.TrimSpace(product.Code) == "" {
			continue
		}

		images, err := h.catalog.GetImagesByCode(ctx, product.Code)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				continue
			}
			return nil, err
		}
		if len(images.Images) == 0 {
			continue
		}

		imageURL := strings.TrimSpace(images.Images[0])
		if imageURL == "" {
			continue
		}
		resolved[i].Image = &imageURL
	}

	return resolved, nil
}

// ListMartinReports returns all available Martin report lists.
func (h *ReportHandler) ListMartinReports(w http.ResponseWriter, r *http.Request) {
	reports, err := h.store.GetAllMartinReportLists(r.Context())
	if err != nil {
		writeRepositoryError(w, err)
		return
	}

	response := make([]dto.MartinReportListResponse, 0, len(reports))
	for _, report := range reports {
		response = append(response, dto.MartinReportListResponse{
			ID:   report.ID,
			Name: report.Name,
		})
	}

	writeJSON(w, http.StatusOK, response)
}

// GetMartinReportByName returns metadata for a single Martin report list.
func (h *ReportHandler) GetMartinReportByName(w http.ResponseWriter, r *http.Request) {
	reportName := strings.TrimSpace(chi.URLParam(r, "name"))
	if reportName == "" {
		writeError(w, http.StatusBadRequest, "report name is required")
		return
	}

	report, err := h.store.GetMartinReportListByName(r.Context(), reportName)
	if err != nil {
		writeRepositoryError(w, err)
		return
	}

	resp := dto.MartinReportListResponse{
		ID:   report.ID,
		Name: report.Name,
	}
	writeJSON(w, http.StatusOK, resp)
}

// CreateMartinReport creates a new Martin report list row using the supplied name.
func (h *ReportHandler) CreateMartinReport(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "report name is required")
		return
	}

	payload.Name = strings.TrimSpace(payload.Name)
	if payload.Name == "" {
		writeError(w, http.StatusBadRequest, "report name is required")
		return
	}

	report, err := h.store.CreateMartinReportList(r.Context(), payload.Name)
	if err != nil {
		writeRepositoryError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, dto.MartinReportListResponse{ID: report.ID, Name: report.Name})
}

// DeleteMartinReport removes a report by its database ID and the linked product rows.
func (h *ReportHandler) DeleteMartinReport(w http.ResponseWriter, r *http.Request) {
	reportID := strings.TrimSpace(chi.URLParam(r, "report_id"))
	id, err := strconv.Atoi(reportID)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid report ID")
		return
	}

	if err := h.store.DeleteMartinReportList(r.Context(), id); err != nil {
		writeRepositoryError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// UpdateMartinReportProducts replaces the product rows attached to one report.
func (h *ReportHandler) UpdateMartinReportProducts(w http.ResponseWriter, r *http.Request) {
	reportID := strings.TrimSpace(chi.URLParam(r, "report_id"))
	id, err := strconv.Atoi(reportID)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid report ID")
		return
	}

	var payload dto.MartinReportProductsRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid product payload")
		return
	}

	products := make([]models.MartinReportProduct, 0, len(payload.Products))
	for idx, item := range payload.Products {
		products = append(products, models.MartinReportProduct{
			ReportID: 		id,
			RowNo:    		item.RowNo,
			Code:     		item.Code,
			Description: 	item.Description,
			Image:    		&item.Image,
			Quantity: item.Quantity,
		})
		if products[len(products)-1].RowNo <= 0 {
			products[len(products)-1].RowNo = idx + 1
		}
	}

	if err := h.store.SetMartinReportProducts(r.Context(), id, products); err != nil {
		writeRepositoryError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

// GetMartinReportProducts returns all products attached to a Martin report.
func (h *ReportHandler) GetMartinReportProducts(w http.ResponseWriter, r *http.Request) {
	reportID := strings.TrimSpace(chi.URLParam(r, "report_id"))

	reportIDInt, err := strconv.Atoi(reportID)
	if err != nil || reportIDInt <= 0 {
		writeError(w, http.StatusBadRequest, "invalid report ID")
		return
	}

	products, err := h.store.GetMartinReportProductsByReportID(r.Context(), reportIDInt)
	if err != nil {
		writeRepositoryError(w, err)
		return
	}

	products, err = h.resolveReportProductImages(r.Context(), products)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to resolve report images")
		return
	}

	resp := make([]dto.MartinReportProductResponse, 0, len(products))
	for _, product := range products {
		resp = append(resp, dto.MartinReportProductResponse{
			RowNo:    			product.RowNo,
			Code:     			product.Code,
			Description:      product.Description,
			Image:    			product.Image,
			Quantity: 			product.Quantity,
		})
	}

	writeJSON(w, http.StatusOK, resp)
}

// generate a PDF Report using prep data
func (h *ReportHandler) GetMartinReport(w http.ResponseWriter, r *http.Request) {
	reportName := strings.TrimSpace(chi.URLParam(r, "name"))
	if reportName == "" {
		writeError(w, http.StatusBadRequest, "report name is required")
		return
	}

	report, err := h.store.GetMartinReportListByName(r.Context(), reportName)
	if err != nil {
		writeRepositoryError(w, err)
		return
	}

	products, err := h.store.GetMartinReportProductsByReportID(r.Context(), report.ID)
	if err != nil {
		writeRepositoryError(w, err)
		return
	}

	products, err = h.resolveReportProductImages(r.Context(), products)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to resolve report images")
		return
	}

	pdfBytes, err := pdf.GeneratePDF(r.Context(), "", report.Name, products)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate PDF")
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(pdfBytes)
}

func resolveImportTemplatePath() (string, error) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", os.ErrNotExist
	}

	serverRoot := filepath.Dir(filepath.Dir(filepath.Dir(currentFile)))
	candidates := []string{
		filepath.Join(serverRoot, "data", "sheet", "template.xlsx"),
		filepath.Join(serverRoot, "..", "data", "sheet", "template.xlsx"),
		filepath.Join(".", "data", "sheet", "template.xlsx"),
		filepath.Join(".", "server", "data", "sheet", "template.xlsx"),
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	return "", os.ErrNotExist
}

func (h *ReportHandler) DownloadImportTemplate(w http.ResponseWriter, r *http.Request) {
	templatePath, err := resolveImportTemplatePath()
	if err != nil {
		writeError(w, http.StatusNotFound, "template file not found")
		return
	}

	data, err := os.ReadFile(templatePath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to read template file")
		return
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=martin-import-template.xlsx")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func (h *ReportHandler) RegisterRoutes(r chi.Router) {
	r.Get("/report/martin", h.ListMartinReports)
	r.Get("/report/martin/template", h.DownloadImportTemplate)
	r.Get("/report/martin/{name}", h.GetMartinReportByName)
	r.Post("/report/martin", h.CreateMartinReport)
	r.Delete("/report/martin/{report_id}", h.DeleteMartinReport)
	r.Put("/report/martin/{report_id}/products", h.UpdateMartinReportProducts)
	r.Get("/report/martin/{name}/pdf", h.GetMartinReport)
	r.Get("/report/martin/all/{report_id}", h.GetMartinReportProducts)
}
