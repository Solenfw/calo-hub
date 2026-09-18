// Package repository contains database access for application domain data.
package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/solenfw/calo-hub/internal/models"
)

// ReportRepository reads and updates Martin report lists and their instrument snapshots.
type ReportRepository struct {
	pool *pgxpool.Pool
}

// NewReportRepository creates a report repository backed by a pgx connection pool.
func NewReportRepository(pool *pgxpool.Pool) *ReportRepository {
	return &ReportRepository{pool: pool}
}

// GetMartinReportListByName fetches one Martin report list by name.
func (r *ReportRepository) GetMartinReportListByName(ctx context.Context, reportName string) (models.MartinReportList, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name FROM martin_report_list WHERE name = $1`, reportName)
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
func (r *ReportRepository) GetAllMartinReportLists(ctx context.Context) ([]models.MartinReportList, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name FROM martin_report_list`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[models.MartinReportList])
}

// CreateMartinReportList inserts a new report and returns the persisted row.
func (r *ReportRepository) CreateMartinReportList(ctx context.Context, reportName string) (models.MartinReportList, error) {
	var report models.MartinReportList
	if err := r.pool.QueryRow(ctx,
		`INSERT INTO martin_report_list (name) VALUES ($1) RETURNING id, name`, reportName,
	).Scan(&report.ID, &report.Name); err != nil {
		return models.MartinReportList{}, err
	}
	return report, nil
}

// DeleteMartinReportList removes a report and all of its linked product rows.
func (r *ReportRepository) DeleteMartinReportList(ctx context.Context, reportID int) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM martin_report_list WHERE id = $1`, reportID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// GetMartinReportProductsByReportID fetches frozen report product snapshots in display order.
// It returns ErrNotFound when the report does not exist and an empty slice when it has no products.
func (r *ReportRepository) GetMartinReportProductsByReportID(ctx context.Context, reportID int) ([]models.MartinReportProduct, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT report_id, row_no, code, description, image, quantity
		FROM martin_report_products
		WHERE report_id = $1
		ORDER BY row_no ASC`, reportID)
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
	if err := r.pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM martin_report_list WHERE id = $1)`, reportID).Scan(&reportExists); err != nil {
		return nil, err
	}
	if !reportExists {
		return nil, ErrNotFound
	}
	return products, nil
}

// ReplaceMartinReportProducts replaces the full instrument snapshot for a report.
func (r *ReportRepository) ReplaceMartinReportProducts(ctx context.Context, reportID int, products []models.MartinReportProduct) error {
	var reportExists bool
	if err := r.pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM martin_report_list WHERE id = $1)`, reportID).Scan(&reportExists); err != nil {
		return err
	}
	if !reportExists {
		return ErrNotFound
	}

	if len(products) == 0 {
		_, err := r.pool.Exec(ctx, `DELETE FROM martin_report_products WHERE report_id = $1`, reportID)
		return err
	}

	batch := &pgx.Batch{}
	batch.Queue(`DELETE FROM martin_report_products WHERE report_id = $1`, reportID)
	for _, product := range products {
		batch.Queue(
			`INSERT INTO martin_report_products (report_id, row_no, code, description, image, quantity) VALUES ($1, $2, $3, $4, $5, $6)`,
			reportID, product.RowNo, product.Code, product.Description, product.Image, product.Quantity,
		)
	}

	batchResults := r.pool.SendBatch(ctx, batch)
	defer batchResults.Close()
	if _, err := batchResults.Exec(); err != nil {
		return err
	}
	for range products {
		if _, err := batchResults.Exec(); err != nil {
			return err
		}
	}
	return nil
}

var _ interface {
	GetMartinReportListByName(context.Context, string) (models.MartinReportList, error)
	GetAllMartinReportLists(context.Context) ([]models.MartinReportList, error)
	GetMartinReportProductsByReportID(context.Context, int) ([]models.MartinReportProduct, error)
	CreateMartinReportList(context.Context, string) (models.MartinReportList, error)
	DeleteMartinReportList(context.Context, int) error
	ReplaceMartinReportProducts(context.Context, int, []models.MartinReportProduct) error
} = (*ReportRepository)(nil)
