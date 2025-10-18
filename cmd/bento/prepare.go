package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"bento/pkg/jubako"
	"bento/pkg/neta"
)

var prepareCmd = &cobra.Command{
	Use:   "prepare [file.bento.json]",
	Short: "Validate a bento file",
	Long: `Prepare validates a .bento.json file without executing it.

This checks:
- JSON syntax is valid
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

	// Use jubako.Parser which includes version validation
	parser := jubako.NewParser()
	def, err := parser.Parse(filename)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	fmt.Printf("✅ %s is valid\n", filename)
	printDefinitionSummary(def)
	return nil
}

func printDefinitionSummary(def neta.Definition) {
	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Version: %s\n", def.Version)
	fmt.Printf("  Type: %s\n", def.Type)
	if def.Name != "" {
		fmt.Printf("  Name: %s\n", def.Name)
	}
	if def.IsGroup() {
		fmt.Printf("  Nodes: %d\n", len(def.Nodes))
	}
}
