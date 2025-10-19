// Package miso provides terminal output "seasoning" - themed styling and progress display.
//
// Theme configuration persistence to ~/.bento/theme.
package miso

import (
	"os"
	"path/filepath"
	"strings"
)

// configDir returns the bento config directory path.
// Mutable var allows mocking in tests.
var configDir = func() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".bento"), nil
}

// themeConfigPath returns the path to the theme config file.
func themeConfigPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "theme"), nil
}

// LoadSavedTheme loads the saved theme variant from disk.
// Returns VariantMaguro (red) as default if no saved theme or on error.
func LoadSavedTheme() Variant {
	path, err := themeConfigPath()
	if err != nil {
		return VariantMaguro
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return VariantMaguro
	}

	variant := Variant(strings.TrimSpace(string(data)))

	// Validate the variant
	for _, v := range AllVariants() {
		if v == variant {
			return variant
		}
	}

	return VariantMaguro
}

// SaveTheme saves the theme variant to disk.
// Creates ~/.bento directory if it doesn't exist.
func SaveTheme(variant Variant) error {
	dir, err := configDir()
	if err != nil {
		return err
	}

	// Create config directory if needed
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	path, err := themeConfigPath()
	if err != nil {
		return err
	}

	return os.WriteFile(path, []byte(string(variant)), 0644)
}
