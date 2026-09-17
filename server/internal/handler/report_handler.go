package handler

// endpoints for /report
import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/solenfw/calo-hub/internal/dto"
	"github.com/solenfw/calo-hub/internal/models"
	"github.com/solenfw/calo-hub/internal/repository"
)




// CatalogReportStore defines the catalog_report operations required by the HTTP handlers.
// Keeping this interface in the consumer package allows handlers to use fakes in unit tests.
type CatalogReportStore interface {
	GetMartinReportListByName(ctx context.Context, name string) (models.MartinReportList, error)
	GetMartinReportProductsByReportID(ctx context.Context, reportID int) ([]models.MartinReportProduct, error)
}

var _ CatalogReportStore = (*repository.CatalogRepository)(nil)

// ReportHandler serves catalog product, image, and Martin report endpoints
type ReportHandler struct {
	store CatalogReportStore
}

// NewReportHandler wires a catalog store into its HTTP handler.
func NewReportHandler(store CatalogReportStore) *ReportHandler {
	return &ReportHandler {
		store: store,
	}
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
		Name: report.Name,
	}
	writeJSON(w, http.StatusOK, resp)
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

	resp := make([]dto.MartinReportProductResponse, 0, len(products))
	for _, product := range products {
		resp = append(resp, dto.MartinReportProductResponse{
			RowNo:    product.RowNo,
			Code:     product.Code,
			Eng:      product.Eng,
			Image:    product.Image,
			Quantity: product.Quantity,
		})
	}

	writeJSON(w, http.StatusOK, resp)
}


func (h* ReportHandler) RegisterRoutes (r chi.Router) {
	r.Get("/catalog/report/martin/{name}", h.GetMartinReportByName)
	r.Get("/catalog/report/martin/all/{report_id}", h.GetMartinReportProducts)
}
