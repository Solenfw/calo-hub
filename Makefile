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
	go run server/app/main.go