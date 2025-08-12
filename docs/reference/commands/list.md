# list command

Lists archives tracked in the MAS registry with powerful filtering and grouping options.

## Synopsis

```bash
7zarch-go list [flags]
```

## Description

The `list` command displays all archives tracked in the MAS registry, organized by status and location. It supports comprehensive filtering to find specific archives and provides both summary and detailed views.

## Flags

| Flag | Type | Description | Default |
|------|------|-------------|---------|
| `--directory` | string | List .7z files in specific directory instead of registry | - |
| `--details` | bool | Show detailed information for each archive | false |
| `--not-uploaded` | bool | Show only archives that haven't been uploaded | false |
| `--pattern` | string | Filter archives by name pattern (glob syntax) | - |
| `--older-than` | string | Show archives older than duration (e.g., '7d', '2w', '168h') | - |
| `--managed` | bool | Show only managed archives | false |
| `--external` | bool | Show only external archives | false |
| `--missing` | bool | Show only missing archives | false |
| `--status` | string | Filter by status (present, missing, deleted) | - |
| `--profile` | string | Filter by compression profile | - |
| `--larger-than` | string | Show archives larger than size (e.g., '100MB', '1GB') | - |

## Output Format

### Standard Output (Tabular)
```
📦 Archives (15 found, 45.2 GB total)
Active: 10 (Managed: 8, External: 2) | Missing: 2 | Deleted: 3

ACTIVE - MANAGED
UID         Name                   Size      Profile    Age       Status
01K2E33XW4  project-backup.7z      2.1 MB    documents  2d ago    ✓
01K2F44YZ5  podcast-103.7z         156 MB    media      1w ago    ✓
01K2G55AB6  code-archive.7z        0.8 MB    documents  2w ago    ⚠️

ACTIVE - EXTERNAL
01K2H66CD7  /backup/old-site.7z    45 MB     balanced   1m ago    ✓
01K2I77EF8  ~/Desktop/temp.7z      2 MB      media      3d ago    ✓

DELETED (auto-purge older than 30 days)
🗑️  deleted-archive.7z - Deleted 2025-08-10 14:30:22
🗑️  old-backup.7z - Deleted 2025-08-09 10:15:33
```

### Detailed Output (`--details`)
```
📦 Archives (15 found, 45.2 GB total)
Active: 10 (Managed: 8, External: 2) | Missing: 2 | Deleted: 3

ACTIVE - MANAGED
📦 project-backup.7z - 📤 Not uploaded
   ID: 01K2E33XW4HTX7RVPS9Y6CRGDY
   Path: ~/.7zarch-go/archives/project-backup.7z
   Size: 2.1 MB
   Created: 2025-08-10 14:30:22
   Profile: documents
   Age: 2d

DELETED (auto-purge older than 30 days)
🗑️  deleted-archive.7z - Deleted 2025-08-10 14:30:22
   ID: 01K2D22VW3HTX7RVPS9Y6CRGDY
   Auto-purge: 28 days (2025-09-07)
   Original: ~/Documents/deleted-archive.7z
   Trash: ~/.7zarch-go/trash/deleted-archive.7z
   Size: 5.2 MB
```

## Filtering Examples

### Basic Filters

**Show only managed archives:**
```bash
7zarch-go list --managed
```

**Show archives not yet uploaded:**
```bash
7zarch-go list --not-uploaded
```

**Show missing archives:**
```bash
7zarch-go list --missing
```

### Pattern Matching

**Filter by name pattern:**
```bash
7zarch-go list --pattern "project-*"
7zarch-go list --pattern "*.backup.7z"
```

### Time-Based Filters

**Show archives older than 7 days:**
```bash
7zarch-go list --older-than 7d     # 7 days
7zarch-go list --older-than 2w     # 2 weeks
7zarch-go list --older-than 720h   # 720 hours (30 days)
```

### Size-Based Filters

**Show large archives:**
```bash
7zarch-go list --larger-than 100MB
7zarch-go list --larger-than 1GB
```

### Status and Profile Filters

**Filter by status:**
```bash
7zarch-go list --status present    # Active archives
7zarch-go list --status deleted    # In trash
7zarch-go list --status missing    # File not found
```

**Filter by compression profile:**
```bash
7zarch-go list --profile media      # Media-optimized archives
7zarch-go list --profile documents  # Document-optimized archives
```

### Combined Filters

**Complex queries:**
```bash
# Large media files not uploaded
7zarch-go list --profile media --larger-than 100MB --not-uploaded

# Old external archives
7zarch-go list --external --older-than 30d

# Missing managed archives
7zarch-go list --managed --status missing
```

## Directory Listing

List .7z files in a specific directory (bypasses registry):
```bash
7zarch-go list --directory ~/Backups/
```

Output:
```
📁 Listing .7z files in: ~/Backups/

Found 3 archive(s):

📦 backup-2025-08-01.7z
📦 backup-2025-08-08.7z
📦 backup-2025-08-15.7z
```

With details:
```bash
7zarch-go list --directory ~/Backups/ --details
```

## Status Indicators

- `✓` - Archive present and accessible
- `⚠️` - Archive missing from expected location
- `📤` - Not uploaded to remote storage
- `✅` - Uploaded to remote storage
- `🗑️` - Deleted (in trash)

## Empty State

When no archives match the filters:
```
No archives found.
💡 Tip: Create archives with '7zarch-go create <path>' to see them here.
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success - archives listed |
| 1 | Registry access error |
| 2 | Invalid filter parameters |
| 3 | Directory not found (with --directory) |

## Configuration

The `list` command respects these configuration settings:

```yaml
ui:
  default_list_limit: 50    # Pagination limit
  show_full_uid: false      # Show first 10 chars by default
  
storage:
  retention_days: 30        # For auto-purge display
```

## Use Cases

### Daily Status Check
```bash
# Morning archive status
7zarch-go list --not-uploaded       # What needs uploading?
7zarch-go list --missing            # Any files disappeared?
7zarch-go list --older-than 30d     # What can be cleaned up?
```

### Storage Management
```bash
# Find space-consuming archives
7zarch-go list --larger-than 1GB --details

# Identify cleanup candidates
7zarch-go list --older-than 90d --external
```

### Upload Workflow
```bash
# List archives ready for upload
7zarch-go list --not-uploaded --managed

# After upload, verify
7zarch-go list --pattern "uploaded-*"
```

### Troubleshooting
```bash
# Find problematic archives
7zarch-go list --status missing --details

# Check specific profile performance
7zarch-go list --profile media --details
```

## Tips

### Performance
- Filters are applied sequentially for optimal performance
- Use specific filters to reduce result sets
- Pattern matching uses glob syntax (faster than regex)

### Organization
- The tabular format aligns columns for easy scanning
- Archives are grouped by status and location
- Deleted archives show auto-purge countdown

### Scripting
```bash
#!/bin/bash
# Find large archives not uploaded
for id in $(7zarch-go list --larger-than 1GB --not-uploaded --format ids); do
    echo "Processing $id"
    7zarch-go show "$id"
done
```

## Related Commands

- **[show](show.md)** - Display details for specific archive
- **[create](create.md)** - Create new archives
- **[delete](delete.md)** - Move archives to trash
- **[trash](trash.md)** - Manage deleted archives

---

The `list` command is your primary tool for understanding the state of your archive collection and finding specific archives based on various criteria.