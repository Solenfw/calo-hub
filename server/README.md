# Calo Hub Go service

This service exposes the catalog and Martin report HTTP API backed by PostgreSQL.

> **Actively under development — mid-migration from a prior Python/FastAPI backend. Some packages are implemented, others are placeholder stubs — see below.**

## Package layout

| Path | Current state |
| --- | --- |
| `cmd/app` | **Partially implemented:** starts the HTTP server, connects with `DATABASE_URL`, wires the catalog repository/handlers/router, and listens on `:8080`; configuration is currently hardcoded/read directly from the environment. |
| `internal/config` | **Stub/placeholder:** `Config` has no fields and `Load` returns `nil, nil`. |
| `internal/db` | **Implemented:** creates and pings a `pgxpool.Pool` from a supplied DSN. |
| `internal/dto` | **Implemented for the catalog API:** defines product, image, Martin report, and report-product response shapes. |
| `internal/handler` | **Partially implemented:** catalog handlers perform validation, repository calls, DTO mapping, and JSON error handling; `convert_handler.go`, `report_handler.go`, and `middleware.go` contain comments only. |
| `internal/models` | **Partially implemented:** catalog database models are defined; `User` is an empty placeholder type. |
| `internal/repository` | **Implemented for the current catalog schema:** product, image, Martin report list, and report-product queries are backed by `pgx`; missing rows map to `ErrNotFound`. |
| `internal/router` | **Implemented for the catalog API only:** registers the six catalog GET routes. |
| `internal/service` | **Stub/placeholder:** `external_api.go` contains only a comment; no external API service exists. |
| `internal/service/pdf` | **Stub/placeholder:** `pdf_drawer.go` contains only a comment; PDF/report generation is not implemented. |

## HTTP API currently registered

The router in `internal/router/router.go` registers:

```text
GET /catalog/products?q=...
GET /catalog/products/{code}
GET /catalog/images/{code}
GET /catalog/report/martin
GET /catalog/report/martin/{name}
GET /catalog/report/martin/all/{report_id}
```

The handlers return JSON. Product searches split the query into whitespace-separated terms, product codes are validated, and repository `ErrNotFound` values become a `404` response with `{"error":"not found"}`. No conversion, PDF, authentication, or general report-generation routes are registered.

## Running locally

From this directory, with PostgreSQL available and the schema applied:

```bash
export DATABASE_URL='postgres://user:password@localhost:5432/database?sslmode=disable'
go run cmd/app/main.go
```

`DATABASE_URL` is the only environment variable read directly by `cmd/app/main.go`. The server listens on `http://localhost:8080`.

The repository root contains `.env.example`, which also lists these Goose CLI variables:

```text
GOOSE_DRIVER
GOOSE_USER
GOOSE_PASSWORD
GOOSE_DB
GOOSE_DBSTRING
GOOSE_MIGRATION_DIR
```

The current Go application does not read those variables. The same example also contains `SECRET_KEY` and `FRONTEND_ORIGINS`; no current Go code consumes them. The sample file does not define `DATABASE_URL`, so add it separately when running the server.

## Migrations

The current schema is in `migrations/20260907221330_init_schema.sql`. It creates `products`, `images`, `martin_report_list`, and `martin_report_products`, with the foreign keys and cascade rules defined in that file.

The `.env.example` file is set up for Goose's environment-driven form: `GOOSE_DRIVER`, `GOOSE_DBSTRING`, and `GOOSE_MIGRATION_DIR` provide the driver, database connection string, and migration directory. From `server/`, the explicit equivalent is:

```bash
goose -dir migrations postgres "$GOOSE_DBSTRING" up
goose -dir migrations postgres "$GOOSE_DBSTRING" down
```

The root Makefile exposes `make migrate-up` and `make migrate-down`, but those recipes invoke only `goose up` and `goose down` from the repository root. They do not pass `-dir`, driver, or database-string flags. With the example `GOOSE_MIGRATION_DIR=./migrations`, that relative path points to a root-level `migrations/` directory rather than this `server/migrations/` directory; use the explicit commands above or correct the environment/Makefile invocation before relying on the root targets.

## Testing

Run all Go tests from `server/`:

```bash
go test ./...
```

The repository currently contains three test tiers:

- `internal/handler/catalog_handler_test.go`: table-driven handler unit tests using fake catalog stores and `net/http/httptest`, covering validation, DTO responses, empty results, route parameters, and repository error mapping.
- `internal/router/router_test.go`: `httptest` wiring tests for route dispatch, path parameters, wrong methods, and unknown paths.
- `internal/repository/catalog_repository_test.go`: PostgreSQL integration tests using testcontainers-go. The tests apply the Goose migration, seed fixtures, exercise catalog/report queries, and verify cascade behavior. They require a usable Docker/testcontainers runtime.

## Root Makefile targets

These targets are defined in the repository root:

| Target | Description |
| --- | --- |
| `make compose` | Starts the development Docker Compose services. |
| `make down` | Stops the development Docker Compose services. |
| `make migrate-up` | Runs `goose up` from the repository root using the exported environment. |
| `make migrate-down` | Runs `goose down` from the repository root using the exported environment. |
| `make run` | Runs `go run server/cmd/app/main.go`. |
| `make test` | Runs `go test ./...` from `server/`. |
| `make vet` | Runs `go vet ./...` from `server/`. |
| `make build` | Runs `go build ./...` from `server/`. |
| `make fmt` | Lists files reported by `gofmt -l` and `goimports -l` from `server/`; it does not rewrite files. |

