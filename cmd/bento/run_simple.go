// Package main implements simple execution for the run command.
//
// This file contains the executeSimple function which runs bento workflows
// with simple single-line progress output for non-TTY environments (CI/CD, pipes, etc).
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Develonaut/bento/pkg/itamae"
	"github.com/Develonaut/bento/pkg/miso"
	"github.com/Develonaut/bento/pkg/neta"
)

// executeSimple executes bento with simple single-line progress (non-TTY mode).
func executeSimple(def *neta.Definition) error {
	// Get theme and palette from miso manager
	manager := miso.NewManager()
	theme := manager.GetTheme()
	palette := manager.GetPalette()

	// Create simple messenger that prints to stdout
	messenger := miso.NewSimpleMessenger(theme, palette)

	// Create pantry and file logger (always log to file)
	p := createPantry()
	logger, logFile, err := createFileLogger()
	if err != nil {
		printError(fmt.Sprintf("Warning: Failed to create log file: %v", err))
		// Continue without file logging
	}
	if logFile != nil {
		defer logFile.Close()
	}

	// Also log to stdout if verbose
	if verboseFlag && logger != nil {
		logger = createDualLogger(logger)
	}

	// Create chef with messenger
	chef := itamae.NewWithMessenger(p, logger, messenger)

	// Execute bento
	ctx, cancel := context.WithTimeout(context.Background(), timeoutFlag)
	defer cancel()

	start := time.Now()
	result, err := chef.Serve(ctx, def)
	duration := time.Since(start)

	if err != nil {
		// Get random error status word
		statusWord := getErrorStatusWord()
		printError(fmt.Sprintf("Oh no! Bento is %s: %v", statusWord, err))
		os.Exit(1)
	}

	printSuccess(fmt.Sprintf("Delicious! Bento executed successfully in %s", formatDuration(duration)))
	fmt.Printf("   ✓ %d nodes executed\n", result.NodesExecuted)
	return nil
}
