// Package repository contains database access for application domain data.
package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/solenfw/calo-hub/internal/models"
)

// ErrNotFound marks a missing catalog row without exposing database-driver
// errors through the application boundary.
var ErrNotFound = errors.New("catalog item not found")

// CatalogRepository reads catalog products and images from the PostgreSQL database.
//
// The repository is intentionally responsible only for persistence and row-to-
// struct mapping. HTTP code and business logic live in the handler layer, while
// the database specifics remain here.
type CatalogRepository struct {
	pool *pgxpool.Pool
}

// NewCatalogRepository creates a repository backed by a pgx connection pool.
func NewCatalogRepository(pool *pgxpool.Pool) *CatalogRepository {
	return &CatalogRepository{
		pool: pool,
	}
}

// GetProductByCode fetches one product by its exact catalog code.
func (r *CatalogRepository) GetProductByCode(ctx context.Context, code string) (models.Products, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT code, eng, viet, alternative, brand 
		 FROM products 
		 WHERE code = $1`,
		code,
	)
	if err != nil {
		return models.Products{}, err
	}
	defer rows.Close()

	product, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.Products])
	if err != nil {
		return models.Products{}, mapRepositoryError(err)
	}

	return product, nil
}

// GetProductsByTerms fetches products whose English or Vietnamese names match every term.
func (r *CatalogRepository) GetProductsByTerms(ctx context.Context, terms []string) ([]models.Products, error) {
	if len(terms) == 0 {
		return []models.Products{}, nil
	}

	query := `
		SELECT code, eng, viet, alternative, brand 
		FROM products
		WHERE concat_ws(' ', eng, viet) ILIKE ALL (
			ARRAY(SELECT '%' || unnest($1::text[]) || '%')
		)
	`

	rows, err := r.pool.Query(ctx, query, terms)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[models.Products])
}

// GetImagesByCode fetches the stored image list for one product code.
func (r *CatalogRepository) GetImagesByCode(ctx context.Context, code string) (models.Images, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT code, images FROM images WHERE code = $1`,
		code,
	)
	if err != nil {
		return models.Images{}, err
	}
	defer rows.Close()

	images, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.Images])
	if err != nil {
		return models.Images{}, mapRepositoryError(err)
	}

	return images, nil
}

// mapRepositoryError normalizes storage errors for callers outside this package.
func mapRepositoryError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	return err
}
