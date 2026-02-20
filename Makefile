.PHONY: build run watch test clean docker-up docker-down migrate

build:
	go build -ldflags="-w -s" -o bin/api ./cmd/api/main.go

run:
	go run ./cmd/api/main.go

watch:
	air

test:
	go test ./... -v -race -cover

clean:
	rm -rf bin/ tmp/

# Docker
docker-up:
	docker compose up -d --build

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f app

# Database migrations (requires golang-migrate)
migrate-up:
	migrate -path migrations -database "$$DATABASE_URL" up

migrate-down:
	migrate -path migrations -database "$$DATABASE_URL" down
