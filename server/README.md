# Calo Hub Go service

This service exposes the catalog and Martin report API backed by PostgreSQL and serves the report export and import flows used by the front end.

> The catalog and Martin report features are the active backend paths. The project is still evolving, but the main API and report generation flow are implemented rather than placeholder-only.

## Current implementation

### API entrypoint

The app boots in `cmd/app/main.go` and wires together:

- database pool from `GOOSE_DBSTRING`
- catalog repository
- report repository
- catalog handler
- report handler
- chi router

The server listens on `:8080` and logs each request status/method/path.

### Registered routes

```text
GET /catalog/products?q=...
GET /catalog/products/{code}
GET /catalog/images/{code}
GET /report/martin
GET /report/martin/template
GET /report/martin/{name}
GET /report/martin/all/{report_id}
GET /report/martin/{name}/pdf
POST /report/martin
PUT /report/martin/{report_id}/products
DELETE /report/martin/{report_id}
```

The report handler resolves product images internally when a snapshot row has no image URL and then generates the PDF using the embedded Roboto font asset.

### PDF generation

The PDF renderer in `internal/service/pdf/pdf_drawer.go` is no longer a stub. It:

- builds a PDF report document from Martin product rows
- inserts report detail rows and a summary table
- fetches remote product images and embeds them when available
- uses an embedded `Roboto-Regular.ttf` asset for UTF-8-friendly text

### Repository layer

The repository package includes both catalog and report persistence logic:

- catalog search and product lookup
- product image lookup
- Martin report list CRUD
- report product snapshot reads and updates

Missing rows map to `ErrNotFound` and are converted to HTTP 404 responses in the handlers.

## Local run

From the `server/` directory:

```bash
go run cmd/app/main.go
```

The environment is expected to provide a valid Postgres connection string, typically via `GOOSE_DBSTRING` in the local environment or user config.

## Migrations

The schema lives under `migrations/` and defines the product, image, report list, and report-product tables used by the catalog/report features.

Typical migration commands:

```bash
goose -dir migrations postgres "$GOOSE_DBSTRING" up
goose -dir migrations postgres "$GOOSE_DBSTRING" down
```

## Testing

Run the backend suite from `server/`:

```bash
go test ./...
```

The repository includes tests for:

- catalog handlers
- report handlers
- router registration and dispatch
- repository behavior for catalog and Martin report data

## Current state notes

- The main development focus is the catalog + Martin report flow.
- The app currently expects frontend-driven report operations and backend-managed PDF generation.
- Some legacy and experimental code may remain elsewhere in the repo, but the current active paths are the catalog and report APIs.
