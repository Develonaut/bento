// Package miso provides terminal output "seasoning" - themed styling and progress display.
//
// Theme configuration persistence to ~/.bento/config/theme.
package miso

import (
	"fmt"
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
	// Store config in ~/.bento/config/ subdirectory for consistency with hangiri storage structure
	return filepath.Join(home, ".bento", "config"), nil
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
// Returns VariantTonkotsu (creamy white) as default if no saved theme or on error.
func LoadSavedTheme() Variant {
	path, err := themeConfigPath()
	if err != nil {
		return VariantTonkotsu
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return VariantTonkotsu
	}

	variant := Variant(strings.TrimSpace(string(data)))

	// Validate the variant
	for _, v := range AllVariants() {
		if v == variant {
			return variant
		}
	}

	return VariantTonkotsu
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

// slowMoConfigPath returns the path to the slowMo config file.
func slowMoConfigPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "slowmo"), nil
}

// LoadSlowMoDelay loads the saved slowMo delay from disk.
// Returns 250ms as default if no saved value or on error.
// SlowMo adds artificial delays between node executions to make animations visible.
func LoadSlowMoDelay() int {
	path, err := slowMoConfigPath()
	if err != nil {
		return 250
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return 250
	}

	value := strings.TrimSpace(string(data))

	// Parse as milliseconds
	var ms int
	_, err = fmt.Sscanf(value, "%d", &ms)
	if err != nil || ms < 0 {
		return 250
	}

	return ms
}

// SaveSlowMoDelay saves the slowMo delay (in milliseconds) to disk.
// Creates ~/.bento directory if it doesn't exist.
func SaveSlowMoDelay(ms int) error {
	dir, err := configDir()
	if err != nil {
		return err
	}

	// Create config directory if needed
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	path, err := slowMoConfigPath()
	if err != nil {
		return err
	}

	return os.WriteFile(path, []byte(fmt.Sprintf("%d", ms)), 0644)
}
