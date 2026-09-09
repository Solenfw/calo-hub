.PHONY: test vet build fmt migrate-up migrate-down compose down run
-include .env
export


compose:
	docker compose up -d

down:
	docker compose down

migrate-up:
	goose up

migrate-down:
	goose down

run:
	go run server/cmd/app/main.go

test:
	cd server && go test ./...

vet:
	cd server && go vet ./...

build:
	cd server && go build ./...

fmt:
	cd server && gofmt -l . && goimports -l .
