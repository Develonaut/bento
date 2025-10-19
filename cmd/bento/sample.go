// Package main implements the sample command for validating bentos.
//
// The sample command loads a bento definition and validates it without executing.
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

var sampleVerboseFlag bool

var sampleCmd = &cobra.Command{
	Use:   "sample [file].bento.json",
	Short: "🥢 Sample a bento (validate without executing)",
	Long: `Validate a bento without executing it.

Sample your bento to check if it tastes right before serving!
This validates structure, neta types, parameters, and connections.

Examples:
  bento sample workflow.bento.json
  bento sample workflow.bento.json --verbose`,
	Args: cobra.ExactArgs(1),
	RunE: runSample,
}

func init() {
	sampleCmd.Flags().BoolVarP(&sampleVerboseFlag, "verbose", "v", false, "Show detailed validation results")
}

// runSample executes the sample command logic.
func runSample(cmd *cobra.Command, args []string) error {
	def, err := loadBentoForSample(args[0])
	if err != nil {
		return err
	}

	if err := validateForSample(def); err != nil {
		return err
	}

	showValidationResults()
	return nil
}

// loadBentoForSample loads a bento and prints status.
func loadBentoForSample(bentoPath string) (*neta.Definition, error) {
	def, err := loadBento(bentoPath)
	if err != nil {
		printError(fmt.Sprintf("Failed to load bento: %v", err))
		return nil, err
	}

	printInfo(fmt.Sprintf("Sampling bento: %s", def.Name))
	return def, nil
}

// validateForSample validates the bento definition.
func validateForSample(def *neta.Definition) error {
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
	if sampleVerboseFlag {
		printCheck("Valid JSON structure")
		printCheck("All neta types recognized")
		printCheck("All edges valid")
		printCheck("Required parameters present")
	}

	printSuccess("Tastes great! Ready to serve.")
}
