.PHONY: build lint test

build:
	go build

lint:
	go vet ./...
	golangci-lint run

test:
	go test ./...
