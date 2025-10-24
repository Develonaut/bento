package miso

import (
	"fmt"
	"os"
	"path/filepath"
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

// Note: detectGoogleDrive, detectDropbox, and detectOneDrive are defined in config.go
// to avoid circular dependency issues (config.go needs them for LoadBentoHome)

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

// CompressPath converts absolute paths to use special markers for portability.
// This is the inverse of ResolvePath - useful for displaying paths in a platform-independent way.
// Example: "/Users/Ryan/Library/CloudStorage/GoogleDrive-email/My Drive/foo" -> "{{GDRIVE}}/foo"
func CompressPath(path string) string {
	if path == "" {
		return ""
	}

	// Clean the path first
	cleaned := filepath.Clean(path)

	// Try to compress special markers in order of specificity
	// Check GDRIVE first
	if gdrivePath := detectGoogleDrive(); gdrivePath != "" {
		cleanedGDrive := filepath.Clean(gdrivePath)
		if strings.HasPrefix(cleaned, cleanedGDrive) {
			rel := strings.TrimPrefix(cleaned, cleanedGDrive)
			rel = strings.TrimPrefix(rel, string(filepath.Separator))
			if rel == "" {
				return "{{GDRIVE}}"
			}
			return "{{GDRIVE}}" + string(filepath.Separator) + rel
		}
	}

	// Check DROPBOX
	if dropboxPath := detectDropbox(); dropboxPath != "" {
		cleanedDropbox := filepath.Clean(dropboxPath)
		if strings.HasPrefix(cleaned, cleanedDropbox) {
			rel := strings.TrimPrefix(cleaned, cleanedDropbox)
			rel = strings.TrimPrefix(rel, string(filepath.Separator))
			if rel == "" {
				return "{{DROPBOX}}"
			}
			return "{{DROPBOX}}" + string(filepath.Separator) + rel
		}
	}

	// Check ONEDRIVE
	if onedrivePath := detectOneDrive(); onedrivePath != "" {
		cleanedOneDrive := filepath.Clean(onedrivePath)
		if strings.HasPrefix(cleaned, cleanedOneDrive) {
			rel := strings.TrimPrefix(cleaned, cleanedOneDrive)
			rel = strings.TrimPrefix(rel, string(filepath.Separator))
			if rel == "" {
				return "{{ONEDRIVE}}"
			}
			return "{{ONEDRIVE}}" + string(filepath.Separator) + rel
		}
	}

	// Check BENTO_HOME
	bentoHome := LoadBentoHome()
	cleanedBentoHome := filepath.Clean(bentoHome)
	if strings.HasPrefix(cleaned, cleanedBentoHome) {
		rel := strings.TrimPrefix(cleaned, cleanedBentoHome)
		rel = strings.TrimPrefix(rel, string(filepath.Separator))
		if rel == "" {
			return "{{BENTO_HOME}}"
		}
		return "{{BENTO_HOME}}" + string(filepath.Separator) + rel
	}

	// If no marker applies, return the original cleaned path
	return cleaned
}
