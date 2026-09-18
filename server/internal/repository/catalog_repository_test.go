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

func seedRepositoryFixtures(t *testing.T) {
	t.Helper()
	pool := repositoryTestPool(t)
	ctx := context.Background()

	_, err := pool.Exec(ctx, `
		TRUNCATE images, products
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
