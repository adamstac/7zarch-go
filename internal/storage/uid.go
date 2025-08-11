package storage

import (
	"crypto/rand"
	"time"

	ulid "github.com/oklog/ulid/v2"
)

// generateUID creates a proper ULID string (26 chars, lexicographically sortable by time)
func generateUID() string {
	entropy := ulid.Monotonic(rand.Reader, 0)
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
	return id.String()
}

