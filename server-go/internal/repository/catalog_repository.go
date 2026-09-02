package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/solenfw/calo-hub/internal/models"
)



type CatalogRepository struct {
	pool *pgxpool.Pool	
}	

func NewCatalogRepository(pool *pgxpool.Pool) *CatalogRepository {
	return &CatalogRepository{
		pool: pool,
	}
}


func (r *CatalogRepository) GetProductByCode(ctx context.Context, code string) (models.Products, error) {
	// Implement the logic to retrieve a Products by its code from the database
	// using the r.pool connection pool.

	rows, err := r.pool.Query(ctx, 
		`SELECT code, eng, viet, alternative, brand FROM products WHERE code = $1`, 
		code,
	)

	if err != nil {
		return models.Products{}, err
	}
	defer rows.Close()

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.Products])
}

func (r *CatalogRepository) GetProductsByTerms(ctx context.Context, terms []string) ([]models.Products, error) {
	if len(terms) == 0 {
		return []models.Products{}, nil
	}

	query := `
		SELECT code, eng, viet, alternative, brand 
		FROM products
		WHERE concat_ws(' ', eng, viet) ILIKE ALL (
			ARRAY(SELECT '%' || unnest($1::text[]) || '%')
		);
	`

	rows, err := r.pool.Query(ctx, query, terms)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[models.Products])
}

func (r *CatalogRepository) GetImagesByCode(ctx context.Context, code string) (models.Images, error) {
	rows, err := r.pool.Query(ctx, 
		`SELECT code, Images FROM images WHERE code = $1`, 
		code,
	)

	if err != nil {
		return models.Images{}, err
	}
	defer rows.Close()

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.Images])
}


func (r *CatalogRepository) GetMartinReportList(ctx context.Context) ([]models.MartinReportList, error) {
	rows, err := r.pool.Query(ctx, 
		`SELECT id, name FROM martin_report_list`, 
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[models.MartinReportList])
}

func (r *CatalogRepository) GetMartinReportProductsByID(ctx context.Context, reportID int) ([]models.MartinReportProduct, error) {
	query := `
		SELECT 
			mrp.report_id,
			mrp.row_no, 
			mrp.code, 
			COALESCE(p.eng, '') AS eng, 
			img.images[1]       AS image, 
			mrp.quantity 
		FROM martin_report_products mrp
		LEFT JOIN martin_products p ON mrp.code = p.code
		LEFT JOIN images img ON p.code = img.code
		WHERE mrp.report_id = $1
		ORDER BY mrp.row_no ASC;
	`

	rows, err := r.pool.Query(ctx, query, reportID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[models.MartinReportProduct])
}