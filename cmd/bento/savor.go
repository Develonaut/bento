// Package main implements the savor command for executing bentos.
//
// The savor command loads a bento definition from a JSON file and executes it
// using the itamae orchestration engine. It provides real-time progress updates
// and displays execution results with fun sushi-themed output.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Develonaut/bento/pkg/itamae"
	"github.com/Develonaut/bento/pkg/neta"
	"github.com/Develonaut/bento/pkg/omakase"
	"github.com/Develonaut/bento/pkg/pantry"
	"github.com/Develonaut/bento/pkg/shoyu"
	"github.com/spf13/cobra"

	editfields "github.com/Develonaut/bento/pkg/neta/library/editfields"
	filesystem "github.com/Develonaut/bento/pkg/neta/library/filesystem"
	group "github.com/Develonaut/bento/pkg/neta/library/group"
	httpneta "github.com/Develonaut/bento/pkg/neta/library/http"
	image "github.com/Develonaut/bento/pkg/neta/library/image"
	loop "github.com/Develonaut/bento/pkg/neta/library/loop"
	parallel "github.com/Develonaut/bento/pkg/neta/library/parallel"
	shellcommand "github.com/Develonaut/bento/pkg/neta/library/shellcommand"
	spreadsheet "github.com/Develonaut/bento/pkg/neta/library/spreadsheet"
	transform "github.com/Develonaut/bento/pkg/neta/library/transform"
)

var (
	verboseFlag bool
	timeoutFlag time.Duration
)

var savorCmd = &cobra.Command{
	Use:   "savor [file].bento.json",
	Short: "🍱 Savor a bento (execute workflow)",
	Long: `Execute a bento workflow from start to finish.

Savor your bento! This command executes all neta in the bento
and reports progress with delicious output.

Examples:
  bento savor workflow.bento.json
  bento savor workflow.bento.json --verbose
  bento savor workflow.bento.json --timeout 30m`,
	Args: cobra.ExactArgs(1),
	RunE: runSavor,
}

func init() {
	savorCmd.Flags().BoolVarP(&verboseFlag, "verbose", "v", false, "Verbose output")
	savorCmd.Flags().DurationVar(&timeoutFlag, "timeout", 10*time.Minute, "Execution timeout")
}

// runSavor executes the savor command logic.
func runSavor(cmd *cobra.Command, args []string) error {
	def, err := loadAndValidate(args[0])
	if err != nil {
		return err
	}

	chef := setupChef()
	return executeBento(chef, def)
}

// loadAndValidate loads and validates a bento.
func loadAndValidate(bentoPath string) (*neta.Definition, error) {
	def, err := loadBento(bentoPath)
	if err != nil {
		printError(fmt.Sprintf("Failed to load bento: %v", err))
		return nil, err
	}

	printInfo(fmt.Sprintf("Savoring bento: %s", def.Name))

	if err := validateBento(def); err != nil {
		printError(fmt.Sprintf("Validation failed: %v", err))
		return nil, err
	}

	return def, nil
}

// setupChef creates and configures the itamae.
func setupChef() *itamae.Itamae {
	p := createPantry()
	logger := createLogger()
	chef := itamae.New(p, logger)
	setupProgress(chef)
	return chef
}

// executeBento executes the bento and reports results.
func executeBento(chef *itamae.Itamae, def *neta.Definition) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutFlag)
	defer cancel()

	start := time.Now()
	result, err := chef.Serve(ctx, def)
	duration := time.Since(start)

	if err != nil {
		printError(fmt.Sprintf("Execution failed: %v", err))
		return err
	}

	printSuccess(fmt.Sprintf("Delicious! Bento savored successfully in %s", formatDuration(duration)))
	fmt.Printf("   %d neta executed\n", result.NodesExecuted)
	return nil
}

// loadBento loads a bento definition from a JSON file.
func loadBento(path string) (*neta.Definition, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var def neta.Definition
	if err := json.Unmarshal(data, &def); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &def, nil
}

// createPantry creates and populates the pantry with all neta types.
func createPantry() *pantry.Pantry {
	p := pantry.New()

	// Register all neta types
	p.RegisterFactory("edit-fields", func() neta.Executable { return editfields.New() })
	p.RegisterFactory("file-system", func() neta.Executable { return filesystem.New() })
	p.RegisterFactory("group", func() neta.Executable { return group.New() })
	p.RegisterFactory("http-request", func() neta.Executable { return httpneta.New() })
	p.RegisterFactory("image", func() neta.Executable { return image.New() })
	p.RegisterFactory("loop", func() neta.Executable { return loop.New() })
	p.RegisterFactory("parallel", func() neta.Executable { return parallel.New() })
	p.RegisterFactory("shell-command", func() neta.Executable { return shellcommand.New() })
	p.RegisterFactory("spreadsheet", func() neta.Executable { return spreadsheet.New() })
	p.RegisterFactory("transform", func() neta.Executable { return transform.New() })

	return p
}

// createLogger creates a logger with appropriate level.
func createLogger() *shoyu.Logger {
	level := shoyu.LevelInfo
	if verboseFlag {
		level = shoyu.LevelDebug
	}

	return shoyu.New(shoyu.Config{
		Level:  level,
		Format: shoyu.FormatConsole,
	})
}

// setupProgress configures progress callbacks for the itamae.
func setupProgress(chef *itamae.Itamae) {
	chef.OnProgress(func(nodeID, status string) {
		if verboseFlag {
			if status == "starting" {
				printProgress(fmt.Sprintf("Savoring neta '%s'...", nodeID))
			} else if status == "completed" {
				printCheck("Delicious!")
			}
		}
	})
}

// validateBento validates the bento definition before execution.
func validateBento(def *neta.Definition) error {
	validator := omakase.New()
	ctx := context.Background()
	return validator.Validate(ctx, def)
}
