package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// VSCodeTheme represents a VS Code theme JSON structure
type VSCodeTheme struct {
	Name   string                 `json:"name"`
	Type   string                 `json:"type"`
	Colors map[string]string      `json:"colors"`
	TokenColors []TokenColor       `json:"tokenColors,omitempty"`
}

// TokenColor represents syntax highlighting colors
type TokenColor struct {
	Name     string                 `json:"name,omitempty"`
	Scope    interface{}           `json:"scope,omitempty"` // Can be string or []string
	Settings map[string]string     `json:"settings"`
}

// ExtensionMetadata represents VS Code extension package.json metadata
type ExtensionMetadata struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Publisher   string `json:"publisher"`
	Author      struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"author"`
	Version     string `json:"version"`
	Repository  struct {
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"repository"`
}

// ThemeInfo holds information about a discoverable theme
type ThemeInfo struct {
	Name        string
	DisplayName string
	Path        string
	Type        string // "dark" or "light"
	ExtensionMetadata *ExtensionMetadata // Optional extension metadata
}

// DiscoverVSCodeThemes finds all VS Code themes in the extensions directory
func DiscoverVSCodeThemes() ([]ThemeInfo, error) {
	extensionsPath, err := getVSCodeExtensionsPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get VS Code extensions path: %w", err)
	}
	var themes []ThemeInfo

	// Walk through all extension directories
	err = filepath.WalkDir(extensionsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// Skip directories we can't access
			return nil
		}

		// Look for theme JSON files in themes directories
		if d.IsDir() || !strings.HasSuffix(strings.ToLower(path), ".json") {
			return nil
		}

		// Check if this is in a themes directory (cross-platform)
		if !isThemesDirectory(path) {
			return nil
		}

		// Try to read and parse the theme file
		themeInfo, err := parseThemeFile(path)
		if err != nil {
			// Skip invalid theme files
			return nil
		}

		themes = append(themes, *themeInfo)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk extensions directory: %w", err)
	}

	return themes, nil
}

// parseThemeFile parses a VS Code theme JSON file and extracts basic info
func parseThemeFile(path string) (*ThemeInfo, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var theme VSCodeTheme
	if err := json.Unmarshal(data, &theme); err != nil {
		return nil, err
	}

	// Skip if no name
	if theme.Name == "" {
		return nil, fmt.Errorf("theme has no name")
	}

	// Generate a clean display name from the file path and theme name
	filename := filepath.Base(path)
	extensionName := extractExtensionName(path)
	
	displayName := theme.Name
	if extensionName != "" {
		displayName = fmt.Sprintf("%s (%s)", theme.Name, extensionName)
	}

	return &ThemeInfo{
		Name:        strings.TrimSuffix(filename, ".json"),
		DisplayName: displayName,
		Path:        path,
		Type:        theme.Type,
	}, nil
}

// extractExtensionName extracts the extension name from the path
func extractExtensionName(path string) string {
	// Extract extension name from path like /path/to/.vscode/extensions/author.extension-version/themes/theme.json
	parts := strings.Split(path, string(filepath.Separator))
	for i, part := range parts {
		if part == "extensions" && i+1 < len(parts) {
			// The next part should be the extension directory name
			extensionDir := parts[i+1]
			// Remove version number (everything after the last dash if it looks like a version)
			lastDash := strings.LastIndex(extensionDir, "-")
			if lastDash != -1 {
				possibleVersion := extensionDir[lastDash+1:]
				// Simple check if it looks like a version (starts with number or contains dots)
				if len(possibleVersion) > 0 && (possibleVersion[0] >= '0' && possibleVersion[0] <= '9' || strings.Contains(possibleVersion, ".")) {
					extensionDir = extensionDir[:lastDash]
				}
			}
			
			// Convert author.extension format to readable name
			if strings.Contains(extensionDir, ".") {
				parts := strings.Split(extensionDir, ".")
				if len(parts) == 2 {
					return fmt.Sprintf("%s by %s", strings.Title(strings.ReplaceAll(parts[1], "-", " ")), strings.Title(parts[0]))
				}
			}
			
			return strings.Title(strings.ReplaceAll(extensionDir, "-", " "))
		}
	}
	return ""
}

// LoadVSCodeTheme loads and parses a VS Code theme file
func LoadVSCodeTheme(path string) (*VSCodeTheme, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read theme file: %w", err)
	}

	var theme VSCodeTheme
	if err := json.Unmarshal(data, &theme); err != nil {
		return nil, fmt.Errorf("failed to parse theme JSON: %w", err)
	}

	return &theme, nil
}

// LoadExtensionMetadata loads the package.json metadata for an extension from a theme path
func LoadExtensionMetadata(themePath string) (*ExtensionMetadata, error) {
	// Find the extension directory from the theme path
	// Theme path is like: /path/to/.vscode/extensions/author.extension-version/themes/theme.json
	// We need to go up to the extension directory and look for package.json
	
	parts := strings.Split(themePath, string(filepath.Separator))
	extensionDirIndex := -1
	
	// Find the extensions directory
	for i, part := range parts {
		if part == "extensions" && i+1 < len(parts) {
			extensionDirIndex = i + 1
			break
		}
	}
	
	if extensionDirIndex == -1 {
		return nil, fmt.Errorf("could not find extensions directory in path: %s", themePath)
	}
	
	// Build the path to the extension directory
	extensionDir := strings.Join(parts[:extensionDirIndex+1], string(filepath.Separator))
	packageJSONPath := filepath.Join(extensionDir, "package.json")
	
	// Try to read the package.json file
	data, err := os.ReadFile(packageJSONPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read package.json: %w", err)
	}
	
	var metadata ExtensionMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse package.json: %w", err)
	}
	
	return &metadata, nil
}

