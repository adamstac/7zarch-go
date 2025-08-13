package tui

import "github.com/charmbracelet/lipgloss"

// Theme defines the color scheme for the TUI
type Theme struct {
	Name        string
	Header      lipgloss.Color // Title/app name
	Foreground  lipgloss.Color // Normal text
	Selection   lipgloss.Color // Selected item background
	SelText     lipgloss.Color // Selected item text
	Metadata    lipgloss.Color // Size/date info
	StatusOK    lipgloss.Color // ✓ Present archives
	StatusMiss  lipgloss.Color // ? Missing archives
	StatusDel   lipgloss.Color // X Deleted archives
	Commands    lipgloss.Color // Command help text
}

// GetTheme returns a theme by name
func GetTheme(name string) Theme {
	themes := map[string]Theme{
		"blue":           BlueTechTheme(),
		"green":          TerminalGreenTheme(),
		"purple":         PurpleGradientTheme(),
		"cyan":           NeonCyanTheme(),
		"charmbracelet":  CharmbraceletTheme(),
		"dracula":        DraculaClassicTheme(),
		"dracula-warm":   DraculaWarmTheme(),
		"dracula-cool":   DraculaCoolTheme(),
		"dracula-minimal": DraculaMinimalTheme(),
	}
	
	if theme, exists := themes[name]; exists {
		return theme
	}
	return DraculaClassicTheme() // Default to Dracula
}

// GetAllThemes returns all available themes
func GetAllThemes() []Theme {
	return []Theme{
		BlueTechTheme(),
		TerminalGreenTheme(),
		PurpleGradientTheme(),
		NeonCyanTheme(),
		CharmbraceletTheme(),
		DraculaClassicTheme(),
		DraculaWarmTheme(),
		DraculaCoolTheme(),
		DraculaMinimalTheme(),
	}
}

// Blue Tech Theme - Professional and clean
func BlueTechTheme() Theme {
	return Theme{
		Name:       "Blue Tech",
		Header:     lipgloss.Color("#00BFFF"),
		Foreground: lipgloss.Color("#FFFFFF"),
		Selection:  lipgloss.Color("#1E3A8A"),
		SelText:    lipgloss.Color("#FFFFFF"),
		Metadata:   lipgloss.Color("#22D3EE"),
		StatusOK:   lipgloss.Color("#10B981"),
		StatusMiss: lipgloss.Color("#F59E0B"),
		StatusDel:  lipgloss.Color("#EF4444"),
		Commands:   lipgloss.Color("#3B82F6"),
	}
}

// Terminal Green Theme - Classic terminal feel
func TerminalGreenTheme() Theme {
	return Theme{
		Name:       "Terminal Green",
		Header:     lipgloss.Color("#00FF00"),
		Foreground: lipgloss.Color("#FFFFFF"),
		Selection:  lipgloss.Color("#065F46"),
		SelText:    lipgloss.Color("#00FF00"),
		Metadata:   lipgloss.Color("#86EFAC"),
		StatusOK:   lipgloss.Color("#00FF00"),
		StatusMiss: lipgloss.Color("#FBBF24"),
		StatusDel:  lipgloss.Color("#FF5555"),
		Commands:   lipgloss.Color("#22C55E"),
	}
}

// Purple Gradient Theme - Modern and vibrant
func PurpleGradientTheme() Theme {
	return Theme{
		Name:       "Purple Gradient",
		Header:     lipgloss.Color("#FF00FF"),
		Foreground: lipgloss.Color("#FFFFFF"),
		Selection:  lipgloss.Color("#7C3AED"),
		SelText:    lipgloss.Color("#F472B6"),
		Metadata:   lipgloss.Color("#C4B5FD"),
		StatusOK:   lipgloss.Color("#EC4899"),
		StatusMiss: lipgloss.Color("#FB923C"),
		StatusDel:  lipgloss.Color("#EF4444"),
		Commands:   lipgloss.Color("#8B5CF6"),
	}
}

// Neon Cyan Theme - High contrast and energetic
func NeonCyanTheme() Theme {
	return Theme{
		Name:       "Neon Cyan",
		Header:     lipgloss.Color("#00FFFF"),
		Foreground: lipgloss.Color("#FFFFFF"),
		Selection:  lipgloss.Color("#0F766E"),
		SelText:    lipgloss.Color("#00FFFF"),
		Metadata:   lipgloss.Color("#FB923C"),
		StatusOK:   lipgloss.Color("#10B981"),
		StatusMiss: lipgloss.Color("#F97316"),
		StatusDel:  lipgloss.Color("#EF4444"),
		Commands:   lipgloss.Color("#06B6D4"),
	}
}

// Charmbracelet Theme - Playful pink theme
func CharmbraceletTheme() Theme {
	return Theme{
		Name:       "Charmbracelet",
		Header:     lipgloss.Color("#FF1493"),
		Foreground: lipgloss.Color("#FFFFFF"),
		Selection:  lipgloss.Color("#EC4899"),
		SelText:    lipgloss.Color("#FFFFFF"),
		Metadata:   lipgloss.Color("#FBCFE8"),
		StatusOK:   lipgloss.Color("#32CD32"),
		StatusMiss: lipgloss.Color("#FFD700"),
		StatusDel:  lipgloss.Color("#FF6347"),
		Commands:   lipgloss.Color("#EC4899"),
	}
}

// Dracula Classic Theme - Enhanced with vibrant selection
func DraculaClassicTheme() Theme {
	return Theme{
		Name:       "Dracula",
		Header:     lipgloss.Color("#BD93F9"),
		Foreground: lipgloss.Color("#F8F8F2"),
		Selection:  lipgloss.Color("#BD93F9"), // Vibrant purple background
		SelText:    lipgloss.Color("#282A36"), // Dark background text for contrast
		Metadata:   lipgloss.Color("#8BE9FD"),
		StatusOK:   lipgloss.Color("#50FA7B"),
		StatusMiss: lipgloss.Color("#F1FA8C"),
		StatusDel:  lipgloss.Color("#FF5555"),
		Commands:   lipgloss.Color("#FF79C6"),
	}
}

// Dracula Warm Theme - Orange accents with vibrant selection
func DraculaWarmTheme() Theme {
	return Theme{
		Name:       "Dracula Warm",
		Header:     lipgloss.Color("#BD93F9"),
		Foreground: lipgloss.Color("#F8F8F2"),
		Selection:  lipgloss.Color("#FFB86C"), // Vibrant orange background
		SelText:    lipgloss.Color("#282A36"), // Dark text for contrast
		Metadata:   lipgloss.Color("#8BE9FD"),
		StatusOK:   lipgloss.Color("#50FA7B"),
		StatusMiss: lipgloss.Color("#F1FA8C"),
		StatusDel:  lipgloss.Color("#FF5555"),
		Commands:   lipgloss.Color("#FF79C6"),
	}
}

// Dracula Cool Theme - Pink/Purple focus with vibrant selection
func DraculaCoolTheme() Theme {
	return Theme{
		Name:       "Dracula Cool",
		Header:     lipgloss.Color("#FF79C6"),
		Foreground: lipgloss.Color("#F8F8F2"),
		Selection:  lipgloss.Color("#FF79C6"), // Vibrant pink background  
		SelText:    lipgloss.Color("#282A36"), // Dark text for contrast
		Metadata:   lipgloss.Color("#BD93F9"),
		StatusOK:   lipgloss.Color("#50FA7B"),
		StatusMiss: lipgloss.Color("#F1FA8C"),
		StatusDel:  lipgloss.Color("#FF5555"),
		Commands:   lipgloss.Color("#BD93F9"),
	}
}

// Dracula Minimal Theme - Enhanced selection with subtle vibrance
func DraculaMinimalTheme() Theme {
	return Theme{
		Name:       "Dracula Minimal",
		Header:     lipgloss.Color("#BD93F9"),
		Foreground: lipgloss.Color("#F8F8F2"),
		Selection:  lipgloss.Color("#6272A4"), // Subtle blue-gray background (brighter than original)
		SelText:    lipgloss.Color("#F8F8F2"), // Bright foreground for readability
		Metadata:   lipgloss.Color("#6272A4"),
		StatusOK:   lipgloss.Color("#50FA7B"),
		StatusMiss: lipgloss.Color("#F1FA8C"),
		StatusDel:  lipgloss.Color("#FF5555"),
		Commands:   lipgloss.Color("#8BE9FD"), // Brighter cyan for commands
	}
}
