# Calo Hub frontend

This Next.js frontend provides the catalog UI and several in-progress workflow screens, calling the Go server's catalog API where integration exists.

> **Actively under development.** The catalog page is connected to the Go API; the dashboard, conversion, extraction, and eSelector pages are currently presentational or demo UI and are not end-to-end backend features.

## Stack

- Next.js `^16.2.6`
- React `19.2.4` and React DOM `19.2.4`
- Tailwind CSS `^4` with `@tailwindcss/postcss`
- Motion `^12.38.0`
- Lucide React `^1.7.0`
- `clsx` and `tailwind-merge` for class composition
- TypeScript `^5` and ESLint 9 with `eslint-config-next`

## Running locally

Install dependencies and start the development server from this directory:

```bash
npm install
npm run dev
```

The app is served at `http://localhost:3000`.

Catalog API requests use `NEXT_PUBLIC_API_URL` when it is set. If it is not set, the client defaults to `http://localhost:8080`:

```bash
NEXT_PUBLIC_API_URL=http://localhost:8080 npm run dev
```

The Go server must be running separately for catalog searches and image requests to return backend data. The development Docker Compose file starts PostgreSQL only; it does not start either application.

Other package scripts are:

```bash
npm run build
npm start
npm run lint
```

`build` creates the Next.js production build, `start` serves that build, and `lint` runs ESLint. No frontend test script is defined.

## Page and route structure

The route group is `src/app/(shell)`. The parentheses are a Next.js route group and do not appear in the URLs.

| URL | Current behavior |
| --- | --- |
| `/` | Renders the `Dashboard` component. The page is a static dashboard with example statistics, activity, and file rows; it does not call the Go API. |
| `/dashboard` | Renders the same `Dashboard` component and is likewise presentational/static. |
| `/catalog` | Renders the interactive catalog. `searchProducts` calls `GET /catalog/products?q=...`; selecting a product calls `GET /catalog/images/{code}`. Search results, brand filtering, language selection, image display, and the lightbox are client-side behavior. |
| `/convert` | Renders `QuickConvert`, a static conversion form/mockup. Its controls do not call a backend endpoint. |
| `/extraction` | Renders the extraction workbench with local page-selection state, hardcoded document previews, and hardcoded extracted text. It does not upload, OCR, or call the Go server. |
| `/eselector` | Renders a static selector form with hardcoded example matches and controls. It does not call a backend endpoint. |

`src/app/(shell)/layout.tsx` supplies the shared top/side/bottom navigation and page transition animation. `src/app/layout.tsx` supplies metadata, fonts, and global styles.

## Go API client wiring

`src/catalog/productHandlers.ts` defines client helpers for these Go routes:

```text
GET /catalog/products?q=...
GET /catalog/products/{code}
GET /catalog/images/{code}
GET /catalog/report/martin
GET /catalog/report/martin/{name}
GET /catalog/report/martin/all/{report_id}
```

The catalog page currently uses the product search and image helpers. The Martin report helpers exist but are not currently used by one of the listed route pages.

The TypeScript response types in `src/types/product.ts` mirror the Go DTO names and JSON fields, including `image: string | null` for Martin report products and `ApiErrorResponse` for server errors.

## Testing

There is currently no `test` script in `package.json`, no frontend test directory, and no frontend test suite in this directory. Use `npm run lint` for the configured static check; it is not a replacement for component or integration tests.

