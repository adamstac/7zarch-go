package main

import (
	"fmt"
	"os"

	"github.com/adamstac/7zarch-go/cmd"
	"github.com/spf13/cobra"
)

var (
	// Version information - set during build
	Version   = "0.1.0-dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "7zarch-go",
		Short: "An intelligent archive management tool",
		Long: `7zarch-go optimizes compression based on content type, 
creating, testing, and managing archives with smart defaults and TrueNAS integration (soon).

Features:
- Intelligent compression profiles for different file types
- Concurrent archive testing for 10x faster verification
- Configuration presets for common workflows
- Comprehensive mode with checksums and metadata
- Single binary distribution`,
		Version: fmt.Sprintf("%s (built %s, commit %s)", Version, BuildTime, GitCommit),
	}

	// Add commands
	rootCmd.AddCommand(cmd.CreateCmd())
	rootCmd.AddCommand(cmd.TestCmd())
	rootCmd.AddCommand(cmd.UploadCmd())
	rootCmd.AddCommand(cmd.ListCmd())
	rootCmd.AddCommand(cmd.ProfilesCmd())
	rootCmd.AddCommand(cmd.ConfigCmd())
	// MAS commands (initial)
	rootCmd.AddCommand(cmd.MasShowCmd())
	rootCmd.AddCommand(cmd.MasDbCmd())

	// Execute
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}