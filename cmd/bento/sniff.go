// Package main implements the sniff command for validating bentos.
//
// The sniff command loads a bento definition and validates it without executing.
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

var sniffVerboseFlag bool

var sniffCmd = &cobra.Command{
	Use:   "sniff [file].bento.json",
	Short: "🍱 Sniff a bento (validate without executing)",
	Long: `Validate a bento without executing it.

Sniff your bento to check if it's fresh and properly configured.
This validates structure, neta types, parameters, and connections.

Examples:
  bento sniff workflow.bento.json
  bento sniff workflow.bento.json --verbose`,
	Args: cobra.ExactArgs(1),
	RunE: runSniff,
}

func init() {
	sniffCmd.Flags().BoolVarP(&sniffVerboseFlag, "verbose", "v", false, "Show detailed validation results")
}

// runSniff executes the sniff command logic.
func runSniff(cmd *cobra.Command, args []string) error {
	def, err := loadBentoForSniff(args[0])
	if err != nil {
		return err
	}

	if err := validateForSniff(def); err != nil {
		return err
	}

	showValidationResults()
	return nil
}

// loadBentoForSniff loads a bento and prints status.
func loadBentoForSniff(bentoPath string) (*neta.Definition, error) {
	def, err := loadBento(bentoPath)
	if err != nil {
		printError(fmt.Sprintf("Failed to load bento: %v", err))
		return nil, err
	}

	printInfo(fmt.Sprintf("Sniffing bento: %s", def.Name))
	return def, nil
}

// validateForSniff validates the bento definition.
func validateForSniff(def *neta.Definition) error {
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
	if sniffVerboseFlag {
		printCheck("Valid JSON structure")
		printCheck("All neta types recognized")
		printCheck("All edges valid")
		printCheck("Required parameters present")
	}

	printSuccess("Smells fresh! Ready to taste.")
}
