package cmd

import (
	"github.com/spf13/cobra"
)

// TrashCmd groups trash-related subcommands.
func TrashCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trash",
		Short: "Manage deleted archives in trash",
	}
	cmd.AddCommand(trashListCmd())
	cmd.AddCommand(trashPurgeCmd())
	return cmd
}

func trashListCmd() *cobra.Command {
	var (
		flagWithinDays int
		flagBefore     string
		flagJSON       bool
	)
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List deleted archives and purge eligibility",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Println("trash list: not implemented yet")
			return nil
		},
	}
	cmd.Flags().IntVar(&flagWithinDays, "within-days", 0, "Show items purging within N days (0=all)")
	cmd.Flags().StringVar(&flagBefore, "before", "", "Only show items deleted before YYYY-MM-DD")
	cmd.Flags().BoolVar(&flagJSON, "json", false, "Output as JSON")
	return cmd
}

func trashPurgeCmd() *cobra.Command {
	var (
		flagAll    bool
		flagForce  bool
		flagDryRun bool
		flagWithin int
	)
	cmd := &cobra.Command{
		Use:   "purge",
		Short: "Permanently delete trashed archives past retention",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Println("trash purge: not implemented yet")
			return nil
		},
	}
	cmd.Flags().BoolVar(&flagAll, "all", false, "Purge all trashed archives (ignore retention)")
	cmd.Flags().BoolVar(&flagForce, "force", false, "Skip confirmation prompts")
	cmd.Flags().BoolVar(&flagDryRun, "dry-run", false, "Show actions without making changes")
	cmd.Flags().IntVar(&flagWithin, "within-days", 0, "Only purge items purging within N days (0=all)")
	return cmd
}
