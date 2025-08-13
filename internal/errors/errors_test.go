package errors

import (
	"fmt"
	"testing"
)

func TestNotFoundError(t *testing.T) {
	err := &NotFoundError{
		Resource:    "Archive",
		ID:          "test123",
		Suggestions: []string{"check spelling", "use full ID"},
	}

	expected := "Archive 'test123' not found. Try: check spelling, use full ID"
	if err.Error() != expected {
		t.Errorf("Expected %q, got %q", expected, err.Error())
	}
}

func TestNotFoundError_NoSuggestions(t *testing.T) {
	err := &NotFoundError{
		Resource: "File",
		ID:       "missing.txt",
	}

	expected := "File 'missing.txt' not found."
	if err.Error() != expected {
		t.Errorf("Expected %q, got %q", expected, err.Error())
	}
}

func TestValidationError(t *testing.T) {
	err := &ValidationError{
		Field:   "profile",
		Value:   "invalid",
		Message: "unknown profile. Available: media, documents, balanced",
	}

	expected := "Invalid profile 'invalid': unknown profile. Available: media, documents, balanced"
	if err.Error() != expected {
		t.Errorf("Expected %q, got %q", expected, err.Error())
	}
}

func TestInvalidOperationError(t *testing.T) {
	err := &InvalidOperationError{
		Operation: "delete",
		Resource:  "archive",
		Reason:    "archive is already deleted",
	}

	expected := "Cannot delete archive: archive is already deleted"
	if err.Error() != expected {
		t.Errorf("Expected %q, got %q", expected, err.Error())
	}
}

func TestFileSystemError(t *testing.T) {
	err := &FileSystemError{
		Path:      "/path/to/file",
		Operation: "read",
		Err:       fmt.Errorf("permission denied"),
	}

	expected := "Failed to read '/path/to/file': permission denied"
	if err.Error() != expected {
		t.Errorf("Expected %q, got %q", expected, err.Error())
	}
}

func TestNewArchiveNotFound(t *testing.T) {
	err := NewArchiveNotFound("abc123")

	// Should be a NotFoundError
	if notFoundErr, ok := err.(*NotFoundError); ok {
		if notFoundErr.Resource != "Archive" {
			t.Errorf("Expected resource 'Archive', got %q", notFoundErr.Resource)
		}
		if notFoundErr.ID != "abc123" {
			t.Errorf("Expected ID 'abc123', got %q", notFoundErr.ID)
		}
		if len(notFoundErr.Suggestions) == 0 {
			t.Error("Expected suggestions to be provided")
		}
	} else {
		t.Errorf("Expected *NotFoundError, got %T", err)
	}

	expected := "Archive 'abc123' not found. Try: use 'list' to see available archives, check the archive ID"
	if err.Error() != expected {
		t.Errorf("Expected %q, got %q", expected, err.Error())
	}
}