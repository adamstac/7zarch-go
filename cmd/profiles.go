package cmd

import (
	"fmt"
	"sort"

	"github.com/adamstac/7zarch-go/internal/archive"
	"github.com/spf13/cobra"
)

func ProfilesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profiles",
		Short: "List available compression profiles",
		Long:  `List all available compression profiles with their settings and recommended use cases.`,
		RunE:  runProfiles,
	}

	return cmd
}

func runProfiles(cmd *cobra.Command, args []string) error {
	profiles := archive.ListProfiles()

	// Sort profiles by name for consistent output
	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].Name < profiles[j].Name
	})

	fmt.Printf("Available Compression Profiles\n")
	fmt.Printf("==============================\n\n")

	for _, profile := range profiles {
		fmt.Printf("📦 %s\n", profile.Name)
		fmt.Printf("   %s\n", profile.Description)
		fmt.Printf("   Settings: Level %d, Dictionary %s, Fast bytes %d",
			profile.Level, profile.DictionarySize, profile.FastBytes)

		if profile.SolidMode {
			fmt.Printf(", Solid mode on\n")
		} else {
			fmt.Printf(", Solid mode off\n")
		}

		// Add usage examples
		switch profile.Name {
		case "Media":
			fmt.Printf("   Best for: Video files, audio files, photos, podcasts\n")
			fmt.Printf("   Example: 7zarch-go create video-project --profile media\n")
		case "Documents":
			fmt.Printf("   Best for: Text files, source code, office documents, PDFs\n")
			fmt.Printf("   Example: 7zarch-go create source-code --profile documents\n")
		case "Balanced":
			fmt.Printf("   Best for: Mixed content, general backups\n")
			fmt.Printf("   Example: 7zarch-go create backup-folder --profile balanced\n")
		}

		fmt.Printf("\n")
	}

	fmt.Printf("Smart Compression (Default Behavior)\n")
	fmt.Printf("====================================\n")
	fmt.Printf("7zarch-go is smart by default - it analyzes your content and automatically\n")
	fmt.Printf("selects the best profile for optimal performance:\n")
	fmt.Printf("   7zarch-go create my-folder\n\n")

	fmt.Printf("Manual Profile Override\n")
	fmt.Printf("=======================\n")
	fmt.Printf("Force a specific profile when you know what you want:\n")
	fmt.Printf("   7zarch-go create my-folder --profile media\n\n")

	fmt.Printf("Traditional Compression Level\n")
	fmt.Printf("=============================\n")
	fmt.Printf("Use traditional compression level (0-9) to disable smart behavior:\n")
	fmt.Printf("   7zarch-go create my-folder --compression 9\n\n")

	fmt.Printf("Presets\n")
	fmt.Printf("=======\n")
	fmt.Printf("Use predefined combinations of settings for common workflows:\n")
	fmt.Printf("   7zarch-go create my-folder --preset podcast\n")
	fmt.Printf("   7zarch-go create my-folder --preset backup\n")
	fmt.Printf("   7zarch-go create my-folder --preset source_code\n")
	fmt.Printf("\nRun '7zarch-go config show' to see available presets in your configuration.\n")

	return nil
}
