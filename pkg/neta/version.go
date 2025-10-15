package neta

import (
	"fmt"
	"strconv"
	"strings"
)

// isCompatibleVersion checks if a version string is compatible
// with the current version. Compatibility is based on major version.
func isCompatibleVersion(v string) bool {
	if v == "" {
		return false
	}

	major, err := parseMajorVersion(v)
	if err != nil {
		return false
	}

	currentMajor, err := parseMajorVersion(CurrentVersion)
	if err != nil {
		return false
	}

	return major == currentMajor
}

// parseMajorVersion extracts the major version number from a version string.
func parseMajorVersion(v string) (int, error) {
	parts := strings.Split(v, ".")
	if len(parts) < 1 {
		return 0, fmt.Errorf("invalid version format")
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid major version: %w", err)
	}

	return major, nil
}

// ValidateVersion checks version and returns descriptive error.
// It ensures the version is present and compatible with CurrentVersion.
func ValidateVersion(v string) error {
	if v == "" {
		return fmt.Errorf("version is required (current version: %s)", CurrentVersion)
	}

	if !isCompatibleVersion(v) {
		return fmt.Errorf("incompatible version %s (current version: %s)", v, CurrentVersion)
	}

	return nil
}
