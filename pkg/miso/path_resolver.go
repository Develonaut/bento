package miso

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// ResolvePath expands special markers and environment variables in a path.
// Supports:
//   - {{GDRIVE}} - Google Drive root
//   - {{DROPBOX}} - Dropbox root
//   - {{ONEDRIVE}} - OneDrive root
//   - {{BENTO_HOME}} - Configured bento home
//   - ${VAR} or $VAR - Environment variables
func ResolvePath(path string) (string, error) {
	if path == "" {
		return "", nil
	}

	// Step 1: Expand special markers
	resolved := path
	resolved = expandSpecialMarkers(resolved)

	// Step 2: Expand environment variables
	resolved = os.ExpandEnv(resolved)

	// Step 3: Clean the path
	resolved = filepath.Clean(resolved)

	return resolved, nil
}

// expandSpecialMarkers replaces special markers with platform-specific paths
func expandSpecialMarkers(path string) string {
	// {{BENTO_HOME}}
	if strings.Contains(path, "{{BENTO_HOME}}") {
		bentoHome := LoadBentoHome()
		path = strings.ReplaceAll(path, "{{BENTO_HOME}}", bentoHome)
	}

	// {{GDRIVE}}
	if strings.Contains(path, "{{GDRIVE}}") {
		if gdrivePath := detectGoogleDrive(); gdrivePath != "" {
			path = strings.ReplaceAll(path, "{{GDRIVE}}", gdrivePath)
		}
	}

	// {{DROPBOX}}
	if strings.Contains(path, "{{DROPBOX}}") {
		if dropboxPath := detectDropbox(); dropboxPath != "" {
			path = strings.ReplaceAll(path, "{{DROPBOX}}", dropboxPath)
		}
	}

	// {{ONEDRIVE}}
	if strings.Contains(path, "{{ONEDRIVE}}") {
		if onedrivePath := detectOneDrive(); onedrivePath != "" {
			path = strings.ReplaceAll(path, "{{ONEDRIVE}}", onedrivePath)
		}
	}

	return path
}

// detectGoogleDrive attempts to detect the Google Drive root path
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

// detectDropbox attempts to detect the Dropbox root path
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

// detectOneDrive attempts to detect the OneDrive root path
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

// ResolvePathsInMap resolves all paths in a string map (useful for variables)
func ResolvePathsInMap(m map[string]string) (map[string]string, error) {
	if m == nil {
		return nil, nil
	}

	resolved := make(map[string]string, len(m))
	for key, value := range m {
		resolvedValue, err := ResolvePath(value)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve path for key %s: %w", key, err)
		}
		resolved[key] = resolvedValue
	}

	return resolved, nil
}
