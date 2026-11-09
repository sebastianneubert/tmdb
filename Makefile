.PHONY: build test run

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

build:
	go build $(LDFLAGS) -o tmdb cmd/main.go

test:
	go test -v ./test/...

run:
	go run cmd/main.go

