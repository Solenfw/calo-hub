// Package config centralizes application configuration and environment-driven
// settings used by the server bootstrap.
package config

import (
	"fmt"
	"os"
)

// Config stores the runtime settings the application reads from configuration
// sources such as environment variables or deployment files. The project is
// still in the early stages of wiring these values through the app, so the
// struct is intentionally small and expandable.
type Config struct {
	// DB DSN, port, external API keys, etc.
	DatabaseURL string
	Port 			string
}

// Load reads configuration values and returns a populated Config instance.
func Load() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = os.Getenv("SUPABASE_DATABASE_URL")
	}

	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required but not set.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		DatabaseURL: dbURL,
		Port: port,
	}, nil
	
}
