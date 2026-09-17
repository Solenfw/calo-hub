// Package db owns the database connection lifecycle for the application.
//
// The repository layer expects a *pgxpool.Pool and the bootstrap code wires it
// here so the rest of the server can work with a shared, pooled connection
// source rather than opening ad-hoc database sessions.
package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)


// Connect creates and validates a pgx connection pool from the supplied DSN.
// The pool is used across repositories so queries can share a connection pool
// while still keeping database access logic separated from HTTP concerns.
func Connect(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}
