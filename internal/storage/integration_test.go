package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestMASWorkflowIntegration tests end-to-end MAS workflows
func TestMASWorkflowIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	t.Run("create_list_show_workflow", func(t *testing.T) {
		// This test documents the expected workflow AC will implement
		manager, tempDir := setupTestManager(t)
		defer os.RemoveAll(tempDir)

		// Step 1: Create archives (simulates `7zarch-go create`)
		archives := []struct {
			name    string
			path    string
			size    int64
			profile string
			managed bool
		}{
			{"project-backup.7z", "/managed/project-backup.7z", 2097152, "media", true},
			{"docs.7z", "/external/docs.7z", 524288, "documents", false},
			{"podcast-103.7z", "/managed/podcast-103.7z", 163577856, "media", true},
		}

		var createdUIDs []string
		for _, arch := range archives {
			err := manager.Add(arch.name, arch.path, arch.size, arch.profile,
				"", "", arch.managed)
			if err != nil {
				t.Fatalf("Failed to create archive %s: %v", arch.name, err)
			}

			// Get the UID for later reference
			created, err := manager.Get(arch.name)
			if err != nil {
				t.Fatalf("Failed to retrieve created archive %s: %v", arch.name, err)
			}
			createdUIDs = append(createdUIDs, created.UID)
		}

		// Step 2: List archives (simulates `7zarch-go list`)
		allArchives, err := manager.List()
		if err != nil {
			t.Fatalf("Failed to list archives: %v", err)
		}

		if len(allArchives) != len(archives) {
			t.Errorf("Expected %d archives, got %d", len(archives), len(allArchives))
		}

		// Verify managed vs external separation
		managedCount := 0
		externalCount := 0
		for _, arch := range allArchives {
			if arch.Managed {
				managedCount++
			} else {
				externalCount++
			}
		}

		if managedCount != 2 || externalCount != 1 {
			t.Errorf("Expected 2 managed, 1 external; got %d managed, %d external",
				managedCount, externalCount)
		}

		// Step 3: Show individual archives (simulates `7zarch-go show <uid>`)
		for i, uid := range createdUIDs {
			// Test resolution by UID (AC's resolver will handle this)
			archive, err := manager.Get(archives[i].name)
			if err != nil {
				t.Fatalf("Failed to show archive %s: %v", uid, err)
			}

			if archive.UID != uid {
				t.Errorf("Expected UID %s, got %s", uid, archive.UID)
			}

			if archive.Name != archives[i].name {
				t.Errorf("Expected name %s, got %s", archives[i].name, archive.Name)
			}
		}

		t.Logf("Create→List→Show workflow completed successfully")
	})

	t.Run("mixed_storage_workflow", func(t *testing.T) {
		manager, tempDir := setupTestManager(t)
		defer os.RemoveAll(tempDir)

		// Create archives in different locations
		testCases := []struct {
			name     string
			path     string
			managed  bool
			location string
		}{
			{"managed-1.7z", "/managed/managed-1.7z", true, "managed"},
			{"managed-2.7z", "/managed/managed-2.7z", true, "managed"},
			{"external-home.7z", "/home/user/external-home.7z", false, "external"},
			{"external-backup.7z", "/backup/external-backup.7z", false, "external"},
			{"external-nas.7z", "/nas/external-nas.7z", false, "external"},
		}

		for _, tc := range testCases {
			err := manager.Add(tc.name, tc.path, 1024, "balanced", "", "", tc.managed)
			if err != nil {
				t.Fatalf("Failed to add %s archive %s: %v", tc.location, tc.name, err)
			}
		}

		// List all archives
		all, err := manager.List()
		if err != nil {
			t.Fatalf("Failed to list mixed archives: %v", err)
		}

		// Verify location tracking
		managedArchives := make([]*Archive, 0)
		externalArchives := make([]*Archive, 0)

		for _, arch := range all {
			if arch.Managed {
				managedArchives = append(managedArchives, arch)
			} else {
				externalArchives = append(externalArchives, arch)
			}
		}

		if len(managedArchives) != 2 {
			t.Errorf("Expected 2 managed archives, got %d", len(managedArchives))
		}

		if len(externalArchives) != 3 {
			t.Errorf("Expected 3 external archives, got %d", len(externalArchives))
		}

		// Verify path tracking
		for _, arch := range managedArchives {
			if !filepath.HasPrefix(arch.Path, "/managed/") {
				t.Errorf("Managed archive has wrong path: %s", arch.Path)
			}
		}

		for _, arch := range externalArchives {
			if filepath.HasPrefix(arch.Path, "/managed/") {
				t.Errorf("External archive has managed path: %s", arch.Path)
			}
		}

		t.Logf("Mixed storage workflow: %d managed, %d external archives tracked",
			len(managedArchives), len(externalArchives))
	})

	t.Run("upload_workflow", func(t *testing.T) {
		manager, tempDir := setupTestManager(t)
		defer os.RemoveAll(tempDir)

		// Create archive
		err := manager.Add("upload-test.7z", "/test/upload-test.7z", 1024, "balanced",
			"checksum123", "", true)
		if err != nil {
			t.Fatalf("Failed to create test archive: %v", err)
		}

		// List not uploaded (simulates `7zarch-go list --not-uploaded`)
		notUploaded, err := manager.ListNotUploaded()
		if err != nil {
			t.Fatalf("Failed to list not uploaded: %v", err)
		}

		if len(notUploaded) != 1 {
			t.Errorf("Expected 1 not uploaded archive, got %d", len(notUploaded))
		}

		if notUploaded[0].Name != "upload-test.7z" {
			t.Errorf("Expected 'upload-test.7z', got '%s'", notUploaded[0].Name)
		}

		// Mark as uploaded (simulates upload process)
		err = manager.MarkUploaded("upload-test.7z", "s3://bucket/upload-test.7z")
		if err != nil {
			t.Fatalf("Failed to mark as uploaded: %v", err)
		}

		// Verify upload status
		uploaded, err := manager.Get("upload-test.7z")
		if err != nil {
			t.Fatalf("Failed to get uploaded archive: %v", err)
		}

		if !uploaded.Uploaded {
			t.Error("Expected archive to be marked as uploaded")
		}

		if uploaded.Destination != "s3://bucket/upload-test.7z" {
			t.Errorf("Expected destination 's3://bucket/upload-test.7z', got '%s'",
				uploaded.Destination)
		}

		if uploaded.UploadedAt == nil {
			t.Error("Expected upload timestamp to be set")
		}

		// List not uploaded again - should be empty
		notUploadedAfter, err := manager.ListNotUploaded()
		if err != nil {
			t.Fatalf("Failed to list not uploaded after upload: %v", err)
		}

		if len(notUploadedAfter) != 0 {
			t.Errorf("Expected 0 not uploaded archives, got %d", len(notUploadedAfter))
		}

		t.Logf("Upload workflow completed: marked as uploaded to %s", uploaded.Destination)
	})

	t.Run("age_based_filtering_workflow", func(t *testing.T) {
		manager, tempDir := setupTestManager(t)
		defer os.RemoveAll(tempDir)

		now := time.Now()

		// Create archives with different ages
		archiveAges := []struct {
			name string
			age  time.Duration
		}{
			{"recent.7z", -1 * time.Hour},          // 1 hour ago
			{"day-old.7z", -24 * time.Hour},        // 1 day ago
			{"week-old.7z", -7 * 24 * time.Hour},   // 1 week ago
			{"month-old.7z", -30 * 24 * time.Hour}, // 1 month ago
		}

		for _, arch := range archiveAges {
			// Create archive with specific timestamp
			testArchive := &Archive{
				UID:     generateUID(),
				Name:    arch.name,
				Path:    "/test/" + arch.name,
				Size:    1024,
				Created: now.Add(arch.age),
				Profile: "balanced",
				Managed: true,
				Status:  "present",
			}

			err := manager.Registry().Add(testArchive)
			if err != nil {
				t.Fatalf("Failed to add archive %s: %v", arch.name, err)
			}
		}

		// Test age-based filtering (simulates `7zarch-go list --older-than 3d`)
		threeDaysAgo := 3 * 24 * time.Hour
		oldArchives, err := manager.ListOlderThan(threeDaysAgo)
		if err != nil {
			t.Fatalf("Failed to list older archives: %v", err)
		}

		expectedOld := []string{"week-old.7z", "month-old.7z"}
		if len(oldArchives) != len(expectedOld) {
			t.Errorf("Expected %d old archives, got %d", len(expectedOld), len(oldArchives))
		}

		// Verify correct archives were returned
		oldNames := make(map[string]bool)
		for _, arch := range oldArchives {
			oldNames[arch.Name] = true
		}

		for _, expectedName := range expectedOld {
			if !oldNames[expectedName] {
				t.Errorf("Expected old archive '%s' not found", expectedName)
			}
		}

		t.Logf("Age-based filtering: found %d archives older than %v",
			len(oldArchives), threeDaysAgo)
	})
}

// TestResolverIntegration tests resolver with real registry data
func TestResolverIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	t.Run("resolver_with_large_dataset", func(t *testing.T) {
		registry, _ := setupTestRegistry(t)

		// Create a realistic dataset with proper ULID format
		// ULIDs are 26 chars: timestamp(10) + randomness(16)
		prefixes := []string{"01JEX4RT2N", "01JEY5SU3O", "01JEZ6TV4P", "01JFA7UW5Q", "01JFB8VX6R"}
		var allArchives []*Archive

		for i, prefix := range prefixes {
			for j := 0; j < 10; j++ {
				// Create full 26-char ULID with proper format
				uid := fmt.Sprintf("%s%04dABCDEFGHIJKLMN", prefix, j)[:26]
				archive := &Archive{
					UID:      uid,
					Name:     fmt.Sprintf("archive-%s-%02d.7z", prefix[:5], j),
					Path:     fmt.Sprintf("/test/archive-%s-%02d.7z", prefix[:5], j),
					Size:     int64(1024 * (j + 1)),
					Created:  time.Now(),
					Checksum: fmt.Sprintf("checksum-%s-%02d", prefix[:5], j),
					Profile:  "balanced",
					Managed:  i%2 == 0,
					Status:   "present",
				}

				err := registry.Add(archive)
				if err != nil {
					t.Fatalf("Failed to add test archive: %v", err)
				}
				allArchives = append(allArchives, archive)
			}
		}

		resolver := NewResolver(registry)
		resolver.MinPrefixLength = 4 // Lower for testing

		// Test unique prefix resolution with exact UID
		testCases := []struct {
			input    string
			expected string
		}{
			{allArchives[0].UID, "archive-01JEX-00.7z"},  // Full UID
			{allArchives[10].UID, "archive-01JEY-00.7z"}, // Full UID
			{allArchives[20].UID, "archive-01JEZ-00.7z"}, // Full UID
		}

		for _, tc := range testCases {
			result, err := resolver.Resolve(tc.input)
			if err != nil {
				t.Errorf("Failed to resolve %s: %v", tc.input, err)
				continue
			}

			if result.Name != tc.expected {
				t.Errorf("Resolved %s to %s, expected %s", tc.input, result.Name, tc.expected)
			}
		}

		// Test ambiguous prefix resolution
		ambiguousTests := []string{
			"01JEX", // Should match multiple archives
			"01JE",  // Should match even more
		}

		for _, input := range ambiguousTests {
			_, err := resolver.Resolve(input)
			if err == nil {
				t.Errorf("Expected ambiguous error for %s, got nil", input)
				continue
			}

			ambiguousErr, ok := err.(*AmbiguousIDError)
			if !ok {
				t.Errorf("Expected AmbiguousIDError for %s, got %T", input, err)
				continue
			}

			t.Logf("Ambiguous resolution for %s: %d matches", input, len(ambiguousErr.Matches))
		}

		t.Logf("Resolver integration test completed with %d archives", len(allArchives))
	})
}

// setupTestManager creates a test manager with registry
func setupTestManager(t *testing.T) (*Manager, string) {
	tempDir, err := os.MkdirTemp("", "7zarch-manager-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	manager, err := NewManager(tempDir)
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to create test manager: %v", err)
	}

	t.Cleanup(func() {
		manager.Close()
		os.RemoveAll(tempDir)
	})

	return manager, tempDir
}

// TestCommandIntegration documents expected command interactions
func TestCommandIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping command integration tests in short mode")
	}

	t.Run("create_command_integration", func(t *testing.T) {
		// This documents how AC's create command should work with MAS
		manager, tempDir := setupTestManager(t)
		defer os.RemoveAll(tempDir)

		// Simulate create command behavior
		scenarios := []struct {
			name          string
			outputFlag    string
			expectManaged bool
		}{
			{"default-managed.7z", "", true},                                 // No --output = managed
			{"custom-external.7z", "/custom/path/custom-external.7z", false}, // --output = external
		}

		for _, scenario := range scenarios {
			// Determine output path (AC's create command logic)
			var finalPath string
			var managed bool

			if scenario.outputFlag == "" {
				// No --output flag = use managed storage
				finalPath = manager.GetManagedPath(scenario.name)
				managed = true
			} else {
				// --output specified = external storage
				finalPath = scenario.outputFlag
				managed = false
			}

			// Register the archive (what create command should do)
			err := manager.Add(scenario.name, finalPath, 1024, "balanced", "", "", managed)
			if err != nil {
				t.Fatalf("Failed to register archive %s: %v", scenario.name, err)
			}

			// Verify registration
			archive, err := manager.Get(scenario.name)
			if err != nil {
				t.Fatalf("Failed to retrieve registered archive %s: %v", scenario.name, err)
			}

			if archive.Managed != scenario.expectManaged {
				t.Errorf("Archive %s: expected managed=%v, got %v",
					scenario.name, scenario.expectManaged, archive.Managed)
			}

			if archive.Path != finalPath {
				t.Errorf("Archive %s: expected path %s, got %s",
					scenario.name, finalPath, archive.Path)
			}
		}

		t.Logf("Create command integration: both managed and external registration work")
	})

	t.Run("list_command_filtering", func(t *testing.T) {
		manager, tempDir := setupTestManager(t)
		defer os.RemoveAll(tempDir)

		// Create mixed dataset
		archives := []struct {
			name     string
			managed  bool
			uploaded bool
			age      time.Duration
		}{
			{"managed-uploaded.7z", true, true, -2 * 24 * time.Hour},
			{"managed-not-uploaded.7z", true, false, -1 * 24 * time.Hour},
			{"external-uploaded.7z", false, true, -5 * 24 * time.Hour},
			{"external-not-uploaded.7z", false, false, -3 * 24 * time.Hour},
		}

		for _, arch := range archives {
			// Create archive with specific timestamp
			testArchive := &Archive{
				UID:      generateUID(),
				Name:     arch.name,
				Path:     "/test/" + arch.name,
				Size:     1024,
				Created:  time.Now().Add(arch.age), // Set the age
				Profile:  "balanced",
				Managed:  arch.managed,
				Status:   "present",
				Checksum: "",
				Uploaded: false,
			}

			err := manager.Registry().Add(testArchive)
			if err != nil {
				t.Fatalf("Failed to add archive %s: %v", arch.name, err)
			}

			if arch.uploaded {
				err = manager.MarkUploaded(arch.name, "s3://test/"+arch.name)
				if err != nil {
					t.Fatalf("Failed to mark %s as uploaded: %v", arch.name, err)
				}
			}
		}

		// Test filtering scenarios (what AC's list command should support)

		// Filter: --not-uploaded
		notUploaded, err := manager.ListNotUploaded()
		if err != nil {
			t.Fatalf("Failed to list not uploaded: %v", err)
		}

		expectedNotUploaded := 2
		if len(notUploaded) != expectedNotUploaded {
			t.Errorf("Expected %d not uploaded, got %d", expectedNotUploaded, len(notUploaded))
		}

		// Filter: --older-than 3d
		threeDays := 3 * 24 * time.Hour
		olderArchives, err := manager.ListOlderThan(threeDays)
		if err != nil {
			t.Fatalf("Failed to list older archives: %v", err)
		}

		// Both external-uploaded.7z (5 days) and external-not-uploaded.7z (3 days) are >= 3 days old
		expectedOlder := 2
		if len(olderArchives) != expectedOlder {
			t.Errorf("Expected %d older archives, got %d", expectedOlder, len(olderArchives))
		}

		t.Logf("List command filtering: %d not uploaded, %d older than 3d",
			len(notUploaded), len(olderArchives))
	})
}
