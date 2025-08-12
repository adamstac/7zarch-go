# 7EP-0001: Trash Management System

**Status:** Draft  
**Author(s):** Adam Stacoviak, Claude Code (CC)  
**Assignment:** AC  
**Difficulty:** 3 (moderate - new commands + workflow integration)  
**Created:** 2025-08-11  
**Updated:** 2025-08-11  

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

**Why now:**
- Delete functionality is complete but incomplete without restore capability
- User testing revealed confusion about deleted archive recovery
- Auto-purge feature needs manual purge counterpart for user control

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
Design a safe, reviewable, and automatable trash lifecycle for soft-deleted archives with intelligent retention policies and user control.

### New Commands

#### `7zarch-go restore <id>`
- **Resolver integration**: Accept ULID, name, checksum prefix
- **Validation**: Only restore archives with `status="deleted"`
- **File operations**: Move from trash back to original location
- **Registry updates**: 
  - `status` → `"present"`
  - `path` → value from `original_path`
  - `deleted_at` → `NULL`
  - Preserve `original_path` for audit trail

#### `7zarch-go trash list [--details]`
- **Filtering**: Only show archives with `status="deleted"`
- **Display**: Similar to regular list but trash-focused
- **Details**: Show deletion date, original location, purge countdown

#### `7zarch-go trash purge [<id>] [--older-than <duration>]`
- **Manual purge**: Permanently delete specific archive or all trash
- **Age-based**: `--older-than 30d` for bulk cleanup
- **File operations**: Remove from filesystem and registry
- **Confirmation**: Require `--force` for permanent deletion

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

## Implementation Plan

### Phase 1: Core Functionality
- [ ] Implement `restore` command with resolver integration
- [ ] Add basic file movement and registry updates
- [ ] Implement `trash list` command
- [ ] Add `trash purge` with confirmation prompts

### Phase 2: Polish & Testing
- [ ] Add bulk operations support
- [ ] Implement `--older-than` filtering for purge
- [ ] Add comprehensive error handling
- [ ] Integration with existing list command improvements

### Dependencies
- Existing resolver system (implemented)
- Current trash infrastructure from delete command (implemented)
- Enhanced list display (recently implemented)

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