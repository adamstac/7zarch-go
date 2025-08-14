# Batch Operations Command

The `batch` command performs operations on multiple archives efficiently with progress tracking and error handling.

## Synopsis

```bash
7zarch-go batch <operation> [flags]
```

## Operations

### move
Move multiple archives to a new location.

```bash
# Move archives using saved query
7zarch-go batch move --query=old-files --to=/archive/old/

# Move specific archives from stdin
echo -e "01J3K4...\n01J3K5..." | 7zarch-go batch move --stdin --to=/backup/
```

### delete
Delete multiple archives with trash integration.

```bash
# Delete archives using saved query (requires confirmation)
7zarch-go batch delete --query=temp-files --confirm

# Delete specific archives from stdin
7zarch-go list --older-than=1y --output=json | jq -r '.[].uid' | 7zarch-go batch delete --stdin --confirm
```

## Selection Methods

### Saved Queries
Use saved query to select archives:

```bash
7zarch-go batch move --query=media-large --to=/external/media/
```

### Standard Input
Read archive UIDs from stdin (one per line):

```bash
7zarch-go list --profile=documents --output=json | jq -r '.[].uid' | 7zarch-go batch move --stdin --to=/backup/docs/
```

### Filters (Future)
Direct filter support will be added in future versions.

## Flags

### Selection Flags
- `--query=<name>` - Use saved query to select archives
- `--stdin` - Read archive UIDs from stdin

### Operation Flags
- `--to=<path>` - Destination path for move operation (required for move)
- `--confirm` - Confirm destructive operations (required for delete)
- `--dry-run` - Show what would be done without executing

### Performance Flags
- `--concurrent=<n>` - Number of concurrent operations (default: 4)
- `--progress` - Show progress during batch operations (default: true)

### Output Flags
- `--output=<format>` - Output format: table, json, csv, yaml (default: table)

## Progress Tracking

Batch operations show real-time progress:

```
Progress: 45/100 (45.0%) - Processing backup-2024.7z
```

On completion:
```
Progress: 100/100 (100.0%) - Complete - Completed in 2m34s
Successfully moved 100 archives to /backup/archives/
```

## Error Handling

- **Partial failures**: Operations continue on individual failures
- **Error collection**: All errors reported at completion
- **Context cancellation**: Operations can be cancelled (Ctrl+C)
- **Rollback**: No automatic rollback (design decision for safety)

## Safety Features

### Confirmation Required
Destructive operations require explicit confirmation:

```bash
# This will fail
7zarch-go batch delete --query=old-files
# Error: delete operation requires --confirm flag for safety

# This will work
7zarch-go batch delete --query=old-files --confirm
```

### Dry Run
Preview operations before execution:

```bash
7zarch-go batch move --query=media-files --to=/external/ --dry-run
```

Output:
```
Selected 25 archive(s) for move operation:
  - vacation-2023.7z (01J3K4M2...)
  - photos-backup.7z (01J3K4M3...)
  ...
Dry run - no operations performed
```

### Overwrite Protection
Move operations prevent accidental overwrites:

```bash
# Will fail if files exist at destination
7zarch-go batch move --query=docs --to=/backup/existing/
```

## Performance Characteristics

- **Concurrency**: 4 operations by default (configurable)
- **Memory usage**: Streams operations, doesn't load all archives into memory
- **Progress updates**: Every 1-2 seconds during processing
- **Cross-device moves**: Automatic fallback to copy+remove

## Examples

### Basic Operations

Move all media files to external storage:
```bash
7zarch-go query save "media-files" --profile=media
7zarch-go batch move --query=media-files --to=/external/media/
```

Clean up old temporary archives:
```bash
7zarch-go query save "old-temp" --pattern="temp-*" --older-than=30d
7zarch-go batch delete --query=old-temp --confirm
```

### Pipeline Integration

Archive cleanup workflow:
```bash
# Find archives to clean up
7zarch-go list --older-than=1y --status=missing --output=json > cleanup.json

# Review the list
cat cleanup.json | jq -r '.[].name'

# Delete them
cat cleanup.json | jq -r '.[].uid' | 7zarch-go batch delete --stdin --confirm
```

### Advanced Usage

High-performance batch move with custom concurrency:
```bash
7zarch-go batch move --query=large-archives --to=/fast-storage/ --concurrent=8 --progress
```

Quiet batch operation (no progress):
```bash
7zarch-go batch delete --query=temp-files --confirm --progress=false
```

## Integration with Other Commands

### Query System
Batch operations work seamlessly with saved queries:

```bash
# Save complex query
7zarch-go query save "cleanup-candidates" --older-than=6m --status=ok --profile=documents

# Use in batch operation
7zarch-go batch move --query=cleanup-candidates --to=/archive/old-docs/
```

### Search Integration
Combine with search results:

```bash
# Search and save as query
7zarch-go search query "project backup" --save-query=project-backups

# Process results
7zarch-go batch move --query=project-backups --to=/project-archive/
```

### List Integration
Pipeline with list command:

```bash
# Export specific archives
7zarch-go list --managed --larger-than=1GB --output=json | jq -r '.[].uid' | 7zarch-go batch move --stdin --to=/big-storage/
```

## Error Scenarios

### Move Failures
```bash
# Permission denied
Error: failed to move archive1.7z: permission denied

# Destination exists
Error: destination file already exists: /backup/archive1.7z

# Cross-device fallback
Warning: Using copy+remove fallback for cross-device move
```

### Delete Failures
```bash
# File in use
Error: failed to delete archive1.7z: file is being used by another process

# Permission denied
Error: failed to delete archive1.7z: permission denied
```

### Partial Success
```bash
Progress: 100/100 (100.0%) - Completed with 3 errors in 1m23s

Errors:
  - failed to process backup-2024.7z: permission denied
  - failed to process temp-archive.7z: file not found
  - failed to process large-file.7z: disk full

Successfully processed 97 archives
```

## Technical Implementation

- **Worker pool**: Configurable concurrency with goroutine pool
- **Thread-safe**: Progress tracking and error collection
- **Context aware**: Supports cancellation and timeouts
- **Resource efficient**: Bounded memory usage regardless of archive count
- **Cross-platform**: Handles filesystem differences (rename vs copy+remove)

## Future Enhancements

- **Filter integration**: Direct filter support without queries
- **Rollback capability**: Optional transaction-like behavior
- **Resume support**: Continue interrupted operations
- **Batch create**: Create multiple archives in one operation
- **Upload integration**: Batch upload to cloud providers