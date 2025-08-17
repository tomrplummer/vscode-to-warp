package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// getPlatformInfo returns information about the current platform
func getPlatformInfo() (string, string) {
	return runtime.GOOS, runtime.GOARCH
}

// getVSCodeExtensionsPath returns the VS Code extensions directory path for the current platform
func getVSCodeExtensionsPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	switch runtime.GOOS {
	case "windows":
		// Windows: %USERPROFILE%\.vscode\extensions
		return filepath.Join(homeDir, ".vscode", "extensions"), nil
	case "darwin":
		// macOS: ~/.vscode/extensions
		return filepath.Join(homeDir, ".vscode", "extensions"), nil
	case "linux":
		// Linux: ~/.vscode/extensions
		return filepath.Join(homeDir, ".vscode", "extensions"), nil
	default:
		// Try the standard Unix path for unknown systems
		return filepath.Join(homeDir, ".vscode", "extensions"), nil
	}
}

// getWarpThemesPath returns the Warp themes directory path for the current platform
func getWarpThemesPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	switch runtime.GOOS {
	case "windows":
		// Warp is not available on Windows
		return "", fmt.Errorf("Warp terminal is not available on Windows. This tool only works on macOS and Linux")
	case "darwin":
		// macOS: ~/.warp/themes
		return filepath.Join(homeDir, ".warp", "themes"), nil
	case "linux":
		// Linux: ~/.warp/themes  
		return filepath.Join(homeDir, ".warp", "themes"), nil
	default:
		// For unknown Unix-like systems, try the standard path
		return filepath.Join(homeDir, ".warp", "themes"), nil
	}
}

// isThemesDirectory checks if a path contains a themes directory (cross-platform)
func isThemesDirectory(path string) bool {
	// Use filepath.Separator to handle both / and \ separators
	pathSeparator := string(filepath.Separator)
	themesDir := pathSeparator + "themes" + pathSeparator
	
	// Also check for themes at the end of path (edge case)
	return contains(path, themesDir) || hasThemesAtEnd(path)
}

// contains checks if a string contains a substring (case-insensitive on Windows)
func contains(str, substr string) bool {
	if runtime.GOOS == "windows" {
		// Case-insensitive on Windows
		return containsCaseInsensitive(str, substr)
	}
	return containsString(str, substr)
}

// containsString is a simple string contains check
func containsString(str, substr string) bool {
	return len(str) >= len(substr) && findSubstring(str, substr)
}

// findSubstring finds a substring in a string
func findSubstring(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// containsCaseInsensitive performs case-insensitive substring search
func containsCaseInsensitive(str, substr string) bool {
	strLower := toLower(str)
	substrLower := toLower(substr)
	return containsString(strLower, substrLower)
}

// toLower converts string to lowercase
func toLower(s string) string {
	result := make([]rune, len(s))
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			result[i] = r + ('a' - 'A')
		} else {
			result[i] = r
		}
	}
	return string(result)
}

// hasThemesAtEnd checks if path ends with themes directory
func hasThemesAtEnd(path string) bool {
	pathSeparator := string(filepath.Separator)
	return hasSuffix(path, pathSeparator+"themes") || hasSuffix(path, pathSeparator+"themes"+pathSeparator)
}

// hasSuffix checks if string has a suffix
func hasSuffix(str, suffix string) bool {
	return len(str) >= len(suffix) && str[len(str)-len(suffix):] == suffix
}

// validatePlatformSupport checks if the current platform is supported
func validatePlatformSupport() error {
	if runtime.GOOS == "windows" {
		return fmt.Errorf(`Warp terminal is not available on Windows.

This tool is designed to convert VS Code themes to Warp themes.
Since Warp only runs on macOS and Linux, this tool cannot be used on Windows.

However, VS Code theme discovery will still work if you want to examine
your installed themes for other purposes.`)
	}
	return nil
}

// getAlternativeTerminalInfo provides information about alternative terminals on Windows
func getAlternativeTerminalInfo() string {
	return `Alternative terminals that support custom themes on Windows:
- Windows Terminal (JSON themes)
- Alacritty (YAML/TOML themes)
- Hyper (JavaScript themes)
- ConEmu (XML themes)`
}
