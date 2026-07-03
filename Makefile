.PHONY: help fmt lint build test coverage check clean

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "  fmt       format all Go files"
	@echo "  lint      run golangci-lint"
	@echo "  build     build the binary to bin/"
	@echo "  test      run all tests"
	@echo "  coverage  generate HTML coverage report"
	@echo "  check     fmt + lint + build + test (pre-PR)"
	@echo "  clean     remove build artifacts"

fmt:
	go fmt ./...

lint:
	golangci-lint run

build:
	go build -o bin/groupie-tracker ./cmd

test:
	go test -v ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

check: fmt lint build test
	@echo "all checks passed"

clean:
	rm -rf bin/ coverage.out coverage.html
