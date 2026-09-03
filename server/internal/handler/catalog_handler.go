package handler

// endpoints for /catalog
import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/solenfw/calo-hub/internal/repository"
	"github.com/solenfw/calo-hub/internal/dto"
)


type CatalogHandler struct {
	repo *repository.CatalogRepository
}

func NewCatalogHandler(repo *repository.CatalogRepository) *CatalogHandler {
	return &CatalogHandler{
		repo: repo,
	}
}

func (h *CatalogHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	search_term := chi.URLParam(r, "code")
	pattern, err := regexp.Compile(`^\d{2}-\d{3}-\d{2}-\d{2}$|^[A-Z]{2}(?:-[A-Z])?\d{3,6}[A-Z]{0,2}$`) 
	
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	
	if pattern.MatchString(search_term) {
		product, err := h.repo.GetProductByCode(r.Context(), search_term)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		resp := dto.ProductResponse {
			Code: product.Code,
			Eng: product.Eng,
			Viet: product.Viet,
			Alternative: product.Alternative,
			Brand: product.Brand,
		}
		json.NewEncoder(w).Encode(resp)
	}

	terms := regexp.MustCompile(`\s+`).Split(search_term, -1)
	products, err := h.repo.GetProductsByTerms(r.Context(), terms)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var resp []dto.ProductResponse
	for _, product := range products {
		resp = append(resp, dto.ProductResponse{
			Code: product.Code,
			Eng: product.Eng,
			Viet: product.Viet,
			Alternative: product.Alternative,
			Brand: product.Brand,
		})
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *CatalogHandler) GetImages(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	images, err := h.repo.GetImagesByCode(r.Context(), code)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	resp := dto.ImageResponse {
		Code: images.Code,
		Images: images.Images,
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *CatalogHandler) GetMartinReportList(w http.ResponseWriter, r *http.Request) {
	reports, err := h.repo.GetMartinReportList(r.Context())
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var resp []dto.MartinReportListResponse
	for _, report := range reports {
		resp = append(resp, dto.MartinReportListResponse{
			Name: report.Name,
		})
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *CatalogHandler) GetMartinReportProducts(w http.ResponseWriter, r *http.Request) {
	reportID := chi.URLParam(r, "report_id")

	reportIDInt, err := strconv.Atoi(reportID)
	if err != nil {
		http.Error(w, "invalid report ID", http.StatusBadRequest)
		return
	}

	products, err := h.repo.GetMartinReportProductsByID(r.Context(), reportIDInt)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var resp []dto.MartinReportProductResponse
	for _, product := range products {
		resp = append(resp, dto.MartinReportProductResponse{
			RowNo: product.RowNo,
			Code: product.Code,
			Eng: product.Eng,
			Image: product.Image,
			Quantity: product.Quantity,
		})
	}
	json.NewEncoder(w).Encode(resp)
}
