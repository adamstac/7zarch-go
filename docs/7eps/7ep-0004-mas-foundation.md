# 7EP-0004: MAS Foundation Implementation

**Status:** In Progress  
**Author(s):** Claude Code (CC)  
**Assignment:** AC (Primary), CC (Supporting)  
**Difficulty:** 4 (complex - foundational system with multiple interdependent components)  
**Created:** 2025-08-12  
**Updated:** 2025-08-12  
**PRs:** [#5 (AC)](https://github.com/adamstac/7zarch-go/pull/5), [#6 (CC)](https://github.com/adamstac/7zarch-go/pull/6)  

## Executive Summary

Implement the core MAS (Managed Archive Storage) foundation consisting of ULID resolution, show command, and enhanced list functionality to enable stable ID-based archive operations.

## Evidence & Reasoning

**User feedback/pain points:**
- Users want to reference archives by short IDs instead of full paths
- Need quick way to inspect archive details and verify integrity
- Current list command lacks filtering and status information
- Registry operations feel disconnected from file system reality

**Current limitations:**
- No ULID resolution system for user-friendly references
- No show command to validate registry data against files
- Basic list command without filtering or status indicators
- Inconsistent error handling across registry operations

**Why now:**
- Registry infrastructure exists but lacks user-facing operations
- Foundation needed before trash management and other advanced features
- Users need confidence in registry accuracy and file tracking
- Clear API patterns required for future command development

## Use Cases

### Primary Use Case: ID-Based Archive Management
```bash
# User creates archive and gets ULID
7zarch-go create project-backup
# Output: Created archive with ID: 01K2E33XW4

# User references by short ID instead of path
7zarch-go show 01K2E33
7zarch-go delete 01K2E33
7zarch-go move 01K2E33 --to /backup/

# Disambiguation when prefix is ambiguous
7zarch-go show 01K2
# Multiple matches found for '01K2':
# [1] 01K2E33 project-backup (managed, 2.1 MB, 2d ago)
# [2] 01K2F44 project-docs (external, 0.5 MB, 1w ago)
# Please specify: 01K2E33, 01K2F44, or full name
```

### Secondary Use Cases
- **Registry validation**: Show command verifies files exist and match checksums
- **Archive discovery**: Enhanced list with filtering by status, location, profile
- **Troubleshooting**: Clear error messages with suggested resolutions
- **Batch operations**: List filters enable bulk operations on specific archive sets

## Technical Design

### Overview
Build user-facing MAS operations on existing registry infrastructure with consistent patterns for resolution, validation, and display.

### Component Architecture

#### 1. ULID Resolution System
```go
// internal/storage/resolver.go
type Resolver struct {
    registry *Registry
}

// Resolution priority order
func (r *Resolver) ResolveID(input string) (*Archive, error) {
    // 1. Exact ULID match (fastest path)
    if archive := r.getByUID(input); archive != nil {
        return archive, nil
    }
    
    // 2. ULID prefix (most common use case)
    matches := r.getByUIDPrefix(input)
    if len(matches) == 1 {
        return matches[0], nil
    } else if len(matches) > 1 {
        return nil, &AmbiguousIDError{ID: input, Matches: matches}
    }
    
    // 3. Checksum prefix
    matches = r.getByChecksumPrefix(input)
    if len(matches) == 1 {
        return matches[0], nil
    } else if len(matches) > 1 {
        return nil, &AmbiguousIDError{ID: input, Matches: matches}
    }
    
    // 4. Name exact match
    if archive := r.getByName(input); archive != nil {
        return archive, nil
    }
    
    return nil, &ArchiveNotFoundError{ID: input}
}

// Interactive disambiguation
func (r *Resolver) HandleAmbiguous(err *AmbiguousIDError) (*Archive, error) {
    // Display options with context
    // Prompt for selection
    // Return selected archive
}
```

#### 2. Show Command Implementation
```go
// cmd/mas_show.go
func runMasShow(cmd *cobra.Command, args []string) error {
    resolver := storage.NewResolver(registry)
    
    archive, err := resolver.ResolveID(args[0])
    if err != nil {
        return handleResolutionError(err)
    }
    
    // Verify file existence and integrity
    status := verifyArchiveStatus(archive)
    
    // Display comprehensive information
    displayArchiveDetails(archive, status)
    
    return nil
}

// Archive status verification
type ArchiveStatus struct {
    FileExists    bool
    ChecksumValid bool
    LastVerified  time.Time
    Issues        []string
}
```

#### 3. Enhanced List Command
```go
// cmd/list.go enhanced filters
type ListFilters struct {
    Status      string   // present, missing, deleted
    Managed     *bool    // true, false, nil (all)
    Profile     string   // media, documents, balanced
    Pattern     string   // glob pattern for names
    OlderThan   time.Duration
    LargerThan  int64
    NotUploaded bool
}

func applyFilters(archives []*storage.Archive, filters ListFilters) []*storage.Archive {
    // Apply each filter sequentially
    // Return filtered results
}
```

### Error Handling Standards

#### Standard Error Types
```go
type ArchiveNotFoundError struct {
    ID string
}

type AmbiguousIDError struct {
    ID      string
    Matches []*Archive
}

type RegistryError struct {
    Operation string
    Cause     error
}

type FileVerificationError struct {
    Archive *Archive
    Issue   string
}
```

#### User-Friendly Messages
```go
func (e *ArchiveNotFoundError) Error() string {
    return fmt.Sprintf("Archive '%s' not found.\n" +
        "💡 Use '7zarch-go list' to see available archives", e.ID)
}

func (e *AmbiguousIDError) Error() string {
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("Multiple archives match '%s':\n", e.ID))
    for i, archive := range e.Matches {
        sb.WriteString(fmt.Sprintf("[%d] %s %s (%s, %.1f MB, %s)\n",
            i+1, archive.UID[:8], archive.Name,
            archiveLocation(archive), 
            float64(archive.Size)/(1024*1024),
            humanizeTime(archive.Created)))
    }
    sb.WriteString("Please specify full ULID or use a longer prefix")
    return sb.String()
}
```

### Display Standards

#### Show Command Output Format
```
Archive: project-backup.7z (01K2E33XW4HTX7RVPS9Y6CRGDY)
Location: ~/.7zarch-go/archives/project-backup.7z (managed)
Status: present ✓
Size: 2.1 MB (compressed from 8.7 MB, 75.9% reduction)
Profile: documents
Created: 2025-08-11 14:30:22 (2 days ago)
Checksum: sha256:a1b2c3d4... (verified ✓)
Upload Status: not uploaded

💡 Use '7zarch-go delete 01K2E33' to move to trash
💡 Use '7zarch-go move 01K2E33 --to /backup/' to relocate
```

#### List Command Output Format
```
📦 Archives (15 found, 45.2 GB total)

MANAGED STORAGE (~/.7zarch-go/archives/):
01K2E33  project-backup.7z     2.1 MB   documents  2d ago   ✓
01K2F44  podcast-103.7z       156 MB   media      1w ago   ✓
01K2G55  code-dump.7z          0.8 MB   documents  1w ago   ⚠️ missing

EXTERNAL STORAGE:
01K2H66  /backup/old-site.7z   45 MB    balanced   1m ago   ✓
01K2I77  ~/Desktop/temp.7z      2 MB    media      3d ago   ✓

💡 Use '7zarch-go show <id>' for details
💡 Use '7zarch-go list --help' for filter options
```

## Implementation Plan

### Phase 1: Core Infrastructure (AC Primary)
- [ ] **ULID Resolution System** (AC)
  - [ ] Implement core resolver with prefix matching
  - [ ] Add disambiguation interface for multiple matches
  - [ ] Create standard error types with helpful messages
  - [ ] Add resolution performance tests (target <50ms)

- [ ] **Show Command Base** (AC)
  - [ ] Basic show command with resolver integration
  - [ ] File existence verification
  - [ ] Checksum validation on demand
  - [ ] Status indicator display (✓/❌/⚠️)

- [ ] **Registry Query Optimization** (AC)
  - [ ] Add database indexes for ULID prefix queries
  - [ ] Implement efficient checksum prefix matching
  - [ ] Optimize name-based lookups
  - [ ] Performance testing with 1000+ archive datasets

### Phase 2: Enhanced Operations (AC Primary)
- [ ] **Enhanced List Command** (AC)
  - [ ] Implement comprehensive filtering system
  - [ ] Add status-based grouping (managed/external/deleted)
  - [ ] Create tabular output with consistent formatting
  - [ ] Add summary statistics and helpful tips

- [ ] **Advanced Show Features** (AC)
  - [ ] Detailed metadata display
  - [ ] File integrity verification
  - [ ] Location-specific information
  - [ ] Suggested actions based on status

### Phase 3: Polish & Testing (CC Supporting)
- [x] **Error Handling Standardization** (CC) - PR #6
  - [x] Consistent error message format across commands
  - [x] Recovery suggestions for common issues
  - [x] Help text improvements
  - [ ] Error message user testing

- [x] **Test Infrastructure** (CC) - PR #6
  - [x] Test helpers for registry creation
  - [x] Archive creation utilities with options
  - [x] Assertion helpers for resolver testing
  - [ ] Cross-platform compatibility verification

- [x] **Documentation Updates** (CC) - PR #6
  - [x] Show command reference documentation
  - [x] Troubleshooting guides
  - [ ] List command documentation
  - [ ] Migration documentation

### Dependencies
- Existing registry infrastructure (implemented)
- ULID generation system (implemented)
- Database schema with indexes (needs optimization)

## Testing Strategy

### Acceptance Criteria
- [ ] Can resolve archives by ULID, ULID prefix, checksum prefix, and name
- [ ] Disambiguation works intuitively for ambiguous inputs
- [ ] Show command displays accurate, helpful information
- [ ] List command supports all documented filters
- [ ] Operations complete in <100ms for typical registries (<1000 archives)
- [ ] Error messages are actionable and helpful
- [ ] File verification detects missing/corrupted archives

### Test Scenarios

#### Resolution Testing
- Empty registry scenarios
- Single character prefixes (high ambiguity)
- Checksum collisions
- Unicode names and special characters
- Large registry performance (10,000+ archives)

#### Show Command Testing
- Missing files (moved/deleted outside 7zarch-go)
- Corrupted archives (checksum mismatch)
- Network storage latency
- Permission errors

#### List Command Testing
- Complex filter combinations
- Large result sets with pagination
- Mixed managed/external archives
- Status transitions during operation

### Performance Benchmarks
- **Resolution operations**: <50ms for any query type
- **Show command**: <100ms including file verification
- **List operations**: <200ms for 1000+ archives with filters
- **Memory usage**: <10MB for typical registry operations

## Migration/Compatibility

### Breaking Changes
None - all new functionality building on existing registry.

### Upgrade Path
Automatic - existing registries work immediately with new commands.

### Backward Compatibility
All existing commands continue working unchanged.

## Alternatives Considered

**Separate resolution library**: Considered extracting resolver to separate package but decided inline implementation reduces complexity for initial version.

**GraphQL-style query interface**: Evaluated complex query syntax but decided simple flags provide better CLI experience.

**Fuzzy matching for names**: Considered Levenshtein distance for typos but decided explicit disambiguation is clearer.

## Implementation Notes

### Key Design Decisions (Learned During Implementation)

#### Error Message Philosophy
The error types implemented follow a user-first approach:
- **Context First**: Tell user what went wrong in their terms, not technical terms
- **Suggestions Always**: Every error includes actionable next steps
- **Visual Indicators**: Use emojis sparingly but effectively (💡 for tips)
- **Progressive Detail**: Simple message first, detailed help available

Example implemented:
```go
func (e *ArchiveNotFoundError) Error() string {
    return fmt.Sprintf("Archive '%s' not found.\n💡 Use '7zarch-go list' to see available archives", e.ID)
}
```

#### Test Infrastructure Design
Test helpers focus on builder patterns for flexibility:
- **Registry Creation**: In-memory SQLite for fast tests
- **Archive Builders**: Functional options pattern for readable test setup
- **Assertion Helpers**: Domain-specific assertions reduce boilerplate

Key insight: Tests should read like specifications:
```go
archive := CreateTestArchive(t, reg, "test.7z", 
    WithSize(2*MB), 
    WithProfile("documents"),
    WithStatus("deleted"))
AssertResolves(t, resolver, "test", archive)
```

#### Resolution Priority Insights
Through test design, the optimal resolution order became clear:
1. **Exact ULID** - Fastest path, most specific
2. **ULID Prefix** - Most common user interaction (copy first 8 chars)
3. **Checksum Prefix** - Power user feature for deduplication
4. **Name Match** - Fallback for human-friendly access

This priority prevents name collisions from breaking ULID resolution.

### Coordination Patterns That Worked

#### Branch Naming Convention
- **AC Branch**: `feature/7ep-0004-mas-foundation`
- **CC Branch**: `cc/7ep-0004-support`

Prefix by role prevents conflicts and clarifies ownership.

#### Task Separation
Clear boundaries in 7EP prevented toe-stepping:
- **AC**: Core business logic (resolver, commands)
- **CC**: Infrastructure (errors, testing, docs)

No file conflicts, clear ownership, parallel development.

#### Cross-PR Communication
- Each PR references the 7EP number
- PRs cross-link in descriptions
- Comments notify of dependencies
- Clear "this provides X for Y" messaging

### Performance Considerations Discovered

#### Registry Query Optimization
Testing revealed key optimization points:
- **Index on uid prefix**: Critical for ULID resolution
- **Index on checksum prefix**: Enables fast deduplication
- **Name index**: Already exists, just needs case handling

#### Memory Management
Test helpers revealed memory patterns:
- **Batch operations**: Need streaming/pagination for large registries
- **Error messages**: Avoid loading all matches for ambiguous errors
- **Test cleanup**: Proper cleanup prevents test database accumulation

### Documentation Insights

#### Show Command Documentation Structure
Most effective documentation pattern:
1. **Quick examples first** - Get users successful fast
2. **Comprehensive flags table** - Reference when needed
3. **Output format examples** - Show don't just tell
4. **Troubleshooting section** - Anticipate problems
5. **Script integration** - Power user examples

#### Error Message Documentation
Users need to see actual error messages in docs:
- Include exact error text they'll encounter
- Show the resolution steps
- Explain why the error occurred

## Future Considerations

- **Full-text search**: Search across all archive metadata fields
- **Saved searches**: Store complex filter combinations
- **Shell completion**: Auto-complete for ULID prefixes
- **Batch operations**: Apply operations to filtered archive sets
- **Interactive disambiguation**: Terminal UI for selecting from ambiguous matches

## References

- Related: Existing registry infrastructure in internal/storage/
- Related: 7EP-0001 Trash Management (depends on show command patterns)
- Related: 7EP-0003 Database Migrations (performance optimization dependency)
- PRs: [#5 (AC implementation)](https://github.com/adamstac/7zarch-go/pull/5), [#6 (CC support)](https://github.com/adamstac/7zarch-go/pull/6)