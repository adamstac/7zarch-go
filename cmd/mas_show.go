package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/adamstac/7zarch-go/internal/config"
	"github.com/adamstac/7zarch-go/internal/storage"
	"github.com/spf13/cobra"
)

func MasShowCmd() *cobra.Command {
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

			printArchive(arc)
			return nil
		},
	}
	return cmd
}

func printArchive(a *storage.Archive) {
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
	if a.Checksum != "" {
		fmt.Printf("Checksum:   %s\n", a.Checksum)
	}
	if a.Profile != "" {
		fmt.Printf("Profile:    %s\n", a.Profile)
	}
	if a.Uploaded {
		fmt.Printf("Uploaded:   %t (%s)\n", a.Uploaded, a.Destination)
	}
}
