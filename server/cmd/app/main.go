package main

import (
	"context"
	"log"
	"os"
	"net/http"

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

	r.Get("/catalog/products/{code}", handler.GetProducts)
	r.Get("/catalog/images/{code}", handler.GetImages)
	r.Get("/catalog/report/martin/{request_id}", handler.GetMartinReportList)
	r.Get("/catalog/report/martin/all/{request_id}", handler.GetMartinReportProducts)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server failed: %v", err)
	}

	pool.Close()
}
