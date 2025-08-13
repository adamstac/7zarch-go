package performance

import (
	"fmt"
	"testing"
	"time"

	"github.com/adamstac/7zarch-go/internal/storage"
	"github.com/adamstac/7zarch-go/test-datasets/helpers"
)

// BenchmarkULIDResolution tests ULID resolution performance with various scenarios
func BenchmarkULIDResolution(b *testing.B) {
	scenarios := []string{
		"small-test",              // 10 archives
		"medium-test",             // 100 archives
		"large-test",              // 1,000 archives
		"disambiguation-stress",   // 1,000 with similar ULIDs
		"scaling-validation",      // 10,000 archives
	}

	for _, scenario := range scenarios {
		b.Run(scenario, func(b *testing.B) {
			helpers.BenchmarkWithScenario(b, scenario, benchmarkResolution)
		})
	}
}

func benchmarkResolution(b *testing.B, reg *storage.Registry, archives []*storage.Archive) {
	resolver := storage.NewResolver(reg)
	resolver.MinPrefixLength = 4 // From 7EP-0006 learning

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		archive := archives[i%len(archives)]
		prefix := archive.UID[:8] // Use 8-character prefix
		
		_, err := resolver.Resolve(prefix)
		if err != nil && !isAmbiguousError(err) {
			b.Fatalf("Resolution failed: %v", err)
		}
	}
}

// BenchmarkListFiltering tests list filtering performance
func BenchmarkListFiltering(b *testing.B) {
	b.Run("scaling-validation", func(b *testing.B) {
		helpers.BenchmarkWithScenario(b, "scaling-validation", benchmarkFiltering)
	})
}

func benchmarkFiltering(b *testing.B, reg *storage.Registry, archives []*storage.Archive) {
	testCases := []struct {
		name   string
		filter storage.ListFilters
	}{
		{
			name:   "status_present",
			filter: storage.ListFilters{Status: "present"},
		},
		{
			name:   "profile_media",
			filter: storage.ListFilters{Profile: "media"},
		},
		{
			name:   "larger_than_100mb",
			filter: storage.ListFilters{LargerThan: 100 * 1024 * 1024},
		},
		{
			name:   "older_than_30d",
			filter: storage.ListFilters{OlderThan: 30 * 24 * time.Hour},
		},
		{
			name: "complex_combination",
			filter: storage.ListFilters{
				Profile:    "media",
				LargerThan: 10 * 1024 * 1024,
				Status:     "present",
			},
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
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

// BenchmarkRegistryOperations tests various registry operations
func BenchmarkRegistryOperations(b *testing.B) {
	scenarios := []struct {
		name  string
		count int
	}{
		{"small", 10},
		{"medium", 100},
		{"large", 1000},
	}

	for _, scenario := range scenarios {
		b.Run(scenario.name+"_add", func(b *testing.B) {
			benchmarkAdd(b, scenario.count)
		})

		b.Run(scenario.name+"_list", func(b *testing.B) {
			benchmarkList(b, scenario.count)
		})

		b.Run(scenario.name+"_get", func(b *testing.B) {
			benchmarkGet(b, scenario.count)
		})
	}
}

func benchmarkAdd(b *testing.B, count int) {
	b.StopTimer()
	reg, archives := helpers.CreateTestRegistry(b, count)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		archive := archives[i%len(archives)]
		// Modify UID to make it unique for each add
		archive.UID = fmt.Sprintf("%s-%d", archive.UID, i)
		
		if err := reg.Add(archive); err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkList(b *testing.B, count int) {
	reg, _ := helpers.CreateTestRegistry(b, count)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		archives, err := reg.List()
		if err != nil {
			b.Fatal(err)
		}
		_ = archives
	}
}

func benchmarkGet(b *testing.B, count int) {
	reg, archives := helpers.CreateTestRegistry(b, count)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		archive := archives[i%len(archives)]
		
		result, err := reg.GetByUID(archive.UID)
		if err != nil {
			b.Fatal(err)
		}
		_ = result
	}
}

// Helper function to check if error is ambiguous
func isAmbiguousError(err error) bool {
	if err == nil {
		return false
	}
	// Check if error message contains "ambiguous"
	return len(err.Error()) > 9 && err.Error()[:9] == "ambiguous"
}