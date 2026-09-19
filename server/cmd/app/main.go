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
	pool, err := db.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}

	catalogRepo := repository.NewCatalogRepository(pool)
	reportRepo := repository.NewReportRepository(pool)

	catalogHandler := handler.NewCatalogHandler(catalogRepo)
	reportHandler := handler.NewReportHandler(reportRepo, catalogRepo)

	r := router.New(catalogHandler, reportHandler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", loggingMiddleware(r)); err != nil {
		log.Fatalf("server failed: %v", err)
	}

	pool.Close()
}
// responseRecorder is a custom wrapper to capture the HTTP status code.
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader intercepts the status code before passing it to the real ResponseWriter.
func (rec *responseRecorder) WriteHeader(statusCode int) {
	rec.statusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

// loggingMiddleware records the status code, method, path, and duration.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the original writer. Default to 200 in case WriteHeader is never called explicitly.
		rec := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK, 
		}

		next.ServeHTTP(rec, r)

		// Now you can see exactly which requests are failing
		log.Printf("[%d] %s %s (%s)", rec.statusCode, r.Method, r.URL.Path, time.Since(start))
	})
}
