package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

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
			// Initialize registry via storage manager (reuses managed path config)
			mgr, err := storage.NewManager(defaultManagedPath())
			if err != nil { return fmt.Errorf("failed to init storage: %w", err) }
			defer mgr.Close()

			arc, err := resolveArchive(mgr, id)
			if err != nil { return err }

			printArchive(arc)
			return nil
		},
	}
	return cmd
}

func resolveArchive(mgr *storage.Manager, sel string) (*storage.Archive, error) {
	// 1) numeric id
	if n, err := strconv.ParseInt(sel, 10, 64); err == nil {
		list, _ := mgr.List()
		for _, a := range list { if a.ID == n { return a, nil } }
	}
	// 2) uid prefix
	list, _ := mgr.List()
	var uidCandidates []*storage.Archive
	for _, a := range list { if strings.HasPrefix(strings.ToLower(a.UID), strings.ToLower(sel)) { uidCandidates = append(uidCandidates, a) } }
	if len(uidCandidates) == 1 { return uidCandidates[0], nil }
	if len(uidCandidates) > 1 { return nil, fmt.Errorf("ambiguous uid prefix; matches %d entries", len(uidCandidates)) }
	// 3) checksum prefix
	var chkCandidates []*storage.Archive
	for _, a := range list { if strings.HasPrefix(strings.ToLower(a.Checksum), strings.ToLower(sel)) { chkCandidates = append(chkCandidates, a) } }
	if len(chkCandidates) == 1 { return chkCandidates[0], nil }
	if len(chkCandidates) > 1 { return nil, fmt.Errorf("ambiguous checksum prefix; matches %d entries", len(chkCandidates)) }
	// 4) name exact
	if a, err := mgr.Get(sel); err == nil { return a, nil }
	return nil, errors.New("archive not found; try numeric id, uid prefix, checksum prefix, or exact name")
}

func printArchive(a *storage.Archive) {
	fmt.Printf("UID:        %s\n", a.UID)
	fmt.Printf("Name:       %s\n", a.Name)
	fmt.Printf("Path:       %s\n", a.Path)
	fmt.Printf("Managed:    %t\n", a.Managed)
	fmt.Printf("Status:     %s\n", a.Status)
	fmt.Printf("Size:       %d\n", a.Size)
	fmt.Printf("Created:    %s\n", a.Created.Format("2006-01-02 15:04:05"))
	if a.Checksum != "" { fmt.Printf("Checksum:   %s\n", a.Checksum) }
	if a.Profile != "" { fmt.Printf("Profile:    %s\n", a.Profile) }
	if a.Uploaded { fmt.Printf("Uploaded:   %t (%s)\n", a.Uploaded, a.Destination) }
}

// defaultManagedPath returns the default from config if available; fallback to ~/.7zarch-go
func defaultManagedPath() string {
	// For now, keep simple and match internal defaults
	return "~/.7zarch-go"
}

