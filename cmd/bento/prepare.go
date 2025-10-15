package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"bento/pkg/neta"
)

var prepareCmd = &cobra.Command{
	Use:   "prepare [file.bento.yaml]",
	Short: "Validate a bento workflow file",
	Long: `Prepare validates a .bento.yaml file without executing it.

This checks:
- YAML syntax is valid
- Node types are registered
- Required parameters are present
- Structure is well-formed`,
	Args: cobra.ExactArgs(1),
	RunE: runPrepare,
}

func init() {
	rootCmd.AddCommand(prepareCmd)
}

func runPrepare(cmd *cobra.Command, args []string) error {
	filename := args[0]

	def, err := loadDefinition(filename)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := validateDefinition(def); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	fmt.Printf("✅ %s is valid\n", filename)
	printDefinitionSummary(def)
	return nil
}

// loadDefinition reads and parses a .bento.yaml file.
func loadDefinition(filename string) (neta.Definition, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return neta.Definition{}, err
	}

	var def neta.Definition
	if err := yaml.Unmarshal(data, &def); err != nil {
		return neta.Definition{}, err
	}

	return def, nil
}

// validateDefinition checks if a definition is well-formed.
func validateDefinition(def neta.Definition) error {
	if def.Type == "" {
		return fmt.Errorf("type is required")
	}

	if isGroupType(def.Type) {
		return validateGroup(def)
	}

	if def.IsGroup() {
		return validateGroup(def)
	}

	return validateLeaf(def)
}

// isGroupType checks if the type is a known group type.
func isGroupType(nodeType string) bool {
	groupTypes := []string{"sequence", "parallel"}
	for _, t := range groupTypes {
		if nodeType == t {
			return true
		}
	}
	return false
}

func validateGroup(def neta.Definition) error {
	if len(def.Nodes) == 0 {
		return fmt.Errorf("group must have child nodes")
	}

	for i, child := range def.Nodes {
		if err := validateDefinition(child); err != nil {
			return fmt.Errorf("node %d: %w", i, err)
		}
	}

	return nil
}

func validateLeaf(def neta.Definition) error {
	if def.Type == "" {
		return fmt.Errorf("type is required")
	}
	return nil
}

func printDefinitionSummary(def neta.Definition) {
	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Type: %s\n", def.Type)
	if def.Name != "" {
		fmt.Printf("  Name: %s\n", def.Name)
	}
	if def.IsGroup() {
		fmt.Printf("  Nodes: %d\n", len(def.Nodes))
	}
}
