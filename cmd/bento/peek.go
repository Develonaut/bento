// Package main implements the peek command for validating bentos.
//
// The peek command loads a bento definition and validates it without executing.
// It checks structure, neta types, parameters, and edges to ensure the bento
// is properly configured.
package main

import (
	"context"
	"fmt"

	"github.com/Develonaut/bento/pkg/neta"
	"github.com/Develonaut/bento/pkg/omakase"
	"github.com/spf13/cobra"
)

var peekVerboseFlag bool

var peekCmd = &cobra.Command{
	Use:   "peek [file].bento.json",
	Short: "🍱 Peek at a bento (validate without executing)",
	Long: `Validate a bento without executing it.

Peek inside your bento box to check if everything looks good!
This validates structure, neta types, parameters, and connections.

Examples:
  bento peek workflow.bento.json
  bento peek workflow.bento.json --verbose`,
	Args: cobra.ExactArgs(1),
	RunE: runPeek,
}

func init() {
	peekCmd.Flags().BoolVarP(&peekVerboseFlag, "verbose", "v", false, "Show detailed validation results")
}

// runPeek executes the peek command logic.
func runPeek(cmd *cobra.Command, args []string) error {
	def, err := loadBentoForPeek(args[0])
	if err != nil {
		return err
	}

	if err := validateForPeek(def); err != nil {
		return err
	}

	showValidationResults()
	return nil
}

// loadBentoForPeek loads a bento and prints status.
func loadBentoForPeek(bentoPath string) (*neta.Definition, error) {
	def, err := loadBento(bentoPath)
	if err != nil {
		printError(fmt.Sprintf("Failed to load bento: %v", err))
		return nil, err
	}

	printInfo(fmt.Sprintf("Peeking at bento: %s", def.Name))
	return def, nil
}

// validateForPeek validates the bento definition.
func validateForPeek(def *neta.Definition) error {
	validator := omakase.New()
	ctx := context.Background()

	if err := validator.Validate(ctx, def); err != nil {
		printError(fmt.Sprintf("Validation failed: %v", err))
		return err
	}

	return nil
}

// showValidationResults displays validation results.
func showValidationResults() {
	if peekVerboseFlag {
		printCheck("Valid JSON structure")
		printCheck("All neta types recognized")
		printCheck("All edges valid")
		printCheck("Required parameters present")
	}

	printSuccess("Looks delicious! Ready to eat.")
}
