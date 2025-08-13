package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	
	if cfg == nil {
		t.Fatal("Expected DefaultConfig to return non-nil config")
	}
	
	// Test default values
	if !cfg.Storage.UseManagedDefault {
		t.Error("Expected UseManagedDefault to be true by default")
	}
	
	if cfg.Storage.RetentionDays != 30 {
		t.Errorf("Expected RetentionDays 30, got %d", cfg.Storage.RetentionDays)
	}
	
	if cfg.Compression.MediaThreshold != 70 {
		t.Errorf("Expected MediaThreshold 70, got %d", cfg.Compression.MediaThreshold)
	}
	
	if cfg.Compression.DocsThreshold != 60 {
		t.Errorf("Expected DocsThreshold 60, got %d", cfg.Compression.DocsThreshold)
	}
}

func TestLoadConfig_NotExists(t *testing.T) {
	// Set HOME to non-existent directory
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", "/nonexistent")
	defer os.Setenv("HOME", oldHome)
	
	cfg, err := Load()
	
	// Should return default config when file doesn't exist (not an error)
	if err != nil {
		t.Errorf("Expected no error when HOME doesn't exist, got: %v", err)
	}
	if cfg == nil {
		t.Error("Expected default config when HOME doesn't exist")
	}
	
	// Should have default values
	if cfg.Storage.RetentionDays != 30 {
		t.Errorf("Expected default RetentionDays 30, got %d", cfg.Storage.RetentionDays)
	}
}

func TestLoadConfig_ValidConfig(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	
	// Create valid config file
	configContent := `storage:
  managed_path: "/custom/path"
  use_managed_default: false
  retention_days: 14
compression:
  media_threshold: 90
  docs_threshold: 10
defaults:
  create:
    comprehensive: true
    force: false
    threads: 8
`
	
	configPath := filepath.Join(tempDir, ".7zarch-go-config")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}
	
	// Set HOME to temp directory
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", oldHome)
	
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	if cfg == nil {
		t.Fatal("Expected config to be non-nil")
	}
	
	// Test loaded values
	if cfg.Storage.ManagedPath != "/custom/path" {
		t.Errorf("Expected ManagedPath '/custom/path', got %q", cfg.Storage.ManagedPath)
	}
	
	if cfg.Storage.UseManagedDefault {
		t.Error("Expected UseManagedDefault to be false")
	}
	
	if cfg.Storage.RetentionDays != 14 {
		t.Errorf("Expected RetentionDays 14, got %d", cfg.Storage.RetentionDays)
	}
	
	if cfg.Compression.MediaThreshold != 90 {
		t.Errorf("Expected MediaThreshold 90, got %d", cfg.Compression.MediaThreshold)
	}
	
	if cfg.Defaults.Create.Threads != 8 {
		t.Errorf("Expected Threads 8, got %d", cfg.Defaults.Create.Threads)
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	
	// Create invalid YAML
	configContent := `invalid: yaml: content: [unclosed
`
	
	configPath := filepath.Join(tempDir, ".7zarch-go-config")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}
	
	// Set HOME to temp directory
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", oldHome)
	
	cfg, err := Load()
	if err != nil {
		t.Errorf("Expected no error (defaults returned), got: %v", err)
	}
	if cfg == nil {
		t.Error("Expected default config to be returned for invalid YAML")
	}
	
	// Should have default values when YAML is invalid
	if cfg.Storage.RetentionDays != 30 {
		t.Errorf("Expected default RetentionDays 30, got %d", cfg.Storage.RetentionDays)
	}
}