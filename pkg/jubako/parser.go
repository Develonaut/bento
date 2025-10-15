package jubako

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"bento/pkg/neta"
)

// Parser handles .bento.yaml file parsing.
type Parser struct{}

// NewParser creates a new parser.
func NewParser() *Parser {
	return &Parser{}
}

// Parse reads and parses a .bento.yaml file.
func (p *Parser) Parse(path string) (neta.Definition, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return neta.Definition{}, fmt.Errorf("read failed: %w", err)
	}

	return p.ParseBytes(data)
}

// ParseBytes parses .bento.yaml from bytes.
func (p *Parser) ParseBytes(data []byte) (neta.Definition, error) {
	var def neta.Definition
	if err := yaml.Unmarshal(data, &def); err != nil {
		return neta.Definition{}, fmt.Errorf("invalid YAML: %w", err)
	}

	if err := validateDefinition(def); err != nil {
		return neta.Definition{}, fmt.Errorf("validation failed: %w", err)
	}

	return def, nil
}

// Format converts a definition to YAML.
func (p *Parser) Format(def neta.Definition) ([]byte, error) {
	data, err := yaml.Marshal(def)
	if err != nil {
		return nil, fmt.Errorf("marshal failed: %w", err)
	}
	return data, nil
}

// validateDefinition ensures a definition is well-formed.
func validateDefinition(def neta.Definition) error {
	// Validate version and type first
	if err := validateStructure(def); err != nil {
		return err
	}

	// Validate node parameters using validation framework
	validator := neta.NewValidator()
	if err := validator.ValidateRecursive(def); err != nil {
		return err
	}

	return nil
}

// validateStructure recursively validates version and type of a definition and its children.
func validateStructure(def neta.Definition) error {
	// Validate version
	if err := neta.ValidateVersion(def.Version); err != nil {
		return fmt.Errorf("version error: %w", err)
	}

	// Validate type is present
	if def.Type == "" {
		return fmt.Errorf("type is required")
	}

	// Recursively validate child nodes
	for i, child := range def.Nodes {
		if err := validateStructure(child); err != nil {
			return fmt.Errorf("node %d: %w", i, err)
		}
	}

	return nil
}
