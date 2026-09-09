// Package main wires the HTTP server, database pool, and application routes.
package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/solenfw/calo-hub/internal/db"
	"github.com/solenfw/calo-hub/internal/handler"
	"github.com/solenfw/calo-hub/internal/repository"
	"github.com/solenfw/calo-hub/internal/router"
)

func main() {
	// wire config, db pool, router, start http server
	pool, err := db.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}

	repo := repository.NewCatalogRepository(pool)
	catalogHandler := handler.NewCatalogHandler(repo)
	r := router.New(catalogHandler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server failed: %v", err)
	}

	pool.Close()
}
