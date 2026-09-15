// Package repository contains database access for application domain data.
package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/solenfw/calo-hub/internal/models"
)

// ErrNotFound marks a missing catalog row without exposing database-driver errors.
var ErrNotFound = errors.New("catalog item not found")

// CatalogRepository reads catalog products, images, and Martin report data.
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

// GetMartinReportListByName fetches one Martin report list by name.
func (r *CatalogRepository) GetMartinReportListByName(ctx context.Context, reportName string) (models.MartinReportList, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name FROM martin_report_list WHERE name = $1`,
		reportName,
	)
	if err != nil {
		return models.MartinReportList{}, err
	}
	defer rows.Close()

	report, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.MartinReportList])
	if err != nil {
		return models.MartinReportList{}, mapRepositoryError(err)
	}

	return report, nil
}

// GetAllMartinReportLists fetches all Martin report lists.
func (r *CatalogRepository) GetAllMartinReportLists(ctx context.Context) ([]models.MartinReportList, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name FROM martin_report_list`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[models.MartinReportList])
}

// mapRepositoryError normalizes storage errors for callers outside this package.
func mapRepositoryError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	return err
}

// GetMartinReportProductsByReportID fetches frozen report product snapshots in display order.
// It returns ErrNotFound when the report does not exist and an empty slice when it exists without products.
func (r *CatalogRepository) GetMartinReportProductsByReportID(ctx context.Context, reportID int) ([]models.MartinReportProduct, error) {
	query := `
		SELECT 
			mrp.report_id,
			mrp.row_no, 
			mrp.code, 
			mrp.eng,
			mrp.image,
			mrp.quantity 
		FROM martin_report_products mrp
		WHERE mrp.report_id = $1
		ORDER BY mrp.row_no ASC;
	`

	rows, err := r.pool.Query(ctx, query, reportID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.MartinReportProduct])
	if err != nil {
		return nil, err
	}
	if len(products) > 0 {
		return products, nil
	}

	var reportExists bool
	if err := r.pool.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM martin_report_list WHERE id = $1)`,
		reportID,
	).Scan(&reportExists); err != nil {
		return nil, err
	}
	if !reportExists {
		return nil, ErrNotFound
	}

	return products, nil
}
