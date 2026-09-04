// Package main wires the HTTP server, database pool, and application routes.
package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/solenfw/calo-hub/internal/db"
	"github.com/solenfw/calo-hub/internal/handler"
	"github.com/solenfw/calo-hub/internal/repository"
)

func main() {
	// wire config, db pool, router, start http server
	pool, err := db.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}

	r := chi.NewRouter()
	repo := repository.NewCatalogRepository(pool)
	handler := handler.NewCatalogHandler(repo)

	// Catalog API routes keep exact lookups and searches separate so response shapes stay predictable.
	r.Get("/catalog/products", handler.SearchProducts)
	r.Get("/catalog/products/{code}", handler.GetProductByCode)
	r.Get("/catalog/images/{code}", handler.GetImages)
	r.Get("/catalog/report/martin", handler.ListMartinReports)
	r.Get("/catalog/report/martin/{name}", handler.GetMartinReportByName)
	r.Get("/catalog/report/martin/all/{report_id}", handler.GetMartinReportProducts)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server failed: %v", err)
	}

	pool.Close()
}
