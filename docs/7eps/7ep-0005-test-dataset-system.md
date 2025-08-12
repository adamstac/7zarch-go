# 7EP-0005: Comprehensive Test Dataset System

**Status:** Draft  
**Author(s):** Claude Code (CC)  
**Assignment:** CC (Primary), AC (Supporting)  
**Difficulty:** 3 (moderate - systematic but well-defined scope)  
**Created:** 2025-08-12  
**Updated:** 2025-08-12  

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
- **Performance Testing Blocked**: Can't validate 7EP-0004 completion criteria (<50ms resolution, <200ms list ops)
- **Integration Testing Limited**: Real-world scenarios not covered
- **Flaky Tests**: Inconsistent test data leads to unreliable results  
- **Developer Experience**: Hard to reproduce issues, no standard test datasets
- **CI/CD Constraints**: No systematic approach to test data in automated pipelines

**Why Now:**
- **7EP-0004 Completion**: Performance testing is the main blocker for MAS Foundation completion
- **Growing Complexity**: More features need more sophisticated test scenarios
- **Test Infrastructure Exists**: Test helpers from PR #6 provide foundation to build on
- **Clean Slate Opportunity**: Early enough in project to establish good patterns

## Technical Design

### Overview

Build a hierarchical test dataset system that generates reproducible, diverse archive collections for comprehensive testing across performance, integration, and edge case scenarios.

### Dataset Organization Structure

```
test-datasets/
├── generators/                 # Dataset creation tools
│   ├── generator.go           # Core generation engine
│   ├── content/               # Content type generators
│   │   ├── documents.go       # Text, PDF, office files
│   │   ├── media.go           # Images, audio, video
│   │   ├── code.go            # Source trees, repositories
│   │   └── mixed.go           # Realistic combinations
│   └── scenarios/             # Specific test scenarios
│       ├── performance.go     # Large-scale benchmarks
│       ├── edge_cases.go      # Unicode, special chars
│       └── integration.go     # End-to-end workflows
├── micro/                     # 1-10 archives (unit tests)
├── small/                     # 10-100 archives (integration tests)
├── medium/                    # 100-1000 archives (stress tests)
├── large/                     # 1000-10000 archives (performance tests)
├── profiles/                  # Profile-specific datasets
│   ├── documents/             # Document-heavy scenarios
│   ├── media/                 # Media-heavy scenarios
│   └── balanced/              # Mixed content scenarios
└── edge-cases/                # Special scenarios
    ├── unicode/               # International characters
    ├── large-files/           # 100MB+ individual files
    ├── deep-nested/           # 20+ directory levels
    └── corruption/            # Intentionally corrupted data
```

### Core Generation System

#### Dataset Generator Interface
```go
// test-datasets/generators/generator.go
package generators

import (
    "context"
    "time"
)

// DatasetSpec defines what kind of dataset to generate
type DatasetSpec struct {
    Name          string
    Size          DatasetSize
    ContentMix    ContentMix
    TimeSpread    time.Duration  // Spread creation times over period
    SizeVariation SizeVariation  // File size distribution
    Structure     StructureType  // Directory organization
    EdgeCases     []EdgeCase     // Special scenarios to include
}

type DatasetSize string
const (
    SizeMicro  DatasetSize = "micro"   // 1-10 archives
    SizeSmall  DatasetSize = "small"   // 10-100 archives
    SizeMedium DatasetSize = "medium"  // 100-1000 archives
    SizeLarge  DatasetSize = "large"   // 1000-10000 archives
)

type ContentMix struct {
    Documents float32 // 0.0-1.0 percentage
    Media     float32
    Code      float32
    Data      float32
}

// Generator creates reproducible test datasets
type Generator struct {
    seed     int64
    basePath string
    registry *storage.Registry
}

func NewGenerator(seed int64, basePath string) *Generator {
    return &Generator{
        seed:     seed,
        basePath: basePath,
    }
}

// Generate creates a dataset according to spec
func (g *Generator) Generate(ctx context.Context, spec DatasetSpec) (*Dataset, error) {
    // Use seed for reproducible randomization
    rng := rand.New(rand.NewSource(g.seed))
    
    dataset := &Dataset{
        Name:     spec.Name,
        Archives: make([]*TestArchive, 0),
        BasePath: filepath.Join(g.basePath, spec.Name),
    }
    
    // Generate content based on spec
    archives := g.generateArchives(rng, spec)
    
    // Create actual files and register archives
    for _, archive := range archives {
        if err := g.createArchiveFiles(archive); err != nil {
            return nil, err
        }
        if err := g.registerArchive(archive); err != nil {
            return nil, err
        }
        dataset.Archives = append(dataset.Archives, archive)
    }
    
    return dataset, nil
}
```

#### Content Type Generators

**Document Generator:**
```go
// test-datasets/generators/content/documents.go
func GenerateDocumentContent(rng *rand.Rand, size int) []byte {
    switch {
    case size < 1024: // Small configs, notes
        return generateConfigFile(rng, size)
    case size < 100*1024: // Medium documents
        return generateTextDocument(rng, size)  
    default: // Large documents, presentations
        return generateStructuredDocument(rng, size)
    }
}

func generateTextDocument(rng *rand.Rand, targetSize int) []byte {
    // Generate realistic text content with:
    // - Lorem ipsum base text
    // - Realistic word/paragraph distribution
    // - Common document structures (headers, lists, etc.)
}
```

**Media Generator:**
```go
// test-datasets/generators/content/media.go
func GenerateImageContent(rng *rand.Rand, size int) []byte {
    // Generate synthetic image data:
    // - Realistic JPEG/PNG headers
    // - Appropriate compression ratios
    // - Various dimensions and color depths
}

func GenerateVideoContent(rng *rand.Rand, size int) []byte {
    // Generate synthetic video file:
    // - MP4 container with realistic headers
    // - Appropriate bitrates for size
    // - Minimal but valid video stream
}
```

### Predefined Dataset Scenarios

#### Performance Testing Datasets
```go
// test-datasets/generators/scenarios/performance.go
var PerformanceScenarios = []DatasetSpec{
    {
        Name: "resolution-benchmark-1k",
        Size: SizeMedium,
        ContentMix: ContentMix{Documents: 0.7, Code: 0.2, Data: 0.1},
        EdgeCases: []EdgeCase{SimilarULIDPrefixes}, // Test disambiguation
    },
    {
        Name: "list-filtering-10k",
        Size: SizeLarge, 
        ContentMix: ContentMix{Documents: 0.4, Media: 0.3, Code: 0.2, Data: 0.1},
        TimeSpread: 365 * 24 * time.Hour, // Year-long spread
        SizeVariation: LogNormalDistribution, // Realistic size distribution
    },
    {
        Name: "show-verification-mixed",
        Size: SizeMedium,
        ContentMix: ContentMix{Media: 0.8, Documents: 0.2}, // Large files
        EdgeCases: []EdgeCase{CorruptedChecksums, MissingFiles},
    },
}
```

#### Integration Testing Datasets
```go
// test-datasets/generators/scenarios/integration.go
var IntegrationScenarios = []DatasetSpec{
    {
        Name: "realistic-user-workflow",
        Size: SizeSmall,
        ContentMix: ContentMix{Documents: 0.5, Code: 0.3, Media: 0.2},
        Structure: MixedHierarchy,
        EdgeCases: []EdgeCase{UnicodeNames, LongPaths},
    },
    {
        Name: "media-heavy-workflow", 
        Size: SizeMedium,
        ContentMix: ContentMix{Media: 0.8, Documents: 0.2},
        SizeVariation: HighVariance, // Mix of tiny and huge files
    },
}
```

### Integration with Existing Test Infrastructure

#### Enhanced Test Helpers
```go
// internal/storage/test_helpers.go - Enhanced to use datasets

// TestRegistryWithDataset creates registry with predefined dataset
func TestRegistryWithDataset(t *testing.T, datasetName string) (*Registry, *Dataset) {
    t.Helper()
    
    reg := TestRegistry(t) // Existing helper
    
    generator := generators.NewGenerator(42, t.TempDir()) // Fixed seed
    spec := generators.GetScenario(datasetName)
    dataset, err := generator.Generate(context.Background(), spec)
    if err != nil {
        t.Fatalf("Failed to generate dataset %s: %v", datasetName, err)
    }
    
    // Register all archives with registry
    for _, archive := range dataset.Archives {
        if err := reg.Register(archive.Archive); err != nil {
            t.Fatalf("Failed to register test archive: %v", err)
        }
    }
    
    t.Cleanup(func() {
        dataset.Cleanup()
        reg.Close()
    })
    
    return reg, dataset
}

// Performance test helpers
func BenchmarkResolver(b *testing.B, datasetSize DatasetSize) {
    reg, dataset := BenchmarkRegistryWithDataset(b, string(datasetSize))
    resolver := storage.NewResolver(reg)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        // Test resolution with random archive from dataset
        archive := dataset.Archives[i%len(dataset.Archives)]
        _, err := resolver.ResolveID(archive.UID[:8])
        if err != nil {
            b.Fatalf("Resolution failed: %v", err)
        }
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

## Implementation Plan

### Phase 1: Core Generator Infrastructure
- [ ] **Dataset Generator Framework**
  - [ ] Base generator with reproducible seeding
  - [ ] DatasetSpec definition and parsing
  - [ ] File creation and registry integration
  - [ ] Cleanup and lifecycle management

- [ ] **Basic Content Generators**
  - [ ] Document content generator (text, structured)
  - [ ] Simple binary content generator
  - [ ] Directory structure creation
  - [ ] Size and timing distribution

### Phase 2: Performance Testing Datasets
- [ ] **Benchmark Dataset Creation**
  - [ ] 1K archive dataset for resolution testing
  - [ ] 10K archive dataset for list filtering
  - [ ] Large file dataset for show/verification
  - [ ] Mixed profile datasets for realistic scenarios

- [ ] **Benchmark Integration**
  - [ ] Enhanced test helpers using datasets
  - [ ] Benchmark suite for 7EP-0004 validation
  - [ ] Performance regression testing
  - [ ] CI integration for automated benchmarks

### Phase 3: Integration & Edge Case Testing
- [ ] **Realistic Integration Scenarios**
  - [ ] User workflow datasets (create→list→show→delete)
  - [ ] Mixed managed/external scenarios
  - [ ] Profile-specific usage patterns
  - [ ] Time-distributed archive creation

- [ ] **Edge Case Coverage**
  - [ ] Unicode and special character handling
  - [ ] Large file and deep hierarchy scenarios
  - [ ] Corruption and missing file scenarios
  - [ ] Network storage and permission scenarios

### Phase 4: Advanced Features
- [ ] **Dynamic Dataset Generation**
  - [ ] On-demand dataset creation for specific tests
  - [ ] Parameterized generation for fuzzing
  - [ ] Incremental dataset updates
  - [ ] Dataset versioning and migration

- [ ] **Test Organization**
  - [ ] Move existing demo-files into proper structure
  - [ ] Integrate with existing integration tests
  - [ ] Documentation for dataset usage patterns
  - [ ] Developer guides for adding new scenarios

## Success Criteria

### Performance Validation (7EP-0004)
- [ ] Can generate 10K+ archive datasets in <30 seconds
- [ ] Benchmark suite validates <50ms resolution requirement
- [ ] List filtering benchmarks validate <200ms requirement
- [ ] Show command benchmarks validate <100ms requirement
- [ ] Memory usage testing validates <10MB requirement

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

## Related Work

### Builds On
- **7EP-0004**: Provides performance testing needed for completion
- **Test Helpers (PR #6)**: Uses existing test infrastructure
- **Integration Tests**: Enhances scenarios from mas-foundation-integration.md

### Enables
- **7EP-0004 Completion**: Provides missing performance testing capability
- **Future 7EPs**: Solid test foundation for trash management, CI, etc.
- **Regression Testing**: Systematic performance monitoring
- **Realistic Testing**: Better coverage of real-world scenarios

## Future Considerations

- **Fuzzing Integration**: Use datasets as basis for property-based testing
- **Cloud Storage Testing**: Remote file scenarios with simulated latency
- **Concurrency Testing**: Multi-user scenarios with shared datasets
- **Migration Testing**: Dataset evolution and backward compatibility
- **Performance Profiling**: Detailed bottleneck analysis with realistic data

---

This comprehensive test dataset system transforms 7zarch-go testing from ad hoc file creation to systematic, scalable, realistic test scenarios that enable confident performance validation and comprehensive integration testing.