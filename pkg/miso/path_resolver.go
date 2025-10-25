package miso

import (
	"github.com/Develonaut/bento/pkg/kombu"
)

// ResolvePath expands special markers and environment variables in a path.
// This is a convenience wrapper around kombu.ResolvePath.
func ResolvePath(path string) (string, error) {
	return kombu.ResolvePath(path)
}

// CompressPath converts absolute paths to use special markers for portability.
// This is a convenience wrapper around kombu.CompressPath.
func CompressPath(path string) string {
	return kombu.CompressPath(path)
}

// ResolvePathsInMap resolves all paths in a string map (useful for variables).
// This is a convenience wrapper around kombu.ResolvePathsInMap.
func ResolvePathsInMap(m map[string]string) (map[string]string, error) {
	return kombu.ResolvePathsInMap(m)
}
