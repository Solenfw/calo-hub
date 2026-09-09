// Package router defines the HTTP route table for the application.
package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/solenfw/calo-hub/internal/handler"
)

// New builds the application router for the catalog handlers.
func New(h *handler.CatalogHandler) chi.Router {
	r := chi.NewRouter()
	r.Get("/catalog/products", h.SearchProducts)
	r.Get("/catalog/products/{code}", h.GetProductByCode)
	r.Get("/catalog/images/{code}", h.GetImages)
	r.Get("/catalog/report/martin", h.ListMartinReports)
	r.Get("/catalog/report/martin/{name}", h.GetMartinReportByName)
	r.Get("/catalog/report/martin/all/{report_id}", h.GetMartinReportProducts)
	return r
}
