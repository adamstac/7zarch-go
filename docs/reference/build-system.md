# Build & Development Environment

## Makefile Targets

```bash
make build          # Build main binary
make test           # Run unit tests
make integration    # Integration tests
make lint           # Code linting (go vet + revive) and formatting
make build-all      # Multi-platform builds
```

## Code Quality & Linting

**Current Linter: revive** (replaced golangci-lint due to CI module resolution issues)

```bash
# Local linting (same as CI)
make lint                           # Run go vet + gofmt check
go install github.com/mgechev/revive@latest
~/go/bin/revive -config revive.toml -formatter friendly ./...
```

**Why revive instead of golangci-lint:**
- **Module Resolution**: golangci-lint has persistent module resolution issues in CI environments with `yaml` and `progressbar` imports
- **Reliability**: revive works consistently in both local and CI environments  
- **Performance**: Faster than golangci-lint with comparable code quality feedback
- **Maintainability**: Simpler configuration, fewer CI environment issues

**Configuration:** `revive.toml` provides reasonable defaults with warnings-only output (non-blocking)

## User Installation Pattern

```bash
# Development symlink approach
ln -sf $(pwd)/7zarch-go /usr/local/bin/7zarch-go
```

## Build Pipeline

**Note:** As of 7EP-0013, the project uses Goreleaser for professional builds with Level 2 reproducibility. See:
- `make dev` - Build with Goreleaser and install to ~/bin
- `make dist` - Build for current platform  
- `make validate` - Validate Goreleaser config
- `make release` - Create release (CI only)

The legacy build system (`make build`) is still available but Goreleaser is preferred for all builds.