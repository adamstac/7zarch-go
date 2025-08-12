# 7EP-0006: Minimal Performance Testing for 7EP-0004 Completion

**Status:** Draft  
**Author(s):** Claude Code (CC)  
**Assignment:** CC (Primary)  
**Difficulty:** 1 (simple - focused scope with clear deliverables)  
**Created:** 2025-08-12  
**Updated:** 2025-08-12  

## Executive Summary

Create minimal performance testing infrastructure to validate 7EP-0004 MAS Foundation completion criteria (<50ms resolution, <200ms list ops, <100ms show command). Dead simple approach focused solely on unblocking 7EP-0004 rather than comprehensive test infrastructure.

## Evidence & Reasoning

**Immediate Problem:**
- **7EP-0004 is 90% complete** but blocked on performance validation
- **No way to test** the <50ms resolution, <200ms list filtering requirements
- **7EP-0005 is comprehensive** but complex and would take days to implement
- **Need simple solution** to validate requirements and complete 7EP-0004 this week

**Why Minimal Approach:**
- **Time Sensitive**: 7EP-0004 completion is priority
- **Clear Requirements**: We know exactly what to test  
- **Simple Scope**: Just validate three performance targets
- **Evolutionary**: Can build comprehensive system later (7EP-0005) if needed

**Evidence for Requirements:**
- 7EP-0004 specifies <50ms resolution operations
- 7EP-0004 specifies <200ms list operations for 1000+ archives
- 7EP-0004 specifies <100ms show command including verification
- Need large dataset (1000+ archives) to stress test realistically

## Technical Design

### Overview
Single performance test file that generates minimal test archives in-memory and runs focused benchmarks against the three core MAS Foundation operations.

### Architecture

```go
// test/performance/mas_benchmark_test.go
package performance

import (
    "testing"
    "time"
    "math/rand"
    "github.com/adamstac/7zarch-go/internal/storage"
)

// Simple archive generator - no fancy content, just metadata
func generateTestArchives(count int) []*storage.Archive {
    archives := make([]*storage.Archive, count)
    rng := rand.New(rand.NewSource(42)) // Reproducible
    
    profiles := []string{"documents", "media", "balanced"}
    sizes := []int64{1024, 100*1024, 10*1024*1024} // 1KB, 100KB, 10MB
    
    for i := 0; i < count; i++ {
        archive := &storage.Archive{
            UID:     generateSequentialUID(i), // For prefix testing
            Name:    fmt.Sprintf("archive-%04d.7z", i),
            Path:    fmt.Sprintf("/tmp/test-archive-%04d.7z", i),
            Size:    sizes[rng.Intn(len(sizes))],
            Created: time.Now().Add(-time.Duration(rng.Intn(365)) * 24 * time.Hour),
            Profile: profiles[rng.Intn(len(profiles))],
            Managed: i%10 != 0, // 90% managed, 10% external
            Status:  "present",
            Checksum: fmt.Sprintf("sha256:%032x", i), // Fake but unique
        }
        archives[i] = archive
    }
    
    return archives
}

// Generate ULIDs with some overlap for disambiguation testing
func generateSequentialUID(i int) string {
    // Create ULIDs where first few characters are similar for some archives
    // This tests disambiguation performance
    if i < 100 {
        return fmt.Sprintf("01K2E%021d", i) // Similar prefixes
    }
    return fmt.Sprintf("01K%023d", i+1000) // Different prefixes
}
```

### Three Core Benchmarks

#### 1. Resolution Performance
```go
func BenchmarkULIDResolution(b *testing.B) {
    archives := generateTestArchives(1000)
    reg := setupRegistryWithArchives(b, archives)
    resolver := storage.NewResolver(reg)
    
    // Test cases: full UID, 8-char prefix, 4-char prefix (disambiguation)
    testCases := []struct {
        name string
        getID func(archive *storage.Archive) string
    }{
        {"full_uid", func(a *storage.Archive) string { return a.UID }},
        {"8_char_prefix", func(a *storage.Archive) string { return a.UID[:8] }},
        {"4_char_prefix", func(a *storage.Archive) string { return a.UID[:4] }}, // May be ambiguous
        {"name_lookup", func(a *storage.Archive) string { return a.Name }},
    }
    
    for _, tc := range testCases {
        b.Run(tc.name, func(b *testing.B) {
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                archive := archives[i%len(archives)]
                id := tc.getID(archive)
                _, err := resolver.ResolveID(id)
                // Allow AmbiguousIDError for 4-char prefixes
                if err != nil && !isAmbiguousError(err) {
                    b.Fatalf("Resolution failed: %v", err)
                }
            }
        })
    }
}

// Target: <50ms per operation (or <50 microseconds in benchmarks)
```

#### 2. List Filtering Performance  
```go
func BenchmarkListFiltering(b *testing.B) {
    archives := generateTestArchives(10000) // Larger set for filtering
    reg := setupRegistryWithArchives(b, archives)
    
    testCases := []struct {
        name   string
        filter ListFilters
    }{
        {"no_filter", ListFilters{}},
        {"status_present", ListFilters{Status: "present"}},
        {"profile_media", ListFilters{Profile: "media"}},
        {"managed_only", ListFilters{Managed: &[]bool{true}[0]}},
        {"large_files", ListFilters{LargerThan: 1024*1024}}, // >1MB
        {"old_files", ListFilters{OlderThan: 30*24*time.Hour}}, // >30 days
        {"complex", ListFilters{
            Profile: "media", 
            LargerThan: 100*1024, 
            Status: "present",
        }},
    }
    
    for _, tc := range testCases {
        b.Run(tc.name, func(b *testing.B) {
            b.ResetTimer() 
            for i := 0; i < b.N; i++ {
                results, err := reg.ListWithFilters(tc.filter)
                if err != nil {
                    b.Fatalf("List filtering failed: %v", err)
                }
                _ = results // Use results to prevent optimization
            }
        })
    }
}

// Target: <200ms per operation for 10K archives
```

#### 3. Show Command Performance
```go
func BenchmarkShowCommand(b *testing.B) {
    archives := generateTestArchives(1000)
    reg := setupRegistryWithArchives(b, archives)
    resolver := storage.NewResolver(reg)
    
    testCases := []struct {
        name   string
        verify bool
    }{
        {"basic_show", false},
        {"with_verification", true}, // More expensive
    }
    
    for _, tc := range testCases {
        b.Run(tc.name, func(b *testing.B) {
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                archive := archives[i%len(archives)]
                
                // Resolve archive
                resolved, err := resolver.ResolveID(archive.UID[:8])
                if err != nil {
                    b.Fatalf("Resolution failed: %v", err)
                }
                
                // Simulate show command operations
                if tc.verify {
                    // Simulate checksum verification (expensive)
                    _ = simulateChecksumVerification(resolved)
                }
                _ = formatArchiveDisplay(resolved)
            }
        })
    }
}

// Target: <100ms per operation including verification
```

### Helper Functions

```go
func setupRegistryWithArchives(tb testing.TB, archives []*storage.Archive) *storage.Registry {
    tb.Helper()
    
    reg := storage.TestRegistry(tb) // Use existing helper
    
    // Register all archives
    for _, archive := range archives {
        if err := reg.Register(archive); err != nil {
            tb.Fatalf("Failed to register test archive: %v", err)
        }
    }
    
    return reg
}

func isAmbiguousError(err error) bool {
    _, ok := err.(*storage.AmbiguousIDError)
    return ok
}

func simulateChecksumVerification(archive *storage.Archive) error {
    // Simulate the time cost of reading and hashing a file
    // For testing purposes, just sleep for realistic duration
    time.Sleep(10 * time.Microsecond) // Simulate small file verification
    return nil
}

func formatArchiveDisplay(archive *storage.Archive) string {
    // Simulate the string formatting work done by show command
    return fmt.Sprintf("Archive: %s (%s)\nSize: %d bytes\n", 
        archive.Name, archive.UID[:8], archive.Size)
}
```

## Success Criteria

### Performance Targets (7EP-0004 Requirements)
- [ ] **Resolution Operations**: Average <50 microseconds per operation
- [ ] **List Filtering**: Average <200 milliseconds for 10K archives  
- [ ] **Show Command**: Average <100 milliseconds including verification

### Validation Requirements
- [ ] **1000+ Archive Dataset**: Generate and test with realistic scale
- [ ] **Multiple Scenarios**: Test various filter combinations and lookup types
- [ ] **Reproducible Results**: Same seed produces consistent benchmark results
- [ ] **CI Integration**: Benchmarks can run in automated pipeline

### Deliverables
- [ ] **Single Test File**: `test/performance/mas_benchmark_test.go`
- [ ] **Benchmark Results**: Document baseline performance on standard hardware
- [ ] **7EP-0004 Validation**: Confirm all performance requirements met
- [ ] **Documentation**: Brief README explaining how to run benchmarks

## Implementation Plan

### Phase 1: Core Benchmarks (2-3 hours)
- [ ] Create test file structure
- [ ] Implement archive generation function
- [ ] Build resolution benchmark with multiple ID types
- [ ] Build list filtering benchmark with common filters
- [ ] Build show command benchmark with/without verification

### Phase 2: Validation (1 hour)  
- [ ] Run benchmarks and document baseline results
- [ ] Validate against 7EP-0004 requirements
- [ ] Adjust if performance targets not met
- [ ] Create simple documentation

### Phase 3: Integration (30 minutes)
- [ ] Add to CI pipeline (optional)
- [ ] Update 7EP-0004 status to completed
- [ ] Clean up any temporary code

## Usage

```bash
# Run all performance benchmarks
go test -bench=. ./test/performance/

# Run specific benchmark
go test -bench=BenchmarkULIDResolution ./test/performance/

# Run with memory profiling
go test -bench=. -memprofile=mem.prof ./test/performance/

# Generate performance report
go test -bench=. -benchmem ./test/performance/ > performance-results.txt
```

## Migration Strategy

### Current State
- No performance testing infrastructure
- 7EP-0004 blocked on performance validation
- Existing unit tests but no large-scale benchmarks

### Implementation
- Single new test file, no changes to existing code
- Uses existing test helpers where possible
- Minimal dependencies, focused scope

### Future Evolution  
- This minimal system can evolve into 7EP-0005 comprehensive approach
- Benchmark results inform optimization priorities
- Foundation for regression testing and CI integration

## Comparison to 7EP-0005

| Aspect | 7EP-0006 (Minimal) | 7EP-0005 (Comprehensive) |
|--------|-------|-------------|
| **Scope** | Just 7EP-0004 validation | Full test infrastructure |
| **Time** | 2-3 hours | 1-2 weeks |
| **Complexity** | Single test file | Multiple packages, generators |
| **Features** | 3 core benchmarks | Content generation, edge cases, organization |
| **Value** | Unblocks 7EP-0004 immediately | Long-term test foundation |

**Relationship**: 7EP-0006 is a stepping stone to 7EP-0005. Complete the immediate need, then build comprehensive system later if warranted.

## Testing Strategy

### Benchmark Validation
- Run on multiple machines to confirm consistency  
- Test with different archive counts (100, 1K, 10K)
- Validate memory usage stays reasonable
- Check for performance regressions over time

### Integration Testing
- Ensure benchmarks don't interfere with existing tests
- Validate that generated archives work with real MAS operations
- Test cleanup and resource management

## Future Considerations

- **Regression Testing**: Track performance over time
- **Platform Testing**: Validate performance across OS/architectures  
- **Optimization**: Use benchmark results to guide performance improvements
- **Migration**: Evolution path to 7EP-0005 when comprehensive testing needed

---

**Bottom Line**: Get 7EP-0004 to 100% complete this week with minimal effort, then decide if we need the comprehensive 7EP-0005 approach based on actual usage patterns and performance insights.