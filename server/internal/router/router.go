// Package router defines the application route table and the helper used to wire
// all HTTP handlers into a single chi router.
//
// The router is intentionally thin: it does not contain business logic itself,
// only the registration step that connects each handler to the correct URL
// pattern and shared middleware stack.
package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)


// routable is the minimal interface an HTTP component must satisfy to be
// attached to the application router. This keeps the router generic and allows
// catalog and report handlers to register their routes without knowing about the
// concrete implementation of the other components.
type routable interface {
	RegisterRoutes(r chi.Router)
}

// New builds a chi router and registers all supplied handlers with the shared
// middleware configuration. It centralizes route assembly so individual
// handlers remain focused on their endpoint logic.
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
