package miso

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// VariablesManager handles storage and retrieval of user-defined variables.
// Unlike secrets (which use OS keychain), variables are stored in JSON
// for easy editing and non-sensitive configuration values.
type VariablesManager struct {
	filePath  string
	variables map[string]string
}

// NewVariablesManager creates a new variables manager.
func NewVariablesManager() (*VariablesManager, error) {
	bentoDir := LoadBentoHome()
	filePath := filepath.Join(bentoDir, "variables.json")

	// Ensure bento directory exists
	if err := os.MkdirAll(bentoDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create bento directory: %w", err)
	}

	mgr := &VariablesManager{
		filePath:  filePath,
		variables: make(map[string]string),
	}

	// Load existing variables if file exists
	if err := mgr.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return mgr, nil
}

// load reads variables from JSON file
func (m *VariablesManager) load() error {
	data, err := os.ReadFile(m.filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &m.variables)
}

// save writes variables to JSON file
func (m *VariablesManager) save() error {
	data, err := json.MarshalIndent(m.variables, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal variables: %w", err)
	}

	if err := os.WriteFile(m.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write variables file: %w", err)
	}

	return nil
}

// Set stores a variable
func (m *VariablesManager) Set(key, value string) error {
	if key == "" {
		return fmt.Errorf("variable key cannot be empty")
	}

	m.variables[key] = value
	return m.save()
}

// Get retrieves a variable
func (m *VariablesManager) Get(key string) (string, error) {
	value, ok := m.variables[key]
	if !ok {
		return "", fmt.Errorf("variable %s not found", key)
	}
	return value, nil
}

// Delete removes a variable
func (m *VariablesManager) Delete(key string) error {
	if _, ok := m.variables[key]; !ok {
		return fmt.Errorf("variable %s not found", key)
	}

	delete(m.variables, key)
	return m.save()
}

// List returns all variable keys
func (m *VariablesManager) List() []string {
	keys := make([]string, 0, len(m.variables))
	for key := range m.variables {
		keys = append(keys, key)
	}
	return keys
}

// GetAll returns all variables as a map
func (m *VariablesManager) GetAll() map[string]string {
	result := make(map[string]string, len(m.variables))
	for k, v := range m.variables {
		result[k] = v
	}
	return result
}
