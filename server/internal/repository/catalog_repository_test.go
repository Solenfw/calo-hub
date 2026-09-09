package repository

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/solenfw/calo-hub/internal/models"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

const postgresTestImage = "postgres:17-alpine"

type repositoryTestEnvironment struct {
	once      sync.Once
	container *postgres.PostgresContainer
	sqlDB     *sql.DB
	pool      *pgxpool.Pool
	err       error
}

var repositoryTestEnv repositoryTestEnvironment

type repositoryFixtures struct {
	reportWithProductsID int
	emptyReportID        int
}

func TestMain(m *testing.M) {
	exitCode := m.Run()

	if repositoryTestEnv.pool != nil {
		repositoryTestEnv.pool.Close()
	}
	if repositoryTestEnv.sqlDB != nil {
		_ = repositoryTestEnv.sqlDB.Close()
	}
	if repositoryTestEnv.container != nil {
		_ = testcontainers.TerminateContainer(repositoryTestEnv.container)
	}

	os.Exit(exitCode)
}

func repositoryTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	testcontainers.SkipIfProviderIsNotHealthy(t)

	repositoryTestEnv.once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		repositoryTestEnv.container, repositoryTestEnv.err = postgres.Run(
			ctx,
			postgresTestImage,
			postgres.WithDatabase("calo_hub_test"),
			postgres.WithUsername("postgres"),
			postgres.WithPassword("postgres"),
			postgres.BasicWaitStrategies(),
		)
		if repositoryTestEnv.err != nil {
			return
		}

		dsn, err := repositoryTestEnv.container.ConnectionString(ctx, "sslmode=disable")
		if err != nil {
			repositoryTestEnv.err = err
			return
		}

		repositoryTestEnv.sqlDB, err = sql.Open("pgx", dsn)
		if err != nil {
			repositoryTestEnv.err = err
			return
		}
		if err := repositoryTestEnv.sqlDB.PingContext(ctx); err != nil {
			repositoryTestEnv.err = err
			return
		}

		_, sourcePath, _, ok := runtime.Caller(0)
		if !ok {
			repositoryTestEnv.err = errors.New("could not locate repository test source")
			return
		}
		migrationDir := filepath.Join(filepath.Dir(sourcePath), "../../migrations")
		provider, err := goose.NewProvider(
			goose.DialectPostgres,
			repositoryTestEnv.sqlDB,
			os.DirFS(migrationDir),
		)
		if err != nil {
			repositoryTestEnv.err = err
			return
		}
		if _, err := provider.Up(ctx); err != nil {
			repositoryTestEnv.err = err
			return
		}

		repositoryTestEnv.pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			repositoryTestEnv.err = err
			return
		}
		repositoryTestEnv.err = repositoryTestEnv.pool.Ping(ctx)
	})

	if repositoryTestEnv.err != nil {
		t.Fatalf("set up PostgreSQL test environment: %v", repositoryTestEnv.err)
	}
	return repositoryTestEnv.pool
}

func seedRepositoryFixtures(t *testing.T) repositoryFixtures {
	t.Helper()
	pool := repositoryTestPool(t)
	ctx := context.Background()

	_, err := pool.Exec(ctx, `
		TRUNCATE martin_report_products, martin_report_list, images, products
		RESTART IDENTITY CASCADE;
	`)
	if err != nil {
		t.Fatalf("truncate fixtures: %v", err)
	}

	_, err = pool.Exec(ctx, `
		INSERT INTO products (code, eng, viet, alternative, brand) VALUES
			('P-001', 'Live Alpha Clamp', 'Kẹp Alpha', 'ALT-001', 'Martin'),
			('P-002', 'Live Valve Holder', 'Giá đỡ van', 'ALT-002', 'Aesculap'),
			('P-003', 'Beta Retractor', 'Viet Beta', 'ALT-003', 'Martin');

		INSERT INTO images (code, images) VALUES
			('P-001', ARRAY['live-alpha-1.png', 'live-alpha-2.png']),
			('P-003', ARRAY[]::text[]);
	`)
	if err != nil {
		t.Fatalf("seed products and images: %v", err)
	}

	var fixtures repositoryFixtures
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

	_, err = pool.Exec(ctx, `
		INSERT INTO martin_report_products (report_id, row_no, code, eng, image, quantity) VALUES
			($1, 1, 'P-001', 'Frozen Alpha Clamp', 'frozen-alpha.png', 3),
			($1, 2, 'P-002', 'Frozen Valve Holder', NULL, 1);
	`, fixtures.reportWithProductsID)
	if err != nil {
		t.Fatalf("seed report products: %v", err)
	}

	return fixtures
}

func expectedProduct(code, eng, viet, alternative, brand string) models.Products {
	return models.Products{
		Code:        code,
		Eng:         eng,
		Viet:        viet,
		Alternative: alternative,
		Brand:       brand,
	}
}

func productCodes(products []models.Products) []string {
	codes := make([]string, 0, len(products))
	for _, product := range products {
		codes = append(codes, product.Code)
	}
	return codes
}

func TestCatalogRepositoryGetProductByCode(t *testing.T) {
	seedRepositoryFixtures(t)
	repo := NewCatalogRepository(repositoryTestPool(t))

	tests := []struct {
		name      string
		code      string
		want      models.Products
		wantError error
	}{
		{
			name: "existing code",
			code: "P-001",
			want: expectedProduct("P-001", "Live Alpha Clamp", "Kẹp Alpha", "ALT-001", "Martin"),
		},
		{
			name:      "nonexistent code",
			code:      "P-404",
			wantError: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetProductByCode(context.Background(), tt.code)
			if tt.wantError != nil {
				if !errors.Is(err, tt.wantError) {
					t.Fatalf("error = %v, want %v", err, tt.wantError)
				}
				if errors.Is(err, pgx.ErrNoRows) {
					t.Fatalf("error leaked pgx.ErrNoRows: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("GetProductByCode() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("product = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestCatalogRepositoryGetProductsByTerms(t *testing.T) {
	seedRepositoryFixtures(t)
	repo := NewCatalogRepository(repositoryTestPool(t))

	tests := []struct {
		name      string
		terms     []string
		wantCodes []string
	}{
		{name: "English match", terms: []string{"alpha"}, wantCodes: []string{"P-001"}},
		{name: "Vietnamese match", terms: []string{"đỡ"}, wantCodes: []string{"P-002"}},
		{name: "all terms across English and Vietnamese", terms: []string{"alpha", "kẹp"}, wantCodes: []string{"P-001"}},
		{name: "does not match any one term", terms: []string{"alpha", "valve"}, wantCodes: []string{}},
		{name: "case insensitive", terms: []string{"ALPHA"}, wantCodes: []string{"P-001"}},
		{name: "no match", terms: []string{"not-present"}, wantCodes: []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetProductsByTerms(context.Background(), tt.terms)
			if err != nil {
				t.Fatalf("GetProductsByTerms() error = %v", err)
			}
			if got == nil {
				t.Fatal("GetProductsByTerms() returned nil slice")
			}
			if !reflect.DeepEqual(productCodes(got), tt.wantCodes) {
				t.Fatalf("product codes = %#v, want %#v", productCodes(got), tt.wantCodes)
			}
		})
	}
}

func TestCatalogRepositoryGetImagesByCode(t *testing.T) {
	seedRepositoryFixtures(t)
	repo := NewCatalogRepository(repositoryTestPool(t))

	tests := []struct {
		name       string
		code       string
		wantImages []string
		wantError  error
	}{
		{name: "populated array", code: "P-001", wantImages: []string{"live-alpha-1.png", "live-alpha-2.png"}},
		{name: "empty array", code: "P-003", wantImages: []string{}},
		{name: "nonexistent code", code: "P-002", wantError: ErrNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetImagesByCode(context.Background(), tt.code)
			if tt.wantError != nil {
				if !errors.Is(err, tt.wantError) {
					t.Fatalf("error = %v, want %v", err, tt.wantError)
				}
				return
			}
			if err != nil {
				t.Fatalf("GetImagesByCode() error = %v", err)
			}
			if got.Code != tt.code {
				t.Fatalf("code = %q, want %q", got.Code, tt.code)
			}
			if got.Images == nil {
				t.Fatal("GetImagesByCode() returned nil image slice")
			}
			if !reflect.DeepEqual(got.Images, tt.wantImages) {
				t.Fatalf("images = %#v, want %#v", got.Images, tt.wantImages)
			}
		})
	}
}

func TestCatalogRepositoryMartinReportLists(t *testing.T) {
	fixtures := seedRepositoryFixtures(t)
	repo := NewCatalogRepository(repositoryTestPool(t))

	t.Run("all lists", func(t *testing.T) {
		got, err := repo.GetAllMartinReportLists(context.Background())
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

	t.Run("nonexistent name", func(t *testing.T) {
		_, err := repo.GetMartinReportListByName(context.Background(), "Missing Report")
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("error = %v, want ErrNotFound", err)
		}
	})

	t.Run("zero lists", func(t *testing.T) {
		pool := repositoryTestPool(t)
		if _, err := pool.Exec(context.Background(), `DELETE FROM martin_report_list`); err != nil {
			t.Fatalf("delete report lists: %v", err)
		}
		got, err := repo.GetAllMartinReportLists(context.Background())
		if err != nil {
			t.Fatalf("GetAllMartinReportLists() error = %v", err)
		}
		if got == nil || len(got) != 0 {
			t.Fatalf("reports = %#v, want non-nil empty slice", got)
		}
	})
}

func TestCatalogRepositoryGetMartinReportProductsByReportID(t *testing.T) {
	fixtures := seedRepositoryFixtures(t)
	repo := NewCatalogRepository(repositoryTestPool(t))
	snapshotImage := "frozen-alpha.png"

	t.Run("uses frozen snapshot fields", func(t *testing.T) {
		got, err := repo.GetMartinReportProductsByReportID(context.Background(), fixtures.reportWithProductsID)
		if err != nil {
			t.Fatalf("GetMartinReportProductsByReportID() error = %v", err)
		}
		want := []models.MartinReportProduct{
			{ReportID: fixtures.reportWithProductsID, RowNo: 1, Code: "P-001", Eng: "Frozen Alpha Clamp", Image: &snapshotImage, Quantity: 3},
			{ReportID: fixtures.reportWithProductsID, RowNo: 2, Code: "P-002", Eng: "Frozen Valve Holder", Image: nil, Quantity: 1},
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("report products = %#v, want %#v", got, want)
		}
	})

	t.Run("existing report with zero products", func(t *testing.T) {
		got, err := repo.GetMartinReportProductsByReportID(context.Background(), fixtures.emptyReportID)
		if err != nil {
			t.Fatalf("GetMartinReportProductsByReportID() error = %v", err)
		}
		if got == nil || len(got) != 0 {
			t.Fatalf("report products = %#v, want non-nil empty slice", got)
		}
	})

	t.Run("nonexistent report", func(t *testing.T) {
		_, err := repo.GetMartinReportProductsByReportID(context.Background(), 404)
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("error = %v, want ErrNotFound", err)
		}
	})
}

func TestCatalogRepositoryCascadeDeletesImagesWhenProductDeleted(t *testing.T) {
	seedRepositoryFixtures(t)
	pool := repositoryTestPool(t)

	if _, err := pool.Exec(context.Background(), `DELETE FROM products WHERE code = 'P-001'`); err != nil {
		t.Fatalf("delete product: %v", err)
	}

	var count int
	if err := pool.QueryRow(context.Background(), `SELECT COUNT(*) FROM images WHERE code = 'P-001'`).Scan(&count); err != nil {
		t.Fatalf("count images: %v", err)
	}
	if count != 0 {
		t.Fatalf("images rows = %d, want 0", count)
	}
}

func TestCatalogRepositoryCascadeDeletesReportProductsWhenReportDeleted(t *testing.T) {
	fixtures := seedRepositoryFixtures(t)
	pool := repositoryTestPool(t)

	if _, err := pool.Exec(context.Background(), `DELETE FROM martin_report_list WHERE id = $1`, fixtures.reportWithProductsID); err != nil {
		t.Fatalf("delete report list: %v", err)
	}

	var count int
	if err := pool.QueryRow(context.Background(), `SELECT COUNT(*) FROM martin_report_products WHERE report_id = $1`, fixtures.reportWithProductsID).Scan(&count); err != nil {
		t.Fatalf("count report products: %v", err)
	}
	if count != 0 {
		t.Fatalf("report product rows = %d, want 0", count)
	}
}

func TestCatalogRepositoryCascadeDeletesReportProductsWhenProductDeleted(t *testing.T) {
	fixtures := seedRepositoryFixtures(t)
	pool := repositoryTestPool(t)

	if _, err := pool.Exec(context.Background(), `DELETE FROM products WHERE code = 'P-001'`); err != nil {
		t.Fatalf("delete product: %v", err)
	}

	var count int
	if err := pool.QueryRow(context.Background(), `
		SELECT COUNT(*)
		FROM martin_report_products
		WHERE report_id = $1 AND code = 'P-001'
	`, fixtures.reportWithProductsID).Scan(&count); err != nil {
		t.Fatalf("count report product rows: %v", err)
	}
	if count != 0 {
		t.Fatalf("deleted product report rows = %d, want 0", count)
	}
}
