// Package miso provides terminal output "seasoning" - themed styling and progress display.
//
// Theme configuration persistence to ~/.bento/config/theme.
package miso

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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
// Returns 5000ms as default if no saved value or on error.
// SlowMo adds artificial delays between node executions to make animations visible.
func LoadSlowMoDelay() int {
	path, err := slowMoConfigPath()
	if err != nil {
		return 5000
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return 5000
	}

	value := strings.TrimSpace(string(data))

	// Parse as milliseconds
	var ms int
	_, err = fmt.Sscanf(value, "%d", &ms)
	if err != nil || ms < 0 {
		return 5000
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

// saveDirConfigPath returns the path to the save directory config file.
func saveDirConfigPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "savedir"), nil
}

// LoadSaveDirectory loads the saved bentos directory from disk.
// Returns ~/.bento as default if no saved value or on error.
func LoadSaveDirectory() string {
	path, err := saveDirConfigPath()
	if err != nil {
		return defaultSaveDir()
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return defaultSaveDir()
	}

	dir := strings.TrimSpace(string(data))
	if dir == "" {
		return defaultSaveDir()
	}

	return dir
}

// SaveSaveDirectory saves the bentos directory to disk.
// Creates ~/.bento directory if it doesn't exist.
func SaveSaveDirectory(dir string) error {
	confDir, err := configDir()
	if err != nil {
		return err
	}

	// Create config directory if needed
	if err := os.MkdirAll(confDir, 0755); err != nil {
		return err
	}

	path, err := saveDirConfigPath()
	if err != nil {
		return err
	}

	return os.WriteFile(path, []byte(dir), 0644)
}

// defaultSaveDir returns the default save directory.
func defaultSaveDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "./.bento"
	}
	return filepath.Join(home, ".bento")
}

// bentoHomeConfigPath returns the path to the bento home config file.
func bentoHomeConfigPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "bentohome"), nil
}

// LoadBentoHome loads the configured bento home directory from disk.
// Returns the default ~/.bento if no custom home is configured.
// Automatically resolves {{GDRIVE}} and other special markers.
func LoadBentoHome() string {
	path, err := bentoHomeConfigPath()
	if err != nil {
		return defaultBentoHome()
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return defaultBentoHome()
	}

	dir := strings.TrimSpace(string(data))
	if dir == "" {
		return defaultBentoHome()
	}

	// Resolve special markers like {{GDRIVE}}
	// Note: Can't call ResolvePath here due to circular dependency
	// (ResolvePath calls LoadBentoHome for {{BENTO_HOME}})
	// So we handle other markers manually
	resolved := dir
	if strings.Contains(resolved, "{{GDRIVE}}") {
		if gdrivePath := detectGoogleDrive(); gdrivePath != "" {
			resolved = strings.ReplaceAll(resolved, "{{GDRIVE}}", gdrivePath)
		}
	}
	if strings.Contains(resolved, "{{DROPBOX}}") {
		if dropboxPath := detectDropbox(); dropboxPath != "" {
			resolved = strings.ReplaceAll(resolved, "{{DROPBOX}}", dropboxPath)
		}
	}
	if strings.Contains(resolved, "{{ONEDRIVE}}") {
		if onedrivePath := detectOneDrive(); onedrivePath != "" {
			resolved = strings.ReplaceAll(resolved, "{{ONEDRIVE}}", onedrivePath)
		}
	}

	return filepath.Clean(resolved)
}

// SaveBentoHome saves the bento home directory to disk.
// Creates ~/.bento directory if it doesn't exist.
func SaveBentoHome(dir string) error {
	confDir, err := configDir()
	if err != nil {
		return err
	}

	// Create config directory if needed
	if err := os.MkdirAll(confDir, 0755); err != nil {
		return err
	}

	path, err := bentoHomeConfigPath()
	if err != nil {
		return err
	}

	return os.WriteFile(path, []byte(dir), 0644)
}

// defaultBentoHome returns the default bento home directory.
func defaultBentoHome() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "./.bento"
	}
	return filepath.Join(home, ".bento")
}

// detectGoogleDrive attempts to detect the Google Drive root path.
// Duplicated from path_resolver.go to avoid circular dependency.
func detectGoogleDrive() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	switch runtime.GOOS {
	case "darwin":
		// Mac: ~/Library/CloudStorage/GoogleDrive-{email}/My Drive
		cloudStorageDir := filepath.Join(homeDir, "Library", "CloudStorage")
		if entries, err := os.ReadDir(cloudStorageDir); err == nil {
			for _, entry := range entries {
				if entry.IsDir() && strings.HasPrefix(entry.Name(), "GoogleDrive-") {
					myDrive := filepath.Join(cloudStorageDir, entry.Name(), "My Drive")
					if stat, err := os.Stat(myDrive); err == nil && stat.IsDir() {
						return myDrive
					}
				}
			}
		}

	case "windows":
		// Windows: Check common drive letters for "My Drive"
		driveLetters := []string{"G", "H", "I", "J", "K", "L", "M", "N", "O"}
		for _, letter := range driveLetters {
			myDrive := filepath.Join(letter+":", "My Drive")
			if stat, err := os.Stat(myDrive); err == nil && stat.IsDir() {
				return myDrive
			}
		}

		// Also check %USERPROFILE%\Google Drive
		googleDrive := filepath.Join(homeDir, "Google Drive")
		if stat, err := os.Stat(googleDrive); err == nil && stat.IsDir() {
			return googleDrive
		}

	case "linux":
		// Linux: ~/Google Drive (older Drive File Stream)
		googleDrive := filepath.Join(homeDir, "Google Drive")
		if stat, err := os.Stat(googleDrive); err == nil && stat.IsDir() {
			return googleDrive
		}
	}

	return ""
}

// detectDropbox attempts to detect the Dropbox root path.
// Duplicated from path_resolver.go to avoid circular dependency.
func detectDropbox() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	dropboxPath := filepath.Join(homeDir, "Dropbox")
	if stat, err := os.Stat(dropboxPath); err == nil && stat.IsDir() {
		return dropboxPath
	}

	return ""
}

// detectOneDrive attempts to detect the OneDrive root path.
// Duplicated from path_resolver.go to avoid circular dependency.
func detectOneDrive() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	switch runtime.GOOS {
	case "darwin":
		// Mac: ~/Library/CloudStorage/OneDrive-{org}
		cloudStorageDir := filepath.Join(homeDir, "Library", "CloudStorage")
		if entries, err := os.ReadDir(cloudStorageDir); err == nil {
			for _, entry := range entries {
				if entry.IsDir() && strings.HasPrefix(entry.Name(), "OneDrive") {
					oneDrivePath := filepath.Join(cloudStorageDir, entry.Name())
					return oneDrivePath
				}
			}
		}

	case "windows":
		// Windows: %USERPROFILE%\OneDrive
		oneDrivePath := filepath.Join(homeDir, "OneDrive")
		if stat, err := os.Stat(oneDrivePath); err == nil && stat.IsDir() {
			return oneDrivePath
		}

	case "linux":
		// Linux: ~/OneDrive
		oneDrivePath := filepath.Join(homeDir, "OneDrive")
		if stat, err := os.Stat(oneDrivePath); err == nil && stat.IsDir() {
			return oneDrivePath
		}
	}

	return ""
}
