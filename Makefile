.PHONY: server-local-run client-local-run compose down migrate-up migrate-down
-include ~/.config/secrets/calohub/.env
export

server-local-run:
	cd server/ && go run cmd/app/main.go

client-local-run:
	cd client/ && npm run dev

compose:
	docker compose up -d

down:
	docker compose down

migrate-up:
	goose "$(DATABASE_URL)" up

migrate-down:
	goose "$(DATABASE_URL)" down

test:
	cd server && go test ./...

vet:
	cd server && go vet ./...

build:
	cd server && go build ./...

fmt:
	cd server && gofmt -l . && goimports -l .
