.PHONY: dev build test lint run clean \
       docker-up docker-down docker-build \
       migrate-up migrate-down migrate-create \
       tidy

# ─── Variables ────────────────────────────────────────────────
APP_NAME    := luxbiss-server
MAIN_PATH   := ./cmd/api
BINARY      := ./bin/$(APP_NAME)
MIGRATION_DIR := ./migrations
DB_URL      ?= postgres://postgres:postgres@localhost:5432/luxbiss?sslmode=disable

# ─── Development ──────────────────────────────────────────────

## dev: Run with hot-reload (Air)
dev:
	@air -c .air.toml

## run: Run without hot-reload
run:
	@go run $(MAIN_PATH)

## build: Compile production binary
build:
	@echo "Building $(APP_NAME)..."
	@CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BINARY) $(MAIN_PATH)
	@echo "Binary: $(BINARY)"

## clean: Remove build artifacts
clean:
	@rm -rf ./bin ./tmp
	@echo "Cleaned."

# ─── Quality ─────────────────────────────────────────────────

## test: Run all tests
test:
	@go test -v -race -count=1 ./...

## test-cover: Run tests with coverage
test-cover:
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## lint: Run linter
lint:
	@golangci-lint run ./...

## tidy: Tidy and verify Go modules
tidy:
	@go mod tidy
	@go mod verify

# ─── Docker ──────────────────────────────────────────────────

## docker-up: Start all services
docker-up:
	@docker compose up -d

## docker-down: Stop all services
docker-down:
	@docker compose down

## docker-build: Build Docker image
docker-build:
	@docker compose build

## docker-logs: View logs
docker-logs:
	@docker compose logs -f api

# ─── Database Migrations ────────────────────────────────────

## migrate-up: Apply all migrations
migrate-up:
	@migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" up

## migrate-down: Rollback last migration
migrate-down:
	@migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" down 1

## migrate-create: Create a new migration (usage: make migrate-create name=create_products)
migrate-create:
	@migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(name)

# ─── Help ────────────────────────────────────────────────────

## help: Show available commands
help:
	@echo "Available commands:"
	@grep -E '^## ' Makefile | sed 's/## /  /'
