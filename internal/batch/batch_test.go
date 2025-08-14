package batch

import (
	"context"
	"testing"
	"time"

	"github.com/adamstac/7zarch-go/internal/storage"
)

// Mock registry for testing
type mockRegistry struct {
	deleteFunc func(uid string) error
	updateFunc func(archive *storage.Archive) error
}

func (m *mockRegistry) Delete(uid string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(uid)
	}
	return nil
}

func (m *mockRegistry) Update(archive *storage.Archive) error {
	if m.updateFunc != nil {
		return m.updateFunc(archive)
	}
	return nil
}

// Mock manager for testing
type mockManager struct {
	registry   *mockRegistry
	basePath   string
	deleteFunc func(uid string) error
}

func (m *mockManager) Delete(uid string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(uid)
	}
	return m.registry.Delete(uid)
}

func (m *mockManager) Registry() *storage.Registry {
	// This is a hack for testing - we can't return a proper *storage.Registry
	// In real usage, this works because storage.Manager.Registry() returns *storage.Registry
	return nil
}

func (m *mockManager) GetBasePath() string {
	return m.basePath
}

func TestProcessor_Move_Success(t *testing.T) {
	// Skip this test since it requires actual file system operations
	// Focus on testing the core batch processing logic instead
	t.Skip("Move operations require file system access, tested in integration tests")
}

func TestProcessor_Delete_Success(t *testing.T) {
	registry := &mockRegistry{}
	manager := &mockManager{
		registry: registry,
		basePath: "/managed",
	}
	processor := NewProcessor(manager)

	archives := []*storage.Archive{
		{UID: "test1", Name: "archive1.7z"},
	}

	var updates []ProgressUpdate
	callback := func(update ProgressUpdate) {
		updates = append(updates, update)
	}

	ctx := context.Background()
	err := processor.Delete(ctx, archives, callback)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(updates) == 0 {
		t.Fatal("expected progress updates")
	}

	finalUpdate := updates[len(updates)-1]
	if finalUpdate.Total != 1 || finalUpdate.Completed != 1 {
		t.Errorf("expected Total=1, Completed=1, got Total=%d, Completed=%d",
			finalUpdate.Total, finalUpdate.Completed)
	}
}

func TestProcessor_Move_WithErrors(t *testing.T) {
	// Skip this test since it requires actual file system operations
	t.Skip("Move operations require file system access, tested in integration tests")
}

func TestProcessor_ContextCancellation(t *testing.T) {
	registry := &mockRegistry{
		deleteFunc: func(uid string) error {
			time.Sleep(100 * time.Millisecond) // Simulate slow operation
			return nil
		},
	}
	manager := &mockManager{
		registry: registry,
		basePath: "/managed",
	}
	processor := NewProcessor(manager)

	archives := []*storage.Archive{
		{UID: "test1", Name: "archive1.7z"},
		{UID: "test2", Name: "archive2.7z"},
		{UID: "test3", Name: "archive3.7z"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := processor.Delete(ctx, archives, nil)

	if err == nil || err != context.DeadlineExceeded {
		t.Fatalf("expected context deadline exceeded, got: %v", err)
	}
}

func TestProcessor_EmptyArchives(t *testing.T) {
	registry := &mockRegistry{}
	manager := &mockManager{
		registry: registry,
		basePath: "/managed",
	}
	processor := NewProcessor(manager)

	var archives []*storage.Archive

	err := processor.Delete(context.Background(), archives, nil)

	if err != nil {
		t.Fatalf("unexpected error with empty archives: %v", err)
	}
}

func TestProcessor_SetConcurrency(t *testing.T) {
	registry := &mockRegistry{}
	manager := &mockManager{
		registry: registry,
		basePath: "/managed",
	}
	processor := NewProcessor(manager)

	// Test default concurrency
	if processor.concurrent != 4 {
		t.Errorf("expected default concurrency 4, got %d", processor.concurrent)
	}

	// Test setting concurrency
	processor.SetConcurrency(8)
	if processor.concurrent != 8 {
		t.Errorf("expected concurrency 8, got %d", processor.concurrent)
	}

	// Test invalid concurrency (should not change)
	processor.SetConcurrency(0)
	if processor.concurrent != 8 {
		t.Errorf("expected concurrency to remain 8, got %d", processor.concurrent)
	}

	processor.SetConcurrency(-1)
	if processor.concurrent != 8 {
		t.Errorf("expected concurrency to remain 8, got %d", processor.concurrent)
	}
}