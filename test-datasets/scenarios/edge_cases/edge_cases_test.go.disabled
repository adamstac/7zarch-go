package edge_cases

import (
	"strings"
	"testing"
	"time"

	"github.com/adamstac/7zarch-go/internal/storage"
	"github.com/adamstac/7zarch-go/test-datasets/helpers"
)

// TestUnicodeNames verifies handling of international characters and emojis
func TestUnicodeNames(t *testing.T) {
	reg, archives := helpers.TestRegistryWithScenario(t, "unicode-names")

	// Verify all archives were created
	helpers.AssertArchiveCount(t, reg, len(archives))

	// Test resolution with Unicode names
	resolver := storage.NewResolver(reg)
	for _, archive := range archives[:5] { // Test first 5
		resolved, err := resolver.Resolve(archive.UID[:8])
		if err != nil {
			t.Errorf("Failed to resolve Unicode-named archive %q: %v", archive.Name, err)
			continue
		}

		if resolved.UID != archive.UID {
			t.Errorf("Resolution mismatch for %q: got %q, want %q", 
				archive.Name, resolved.UID, archive.UID)
		}
	}

	// Test listing with Unicode names
	listed, err := reg.List()
	if err != nil {
		t.Fatalf("Failed to list archives with Unicode names: %v", err)
	}

	// Verify Unicode characters are preserved
	unicodeFound := false
	emojiFound := false
	specialFound := false

	for _, archive := range listed {
		if strings.Contains(archive.Name, "测试") || strings.Contains(archive.Name, "файл") {
			unicodeFound = true
		}
		if strings.Contains(archive.Name, "🚀") || strings.Contains(archive.Name, "📝") {
			emojiFound = true
		}
		if strings.Contains(archive.Name, "spaces") || strings.Contains(archive.Name, "[brackets]") {
			specialFound = true
		}
	}

	if !unicodeFound {
		t.Error("No Unicode filenames found in listing")
	}
	if !emojiFound {
		t.Error("No emoji filenames found in listing")
	}
	if !specialFound {
		t.Error("No special character filenames found in listing")
	}
}

// TestBoundaryConditions verifies handling of extreme values
func TestBoundaryConditions(t *testing.T) {
	reg, archives := helpers.TestRegistryWithScenario(t, "boundary-conditions")

	// Verify all archives were created
	helpers.AssertArchiveCount(t, reg, len(archives))

	// Check for extreme file sizes
	var zeroSizeCount, largeSizeCount int
	for _, archive := range archives {
		if archive.Size == 0 {
			zeroSizeCount++
		}
		if archive.Size >= 1024*1024*1024 { // 1GB or larger
			largeSizeCount++
		}
	}

	if zeroSizeCount == 0 {
		t.Error("No zero-size files found in boundary conditions")
	}
	if largeSizeCount == 0 {
		t.Error("No large files (>=1GB) found in boundary conditions")
	}

	// Test filtering with extreme sizes
	filters := storage.ListFilters{
		LargerThan: 500 * 1024 * 1024, // 500MB
	}
	
	largeArchives, err := reg.ListWithFilters(filters)
	if err != nil {
		t.Fatalf("Failed to filter large archives: %v", err)
	}

	if len(largeArchives) == 0 {
		t.Error("No archives found with size > 500MB")
	}

	// Test resolution with collision patterns
	resolver := storage.NewResolver(reg)
	ambiguousCount := 0
	
	for _, archive := range archives[:10] {
		_, err := resolver.Resolve(archive.UID[:4]) // Very short prefix
		if err != nil && strings.Contains(err.Error(), "ambiguous") {
			ambiguousCount++
		}
	}

	if ambiguousCount == 0 {
		t.Error("Expected some ambiguous resolutions with short prefixes")
	}
}

// TestMixedStorage verifies handling of managed vs external archives
func TestMixedStorage(t *testing.T) {
	reg, archives := helpers.TestRegistryWithScenario(t, "mixed-storage-scenario")

	// Count managed vs external
	var managedCount, externalCount int
	for _, archive := range archives {
		if archive.Managed {
			managedCount++
		} else {
			externalCount++
		}
	}

	// Verify we have a mix (roughly 70/30 based on scenario)
	managedRatio := float64(managedCount) / float64(len(archives))
	if managedRatio < 0.6 || managedRatio > 0.8 {
		t.Errorf("Unexpected managed ratio: %.2f (expected ~0.7)", managedRatio)
	}

	// Test filtering by managed status
	managedFilter := storage.ListFilters{
		Managed: &[]bool{true}[0], // Pointer to true
	}
	
	managedArchives, err := reg.ListWithFilters(managedFilter)
	if err != nil {
		t.Fatalf("Failed to filter managed archives: %v", err)
	}

	if len(managedArchives) != managedCount {
		t.Errorf("Filtered managed count mismatch: got %d, expected %d",
			len(managedArchives), managedCount)
	}

	// Test filtering by external status
	externalFilter := storage.ListFilters{
		Managed: &[]bool{false}[0], // Pointer to false
	}
	
	externalArchives, err := reg.ListWithFilters(externalFilter)
	if err != nil {
		t.Fatalf("Failed to filter external archives: %v", err)
	}

	if len(externalArchives) != externalCount {
		t.Errorf("Filtered external count mismatch: got %d, expected %d",
			len(externalArchives), externalCount)
	}
}

// TestTimeSeriesArchives verifies handling of archives over time
func TestTimeSeriesArchives(t *testing.T) {
	reg, archives := helpers.TestRegistryWithScenario(t, "time-series-archives")

	// Verify we have one archive per day
	if len(archives) != 365 {
		t.Errorf("Expected 365 archives (one per day), got %d", len(archives))
	}

	// Verify time spread
	if len(archives) > 1 {
		firstTime := archives[0].Created
		lastTime := archives[len(archives)-1].Created
		
		duration := lastTime.Sub(firstTime)
		expectedDuration := 364 * 24 * time.Hour // 364 days between first and last
		
		// Allow 1 hour tolerance
		if duration < expectedDuration-time.Hour || duration > expectedDuration+time.Hour {
			t.Errorf("Unexpected time spread: %v (expected ~%v)", duration, expectedDuration)
		}
	}

	// Test filtering by time ranges
	now := time.Now()
	thirtyDaysAgo := now.AddDate(0, 0, -30)
	
	recentFilter := storage.ListFilters{
		NewerThan: thirtyDaysAgo,
	}
	
	recentArchives, err := reg.ListWithFilters(recentFilter)
	if err != nil {
		t.Fatalf("Failed to filter recent archives: %v", err)
	}

	// Should have roughly 30 archives from the last 30 days
	if len(recentArchives) < 25 || len(recentArchives) > 35 {
		t.Errorf("Unexpected number of recent archives: %d (expected ~30)", len(recentArchives))
	}
}