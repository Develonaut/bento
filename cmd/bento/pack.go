package main

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"bento/pkg/itamae"
	"bento/pkg/neta"
)

var (
	packTimeout time.Duration
	packDryRun  bool
)

var packCmd = &cobra.Command{
	Use:   "pack [file.bento.yaml]",
	Short: "Execute a bento workflow",
	Long: `Pack executes a .bento.yaml workflow file.

This runs all nodes in the workflow and reports results.`,
	Args: cobra.ExactArgs(1),
	RunE: runPack,
}

func init() {
	rootCmd.AddCommand(packCmd)

	packCmd.Flags().DurationVar(&packTimeout, "timeout", 5*time.Minute, "execution timeout")
	packCmd.Flags().BoolVar(&packDryRun, "dry-run", false, "validate without executing")
}

func runPack(cmd *cobra.Command, args []string) error {
	filename := args[0]

	def, err := loadDefinition(filename)
	if err != nil {
		return fmt.Errorf("failed to load: %w", err)
	}

	if packDryRun {
		return runPrepare(cmd, args)
	}

	result, err := executeBento(def, packTimeout)
	if err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	printResult(result)
	return nil
}

// executeBento runs the workflow with the itamae orchestrator.
func executeBento(def neta.Definition, timeout time.Duration) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	p := initializePantry()
	chef := itamae.New(p)

	result, err := chef.Execute(ctx, def)
	if err != nil {
		return nil, err
	}

	return result.Output, nil
}

func printResult(result interface{}) {
	fmt.Println("\n✅ Execution complete!")
	fmt.Printf("\nResult:\n%v\n", result)
}
