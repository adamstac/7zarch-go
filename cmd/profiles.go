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

	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "Available Compression Profiles\n")
	fmt.Fprintf(out, "==============================\n\n")

	if len(profiles) == 0 {
		fmt.Fprintf(out, "No compression profiles found.\n")
		fmt.Fprintf(out, "Tip: Use smart defaults or specify --compression to override.\n")
		return nil
	}

	for _, profile := range profiles {
<<<<<<< HEAD
		fmt.Fprintf(out, "📦 %s\n", profile.Name)
		fmt.Fprintf(out, "   %s\n", profile.Description)
		fmt.Fprintf(out, "   Settings: Level %d, Dictionary %s, Fast bytes %d",
=======
		fmt.Printf("📦 %s\n", profile.Name)
		fmt.Printf("   %s\n", profile.Description)
		fmt.Printf("   Settings: Level %d, Dictionary %s, Fast bytes %d",
>>>>>>> origin/main
			profile.Level, profile.DictionarySize, profile.FastBytes)

		if profile.SolidMode {
			fmt.Fprintf(out, ", Solid mode on\n")
		} else {
			fmt.Fprintf(out, ", Solid mode off\n")
		}

		// Add usage examples
		switch profile.Name {
		case "Media":
			fmt.Fprintf(out, "   Best for: Video files, audio files, photos, podcasts\n")
			fmt.Fprintf(out, "   Example: 7zarch-go create video-project --profile media\n")
		case "Documents":
			fmt.Fprintf(out, "   Best for: Text files, source code, office documents, PDFs\n")
			fmt.Fprintf(out, "   Example: 7zarch-go create source-code --profile documents\n")
		case "Balanced":
			fmt.Fprintf(out, "   Best for: Mixed content, general backups\n")
			fmt.Fprintf(out, "   Example: 7zarch-go create backup-folder --profile balanced\n")
		}

<<<<<<< HEAD
		fmt.Fprintf(out, "\n")
	}

	fmt.Fprintf(out, "Smart Compression (Default Behavior)\n")
	fmt.Fprintf(out, "====================================\n")
	fmt.Fprintf(out, "7zarch-go is smart by default - it analyzes your content and automatically\n")
	fmt.Fprintf(out, "selects the best profile for optimal performance:\n")
	fmt.Fprintf(out, "   7zarch-go create my-folder\n\n")

	fmt.Fprintf(out, "Manual Profile Override\n")
	fmt.Fprintf(out, "=======================\n")
	fmt.Fprintf(out, "Force a specific profile when you know what you want:\n")
	fmt.Fprintf(out, "   7zarch-go create my-folder --profile media\n\n")

	fmt.Fprintf(out, "Traditional Compression Level\n")
	fmt.Fprintf(out, "=============================\n")
	fmt.Fprintf(out, "Use traditional compression level (0-9) to disable smart behavior:\n")
	fmt.Fprintf(out, "   7zarch-go create my-folder --compression 9\n\n")

	fmt.Fprintf(out, "Presets\n")
	fmt.Fprintf(out, "=======\n")
	fmt.Fprintf(out, "Use predefined combinations of settings for common workflows:\n")
	fmt.Fprintf(out, "   7zarch-go create my-folder --preset podcast\n")
	fmt.Fprintf(out, "   7zarch-go create my-folder --preset backup\n")
	fmt.Fprintf(out, "   7zarch-go create my-folder --preset source_code\n")
	fmt.Fprintf(out, "\nRun '7zarch-go config show' to see available presets in your configuration.\n")
=======
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
>>>>>>> origin/main

	return nil
}
