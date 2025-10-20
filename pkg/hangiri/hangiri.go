// Package hangiri provides persistent storage for bento data.
//
// "Hangiri" (半切り - wooden tub for sushi rice) stores bento-related data
// on disk as JSON files, allowing bentos, secrets, and other data to be saved,
// loaded, and reused.
//
// Storage structure:
//
//	~/.bento/
//	  bentos/     - User-created workflow definitions
//	  secrets/    - API keys, credentials, etc. (managed by wasabi)
//	  templates/  - Reusable workflow templates
//	  config/     - Configuration files (themes, preferences)
//	  cache/      - Temporary/cached data
//
// File format: <name>.bento.json (for bentos), <name>.json (for others)
//
// # Usage
//
//	// Create a storage instance
//	storage := hangiri.NewDefaultStorage()
//
//	// Save a bento
//	err := storage.SaveBento(ctx, "my-workflow", definition)
//
//	// Load a bento by name
//	def, err := storage.LoadBento(ctx, "my-workflow")
//
//	// List all bentos
//	names, err := storage.ListBentos(ctx)
//
//	// Delete a bento
//	err := storage.DeleteBento(ctx, "my-workflow")
//
// Security: All names are validated to prevent directory traversal attacks.
package hangiri

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Develonaut/bento/pkg/neta"
)

// StorageType represents different types of storage subdirectories.
type StorageType string

const (
	StorageTypeBentos    StorageType = "bentos"
	StorageTypeSecrets   StorageType = "secrets"
	StorageTypeTemplates StorageType = "templates"
	StorageTypeConfig    StorageType = "config"
	StorageTypeCache     StorageType = "cache"
)

// Storage manages persistent storage of bento-related data.
type Storage struct {
	baseDir string
}

// New creates a new Storage instance with a custom base directory.
//
// baseDir is the root directory (typically ~/.bento/)
func New(baseDir string) *Storage {
	return &Storage{baseDir: expandHome(baseDir)}
}

// NewDefaultStorage creates a Storage instance using the default ~/.bento/ directory.
func NewDefaultStorage() *Storage {
	return New("~/.bento")
}

// expandHome expands ~ to the user's home directory.
func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}

// getStorageDir returns the full path to a storage subdirectory.
func (s *Storage) getStorageDir(storageType StorageType) string {
	return filepath.Join(s.baseDir, string(storageType))
}

// ensureStorageDir creates a storage subdirectory if it doesn't exist.
func (s *Storage) ensureStorageDir(storageType StorageType) error {
	dir := s.getStorageDir(storageType)
	return os.MkdirAll(dir, 0755)
}

// getBentoPath returns the full file path for a bento name.
func (s *Storage) getBentoPath(name string) string {
	// Strip .bento.json extension if present
	name = strings.TrimSuffix(name, ".bento.json")
	return filepath.Join(s.getStorageDir(StorageTypeBentos), name+".bento.json")
}

// getGenericPath returns the full file path for a generic JSON file.
// Reserved for future storage types (templates, cache, etc.).
// Currently unused but kept for consistency with storage architecture.
func (s *Storage) getGenericPath(storageType StorageType, name string) string {
	// Strip .json extension if present
	name = strings.TrimSuffix(name, ".json")
	return filepath.Join(s.getStorageDir(storageType), name+".json")
}

// SaveBento saves a bento definition to ~/.bento/bentos/
//
// The bento is saved as <name>.bento.json in the bentos directory.
// Returns an error if the name is invalid or if writing fails.
func (s *Storage) SaveBento(ctx context.Context, name string, def *neta.Definition) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if err := validateName(name); err != nil {
		return err
	}

	if err := s.ensureStorageDir(StorageTypeBentos); err != nil {
		return err
	}

	data, err := s.marshal(def)
	if err != nil {
		return fmt.Errorf("failed to serialize bento '%s': %w", name, err)
	}

	path := s.getBentoPath(name)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write bento '%s': %w", name, err)
	}

	return nil
}

// LoadBento loads a bento definition from ~/.bento/bentos/
//
// Returns an error if the bento doesn't exist or cannot be parsed.
func (s *Storage) LoadBento(ctx context.Context, name string) (*neta.Definition, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if err := validateName(name); err != nil {
		return nil, err
	}

	path := s.getBentoPath(name)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("bento '%s' not found", name)
		}
		return nil, fmt.Errorf("failed to read bento '%s': %w", name, err)
	}

	def, err := s.unmarshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bento '%s': %w", name, err)
	}

	return def, nil
}

// ListBentos returns all saved bento names from ~/.bento/bentos/
func (s *Storage) ListBentos(ctx context.Context) ([]string, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if err := s.ensureStorageDir(StorageTypeBentos); err != nil {
		return nil, err
	}

	dir := s.getStorageDir(StorageTypeBentos)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to list bentos: %w", err)
	}

	var names []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".bento.json") {
			name := strings.TrimSuffix(entry.Name(), ".bento.json")
			names = append(names, name)
		}
	}
	return names, nil
}

// DeleteBento removes a bento from ~/.bento/bentos/
func (s *Storage) DeleteBento(ctx context.Context, name string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if err := validateName(name); err != nil {
		return err
	}

	path := s.getBentoPath(name)
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("bento '%s' not found", name)
		}
		return fmt.Errorf("failed to delete bento '%s': %w", name, err)
	}

	return nil
}

// BentoExists checks if a bento exists in storage.
func (s *Storage) BentoExists(ctx context.Context, name string) bool {
	if err := validateName(name); err != nil {
		return false
	}
	path := s.getBentoPath(name)
	_, err := os.Stat(path)
	return err == nil
}

// Legacy methods for backward compatibility
// These maintain the old API but delegate to the new bento-specific methods.

// Save saves a bento definition using the legacy API.
// Deprecated: Use SaveBento instead.
func (s *Storage) Save(ctx context.Context, name string, def *neta.Definition) error {
	return s.SaveBento(ctx, name, def)
}

// Load loads a bento definition using the legacy API.
// Deprecated: Use LoadBento instead.
func (s *Storage) Load(ctx context.Context, name string) (*neta.Definition, error) {
	return s.LoadBento(ctx, name)
}

// List returns all saved bento names using the legacy API.
// Deprecated: Use ListBentos instead.
func (s *Storage) List(ctx context.Context) ([]string, error) {
	return s.ListBentos(ctx)
}

// Delete removes a bento using the legacy API.
// Deprecated: Use DeleteBento instead.
func (s *Storage) Delete(ctx context.Context, name string) error {
	return s.DeleteBento(ctx, name)
}

// marshal serializes a definition to JSON with indentation.
func (s *Storage) marshal(def *neta.Definition) ([]byte, error) {
	return json.MarshalIndent(def, "", "  ")
}

// unmarshal deserializes JSON data to a definition.
func (s *Storage) unmarshal(data []byte) (*neta.Definition, error) {
	var def neta.Definition
	if err := json.Unmarshal(data, &def); err != nil {
		return nil, err
	}
	return &def, nil
}
