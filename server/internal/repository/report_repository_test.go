package repository

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/solenfw/calo-hub/internal/models"
)

type reportRepositoryFixtures struct {
	reportWithProductsID int
	emptyReportID        int
}

func seedReportRepositoryFixtures(t *testing.T) reportRepositoryFixtures {
	t.Helper()
	pool := repositoryTestPool(t)
	ctx := context.Background()

	if _, err := pool.Exec(ctx, `
		TRUNCATE martin_report_products, martin_report_list, images, products
		RESTART IDENTITY CASCADE;
	`); err != nil {
		t.Fatalf("truncate fixtures: %v", err)
	}

	if _, err := pool.Exec(ctx, `
		INSERT INTO products (code, eng, viet, alternative, brand) VALUES
			('P-001', 'Live Alpha Clamp', 'Kẹp Alpha', 'ALT-001', 'Martin'),
			('P-002', 'Live Valve Holder', 'Giá đỡ van', 'ALT-002', 'Aesculap'),
			('P-003', 'Beta Retractor', 'Viet Beta', 'ALT-003', 'Martin');
	`); err != nil {
		t.Fatalf("seed products: %v", err)
	}

	var fixtures reportRepositoryFixtures
	if err := pool.QueryRow(ctx,
		`INSERT INTO martin_report_list (name) VALUES ('Report With Products') RETURNING id`,
	).Scan(&fixtures.reportWithProductsID); err != nil {
		t.Fatalf("seed populated report: %v", err)
	}
	if err := pool.QueryRow(ctx,
		`INSERT INTO martin_report_list (name) VALUES ('Empty Report') RETURNING id`,
	).Scan(&fixtures.emptyReportID); err != nil {
		t.Fatalf("seed empty report: %v", err)
	}

	if _, err := pool.Exec(ctx, `
		INSERT INTO martin_report_products (report_id, row_no, code, description, image, quantity) VALUES
			($1, 1, 'P-001', 'Frozen Alpha Clamp', 'frozen-alpha.png', 3),
			($1, 2, 'P-002', 'Frozen Valve Holder', NULL, 1);
	`, fixtures.reportWithProductsID); err != nil {
		t.Fatalf("seed report products: %v", err)
	}

	return fixtures
}

func TestReportRepositoryMartinReportLists(t *testing.T) {
	fixtures := seedReportRepositoryFixtures(t)
	repo := NewReportRepository(repositoryTestPool(t))
	ctx := context.Background()

	t.Run("all lists", func(t *testing.T) {
		got, err := repo.GetAllMartinReportLists(ctx)
		if err != nil {
			t.Fatalf("GetAllMartinReportLists() error = %v", err)
		}
		want := []models.MartinReportList{
			{ID: fixtures.reportWithProductsID, Name: "Report With Products"},
			{ID: fixtures.emptyReportID, Name: "Empty Report"},
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("reports = %#v, want %#v", got, want)
		}
	})

	t.Run("find by name", func(t *testing.T) {
		got, err := repo.GetMartinReportListByName(ctx, "Report With Products")
		if err != nil {
			t.Fatalf("GetMartinReportListByName() error = %v", err)
		}
		if got.ID != fixtures.reportWithProductsID || got.Name != "Report With Products" {
			t.Fatalf("report = %#v, want populated fixture", got)
		}
	})

	t.Run("missing name", func(t *testing.T) {
		_, err := repo.GetMartinReportListByName(ctx, "Missing Report")
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("error = %v, want ErrNotFound", err)
		}
	})

	t.Run("zero lists", func(t *testing.T) {
		if _, err := repositoryTestPool(t).Exec(ctx, `DELETE FROM martin_report_list`); err != nil {
			t.Fatalf("delete report lists: %v", err)
		}
		got, err := repo.GetAllMartinReportLists(ctx)
		if err != nil {
			t.Fatalf("GetAllMartinReportLists() error = %v", err)
		}
		if got == nil || len(got) != 0 {
			t.Fatalf("reports = %#v, want non-nil empty slice", got)
		}
	})
}

func TestReportRepositoryProducts(t *testing.T) {
	fixtures := seedReportRepositoryFixtures(t)
	repo := NewReportRepository(repositoryTestPool(t))
	ctx := context.Background()
	snapshotImage := "frozen-alpha.png"

	t.Run("reads frozen snapshot fields in row order", func(t *testing.T) {
		got, err := repo.GetMartinReportProductsByReportID(ctx, fixtures.reportWithProductsID)
		if err != nil {
			t.Fatalf("GetMartinReportProductsByReportID() error = %v", err)
		}
		want := []models.MartinReportProduct{
			{ReportID: fixtures.reportWithProductsID, RowNo: 1, Code: "P-001", Description: "Frozen Alpha Clamp", Image: &snapshotImage, Quantity: 3},
			{ReportID: fixtures.reportWithProductsID, RowNo: 2, Code: "P-002", Description: "Frozen Valve Holder", Image: nil, Quantity: 1},
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("report products = %#v, want %#v", got, want)
		}
	})

	t.Run("existing report with zero products", func(t *testing.T) {
		got, err := repo.GetMartinReportProductsByReportID(ctx, fixtures.emptyReportID)
		if err != nil {
			t.Fatalf("GetMartinReportProductsByReportID() error = %v", err)
		}
		if got == nil || len(got) != 0 {
			t.Fatalf("report products = %#v, want non-nil empty slice", got)
		}
	})

	t.Run("missing report", func(t *testing.T) {
		_, err := repo.GetMartinReportProductsByReportID(ctx, 404)
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("error = %v, want ErrNotFound", err)
		}
	})
}

func TestReportRepositoryCRUD(t *testing.T) {
	fixtures := seedReportRepositoryFixtures(t)
	repo := NewReportRepository(repositoryTestPool(t))
	ctx := context.Background()

	t.Run("create report list", func(t *testing.T) {
		report, err := repo.CreateMartinReportList(ctx, "New Surgical Report")
		if err != nil {
			t.Fatalf("CreateMartinReportList() error = %v", err)
		}
		if report.ID <= 0 || report.Name != "New Surgical Report" {
			t.Fatalf("report = %#v, want generated ID and requested name", report)
		}
	})

	t.Run("replace products", func(t *testing.T) {
		products := []models.MartinReportProduct{
			{ReportID: fixtures.reportWithProductsID, RowNo: 1, Code: "P-003", Description: "Beta Retractor", Quantity: 5},
			{ReportID: fixtures.reportWithProductsID, RowNo: 2, Code: "P-001", Description: "Live Alpha Clamp", Quantity: 1},
		}
		if err := repo.ReplaceMartinReportProducts(ctx, fixtures.reportWithProductsID, products); err != nil {
			t.Fatalf("ReplaceMartinReportProducts() error = %v", err)
		}
		got, err := repo.GetMartinReportProductsByReportID(ctx, fixtures.reportWithProductsID)
		if err != nil {
			t.Fatalf("GetMartinReportProductsByReportID() error = %v", err)
		}
		if !reflect.DeepEqual(got, products) {
			t.Fatalf("products = %#v, want %#v", got, products)
		}

		if err := repo.ReplaceMartinReportProducts(ctx, fixtures.reportWithProductsID, nil); err != nil {
			t.Fatalf("ReplaceMartinReportProducts(empty) error = %v", err)
		}
		got, err = repo.GetMartinReportProductsByReportID(ctx, fixtures.reportWithProductsID)
		if err != nil || len(got) != 0 {
			t.Fatalf("products after empty replace = %#v, error = %v", got, err)
		}
	})

	t.Run("delete report cascades products", func(t *testing.T) {
		if err := repo.DeleteMartinReportList(ctx, fixtures.reportWithProductsID); err != nil {
			t.Fatalf("DeleteMartinReportList() error = %v", err)
		}
		var reportCount, productCount int
		pool := repositoryTestPool(t)
		if err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM martin_report_list WHERE id = $1`, fixtures.reportWithProductsID).Scan(&reportCount); err != nil {
			t.Fatalf("count report rows: %v", err)
		}
		if err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM martin_report_products WHERE report_id = $1`, fixtures.reportWithProductsID).Scan(&productCount); err != nil {
			t.Fatalf("count product rows: %v", err)
		}
		if reportCount != 0 || productCount != 0 {
			t.Fatalf("report count = %d, product count = %d, want both zero", reportCount, productCount)
		}
	})

	t.Run("missing report errors", func(t *testing.T) {
		if err := repo.DeleteMartinReportList(ctx, 99999); !errors.Is(err, ErrNotFound) {
			t.Fatalf("DeleteMartinReportList() error = %v, want ErrNotFound", err)
		}
		if err := repo.ReplaceMartinReportProducts(ctx, 99999, nil); !errors.Is(err, ErrNotFound) {
			t.Fatalf("ReplaceMartinReportProducts() error = %v, want ErrNotFound", err)
		}
	})
}
