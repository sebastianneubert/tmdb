.PHONY: build test run

build:
	go build -o tmdb cmd/main.go

test:
	go test ./...

run:
	go run cmd/main.go

