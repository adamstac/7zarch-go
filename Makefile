# 7zarch-go Makefile

VERSION ?= 0.1.0-dev
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)

# Test and coverage settings
COVERAGE_DIR := coverage
COVERAGE_PROFILE := $(COVERAGE_DIR)/coverage.out
COVERAGE_HTML := $(COVERAGE_DIR)/coverage.html

# Build targets
.PHONY: build build-all clean test install deps help
.PHONY: test-all test-unit test-integration test-bench test-coverage test-coverage-html
.PHONY: test-mas test-resolver test-registry test-edge-cases
.PHONY: dev-tools debug-registry mas-inspect mas-stats
.PHONY: dev dist release validate goreleaser-build

build: ## Build for current platform
	go build -ldflags "$(LDFLAGS)" -o 7zarch-go

build-all: ## Build for all platforms (legacy method)
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/7zarch-go-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/7zarch-go-darwin-arm64
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/7zarch-go-linux-amd64
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/7zarch-go-linux-arm64
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/7zarch-go-windows-amd64.exe

##@ Goreleaser Targets

dev: ## Local build and install to ~/bin (using Goreleaser)
	goreleaser build --single-target --clean --snapshot
	mkdir -p ~/bin && cp dist/7zarch-go_*/7zarch-go ~/bin/

dist: ## Build for current platform (using Goreleaser)
	goreleaser build --single-target --clean --snapshot

release: ## Create release (CI only - requires git tag)
	goreleaser release --clean

validate: ## Validate Goreleaser configuration
	goreleaser check

goreleaser-build: ## Build all platforms using Goreleaser (snapshot mode)
	goreleaser build --clean --snapshot

clean: ## Clean build artifacts
	rm -f 7zarch-go
	rm -rf dist/
	rm -rf $(COVERAGE_DIR)/
	rm -f *.7z *.log *.sha256

##@ Testing Targets

test: ## Run basic tests
	go test ./...

test-all: test-unit test-integration test-edge-cases ## Run all test suites

test-unit: ## Run unit tests only
	go test -short ./...

test-integration: ## Run integration tests
	go test -run Integration ./...

test-mas: ## Run MAS-specific tests
	go test -v ./internal/storage/...

test-resolver: ## Run resolver tests
	go test -v -run Resolve ./internal/storage/

test-registry: ## Run registry tests  
	go test -v -run Registry ./internal/storage/

test-edge-cases: ## Run edge case tests
	go test -v -run Edge ./internal/storage/

test-verbose: ## Run tests with verbose output
	go test -v ./...

test-race: ## Run tests with race detection
	go test -race ./...

test-coverage: ## Generate test coverage report
	@mkdir -p $(COVERAGE_DIR)
	go test -coverprofile=$(COVERAGE_PROFILE) -covermode=atomic ./...
	go tool cover -func=$(COVERAGE_PROFILE)

test-coverage-html: test-coverage ## Generate HTML coverage report
	go tool cover -html=$(COVERAGE_PROFILE) -o $(COVERAGE_HTML)
	@echo "Coverage report: $(COVERAGE_HTML)"

##@ DDD Framework Validation (7EP-0020)

validate-framework: ## Run complete DDD framework validation suite
	@echo "🔍 Running DDD Framework Validation Suite..."
	@echo "============================================="
	@cd scripts && go build validate-framework.go && ./validate-framework ..
	@echo ""
	@cd scripts && go build validate-consistency.go && ./validate-consistency ..
	@echo ""
	@cd scripts && go build validate-git-patterns.go && ./validate-git-patterns ..
	@echo ""
	@echo "📊 DDD Framework validation complete"
	@echo ""
	@echo "📊 Framework Health Summary:"
	@./scripts/framework-health.sh | tail -5

validate-framework-roles: ## Validate role files only
	@echo "🔍 Validating role files only..."
	@cd scripts && go build validate-framework.go && ./validate-framework ..

validate-framework-consistency: ## Validate cross-document consistency
	@echo "🔍 Validating cross-document consistency..."
	@cd scripts && go build validate-consistency.go && ./validate-consistency ..

validate-framework-git: ## Validate git patterns
	@echo "🔍 Validating git patterns..."
	@cd scripts && go build validate-git-patterns.go && ./validate-git-patterns ..

validate-framework-integration: ## Test complete framework integration
	@echo "🔍 Testing framework integration..."
	@./scripts/test-agent-lifecycle.sh
	@echo ""
	@./scripts/test-workflows.sh

framework-health: ## Generate DDD framework health dashboard
	@echo "📊 Generating framework health dashboard..."
	@./scripts/framework-health.sh

##@ Benchmarking Targets  

bench: ## Run all benchmarks
	go test -bench=. ./...

bench-mas: ## Run MAS benchmarks
	go test -bench=. ./internal/storage/

bench-resolver: ## Run resolver benchmarks
	go test -bench=BenchmarkResolver ./internal/storage/

bench-registry: ## Run registry benchmarks
	go test -bench=BenchmarkRegistry ./internal/storage/

bench-scalability: ## Run scalability benchmarks
	go test -bench=BenchmarkScalability ./internal/storage/

bench-memory: ## Run memory usage benchmarks
	go test -bench=BenchmarkMemory -benchmem ./internal/storage/

##@ Development and Utilities

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

dev-legacy: build symlink ## Build and symlink for development (legacy method)

dev-tools: build mas-inspect mas-stats ## Build development utilities

format: ## Format code
	go fmt ./...

lint: ## Run linters
	go vet ./...

vet-shadow: ## Check for variable shadowing
	go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow@latest
	go vet -vettool=$(shell which shadow) ./...

##@ MAS Development and Debugging

mas-inspect: build ## Inspect MAS registry database
	@echo "#!/bin/bash" > mas-inspect
	@echo "# MAS Registry Inspector" >> mas-inspect
	@echo "DB_PATH=\$${1:-~/.7zarch-go/registry.db}" >> mas-inspect
	@echo "if [ ! -f \"\$$DB_PATH\" ]; then" >> mas-inspect
	@echo "  echo \"Registry not found: \$$DB_PATH\"" >> mas-inspect
	@echo "  exit 1" >> mas-inspect
	@echo "fi" >> mas-inspect
	@echo "echo \"=== MAS Registry: \$$DB_PATH ===\"" >> mas-inspect
	@echo "sqlite3 \"\$$DB_PATH\" \".schema archives\"" >> mas-inspect
	@echo "echo" >> mas-inspect
	@echo "sqlite3 \"\$$DB_PATH\" \"SELECT COUNT(*) as total_archives FROM archives;\"" >> mas-inspect
	@echo "sqlite3 \"\$$DB_PATH\" \"SELECT managed, COUNT(*) FROM archives GROUP BY managed;\"" >> mas-inspect
	@echo "sqlite3 \"\$$DB_PATH\" \"SELECT status, COUNT(*) FROM archives GROUP BY status;\"" >> mas-inspect
	@chmod +x mas-inspect
	@echo "Created mas-inspect utility"

mas-stats: build ## Generate MAS statistics
	@echo "#!/bin/bash" > mas-stats
	@echo "# MAS Statistics Generator" >> mas-stats
	@echo "DB_PATH=\$${1:-~/.7zarch-go/registry.db}" >> mas-stats
	@echo "if [ ! -f \"\$$DB_PATH\" ]; then" >> mas-stats
	@echo "  echo \"Registry not found: \$$DB_PATH\"" >> mas-stats
	@echo "  exit 1" >> mas-stats
	@echo "fi" >> mas-stats
	@echo "echo \"=== MAS Registry Statistics ===\"" >> mas-stats
	@echo "sqlite3 \"\$$DB_PATH\" -header -column \"SELECT" >> mas-stats
	@echo "  COUNT(*) as Total," >> mas-stats
	@echo "  SUM(CASE WHEN managed = 1 THEN 1 ELSE 0 END) as Managed," >> mas-stats
	@echo "  SUM(CASE WHEN managed = 0 THEN 1 ELSE 0 END) as External," >> mas-stats
	@echo "  SUM(CASE WHEN uploaded = 1 THEN 1 ELSE 0 END) as Uploaded," >> mas-stats
	@echo "  SUM(size) as TotalBytes" >> mas-stats
	@echo "FROM archives;\"" >> mas-stats
	@echo "echo" >> mas-stats
	@echo "echo \"=== Recent Archives ===\"" >> mas-stats
	@echo "sqlite3 \"\$$DB_PATH\" -header -column \"SELECT uid, name, managed, status, created FROM archives ORDER BY created DESC LIMIT 10;\"" >> mas-stats
	@echo "echo" >> mas-stats
	@echo "echo \"=== Profile Distribution ===\"" >> mas-stats
	@echo "sqlite3 \"\$$DB_PATH\" -header -column \"SELECT profile, COUNT(*) as count FROM archives GROUP BY profile ORDER BY count DESC;\"" >> mas-stats
	@chmod +x mas-stats
	@echo "Created mas-stats utility"

debug-registry: build ## Debug registry issues
	@echo "=== Registry Debug Information ==="
	@echo "Current directory: $(PWD)"
	@echo "Expected registry: ~/.7zarch-go/registry.db"
	@ls -la ~/.7zarch-go/ 2>/dev/null || echo "No ~/.7zarch-go directory found"
	@echo
	@echo "=== Testing registry creation ==="
	@mkdir -p /tmp/7zarch-test
	@echo "Testing registry in /tmp/7zarch-test"
	./7zarch-go create --output /tmp/7zarch-test/test.7z --dry-run demo-files 2>/dev/null || echo "Registry test completed"

##@ Quick Testing

run-create: build ## Test create command
	./7zarch-go create --dry-run demo-files || echo "Create demo-files directory first"

run-test: build ## Test archive testing command
	./7zarch-go test --dry-run *.7z 2>/dev/null || echo "No .7z files found"

run-list: build ## Test list command
	./7zarch-go list || echo "No archives in registry yet"

run-mas-show: build ## Test mas show command
	@echo "Testing mas show command..."
	./7zarch-go mas show test 2>/dev/null || echo "No archives to show yet"

# Help target
help: ## Show this help message
	@echo "7zarch-go Build Commands:"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

# Default target
.DEFAULT_GOAL := help