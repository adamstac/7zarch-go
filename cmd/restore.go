package cmd

import (
	"github.com/spf13/cobra"
)

// RestoreCmd returns the `restore` command scaffold.
// Usage: 7zarch-go restore <id>
func RestoreCmd() *cobra.Command {
	var (
		flagForce  bool
		flagDryRun bool
	)

	cmd := &cobra.Command{
		Use:   "restore <id>",
		Short: "Restore a deleted archive from trash",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Placeholder implementation; real logic will be added in subsequent commits
			cmd.Println("restore: not implemented yet")
			return nil
		},
	}

	cmd.Flags().BoolVar(&flagForce, "force", false, "Overwrite existing file if it exists at original location")
	cmd.Flags().BoolVar(&flagDryRun, "dry-run", false, "Show actions without making changes")

	return cmd
}
