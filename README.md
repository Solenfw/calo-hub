# Calo Hub

Calo Hub is a medical instrument catalog and Martin report-management tool with a Next.js interface and a Go/PostgreSQL backend.

> **Status: actively under development.** The backend is mid-migration from Python/FastAPI to Go — some endpoints and features are implemented, others are still placeholders. See [Current state](#current-state) before assuming any feature works end-to-end.

## Stack

- **Backend:** Go 1.26.3, [chi](https://github.com/go-chi/chi), [pgx](https://github.com/jackc/pgx), [goose](https://github.com/pressly/goose), and PostgreSQL 17.
- **Frontend:** Next.js 16, React 19, Tailwind CSS 4, Motion, and Lucide React.
- **Testing:** Go `httptest`, table-driven handler/router tests, and PostgreSQL integration tests using testcontainers-go.
- **Infrastructure:** Docker Compose.

The current backend is Go. The repository still contains migration-era references to the former Python/FastAPI application, including comments describing replacements or ports and the older root README content. No Python backend implementation is part of the current source tree.

## Current state

### Implemented

- The Go application starts from `server/cmd/app/main.go`, connects to PostgreSQL using `DATABASE_URL`, constructs the catalog repository and handlers, and listens on port `8080`.
- Catalog routes are registered in `server/internal/router`:

  - `GET /catalog/products?q=...` searches products by all whitespace-separated terms.
  - `GET /catalog/products/{code}` returns one product after code validation.
  - `GET /catalog/images/{code}` returns the product image list.
  - `GET /catalog/report/martin` lists Martin report lists.
  - `GET /catalog/report/martin/{name}` returns one Martin report list.
  - `GET /catalog/report/martin/all/{report_id}` returns the report's product snapshots.

- Catalog handlers map database models to DTOs, return JSON responses, validate route/query parameters, and normalize missing rows to `404 {"error":"not found"}`.
- The repository implements product lookup/search, image lookup, Martin report list lookup/listing, and report-product lookup. Report-product reads use the stored `eng` and `image` snapshot columns rather than joining live product descriptions.
- The Goose migration creates `products`, `images`, `martin_report_list`, and `martin_report_products`, including the current foreign keys and cascade-delete behavior.
- The frontend has these Next.js routes: `/`, `/dashboard`, `/catalog`, `/convert`, `/extraction`, and `/eselector`.
- The catalog UI calls the Go API for product search and product images. Client helpers also exist for Martin report endpoints.

### Not implemented or still placeholder

- `internal/config.Load` is an empty configuration stub that returns `nil, nil`; the entrypoint currently reads `DATABASE_URL` directly and hardcodes port `8080`.
- `internal/service/external_api.go` contains only a placeholder comment; there is no external API integration.
- `internal/service/pdf/pdf_drawer.go` contains only a placeholder comment; PDF/report generation is not implemented.
- `internal/handler/convert_handler.go`, `internal/handler/report_handler.go`, and `internal/handler/middleware.go` contain package declarations and comments only. There are no Go conversion, report-generation, or middleware endpoints wired into the router.
- The Convert, eSelector, Extraction, and most Dashboard UI is static/presentational code with hardcoded example content and controls that are not connected to backend features.
- There is no data-import command in the Go application. CSV files are present under `server/data`, but the migrations do not seed those files automatically.
- `docker-compose.prod.yml` declares `db`, `api`, and `frontend` services, but the repository has only `client/Dockerfile`; there is no root or server Dockerfile, and the production Compose build entries do not specify build contexts. Treat the production Compose file as incomplete until its build setup is supplied.

## Database schema

The current migration is `server/migrations/20260907221330_init_schema.sql`:

- `products` stores catalog codes, English/Vietnamese descriptions, alternative codes, and brands.
- `images` stores one `TEXT[]` image list per product and cascades when the product is deleted.
- `martin_report_list` stores report-list names.
- `martin_report_products` stores ordered report rows, including frozen English and image snapshot values, and cascades when either the report list or referenced product is deleted.

The local Compose database uses PostgreSQL 17 and exposes port `5432`.

## Running locally

The root Makefile is the source of truth for the available backend/database commands. It includes `.env` when present and exports its variables.

Start the local PostgreSQL service:

```bash
docker compose up -d
# equivalent Make target:
make compose
```

Run the Go API after providing a valid `DATABASE_URL` in the environment or root `.env`:

```bash
make run
```

The API listens on `http://localhost:8080`. Apply or roll back migrations with the existing Make targets:

```bash
make migrate-up
make migrate-down
```

These targets invoke the `goose` CLI as `goose up` and `goose down`; the CLI must be installed and configured for the current environment.

Run the Next.js development app from `client/`:

```bash
cd client
npm install
npm run dev
```

The frontend is served on `http://localhost:3000`. Catalog API helpers default to `http://localhost:8080`; set `NEXT_PUBLIC_API_URL` when the API is elsewhere.

The development Compose file starts only PostgreSQL. It does not start the Go API or the frontend.

## Running tests and checks

From the repository root:

```bash
make test   # cd server && go test ./...
make vet    # cd server && go vet ./...
make build  # cd server && go build ./...
make fmt    # cd server && gofmt -l . && goimports -l .
```

Current Go test coverage includes:

- Table-driven catalog handler tests using fake stores and `net/http/httptest`, covering validation, DTO responses, JSON errors, and repository-error mapping.
- Router wiring tests using `httptest`, including route parameters, method rejection, and unknown paths.
- PostgreSQL repository integration tests using testcontainers-go. They apply the Goose migration, seed small fixtures, test product/search/image/report behavior, and verify the configured cascade deletes. These tests require a usable Docker/testcontainers runtime.

There is currently no frontend test script or frontend test suite in `client/package.json`.

## Project structure

```text
calo-hub/
├── client/
│   ├── src/app/                  # Next.js routes and layouts
│   ├── src/catalog/              # Catalog UI and Go API client helpers
│   ├── src/components/           # Shared navigation/loading components
│   ├── src/convert/               # Presentational conversion UI
│   ├── src/dashboard/             # Dashboard UI
│   ├── src/eselector/             # Presentational selector UI
│   ├── src/extraction/            # Presentational extraction workbench UI
│   ├── src/types/                 # API-aligned TypeScript types
│   ├── Dockerfile
│   └── package.json
├── server/
│   ├── cmd/app/main.go            # Go API entrypoint
│   ├── internal/config/           # Configuration stub
│   ├── internal/db/               # pgx pool connection helper
│   ├── internal/dto/              # API response DTOs
│   ├── internal/handler/          # Implemented catalog handlers and stubs
│   ├── internal/models/            # Database-backed models
│   ├── internal/repository/        # PostgreSQL catalog repository
│   ├── internal/router/            # chi route table
│   ├── internal/service/            # Service placeholders
│   ├── migrations/                 # Goose schema migration
│   ├── data/                       # Catalog/report CSV source files
│   ├── go.mod
│   └── go.sum
├── docker-compose.yml              # Development PostgreSQL service
├── docker-compose.prod.yml         # Incomplete production service definition
├── Makefile
└── .env.example
```
