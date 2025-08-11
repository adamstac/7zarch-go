# 7zarch-go Makefile

VERSION ?= 0.1.0-dev
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)

# Build targets
.PHONY: build build-all clean test install deps help

build: ## Build for current platform
	go build -ldflags "$(LDFLAGS)" -o 7zarch-go

build-all: ## Build for all platforms
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/7zarch-go-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/7zarch-go-darwin-arm64
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/7zarch-go-linux-amd64
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/7zarch-go-linux-arm64
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/7zarch-go-windows-amd64.exe

clean: ## Clean build artifacts
	rm -f 7zarch-go
	rm -rf dist/
	rm -f *.7z *.log *.sha256

test: ## Run tests
	go test ./...

test-verbose: ## Run tests with verbose output
	go test -v ./...

test-race: ## Run tests with race detection
	go test -race ./...

bench: ## Run benchmarks
	go test -bench=. ./...

deps: ## Download dependencies
	go mod download
	go mod tidy

install: build ## Install to ~/bin
	mkdir -p ~/bin
	cp 7zarch-go ~/bin/
	@echo "Installed to ~/bin/7zarch-go"

symlink: build ## Create symlink in ~/bin
	mkdir -p ~/bin
	ln -sf $(PWD)/7zarch-go ~/bin/7zarch-go
	@echo "Symlinked $(PWD)/7zarch-go to ~/bin/7zarch-go"

# Development helpers
dev: build symlink ## Build and symlink for development

run-create: build ## Test create command
	./7zarch-go create --dry-run test-data || echo "Create test-data directory first"

run-test: build ## Test archive testing command
	./7zarch-go test --dry-run *.7z 2>/dev/null || echo "No .7z files found"

format: ## Format code
	go fmt ./...

lint: ## Run linters
	go vet ./...

# Help target
help: ## Show this help message
	@echo "7zarch-go Build Commands:"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

# Default target
.DEFAULT_GOAL := help