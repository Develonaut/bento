// Package main implements TUI execution for the run command.
//
// This file contains the executeTUI function which runs bento workflows
// with a real-time Bubbletea terminal UI display.
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Develonaut/bento/pkg/itamae"
	"github.com/Develonaut/bento/pkg/miso"
	"github.com/Develonaut/bento/pkg/neta"
)

// tuiContext holds the components needed for TUI execution.
type tuiContext struct {
	program  *tea.Program
	chef     *itamae.Itamae
	logFile  *os.File
	result   *itamae.Result
	err      error
	duration time.Duration
}

// executeTUI executes bento with Bubbletea TUI display.
func executeTUI(def *neta.Definition) error {
	ctx := setupTUIContext(def)
	if ctx.logFile != nil {
		defer ctx.logFile.Close()
	}

	runTUIExecution(ctx, def)

	return handleTUICompletion(ctx)
}

// setupTUIContext creates all components needed for TUI execution.
func setupTUIContext(def *neta.Definition) *tuiContext {
	manager := miso.NewManager()
	theme := manager.GetTheme()
	palette := manager.GetPalette()

	logger, logFile, err := createFileLogger()
	if err != nil {
		printError(fmt.Sprintf("Warning: Failed to create log file: %v", err))
	}

	p := createPantry()
	model := miso.NewExecutor(def, theme, palette)
	program := tea.NewProgram(model)
	messenger := miso.NewBubbletMessenger(program)
	chef := itamae.NewWithMessenger(p, logger, messenger)

	return &tuiContext{
		program: program,
		chef:    chef,
		logFile: logFile,
	}
}

// runTUIExecution starts background execution and runs the TUI program.
func runTUIExecution(ctx *tuiContext, def *neta.Definition) {
	go executeInBackground(ctx, def)

	finalModel, err := ctx.program.Run()
	if err != nil {
		printError(fmt.Sprintf("TUI error: %v", err))
		os.Exit(1)
	}

	checkExecutionSuccess(finalModel, ctx)
}

// executeInBackground runs bento execution in a goroutine.
func executeInBackground(ctx *tuiContext, def *neta.Definition) {
	time.Sleep(100 * time.Millisecond)
	ctx.program.Send(miso.ExecutionInitMsg{Definition: def})

	execCtx, cancel := context.WithTimeout(context.Background(), timeoutFlag)
	defer cancel()

	start := time.Now()
	result, err := ctx.chef.Serve(execCtx, def)
	ctx.duration = time.Since(start)
	ctx.result = result
	ctx.err = err

	ctx.program.Send(miso.ExecutionCompleteMsg{
		Success: err == nil,
		Error:   err,
	})
}

// checkExecutionSuccess validates the final model and checks success.
func checkExecutionSuccess(finalModel tea.Model, ctx *tuiContext) {
	finalExecutor, ok := finalModel.(miso.Executor)
	if !ok || !finalExecutor.Success() {
		if ctx.err != nil {
			statusWord := getErrorStatusWord()
			printError(fmt.Sprintf("Oh no! Bento is %s: %v", statusWord, ctx.err))
		}
		os.Exit(1)
	}
}

// handleTUICompletion prints success summary.
func handleTUICompletion(ctx *tuiContext) error {
	if ctx.result != nil {
		printSuccess(fmt.Sprintf("Delicious! Bento executed successfully in %s", formatDuration(ctx.duration)))
		fmt.Printf("   ✓ %d nodes executed\n", ctx.result.NodesExecuted)
	}
	return nil
}
