// Package main wires the HTTP server, database pool, and application routes.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
	
	"github.com/joho/godotenv"
	"github.com/solenfw/calo-hub/internal/db"
	"github.com/solenfw/calo-hub/internal/handler"
	"github.com/solenfw/calo-hub/internal/repository"
	"github.com/solenfw/calo-hub/internal/router"
)

func main() {
	homeDir, _ := os.UserHomeDir()
	_ = godotenv.Load(homeDir + "/.config/secrets/calohub/.env")

	// wire config, db pool, router, start http server
	pool, err := db.Connect(context.Background(), os.Getenv("GOOSE_DBSTRING"))
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}

	repo := repository.NewCatalogRepository(pool)

	catalogHandler := handler.NewCatalogHandler(repo)
	reportHandler := handler.NewReportHandler(repo)

	r := router.New(catalogHandler, reportHandler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", loggingMiddleware(r)); err != nil {
		log.Fatalf("server failed: %v", err)
	}

	pool.Close()
}


func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}