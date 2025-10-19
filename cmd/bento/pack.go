// Package main implements the pack command for creating template bentos.
//
// The pack command creates a new bento template file with a basic structure
// that users can customize. It provides a quick way to start a new workflow.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Develonaut/bento/pkg/neta"
	"github.com/spf13/cobra"
)

var templateType string

var packCmd = &cobra.Command{
	Use:   "pack [name]",
	Short: "🍱 Create a new bento template",
	Long: `Create a new bento template file.

Pack a fresh bento box with a template you can customize.

Examples:
  bento pack my-workflow
  bento pack my-workflow --type simple`,
	Args: cobra.ExactArgs(1),
	RunE: runPack,
}

func init() {
	packCmd.Flags().StringVar(&templateType, "type", "simple", "Template type (simple, loop, parallel)")
}

// runPack executes the pack command logic.
func runPack(cmd *cobra.Command, args []string) error {
	name := args[0]
	fileName := name + ".bento.json"

	if err := checkFileExists(fileName); err != nil {
		return err
	}

	if err := createBentoFile(name, fileName); err != nil {
		return err
	}

	showNextSteps(fileName)
	return nil
}

// checkFileExists checks if file already exists.
func checkFileExists(fileName string) error {
	if _, err := os.Stat(fileName); err == nil {
		printError(fmt.Sprintf("File '%s' already exists", fileName))
		return fmt.Errorf("file already exists: %s", fileName)
	}
	return nil
}

// createBentoFile creates a new bento template file.
func createBentoFile(name, fileName string) error {
	printInfo(fmt.Sprintf("Packing new bento: %s", name))

	template := createTemplate(name)
	if err := writeTemplate(fileName, template); err != nil {
		printError(fmt.Sprintf("Failed to create bento: %v", err))
		return err
	}

	printSuccess(fmt.Sprintf("Created: %s", fileName))
	return nil
}

// showNextSteps displays next steps after creating bento.
func showNextSteps(fileName string) {
	fmt.Println("\nNext steps:")
	fmt.Printf("  1. Edit %s\n", fileName)
	fmt.Printf("  2. Run: bento sniff %s\n", fileName)
	fmt.Printf("  3. Run: bento taste %s\n", fileName)
}

// createTemplate creates a template bento definition.
func createTemplate(name string) *neta.Definition {
	return &neta.Definition{
		ID:      name,
		Type:    "group",
		Version: "1.0.0",
		Name:    formatName(name),
		Position: neta.Position{
			X: 0,
			Y: 0,
		},
		Metadata: neta.Metadata{
			Tags: []string{"template"},
		},
		Parameters:  make(map[string]interface{}),
		InputPorts:  []neta.Port{},
		OutputPorts: []neta.Port{},
		Nodes:       []neta.Definition{createSampleNode()},
		Edges:       []neta.Edge{},
	}
}

// createSampleNode creates a sample edit-fields node.
func createSampleNode() neta.Definition {
	return neta.Definition{
		ID:      "sample-1",
		Type:    "edit-fields",
		Version: "1.0.0",
		Name:    "Sample Node",
		Position: neta.Position{
			X: 100,
			Y: 100,
		},
		Metadata: neta.Metadata{},
		Parameters: map[string]interface{}{
			"values": map[string]interface{}{
				"message": "Hello from bento!",
			},
		},
		InputPorts:  []neta.Port{},
		OutputPorts: []neta.Port{},
	}
}

// formatName converts a kebab-case or snake_case name to Title Case.
func formatName(name string) string {
	if len(name) == 0 {
		return name
	}
	return strings.ToUpper(string(name[0])) + name[1:]
}

// writeTemplate writes the template to a JSON file.
func writeTemplate(fileName string, template *neta.Definition) error {
	data, err := json.MarshalIndent(template, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize template: %w", err)
	}

	if err := os.WriteFile(fileName, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
