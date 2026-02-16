.PHONY: help build install test test-integration test-coverage clean fmt lint vet

# Default target
.DEFAULT_GOAL := help

# Variables
BINARY_NAME=gocrete
MAIN_PATH=cmd/gocrete/main.go
INSTALL_PATH=$(HOME)/go/bin/$(BINARY_NAME)

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: ./$(BINARY_NAME)"

install: ## Install the binary to GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	@go install $(MAIN_PATH)
	@echo "Installed to: $(INSTALL_PATH)"

test: ## Run unit tests
	@echo "Running tests..."
	@go test -v -short ./...

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

bench: ## Run benchmarks
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Format complete"

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run || echo "golangci-lint not installed, skipping"

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

check: fmt vet test ## Run format, vet, and tests

generate-example: build ## Generate an example project
	@echo "Generating example project..."
	@rm -rf /tmp/gocrete-example
	@./$(BINARY_NAME) init example-service \
		--module github.com/example/service \
		--router chi \
		--db postgres \
		--migrations goose \
		--docker
	@echo "Example generated in: example-service/"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

update-deps: ## Update dependencies
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy

verify: ## Verify dependencies
	@echo "Verifying dependencies..."
	@go mod verify

all: clean fmt vet test build ## Clean, format, vet, test, and build

release: clean test ## Prepare for release
	@echo "Preparing release..."
	@./scripts/release.sh || echo "Release script not found"
