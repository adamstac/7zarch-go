# create command

Creates archives with intelligent compression optimization.

## Synopsis

```bash
7zarch-go create [flags] <path>
```

## Description

The `create` command compresses files and directories into 7z archives with intelligent optimization based on content analysis. It supports both managed storage (default) and custom output locations.

## Flags

| Flag | Type | Description | Default |
|------|------|-------------|---------|
| `--profile` | string | Compression profile (media/documents/balanced) | auto-detected |
| `--preset` | string | Use saved preset from configuration | none |
| `--compression` | int | Manual compression level (0-9, disables smart compression) | auto |
| `--comprehensive` | bool | Create archive with checksums and metadata | false |
| `--output` | string | Specify output location (bypasses managed storage) | managed |
| `--force` | bool | Overwrite existing archive without confirmation | false |
| `--dry-run` | bool | Show what would be done without executing | false |
| `--no-managed` | bool | Disable managed storage for this operation | false |
| `--exclude` | strings | Patterns to exclude from archive | none |
| `--threads` | int | Number of compression threads (0 = auto) | 0 |

## Examples

### Basic Usage

**Create archive with smart compression:**
```bash
7zarch-go create ~/Documents/project
```
- Analyzes content automatically
- Stores in managed storage as `project.7z`
- Registers in local database

**Create archive in specific location:**
```bash
7zarch-go create ~/Documents/project --output ~/Backups/project-backup.7z
```
- Bypasses managed storage
- Creates archive at specified path

### Compression Profiles

**Use media profile for fast compression:**
```bash
7zarch-go create ~/Videos/vacation --profile media
```
- Optimized for already-compressed content
- 3-5x faster than maximum compression
- ~10% larger files

**Use documents profile for maximum compression:**
```bash
7zarch-go create ~/Code/website --profile documents
```
- Optimized for text and code files
- Maximum compression ratio
- Slower but smallest files

**Use balanced profile for mixed content:**
```bash
7zarch-go create ~/Mixed/content --profile balanced
```
- Good balance of speed and compression
- Works well for varied file types

### Advanced Options

**Create comprehensive archive with metadata:**
```bash
7zarch-go create important-data --comprehensive
```
- Generates SHA256 checksums
- Creates compression log
- Includes metadata file

**Manual compression level:**
```bash
7zarch-go create large-dataset --compression 3
```
- Disables smart compression
- Uses specified level (0=store, 9=maximum)

**Exclude patterns:**
```bash
7zarch-go create ~/Code/project --exclude "*.log" --exclude ".git"
```
- Excludes log files and git directory
- Supports glob patterns

### Presets

**Create configuration preset:**
```yaml
# ~/.7zarch-go-config
presets:
  podcast:
    profile: media
    comprehensive: true
    exclude:
      - "*.DS_Store"
      - "*.log"
```

**Use preset:**
```bash
7zarch-go create episode-105 --preset podcast
```

### Overwrite Protection

**Force overwrite existing archive:**
```bash
7zarch-go create project --force
```

**Preview without executing:**
```bash
7zarch-go create project --dry-run
```

## Output

### Success Output
```
🔍 Analyzing content in /Users/adam/Documents/project...
📊 Detected: 85% text files, 15% images
🎯 Selected profile: documents (maximum compression)
📦 Creating archive: project.7z
✅ Archive created successfully
📋 Registered in managed storage
   Size: 15.4 MB → 2.8 MB (82% reduction)
   Location: ~/.7zarch-go/archives/project.7z
   ID: 01K2E33XW4HTX7RVPS9Y6CRGDY
```

### Comprehensive Mode Output
```
📦 Creating comprehensive archive...
✅ Archive created: project.7z
✅ Checksums generated: project.7z.sha256
✅ Compression log: project.7z.log
✅ Metadata file: project.7z.meta
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Archive creation failed |
| 2 | Input path not found |
| 3 | Permission denied |
| 4 | Archive already exists (without --force) |
| 5 | Configuration error |
| 10 | 7z command not found |

## Related Commands

- **[test](test.md)** - Verify archive integrity
- **[list](list.md)** - List managed archives
- **[config](config.md)** - Manage configuration and presets

## Tips

### Performance Optimization
- Use `--profile media` for video/audio files
- Set `--threads 4` on multi-core systems for large archives
- Use `--compression 5` for balance of speed and size

### Automation
```bash
# Backup script example
7zarch-go create ~/Documents --preset daily-backup --force
```

### Troubleshooting
- **"7z command not found"**: Install p7zip package
- **"Permission denied"**: Check file/directory permissions
- **"Archive exists"**: Use `--force` to overwrite or choose different name

## Configuration Integration

The `create` command respects these configuration settings:

```yaml
defaults:
  comprehensive: true     # Always create checksums
  compression: 5         # Default compression level
  profile: balanced      # Default profile when auto-detection fails

storage:
  use_managed_default: true  # Use managed storage by default
  managed_path: ~/.7zarch-go # Managed storage location
```

See [Configuration Guide](../user-guide/configuration.md) for complete configuration options.