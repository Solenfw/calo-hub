# Calo Hub

Calo Hub is a full-stack medical instrument catalog and Martin report workflow. The backend is implemented in Go and exposes the catalog/report API, while the frontend is a Next.js app that handles product search, report management, spreadsheet imports, and PDF export.

> Status: the core catalog and Martin report flows are in place and the backend/frontend are wired together. The project is still evolving, but the main report and catalog paths are active instead of placeholder-only.

## Stack

- Backend: Go, chi, pgx, PostgreSQL, goose
- Frontend: Next.js, React, TypeScript, Tailwind CSS, xlsx, Lucide icons
- PDF generation: embedded Roboto TTF with go-pdf/fpdf
- Local tooling: Docker Compose, Go tests, and npm scripts

## Current working flows

### Backend

The Go API boots from [server/cmd/app/main.go](server/cmd/app/main.go) and wires together the database pool, catalog repository, report repository, and handlers.

The router registers the main routes:

- `GET /catalog/products?q=...`
- `GET /catalog/products/{code}`
- `GET /catalog/images/{code}`
- `GET /report/martin`
- `GET /report/martin/{name}`
- `GET /report/martin/all/{report_id}`
- `GET /report/martin/{name}/pdf`
- `GET /report/martin/template`
- `POST /report/martin`
- `PUT /report/martin/{report_id}/products`
- `DELETE /report/martin/{report_id}`

The report flow includes backend image resolution: when a report row has no explicit image URL, the backend resolves it from catalog data before generating the PDF.

### Frontend

The active UI is centered on the report page at `/report`.

- The client-side report screen lets the user create, load, save, update, import, and export Martin reports.
- Spreadsheet imports use `xlsx` and accept tabular data in the expected code/description/quantity format.
- PDF export calls the backend report export endpoint and downloads the resulting file.
- The old selector-only flow has been replaced by the report workflow.

## Local setup

### Start the database

```bash
docker compose up -d
```

### Run the Go API

From the project root or inside `server/`:

```bash
cd server
go run cmd/app/main.go
```

The API listens on `http://localhost:8080`.

### Run the frontend

```bash
cd client
npm install
npm run dev
```

The frontend is served at `http://localhost:3000`.

If the backend is not on the default port, set `NEXT_PUBLIC_API_URL` before starting the app.

## Testing

Run the backend tests from `server/`:

```bash
cd server
go test ./...
```

Run the client lint check from `client/`:

```bash
cd client
npm run lint
```

## Repository layout

```text
calo-hub/
├── client/                      # Next.js frontend
│   ├── src/app/(shell)/report/ # report page route
│   ├── src/report/              # report UI + handlers
│   ├── src/catalog/             # catalog UI + API calls
│   ├── src/components/          # shared UI
│   ├── src/types/               # API TypeScript models
│   ├── package.json
│   └── README.md
├── server/                      # Go backend
│   ├── cmd/app/main.go          # entrypoint
│   ├── internal/handler/        # catalog + report handlers
│   ├── internal/repository/      # PostgreSQL repositories
│   ├── internal/service/pdf/    # PDF renderer and font embedding
│   ├── internal/router/          # chi route registration
│   ├── migrations/              # database migrations
│   ├── data/                    # CSV/XLSX resources
│   └── README.md
├── docker-compose.yml
├── docker-compose.prod.yml
├── Makefile
├── README.md
└── .env.example
```

## Notes

- This monorepo is still being refined, but the catalog/report backend and the report frontend are the current active workstreams.
- The PDF generator is expected to work with UTF-8-friendly content using the embedded Roboto font, and the backend resolves missing product images internally.
- The project is not yet a fully polished production deployment; it is a working development codebase with an active report feature set.
