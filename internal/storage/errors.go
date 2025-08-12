package storage

import (
	"fmt"
	"strings"
	"time"
)

// Standard error types for MAS operations per 7EP-0004

// ArchiveNotFoundError indicates the requested archive doesn't exist
type ArchiveNotFoundError struct {
	ID string
}

func (e *ArchiveNotFoundError) Error() string {
	return fmt.Sprintf("Archive '%s' not found.\n💡 Use '7zarch-go list' to see available archives", e.ID)
}

// AmbiguousIDError indicates multiple archives match the given ID
type AmbiguousIDError struct {
	ID      string
	Matches []*Archive
}

func (e *AmbiguousIDError) Error() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Multiple archives match '%s':\n", e.ID))
	for i, archive := range e.Matches {
		location := "managed"
		if !archive.Managed {
			location = "external"
		}
		age := humanizeAge(archive.Created)
		size := humanizeSize(archive.Size)
		
		sb.WriteString(fmt.Sprintf("[%d] %s %s (%s, %s, %s)\n",
			i+1, 
			archive.UID[:8], 
			archive.Name,
			location,
			size,
			age))
	}
	sb.WriteString("\nPlease specify full ULID or use a longer prefix")
	return sb.String()
}

// RegistryError indicates a database operation failure
type RegistryError struct {
	Operation string
	Cause     error
}

func (e *RegistryError) Error() string {
	return fmt.Sprintf("Registry operation '%s' failed: %v\n💡 Run '7zarch-go db status' to check registry health", 
		e.Operation, e.Cause)
}

// FileVerificationError indicates archive file issues
type FileVerificationError struct {
	Archive *Archive
	Issue   string
}

func (e *FileVerificationError) Error() string {
	return fmt.Sprintf("Archive '%s' verification failed: %s\n💡 Use '7zarch-go db verify' to check all archives",
		e.Archive.Name, e.Issue)
}

// Helper functions for error messages

func humanizeAge(t time.Time) string {
	dur := time.Since(t)
	switch {
	case dur < time.Hour:
		return fmt.Sprintf("%dm ago", int(dur.Minutes()))
	case dur < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(dur.Hours()))
	case dur < 7*24*time.Hour:
		return fmt.Sprintf("%dd ago", int(dur.Hours()/24))
	case dur < 30*24*time.Hour:
		return fmt.Sprintf("%dw ago", int(dur.Hours()/(24*7)))
	default:
		return fmt.Sprintf("%dm ago", int(dur.Hours()/(24*30)))
	}
}

func humanizeSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// Error helpers for consistent messaging

// NotFoundHelp returns helpful suggestions for not found errors
func NotFoundHelp(id string) string {
	return fmt.Sprintf(`Archive not found: %s

Possible causes:
• Archive was deleted (check: 7zarch-go trash list)
• Typo in the ID (check: 7zarch-go list)
• Archive exists but not registered (fix: 7zarch-go import <path>)`, id)
}

// AmbiguousHelp returns helpful suggestions for ambiguous matches
func AmbiguousHelp() string {
	return `💡 Tips for resolving ambiguous matches:
• Use more characters of the ULID (e.g., 01K2E33 instead of 01K)
• Use the full ULID shown in brackets
• Use the archive name if unique`
}

// VerificationHelp returns helpful suggestions for verification failures
func VerificationHelp(issue string) string {
	switch issue {
	case "file not found":
		return `File not found. Possible fixes:
• If moved: 7zarch-go move <id> --reattach <new-path>
• If deleted: 7zarch-go list --missing to see all missing archives
• If on network drive: Check mount status`
	case "checksum mismatch":
		return `Checksum mismatch detected. This archive may be corrupted.
• Create new archive: 7zarch-go create <source>
• Restore from backup if available
• Run full verification: 7zarch-go test <archive>`
	default:
		return "Run '7zarch-go db verify' for comprehensive registry check"
	}
}