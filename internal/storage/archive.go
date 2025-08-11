package storage

import (
	"time"
)

// Archive represents an archive entry in the registry
type Archive struct {
	ID           int64      `json:"id"`
	UID          string     `json:"uid"`
	Name         string     `json:"name"`
	Path         string     `json:"path"`
	Size         int64      `json:"size"`
	Created      time.Time  `json:"created"`
	Checksum     string     `json:"checksum,omitempty"`
	Profile      string     `json:"profile,omitempty"`     // compression profile used
	Managed      bool       `json:"managed"`
	Status       string     `json:"status"`                // present | missing | deleted
	LastSeen     *time.Time `json:"last_seen,omitempty"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
	OriginalPath string     `json:"original_path,omitempty"`
	Uploaded     bool       `json:"uploaded"`
	Destination  string     `json:"destination,omitempty"` // where it was uploaded
	UploadedAt   *time.Time `json:"uploaded_at,omitempty"`
	Metadata     string     `json:"metadata,omitempty"`    // JSON blob for extensibility
}

// IsManaged returns true if this archive is in managed storage
func (a *Archive) IsManaged() bool {
	// Archive is managed if it's in the managed storage path
	// This will be determined by the storage manager
	return a.Path != "" && len(a.Path) > 0
}

// Age returns the age of the archive
func (a *Archive) Age() time.Duration {
	return time.Since(a.Created)
}