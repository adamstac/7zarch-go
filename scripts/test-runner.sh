#!/bin/bash
# MAS Testing and Development Utilities
# Comprehensive test runner for 7zarch-go MAS development

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
COVERAGE_DIR="coverage"
COVERAGE_PROFILE="$COVERAGE_DIR/coverage.out"
COVERAGE_HTML="$COVERAGE_DIR/coverage.html"
TEST_TIMEOUT="10m"

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_dependencies() {
    log_info "Checking dependencies..."
    
    if ! command -v go &> /dev/null; then
        log_error "Go is not installed or not in PATH"
        exit 1
    fi
    
    if ! command -v sqlite3 &> /dev/null; then
        log_warning "sqlite3 not found - some debug utilities may not work"
    fi
    
    log_success "Dependencies check passed"
}

setup_coverage() {
    log_info "Setting up coverage directory..."
    mkdir -p "$COVERAGE_DIR"
}

run_unit_tests() {
    log_info "Running unit tests..."
    go test -short -timeout "$TEST_TIMEOUT" ./...
    log_success "Unit tests completed"
}

run_integration_tests() {
    log_info "Running integration tests..."
    go test -run Integration -timeout "$TEST_TIMEOUT" ./...
    log_success "Integration tests completed"
}

run_mas_tests() {
    log_info "Running MAS-specific tests..."
    go test -v -timeout "$TEST_TIMEOUT" ./internal/storage/...
    log_success "MAS tests completed"
}

run_resolver_tests() {
    log_info "Running resolver tests..."
    go test -v -run Resolve -timeout "$TEST_TIMEOUT" ./internal/storage/
    log_success "Resolver tests completed"
}

run_registry_tests() {
    log_info "Running registry tests..."
    go test -v -run Registry -timeout "$TEST_TIMEOUT" ./internal/storage/
    log_success "Registry tests completed"
}

run_edge_case_tests() {
    log_info "Running edge case tests..."
    go test -v -run Edge -timeout "$TEST_TIMEOUT" ./internal/storage/
    log_success "Edge case tests completed"
}

run_race_tests() {
    log_info "Running tests with race detection..."
    go test -race -short -timeout "$TEST_TIMEOUT" ./...
    log_success "Race detection tests completed"
}

run_benchmarks() {
    log_info "Running benchmarks..."
    go test -bench=. -timeout "$TEST_TIMEOUT" ./internal/storage/ | tee benchmarks.log
    log_success "Benchmarks completed - results saved to benchmarks.log"
}

run_coverage() {
    log_info "Generating coverage report..."
    setup_coverage
    go test -coverprofile="$COVERAGE_PROFILE" -covermode=atomic -timeout "$TEST_TIMEOUT" ./...
    
    if [ -f "$COVERAGE_PROFILE" ]; then
        go tool cover -func="$COVERAGE_PROFILE"
        go tool cover -html="$COVERAGE_PROFILE" -o "$COVERAGE_HTML"
        log_success "Coverage report generated: $COVERAGE_HTML"
    else
        log_error "Coverage profile not generated"
        return 1
    fi
}

run_comprehensive_test() {
    log_info "Running comprehensive test suite..."
    echo "========================================"
    
    check_dependencies
    setup_coverage
    
    echo "========================================"
    run_unit_tests
    echo "========================================"
    run_integration_tests
    echo "========================================"
    run_mas_tests
    echo "========================================"
    run_resolver_tests
    echo "========================================"
    run_registry_tests
    echo "========================================"
    run_edge_case_tests
    echo "========================================"
    run_race_tests
    echo "========================================"
    run_coverage
    echo "========================================"
    
    log_success "Comprehensive test suite completed!"
    echo "Coverage report: $COVERAGE_HTML"
    
    if [ -f benchmarks.log ]; then
        echo "Benchmark results: benchmarks.log"
    fi
}

validate_mas_implementation() {
    log_info "Validating MAS implementation readiness..."
    
    # Check for expected files
    local missing_files=()
    
    if [ ! -f "internal/storage/resolver.go" ]; then
        missing_files+=("internal/storage/resolver.go")
    fi
    
    if [ ! -f "internal/storage/registry.go" ]; then
        missing_files+=("internal/storage/registry.go")
    fi
    
    if [ ! -f "cmd/mas_show.go" ]; then
        missing_files+=("cmd/mas_show.go")
    fi
    
    if [ ${#missing_files[@]} -gt 0 ]; then
        log_warning "Some MAS files are missing (expected during development):"
        for file in "${missing_files[@]}"; do
            echo "  - $file"
        done
    else
        log_success "All expected MAS files are present"
    fi
    
    # Run tests to validate what's implemented
    log_info "Running implementation validation tests..."
    
    if go test -run TestResolve ./internal/storage/ > /dev/null 2>&1; then
        log_success "Resolver implementation tests pass"
    else
        log_warning "Resolver implementation tests not passing (expected during development)"
    fi
    
    if go test -run TestRegistry ./internal/storage/ > /dev/null 2>&1; then
        log_success "Registry implementation tests pass"
    else
        log_warning "Registry implementation tests not passing (expected during development)"
    fi
}

show_help() {
    echo "MAS Testing and Development Utilities"
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  unit           Run unit tests only"
    echo "  integration    Run integration tests"
    echo "  mas            Run MAS-specific tests"
    echo "  resolver       Run resolver tests"
    echo "  registry       Run registry tests"
    echo "  edge-cases     Run edge case tests"
    echo "  race           Run tests with race detection"
    echo "  bench          Run benchmarks"
    echo "  coverage       Generate coverage report"
    echo "  comprehensive  Run full test suite"
    echo "  validate       Validate MAS implementation readiness"
    echo "  help           Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 comprehensive    # Run everything"
    echo "  $0 mas              # Test MAS functionality"
    echo "  $0 coverage         # Generate coverage report"
    echo "  $0 validate         # Check implementation status"
}

# Main command handling
case "${1:-help}" in
    unit)
        check_dependencies
        run_unit_tests
        ;;
    integration)
        check_dependencies
        run_integration_tests
        ;;
    mas)
        check_dependencies
        run_mas_tests
        ;;
    resolver)
        check_dependencies
        run_resolver_tests
        ;;
    registry)
        check_dependencies
        run_registry_tests
        ;;
    edge-cases)
        check_dependencies
        run_edge_case_tests
        ;;
    race)
        check_dependencies
        run_race_tests
        ;;
    bench)
        check_dependencies
        run_benchmarks
        ;;
    coverage)
        check_dependencies
        run_coverage
        ;;
    comprehensive)
        run_comprehensive_test
        ;;
    validate)
        validate_mas_implementation
        ;;
    help|*)
        show_help
        ;;
esac