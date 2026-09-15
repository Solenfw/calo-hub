// Package router defines the HTTP route table for the application.
package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/solenfw/calo-hub/internal/handler"
	"github.com/go-chi/cors"
)

// New builds the application router for the catalog handlers.
func New(h *handler.CatalogHandler) chi.Router {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Get("/catalog/products", h.SearchProducts)
	r.Get("/catalog/products/{code}", h.GetProductByCode)
	r.Get("/catalog/images/{code}", h.GetImages)
	r.Get("/catalog/report/martin", h.ListMartinReports)
	r.Get("/catalog/report/martin/{name}", h.GetMartinReportByName)
	r.Get("/catalog/report/martin/all/{report_id}", h.GetMartinReportProducts)
	return r
}
