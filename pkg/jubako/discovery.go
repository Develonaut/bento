package jubako

import (
	"os"
	"path/filepath"
	"strings"
)

// Discovery finds workflow files.
type Discovery struct {
	searchPaths []string
}

// NewDiscovery creates a new discovery instance.
func NewDiscovery(paths ...string) *Discovery {
	if len(paths) == 0 {
		paths = defaultPaths()
	}

	return &Discovery{
		searchPaths: paths,
	}
}

// Find searches for .bento.yaml files.
func (d *Discovery) Find() ([]string, error) {
	found := []string{}

	for _, path := range d.searchPaths {
		files, err := findInPath(path)
		if err != nil {
			continue // Skip inaccessible paths
		}
		found = append(found, files...)
	}

	return found, nil
}

// Watch monitors directories for changes (optional, Phase 5+).
func (d *Discovery) Watch() (<-chan string, error) {
	// Future: fsnotify integration
	return nil, nil
}

// findInPath searches a single path for .bento.yaml files.
func findInPath(root string) ([]string, error) {
	found := []string{}
	err := filepath.Walk(root, walkFunc(&found))
	return found, err
}

// walkFunc returns a filepath.WalkFunc that collects .bento.yaml files.
func walkFunc(found *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if isBentoFile(path) {
			*found = append(*found, path)
		}
		return nil
	}
}

// isBentoFile checks if a file is a .bento.yaml file.
func isBentoFile(path string) bool {
	base := filepath.Base(path)
	return strings.HasSuffix(base, ".bento.yaml")
}

// defaultPaths returns default search paths.
func defaultPaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return []string{"."}
	}

	return []string{
		".",
		filepath.Join(home, ".bento"),
		filepath.Join(home, "bento"),
	}
}
