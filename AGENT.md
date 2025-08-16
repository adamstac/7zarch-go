# 7zarch-go Agent Guide

## Build/Test Commands
- **Build**: `make dev` (builds and installs to ~/bin) or `make build` (local only)
- **Test**: `go test ./...` (all tests), `go test -v ./internal/storage/` (single package)
- **Lint**: `go vet ./...` (basic) or `golangci-lint run` (full, config in .golangci.yml)
- **Coverage**: `make test-coverage-html` (generates coverage.html)
- **Single test**: `go test -run TestFunctionName ./path/to/package/`

## DDD Framework Validation Commands (7EP-0020)
- **Complete validation**: `make validate-framework` (all validation systems)
- **Role files only**: `make validate-framework-roles` (role standardization)
- **Cross-document consistency**: `make validate-framework-consistency` (coordination sync)
- **Git patterns**: `make validate-framework-git` (session logs, commits, branches)
- **Integration testing**: `make validate-framework-integration` (end-to-end lifecycle)
- **Framework health**: `make framework-health` (continuous monitoring dashboard)

## Architecture
- **Cobra CLI**: Commands in `cmd/`, main entry in `main.go`, extensive subcommands
- **Database**: SQLite registry (`~/.7zarch-go/registry.db`) with ULID-based archive tracking
- **Core packages**: `internal/storage` (registry/db), `internal/tui` (bubbletea), `internal/batch` (operations)
- **TUI**: Bubbletea framework with 9 themes, multiple entry points (browse/ui/i/tui)

## Code Style (from existing patterns)
- **Imports**: Standard → external → internal (grouped with blank lines)
- **Errors**: Wrap with context (`fmt.Errorf("operation failed: %w", err)`)
- **Functions**: `runCommandName(cmd *cobra.Command, args []string) error` pattern for cobra
- **Types**: PascalCase exported, camelCase unexported, descriptive names
- **Database**: Use prepared statements, proper transactions, defer close pattern
- **Testing**: `_test.go` suffix, table-driven tests preferred, use testify assertions

## Key Conventions
- ULID for user-facing IDs, SQLite auto-increment for internal IDs
- Managed vs external storage distinction (boolean `managed` field)
- Soft deletes with `status` field (present/missing/deleted)
- Profile-based compression (media/documents/balanced/maximum)
