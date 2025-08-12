# 7EP-0001: Trash Management System

**Status:** Draft  
**Author(s):** Adam Stacoviak, Claude Code (CC)  
**Assignment:** AC  
**Difficulty:** 3 (moderate - new commands + workflow integration)  
**Created:** 2025-08-11  
**Updated:** 2025-08-11  
**Revised:** 2025-08-12 (post MAS Foundation completion)  

## Executive Summary

Implement comprehensive trash management for deleted archives including restore functionality, permanent purge operations, and trash-specific listing commands to complete the archive deletion lifecycle.

## Evidence & Reasoning

**User feedback/pain points:**
- Current `delete` command moves files to trash but provides no way to restore accidentally deleted archives
- No visibility into what's in trash or when items will be auto-purged
- Users cannot manually clean up trash before auto-purge period expires

**Current limitations:**
- Soft-deleted archives appear in regular `list` output without clear deleted status (addressed in recent updates)
- No `restore` command to undo deletions
- No manual purge capability for immediate cleanup
- No trash-specific operations

**Why now (Post 7EP-0004):**
- **MAS Foundation Complete**: ULID resolution and registry operations provide solid foundation
- **Delete workflow incomplete**: Current delete moves to trash but no restore capability
- **User workflow gap**: Accidental deletions have no recovery path
- **Auto-purge needs control**: Manual purge counterpart required for user control
- **Strategic timing**: Build on proven MAS patterns while they're fresh

## Use Cases

### Primary Use Case: Accidental Deletion Recovery
```bash
# User accidentally deletes archive
7zarch-go delete important-backup.7z

# User realizes mistake and restores
7zarch-go restore important-backup.7z
# OR
7zarch-go restore 01K2E33XW4HTX7RVPS9Y6CRGDY
```

### Secondary Use Cases
- **Trash inspection**: `7zarch-go trash list` to see what's deleted
- **Manual cleanup**: `7zarch-go trash purge` before auto-purge period
- **Selective purge**: `7zarch-go trash purge <id>` for specific items
- **Bulk operations**: Restore or purge multiple items

## Technical Design

### Overview
Build on 7EP-0004's proven MAS Foundation patterns to implement a safe, reviewable, and automatable trash lifecycle for soft-deleted archives with intelligent retention policies and user control.

### Building on MAS Foundation (7EP-0004)
**Leveraging Completed Infrastructure:**
- ✅ **ULID Resolution**: All trash commands will use established resolver patterns
- ✅ **Registry Operations**: Proven `reg.Add()`, `reg.Update()` patterns
- ✅ **Error Handling**: Standard error types with user-friendly messages
- ✅ **List Filtering**: Extend existing list filtering for trash-specific views
- ✅ **Show Command**: Trash items will work with existing show functionality

**Implementation Strategy:**
- **Consistent API**: Follow 7EP-0004's command patterns and flag conventions
- **Resolver Integration**: All trash operations accept ULID, ULID prefix, checksum prefix, name
- **Error Messages**: Use established error types (`ArchiveNotFoundError`, etc.)
- **Performance**: Build on validated O(1) database operations

### New Commands

#### `7zarch-go restore <id>`
**Implementation using MAS Foundation patterns:**
```go
func runMasRestore(cmd *cobra.Command, args []string) error {
    resolver := storage.NewResolver(registry) // 7EP-0004 pattern
    
    // Resolve archive (supports all MAS Foundation identifier types)
    archive, err := resolver.Resolve(args[0])
    if err != nil {
        return handleResolutionError(err) // 7EP-0004 error handling
    }
    
    // Validation: Only restore deleted archives
    if archive.Status != "deleted" {
        return &storage.InvalidOperationError{
            Operation: "restore",
            Archive:   archive,
            Reason:    "Archive is not deleted",
        }
    }
    
    // Move file from trash back to original location
    return restoreArchiveFile(archive)
}
```

**Registry Updates Pattern:**
- `status` → `"present"`
- `path` → value from `original_path` 
- `deleted_at` → `NULL`
- Preserve `original_path` for audit trail

#### `7zarch-go trash list [--details]`
**Extends 7EP-0004 list filtering:**
```go  
func runMasTrashList(cmd *cobra.Command, args []string) error {
    // Use existing list filtering with trash-specific filter
    filter := storage.ListFilters{
        Status: "deleted", // Only deleted archives
    }
    
    archives, err := registry.ListWithFilters(filter) // 7EP-0004 pattern
    if err != nil {
        return err
    }
    
    // Use enhanced display from 7EP-0004 with trash-specific formatting
    return displayTrashList(archives, options.Details)
}
```

**Display Features:**
- Reuses 7EP-0004's tabular formatting patterns
- Shows deletion date, original location, purge countdown
- Status indicators: 🗑️ for deleted items

#### `7zarch-go trash purge [<id>] [--older-than <duration>]`
**Builds on resolver and error handling patterns:**
```go
func runMasTrashPurge(cmd *cobra.Command, args []string) error {
    var archives []*storage.Archive
    var err error
    
    if len(args) > 0 {
        // Specific archive - use resolver
        resolver := storage.NewResolver(registry)
        archive, err := resolver.Resolve(args[0])
        if err != nil {
            return handleResolutionError(err)
        }
        archives = []*storage.Archive{archive}
    } else {
        // Age-based bulk purge - use filtering
        cutoff := time.Now().Add(-options.OlderThan)
        archives, err = findArchivesOlderThan(cutoff, "deleted")
    }
    
    return executePurge(archives, options)
}
```

**Safety Features:**
- Require `--force` for permanent deletion (follows 7EP-0004 safety patterns)
- Clear confirmation prompts with archive details
- Dry-run support for bulk operations

### API Changes
```bash
# New top-level commands
7zarch-go restore <id>
7zarch-go trash list [--details]
7zarch-go trash purge [<id>] [--older-than <duration>] [--force]
```

### Trash Behavior Model

#### Soft Delete (Default)
- Move file to `~/.7zarch-go/trash/` (preserve filename; optional subfolders by date)
- Registry: `status=deleted`, `managed` unchanged, record `original_path` in metadata
- Record `deleted_at` timestamp for retention logic

#### Force Delete (`--force`)
- Physically remove file and set `status=deleted`; record `deleted_at`

#### External Files (Non-managed)
- Default to DB-only delete (no file move) unless `--force` provided
- Still record `deleted_at` for retention tracking

#### Restore Operation
- `7zarch-go restore <id>` moves trashed file back to its `original_path`
- If original location invalid, fallback to MAS
- Update `status=present`

### Configuration Options
```yaml
trash:
  retention_days: 30          # default auto-purge horizon
  max_size_gb: 10             # optional overall trash cap (oldest-first purge)
  auto_purge_on_start: true   # run quick purge on CLI start
  layout: by_date             # flat | by_date (YYYY/MM)
  confirm_purge: true         # ask before large purges
```

### Auto-Purge Policy
- Triggered on startup if `auto_purge_on_start` is true and budget exceeded:
  - If `retention_days` elapsed → eligible
  - If trash exceeds `max_size_gb` → purge oldest first until within budget
- Always support `--dry-run` and clear logs of what would be removed

### Data Model Changes
Leverages existing `status`, `deleted_at`, `original_path` fields added in recent schema migration.

## Implementation Plan (Building on MAS Foundation)

### Phase 1: Core Commands (Estimated: 4-6 hours)
- [ ] **Restore Command Implementation**
  - [ ] Create `cmd/mas_restore.go` following 7EP-0004 command patterns
  - [ ] Integrate with existing resolver system (`storage.NewResolver`)
  - [ ] Implement file movement from trash to original location
  - [ ] Add registry updates using proven `registry.Update()` pattern
  - [ ] Use standard error types and user-friendly messages

- [ ] **Trash List Command**
  - [ ] Create `cmd/mas_trash.go` with list subcommand
  - [ ] Extend existing `ListFilters` to include status filtering
  - [ ] Reuse tabular display patterns from enhanced list command
  - [ ] Add trash-specific formatting (deletion dates, purge countdown)

### Phase 2: Purge Operations (Estimated: 3-4 hours)
- [ ] **Basic Purge Implementation**
  - [ ] Add purge subcommand to `cmd/mas_trash.go`
  - [ ] Implement single archive purge with resolver integration
  - [ ] Add permanent deletion (file + registry removal)
  - [ ] Implement safety confirmations and `--force` flag

- [ ] **Bulk Purge Features**
  - [ ] Add `--older-than` duration parsing (reuse existing patterns)
  - [ ] Implement age-based filtering for bulk operations  
  - [ ] Add batch processing with progress reporting
  - [ ] Include `--dry-run` support for safety

### Phase 3: Integration & Polish (Estimated: 2-3 hours)
- [ ] **Command Registration**
  - [ ] Update `main.go` to register new trash commands
  - [ ] Add help text and usage examples
  - [ ] Integrate with existing command structure

- [ ] **Testing & Validation**
  - [ ] Unit tests for core restore/purge functions
  - [ ] Integration tests using existing test patterns
  - [ ] Edge case testing (missing files, permission errors)
  - [ ] Performance validation with large trash sets

**Total Estimated Time: 9-13 hours**

### Dependencies ✅ All Complete
- **7EP-0004 MAS Foundation**: ✅ ULID resolution, registry operations, error handling
- **7EP-0006 Performance Testing**: ✅ Performance patterns validated  
- **Enhanced list display**: ✅ Tabular formatting and filtering implemented
- **Current delete/trash infrastructure**: ✅ Database schema and file operations ready

**Ready to Start**: All dependencies complete, implementation can begin immediately.

## Testing Strategy

### Acceptance Criteria
- [ ] Can restore deleted archive to original location
- [ ] Restore command works with ULID, name, and checksum prefixes
- [ ] Trash list shows only deleted archives with proper formatting
- [ ] Manual purge permanently removes archives and files
- [ ] Age-based purge respects time filters
- [ ] Cannot restore non-deleted archives
- [ ] Cannot restore if original location has conflicts

### Test Scenarios
- Restore managed vs external archives
- Handle file conflicts during restore
- Bulk purge operations
- Cross-device restore operations
- Permission handling for purge operations

## Migration/Compatibility

### Breaking Changes
None - all new functionality.

### Upgrade Path
No migration required - trash system already exists.

### Backward Compatibility
Fully compatible - extends existing delete functionality.

## Alternatives Considered

**Single `trash` command with subcommands**: Considered `7zarch-go trash restore <id>` but decided on top-level `restore` for discoverability and muscle memory alignment with common CLI patterns.

**Auto-restore on list**: Considered showing restore hints in regular list output but determined dedicated trash commands provide clearer workflow separation.

## Future Considerations

- **Trash compression**: Archive old trash items to save space
- **Scheduled purge**: Automated cleanup via cron/systemd integration  
- **Recovery metadata**: Store deletion reason/context for audit
- **Cross-system sync**: Sync trash state across devices

## References

- Related: Enhanced list display for deleted archives (implemented)
- Related: Auto-purge configuration via `retention_days` (implemented)
- GitHub Issue: TBD (to be created for discussion)