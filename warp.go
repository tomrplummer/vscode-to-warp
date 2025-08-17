package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// WarpTheme represents a Warp theme YAML structure
type WarpTheme struct {
	Accent     string            `yaml:"accent"`
	Background string            `yaml:"background"`
	Details    string            `yaml:"details"`
	Foreground string            `yaml:"foreground"`
	TerminalColors TerminalColors `yaml:"terminal_colors"`
}

// TerminalColors represents the terminal color palette
type TerminalColors struct {
	Normal ColorPalette `yaml:"normal"`
	Bright ColorPalette `yaml:"bright"`
}

// ColorPalette represents a set of 8 colors
type ColorPalette struct {
	Black   string `yaml:"black"`
	Red     string `yaml:"red"`
	Green   string `yaml:"green"`
	Yellow  string `yaml:"yellow"`
	Blue    string `yaml:"blue"`
	Magenta string `yaml:"magenta"`
	Cyan    string `yaml:"cyan"`
	White   string `yaml:"white"`
}

// ConvertVSCodeToWarp converts a VS Code theme to Warp theme format
func ConvertVSCodeToWarp(vscodeTheme *VSCodeTheme) (*WarpTheme, error) {
	warpTheme := &WarpTheme{}

	// Set basic properties
	warpTheme.Background = getColorOrDefault(vscodeTheme.Colors, "editor.background", "#1e1e1e")
	warpTheme.Foreground = getColorOrDefault(vscodeTheme.Colors, "editor.foreground", "#d4d4d4")
	
	// Try to find an accent color from various VS Code color keys
	accentKeys := []string{
		"focusBorder",
		"button.background",
		"progressBar.background",
		"textLink.foreground",
		"editorCursor.foreground",
		"terminal.ansiBlue",
	}
	
	warpTheme.Accent = findFirstColor(vscodeTheme.Colors, accentKeys, "#007acc")
	
	// Determine if theme is dark or light and set details accordingly
	if vscodeTheme.Type == "light" {
		warpTheme.Details = "lighter"
	} else {
		warpTheme.Details = "darker"
	}

	// Convert terminal colors
	warpTheme.TerminalColors = convertTerminalColors(vscodeTheme.Colors)

	return warpTheme, nil
}

// convertTerminalColors maps VS Code terminal colors to Warp format
func convertTerminalColors(colors map[string]string) TerminalColors {
	return TerminalColors{
		Normal: ColorPalette{
			Black:   getColorOrDefault(colors, "terminal.ansiBlack", "#1e1e1e"),
			Red:     getColorOrDefault(colors, "terminal.ansiRed", "#f44747"),
			Green:   getColorOrDefault(colors, "terminal.ansiGreen", "#6a9955"),
			Yellow:  getColorOrDefault(colors, "terminal.ansiYellow", "#dcdcaa"),
			Blue:    getColorOrDefault(colors, "terminal.ansiBlue", "#569cd6"),
			Magenta: getColorOrDefault(colors, "terminal.ansiMagenta", "#c586c0"),
			Cyan:    getColorOrDefault(colors, "terminal.ansiCyan", "#9cdcfe"),
			White:   getColorOrDefault(colors, "terminal.ansiWhite", "#d4d4d4"),
		},
		Bright: ColorPalette{
			Black:   getColorOrDefault(colors, "terminal.ansiBrightBlack", "#686868"),
			Red:     getColorOrDefault(colors, "terminal.ansiBrightRed", "#f44747"),
			Green:   getColorOrDefault(colors, "terminal.ansiBrightGreen", "#6a9955"),
			Yellow:  getColorOrDefault(colors, "terminal.ansiBrightYellow", "#dcdcaa"),
			Blue:    getColorOrDefault(colors, "terminal.ansiBrightBlue", "#569cd6"),
			Magenta: getColorOrDefault(colors, "terminal.ansiBrightMagenta", "#c586c0"),
			Cyan:    getColorOrDefault(colors, "terminal.ansiBrightCyan", "#9cdcfe"),
			White:   getColorOrDefault(colors, "terminal.ansiBrightWhite", "#ffffff"),
		},
	}
}

// getColorOrDefault returns a color from the map or a default value
func getColorOrDefault(colors map[string]string, key, defaultValue string) string {
	if color, exists := colors[key]; exists && color != "" {
		return cleanColor(color)
	}
	return defaultValue
}

// findFirstColor finds the first available color from a list of keys
func findFirstColor(colors map[string]string, keys []string, defaultValue string) string {
	for _, key := range keys {
		if color, exists := colors[key]; exists && color != "" {
			return cleanColor(color)
		}
	}
	return defaultValue
}

// cleanColor removes any alpha channel or invalid characters from hex colors
func cleanColor(color string) string {
	// Remove any whitespace
	color = strings.TrimSpace(color)
	
	// Ensure it starts with #
	if !strings.HasPrefix(color, "#") {
		return color
	}
	
	// Remove alpha channel if present (8-character hex codes)
	if len(color) == 9 {
		return color[:7]
	}
	
	// Handle short hex codes
	if len(color) == 4 {
		// Convert #abc to #aabbcc
		return fmt.Sprintf("#%c%c%c%c%c%c", color[1], color[1], color[2], color[2], color[3], color[3])
	}
	
	return color
}

// SaveWarpTheme saves a Warp theme to the appropriate directory
func SaveWarpTheme(theme *WarpTheme, name string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	themesDir := filepath.Join(homeDir, ".warp", "themes")
	
	// Create themes directory if it doesn't exist
	if err := os.MkdirAll(themesDir, 0755); err != nil {
		return fmt.Errorf("failed to create themes directory: %w", err)
	}

	// Clean the name for filename
	filename := cleanFilename(name) + ".yaml"
	themePath := filepath.Join(themesDir, filename)

	// Marshal to YAML
	yamlData, err := yaml.Marshal(theme)
	if err != nil {
		return fmt.Errorf("failed to marshal theme to YAML: %w", err)
	}

	// Write to file
	if err := os.WriteFile(themePath, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write theme file: %w", err)
	}

	return nil
}

// cleanFilename removes or replaces characters that aren't suitable for filenames
func cleanFilename(name string) string {
	// Replace spaces and special characters with underscores
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "\\", "_")
	name = strings.ReplaceAll(name, ":", "_")
	name = strings.ReplaceAll(name, "*", "_")
	name = strings.ReplaceAll(name, "?", "_")
	name = strings.ReplaceAll(name, "\"", "_")
	name = strings.ReplaceAll(name, "<", "_")
	name = strings.ReplaceAll(name, ">", "_")
	name = strings.ReplaceAll(name, "|", "_")
	
	// Remove parentheses
	name = strings.ReplaceAll(name, "(", "")
	name = strings.ReplaceAll(name, ")", "")
	
	// Convert to lowercase
	name = strings.ToLower(name)
	
	// Remove multiple consecutive underscores
	for strings.Contains(name, "__") {
		name = strings.ReplaceAll(name, "__", "_")
	}
	
	// Remove leading/trailing underscores
	name = strings.Trim(name, "_")
	
	return name
}
