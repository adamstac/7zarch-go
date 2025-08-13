# Test Dataset System

Comprehensive test dataset generation system for 7zarch-go, built on the proven metadata-only approach from 7EP-0006.

## Overview

This system provides:
- **Reproducible test scenarios** with fixed seeds
- **Performance benchmarks** for scaling validation
- **Edge case coverage** for Unicode, boundaries, and special conditions
- **Integration test support** for end-to-end workflows

## Quick Start

### Using in Tests

```go
import "github.com/adamstac/7zarch-go/test-datasets/helpers"

func TestMyFeature(t *testing.T) {
    // Create registry with predefined scenario
    reg, archives := helpers.TestRegistryWithScenario(t, "small-test")
    
    // Your test logic here
}
```

### Running Benchmarks

```bash
# Run all performance benchmarks
go test -bench=. ./test-datasets/scenarios/performance/

# Run specific benchmark
go test -bench=BenchmarkULIDResolution ./test-datasets/scenarios/performance/
```

### Running Edge Case Tests

```bash
# Test all edge cases
go test ./test-datasets/scenarios/edge_cases/

# Test specific edge case
go test -run TestUnicodeNames ./test-datasets/scenarios/edge_cases/
```

## Available Scenarios

### Performance Testing
- `disambiguation-stress` - 1,000 archives with similar ULIDs
- `scaling-validation` - 10,000 archives for scale testing
- `resolution-stress` - 5,000 archives with intentional collisions

### Integration Testing
- `create-list-show-delete` - 50 archives for workflow testing
- `mixed-storage-scenario` - 100 archives, mixed managed/external
- `time-series-archives` - 365 archives distributed over a year

### Edge Cases
- `unicode-names` - 25 archives with international characters and emojis
- `boundary-conditions` - 30 archives testing size/path limits

### Unit Testing
- `small-test` - 10 archives for quick tests
- `medium-test` - 100 archives for moderate tests
- `large-test` - 1,000 archives for comprehensive tests

## Architecture

```
test-datasets/
├── generators/           # Core generation engine
│   ├── generator.go     # Main generator logic
│   └── scenarios.go     # Predefined scenarios
├── helpers/             # Test integration utilities
│   └── test_helpers.go  # Registry helpers
├── scenarios/           # Test implementations
│   ├── performance/     # Benchmark tests
│   ├── integration/     # Workflow tests
│   └── edge_cases/      # Edge case tests
└── datasets/            # Generated data cache (optional)
```

## Key Features

### Metadata-Only Approach
Following 7EP-0006's success, all test data is metadata-only - no actual files are created. This provides:
- **Fast generation** - 10K archives in <1 second
- **Minimal storage** - No disk space for test files
- **Pure performance** - Focus on registry operations

### Reproducible Seeds
All generators use fixed seeds (default: 42) ensuring:
- **Consistent results** across test runs
- **Debuggable failures** - Same data every time
- **Reliable benchmarks** - No variance from data changes

### ULID Patterns
Three ULID generation patterns for different test needs:
- **Unique** - All different, for scaling tests
- **Similar** - Controlled similarity for disambiguation
- **Collisions** - Intentional overlaps for stress testing

## Creating Custom Scenarios

```go
spec := generators.ScenarioSpec{
    Name:        "my-scenario",
    Count:       500,
    ULIDPattern: generators.ULIDSimilar,
    Profiles: []generators.ProfileDistribution{
        {Profile: "documents", Weight: 0.7},
        {Profile: "media", Weight: 0.3},
    },
    SizePattern: generators.SizeRealistic,
    TimeSpread:  30 * 24 * time.Hour,
    EdgeCases: []generators.EdgeCase{
        generators.UnicodeFilenames,
        generators.MixedManagedExternal,
    },
}

generator := generators.NewGenerator(42)
archives := generator.GenerateScenario(t, spec)
```

## Performance Characteristics

Based on 7EP-0006 validation:
- **Generation**: <1ms per archive (metadata only)
- **Registry operations**: O(1) performance
- **Memory usage**: <10MB for 10K archives
- **Benchmark overhead**: Negligible

## Integration with CI/CD

The test dataset system integrates with the GitHub Actions workflows:

```yaml
- name: Run performance benchmarks
  run: go test -bench=. ./test-datasets/scenarios/performance/
  
- name: Run edge case tests
  run: go test ./test-datasets/scenarios/edge_cases/
```

## Migration from Old Test Data

This system replaces the ad hoc `demo-files/` and `friends-demo/` directories with structured, reproducible test scenarios. To migrate existing tests:

1. Identify test data needs
2. Choose appropriate scenario or create custom
3. Replace file operations with `TestRegistryWithScenario()`
4. Remove old test directories

## Related Documentation

- [7EP-0005: Comprehensive Test Dataset System](../docs/7eps/7ep-0005-test-dataset-system.md)
- [7EP-0006: Minimal Performance Testing](../docs/7eps/7ep-0006-minimal-performance-testing.md)
- [Testing Guide](../docs/testing-guide.md)