.PHONY: build run watch test clean

build:
	go build -o bin/api cmd/api/main.go

run:
	go run cmd/api/main.go

watch:
	air

test:
	go test ./...

clean:
	rm -rf bin/
