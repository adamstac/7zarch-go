# 7EP-0007: Enhanced MAS Operations

**Status:** Draft  
**Author(s):** Claude Code (CC)  
**Assignment:** AC (Primary), CC (Supporting)  
**Difficulty:** 3 (moderate - building on MAS Foundation with new operations)  
**Created:** 2025-08-12  
**Updated:** 2025-08-12  

## Executive Summary

Enhance MAS (Managed Archive Storage) with advanced operations including move/relocate functionality, import workflows for external archives, batch operations, and archive replication to complete the comprehensive archive management suite.

## Evidence & Reasoning

**Current State (Post 7EP-0004):**
- ✅ Core MAS Foundation complete with ULID resolution, show, list operations
- ✅ Performance validated (100-2,941x faster than requirements)
- ✅ Create, delete, and basic registry operations working
- ❌ **Missing**: Advanced operations for comprehensive archive lifecycle

**User Needs Identified:**
- **Archive Mobility**: Move archives between locations without re-compression
- **External Integration**: Import existing .7z files into MAS registry
- **Batch Operations**: Operate on multiple archives efficiently
- **Location Management**: Organize archives across storage locations
- **Registry Maintenance**: Repair, verify, and optimize registry data

**Why Now:**
- **MAS Foundation Stable**: Core infrastructure proven and performant
- **User Workflow Gaps**: Basic operations leave gaps in complete archive management
- **Building Momentum**: Continue MAS development while patterns are fresh
- **Strategic Timing**: Before trash management (7EP-0001) adds complexity

## Use Cases

### Primary Use Case: Archive Organization and Mobility
```bash
# Move archive to different location
7zarch-go move 01K2E33 --to /backup/archives/

# Move with rename
7zarch-go move project-backup.7z --to /archive/project-backup-final.7z

# Relocate managed archive to external storage
7zarch-go move 01K2E33 --to /external/backup/ --external

# Move multiple archives
7zarch-go move --pattern "project-*" --to /completed/
```

### Secondary Use Cases
- **External Archive Import**: `7zarch-go import /downloads/*.7z` - Add existing archives to registry
- **Batch Operations**: `7zarch-go batch --profile media --action verify` - Mass operations
- **Archive Verification**: `7zarch-go verify --all` - Registry consistency checking  
- **Storage Optimization**: `7zarch-go organize --by-profile` - Auto-organize by criteria

## Technical Design

### Overview
Build on 7EP-0004's proven ULID resolution and registry patterns to add advanced operations that maintain registry consistency while providing powerful archive management capabilities.

### Component Architecture

#### 1. Move/Relocate Operations
```go
// cmd/mas_move.go - Enhanced version of existing move command
type MoveOptions struct {
    Source      string   // ULID, name, or pattern
    Destination string   // Target path or directory
    Rename      string   // Optional new name
    External    bool     // Move to external storage
    Copy        bool     // Copy instead of move (replication)
    Verify      bool     // Verify after move
    Batch       bool     // Allow batch operations
}

func runMasMove(cmd *cobra.Command, args []string) error {
    resolver := storage.NewResolver(registry)
    
    // Resolve source archive(s)
    archives, err := resolver.ResolvePattern(args[0])
    if err != nil {
        return handleResolutionError(err)
    }
    
    // Process each archive
    for _, archive := range archives {
        if err := executeMove(archive, options); err != nil {
            return fmt.Errorf("failed to move %s: %w", archive.Name, err)
        }
    }
    
    return nil
}

func executeMove(archive *storage.Archive, opts MoveOptions) error {
    // Verify source exists and is accessible
    if err := verifyArchiveAccess(archive); err != nil {
        return err
    }
    
    // Determine target path
    targetPath, err := resolveTargetPath(archive, opts)
    if err != nil {
        return err
    }
    
    // Perform filesystem move/copy
    if opts.Copy {
        err = copyArchiveFile(archive.Path, targetPath)
    } else {
        err = moveArchiveFile(archive.Path, targetPath)
    }
    if err != nil {
        return err
    }
    
    // Update registry
    updatedArchive := *archive
    updatedArchive.Path = targetPath
    updatedArchive.Managed = !opts.External
    updatedArchive.Name = filepath.Base(targetPath)
    
    return registry.Update(&updatedArchive)
}
```

#### 2. Import Operations
```go
// cmd/mas_import.go
type ImportOptions struct {
    Pattern     string   // File pattern to import
    Profile     string   // Compression profile to assign
    Managed     bool     // Import as managed (copy) or external (link)
    Verify      bool     // Verify checksums during import
    Batch       bool     // Process multiple files
    Interactive bool     // Prompt for each file
}

func runMasImport(cmd *cobra.Command, args []string) error {
    // Discover archive files
    files, err := discoverArchiveFiles(args[0], options.Pattern)
    if err != nil {
        return err
    }
    
    fmt.Printf("Found %d archive files to import\n", len(files))
    
    // Process each file
    var imported, skipped, failed int
    for _, file := range files {
        result := processImport(file, options)
        switch result.Status {
        case ImportSuccess:
            imported++
            fmt.Printf("✓ Imported: %s → %s\n", result.Source, result.Archive.UID[:8])
        case ImportSkipped:
            skipped++
            fmt.Printf("- Skipped: %s (%s)\n", result.Source, result.Reason)
        case ImportFailed:
            failed++
            fmt.Printf("✗ Failed: %s (%s)\n", result.Source, result.Error)
        }
    }
    
    fmt.Printf("\nImport complete: %d imported, %d skipped, %d failed\n", 
               imported, skipped, failed)
    return nil
}

func processImport(filePath string, opts ImportOptions) ImportResult {
    // Check if already in registry
    existing, err := findExistingArchive(filePath)
    if err == nil {
        return ImportResult{Status: ImportSkipped, Reason: "Already in registry"}
    }
    
    // Analyze archive
    info, err := analyzeArchiveFile(filePath)
    if err != nil {
        return ImportResult{Status: ImportFailed, Error: err.Error()}
    }
    
    // Create registry entry
    archive := &storage.Archive{
        UID:      generateUID(),
        Name:     filepath.Base(filePath),
        Path:     filePath, // External path initially
        Size:     info.Size,
        Created:  info.Created,
        Checksum: info.Checksum,
        Profile:  opts.Profile,
        Managed:  false, // Start as external
        Status:   "present",
    }
    
    // If importing as managed, copy to managed storage
    if opts.Managed {
        managedPath, err := copyToManagedStorage(filePath, archive)
        if err != nil {
            return ImportResult{Status: ImportFailed, Error: err.Error()}
        }
        archive.Path = managedPath
        archive.Managed = true
    }
    
    // Register archive
    if err := registry.Add(archive); err != nil {
        return ImportResult{Status: ImportFailed, Error: err.Error()}
    }
    
    return ImportResult{Status: ImportSuccess, Archive: archive}
}
```

#### 3. Batch Operations Framework
```go
// cmd/mas_batch.go
type BatchOperation string
const (
    BatchVerify    BatchOperation = "verify"
    BatchMove      BatchOperation = "move"
    BatchReprofile BatchOperation = "reprofile"
    BatchUpdate    BatchOperation = "update"
)

type BatchOptions struct {
    Filter      ListFilters  // Reuse list filtering from 7EP-0004
    Operation   BatchOperation
    Parameters  map[string]string
    DryRun      bool
    Parallel    int // Concurrent operations
    Continue    bool // Continue on errors
}

func runMasBatch(cmd *cobra.Command, args []string) error {
    // Apply filters to find target archives
    archives, err := registry.ListWithFilters(options.Filter)
    if err != nil {
        return err
    }
    
    fmt.Printf("Found %d archives matching criteria\n", len(archives))
    
    if options.DryRun {
        fmt.Printf("DRY RUN: Would perform '%s' operation on:\n", options.Operation)
        for _, archive := range archives {
            fmt.Printf("  - %s (%s)\n", archive.Name, archive.UID[:8])
        }
        return nil
    }
    
    // Execute batch operation
    return executeBatchOperation(archives, options)
}

func executeBatchOperation(archives []*storage.Archive, opts BatchOptions) error {
    // Setup worker pool for parallel processing
    jobs := make(chan *storage.Archive, len(archives))
    results := make(chan BatchResult, len(archives))
    
    // Start workers
    for w := 0; w < opts.Parallel; w++ {
        go batchWorker(jobs, results, opts)
    }
    
    // Send jobs
    for _, archive := range archives {
        jobs <- archive
    }
    close(jobs)
    
    // Collect results
    var success, failed int
    for i := 0; i < len(archives); i++ {
        result := <-results
        if result.Error != nil {
            failed++
            fmt.Printf("✗ %s: %v\n", result.Archive.Name, result.Error)
            if !opts.Continue {
                return result.Error
            }
        } else {
            success++
            fmt.Printf("✓ %s: %s\n", result.Archive.Name, result.Message)
        }
    }
    
    fmt.Printf("Batch operation complete: %d success, %d failed\n", success, failed)
    return nil
}
```

#### 4. Registry Maintenance Operations
```go
// cmd/mas_verify.go
type VerifyOptions struct {
    All         bool     // Verify all archives
    Pattern     string   // Verify archives matching pattern
    Checksums   bool     // Verify file checksums
    Registry    bool     // Verify registry consistency
    Repair      bool     // Attempt automatic repairs
}

func runMasVerify(cmd *cobra.Command, args []string) error {
    var archives []*storage.Archive
    var err error
    
    if options.All {
        archives, err = registry.List()
    } else if len(args) > 0 {
        resolver := storage.NewResolver(registry)
        archives, err = resolver.ResolvePattern(args[0])
    } else {
        return fmt.Errorf("specify --all or archive pattern")
    }
    
    if err != nil {
        return err
    }
    
    fmt.Printf("Verifying %d archives...\n", len(archives))
    
    var issues []VerificationIssue
    
    for _, archive := range archives {
        archiveIssues := verifyArchive(archive, options)
        issues = append(issues, archiveIssues...)
    }
    
    // Report results
    if len(issues) == 0 {
        fmt.Printf("✓ All archives verified successfully\n")
        return nil
    }
    
    fmt.Printf("Found %d issues:\n", len(issues))
    for _, issue := range issues {
        fmt.Printf("  %s: %s (%s)\n", issue.Archive.Name, issue.Description, issue.Severity)
    }
    
    // Attempt repairs if requested
    if options.Repair {
        return attemptRepairs(issues)
    }
    
    return nil
}
```

## Implementation Plan

### Phase 1: Move Operations Enhancement
- [ ] **Enhanced Move Command**
  - [ ] Multi-target move operations
  - [ ] Pattern-based moves (`--pattern "project-*"`)
  - [ ] External/managed storage transitions
  - [ ] Move verification and rollback

- [ ] **Path Resolution**
  - [ ] Smart destination path resolution
  - [ ] Conflict detection and handling  
  - [ ] Directory creation as needed
  - [ ] Permission and space validation

### Phase 2: Import Workflow
- [ ] **Archive Discovery**
  - [ ] File pattern matching and filtering
  - [ ] Recursive directory scanning
  - [ ] Duplicate detection (by checksum)
  - [ ] Archive format validation

- [ ] **Import Processing**
  - [ ] Batch import with progress reporting
  - [ ] Managed vs external import modes
  - [ ] Profile assignment and detection
  - [ ] Interactive import with user prompts

### Phase 3: Batch Operations
- [ ] **Batch Framework**
  - [ ] Filter-based archive selection
  - [ ] Parallel processing with worker pools
  - [ ] Progress reporting and cancellation
  - [ ] Error handling and continuation strategies

- [ ] **Core Batch Operations**
  - [ ] Mass verification and checksum validation
  - [ ] Bulk move and reorganization
  - [ ] Profile reassignment and optimization
  - [ ] Metadata updates and corrections

### Phase 4: Registry Maintenance
- [ ] **Verification System**
  - [ ] File existence and accessibility checks
  - [ ] Checksum validation and corruption detection
  - [ ] Registry consistency verification
  - [ ] Missing file detection and reporting

- [ ] **Repair and Maintenance**
  - [ ] Automatic issue detection and repair
  - [ ] Registry cleanup and optimization
  - [ ] Orphaned file detection and cleanup
  - [ ] Database integrity maintenance

## Success Criteria

### Move Operations
- [ ] Can move archives between any storage locations
- [ ] Pattern-based moves work reliably for batch operations
- [ ] Registry consistency maintained during all move operations
- [ ] Rollback capability for failed moves

### Import Workflow  
- [ ] Can import existing archive collections efficiently
- [ ] Duplicate detection prevents registry pollution
- [ ] Batch imports process 100+ files reliably
- [ ] Both managed and external import modes work correctly

### Batch Operations
- [ ] Can process 1000+ archives in batch operations
- [ ] Parallel processing improves performance significantly
- [ ] Error handling allows continuation of batch jobs
- [ ] Progress reporting provides clear status

### Registry Maintenance
- [ ] Verification detects all common registry inconsistencies
- [ ] Repair operations resolve issues automatically where possible
- [ ] Performance remains good even with 10K+ archives
- [ ] Database integrity maintained under all operations

## Related Work

### Builds On
- **7EP-0004 (Completed)**: MAS Foundation provides ULID resolution, registry operations
- **7EP-0006 (Completed)**: Performance validation ensures operations scale
- **Existing Commands**: Current move, delete commands provide foundation

### Integrates With
- **7EP-0001**: Trash management needs move operations for restore functionality
- **7EP-0005**: Test scenarios will validate batch operations and edge cases
- **7EP-0002**: CI integration will test advanced operations automatically

### Enables
- **Complete Archive Lifecycle**: Full CRUD operations with advanced workflows
- **Enterprise Usage**: Batch operations and maintenance for large deployments
- **Migration Workflows**: Import existing archives and reorganize storage
- **Operations Reliability**: Verification and repair for production use

## Future Considerations

- **Storage Backend Abstraction**: Support for S3, Azure, etc.
- **Archive Deduplication**: Identify and manage duplicate archives
- **Incremental Operations**: Resume interrupted batch operations
- **Audit Logging**: Track all operations for compliance
- **Plugin Architecture**: Extensible batch operations

---

This enhanced MAS operations suite transforms 7zarch-go from basic archive management into a comprehensive, enterprise-ready archive lifecycle management system.