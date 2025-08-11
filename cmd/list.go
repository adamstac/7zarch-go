package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List archives on TrueNAS",
		Long:  `List all archives stored on TrueNAS.`,
		RunE:  runList,
	}

	// Add flags
	cmd.Flags().String("path", "/", "Path to list")
	cmd.Flags().String("storage", "truenas", "Storage backend to use")
	cmd.Flags().Bool("details", false, "Show detailed information")

	return cmd
}

func runList(cmd *cobra.Command, args []string) error {
	path, _ := cmd.Flags().GetString("path")
	details, _ := cmd.Flags().GetBool("details")
	
	fmt.Printf("Listing archives on TrueNAS...\n")
	fmt.Printf("Path: %s\n", path)
	
	if details {
		fmt.Printf("(Detailed mode)\n")
	}
	
	// TODO: Implement TrueNAS listing
	fmt.Printf("\n⚠️  List functionality coming soon!\n")
	fmt.Printf("This will connect to TrueNAS and list archives.\n")
	
	return nil
}