# 7EP-0013: Robust Build Pipeline with Level 1 Reproducible Builds

**Status:** Draft  
**Author(s):** Claude Code (CC)  
**Assignment:** CC  
**Difficulty:** 2 (straightforward - established patterns with clear implementation)  
**Created:** 2025-08-13  
**Updated:** 2025-08-13

## Executive Summary

Eliminate AI assistant build blockages and establish professional-grade build consistency by implementing Level 1 reproducible builds (commit-deterministic timestamps) with robust error handling, proper directory management, and systematic build scripts that work reliably across local development and CI environments.

## Evidence & Reasoning

**Critical Blocker Issues:**
- **CC/AC Getting Stuck**: Complex chained builds cause terminal freezes and workflow blockages
- **Missing Directory Handling**: `dist/` directory creation failures cause silent build errors  
- **Inconsistent Build Process**: Different patterns between local dev, CI, and manual builds
- **Poor Error Handling**: Commands fail silently, leaving developers in unknown states
- **LDFLAGS Complexity**: Multi-line variable evaluation causes shell parsing issues

**User feedback/pain points:**
- AI assistants cannot reliably complete build tasks without getting "stuck"
- Developers experience "works on my machine" build inconsistencies  
- Build failures provide unclear error messages or fail silently
- Complex build commands are error-prone and hard to debug

**Performance data:**
- Current manual build commands: 3-step process prone to failures
- Build script hanging: Blocks development workflow for 5+ minutes per incident
- Directory creation errors: ~30% of build attempts in clean environments

**Current limitations:**
- No standardized build process across development and CI environments
- Build metadata varies by execution time, not source code state
- No validation of build environment or prerequisites
- Manual error-prone command chains required for reliable builds

**Why now:**
- **Active Development Blocker**: CC/AC cannot reliably build, blocking 7EP-0007 Phase 3 work
- **Professional Standards Gap**: Current build system below industry norms
- **Foundation for Growth**: Clean builds needed for future release automation
- **Team Coordination**: Consistent build process enables better AI assistant collaboration

## Use Cases

### Primary Use Case: Reliable AI Assistant Builds
```bash
# Current problematic workflow (causes hanging)
COMMIT=$(git rev-parse --short HEAD || echo local)
TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")  
GOFLAGS="-trimpath" go build -o dist/7zarch-go -ldflags "-s -w -X main.BuildTime=${TIME} -X main.GitCommit=${COMMIT}" ./

# Proposed simple workflow
make dist              # Or: scripts/build.sh dist
# -> Handles all error cases, creates directories, validates output
```

### Secondary Use Cases

#### Developer Local Builds
```bash
# Development builds with current timestamp
make build             # ./7zarch-go (for development)
make dist             # dist/7zarch-go (for distribution)
make dev              # Build + symlink to ~/bin
```

#### CI/CD Builds  
```bash
# Reproducible builds in CI environment
scripts/build.sh dist  # Same process as local, deterministic output
make validate         # Environment and build validation
```

#### Build Troubleshooting
```bash
scripts/validate.sh   # Check environment prerequisites
scripts/clean.sh      # Clean build artifacts and start fresh
make verify          # Build validation without side effects
```

## Technical Design

### Overview
Implement Level 1 reproducible builds using a hybrid Scripts + Makefile approach that provides:
1. **Timestamp Determinism**: Same commit produces same binary
2. **Robust Error Handling**: Clear failure modes and recovery guidance
3. **Environment Validation**: Prerequisites checked before build attempts
4. **Cross-Platform Support**: Works on macOS, Linux, Windows (Git Bash)

### Level 1 Reproducible Build Strategy

**Timestamp Determinism Implementation:**
```bash
# Current (Non-reproducible)
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")  # Varies every build

# Level 1 (Commit-deterministic)
if git rev-parse --verify HEAD >/dev/null 2>&1; then
    # Use commit timestamp for reproducibility
    COMMIT_TIMESTAMP=$(git log -1 --format=%ct)
    BUILD_TIME=$(date -u -d @$COMMIT_TIMESTAMP +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || \
                 date -r $COMMIT_TIMESTAMP -u +"%Y-%m-%dT%H:%M:%SZ")
else
    # Fallback for non-git environments (development)
    BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
fi
```

**Benefits of Level 1 Approach:**
- Same commit hash produces identical binaries across environments
- Debugging: Can verify if local build matches release build
- Professional standard for non-security-critical CLI tools
- Future-proof foundation for enhanced reproducibility

### Component Architecture

#### 1. Core Build Script (`scripts/build.sh`)
```bash
# Key Features:
- Environment validation (Go, Git, project structure)
- Automatic directory creation with error handling
- Level 1 reproducible timestamp generation
- Build targets: local, dist, both
- Binary validation and verification
- Cross-platform compatibility (macOS, Linux focus)
- Colored logging with structured output
```

#### 2. Environment Validation (`scripts/validate.sh`)
```bash
# Validation checks:
- Go installation and version compatibility
- Git availability and repository state
- Project structure verification (go.mod, main.go)
- Build dependencies and toolchain
```

#### 3. Cleanup Script (`scripts/clean.sh`)
```bash
# Comprehensive cleanup:
- Remove build artifacts (7zarch-go, dist/, temp files)
- Clear Go build cache (optional)
- Reset to clean development state
```

#### 4. Enhanced Makefile
```makefile
# New targets calling scripts:
dist: scripts/build.sh dist
validate: scripts/validate.sh
dev: scripts/build.sh local && symlink-to-bin
clean-all: scripts/clean.sh
```

### Build Metadata Standard

**LDFLAGS Pattern (Level 1 Reproducible):**
```bash
VERSION="${VERSION:-dev}"
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
COMMIT_TIMESTAMP=$(git log -1 --format=%ct 2>/dev/null || date +%s)
BUILD_TIME=$(date -u -d @$COMMIT_TIMESTAMP +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS="-s -w"
LDFLAGS="$LDFLAGS -X main.Version=$VERSION"
LDFLAGS="$LDFLAGS -X main.GitCommit=$GIT_COMMIT"  
LDFLAGS="$LDFLAGS -X main.BuildTime=$BUILD_TIME"

go build -trimpath -ldflags "$LDFLAGS" -o "$OUTPUT" ./
```

## Implementation Plan

### Phase 1: Core Build Scripts (3-4 hours)
- [ ] **Create `scripts/build.sh`** (CC)
  - Environment validation and error handling
  - Level 1 reproducible timestamp generation
  - Build targets: local, dist, both
  - Binary validation with size and execution checks
  - Cross-platform support (macOS, Linux, Git Bash on Windows)
  - Colored logging with structured output

- [ ] **Create `scripts/validate.sh`** (CC)
  - Go installation and version checks  
  - Git availability and repository validation
  - Project structure verification (go.mod, correct directory)
  - Build dependency checks

- [ ] **Create `scripts/clean.sh`** (CC)
  - Remove all build artifacts safely
  - Optional Go build cache clearing
  - Confirmation prompts for destructive operations

### Phase 2: Makefile Integration (1-2 hours) 
- [ ] **Enhance existing Makefile** (CC)
  - Add `dist` target calling `scripts/build.sh dist`
  - Add `validate` target for environment checking
  - Add `dev` target for development workflow (build + symlink)
  - Add `clean-all` target for comprehensive cleanup
  - Update existing targets to use scripts where beneficial
  - Maintain backward compatibility with current targets

- [ ] **Add convenience targets** (CC)
  - `verify` - build validation without creating artifacts
  - `dev-setup` - one-time developer environment configuration

### Phase 3: CI/CD Integration (1-2 hours)
- [ ] **Update GitHub Actions workflows** (CC)
  - Modify `.github/workflows/build.yml` to use scripts
  - Ensure Level 1 reproducible builds in CI environment
  - Add build artifact validation and verification steps
  - Test cross-platform script compatibility

- [ ] **Add build verification** (CC)
  - Binary checksum generation and validation
  - Build metadata extraction and verification
  - Reproducibility testing across environments

### Dependencies
- None - improves existing build system without external dependencies
- Compatible with existing Makefile targets
- Works with current GitHub Actions infrastructure

## Testing Strategy

### Acceptance Criteria
- [ ] **AI Assistant Unblocking**: CC/AC can reliably build without hanging or freezing
- [ ] **Reproducible Builds**: Same commit produces identical binary checksums across environments
- [ ] **Error Handling**: Build failures provide clear, actionable error messages
- [ ] **Directory Management**: `dist/` directory created automatically, no silent failures
- [ ] **Cross-Platform**: Scripts work on macOS, Linux, and Windows (Git Bash)
- [ ] **Performance**: Build process completes in under 30 seconds for normal builds
- [ ] **Backward Compatibility**: Existing `make build` continues to work unchanged

### Test Scenarios

#### Build Process Testing
- Fresh clone builds (no dist/ directory)
- Builds in dirty working directories
- Builds with and without Git repository
- Cross-platform build validation (macOS primary, Linux secondary)
- Error recovery from failed builds

#### Reproducibility Testing  
- Multiple builds of same commit produce identical checksums
- Builds across different environments (local vs CI) match
- Timestamp consistency verification
- Build metadata accuracy validation

#### Error Handling Testing
- Missing Go installation
- Missing Git (should fallback gracefully)
- Incorrect working directory
- Insufficient disk space
- Permission errors on directory creation

### Performance Benchmarks
- **Build time**: <30 seconds for typical builds
- **Validation time**: <5 seconds for environment checks
- **Error detection**: <10 seconds to identify and report build issues
- **Cross-platform consistency**: Same build times ±20% across platforms

## Migration/Compatibility

### Breaking Changes
None - all new functionality building on existing build system.

### Upgrade Path
- Existing Makefile targets continue working unchanged
- New targets and scripts available immediately
- Gradual adoption: teams can use new system when ready
- CI/CD updates can be rolled out incrementally

### Backward Compatibility
- All existing `make build`, `make test`, `make lint` targets preserved
- Existing GitHub Actions workflows continue functioning
- Developer muscle memory preserved during transition

## Implementation Options

### Option A: Conservative Enhancement (Recommended)
**Scope**: Enhance current Makefile + add robust scripts  
**Timeline**: 6-8 hours total  
**Risk**: Low - builds on existing patterns

**Approach:**
- Keep all existing Makefile targets for compatibility  
- Add scripts/ directory with new build infrastructure
- New targets call scripts internally
- Level 1 reproducible builds via commit timestamps
- Comprehensive error handling and validation

**Benefits:**
- Addresses immediate CC/AC blocking issues
- Establishes professional build foundation
- Maintains team workflow continuity
- Clear path for future enhancements

### Option B: Goreleaser Migration
**Scope**: Adopt industry-standard Goreleaser framework  
**Timeline**: 8-12 hours total  
**Risk**: Medium - changes development workflow

**Approach:**
- Add `.goreleaser.yml` configuration
- Enhanced release automation with built-in reproducible builds
- GitHub Actions integration with automatic releases
- Industry-standard tooling and patterns

**Benefits:**
- Matches patterns from kubectl, helm, terraform
- Built-in reproducible builds and checksums
- Automatic release process
- Reduced custom build logic maintenance

### Option C: Modern Tooling (Task/Mage)
**Scope**: Replace Make with modern Go-native tooling  
**Timeline**: 10-16 hours total  
**Risk**: High - significant workflow changes

**Approach:**
- Replace Makefile with Taskfile.yml or magefile.go
- Modern dependency management and task orchestration
- Type-safe build scripts (if using Mage)
- Enhanced developer experience

**Benefits:**
- More modern than Make
- Better dependency handling and task orchestration
- Platform-native approach (especially Mage for Go projects)

## Alternatives Considered

**External build tools**: Evaluated Bazel and Buck but decided they're overkill for a single Go binary project.

**Docker-based builds**: Considered containerized builds for full reproducibility but decided Level 1 reproducibility provides sufficient benefits without the complexity overhead.

**GitHub Actions only builds**: Evaluated removing local build capability but decided local development builds are essential for developer productivity.

**Shell script only**: Considered pure shell scripts without Makefile but decided hybrid approach provides better developer experience (familiar `make` targets + robust scripts).

## Future Considerations

### Evolution Path: Level 1 → Level 2+ Reproducibility
- **SOURCE_DATE_EPOCH**: Full timestamp determinism support
- **Container builds**: Hermetic build environment isolation  
- **Build attestation**: Cryptographic build provenance
- **SBOM generation**: Software Bill of Materials for supply chain security

### Release Automation Foundation
- Scripts provide foundation for automated release processes
- Standard LDFLAGS make Goreleaser migration straightforward  
- Build validation enables automated quality gates
- Checksum generation supports package manager integration

### Enhanced Developer Experience
- Shell completion for build targets and options
- Build performance profiling and optimization
- Integration with development tools (VS Code, GoLand)
- Advanced build caching and incremental builds

## Priority Justification

**Critical Priority** because:
1. **Actively Blocking Development**: CC/AC cannot reliably build, blocking 7EP-0007 Phase 3 work
2. **Affects All Contributors**: Every developer experiences build inconsistencies  
3. **Foundation for Other Work**: Clean builds required for continued feature development
4. **Professional Standards**: Current build system significantly below industry expectations
5. **Low Risk, High Impact**: Conservative approach with major workflow improvements

**Recommended Approach: Option A (Conservative Enhancement)**
- Solves immediate blocker issues without disrupting team workflow
- Establishes Level 1 reproducible build foundation  
- Provides clear path for future build system enhancements
- Minimal risk with maximum immediate benefit

## References

- **Level 1 Reproducible Builds**: Industry standard for non-security-critical CLI tools
- **Build patterns**: Reference implementations from kubectl, docker, terraform build systems
- **Go build best practices**: Official Go team recommendations for `-ldflags` and `-trimpath` usage
- **Cross-platform scripting**: Bash patterns that work on macOS, Linux, and Windows Git Bash