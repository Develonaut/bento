// Package main implements the run command for executing bentos.
//
// The run command loads a bento definition from a JSON file and executes it
// using the itamae orchestration engine. It provides real-time progress updates
// and displays execution results.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Develonaut/bento/pkg/hangiri"
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
	dryRunFlag  bool
)

var runCmd = &cobra.Command{
	Use:   "run [file].bento.json",
	Short: "Execute a bento workflow",
	Long: `Execute a bento workflow from start to finish.

This command executes all nodes in the bento workflow
and reports progress and results.

Examples:
  bento run workflow.bento.json
  bento run workflow.bento.json --verbose
  bento run workflow.bento.json --timeout 30m`,
	Args: cobra.ExactArgs(1),
	RunE: runRun,
}

func init() {
	runCmd.Flags().BoolVarP(&verboseFlag, "verbose", "v", false, "Verbose output")
	runCmd.Flags().DurationVar(&timeoutFlag, "timeout", 10*time.Minute, "Execution timeout")
	runCmd.Flags().BoolVar(&dryRunFlag, "dry-run", false, "Show what would be executed without running")
}

// runRun executes the run command logic.
func runRun(cmd *cobra.Command, args []string) error {
	def, err := loadAndValidate(args[0])
	if err != nil {
		return err
	}

	// If dry run, show what would be executed and exit
	if dryRunFlag {
		return showDryRun(def)
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

	printInfo(fmt.Sprintf("Running bento: %s", def.Name))

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

	printSuccess(fmt.Sprintf("Delicious! Bento executed successfully in %s", formatDuration(duration)))
	fmt.Printf("   %d nodes executed\n", result.NodesExecuted)
	return nil
}

// loadBento loads a bento definition from a file path or from storage.
//
// If the path is a valid file, it loads from that file.
// If the path doesn't exist as a file, it tries to load from ~/.bento/bentos/
// This allows users to run: bento run my-workflow
// instead of: bento run ~/.bento/bentos/my-workflow.bento.json
func loadBento(path string) (*neta.Definition, error) {
	// First, try to load as a direct file path
	if isValidFilePath(path) {
		return loadBentoFromFile(path)
	}

	// If not a valid file path, try loading from hangiri storage
	return loadBentoFromStorage(path)
}

// isValidFilePath checks if the path exists as a file.
func isValidFilePath(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// loadBentoFromFile loads a bento from a specific file path.
func loadBentoFromFile(path string) (*neta.Definition, error) {
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

// loadBentoFromStorage loads a bento from hangiri storage by name.
func loadBentoFromStorage(name string) (*neta.Definition, error) {
	// Strip .bento.json extension if provided
	name = strings.TrimSuffix(name, ".bento.json")

	storage := hangiri.NewDefaultStorage()
	ctx := context.Background()

	def, err := storage.LoadBento(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to load bento from storage: %w", err)
	}

	return def, nil
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

// createLogger creates a logger with appropriate level and streaming callback.
func createLogger() *shoyu.Logger {
	level := shoyu.LevelInfo
	if verboseFlag {
		level = shoyu.LevelDebug
	}

	return shoyu.New(shoyu.Config{
		Level: level,
		// Enable streaming output for long-running processes
		// This outputs lines from shell-command neta in real-time
		OnStream: func(line string) {
			fmt.Println(line)
		},
	})
}

// setupProgress configures progress callbacks for the itamae.
func setupProgress(chef *itamae.Itamae) {
	chef.OnProgress(func(nodeID, status string) {
		if verboseFlag {
			if status == "starting" {
				printProgress(fmt.Sprintf("Executing node '%s'...", nodeID))
			} else if status == "completed" {
				printCheck("Complete")
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

// showDryRun displays what would be executed without running.
func showDryRun(def *neta.Definition) error {
	printInfo("DRY RUN MODE - No execution will occur")
	fmt.Printf("\nWould execute bento: %s\n", def.Name)
	fmt.Printf("Total nodes to execute: %d\n\n", len(def.Nodes))

	if verboseFlag {
		printInfo("Nodes that would be executed:")
		for i, node := range def.Nodes {
			fmt.Printf("  %d. [%s] %s (type: %s)\n", i+1, node.ID, node.Name, node.Type)
		}
	}

	fmt.Println("\nValidation: ✓ Passed")
	printSuccess("Dry run complete. Use 'bento run' without --dry-run to execute.")
	return nil
}
