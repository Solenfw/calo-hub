.DEFAULT_GOAL := help

# Local secrets/configuration are optional. `make` also reads variables from
# the invoking shell, so CI can provide these without a local secrets file.
-include $(HOME)/.config/secrets/calohub/.env
export

# Override these on the command line when using a different setup, e.g.
# `make migrate-up MIGRATIONS_DIR=server/migrations GOOSE=goose`.
GOOSE ?= goose
GOOSE_DRIVER ?= postgres
MIGRATIONS_DIR ?= server/migrations
DB_URL ?= $(if $(strip $(DATABASE_URL)),$(DATABASE_URL),$(GOOSE_DBSTRING))
GO_VERSION_PACKAGE := ./...
SERVER_BINARY ?= dist/calo-hub

.PHONY: help server-local-run client-local-run client-install \
	compose compose-build compose-logs compose-ps down \
	migrate-create migrate-create-go migrate-up migrate-up-by-one migrate-up-to migrate-down \
	migrate-down-to migrate-down-all migrate-redo migrate-reset migrate-status \
	migrate-version migrate-validate migrate-fix install-goose \
	test test-race test-cover vet fmt fmt-check client-lint lint check \
	build build-server build-client clean

## Show the available development, build, test, and migration targets.
help:
	@echo 'Calo Hub development targets:'
	@echo '  server-local-run       Run the Go API locally (requires DATABASE_URL)'
	@echo '  client-local-run       Run the Next.js development server'
	@echo '  client-install         Install frontend npm dependencies'
	@echo '  compose                Build and start local Docker Compose services'
	@echo '  compose-build          Build Compose images without starting services'
	@echo '  compose-logs           Follow Compose service logs'
	@echo '  compose-ps             Show Compose service status'
	@echo '  down                   Stop Compose services'
	@echo '  migrate-create         Create SQL migration (use name=description)'
	@echo '  migrate-create-go      Create Go migration (use name=description)'
	@echo '  migrate-up             Apply all pending migrations'
	@echo '  migrate-up-by-one      Apply the next pending migration'
	@echo '  migrate-up-to          Apply migrations through version (use version=N)'
	@echo '  migrate-down           Roll back the most recent migration'
	@echo '  migrate-down-to        Roll back to version (use version=N)'
	@echo '  migrate-down-all       Roll back all migrations (destructive)'
	@echo '  migrate-redo           Re-run the most recent migration'
	@echo '  migrate-reset          Roll back and re-apply all migrations (destructive)'
	@echo '  migrate-status         Show applied and pending migrations'
	@echo '  migrate-version        Print the current database migration version'
	@echo '  migrate-validate       Validate migration files without changing the database'
	@echo '  migrate-fix            Renumber migration files to remove gaps'
	@echo '  install-goose          Install the pinned Goose CLI with go install'
	@echo '  test / test-race       Run backend tests, optionally with the race detector'
	@echo '  test-cover             Run backend tests with coverage summary'
	@echo '  vet / client-lint      Run Go vet / frontend ESLint'
	@echo '  fmt / fmt-check        Format Go files / check Go formatting'
	@echo '  build                  Build Go API and Next.js frontend'
	@echo '  build-server           Build Go API binary to $(SERVER_BINARY)'
	@echo '  build-client           Build the Next.js frontend'
	@echo '  clean                  Remove generated Go build output'
	@echo ''
	@echo 'Migration connection: DATABASE_URL (or GOOSE_DBSTRING); driver: GOOSE_DRIVER.'
	@echo 'Goose CLI: $(GOOSE); migration directory: $(MIGRATIONS_DIR).'

## Run the API using the current shell environment and local Go modules.
server-local-run:
	cd server && go run ./cmd/app

## Start the frontend in Next.js development mode.
client-local-run:
	cd client && npm run dev

## Install frontend dependencies from package-lock.json.
client-install:
	cd client && npm install

## Build images if needed and start the local database and API in the background.
compose:
	docker compose up -d --build

## Build the images declared by docker-compose.yml without starting containers.
compose-build:
	docker compose build

## Follow logs from all local Compose services.
compose-logs:
	docker compose logs -f

## Display local Compose service/container status.
compose-ps:
	docker compose ps

## Stop and remove local Compose containers and networks (keeps named volumes).
down:
	docker compose down

# Goose always receives an explicit migration directory and dialect. The
# database URL is taken from DATABASE_URL first, then GOOSE_DBSTRING.
## Create a timestamped SQL migration; provide the suffix with `name=...`.
migrate-create:
	@test -n "$(strip $(name))" || { echo 'Usage: make migrate-create name=add_products' >&2; exit 2; }
	$(GOOSE) -dir "$(MIGRATIONS_DIR)" create "$(name)" sql

## Create a timestamped Go migration; provide the suffix with `name=...`.
migrate-create-go:
	@test -n "$(strip $(name))" || { echo 'Usage: make migrate-create-go name=backfill_products' >&2; exit 2; }
	$(GOOSE) -dir "$(MIGRATIONS_DIR)" create "$(name)" go

## Apply all pending database migrations.
migrate-up:
	@test -n "$(strip $(DB_URL))" || { echo 'Set DATABASE_URL or GOOSE_DBSTRING before running migrations.' >&2; exit 2; }
	$(GOOSE) -dir "$(MIGRATIONS_DIR)" "$(GOOSE_DRIVER)" "$(DB_URL)" up

## Apply only the next pending database migration.
migrate-up-by-one:
	@test -n "$(strip $(DB_URL))" || { echo 'Set DATABASE_URL or GOOSE_DBSTRING before running migrations.' >&2; exit 2; }
	$(GOOSE) -dir "$(MIGRATIONS_DIR)" "$(GOOSE_DRIVER)" "$(DB_URL)" up-by-one

## Apply migrations through the requested version; provide `version=N`.
migrate-up-to:
	@test -n "$(strip $(version))" || { echo 'Usage: make migrate-up-to version=20260907221330' >&2; exit 2; }
	@test -n "$(strip $(DB_URL))" || { echo 'Set DATABASE_URL or GOOSE_DBSTRING before running migrations.' >&2; exit 2; }
	$(GOOSE) -dir "$(MIGRATIONS_DIR)" "$(GOOSE_DRIVER)" "$(DB_URL)" up-to "$(version)"

## Roll back the most recently applied database migration.
migrate-down:
	@test -n "$(strip $(DB_URL))" || { echo 'Set DATABASE_URL or GOOSE_DBSTRING before running migrations.' >&2; exit 2; }
	$(GOOSE) -dir "$(MIGRATIONS_DIR)" "$(GOOSE_DRIVER)" "$(DB_URL)" down

## Roll back until the database reaches the requested migration version.
migrate-down-to:
	@test -n "$(strip $(version))" || { echo 'Usage: make migrate-down-to version=20260907221330' >&2; exit 2; }
	@test -n "$(strip $(DB_URL))" || { echo 'Set DATABASE_URL or GOOSE_DBSTRING before running migrations.' >&2; exit 2; }
	$(GOOSE) -dir "$(MIGRATIONS_DIR)" "$(GOOSE_DRIVER)" "$(DB_URL)" down-to "$(version)"

## Roll back every applied migration. This drops the project schema/data.
migrate-down-all:
	@test -n "$(strip $(DB_URL))" || { echo 'Set DATABASE_URL or GOOSE_DBSTRING before running migrations.' >&2; exit 2; }
	$(GOOSE) -dir "$(MIGRATIONS_DIR)" "$(GOOSE_DRIVER)" "$(DB_URL)" down-to 0

## Roll back and then re-apply the latest migration.
migrate-redo:
	@test -n "$(strip $(DB_URL))" || { echo 'Set DATABASE_URL or GOOSE_DBSTRING before running migrations.' >&2; exit 2; }
	$(GOOSE) -dir "$(MIGRATIONS_DIR)" "$(GOOSE_DRIVER)" "$(DB_URL)" redo

## Roll back all migrations and apply them again; all schema data is lost.
migrate-reset:
	@test -n "$(strip $(DB_URL))" || { echo 'Set DATABASE_URL or GOOSE_DBSTRING before running migrations.' >&2; exit 2; }
	$(GOOSE) -dir "$(MIGRATIONS_DIR)" "$(GOOSE_DRIVER)" "$(DB_URL)" reset

## List database migration status, including pending files.
migrate-status:
	@test -n "$(strip $(DB_URL))" || { echo 'Set DATABASE_URL or GOOSE_DBSTRING before checking migrations.' >&2; exit 2; }
	$(GOOSE) -dir "$(MIGRATIONS_DIR)" "$(GOOSE_DRIVER)" "$(DB_URL)" status

## Print the database's current migration version.
migrate-version:
	@test -n "$(strip $(DB_URL))" || { echo 'Set DATABASE_URL or GOOSE_DBSTRING before checking migrations.' >&2; exit 2; }
	$(GOOSE) -dir "$(MIGRATIONS_DIR)" "$(GOOSE_DRIVER)" "$(DB_URL)" version

## Validate migration SQL/Go files without connecting to the database.
migrate-validate:
	$(GOOSE) -dir "$(MIGRATIONS_DIR)" validate

## Renumber migration files to remove gaps in their version sequence.
migrate-fix:
	$(GOOSE) -dir "$(MIGRATIONS_DIR)" fix

## Install the Goose CLI matching the version used by the Go module.
install-goose:
	cd server && go install github.com/pressly/goose/v3/cmd/goose@v3.28.0

## Run all Go backend tests.
test:
	cd server && go test $(GO_VERSION_PACKAGE)

## Run backend tests with Go's race detector.
test-race:
	cd server && go test -race $(GO_VERSION_PACKAGE)

## Run backend tests and print package coverage summaries.
test-cover:
	cd server && go test -cover $(GO_VERSION_PACKAGE)

## Run Go static analysis over backend packages.
vet:
	cd server && go vet $(GO_VERSION_PACKAGE)

## Format Go source files in the backend.
fmt:
	cd server && go fmt ./...

## Fail if any Go source file is not gofmt-formatted.
fmt-check:
	cd server && files=$$(find . -type f -name '*.go' -not -path './vendor/*'); test -z "$$(gofmt -l $$files)" || { echo 'Go files need formatting; run make fmt.' >&2; exit 1; }

## Run the frontend ESLint script.
client-lint:
	cd client && npm run lint

## Run Go vet and frontend lint checks.
lint: vet client-lint

## Run backend tests, lint checks, and both production builds.
check: test lint build

## Build both the API binary and the frontend production bundle.
build: build-server build-client

## Compile the Go API to the configured output path.
build-server:
	mkdir -p "$(dir $(abspath $(SERVER_BINARY)))"
	cd server && go build -o "$(abspath $(SERVER_BINARY))" ./cmd/app

## Build the Next.js production frontend.
build-client:
	cd client && npm run build

## Remove generated Go build output.
clean:
	rm -f "$(SERVER_BINARY)"
