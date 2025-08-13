package cmdutil

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/adamstac/7zarch-go/internal/storage"
)

func TestInitStorageManager(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	
	// Set environment variable for config path
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", oldHome)
	
	// Create minimal config file
	configContent := `storage:
  managed_path: "` + filepath.Join(tempDir, "archives") + `"
  use_managed_default: true
`
	configPath := filepath.Join(tempDir, ".7zarch-go-config")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	cfg, mgr, cleanup, err := InitStorageManager()
	if err != nil {
		t.Fatalf("InitStorageManager failed: %v", err)
	}
	defer cleanup()

	if cfg == nil {
		t.Error("Expected config to be non-nil")
	}
	if mgr == nil {
		t.Error("Expected manager to be non-nil")
	}
	if cleanup == nil {
		t.Error("Expected cleanup function to be non-nil")
	}

	// Test cleanup doesn't panic
	cleanup()
}

func TestLoadConfigOrDefault(t *testing.T) {
	// Test with invalid config path (should return defaults)
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", "/nonexistent")
	defer os.Setenv("HOME", oldHome)

	cfg := LoadConfigOrDefault()
	if cfg == nil {
		t.Error("Expected default config to be non-nil")
	}
	
	// Should have default values
	if cfg.Storage.UseManagedDefault != true {
		t.Error("Expected default config to have UseManagedDefault=true")
	}
}

func TestHandleResolverError(t *testing.T) {
	tests := []struct {
		name        string
		inputError  error
		id          string
		expectError string
	}{
		{
			name:        "AmbiguousIDError",
			inputError:  &storage.AmbiguousIDError{ID: "abc", Matches: nil},
			id:          "abc",
			expectError: "Invalid archive ID 'abc': matches multiple archives. Use a longer prefix or full UID",
		},
		{
			name:        "ArchiveNotFoundError", 
			inputError:  &storage.ArchiveNotFoundError{ID: "xyz"},
			id:          "xyz",
			expectError: "Archive 'xyz' not found. Try: use 'list' to see available archives, check the archive ID",
		},
		{
			name:        "UnknownError",
			inputError:  fmt.Errorf("unknown error"),
			id:          "test",
			expectError: "unknown error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := HandleResolverError(tt.inputError, tt.id)
			if err == nil {
				t.Error("Expected error to be non-nil")
				return
			}
			if err.Error() != tt.expectError {
				t.Errorf("Expected error message %q, got %q", tt.expectError, err.Error())
			}
		})
	}
}