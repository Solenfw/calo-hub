// Package handler contains HTTP handlers for the API.
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/solenfw/calo-hub/internal/dto"
	"github.com/solenfw/calo-hub/internal/models"
	"github.com/solenfw/calo-hub/internal/repository"
)

var productCodePattern = regexp.MustCompile(`^\d{2}-\d{3}-\d{2}-\d{2}$|^[A-Z]{2}(?:-[A-Z])?\d{3,6}[A-Z]{0,2}$`)

// CatalogStore defines the catalog operations required by the HTTP handlers.
// Keeping this interface in the consumer package allows handlers to use fakes in unit tests.
type CatalogStore interface {
	GetProductByCode(ctx context.Context, code string) (models.Products, error)
	GetProductsByTerms(ctx context.Context, terms []string) ([]models.Products, error)
	GetImagesByCode(ctx context.Context, code string) (models.Images, error)
	GetAllMartinReportLists(ctx context.Context) ([]models.MartinReportList, error)
}

var _ CatalogStore = (*repository.CatalogRepository)(nil)

// CatalogHandler serves catalog product, image, and Martin report endpoints.
type CatalogHandler struct {
	store CatalogStore
}

// NewCatalogHandler wires a catalog store into its HTTP handler.
func NewCatalogHandler(store CatalogStore) *CatalogHandler {
	return &CatalogHandler{
		store: store,
	}
}

// GetProductByCode returns a single product for an exact catalog code.
func (h *CatalogHandler) GetProductByCode(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimSpace(chi.URLParam(r, "code"))
	if code == "" {
		writeError(w, http.StatusBadRequest, "empty input argument")
		return
	}

	if !productCodePattern.MatchString(code) {
		return
	}

	product, err := h.store.GetProductByCode(r.Context(), code)
	if err != nil {
		writeRepositoryError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toProductResponse(product))
}

// SearchProducts returns products that match all query terms in the product text.
func (h *CatalogHandler) SearchProducts(w http.ResponseWriter, r *http.Request) {
	searchTerm := strings.TrimSpace(r.URL.Query().Get("q"))
	if searchTerm == "" {
		searchTerm = strings.TrimSpace(r.URL.Query().Get("search"))
	}

	terms := strings.Fields(searchTerm)
	if len(terms) == 0 {
		writeError(w, http.StatusBadRequest, "empty search query")
		return
	}

	products, err := h.store.GetProductsByTerms(r.Context(), terms)
	if err != nil {
		writeRepositoryError(w, err)
		return
	}

	resp := make([]dto.ProductResponse, 0, len(products))
	for _, product := range products {
		resp = append(resp, toProductResponse(product))
	}

	writeJSON(w, http.StatusOK, resp)
}

// GetImages returns all stored image URLs for a product code.
func (h *CatalogHandler) GetImages(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimSpace(chi.URLParam(r, "code"))
	if code == "" {
		writeError(w, http.StatusBadRequest, "product code is required")
		return
	}

	images, err := h.store.GetImagesByCode(r.Context(), code)
	if err != nil {
		return
	}

	imageURLs := images.Images
	if imageURLs == nil {
		imageURLs = []string{}
	}

	resp := dto.ImageResponse{
		Code:   images.Code,
		Images: imageURLs,
	}
	writeJSON(w, http.StatusOK, resp)
}

// ListMartinReports returns the available Martin report lists.
func (h *CatalogHandler) ListMartinReports(w http.ResponseWriter, r *http.Request) {
	reports, err := h.store.GetAllMartinReportLists(r.Context())
	if err != nil {
		writeRepositoryError(w, err)
		return
	}

	resp := make([]dto.MartinReportListResponse, 0, len(reports))
	for _, report := range reports {
		resp = append(resp, dto.MartinReportListResponse{
			Name: report.Name,
		})
	}

	writeJSON(w, http.StatusOK, resp)
}

// toProductResponse keeps database model fields from leaking directly into API responses.
func toProductResponse(product models.Products) dto.ProductResponse {
	return dto.ProductResponse{
		Code:        product.Code,
		Eng:         product.Eng,
		Viet:        product.Viet,
		Alternative: product.Alternative,
		Brand:       product.Brand,
	}
}

// writeRepositoryError maps repository-level errors to API-safe HTTP responses.
func writeRepositoryError(w http.ResponseWriter, err error) {
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	writeError(w, http.StatusInternalServerError, "internal server error")
}

// writeError sends a consistent JSON error payload.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// writeJSON sends a JSON response with a single, predictable content type.
func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}


func (h* CatalogHandler) RegisterRoutes(r chi.Router) {
	r.Get("/catalog/products", h.SearchProducts)
	r.Get("/catalog/products/{code}", h.GetProductByCode)
	r.Get("/catalog/images/{code}", h.GetImages)
	r.Get("/catalog/report/martin", h.ListMartinReports)
}