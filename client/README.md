# Calo Hub frontend

This frontend is the user-facing layer for the catalog and Martin report workflows. It talks to the Go backend for search, report management, import/export, and PDF generation.

> The report workflow is the active feature set. The older selector-focused flow is no longer the primary app path.

## Stack

- Next.js 16
- React 19
- TypeScript 5
- Tailwind CSS 4
- xlsx for spreadsheet import/export
- Lucide icons for the report controls

## Local run

```bash
cd client
npm install
npm run dev
```

The app runs at `http://localhost:3000`.

If the backend is not running on the default port, set the API base explicitly:

```bash
NEXT_PUBLIC_API_URL=http://localhost:8080 npm run dev
```

## Current routes

The shell route group is under `src/app/(shell)` and is used for the main application pages.

- `/` -> landing/dashboard page
- `/dashboard` -> dashboard shell
- `/catalog` -> searchable catalog UI
- `/convert` -> conversion utility screen
- `/extraction` -> extraction flow placeholder
- `/report` -> Martin report management screen

The active report UI is implemented in `src/report/Report.tsx` and mounted by `src/app/(shell)/report/page.tsx`.

## Report feature set

The report page supports:

- listing existing Martin reports
- creating a new report
- selecting and deleting a report
- editing rows in the report table
- saving updated report products
- importing spreadsheet rows into a report
- downloading the import template file
- exporting the current report as a PDF

The client helper file `src/report/reportHandlers.ts` wraps the backend calls for this flow.

## API client wiring

The client uses Go API routes such as:

```text
GET /catalog/products?q=...
GET /catalog/products/{code}
GET /catalog/images/{code}
GET /report/martin
GET /report/martin/{name}
GET /report/martin/all/{report_id}
GET /report/martin/{name}/pdf
GET /report/martin/template
POST /report/martin
PUT /report/martin/{report_id}/products
DELETE /report/martin/{report_id}
```

The response models in `src/types` align with the backend DTO output and include the Martin report product shape used by the report table.

## Scripts

```bash
npm run dev
npm run build
npm start
npm run lint
```

There is currently no dedicated frontend test suite; the primary verification step is linting plus the backend test suite run separately.
