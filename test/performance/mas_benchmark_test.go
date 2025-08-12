package performance

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/adamstac/7zarch-go/internal/storage"
)

// generateTestArchives creates count test archives with realistic metadata
func generateTestArchives(count int) []*storage.Archive {
	archives := make([]*storage.Archive, count)
	rng := rand.New(rand.NewSource(42)) // Reproducible

	profiles := []string{"documents", "media", "balanced"}
	sizes := []int64{1024, 100 * 1024, 10 * 1024 * 1024} // 1KB, 100KB, 10MB

	for i := 0; i < count; i++ {
		archive := &storage.Archive{
			UID:      generateSequentialUID(i), // For prefix testing
			Name:     fmt.Sprintf("archive-%04d.7z", i),
			Path:     fmt.Sprintf("/tmp/test-archive-%04d.7z", i),
			Size:     sizes[rng.Intn(len(sizes))],
			Created:  time.Now().Add(-time.Duration(rng.Intn(365)) * 24 * time.Hour),
			Profile:  profiles[rng.Intn(len(profiles))],
			Managed:  i%10 != 0, // 90% managed, 10% external
			Status:   "present",
			Checksum: fmt.Sprintf("sha256:%032x", i), // Fake but unique
		}
		archives[i] = archive
	}

	return archives
}

// generateSequentialUID creates realistic ULIDs with some overlap for disambiguation testing
func generateSequentialUID(i int) string {
	// Create proper 26-character ULIDs that are mostly unique but with some similar prefixes
	// This tests disambiguation performance while maintaining ULID format
	if i < 100 {
		// Similar prefixes for first 100: 01K2E00, 01K2E01, 01K2E02, etc.
		return fmt.Sprintf("01K2E%02d%012d%08d", i/10, i, i*17) // Group by tens
	}
	// Different patterns for others to avoid too much similarity
	return fmt.Sprintf("01K2F%02d%012d%08d", (i-100)/100, i, i*23)
}

// setupRegistryWithArchives creates a test registry and populates it
func setupRegistryWithArchives(tb testing.TB, archives []*storage.Archive) *storage.Registry {
	tb.Helper()

	// Create temp directory for test database
	tmpDir := tb.TempDir()
	dbPath := fmt.Sprintf("%s/test.db", tmpDir)
	
	reg, err := storage.NewRegistry(dbPath)
	if err != nil {
		tb.Fatalf("Failed to create test registry: %v", err)
	}
	
	tb.Cleanup(func() {
		reg.Close()
	})

	// Add all archives
	for _, archive := range archives {
		if err := reg.Add(archive); err != nil {
			tb.Fatalf("Failed to add test archive: %v", err)
		}
	}

	return reg
}

// isAmbiguousError checks if error is disambiguation error
func isAmbiguousError(err error) bool {
	_, ok := err.(*storage.AmbiguousIDError)
	return ok
}

// simulateChecksumVerification simulates the time cost of file verification
func simulateChecksumVerification(archive *storage.Archive) error {
	// Simulate the time cost of reading and hashing a file
	time.Sleep(10 * time.Microsecond) // Simulate small file verification
	return nil
}

// formatArchiveDisplay simulates show command formatting work
func formatArchiveDisplay(archive *storage.Archive) string {
	return fmt.Sprintf("Archive: %s (%s)\nSize: %d bytes\n",
		archive.Name, archive.UID[:8], archive.Size)
}

// BenchmarkULIDResolution tests resolution performance with 1000 archives
// Target: <50 microseconds per operation
func BenchmarkULIDResolution(b *testing.B) {
	archives := generateTestArchives(1000)
	reg := setupRegistryWithArchives(b, archives)
	resolver := storage.NewResolver(reg)
	resolver.MinPrefixLength = 4 // Allow shorter prefixes for testing

	// Test cases: full UID, 8-char prefix, 4-char prefix (disambiguation)
	testCases := []struct {
		name  string
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
				_, err := resolver.Resolve(id)
				// Allow AmbiguousIDError for 4-char prefixes
				if err != nil && !isAmbiguousError(err) {
					b.Fatalf("Resolution failed: %v", err)
				}
			}
		})
	}
}

// BenchmarkListFiltering tests list performance with 10K archives
// Target: <200 milliseconds per operation for 10K archives
func BenchmarkListFiltering(b *testing.B) {
	archives := generateTestArchives(10000) // Larger set for filtering
	reg := setupRegistryWithArchives(b, archives)

	// Define filter types based on the list command implementation
	testCases := []struct {
		name   string
		filter func() ([]*storage.Archive, error)
	}{
		{"no_filter", func() ([]*storage.Archive, error) {
			return reg.List()
		}},
		{"status_present", func() ([]*storage.Archive, error) {
			all, err := reg.List()
			if err != nil {
				return nil, err
			}
			var filtered []*storage.Archive
			for _, a := range all {
				if a.Status == "present" {
					filtered = append(filtered, a)
				}
			}
			return filtered, nil
		}},
		{"profile_media", func() ([]*storage.Archive, error) {
			all, err := reg.List()
			if err != nil {
				return nil, err
			}
			var filtered []*storage.Archive
			for _, a := range all {
				if a.Profile == "media" {
					filtered = append(filtered, a)
				}
			}
			return filtered, nil
		}},
		{"managed_only", func() ([]*storage.Archive, error) {
			all, err := reg.List()
			if err != nil {
				return nil, err
			}
			var filtered []*storage.Archive
			for _, a := range all {
				if a.Managed {
					filtered = append(filtered, a)
				}
			}
			return filtered, nil
		}},
		{"large_files", func() ([]*storage.Archive, error) {
			all, err := reg.List()
			if err != nil {
				return nil, err
			}
			var filtered []*storage.Archive
			for _, a := range all {
				if a.Size > 1024*1024 { // >1MB
					filtered = append(filtered, a)
				}
			}
			return filtered, nil
		}},
		{"old_files", func() ([]*storage.Archive, error) {
			all, err := reg.List()
			if err != nil {
				return nil, err
			}
			cutoff := time.Now().Add(-30 * 24 * time.Hour) // 30 days ago
			var filtered []*storage.Archive
			for _, a := range all {
				if a.Created.Before(cutoff) {
					filtered = append(filtered, a)
				}
			}
			return filtered, nil
		}},
		{"complex", func() ([]*storage.Archive, error) {
			all, err := reg.List()
			if err != nil {
				return nil, err
			}
			var filtered []*storage.Archive
			for _, a := range all {
				if a.Profile == "media" && a.Size > 100*1024 && a.Status == "present" {
					filtered = append(filtered, a)
				}
			}
			return filtered, nil
		}},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				results, err := tc.filter()
				if err != nil {
					b.Fatalf("List filtering failed: %v", err)
				}
				_ = results // Use results to prevent optimization
			}
		})
	}
}

// BenchmarkShowCommand tests show command performance
// Target: <100 milliseconds per operation including verification
func BenchmarkShowCommand(b *testing.B) {
	archives := generateTestArchives(1000)
	reg := setupRegistryWithArchives(b, archives)
	resolver := storage.NewResolver(reg)
	resolver.MinPrefixLength = 4 // Allow shorter prefixes for testing

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

				// Resolve archive (use full UID to avoid ambiguity in show benchmark)
				resolved, err := resolver.Resolve(archive.UID)
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

// BenchmarkScaling tests performance at different dataset sizes
func BenchmarkScaling(b *testing.B) {
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("resolution_%d_archives", size), func(b *testing.B) {
			archives := generateTestArchives(size)
			reg := setupRegistryWithArchives(b, archives)
			resolver := storage.NewResolver(reg)
			resolver.MinPrefixLength = 4 // Allow shorter prefixes for testing

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				archive := archives[i%len(archives)]
				_, err := resolver.Resolve(archive.UID) // Use full UID for scaling test
				if err != nil {
					b.Fatalf("Resolution failed: %v", err)
				}
			}
		})
	}
}