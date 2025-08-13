// Package helpers provides test utilities that integrate the dataset system with storage
package helpers

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/adamstac/7zarch-go/internal/storage"
	"github.com/adamstac/7zarch-go/test-datasets/generators"
)

// TestRegistryWithScenario creates a registry populated with a predefined scenario
// Builds directly on 7EP-0006's successful setupRegistryWithArchives pattern
func TestRegistryWithScenario(tb testing.TB, scenarioName string) (*storage.Registry, []*storage.Archive) {
	tb.Helper()

	// Create temporary registry
	tmpDir := tb.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	reg, err := storage.NewRegistry(dbPath)
	if err != nil {
		tb.Fatalf("Failed to create test registry: %v", err)
	}

	tb.Cleanup(func() {
		if err := reg.Close(); err != nil {
			tb.Errorf("Failed to close registry: %v", err)
		}
	})

	// Get scenario specification
	spec, err := generators.GetScenario(scenarioName)
	if err != nil {
		tb.Fatalf("Failed to get scenario %q: %v", scenarioName, err)
	}

	// Generate archives using scenario system
	generator := generators.NewGenerator(42) // Fixed seed like 7EP-0006
	testArchives := generator.GenerateScenario(tb, spec)

	// Convert to storage archives and add to registry
	archives := make([]*storage.Archive, len(testArchives))
	for i, ta := range testArchives {
		archive := &storage.Archive{
			UID:          ta.UID,
			Name:         ta.Name,
			Path:         ta.Path,
			Size:         ta.Size,
			Created:      ta.Created,
			Profile:      ta.Profile,
			Managed:      ta.Managed,
			Status:       ta.Status,
			Checksum:     ta.Checksum,
			Uploaded:     ta.Uploaded,
			UploadedAt:   ta.UploadedAt,
			DeletedAt:    ta.DeletedAt,
			OriginalPath: ta.OriginalPath,
		}

		if err := reg.Add(archive); err != nil {
			tb.Fatalf("Failed to add test archive %d: %v", i, err)
		}
		archives[i] = archive
	}

	return reg, archives
}

// BenchmarkWithScenario runs a benchmark with a predefined scenario
func BenchmarkWithScenario(b *testing.B, scenarioName string,
	benchmarkFn func(*testing.B, *storage.Registry, []*storage.Archive)) {

	reg, archives := TestRegistryWithScenario(b, scenarioName)

	b.ResetTimer()
	benchmarkFn(b, reg, archives)
}

// CreateTestRegistry creates a simple test registry with the specified number of archives
// Simplified version for basic tests
func CreateTestRegistry(tb testing.TB, count int) (*storage.Registry, []*storage.Archive) {
	tb.Helper()

	spec := generators.ScenarioSpec{
		Name:        "test",
		Count:       count,
		ULIDPattern: generators.ULIDUnique,
		Profiles: []generators.ProfileDistribution{
			{Profile: "balanced", Weight: 1.0},
		},
		SizePattern: generators.SizeUniform,
	}

	tmpDir := tb.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	reg, err := storage.NewRegistry(dbPath)
	if err != nil {
		tb.Fatalf("Failed to create test registry: %v", err)
	}

	tb.Cleanup(func() {
		if err := reg.Close(); err != nil {
			tb.Errorf("Failed to close registry: %v", err)
		}
	})

	generator := generators.NewGenerator(42)
	testArchives := generator.GenerateScenario(tb, spec)

	archives := make([]*storage.Archive, len(testArchives))
	for i, ta := range testArchives {
		archive := &storage.Archive{
			UID:      ta.UID,
			Name:     ta.Name,
			Path:     ta.Path,
			Size:     ta.Size,
			Created:  ta.Created,
			Profile:  ta.Profile,
			Managed:  ta.Managed,
			Status:   ta.Status,
			Checksum: ta.Checksum,
		}

		if err := reg.Add(archive); err != nil {
			tb.Fatalf("Failed to add test archive: %v", err)
		}
		archives[i] = archive
	}

	return reg, archives
}

// AssertArchiveCount verifies the expected number of archives in the registry
func AssertArchiveCount(tb testing.TB, reg *storage.Registry, expected int) {
	tb.Helper()

	archives, err := reg.List()
	if err != nil {
		tb.Fatalf("Failed to list archives: %v", err)
	}

	if len(archives) != expected {
		tb.Errorf("Expected %d archives, got %d", expected, len(archives))
	}
}

// AssertResolves verifies that a ULID prefix resolves to the expected archive
func AssertResolves(tb testing.TB, resolver *storage.Resolver, prefix string, expectedUID string) {
	tb.Helper()

	resolved, err := resolver.Resolve(prefix)
	if err != nil {
		tb.Fatalf("Failed to resolve %q: %v", prefix, err)
	}

	if resolved.UID != expectedUID {
		tb.Errorf("Expected resolution of %q to %q, got %q", prefix, expectedUID, resolved.UID)
	}
}

// AssertAmbiguous verifies that a ULID prefix is ambiguous
func AssertAmbiguous(tb testing.TB, resolver *storage.Resolver, prefix string) {
	tb.Helper()

	_, err := resolver.Resolve(prefix)
	if err == nil {
		tb.Errorf("Expected ambiguous error for prefix %q, but resolution succeeded", prefix)
	}

	// Check if it's actually an ambiguous error
	if err != nil && err.Error() != fmt.Sprintf("ambiguous prefix %q", prefix) {
		tb.Errorf("Expected ambiguous error for prefix %q, got: %v", prefix, err)
	}
}
