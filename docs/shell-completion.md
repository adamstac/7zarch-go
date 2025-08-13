# Shell Completion for 7zarch-go

7zarch-go provides intelligent tab completion for all major shells including Bash, Zsh, Fish, and PowerShell.

## Features

- **Fast completion** (<50ms response time)
- **Archive ID completion** - Tab complete ULID prefixes, archive names, and checksum prefixes
- **Context-aware completion** - Different commands show different archives:
  - `show`, `delete`, `move`: All archives
  - `restore`: Only deleted archives
- **Performance optimized** - Concurrent lookups with timeout protection
- **Cross-platform** - Works on Linux, macOS, and Windows

## Installation

### Bash

**Option 1: Source directly in bashrc**
```bash
# Add to ~/.bashrc
source <(7zarch-go completion bash)
```

**Option 2: System-wide installation**
```bash
# Generate completion file
7zarch-go completion bash | sudo tee /etc/bash_completion.d/7zarch-go > /dev/null

# Or on macOS with Homebrew:
7zarch-go completion bash > $(brew --prefix)/etc/bash_completion.d/7zarch-go
```

### Zsh

**Option 1: Source directly in zshrc**
```zsh
# Add to ~/.zshrc
autoload -U compinit && compinit
source <(7zarch-go completion zsh)
```

**Option 2: Install to zsh completions directory**
```bash
# Create completions directory if it doesn't exist
mkdir -p ~/.zsh/completions

# Generate completion file
7zarch-go completion zsh > ~/.zsh/completions/_7zarch-go

# Add to ~/.zshrc if not already present
echo 'fpath=(~/.zsh/completions $fpath)' >> ~/.zshrc
echo 'autoload -U compinit && compinit' >> ~/.zshrc
```

### Fish

```fish
# Add to ~/.config/fish/config.fish
7zarch-go completion fish | source

# Or save permanently
7zarch-go completion fish > ~/.config/fish/completions/7zarch-go.fish
```

### PowerShell

```powershell
# Add to PowerShell profile
7zarch-go completion powershell | Out-String | Invoke-Expression

# To find your profile path:
echo $PROFILE
```

## Usage Examples

### Archive ID Completion

**UID Prefix Completion**
```bash
7zarch-go show 01K2<TAB>
# Shows: 01K2E3BEJV6GSKHZWWBSKVXEYT (test-pod-2.7z)
```

**Archive Name Completion**
```bash
7zarch-go delete test<TAB>
# Shows: test-pod-2.7z, test-pod.7z
```

**Context-Aware Completion**
```bash
# Show command - lists all archives
7zarch-go show <TAB>

# Restore command - only shows deleted archives  
7zarch-go restore <TAB>
```

### Checksum Prefix Completion
```bash
7zarch-go show ddd2a6a0<TAB>
# Shows: ddd2a6a05c24... (test-pod-2.7z) 01K2E3BE
```

## Performance

- **Response time**: <50ms for typical registries
- **Concurrent lookups**: UID, name, and checksum completion run in parallel
- **Timeout protection**: 50ms timeout prevents blocking
- **Graceful degradation**: Works efficiently with 1000+ archives
- **Smart limits**: Returns top 25 matches for optimal UX

## Troubleshooting

### Completion not working

1. **Verify installation**:
   ```bash
   7zarch-go completion bash --help
   ```

2. **Check shell configuration**:
   ```bash
   # For bash
   echo $BASH_COMPLETION_COMPAT_DIR
   
   # For zsh
   echo $fpath
   ```

3. **Test completion manually**:
   ```bash
   7zarch-go __complete show ""
   ```

### Slow completion

- Completion has a 50ms timeout - returns partial results if exceeded
- Check registry size with `7zarch-go list | wc -l`
- Large registries (1000+ archives) may see reduced completion results

### No archive completions shown

1. **Verify registry has archives**:
   ```bash
   7zarch-go list
   ```

2. **Check storage configuration**:
   ```bash
   7zarch-go config
   ```

3. **Test with verbose output**:
   ```bash
   BASH_COMP_DEBUG_FILE=/tmp/completion.log 7zarch-go __complete show ""
   cat /tmp/completion.log
   ```

## Architecture

The completion system uses:

- **Cobra framework** for shell completion generation
- **Concurrent goroutines** for parallel archive lookups
- **Context timeout** for performance guarantees
- **Filter functions** for command-specific completion
- **Smart caching** via the registry system

Commands with completion:
- `show <TAB>` - All archives
- `delete <TAB>` - All archives  
- `move <TAB>` - All archives
- `restore <TAB>` - Only deleted archives

The system automatically handles ULID prefix expansion, archive name matching, and checksum prefix completion with intelligent descriptions.
