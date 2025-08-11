package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Compression CompressionConfig `yaml:"compression"`
	Defaults    DefaultsConfig    `yaml:"defaults"`
	UI          UIConfig          `yaml:"ui"`
	Profiles    map[string]CustomProfile `yaml:"profiles"`
	TrueNAS     TrueNASConfig     `yaml:"truenas"`
	Presets     map[string]PresetConfig `yaml:"presets"`
	Storage     StorageConfig     `yaml:"storage"`
}

type CompressionConfig struct {
	SmartDefault    bool `yaml:"smart_default"`
	Level           int  `yaml:"level"`
	MediaThreshold  int  `yaml:"media_threshold"`
	DocsThreshold   int  `yaml:"docs_threshold"`
}

type DefaultsConfig struct {
	Create CreateDefaults `yaml:"create"`
	Test   TestDefaults   `yaml:"test"`
}

type CreateDefaults struct {
	Comprehensive bool `yaml:"comprehensive"`
	Force         bool `yaml:"force"`
	Threads       int  `yaml:"threads"`
}

type TestDefaults struct {
	Concurrent int  `yaml:"concurrent"`
	Verbose    bool `yaml:"verbose"`
}

type UIConfig struct {
	ShowAnalysis   bool   `yaml:"show_analysis"`
	ShowTips       bool   `yaml:"show_tips"`
	ProgressStyle  string `yaml:"progress_style"`
	Emojis         bool   `yaml:"emojis"`
}

type CustomProfile struct {
	Name         string `yaml:"name"`
	Description  string `yaml:"description"`
	Level        int    `yaml:"level"`
	Dictionary   string `yaml:"dictionary"`
	FastBytes    int    `yaml:"fast_bytes"`
	SolidMode    bool   `yaml:"solid_mode"`
	Algorithm    string `yaml:"algorithm"`
}

type TrueNASConfig struct {
	DefaultHost string `yaml:"default_host"`
	UploadPath  string `yaml:"upload_path"`
	VerifySSL   bool   `yaml:"verify_ssl"`
	Timeout     int    `yaml:"timeout"`
}

type PresetConfig struct {
	Profile       string   `yaml:"profile"`
	Comprehensive bool     `yaml:"comprehensive"`
	Force         bool     `yaml:"force"`
	Output        string   `yaml:"output"`
	Threads       int      `yaml:"threads"`
	Exclude       []string `yaml:"exclude"`
}

type StorageConfig struct {
	ManagedPath        string `yaml:"managed_path"`
	UseManagedDefault  bool   `yaml:"use_managed_default"`
	AutoOrganize       string `yaml:"auto_organize"` // flat, by_date, by_type
	RetentionDays      int    `yaml:"retention_days"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Compression: CompressionConfig{
			SmartDefault:   true,
			Level:          9,
			MediaThreshold: 70,
			DocsThreshold:  60,
		},
		Defaults: DefaultsConfig{
			Create: CreateDefaults{
				Comprehensive: false,
				Force:         false,
				Threads:       0,
			},
			Test: TestDefaults{
				Concurrent: 5,
				Verbose:    false,
			},
		},
		UI: UIConfig{
			ShowAnalysis:  true,
			ShowTips:      true,
			ProgressStyle: "spinner",
			Emojis:        true,
		},
		Profiles: map[string]CustomProfile{},
		TrueNAS: TrueNASConfig{
			DefaultHost: "truenas-homelab.local",
			UploadPath:  "/mnt/tank/archives",
			VerifySSL:   true,
			Timeout:     300,
		},
		Storage: StorageConfig{
			ManagedPath:       "~/.7zarch-go",
			UseManagedDefault: true,
			AutoOrganize:      "flat",
			RetentionDays:     30,
		},
		Presets: map[string]PresetConfig{
			"podcast": {
				Profile:       "media",
				Comprehensive: true,
				Output:        "~/Archives/Podcasts",
			},
			"backup": {
				Profile:       "balanced",
				Comprehensive: true,
				Force:         true,
			},
			"source_code": {
				Profile:       "documents",
				Comprehensive: true,
				Exclude:       []string{"node_modules", ".git", "target", "build"},
			},
		},
	}
}

// Load loads configuration from ~/.7zarch-go-config
func Load() (*Config, error) {
	config := DefaultConfig()
	
	// Get config file path
	home, err := os.UserHomeDir()
	if err != nil {
		return config, nil // Return defaults if we can't find home
	}
	
	configPath := filepath.Join(home, ".7zarch-go-config")
	
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, nil // Return defaults if config doesn't exist
	}
	
	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, nil // Return defaults on read error
	}
	
	// Parse YAML
	if err := yaml.Unmarshal(data, config); err != nil {
		return config, nil // Return defaults on parse error
	}
	
	return config, nil
}

// ConfigPath returns the path to the config file
func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".7zarch-go-config"), nil
}