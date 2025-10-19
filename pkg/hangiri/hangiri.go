// Package hangiri provides persistent storage for bento workflow definitions.
//
// "Hangiri" (半切り - wooden tub for sushi rice) stores workflow definitions
// on disk as JSON files, allowing bentos to be saved, loaded, and reused.
//
// Storage location: ~/.bento/workflows/
// File format: <name>.json
//
// # Usage
//
//	storage := hangiri.New("~/.bento/workflows")
//
//	// Save a workflow
//	err := storage.Save(ctx, "my-workflow", definition)
//
//	// Load a workflow
//	def, err := storage.Load(ctx, "my-workflow")
//
//	// List all workflows
//	names, err := storage.List(ctx)
//
//	// Delete a workflow
//	err := storage.Delete(ctx, "my-workflow")
//
// Security: Workflow names are validated to prevent directory traversal attacks.
package hangiri

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Develonaut/bento/pkg/neta"
)

// Storage manages persistent storage of bento definitions.
type Storage struct {
	baseDir string
}

// New creates a new Storage instance.
//
// baseDir is the directory where workflows will be stored.
// Typically ~/.bento/workflows/
func New(baseDir string) *Storage {
	return &Storage{baseDir: baseDir}
}

// Save saves a bento definition to disk.
//
// The workflow is saved as <name>.json in the storage directory.
// Returns an error if the name is invalid or if writing fails.
func (s *Storage) Save(ctx context.Context, name string, def *neta.Definition) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if err := validateName(name); err != nil {
		return err
	}

	if err := s.ensureDir(); err != nil {
		return err
	}

	return s.saveToFile(name, def)
}

// saveToFile serializes and writes definition to file.
func (s *Storage) saveToFile(name string, def *neta.Definition) error {
	data, err := s.marshal(def)
	if err != nil {
		return fmt.Errorf("failed to serialize workflow '%s': %w", name, err)
	}

	if err := s.writeFile(name, data); err != nil {
		return fmt.Errorf("failed to write workflow '%s': %w", name, err)
	}

	return nil
}

// Load loads a bento definition from disk.
//
// Returns an error if the workflow doesn't exist or cannot be parsed.
func (s *Storage) Load(ctx context.Context, name string) (*neta.Definition, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if err := validateName(name); err != nil {
		return nil, err
	}

	return s.loadFromFile(name)
}

// loadFromFile reads and deserializes definition from file.
func (s *Storage) loadFromFile(name string) (*neta.Definition, error) {
	data, err := s.readFile(name)
	if err != nil {
		return nil, err
	}

	def, err := s.unmarshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse workflow '%s': %w", name, err)
	}

	return def, nil
}

// List returns all saved workflow names.
func (s *Storage) List(ctx context.Context) ([]string, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if err := s.ensureDir(); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(s.baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to list workflows: %w", err)
	}

	names := s.extractNames(entries)
	return names, nil
}

// Delete removes a workflow from disk.
func (s *Storage) Delete(ctx context.Context, name string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if err := validateName(name); err != nil {
		return err
	}

	path := s.getPath(name)
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("workflow '%s' not found", name)
		}
		return fmt.Errorf("failed to delete workflow '%s': %w", name, err)
	}

	return nil
}

// ensureDir creates the storage directory if it doesn't exist.
func (s *Storage) ensureDir() error {
	return os.MkdirAll(s.baseDir, 0755)
}

// getPath returns the full file path for a workflow name.
func (s *Storage) getPath(name string) string {
	return filepath.Join(s.baseDir, name+".json")
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

// readFile reads a workflow file from disk.
func (s *Storage) readFile(name string) ([]byte, error) {
	path := s.getPath(name)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("workflow '%s' not found", name)
		}
		return nil, fmt.Errorf("failed to read workflow '%s': %w", name, err)
	}
	return data, nil
}

// writeFile writes data to a workflow file.
func (s *Storage) writeFile(name string, data []byte) error {
	path := s.getPath(name)
	return os.WriteFile(path, data, 0644)
}

// extractNames extracts workflow names from directory entries.
func (s *Storage) extractNames(entries []os.DirEntry) []string {
	var names []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			name := entry.Name()[:len(entry.Name())-5]
			names = append(names, name)
		}
	}
	return names
}
