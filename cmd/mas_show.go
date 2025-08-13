package cmd

import (
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/adamstac/7zarch-go/internal/config"
	"github.com/adamstac/7zarch-go/internal/storage"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func MasShowCmd() *cobra.Command {
	var (
		verify bool
		output string
	)
	cmd := &cobra.Command{
		Use:   "show <id>",
		Short: "Show archive details by ID (uid, checksum prefix, numeric id, or name)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			cfg, _ := config.Load()
			mgr, err := storage.NewManager(cfg.Storage.ManagedPath)
			if err != nil {
				return fmt.Errorf("failed to init storage (path=%q): %w", cfg.Storage.ManagedPath, err)
			}
			defer mgr.Close()

			resolver := storage.NewResolver(mgr.Registry())
			arc, err := resolver.Resolve(id)
			if err != nil {
				if amb, ok := err.(*storage.AmbiguousIDError); ok {
					printAmbiguousOptions(amb)
				}
				return err
			}

			// File existence verification + last_seen/status update
			now := time.Now()
			if _, statErr := os.Stat(arc.Path); statErr == nil {
				arc.Status = "present"
			} else {
				arc.Status = "missing"
			}
			arc.LastSeen = &now
			_ = mgr.Registry().Update(arc)

			if output != "" {
				return outputArchive(arc, output, verify)
			}

			printArchive(arc, verify)
			return nil
		},
	}
	cmd.Flags().BoolVar(&verify, "verify", false, "Verify checksum against file (slower)")
	cmd.Flags().StringVar(&output, "output", "", "Output format: json|csv|yaml (default: human-readable)")
	return cmd
}

func printArchive(a *storage.Archive, verify bool) {
	status := a.Status
	switch status {
	case "present":
		status += " ✓"
	case "missing":
		status += " ⚠️"
	}
	fmt.Printf("UID:        %s\n", a.UID)
	fmt.Printf("Name:       %s\n", a.Name)
	fmt.Printf("Path:       %s\n", a.Path)
	fmt.Printf("Managed:    %t\n", a.Managed)
	fmt.Printf("Status:     %s\n", status)
	fmt.Printf("Size:       %d\n", a.Size)
	fmt.Printf("Created:    %s\n", a.Created.Format("2006-01-02 15:04:05"))
	// Checksum line
	if a.Checksum == "" {
		fmt.Printf("Checksum:   (none)\n")
	} else if verify && a.Status == "present" {
		computed, err := computeSHA256(a.Path)
		if err == nil && computed == a.Checksum {
			fmt.Printf("Checksum:   %s (verified ✓)\n", a.Checksum)
		} else if err == nil {
			fmt.Printf("Checksum:   %s (mismatch ⚠️)\n", a.Checksum)
		} else {
			fmt.Printf("Checksum:   %s (verify error: %v)\n", a.Checksum, err)
		}
	} else {
		fmt.Printf("Checksum:   %s\n", a.Checksum)
	}
	if a.Profile != "" {
		fmt.Printf("Profile:    %s\n", a.Profile)
	}
	if a.Uploaded {
		fmt.Printf("Uploaded:   %t (%s)\n", a.Uploaded, a.Destination)
	}
}

func computeSHA256(path string) (string, error) {
	// #nosec G304: path originates from registry-managed archive object
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func printAmbiguousOptions(amb *storage.AmbiguousIDError) {
	fmt.Printf("Multiple archives match '%s':\n", amb.ID)
	for i, a := range amb.Matches {
		loc := "external"
		if a.Managed {
			loc = "managed"
		}
		age := time.Since(a.Created).Round(time.Hour)
		fmt.Printf("[%d] %s  %s  (%s, %.1f MB, %s ago)\n", i+1, safePrefix(a.UID, 8), a.Name, loc, float64(a.Size)/(1024*1024), age)
	}
	fmt.Println("Please specify a longer prefix or the full UID.")
}

func safePrefix(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

// outputArchive outputs single archive in machine-readable format
func outputArchive(a *storage.Archive, format string, verify bool) error {
	// Create enriched archive with verified checksum if requested
	archiveData := *a
	if verify && a.Status == "present" && a.Checksum != "" {
		if computed, err := computeSHA256(a.Path); err == nil {
			metadata := fmt.Sprintf(`{"checksum_verified":%t,"computed_checksum":"%s"}`, computed == a.Checksum, computed)
			archiveData.Metadata = metadata
		}
	}

	switch format {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(&archiveData)
	case "csv":
		return outputArchiveCSV(&archiveData)
	case "yaml":
		enc := yaml.NewEncoder(os.Stdout)
		defer enc.Close()
		return enc.Encode(&archiveData)
	default:
		return fmt.Errorf("unsupported output format: %s (supported: json, csv, yaml)", format)
	}
}

func outputArchiveCSV(a *storage.Archive) error {
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{
		"uid", "name", "path", "size", "created", "checksum", "profile",
		"managed", "status", "last_seen", "deleted_at", "original_path",
		"uploaded", "destination", "uploaded_at",
	}); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data row
	var lastSeen, deletedAt, uploadedAt string
	if a.LastSeen != nil {
		lastSeen = a.LastSeen.Format(time.RFC3339)
	}
	if a.DeletedAt != nil {
		deletedAt = a.DeletedAt.Format(time.RFC3339)
	}
	if a.UploadedAt != nil {
		uploadedAt = a.UploadedAt.Format(time.RFC3339)
	}

	row := []string{
		a.UID,
		a.Name,
		a.Path,
		fmt.Sprintf("%d", a.Size),
		a.Created.Format(time.RFC3339),
		a.Checksum,
		a.Profile,
		fmt.Sprintf("%t", a.Managed),
		a.Status,
		lastSeen,
		deletedAt,
		a.OriginalPath,
		fmt.Sprintf("%t", a.Uploaded),
		a.Destination,
		uploadedAt,
	}

	return writer.Write(row)
}
