# 7EP-0005: Comprehensive Test Dataset System

**Status:** Draft  
**Author(s):** Claude Code (CC)  
**Assignment:** CC (Primary), AC (Supporting)  
**Difficulty:** 3 (moderate - systematic but well-defined scope)  
**Created:** 2025-08-12  
**Updated:** 2025-08-12  
**Revised:** 2025-08-12 (based on 7EP-0006 implementation learnings)  

## Executive Summary

Implement a comprehensive, scalable test dataset system to enable systematic performance testing, realistic integration scenarios, and proper test organization. Replace ad hoc test files with a structured approach that supports everything from unit tests to large-scale performance benchmarks.

## Evidence & Reasoning

**Current Problems:**
- **Disorganized Test Files**: `demo-files/` and `friends-demo/` scattered in project root
- **No Performance Test Data**: 7EP-0004 requires testing with 1000+ archives, no systematic way to generate this
- **Limited Test Scenarios**: Existing test helpers create minimal archives, not realistic diversity
- **Inconsistent Test Data**: Each test creates its own data, no reusable patterns
- **Missing Edge Cases**: No systematic coverage of Unicode names, large files, deep hierarchies

**Impact on Development:**
- **Integration Testing Limited**: Real-world scenarios not covered beyond basic performance
- **Edge Case Coverage**: Systematic testing of Unicode, large files, deep hierarchies needed
- **Developer Experience**: Hard to reproduce complex bugs, no standard diverse test datasets
- **CI/CD Enhancement**: Need comprehensive test scenarios for robust automated testing
- **Future 7EP Support**: Upcoming features (trash management, migrations) need sophisticated test data

**Why Now (Post 7EP-0006 Insights):**
- **7EP-0006 Success**: Minimal approach proved metadata-only testing is highly effective
- **Performance Foundation Solid**: Core performance validated, can focus on comprehensive scenarios
- **Proven Patterns**: 7EP-0006 established working patterns for dataset generation and benchmarking
- **Growing Feature Set**: With MAS Foundation complete, more complex features need richer test data
- **Quality Assurance**: Need systematic edge case coverage that 7EP-0006's minimal approach intentionally skipped

## Learnings from 7EP-0006 Implementation

**Key Insights from Minimal Performance Testing:**

**✅ What Worked Extremely Well:**
- **Metadata-Only Approach**: No actual files needed for performance testing - metadata generation is sufficient and fast
- **Fixed Seed Reproducibility**: `rand.New(rand.NewSource(42))` ensures identical results across runs
- **Controlled ULID Generation**: Strategic similarity patterns enable disambiguation testing without chaos
- **Incremental Complexity**: Simple patterns scale to complex scenarios (100 → 1K → 10K archives)
- **Registry Integration**: Direct `reg.Add()` calls work perfectly for populating test databases
- **Benchmark Suite Pattern**: Multiple sub-benchmarks in single function provides comprehensive coverage

**🔧 What Needed Refinement:**
- **ULID Collision Management**: Initial approach created too many duplicates, refined pattern worked well
- **Type Assertions**: Converting `testing.TB` to `*testing.T` required careful handling  
- **API Method Names**: Need to match actual implementation (`reg.Add()` not `reg.Register()`)
- **Error Type Handling**: Duplicate error definitions caused build conflicts

**🚀 Performance Insights:**
- **O(1) Database Performance**: Registry operations scale constantly, not linearly
- **Massive Performance Margins**: Actual performance 100-2,941x better than requirements
- **Minimal Resource Usage**: 10K archives processed with negligible memory overhead
- **Simple is Sufficient**: Complex content generation unnecessary for most test scenarios

## Technical Design

### Overview

Build on 7EP-0006's proven metadata-only approach to create a comprehensive, hierarchical test dataset system that generates reproducible, diverse archive collections for integration testing, edge case coverage, and quality assurance.

### Dataset Organization Structure (Revised)

**Simplified hierarchy based on 7EP-0006 learnings:**

```
test-datasets/
├── generators/                 # Dataset creation tools  
│   ├── generator.go           # Core generation engine (builds on 7EP-0006 patterns)
│   ├── scenarios.go           # Pre-defined test scenarios  
│   └── ulid.go                # ULID generation patterns (refined from 7EP-0006)
├── scenarios/                 # Named scenario implementations
│   ├── performance/           # 7EP-0006 style benchmarks
│   │   ├── resolution.go      # ULID resolution testing
│   │   ├── filtering.go       # List filtering scenarios
│   │   └── scaling.go         # Multi-size datasets
│   ├── integration/           # End-to-end workflow testing
│   │   ├── user_workflows.go  # Create→list→show→delete flows
│   │   ├── mixed_storage.go   # Managed + external scenarios
│   │   └── time_series.go     # Archives across time periods
│   └── edge_cases/            # Systematic edge case coverage
│       ├── unicode.go         # International characters, emojis
│       ├── boundaries.go      # Size limits, name lengths, deep paths
│       ├── corruption.go      # Metadata consistency, missing files
│       └── migration.go       # Database schema changes
└── datasets/                  # Generated dataset cache (optional)
    ├── README.md              # Usage documentation
    └── .gitignore             # Exclude large generated files
```

**Key Changes from Original Design:**
- **Metadata-First**: No complex content generation - follow 7EP-0006's successful metadata-only approach
- **Scenario-Driven**: Organize by test purpose rather than arbitrary size buckets
- **Build on Proven Patterns**: Extend 7EP-0006's ULID generation and registry integration
- **Practical File Structure**: Fewer directories, clearer organization, easier maintenance

### Core Generation System (Revised)

#### Dataset Generator Interface
```go
// test-datasets/generators/generator.go  
package generators

import (
    "math/rand"
    "testing"
    "time"
    "github.com/adamstac/7zarch-go/internal/storage"
)

// ScenarioSpec defines what kind of test scenario to generate
// Simplified from DatasetSpec based on 7EP-0006 learnings
type ScenarioSpec struct {
    Name         string
    Count        int                    // Number of archives
    ULIDPattern  ULIDPattern           // How to generate ULIDs  
    Profiles     []ProfileDistribution // Profile mix
    TimeSpread   time.Duration         // Spread creation times over period
    SizePattern  SizePattern           // Size distribution
    EdgeCases    []EdgeCase            // Special scenarios to include
}

type ULIDPattern string
const (
    ULIDUnique      ULIDPattern = "unique"      // All unique (7EP-0006 scaling tests)
    ULIDSimilar     ULIDPattern = "similar"     // Controlled similarity (7EP-0006 disambiguation)  
    ULIDCollisions  ULIDPattern = "collisions"  // Intentional prefix overlaps
)

type ProfileDistribution struct {
    Profile string
    Weight  float32 // 0.0-1.0
}

// Generator creates reproducible test scenarios using 7EP-0006 patterns
type Generator struct {
    seed     int64
    registry *storage.Registry
}

func NewGenerator(seed int64) *Generator {
    return &Generator{seed: seed}
}

// GenerateScenario creates archives based on scenario spec
// Builds directly on 7EP-0006's successful generateTestArchives() pattern
func (g *Generator) GenerateScenario(tb testing.TB, spec ScenarioSpec) []*storage.Archive {
    tb.Helper()
    
    rng := rand.New(rand.NewSource(g.seed))
    archives := make([]*storage.Archive, spec.Count)
    
    profiles := extractProfiles(spec.Profiles)
    sizes := generateSizeDistribution(spec.SizePattern)
    
    for i := 0; i < spec.Count; i++ {
        archive := &storage.Archive{
            UID:      g.generateUID(i, spec.ULIDPattern, rng),
            Name:     fmt.Sprintf("%s-%04d.7z", spec.Name, i),
            Path:     fmt.Sprintf("/tmp/test-%s-%04d.7z", spec.Name, i),
            Size:     selectSize(sizes, rng),
            Created:  g.generateCreationTime(spec.TimeSpread, rng),
            Profile:  selectProfile(profiles, rng),
            Managed:  rng.Float32() < 0.9, // 90% managed (7EP-0006 pattern)
            Status:   "present",
            Checksum: fmt.Sprintf("sha256:%032x", i),
        }
        
        // Apply edge case modifications
        g.applyEdgeCases(archive, spec.EdgeCases, i, rng)
        
        archives[i] = archive
    }
    
    return archives
}
```

#### ULID Generation Patterns (Refined from 7EP-0006)

```go
// test-datasets/generators/ulid.go
// Builds on 7EP-0006's successful generateSequentialUID() approach

func (g *Generator) generateUID(i int, pattern ULIDPattern, rng *rand.Rand) string {
    switch pattern {
    case ULIDUnique:
        // Each ULID is completely unique - for scaling tests
        return fmt.Sprintf("01K%02d%05d%012d%05d", 
            rng.Intn(26), i, rng.Int63n(999999999999), rng.Intn(99999))
            
    case ULIDSimilar:
        // Controlled similarity for disambiguation testing (7EP-0006 pattern)
        if i < 100 {
            // Group by tens for first 100: 01K2E00, 01K2E01, 01K2E02, etc.
            return fmt.Sprintf("01K2E%02d%012d%08d", i/10, i, i*17)
        }
        return fmt.Sprintf("01K2F%02d%012d%08d", (i-100)/100, i, i*23)
        
    case ULIDCollisions:
        // Intentional prefix collisions for stress testing resolution
        prefixCount := min(i/10 + 1, 5) // Group into 5 prefix buckets max
        return fmt.Sprintf("01K2G%02d%012d%08d", prefixCount, i, rng.Intn(99999999))
        
    default:
        return ULIDUnique // Safe default
    }
}
```

#### Size and Time Distribution

```go
// test-datasets/generators/scenarios.go
// Simplified approach - metadata characteristics only

type SizePattern string
const (
    SizeUniform     SizePattern = "uniform"     // Even distribution
    SizeRealistic   SizePattern = "realistic"   // Log-normal (most files small, few large)
    SizeLargeFiles  SizePattern = "large"       // Focus on large file scenarios
)

func generateSizeDistribution(pattern SizePattern) []int64 {
    switch pattern {
    case SizeUniform:
        return []int64{1024, 100*1024, 10*1024*1024} // 1KB, 100KB, 10MB
    case SizeRealistic:
        // 70% small, 25% medium, 5% large (realistic distribution)
        return []int64{1024, 1024, 1024, 1024, 1024, 1024, 1024, // 70%
                      100*1024, 100*1024, 10*1024*1024}           // 25%, 5%
    case SizeLargeFiles:
        return []int64{100*1024*1024, 500*1024*1024, 1024*1024*1024} // 100MB, 500MB, 1GB
    default:
        return []int64{1024, 100*1024, 10*1024*1024} // Safe default
    }
}

### Predefined Test Scenarios (Revised)

#### Performance Testing Scenarios (Building on 7EP-0006)
```go
// test-datasets/scenarios/performance/resolution.go
var ResolutionScenarios = []ScenarioSpec{
    {
        Name: "disambiguation-stress",
        Count: 1000,
        ULIDPattern: ULIDSimilar, // Creates controlled similarity like 7EP-0006
        Profiles: []ProfileDistribution{
            {Profile: "documents", Weight: 0.7},
            {Profile: "media", Weight: 0.2}, 
            {Profile: "balanced", Weight: 0.1},
        },
        SizePattern: SizeUniform,
        TimeSpread: 30 * 24 * time.Hour, // 30 days
    },
    {
        Name: "scaling-validation",
        Count: 10000,
        ULIDPattern: ULIDUnique, // Each unique for pure scaling test
        Profiles: []ProfileDistribution{
            {Profile: "documents", Weight: 0.4},
            {Profile: "media", Weight: 0.3}, 
            {Profile: "balanced", Weight: 0.3},
        },
        SizePattern: SizeRealistic, // Most small, few large
        TimeSpread: 365 * 24 * time.Hour, // Full year spread
    },
}
```

#### Integration Testing Scenarios  
```go  
// test-datasets/scenarios/integration/user_workflows.go
var WorkflowScenarios = []ScenarioSpec{
    {
        Name: "create-list-show-delete",
        Count: 50,
        ULIDPattern: ULIDUnique,
        Profiles: []ProfileDistribution{
            {Profile: "documents", Weight: 0.5},
            {Profile: "media", Weight: 0.3},
            {Profile: "balanced", Weight: 0.2},
        },
        SizePattern: SizeRealistic,
        TimeSpread: 7 * 24 * time.Hour, // Week timeline
        EdgeCases: []EdgeCase{MixedManagedExternal, TimeSequencing},
    },
    {
        Name: "mixed-storage-scenario",
        Count: 100, 
        ULIDPattern: ULIDSimilar,
        Profiles: []ProfileDistribution{{Profile: "balanced", Weight: 1.0}},
        SizePattern: SizeUniform,
        EdgeCases: []EdgeCase{ManagedExternalMix, CrossProfileFiltering},
    },
}
```

#### Edge Case Testing Scenarios
```go
// test-datasets/scenarios/edge_cases/unicode.go
var EdgeCaseScenarios = []ScenarioSpec{
    {
        Name: "unicode-names",
        Count: 25,
        ULIDPattern: ULIDUnique,
        Profiles: []ProfileDistribution{{Profile: "documents", Weight: 1.0}},
        SizePattern: SizeUniform,  
        EdgeCases: []EdgeCase{
            UnicodeFilenames,    // 测试文件.7z, файл.7z, ファイル.7z
            EmojiFilenames,      // 🚀project.7z, 📝notes.7z
            SpecialCharacters,   // file with spaces.7z, file[brackets].7z
            LongFilenames,       // 255-character filenames
        },
    },
    {
        Name: "boundary-conditions",
        Count: 30,
        ULIDPattern: ULIDCollisions, // Stress test resolution
        Profiles: []ProfileDistribution{{Profile: "balanced", Weight: 1.0}},
        SizePattern: SizeLargeFiles, // Large file edge cases
        EdgeCases: []EdgeCase{
            MaxPathLength,       // Deep directory hierarchies
            MinMaxFileSizes,     // 0 byte and multi-GB files
            TimeBoundaries,      // Unix epoch, far future dates
        },
    },
}
```

### Integration with Existing Test Infrastructure (Based on 7EP-0006 Success)

#### Enhanced Test Helpers  
```go
// internal/storage/test_helpers.go - Enhanced with 7EP-0005 scenarios

// TestRegistryWithScenario creates registry with predefined scenario
// Builds directly on 7EP-0006's successful setupRegistryWithArchives pattern
func TestRegistryWithScenario(tb testing.TB, scenarioName string) (*storage.Registry, []*storage.Archive) {
    tb.Helper()
    
    // Create registry using 7EP-0006 pattern
    tmpDir := tb.TempDir() 
    dbPath := fmt.Sprintf("%s/test.db", tmpDir)
    reg, err := storage.NewRegistry(dbPath)
    if err != nil {
        tb.Fatalf("Failed to create test registry: %v", err)
    }
    
    tb.Cleanup(func() { reg.Close() })
    
    // Generate archives using scenario system
    generator := generators.NewGenerator(42) // Fixed seed like 7EP-0006
    spec := generators.GetScenario(scenarioName)
    archives := generator.GenerateScenario(tb, spec)
    
    // Add archives using proven pattern
    for _, archive := range archives {
        if err := reg.Add(archive); err != nil {
            tb.Fatalf("Failed to add test archive: %v", err)
        }
    }
    
    return reg, archives
}

// Benchmark helpers extending 7EP-0006 approach
func BenchmarkWithScenario(b *testing.B, scenarioName string, 
                           benchmarkFn func(*storage.Registry, []*storage.Archive)) {
    reg, archives := TestRegistryWithScenario(b, scenarioName)
    
    b.ResetTimer()
    benchmarkFn(reg, archives) 
}

// Example usage - extending 7EP-0006's benchmark patterns
func BenchmarkResolutionWithEdgeCases(b *testing.B) {
    scenarios := []string{
        "disambiguation-stress",
        "unicode-names", 
        "boundary-conditions",
    }
    
    for _, scenario := range scenarios {
        b.Run(scenario, func(b *testing.B) {
            BenchmarkWithScenario(b, scenario, func(reg *storage.Registry, archives []*storage.Archive) {
                resolver := storage.NewResolver(reg)
                resolver.MinPrefixLength = 4 // From 7EP-0006 learning
                
                for i := 0; i < b.N; i++ {
                    archive := archives[i%len(archives)]
                    _, err := resolver.Resolve(archive.UID[:8])
                    if err != nil && !isAmbiguousError(err) { // 7EP-0006 pattern
                        b.Fatalf("Resolution failed: %v", err)
                    }
                }
            })
        })
    }
}
```

### Performance Testing Integration

#### Benchmark Suite
```go
// test-datasets/benchmarks/mas_foundation_bench_test.go
func BenchmarkULIDResolution(b *testing.B) {
    testCases := []struct {
        name string
        size DatasetSize
    }{
        {"100_archives", SizeSmall},
        {"1k_archives", SizeMedium}, 
        {"10k_archives", SizeLarge},
    }
    
    for _, tc := range testCases {
        b.Run(tc.name, func(b *testing.B) {
            BenchmarkResolver(b, tc.size)
        })
    }
}

func BenchmarkListFiltering(b *testing.B) {
    reg, _ := BenchmarkRegistryWithDataset(b, "list-filtering-10k")
    
    testCases := []struct {
        name   string
        filter ListFilters
    }{
        {"status_present", ListFilters{Status: "present"}},
        {"profile_media", ListFilters{Profile: "media"}},
        {"larger_than_100mb", ListFilters{LargerThan: 100 * 1024 * 1024}},
        {"older_than_30d", ListFilters{OlderThan: 30 * 24 * time.Hour}},
        {"complex_combination", ListFilters{
            Profile: "media", 
            LargerThan: 10*1024*1024, 
            Status: "present",
        }},
    }
    
    for _, tc := range testCases {
        b.Run(tc.name, func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                archives, err := reg.ListWithFilters(tc.filter)
                if err != nil {
                    b.Fatal(err)
                }
                _ = archives // Use result
            }
        })
    }
}
```

### Dataset Management

#### Lifecycle Management
```go
type Dataset struct {
    Name     string
    Archives []*TestArchive
    BasePath string
    Metadata DatasetMetadata
}

type DatasetMetadata struct {
    Created     time.Time
    Spec        DatasetSpec
    TotalSize   int64
    FileCount   int
    Checksums   map[string]string // For integrity verification
}

// Cleanup removes all generated files and registry entries
func (d *Dataset) Cleanup() error {
    return os.RemoveAll(d.BasePath)
}

// Verify checks dataset integrity
func (d *Dataset) Verify() error {
    for _, archive := range d.Archives {
        if err := archive.Verify(); err != nil {
            return fmt.Errorf("archive %s failed verification: %w", archive.Name, err)
        }
    }
    return nil
}
```

### Edge Case Testing

#### Special Scenarios
```go
// test-datasets/generators/edge_cases.go
func GenerateUnicodeDataset() DatasetSpec {
    return DatasetSpec{
        Name: "unicode-names",
        Size: SizeMicro,
        EdgeCases: []EdgeCase{
            UnicodeNames,      // 测试文件.7z, файл.7z, ファイル.7z
            EmojiNames,        // 🚀project.7z, 📝notes.7z
            SpecialChars,      // file with spaces.7z, file[brackets].7z
            LongNames,         // 255-character filenames
        },
    }
}

func GenerateCorruptionDataset() DatasetSpec {
    return DatasetSpec{
        Name: "corruption-scenarios",
        Size: SizeSmall,
        EdgeCases: []EdgeCase{
            CorruptedHeaders,   // Invalid 7z headers
            CorruptedChecksums, // Wrong SHA-256 values
            TruncatedFiles,     // Files cut off mid-stream
            MissingFiles,       // Registry entries with no files
        },
    }
}
```

## Implementation Plan (Revised Based on 7EP-0006 Learnings)

### Phase 1: Core Generator Infrastructure (Building on 7EP-0006)
- [ ] **Scenario Generator Framework**
  - [ ] ScenarioSpec definition (simplified from DatasetSpec)
  - [ ] ULID generation patterns (refined from 7EP-0006's successful approach)
  - [ ] Registry integration using proven `reg.Add()` pattern
  - [ ] Reproducible seeding with fixed seeds (42 default)

- [ ] **Metadata Generation Patterns**
  - [ ] Size distribution patterns (uniform, realistic, large-files)
  - [ ] Profile distribution system with weights
  - [ ] Time spread generation across configurable periods
  - [ ] Edge case modification system

### Phase 2: Core Scenarios (Extending 7EP-0006 Success)
- [ ] **Performance Scenario Library**
  - [ ] Resolution disambiguation scenarios (builds on 7EP-0006's ULIDSimilar)
  - [ ] Scaling validation scenarios (extends 7EP-0006's scaling tests)
  - [ ] List filtering stress scenarios (10K+ archives with complex filters)
  - [ ] Show command verification scenarios

- [ ] **Integration Test Helpers**
  - [ ] TestRegistryWithScenario helper (extends setupRegistryWithArchives pattern)
  - [ ] BenchmarkWithScenario helper (generalizes 7EP-0006's benchmark pattern)
  - [ ] Scenario-based test assertions (extends AssertResolves pattern)

### Phase 3: Edge Case & Integration Coverage
- [ ] **Edge Case Scenarios**
  - [ ] Unicode filename handling (systematic coverage beyond 7EP-0006)
  - [ ] Boundary condition testing (path lengths, file sizes)
  - [ ] Metadata corruption and consistency scenarios
  - [ ] Cross-platform compatibility scenarios

- [ ] **Workflow Integration Scenarios**
  - [ ] End-to-end user workflow scenarios (create→list→show→delete)
  - [ ] Mixed managed/external storage scenarios
  - [ ] Time-series archive management scenarios
  - [ ] Migration and schema change scenarios

### Phase 4: Quality Assurance & Maintenance
- [ ] **Test Organization Cleanup**
  - [ ] Migrate demo-files/ to structured scenarios
  - [ ] Replace ad hoc test data with systematic scenarios
  - [ ] Documentation for scenario usage patterns
  - [ ] Developer guides for adding new scenarios

- [ ] **Advanced Features**
  - [ ] Scenario composition (combine multiple edge cases)
  - [ ] Performance regression monitoring
  - [ ] CI integration with systematic test coverage
  - [ ] Scenario versioning for compatibility testing

## Success Criteria (Updated with 7EP-0006 Insights)

### Performance Validation (Building on 7EP-0006 Success) 
- [x] **7EP-0006 Baseline Established**: Metadata-only approach proven effective
- [ ] Can generate 10K+ archive scenarios in <5 seconds (metadata-only, no files)
- [ ] Edge case scenarios validate resolution robustness beyond basic performance
- [ ] Integration scenarios validate end-to-end workflow performance
- [ ] Memory usage remains <10MB even with complex scenarios

**Note**: Core performance requirements already validated by 7EP-0006. Focus shifts to comprehensive scenario coverage.

### Test Coverage
- [ ] All major file types represented (documents, media, code, data)
- [ ] Size distribution covers tiny (1KB) to huge (1GB+) files
- [ ] Edge cases systematically covered (Unicode, corruption, etc.)
- [ ] Realistic usage patterns modeled
- [ ] Both unit and integration test scenarios supported

### Developer Experience
- [ ] Simple API for creating test datasets in any test
- [ ] Reproducible datasets (same seed = same data)
- [ ] Fast dataset generation (<5s for small datasets)
- [ ] Automatic cleanup with no manual intervention
- [ ] Clear documentation and examples

## Migration Strategy

### Current Test File Organization
```bash
# Current (messy)
7zarch-go/
├── demo-files/
├── friends-demo/
└── various test files scattered

# Target (organized)
7zarch-go/
├── test-datasets/
│   ├── generators/
│   ├── micro/small/medium/large/
│   └── profiles/edge-cases/
└── clean root directory
```

### Migration Steps
1. **Create new structure** alongside existing files
2. **Migrate existing tests** to use dataset system
3. **Generate replacement datasets** for demo-files/friends-demo
4. **Update CI/CD** to use new dataset paths
5. **Remove old files** once fully migrated

## Testing Strategy

### Generator Testing
- Verify reproducibility (same seed → same dataset)
- Test content type generators produce valid files
- Validate size and timing distributions
- Check edge case coverage

### Integration Testing
- Test with existing test helpers and registry system
- Validate performance benchmarks run correctly
- Ensure cleanup works properly
- Test cross-platform compatibility

### Performance Testing
- Benchmark dataset generation itself
- Validate that generated datasets meet performance requirements
- Test with various dataset sizes and complexity

## Related Work (Updated with 7EP-0006 Relationship)

### Builds On (Now Including 7EP-0006 Success)
- **7EP-0006 (Completed)**: Proven metadata-only approach for performance testing
- **7EP-0004 (Completed)**: MAS Foundation provides solid performance baseline
- **Test Helpers (PR #6)**: Existing test infrastructure and patterns
- **Integration Tests**: Scenarios from mas-foundation-integration.md

### Relationship to 7EP-0006
- **Complementary, Not Competing**: 7EP-0006 solved immediate performance validation needs
- **Extension Strategy**: 7EP-0005 extends 7EP-0006's successful patterns to comprehensive scenarios  
- **Proven Foundation**: Builds on 7EP-0006's working ULID generation, registry integration
- **Different Purpose**: 7EP-0006 = minimal performance validation, 7EP-0005 = comprehensive test infrastructure

### Enables (Post 7EP-0006)
- **Future 7EP Testing**: Rich test scenarios for trash management (7EP-0001), CI (7EP-0002), migrations (7EP-0003)
- **Edge Case Coverage**: Systematic testing of Unicode, boundaries, corruption beyond basic performance
- **Integration Validation**: End-to-end workflow testing with realistic archive diversity
- **Quality Assurance**: Regression testing and comprehensive scenario coverage

## Future Considerations

- **Fuzzing Integration**: Use datasets as basis for property-based testing
- **Cloud Storage Testing**: Remote file scenarios with simulated latency
- **Concurrency Testing**: Multi-user scenarios with shared datasets
- **Migration Testing**: Dataset evolution and backward compatibility
- **Performance Profiling**: Detailed bottleneck analysis with realistic data

## Summary of Key Changes (Post 7EP-0006)

**Major Revisions Based on Implementation Learnings:**

1. **Metadata-First Architecture**: Eliminates complex content generation in favor of 7EP-0006's proven metadata-only approach
2. **Scenario-Driven Organization**: Replaces arbitrary size buckets with purpose-driven scenarios (performance, integration, edge-cases)
3. **Proven Pattern Extension**: Builds directly on 7EP-0006's successful ULID generation, registry integration, and benchmark patterns
4. **Simplified API**: ScenarioSpec replaces complex DatasetSpec, reducing implementation complexity
5. **Focused Scope**: Shifts from "comprehensive system" to "systematic edge case and integration coverage"
6. **Practical File Structure**: Streamlined directory hierarchy, easier maintenance and understanding

**Strategic Position**: 7EP-0005 now serves as the logical evolution of 7EP-0006's success, focusing on areas where comprehensive coverage adds value beyond basic performance validation.

---

This refined test dataset system builds on 7EP-0006's proven approach to provide systematic edge case coverage, integration testing, and quality assurance while maintaining the simplicity and effectiveness of the metadata-only pattern.