// Package handler contains the HTTP-facing request handlers for catalog-related
// endpoints.
//
// These handlers translate request parameters into repository calls, normalize
// errors into API-safe JSON payloads, and shape repository records into the DTOs
// expected by the client. This keeps database concerns out of the transport
// layer while leaving the route registration in one place for the application.
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
//
// Keeping this interface in the consumer package lets the handlers remain
// decoupled from the concrete repository implementation and makes unit tests use
// small in-memory fakes instead of a live database.
type CatalogStore interface {
	GetProductByCode(ctx context.Context, code string) (models.Products, error)
	GetProductsByTerms(ctx context.Context, terms []string) ([]models.Products, error)
	GetImagesByCode(ctx context.Context, code string) (models.Images, error)
	GetAllMartinReportLists(ctx context.Context) ([]models.MartinReportList, error)
}

var _ CatalogStore = (*repository.CatalogRepository)(nil)

// CatalogHandler serves catalog product, image, and list-oriented Martin report endpoints.
//
// The caller is expected to provide a repository-level implementation of
// CatalogStore, which keeps the handler focused on request parsing and response
// shaping rather than database details.
type CatalogHandler struct {
	store CatalogStore
}

// NewCatalogHandler wires a catalog store into its HTTP handler.
func NewCatalogHandler(store CatalogStore) *CatalogHandler {
	return &CatalogHandler{
		store: store,
	}
}

// GetProductByCode resolves one catalog product by code and returns the public
// DTO response.
//
// The current implementation validates the route parameter before contacting the
// repository so malformed codes do not trigger a database lookup and so the
// transport layer keeps the input contract explicit.
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

// SearchProducts returns any catalog products whose text matches all search
// terms. The handler accepts either q or search as the user-facing query key and
// converts the resulting database model rows into the DTO structure used by the
// client.
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

// GetImages fetches all stored image URLs for a product and returns them in a
// stable response structure. This handler intentionally treats missing image
// data as an empty slice rather than failing the whole request, keeping the API
// resilient when a product has no image metadata yet.
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

// ListMartinReports returns the available Martin report names as a lightweight
// list payload. This endpoint is intentionally simple and does not expose the
// internal report IDs to the client, only the metadata needed to navigate to the
// details view or PDF route.
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

// toProductResponse keeps database model fields from leaking directly into the
// public API response and preserves only the subset the client actually needs.
func toProductResponse(product models.Products) dto.ProductResponse {
	return dto.ProductResponse{
		Code:        product.Code,
		Eng:         product.Eng,
		Viet:        product.Viet,
		Alternative: product.Alternative,
		Brand:       product.Brand,
	}
}

// writeRepositoryError translates repository-layer errors into the JSON error
// responses the API exposes. This centralizes the not-found mapping and keeps
// handler methods focused on request flow rather than storage semantics.
func writeRepositoryError(w http.ResponseWriter, err error) {
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	writeError(w, http.StatusInternalServerError, "internal server error")
}

// writeError sends a consistent JSON error payload so callers receive the same
// response shape regardless of which handler rejected the request.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// writeJSON writes a JSON payload with the single expected content type for
// these handlers, ensuring the response is consistent across endpoints.
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