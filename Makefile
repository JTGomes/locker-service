ifneq (,$(wildcard .env))
	include .env
	export
endif

DATABASE_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(POSTGRES_DB)?sslmode=disable

.PHONY: run install test test-race  \
        migrate-up migrate-down  \
        docker-run docker-down

run:
	go run ./cmd/api

test:
	go test ./...

install:
	go mod download

test-race:
	go test -race ./...

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down: 
	migrate -path migrations -database "$(DATABASE_URL)" down 1

docker-run: 
	docker compose up --build

docker-down:
	docker compose down