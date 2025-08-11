package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRegistryCRUD(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "registry.db")
	r, err := NewRegistry(dbPath)
	if err != nil {
		t.Fatalf("NewRegistry: %v", err)
	}
	t.Cleanup(func() { r.Close(); os.RemoveAll(dir) })

	a := &Archive{Name: "test.7z", Path: "/tmp/test.7z", Size: 123, Profile: "Media"}
	if err := r.Add(a); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if a.ID == 0 { t.Fatal("expected ID to be set") }

	got, err := r.Get("test.7z")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Name != a.Name || got.Path != a.Path || got.Size != a.Size {
		t.Fatalf("Get mismatch: got %+v", got)
	}

	got.Size = 456
	if err := r.Update(got); err != nil {
		t.Fatalf("Update: %v", err)
	}

	list, err := r.List()
	if err != nil || len(list) == 0 { t.Fatalf("List: %v, n=%d", err, len(list)) }

	if err := r.Delete("test.7z"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := r.Get("test.7z"); err == nil {
		t.Fatal("expected error after delete")
	}
}

