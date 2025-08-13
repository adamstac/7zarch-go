package errors

import (
	"fmt"
	"strings"
)

// ValidationError represents invalid input validation
type ValidationError struct {
	Field   string
	Value   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Invalid %s '%s': %s", e.Field, e.Value, e.Message)
}

// NotFoundError represents a resource that couldn't be located
type NotFoundError struct {
	Resource string
	ID       string
	Suggestions []string
}

func (e *NotFoundError) Error() string {
	base := fmt.Sprintf("%s '%s' not found.", e.Resource, e.ID)
	if len(e.Suggestions) > 0 {
		base += fmt.Sprintf(" Try: %s", strings.Join(e.Suggestions, ", "))
	}
	return base
}

// InvalidOperationError represents an operation that cannot be performed
type InvalidOperationError struct {
	Operation string
	Resource  string
	Reason    string
}

func (e *InvalidOperationError) Error() string {
	return fmt.Sprintf("Cannot %s %s: %s", e.Operation, e.Resource, e.Reason)
}

// FileSystemError represents file operation failures
type FileSystemError struct {
	Path      string
	Operation string
	Err       error
}

func (e *FileSystemError) Error() string {
	return fmt.Sprintf("Failed to %s '%s': %v", e.Operation, e.Path, e.Err)
}

func (e *FileSystemError) Unwrap() error {
	return e.Err
}

// DatabaseError represents database operation failures
type DatabaseError struct {
	Operation string
	Table     string
	Err       error
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("Database %s failed on %s: %v", e.Operation, e.Table, e.Err)
}

func (e *DatabaseError) Unwrap() error {
	return e.Err
}

// ConfigurationError represents configuration issues
type ConfigurationError struct {
	Setting string
	Value   string
	Message string
}

func (e *ConfigurationError) Error() string {
	if e.Value != "" {
		return fmt.Sprintf("Configuration error for %s='%s': %s", e.Setting, e.Value, e.Message)
	}
	return fmt.Sprintf("Configuration error for %s: %s", e.Setting, e.Message)
}

// Helper functions for common error scenarios

// NewArchiveNotFound creates a standard archive not found error
func NewArchiveNotFound(id string) error {
	return &NotFoundError{
		Resource:    "Archive",
		ID:          id,
		Suggestions: []string{"use 'list' to see available archives", "check the archive ID"},
	}
}

// NewInvalidPath creates a standard invalid path error
func NewInvalidPath(path string, reason string) error {
	return &ValidationError{
		Field:   "path",
		Value:   path,
		Message: reason,
	}
}

// NewPermissionDenied creates a standard permission error
func NewPermissionDenied(path string, operation string) error {
	return &FileSystemError{
		Path:      path,
		Operation: operation,
		Err:       fmt.Errorf("permission denied"),
	}
}