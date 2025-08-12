package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/adamstac/7zarch-go/internal/config"
	"github.com/adamstac/7zarch-go/internal/storage"
	"github.com/spf13/cobra"
)

func MasShowCmd() *cobra.Command {
	var verify bool
	cmd := &cobra.Command{
		Use:   "show <id>",
		Short: "Show archive details by ID (uid, checksum prefix, numeric id, or name)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			cfg, _ := config.Load()
			mgr, err := storage.NewManager(cfg.Storage.ManagedPath)
			if err != nil {
				return fmt.Errorf("failed to init storage: %w", err)
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

			printArchive(arc, verify)
			return nil
		},
	}
	cmd.Flags().BoolVar(&verify, "verify", false, "Verify checksum against file (slower)")
	return cmd
}

func printArchive(a *storage.Archive, verify bool) {
	status := a.Status
	if status == "present" {
		status += " ✓"
	} else if status == "missing" {
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
