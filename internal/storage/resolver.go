package storage

import (
	"strconv"
	"strings"
)

// Error types defined in errors.go

// Resolver resolves user input into a single registry archive
// Resolution order: numeric ID -> exact UID -> UID prefix -> checksum prefix -> exact name
// Default minimum prefix length is 4 to avoid accidental broad matches.

type Resolver struct {
	reg             *Registry
	MinPrefixLength int
	MaxCandidates   int
}

func NewResolver(reg *Registry) *Resolver {
	// Default min ULID prefix length: 12 (covers full timestamp + 2 chars randomness)
	return &Resolver{reg: reg, MinPrefixLength: 12, MaxCandidates: 50}
}

func (r *Resolver) Resolve(input string) (*Archive, error) {
	trim := strings.TrimSpace(input)
	if trim == "" {
		return nil, &ArchiveNotFoundError{ID: input}
	}

	// 1) Numeric ID
	if isAllDigits(trim) {
		if id, err := strconv.ParseInt(trim, 10, 64); err == nil {
			if a, err := r.reg.GetByID(id); err == nil {
				return a, nil
			}
		}
	}

	// 2) Exact UID
	if a, err := r.reg.GetByUID(trim); err == nil {
		return a, nil
	}

	// 3) UID prefix (min length)
	if len(trim) >= r.MinPrefixLength {
		if matches, err := r.reg.FindByUIDPrefix(trim, r.MaxCandidates); err == nil {
			switch len(matches) {
			case 0:
				// continue
			case 1:
				return matches[0], nil
			default:
				return nil, &AmbiguousIDError{ID: input, Matches: matches}
			}
		}
	}

	// 4) Checksum prefix (min length)
	if len(trim) >= r.MinPrefixLength {
		if matches, err := r.reg.FindByChecksumPrefix(trim, r.MaxCandidates); err == nil {
			switch len(matches) {
			case 0:
				// continue
			case 1:
				return matches[0], nil
			default:
				return nil, &AmbiguousIDError{ID: input, Matches: matches}
			}
		}
	}

	// 5) Exact name
	if a, err := r.reg.Get(trim); err == nil {
		return a, nil
	}

	return nil, &ArchiveNotFoundError{ID: input}
}

// Registry returns the underlying registry for use by other components
func (r *Resolver) Registry() *Registry {
	return r.reg
}

func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}
