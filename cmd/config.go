package cmd

import (
	"fmt"
	"os"

	"github.com/adamstac/7zarch-go/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func ConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
		Long:  `Manage 7zarch-go configuration file and settings.`,
	}

	cmd.AddCommand(configInitCmd())
	cmd.AddCommand(configShowCmd())

	return cmd
}

func configInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create default config file",
		Long:  `Create a default configuration file at ~/.7zarch-go-config with examples and documentation.`,
		RunE:  runConfigInit,
	}

	return cmd
}

func configShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		Long:  `Display the current configuration including defaults and any custom settings.`,
		RunE:  runConfigShow,
	}

	return cmd
}

func runConfigInit(cmd *cobra.Command, args []string) error {
	configPath, err := config.ConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("Config file already exists: %s\n", configPath)
		fmt.Printf("Remove it first if you want to recreate it.\n")
		return nil
	}

	// Create default config with comments
	configContent := `# 7zarch-go Configuration
# Convention over configuration - smart defaults work great!
# Only customize what you need.

# Compression behavior
compression:
  # Smart compression analyzes your content and picks optimal settings
  # Set to false to always use manual compression levels
  smart_default: true
  
  # Default compression level when smart mode is disabled (0-9)
  level: 9
  
  # Override smart recommendations for specific content types
  media_threshold: 70    # % media files needed to trigger media profile
  docs_threshold: 60     # % document files needed to trigger docs profile

# Default flags for commands
defaults:
  create:
    comprehensive: false   # Create .log and .sha256 files by default
    force: false          # Overwrite existing archives by default
    threads: 0            # 0 = auto-detect CPU cores
  
  test:
    concurrent: 5         # Default concurrent archive tests
    verbose: false        # Show detailed test output by default

# Output and display
ui:
  # Show educational content analysis on every create
  show_analysis: true
  
  # Show optimization tips when not using optimal settings
  show_tips: true
  
  # Progress display style: "spinner", "bar", "minimal"
  progress_style: "spinner"
  
  # Use emojis in output (disable for CI/automation)
  emojis: true

# Custom compression profiles (extend built-in ones)
profiles:
  # Example: Ultra-fast profile for development/testing
  fast:
    name: "Fast"
    description: "Minimal compression for development and testing"
    level: 1
    dictionary: "4m"
    fast_bytes: 32
    solid_mode: false
    algorithm: "lzma2"
  
  # Example: Maximum compression profile
  maximum:
    name: "Maximum"
    description: "Slowest but highest compression ratio"
    level: 9
    dictionary: "128m"
    fast_bytes: 273
    solid_mode: true
    algorithm: "lzma2"

# TrueNAS integration (when implemented)
truenas:
  default_host: "truenas-homelab.local"
  upload_path: "/mnt/tank/archives"
  verify_ssl: true
  timeout: 300

# Presets - saved combinations of common flags
presets:
  podcast:
    profile: "media"
    comprehensive: true
    output: "~/Archives"
    exclude: ["__MACOSX", ".DS_Store", "*.pkf"]
  
  backup:
    profile: "balanced"
    comprehensive: true
    force: true
  
  source_code:
    profile: "documents"
    comprehensive: true
    exclude: ["node_modules", ".git", "target", "build"]
`

	// Write config file
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("✅ Created config file: %s\n", configPath)
	fmt.Printf("\nEdit this file to customize your 7zarch-go defaults and presets.\n")
	fmt.Printf("Run '7zarch-go config show' to see your current settings.\n")

	return nil
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	configPath, _ := config.ConfigPath()
	fmt.Printf("Configuration loaded from: %s\n\n", configPath)

	// Convert to YAML and display
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	fmt.Printf("Current Configuration:\n")
	fmt.Printf("=====================\n")
	fmt.Print(string(data))

	return nil
}
