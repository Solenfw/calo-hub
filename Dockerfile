# Go API image. Build context is the repo root (see docker-compose.yml).

FROM golang:1.26-alpine AS builder
WORKDIR /src

# Cache module downloads separately from source changes.
COPY server/go.mod server/go.sum ./
RUN go mod download

COPY server/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/calo-hub ./cmd/app

# ==========================================
# distroless/static ships CA certs (needed for PDF image fetches over HTTPS)
# and runs as a non-root user.
FROM gcr.io/distroless/static-debian12:nonroot AS prod
WORKDIR /app

COPY --from=builder /out/calo-hub ./calo-hub
# The import template is resolved from ./data/sheet/template.xlsx at runtime.
COPY --from=builder /src/data/sheet/template.xlsx ./data/sheet/template.xlsx

ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["/app/calo-hub"]
