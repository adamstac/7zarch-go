package storage

import (
	"crypto/rand"
	"encoding/base32"
)

// generateUID creates a lexicographically sortable, ULID-like ID without pulling a new dep.
// Not a full ULID implementation, but sufficient for local uniqueness and CLI ergonomics.
func generateUID() string {
	// 16 random bytes; base32 no padding gives ~26 chars, upper-case.
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	enc := base32.StdEncoding.WithPadding(base32.NoPadding)
	return enc.EncodeToString(b)
}

