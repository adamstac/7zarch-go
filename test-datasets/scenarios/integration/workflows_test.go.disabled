package integration

import (
	"testing"
	"time"

	"github.com/adamstac/7zarch-go/internal/storage"
	"github.com/adamstac/7zarch-go/test-datasets/helpers"
)

// TestCreateListShowDeleteWorkflow tests the complete archive lifecycle
func TestCreateListShowDeleteWorkflow(t *testing.T) {
	reg, archives := helpers.TestRegistryWithScenario(t, "create-list-show-delete")

	// Step 1: Verify all archives were created
	helpers.AssertArchiveCount(t, reg, len(archives))

	// Step 2: Test listing
	listed, err := reg.List()
	if err != nil {
		t.Fatalf("Failed to list archives: %v", err)
	}

	if len(listed) != len(archives) {
		t.Errorf("List count mismatch: got %d, expected %d", len(listed), len(archives))
	}

	// Step 3: Test showing individual archives
	resolver := storage.NewResolver(reg)
	for i, archive := range archives[:5] { // Test first 5
		// Test resolution by ULID prefix
		resolved, err := resolver.Resolve(archive.UID[:8])
		if err != nil {
			t.Errorf("Failed to resolve archive %d: %v", i, err)
			continue
		}

		if resolved.UID != archive.UID {
			t.Errorf("Resolution mismatch for archive %d: got %q, want %q",
				i, resolved.UID, archive.UID)
		}

		// Test getting by full UID
		fetched, err := reg.GetByUID(archive.UID)
		if err != nil {
			t.Errorf("Failed to get archive %d by UID: %v", i, err)
			continue
		}

		if fetched.Name != archive.Name {
			t.Errorf("Name mismatch for archive %d: got %q, want %q",
				i, fetched.Name, archive.Name)
		}
	}

	// Step 4: Test soft delete
	toDelete := archives[0]
	now := time.Now()
	toDelete.Status = "deleted"
	toDelete.DeletedAt = &now
	toDelete.OriginalPath = toDelete.Path
	toDelete.Path = "/trash/" + toDelete.Name

	if err := reg.Update(toDelete); err != nil {
		t.Fatalf("Failed to soft delete archive: %v", err)
	}

	// Verify deletion
	deleted, err := reg.GetByUID(toDelete.UID)
	if err != nil {
		t.Fatalf("Failed to get deleted archive: %v", err)
	}

	if deleted.Status != "deleted" {
		t.Errorf("Archive status not updated: got %q, want %q", deleted.Status, "deleted")
	}

	if deleted.DeletedAt == nil {
		t.Error("DeletedAt timestamp not set")
	}

	// Step 5: Test filtering to exclude deleted
	activeFilter := storage.ListFilters{
		Status: "present",
	}

	activeArchives, err := reg.ListWithFilters(activeFilter)
	if err != nil {
		t.Fatalf("Failed to filter active archives: %v", err)
	}

	if len(activeArchives) != len(archives)-1 {
		t.Errorf("Active archive count incorrect: got %d, expected %d",
			len(activeArchives), len(archives)-1)
	}

	// Verify deleted archive is not in active list
	for _, active := range activeArchives {
		if active.UID == toDelete.UID {
			t.Error("Deleted archive found in active list")
		}
	}
}

// TestMixedStorageWorkflow tests workflows with both managed and external archives
func TestMixedStorageWorkflow(t *testing.T) {
	reg, archives := helpers.TestRegistryWithScenario(t, "mixed-storage-scenario")

	// Separate managed and external archives
	var managed, external []*storage.Archive
	for _, archive := range archives {
		if archive.Managed {
			managed = append(managed, archive)
		} else {
			external = append(external, archive)
		}
	}

	t.Logf("Test scenario: %d managed, %d external archives", len(managed), len(external))

	// Test operations on managed archives
	if len(managed) > 0 {
		testArchive := managed[0]
		
		// Simulate upload
		uploadTime := time.Now()
		testArchive.Uploaded = true
		testArchive.UploadedAt = &uploadTime
		
		if err := reg.Update(testArchive); err != nil {
			t.Fatalf("Failed to update managed archive: %v", err)
		}

		// Verify upload status
		updated, err := reg.GetByUID(testArchive.UID)
		if err != nil {
			t.Fatalf("Failed to get updated archive: %v", err)
		}

		if !updated.Uploaded {
			t.Error("Managed archive upload status not updated")
		}

		if updated.UploadedAt == nil {
			t.Error("Upload timestamp not set")
		}
	}

	// Test filtering by managed status
	managedFilter := storage.ListFilters{
		Managed: &[]bool{true}[0],
	}

	managedList, err := reg.ListWithFilters(managedFilter)
	if err != nil {
		t.Fatalf("Failed to filter managed archives: %v", err)
	}

	if len(managedList) != len(managed) {
		t.Errorf("Managed filter count mismatch: got %d, expected %d",
			len(managedList), len(managed))
	}

	// Test profile distribution across managed/external
	profiles := make(map[string]int)
	for _, archive := range archives {
		profiles[archive.Profile]++
	}

	t.Logf("Profile distribution: %+v", profiles)

	// Verify cross-profile filtering works
	for profile := range profiles {
		filter := storage.ListFilters{
			Profile: profile,
		}

		profileArchives, err := reg.ListWithFilters(filter)
		if err != nil {
			t.Errorf("Failed to filter by profile %q: %v", profile, err)
			continue
		}

		if len(profileArchives) != profiles[profile] {
			t.Errorf("Profile %q count mismatch: got %d, expected %d",
				profile, len(profileArchives), profiles[profile])
		}
	}
}

// TestTimeSeriesWorkflow tests operations on time-distributed archives
func TestTimeSeriesWorkflow(t *testing.T) {
	reg, archives := helpers.TestRegistryWithScenario(t, "time-series-archives")

	// Test time-based queries
	now := time.Now()
	periods := []struct {
		name     string
		duration time.Duration
	}{
		{"last_week", 7 * 24 * time.Hour},
		{"last_month", 30 * 24 * time.Hour},
		{"last_quarter", 90 * 24 * time.Hour},
		{"last_year", 365 * 24 * time.Hour},
	}

	for _, period := range periods {
		t.Run(period.name, func(t *testing.T) {
			cutoff := now.Add(-period.duration)
			
			filter := storage.ListFilters{
				NewerThan: cutoff,
			}

			recent, err := reg.ListWithFilters(filter)
			if err != nil {
				t.Fatalf("Failed to filter archives newer than %v: %v", cutoff, err)
			}

			t.Logf("%s: found %d archives", period.name, len(recent))

			// Verify all returned archives are actually newer than cutoff
			for _, archive := range recent {
				if archive.Created.Before(cutoff) {
					t.Errorf("Archive %q created at %v is before cutoff %v",
						archive.Name, archive.Created, cutoff)
				}
			}
		})
	}

	// Test combining time and profile filters
	combinedFilter := storage.ListFilters{
		NewerThan: now.Add(-30 * 24 * time.Hour),
		Profile:   "documents",
	}

	combinedResults, err := reg.ListWithFilters(combinedFilter)
	if err != nil {
		t.Fatalf("Failed to apply combined filters: %v", err)
	}

	// Verify all results match both criteria
	for _, archive := range combinedResults {
		if archive.Profile != "documents" {
			t.Errorf("Archive %q has wrong profile: %q", archive.Name, archive.Profile)
		}

		if archive.Created.Before(now.Add(-30 * 24 * time.Hour)) {
			t.Errorf("Archive %q is too old: %v", archive.Name, archive.Created)
		}
	}

	t.Logf("Combined filter returned %d documents from last 30 days", len(combinedResults))
}