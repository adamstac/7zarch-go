# Getting Started with 7zarch-go

Welcome! This guide will get you up and running with 7zarch-go in just a few minutes.

## What is 7zarch-go?

7zarch-go is an intelligent archive management tool that:
- **Optimizes compression** automatically based on your content
- **Manages archives** with a local registry for easy tracking
- **Tests archives** concurrently for fast verification
- **Provides smart defaults** while remaining highly configurable

Perfect for backup workflows, content archival, and automated compression tasks.

## Prerequisites

7zarch-go requires the `7z` command-line tool:

### macOS
```bash
brew install p7zip
```

### Ubuntu/Debian
```bash
sudo apt install p7zip-full
```

### Windows
Download and install [7-Zip](https://www.7-zip.org/download.html).

## Installation

### Option 1: Download Binary (Recommended)
```bash
# Download latest release (replace with actual URL when available)
curl -L https://github.com/adamstac/7zarch-go/releases/latest/download/7zarch-go-linux -o 7zarch-go
chmod +x 7zarch-go
sudo mv 7zarch-go /usr/local/bin/
```

### Option 2: Build from Source
```bash
git clone https://github.com/adamstac/7zarch-go.git
cd 7zarch-go
go build -o 7zarch-go .
sudo cp 7zarch-go /usr/local/bin/
```

### Verify Installation
```bash
7zarch-go --version
```

## Your First Archive

Let's create your first archive with smart compression:

### 1. Create a Test Directory
```bash
mkdir ~/test-project
echo "Hello, 7zarch-go!" > ~/test-project/readme.txt
echo "console.log('test')" > ~/test-project/app.js
```

### 2. Create Your First Archive
```bash
7zarch-go create ~/test-project
```

This command:
- ✅ Analyzes your content (text files detected)
- ✅ Chooses optimal compression settings
- ✅ Creates `test-project.7z` in managed storage
- ✅ Registers the archive for easy tracking

### 3. List Your Archives
```bash
7zarch-go list
```

You should see:
```
📦 Archives (1 found)
Active: 1 (Managed: 1, External: 0) | Missing: 0 | Deleted: 0

ACTIVE - MANAGED
📦 test-project - 📤 Not uploaded
```

### 4. Get Detailed Information
```bash
7zarch-go list --details
```

Shows complete archive details including size, creation date, and unique ID.

### 5. Test Archive Integrity
```bash
7zarch-go test test-project
```

Verifies your archive is valid and can be extracted successfully.

## What Just Happened?

**Managed Storage**: Your archive was created in `~/.7zarch-go/archives/` and registered in a local SQLite database for easy management.

**Smart Compression**: 7zarch-go detected text files and applied the "documents" profile for maximum compression.

**ULID Tracking**: Each archive gets a unique, sortable ID for easy reference.

## Next Steps

### Learn Core Commands
- **[Basic Usage](basic-usage.md)** - Essential commands and workflows
- **[Configuration](configuration.md)** - Customize defaults and create presets

### Explore Advanced Features
- **[Advanced Usage](advanced-usage.md)** - Power user features
- **[Workflows](workflows/)** - Real-world usage patterns

### Get Help
- **[Troubleshooting](troubleshooting.md)** - Common issues and solutions
- **[Commands Reference](../reference/commands/)** - Complete command documentation

## Common Next Tasks

**Set up a configuration file:**
```bash
7zarch-go config init
```

**Create a preset for your workflow:**
```yaml
# Edit ~/.7zarch-go-config
presets:
  my-backup:
    profile: balanced
    comprehensive: true
```

**Use your preset:**
```bash
7zarch-go create important-data --preset my-backup
```

---

**Ready to learn more?** Continue with [Basic Usage](basic-usage.md) to master the core commands.