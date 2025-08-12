# show command

Displays detailed information about a specific archive using its ID, name, or checksum prefix.

## Synopsis

```bash
7zarch-go show <identifier>
```

## Description

The `show` command resolves an archive by its identifier and displays comprehensive information including location, status, metadata, and suggested actions. It verifies file existence and can optionally validate checksums.

## Arguments

| Argument | Description |
|----------|-------------|
| `identifier` | Archive identifier - can be ULID, ULID prefix, checksum prefix, or name |

## Flags

| Flag | Type | Description | Default |
|------|------|-------------|---------|
| `--verify` | bool | Verify archive checksum (slower for large files) | false |
| `--json` | bool | Output in JSON format for scripting | false |
| `--full-uid` | bool | Display complete ULID instead of short version | false |

## Examples

### Basic Usage

**Show archive by ULID prefix:**
```bash
7zarch-go show 01K2E33
```

**Show archive by full ULID:**
```bash
7zarch-go show 01K2E33XW4HTX7RVPS9Y6CRGDY
```

**Show archive by name:**
```bash
7zarch-go show project-backup.7z
```

**Show archive by checksum prefix:**
```bash
7zarch-go show a1b2c3d4
```

### Advanced Options

**Verify checksum integrity:**
```bash
7zarch-go show 01K2E33 --verify
```
- Reads entire archive to verify SHA-256 checksum
- Shows ✓ if valid, ❌ if corrupted

**JSON output for scripting:**
```bash
7zarch-go show 01K2E33 --json | jq '.size'
```

**Display full ULID:**
```bash
7zarch-go show project --full-uid
```

## Output Format

### Standard Output
```
Archive: project-backup.7z (01K2E33XW4)
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

### Status Indicators
- `✓` - File exists and accessible
- `⚠️` - File missing from expected location  
- `❌` - File corrupted (checksum mismatch)

### Location Types
- `managed` - Stored in MAS managed storage
- `external` - Stored outside managed storage
- `deleted` - In trash, awaiting purge
- `missing` - Registry entry exists but file not found

### JSON Output
```json
{
  "uid": "01K2E33XW4HTX7RVPS9Y6CRGDY",
  "name": "project-backup.7z",
  "path": "/Users/adam/.7zarch-go/archives/project-backup.7z",
  "size": 2202009,
  "created": "2025-08-11T14:30:22Z",
  "checksum": "sha256:a1b2c3d4e5f6789012345678901234567890abcdef",
  "profile": "documents",
  "managed": true,
  "status": "present",
  "uploaded": false,
  "verified": true,
  "original_size": 9134592,
  "compression_ratio": 0.759
}
```

## Resolution Process

The `show` command resolves identifiers in this priority order:

1. **Exact ULID match** - Full 26-character ULID
2. **ULID prefix** - Minimum 4 characters (configurable)
3. **Checksum prefix** - Minimum 8 characters
4. **Exact name match** - Complete filename
5. **Name without extension** - Filename without .7z

### Ambiguous Matches

When multiple archives match the identifier:
```
Multiple archives match '01K2':
[1] 01K2E33XW4 project-backup.7z (managed, 2.1 MB, 2d ago)
[2] 01K2F44YZ5 project-docs.7z (external, 0.5 MB, 1w ago)

Please specify full ULID or use a longer prefix
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success - archive found and displayed |
| 1 | Archive not found |
| 2 | Ambiguous identifier (multiple matches) |
| 3 | Registry access error |
| 4 | File verification failed |

## Related Commands

- **[list](list.md)** - List all archives with filters
- **[delete](delete.md)** - Move archive to trash
- **[move](move.md)** - Relocate archive
- **[test](test.md)** - Verify archive integrity

## Tips

### Performance
- Basic show is instant (registry lookup only)
- `--verify` reads entire file (slower for large archives)
- Use short ULID prefixes (8 chars) for quick typing

### Troubleshooting

**"Archive not found"**:
```bash
# Check if archive exists
7zarch-go list

# Check if in trash
7zarch-go trash list
```

**"Multiple archives match"**:
```bash
# Use more characters
7zarch-go show 01K2E33X  # Instead of 01K2

# Use full ULID shown in error
7zarch-go show 01K2E33XW4HTX7RVPS9Y6CRGDY
```

**"File missing"**:
```bash
# Archive moved outside 7zarch-go
7zarch-go move 01K2E33 --reattach /new/path/archive.7z

# Or mark as missing
7zarch-go db update 01K2E33 --status missing
```

## Configuration

The `show` command respects these configuration settings:

```yaml
ui:
  show_full_uid: false     # Show first 10 chars by default
  verify_on_show: false    # Don't verify checksums by default
  
storage:
  update_last_seen: true   # Update last_seen timestamp on show
```

## Use Cases

### Quick Archive Inspection
```bash
# See what's in your latest backup
7zarch-go list --limit 1
7zarch-go show 01K2E33  # Use the ID shown
```

### Verify Before Upload
```bash
# Ensure integrity before cloud upload
7zarch-go show important-backup --verify
# If verified ✓, safe to upload
```

### Troubleshoot Missing Files
```bash
# Find archives that point to missing files
7zarch-go list --missing
# Inspect each one
7zarch-go show 01K2E33
# Shows last known location and suggestions
```

### Script Integration
```bash
#!/bin/bash
# Get archive size in bytes
SIZE=$(7zarch-go show "$1" --json | jq -r '.size')
if [ $SIZE -gt 1073741824 ]; then
  echo "Archive larger than 1GB, uploading to cold storage"
fi
```

---

The `show` command is your primary tool for inspecting individual archives and understanding their current state in the MAS system.