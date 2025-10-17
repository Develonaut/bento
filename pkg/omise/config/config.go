package config

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Config holds user configuration settings
type Config struct {
	SaveDirectory string
	SlowMoDelayMs int // 0 = off, or delay in milliseconds (250, 500, 1000, 2000, 4000, 8000)
}

// Default returns the default configuration
func Default() Config {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{SaveDirectory: ".bento", SlowMoDelayMs: 250}
	}
	return Config{
		SaveDirectory: filepath.Join(home, ".bento"),
		SlowMoDelayMs: 250,
	}
}

// configDir returns the bento config directory
func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".bento"), nil
}

// configPath returns the path to the config file
func configPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config"), nil
}

// Load loads the saved configuration
func Load() Config {
	path, err := configPath()
	if err != nil {
		return Default()
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Default()
	}

	lines := strings.Split(string(data), "\n")
	cfg := Default()

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "save_directory":
			cfg.SaveDirectory = expandHome(value)
		case "slow_mo_delay_ms":
			// Parse integer, default to 0 on error
			if delayMs, err := parseInt(value); err == nil {
				cfg.SlowMoDelayMs = delayMs
			}
		}
	}

	return cfg
}

// Save saves the configuration
func Save(cfg Config) error {
	dir, err := configDir()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	path, err := configPath()
	if err != nil {
		return err
	}

	content := "# Bento Configuration\n"
	content += "save_directory=" + cfg.SaveDirectory + "\n"
	content += "slow_mo_delay_ms=" + formatInt(cfg.SlowMoDelayMs) + "\n"

	return os.WriteFile(path, []byte(content), 0644)
}

// expandHome expands ~ to the user's home directory
func expandHome(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	if path == "~" {
		return home
	}

	return filepath.Join(home, path[2:])
}

// contractHome contracts the user's home directory to ~
func contractHome(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	if strings.HasPrefix(path, home) {
		return "~" + path[len(home):]
	}

	return path
}

// GetSaveDirectory returns the save directory with ~ contraction
func (c Config) GetSaveDirectory() string {
	return contractHome(c.SaveDirectory)
}

// parseInt parses an integer from a string
func parseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// formatInt formats an integer as a string
func formatInt(i int) string {
	return strconv.Itoa(i)
}
