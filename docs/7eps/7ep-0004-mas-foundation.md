# 7EP-0004: MAS Foundation Implementation

**Status:** In Progress (90% complete)  
**Author(s):** Claude Code (CC)  
**Assignment:** AC (Primary), CC (Supporting)  
**Difficulty:** 4 (complex - foundational system with multiple interdependent components)  
**Created:** 2025-08-12  
**Updated:** 2025-08-12

## Current Status (August 12, 2025)

**Implementation Progress: 90% Complete**

**AC's Implementation (PR #5):**
- ✅ Complete ULID resolution system with prefix matching
- ✅ Full show command with file verification and integrity checks
- ✅ Enhanced list command with comprehensive filtering
- ✅ Status-based grouping and tabular output formatting
- ✅ Human-friendly duration and size parsing (`7d`, `100MB`)
- 🔄 Implementing `config.Load` error handling and status validation per CodeRabbit feedback

**CC's Support Infrastructure (PR #6 - Merged):**
- ✅ Standard error types with user-friendly messages
- ✅ Test infrastructure with builder patterns
- ✅ Show and list command documentation with examples
- ✅ Error handling patterns and help text standards

**Next Steps:**
- AC finishing PR #5 based on code review feedback
- Performance testing with large registry datasets
- Integration testing of complete workflow  

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

### Phase 1: Core Infrastructure (AC Primary) - IN PROGRESS
- [x] **ULID Resolution System** (AC - PR #5)
  - [x] Implement core resolver with prefix matching
  - [x] Add disambiguation interface for multiple matches
  - [x] Create standard error types with helpful messages (CC - PR #6)
  - [ ] Add resolution performance tests (target <50ms) - AC implementing

- [x] **Show Command Base** (AC - PR #5)
  - [x] Basic show command with resolver integration
  - [x] File existence verification
  - [x] Checksum validation on demand
  - [x] Status indicator display (✓/❌/⚠️)

- [ ] **Registry Query Optimization** (AC - PR #5)
  - [x] Add database indexes for ULID prefix queries
  - [x] Implement efficient checksum prefix matching
  - [x] Optimize name-based lookups
  - [ ] Performance testing with 1000+ archive datasets - pending

### Phase 2: Enhanced Operations (AC Primary) - IN PROGRESS
- [x] **Enhanced List Command** (AC - PR #5)
  - [x] Implement comprehensive filtering system
  - [x] Add status-based grouping (managed/external/deleted)
  - [x] Create tabular output with consistent formatting
  - [x] Add summary statistics and helpful tips

- [x] **Advanced Show Features** (AC - PR #5)
  - [x] Detailed metadata display
  - [x] File integrity verification
  - [x] Location-specific information
  - [x] Suggested actions based on status

### Phase 3: Polish & Testing (CC Supporting) - IN PROGRESS
- [x] **Error Handling Standardization** (CC - PR #6)
  - [x] Consistent error message format across commands
  - [x] Recovery suggestions for common issues
  - [x] Help text improvements
  - [ ] Error message user testing - pending

- [ ] **Comprehensive Testing** (CC)
  - [x] Test infrastructure with builder patterns (CC - PR #6)
  - [ ] Resolution edge cases (empty registry, corruption) - pending
  - [ ] Cross-platform compatibility - pending
  - [ ] Performance benchmarks - pending
  - [ ] User workflow integration tests - pending

- [x] **Documentation Updates** (CC - PR #6)
  - [x] Command reference updates (show.md, list.md)
  - [x] User workflow examples
  - [x] Troubleshooting guides
  - [ ] Migration documentation - pending

### Dependencies
- Existing registry infrastructure (implemented)
- ULID generation system (implemented)
- Database schema with indexes (needs optimization)

## Testing Strategy

### Acceptance Criteria
- [x] Can resolve archives by ULID, ULID prefix, checksum prefix, and name
- [x] Disambiguation works intuitively for ambiguous inputs
- [x] Show command displays accurate, helpful information
- [x] List command supports all documented filters
- [ ] Operations complete in <100ms for typical registries (<1000 archives) - performance testing pending
- [x] Error messages are actionable and helpful
- [x] File verification detects missing/corrupted archives

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

## Future Considerations

- **Full-text search**: Search across all archive metadata fields
- **Saved searches**: Store complex filter combinations
- **Shell completion**: Auto-complete for ULID prefixes
- **Batch operations**: Apply operations to filtered archive sets

## References

- Related: Existing registry infrastructure in internal/storage/
- Related: 7EP-0001 Trash Management (depends on show command patterns)
- Related: 7EP-0003 Database Migrations (performance optimization dependency)