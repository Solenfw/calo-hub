// Package main wires the HTTP server, database pool, and application routes.
//
// The application bootstrap lives here so the repository, route table, and
// request logging can be assembled in one place before the process starts
// serving HTTP traffic. The code keeps startup concerns separate from the
// domain logic in the handler and repository layers.
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

// main boots the application by creating the database pool, wiring the
// repository and handlers, and starting the HTTP server with the router and request logger.
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

// loggingMiddleware records the HTTP method, path, and request duration so each
// route invocation can be observed in the process logs without affecting the
// handler logic itself.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}