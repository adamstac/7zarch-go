# 7zarch-go

An intelligent archive management tool that optimizes compression based on content type. Create, test, and manage archives with smart defaults and powerful configuration options.

## Features

- **Smart Compression** - Automatically detects file types and optimizes compression settings
- **Concurrent Testing** - Test multiple archives in parallel for 10x faster verification
- **Configuration Presets** - Save and reuse common archive settings
- **Comprehensive Mode** - Create archives with checksums and metadata in one command
- **Single Binary** - No dependencies, runs anywhere Go runs

## Installation

### Prerequisites

```bash
# Install 7-Zip (required)
brew install p7zip    # macOS
apt install p7zip     # Ubuntu/Debian
```

### Build from Source

```bash
git clone https://github.com/adamstac/7zarch-go.git
cd 7zarch-go
go build -o 7zarch-go .
```

### Install System-wide

```bash
# Option 1: Copy to PATH
cp 7zarch-go /usr/local/bin/

# Option 2: Create symlink (for development)
ln -s $(pwd)/7zarch-go /usr/local/bin/7zarch-go
```

## Quick Start

### Create Your First Archive

```bash
# Basic archive
7zarch-go create my-project

# Smart compression is the default (analyzes content, picks optimal settings)
7zarch-go create my-videos

# Use a specific compression profile
7zarch-go create podcast-episode --profile media  # 3x faster for media files
7zarch-go create source-code --profile documents   # Maximum compression for text

# Comprehensive mode (archive + checksums + metadata)
7zarch-go create important-data --comprehensive
```

### Test Archive Integrity

```bash
# Test single archive
7zarch-go test archive.7z

# Test all archives in directory (concurrent)
7zarch-go test --directory /path/to/archives --concurrent 10
```

## Configuration

### Create a Config File

```bash
7zarch-go config init
```

This creates `~/.7zarch-go-config` with your preferences:

```yaml
defaults:
  comprehensive: true    # Always create checksums and logs
  compression: 5         # Default compression level
  threads: 0            # 0 = use all CPU cores

presets:
  podcast:
    profile: media       # Use media profile for fast compression
    comprehensive: true  # Include checksums and logs
    excludes:
      - "*.DS_Store"
      - "*.log"
      - ".git"
```

### Use Presets

```bash
# Use your podcast preset
7zarch-go create episode-103 --preset podcast

# View current config
7zarch-go config show
```

## Compression Profiles

7zarch-go includes intelligent compression profiles optimized for different content types:

### Media Profile
- **Best for**: Videos, audio, images
- **Speed**: 3-5x faster than default
- **Size**: ~10% larger than maximum compression
- **Use case**: Podcast episodes, video projects, photo archives

### Documents Profile  
- **Best for**: Text, code, office documents
- **Speed**: Slower but maximum compression
- **Size**: Smallest possible archive
- **Use case**: Source code, documentation, spreadsheets

### Balanced Profile
- **Best for**: Mixed content
- **Speed**: Good balance
- **Size**: Good compression ratio
- **Use case**: General backups, mixed file types

### View Available Profiles

```bash
7zarch-go profiles
```

## Commands Reference

### create

Create an archive from a directory or file.

```bash
7zarch-go create [flags] <path>
```

**Flags:**
- `--profile <name>` - Use compression profile (media/documents/balanced)
- `--preset <name>` - Use saved preset from config
- `--compression <0-9>` - Manual compression level (disables smart behavior)
- `--comprehensive` - Create archive with checksums and metadata
- `--output <path>` - Specify output location
- `--force` - Overwrite existing archive
- `--dry-run` - Show what would be done without doing it

**Examples:**

```bash
# Smart compression with content analysis
7zarch-go create my-project --smart-compression

# Use preset from config
7zarch-go create podcast-103 --preset podcast

# Specify output location
7zarch-go create data --output /backups/data.7z

# Comprehensive with forced overwrite
7zarch-go create important --comprehensive --force
```

### test

Test archive integrity.

```bash
7zarch-go test [flags] <archive or directory>
```

**Flags:**
- `--directory` - Test all archives in directory
- `--concurrent <n>` - Number of parallel tests (default: 10)
- `--dry-run` - Show what would be tested

**Examples:**

```bash
# Test single archive
7zarch-go test backup.7z

# Test directory with 5 parallel workers
7zarch-go test --directory /archives --concurrent 5
```

### profiles

List available compression profiles.

```bash
7zarch-go profiles
```

### config

Manage configuration.

```bash
7zarch-go config <subcommand>
```

**Subcommands:**
- `init` - Create default config file
- `show` - Display current configuration

## Real-World Examples

### Podcast Production Workflow

```bash
# Configure once
7zarch-go config init
# Edit ~/.7zarch-go-config to add podcast preset

# Archive each episode
7zarch-go create friends-103 --preset podcast
7zarch-go create friends-104 --preset podcast

# Verify integrity
7zarch-go test --directory ~/archives/friends
```

### Code Backup Workflow

```bash
# Maximum compression for source code
7zarch-go create ~/Code/my-project --profile documents --comprehensive

# Quick backup with smart detection
7zarch-go create ~/Code/website --smart-compression
```

### Media Archive Workflow

```bash
# Fast compression for video files
7zarch-go create video-project --profile media

# Archive photos with metadata
7zarch-go create photos-2024 --profile media --comprehensive
```

## Performance Tips

### Concurrent Testing
When testing multiple archives, use `--concurrent` to dramatically reduce time:

```bash
# Test 100 archives in the time it takes to test 10
7zarch-go test --directory /backups --concurrent 10
```

### Smart Compression
Smart compression is enabled by default; it analyzes your content and chooses optimal settings:

```bash
# Analyzes content, shows recommendation, applies best profile
7zarch-go create mixed-content
```

### Profile Selection
Choose the right profile for your content:
- **Media files**: Use `--profile media` for 3-5x faster compression
- **Text/code**: Use `--profile documents` for maximum compression
- **Mixed**: Use `--profile balanced` or `--smart-compression`

## Troubleshooting

### Archive Already Exists
Use `--force` to overwrite:
```bash
7zarch-go create project --force
```

### Out of Memory on Large Archives
Reduce compression level or dictionary size:
```bash
7zarch-go create huge-dataset --compression 5
```

### Can't Find 7z Command
Install p7zip:
```bash
brew install p7zip  # macOS
apt install p7zip   # Ubuntu/Debian
```

## Development

### Building
```bash
go build -o 7zarch-go .
```

### Testing
```bash
go test ./...
```

### Project Structure
```
7zarch-go/
├── main.go                 # Entry point
├── cmd/                    # CLI commands
│   ├── create.go          # Archive creation
│   ├── test.go            # Archive testing
│   ├── profiles.go        # Profile listing
│   └── config.go          # Configuration
└── internal/              # Core logic
    ├── archive/           # Compression logic
    └── config/            # Config management
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## License

MIT License - see LICENSE file for details

## Author

Adam Stacoviak ([@adamstac](https://github.com/adamstac))