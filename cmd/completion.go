package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/adamstac/7zarch-go/internal/config"
	"github.com/adamstac/7zarch-go/internal/storage"
	"github.com/spf13/cobra"
)

func CompletionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts for 7zarch-go.

The completion script for bash can be generated with:
	7zarch-go completion bash

Add to your ~/.bashrc:
	source <(7zarch-go completion bash)

For zsh, add to ~/.zshrc:
	autoload -U compinit && compinit
	source <(7zarch-go completion zsh)

For fish, add to ~/.config/fish/config.fish:
	7zarch-go completion fish | source

For PowerShell, add to your profile:
	7zarch-go completion powershell | Out-String | Invoke-Expression`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				return cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				return cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletion(os.Stdout)
			default:
				return fmt.Errorf("unsupported shell: %s", args[0])
			}
		},
	}

	return cmd
}

// completeArchiveIDs provides completion for archive IDs (UID prefixes, names)
func completeArchiveIDs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return completeArchiveIDsWithFilter(cmd, args, toComplete, nil)
}

// completeDeletedArchiveIDs provides completion for deleted archive IDs only (for restore command)
func completeDeletedArchiveIDs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return completeArchiveIDsWithFilter(cmd, args, toComplete, func(a *storage.Archive) bool {
		return a.Status == "deleted"
	})
}

// completeArchiveIDsWithFilter provides completion for archive IDs with optional filtering
func completeArchiveIDsWithFilter(cmd *cobra.Command, args []string, toComplete string, filter func(*storage.Archive) bool) ([]string, cobra.ShellCompDirective) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	mgr, err := storage.NewManager(cfg.Storage.ManagedPath)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	defer mgr.Close()

	completions := make([]string, 0, 50)

	// Use channel to collect results and handle timeout
	resultCh := make(chan []string, 3)
	errorCh := make(chan error, 3)

	// Launch concurrent lookups
	go func() {
		if matches := getUIDPrefixCompletions(mgr, toComplete, filter); matches != nil {
			select {
			case resultCh <- matches:
			case <-ctx.Done():
			}
		}
	}()

	go func() {
		if matches := getNameCompletions(mgr, toComplete, filter); matches != nil {
			select {
			case resultCh <- matches:
			case <-ctx.Done():
			}
		}
	}()

	go func() {
		if matches := getChecksumPrefixCompletions(mgr, toComplete, filter); matches != nil {
			select {
			case resultCh <- matches:
			case <-ctx.Done():
			}
		}
	}()

	// Collect results with timeout
	for i := 0; i < 3; i++ {
		select {
		case matches := <-resultCh:
			completions = append(completions, matches...)
		case err := <-errorCh:
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
		case <-ctx.Done():
			// Timeout - return what we have so far
			break
		}
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

// getUIDPrefixCompletions returns ULID prefix matches
func getUIDPrefixCompletions(mgr *storage.Manager, toComplete string, filter func(*storage.Archive) bool) []string {
	if len(toComplete) < 4 { // Too short for ULID prefix
		return nil
	}

	resolver := storage.NewResolver(mgr.Registry())
	// Temporarily lower min prefix length for completion
	resolver.MinPrefixLength = 4

	matches, err := mgr.Registry().FindByUIDPrefix(toComplete, 25)
	if err != nil {
		return nil
	}

	completions := make([]string, 0, len(matches))
	for _, archive := range matches {
		if archive.UID != "" && (filter == nil || filter(archive)) {
			// Provide UID with description
			desc := fmt.Sprintf("%s\t%s", archive.UID, archive.Name)
			completions = append(completions, desc)
		}
	}

	return completions
}

// getNameCompletions returns archive name matches
func getNameCompletions(mgr *storage.Manager, toComplete string, filter func(*storage.Archive) bool) []string {
	if toComplete == "" {
		return nil
	}

	archives, err := mgr.Registry().List()
	if err != nil {
		return nil
	}

	completions := make([]string, 0, 25)
	for _, archive := range archives {
		if strings.HasPrefix(archive.Name, toComplete) && (filter == nil || filter(archive)) {
			desc := fmt.Sprintf("%s\t%s", archive.Name, safePrefix(archive.UID, 8))
			completions = append(completions, desc)
			if len(completions) >= 25 {
				break
			}
		}
	}

	return completions
}

// getChecksumPrefixCompletions returns checksum prefix matches
func getChecksumPrefixCompletions(mgr *storage.Manager, toComplete string, filter func(*storage.Archive) bool) []string {
	if len(toComplete) < 8 { // Too short for checksum prefix
		return nil
	}

	matches, err := mgr.Registry().FindByChecksumPrefix(toComplete, 10)
	if err != nil {
		return nil
	}

	completions := make([]string, 0, len(matches))
	for _, archive := range matches {
		if archive.Checksum != "" && (filter == nil || filter(archive)) {
			desc := fmt.Sprintf("%s\t%s (%s)",
				safePrefix(archive.Checksum, 12),
				archive.Name,
				safePrefix(archive.UID, 8))
			completions = append(completions, desc)
		}
	}

	return completions
}

// completeCommands provides completion for subcommands
func completeCommands(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var completions []string

	for _, subCmd := range cmd.Root().Commands() {
		if strings.HasPrefix(subCmd.Name(), toComplete) {
			completions = append(completions, subCmd.Name()+"\t"+subCmd.Short)
		}
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}
