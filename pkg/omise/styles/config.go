package styles

import (
	"os"
	"path/filepath"
	"strings"
)

// configDir returns the bento config directory
func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".bento"), nil
}

// themeConfigPath returns the path to the theme config file
func themeConfigPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "theme"), nil
}

// LoadSavedTheme loads the saved theme variant
func LoadSavedTheme() Variant {
	path, err := themeConfigPath()
	if err != nil {
		return VariantMaguro // Default
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return VariantMaguro // Default if no saved theme
	}

	variant := Variant(strings.TrimSpace(string(data)))
	// Validate the variant
	for _, v := range AllVariants() {
		if v == variant {
			return variant
		}
	}

	return VariantMaguro // Default if invalid variant
}

// SaveTheme saves the current theme variant
func SaveTheme(variant Variant) error {
	dir, err := configDir()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	path, err := themeConfigPath()
	if err != nil {
		return err
	}

	return os.WriteFile(path, []byte(string(variant)), 0644)
}
