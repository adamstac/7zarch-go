package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Manager handles the managed storage workspace
type Manager struct {
	basePath string
	registry *Registry
}

// NewManager creates a new storage manager
func NewManager(basePath string) (*Manager, error) {
	// Expand tilde to home directory
	if strings.HasPrefix(basePath, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		basePath = filepath.Join(home, basePath[2:])
	}

	// Create the managed storage directory
	archivesPath := filepath.Join(basePath, "archives")
	if err := os.MkdirAll(archivesPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create managed storage directory: %w", err)
	}

	// Initialize the registry
	dbPath := filepath.Join(basePath, "registry.db")
	registry, err := NewRegistry(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize registry: %w", err)
	}

	// Backfill missing UIDs
	_ = registry.BackfillUIDs(func() string { return generateUID() })

	return &Manager{
		basePath: basePath,
		registry: registry,
	}, nil
}

// GetManagedPath returns the path where an archive should be stored
func (m *Manager) GetManagedPath(archiveName string) string {
	// For now, use flat organization
	// Later we can add date-based or type-based organization
	return filepath.Join(m.basePath, "archives", archiveName)
}

// Add registers a new archive in the registry
// checksum and metadata are optional; pass empty strings if not available
// managed indicates whether the file is stored under the MAS path
func (m *Manager) Add(name, path string, size int64, profile string, checksum string, metadata string, managed bool) error {
	archive := &Archive{
		UID:      generateUID(),
		Name:     name,
		Path:     path,
		Size:     size,
		Created:  time.Now(),
		Profile:  profile,
		Checksum: checksum,
		Managed:  managed,
		Status:   "present",
		Metadata: metadata,
	}
	return m.registry.Add(archive)
}

// List returns all managed archives
func (m *Manager) List() ([]*Archive, error) {
	return m.registry.List()
}

// ListNotUploaded returns archives that haven't been uploaded
func (m *Manager) ListNotUploaded() ([]*Archive, error) {
	return m.registry.ListNotUploaded()
}

// ListOlderThan returns archives older than the specified duration
func (m *Manager) ListOlderThan(duration time.Duration) ([]*Archive, error) {
	return m.registry.ListOlderThan(duration)
}

// Get retrieves an archive by name
func (m *Manager) Get(name string) (*Archive, error) {
	return m.registry.Get(name)
}

// MarkUploaded marks an archive as uploaded
func (m *Manager) MarkUploaded(name string, destination string) error {
	archive, err := m.registry.Get(name)
	if err != nil {
		return err
	}

	now := time.Now()
	archive.Uploaded = true
	archive.Destination = destination
	archive.UploadedAt = &now

	return m.registry.Update(archive)
}

// Delete removes an archive from the registry (does not delete the file)
func (m *Manager) Delete(name string) error {
	return m.registry.Delete(name)
}

// Close closes the storage manager and its registry
func (m *Manager) Close() error {
	if m.registry != nil {
		return m.registry.Close()
	}
	return nil
}

// GetBasePath returns the base path for managed storage
func (m *Manager) GetBasePath() string {
	return m.basePath
}

// GetArchivesPath returns the path where archives are stored
func (m *Manager) GetArchivesPath() string {
	return filepath.Join(m.basePath, "archives")
}

// Exists checks if an archive exists in the registry
func (m *Manager) Exists(name string) bool {
	_, err := m.registry.Get(name)
	return err == nil
}