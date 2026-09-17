// Package router defines the HTTP route table for the application.
package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)


type routable interface {
	RegisterRoutes(r chi.Router)
}

func New(handlers ...routable) chi.Router {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	for _, h := range handlers {
		h.RegisterRoutes(r)
	}

	return r
}
