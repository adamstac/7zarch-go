# Managed Archive Storage (MAS)

Managed Archive Storage (MAS) is 7zarch-go's intelligent local workspace that organizes your archives and tracks metadata in a lightweight SQLite registry.

## What is MAS?

MAS provides a centralized location for your archives with automatic tracking, making it easy to:
- **Find archives** without remembering paths
- **Track metadata** like creation date, size, and compression profile
- **Manage lifecycle** with soft deletion and auto-purge
- **Organize archives** with consistent naming and structure

Think of it as a "local archive database" that knows about all your compressed files.

## Why Use Managed Storage?

### Before MAS: Manual Archive Management
```bash
# Where did I put that backup?
find ~ -name "project-*.7z" 2>/dev/null
# What compression settings did I use?
# When was this created?
# Is this the latest version?
```

### With MAS: Intelligent Archive Management
```bash
# List all archives with metadata
7zarch-go list --details
# Find specific archive by name or ID
7zarch-go show project
# Restore deleted archives
7zarch-go restore project
```

## How MAS Works

### Directory Structure
```
~/.7zarch-go/
├── archives/           # Archive storage
│   ├── project.7z
│   ├── photos-2024.7z
│   └── backup-docs.7z
├── trash/             # Soft-deleted archives
│   └── old-project.7z
├── registry.db        # SQLite metadata database
└── config            # Local configuration
```

### Registry Database
The SQLite registry (`registry.db`) tracks:
- **Archive metadata**: name, size, checksum, creation date
- **Compression info**: profile used, compression ratio
- **Status tracking**: present, deleted, missing, uploaded
- **Relationships**: original paths, deletion timestamps
- **User metadata**: custom tags and descriptions

### File Permissions
- **Registry**: `600` (owner read/write only)
- **Archives**: Inherit from user umask
- **Trash**: Same as archives directory

## Archive Lifecycle

### 1. Creation
```bash
7zarch-go create ~/Documents/project
```
- Archive created in `~/.7zarch-go/archives/project.7z`
- Metadata registered in database
- ULID assigned for unique identification

### 2. Active Management
```bash
# List all archives
7zarch-go list

# Show detailed information
7zarch-go show project

# Test integrity
7zarch-go test project
```

### 3. Soft Deletion
```bash
7zarch-go delete project
```
- Archive moved to `~/.7zarch-go/trash/`
- Status changed to "deleted" in registry
- Original path preserved for restoration

### 4. Restoration or Purge
```bash
# Restore from trash
7zarch-go restore project

# Or permanently delete
7zarch-go trash purge project
```

## Configuration Options

### Enable/Disable MAS
```yaml
# ~/.7zarch-go-config
storage:
  use_managed_default: true    # Use MAS by default
  managed_path: ~/.7zarch-go   # Custom MAS location
  retention_days: 30           # Auto-purge deleted archives
```

### Bypass MAS for Specific Operations
```bash
# Create archive outside MAS
7zarch-go create project --output ~/Backups/project.7z

# Temporarily disable MAS
7zarch-go create project --no-managed
```

## Advanced Usage

### External Archive Registration
Register archives created outside MAS:
```bash
# Register existing archive
7zarch-go register ~/Backups/old-archive.7z

# List shows both managed and external
7zarch-go list
```

### Archive Migration
Move external archives into MAS:
```bash
# Import external archive
7zarch-go move ~/Backups/project.7z
```

### Batch Operations
```bash
# Test all managed archives
7zarch-go test --managed

# List archives older than 30 days
7zarch-go list --older-than 30d

# Purge old deleted archives
7zarch-go trash purge --older-than 7d
```

## Querying and Filtering

### Basic Queries
```bash
# All archives
7zarch-go list

# Only managed archives
7zarch-go list --managed

# Only external archives
7zarch-go list --external

# Missing archives (files deleted outside 7zarch-go)
7zarch-go list --missing
```

### Advanced Filtering
```bash
# Archives matching pattern
7zarch-go list --pattern "project-*"

# Archives older than duration
7zarch-go list --older-than 1w

# Not uploaded archives
7zarch-go list --not-uploaded

# Detailed view with all metadata
7zarch-go list --details
```

### ULID Resolution
Every archive gets a ULID (Universally Unique Lexicographically Sortable Identifier):
```bash
# Reference by name
7zarch-go show project

# Reference by ULID
7zarch-go show 01K2E33XW4HTX7RVPS9Y6CRGDY

# Reference by ULID prefix
7zarch-go show 01K2E33
```

## Backup and Recovery

### Backing Up MAS
```bash
# Backup entire MAS directory
cp -r ~/.7zarch-go ~/Backups/7zarch-go-backup

# Backup just the registry
cp ~/.7zarch-go/registry.db ~/Backups/registry-backup.db
```

### Restoring MAS
```bash
# Restore entire MAS
cp -r ~/Backups/7zarch-go-backup ~/.7zarch-go

# Restore just registry (if archives are intact)
cp ~/Backups/registry-backup.db ~/.7zarch-go/registry.db
```

### Registry Repair
```bash
# Verify registry integrity
7zarch-go db verify

# Rebuild registry from existing archives
7zarch-go db rebuild

# Check for missing files
7zarch-go list --missing
```

## Performance Considerations

### Storage Requirements
- **Registry**: ~1KB per archive entry
- **Overhead**: Minimal (SQLite is very efficient)
- **Archive storage**: Same as original files

### Query Performance
- **Local queries**: Instant (SQLite index)
- **File operations**: Limited by disk I/O
- **Network**: No network calls for MAS operations

### Scalability
- **Archive count**: Tested with 10,000+ archives
- **Database size**: <10MB for 10,000 entries
- **Query speed**: Sub-millisecond for most operations

## Troubleshooting

### Common Issues

**"Registry not found" error:**
```bash
# Initialize MAS
7zarch-go db init
```

**"Permission denied" on registry:**
```bash
# Fix permissions
chmod 600 ~/.7zarch-go/registry.db
```

**Registry corruption:**
```bash
# Backup and rebuild
cp ~/.7zarch-go/registry.db ~/.7zarch-go/registry.db.backup
7zarch-go db rebuild
```

**Archives show as missing:**
```bash
# Check if files were moved/deleted outside 7zarch-go
7zarch-go list --missing

# Re-register if archives exist elsewhere
7zarch-go register /path/to/moved/archive.7z
```

### Diagnostics
```bash
# Show MAS status
7zarch-go db status

# Verify database integrity
7zarch-go db verify

# Show configuration
7zarch-go config show
```

## Migration from Manual Management

### Import Existing Archives
```bash
# Register existing .7z files
find ~/Archives -name "*.7z" -exec 7zarch-go register {} \;

# Or move them into MAS
find ~/Archives -name "*.7z" -exec 7zarch-go move {} \;
```

### Gradual Migration
You can use both managed and external archives:
- New archives: Use MAS by default
- Existing archives: Register or leave external
- Gradually move important archives into MAS

## Integration with Other Tools

### Backup Scripts
```bash
#!/bin/bash
# Backup script using MAS
7zarch-go create ~/Documents --preset backup
7zarch-go create ~/Code --preset backup

# Upload archives marked as not-uploaded
for archive in $(7zarch-go list --not-uploaded --format ids); do
    upload_to_cloud "$archive"
    7zarch-go db mark-uploaded "$archive"
done
```

### CI/CD Integration
```yaml
# GitHub Actions example
- name: Create release archive
  run: 7zarch-go create ./build --profile balanced --output release.7z
```

## Best Practices

### Naming Conventions
- Use descriptive names: `project-website` not `proj1`
- Include dates for time-series: `backup-2024-01-15`
- Use consistent prefixes: `backup-`, `archive-`, `release-`

### Organization Strategies
- **By purpose**: `backup-*`, `release-*`, `archive-*`
- **By date**: `YYYY-MM-DD-*` for chronological sorting
- **By project**: `projectname-*` for grouping

### Maintenance
- **Regular cleanup**: Use auto-purge or manual `trash purge`
- **Integrity checks**: Periodic `7zarch-go test --managed`
- **Storage monitoring**: Check `7zarch-go db status` for size

---

MAS transforms ad-hoc archive creation into a systematic, trackable process that scales from personal use to enterprise backup workflows.