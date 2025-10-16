package screens

import (
	"context"
	"fmt"
	"time"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"

	"bento/pkg/itamae"
	"bento/pkg/jubako"
	"bento/pkg/neta/conditional"
	"bento/pkg/neta/group"
	"bento/pkg/neta/http"
	"bento/pkg/neta/loop"
	"bento/pkg/neta/transform"
	"bento/pkg/pantry"
)

// CopyResultCmd copies result to clipboard and returns feedback message
func CopyResultCmd(result, bentoName, errorMsg string, success bool) tea.Msg {
	// Build content based on what's available
	var content string

	if success && result != "" {
		// Success case with output
		content = fmt.Sprintf("Bento: %s\n\nStatus: Success\n\nOutput:\n%s", bentoName, result)
	} else if !success && errorMsg != "" {
		// Error case
		content = fmt.Sprintf("Bento: %s\n\nStatus: Failed\n\nError:\n%s", bentoName, errorMsg)
	} else if result != "" {
		// Has result but not explicitly success
		content = fmt.Sprintf("Bento: %s\n\nOutput:\n%s", bentoName, result)
	} else if errorMsg != "" {
		// Has error message only
		content = fmt.Sprintf("Bento: %s\n\nError:\n%s", bentoName, errorMsg)
	} else {
		// No content at all
		return CopyResultMsg("No output or error to copy")
	}

	if err := clipboard.WriteAll(content); err != nil {
		return CopyResultMsg(fmt.Sprintf("Failed to copy: %s", err.Error()))
	}

	return CopyResultMsg("✓ Copied to clipboard!")
}

// executionState tracks ongoing execution
var executionState struct {
	running bool
	done    chan ExecutionCompleteMsg
}

// ExecuteBentoCmd creates a command that executes a bento by name
func ExecuteBentoCmd(bentoName string, workDir string) tea.Cmd {
	// Reset and mark execution as running
	executionState.running = true
	executionState.done = make(chan ExecutionCompleteMsg, 1)

	// Start execution in goroutine
	go func() {
		// Load the bento definition
		store, err := jubako.NewStore(workDir)
		if err != nil {
			executionState.done <- ExecutionCompleteMsg{Success: false, Error: err}
			executionState.running = false
			return
		}

		def, err := store.Load(bentoName)
		if err != nil {
			executionState.done <- ExecutionCompleteMsg{Success: false, Error: err}
			executionState.running = false
			return
		}

		// Create pantry and register all standard node types
		registry := pantry.New()
		chef := itamae.New(registry)

		// Register all standard neta types with fully qualified names
		_ = registry.Register("http", http.New())
		_ = registry.Register("transform.jq", transform.NewJQ())
		_ = registry.Register("group.sequence", group.NewSequence(chef))
		_ = registry.Register("group.parallel", group.NewParallel(chef))
		_ = registry.Register("conditional.if", conditional.NewIf(chef))
		_ = registry.Register("loop.for", loop.NewFor(chef))

		// Execute with context
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		result, err := chef.Execute(ctx, def)
		executionState.running = false
		if err != nil {
			executionState.done <- ExecutionCompleteMsg{Success: false, Error: err}
			return
		}

		executionState.done <- ExecutionCompleteMsg{Success: true, Result: result}
	}()

	// Return initial progress message
	return func() tea.Msg {
		return ExecutionProgressMsg{
			Status:   "Loading bento definition...",
			Progress: 0.1,
		}
	}
}

// WaitForExecutionCmd checks if execution is complete (non-blocking)
func WaitForExecutionCmd() tea.Cmd {
	return func() tea.Msg {
		// Non-blocking check if there's a completion message ready
		select {
		case msg := <-executionState.done:
			return msg
		case <-time.After(50 * time.Millisecond):
			// Not done yet, return a message indicating we're still running
			return executionStillRunningMsg{}
		}
	}
}

// executionStillRunningMsg indicates execution is still in progress
type executionStillRunningMsg struct{}

// ProgressTickCmd generates periodic progress updates during execution
func ProgressTickCmd(currentProgress float64) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(150 * time.Millisecond)

		// Increment progress slowly, cap at 90% to leave room for completion
		newProgress := currentProgress + 0.05
		if newProgress > 0.9 {
			newProgress = 0.9
		}

		return ExecutionProgressMsg{
			Status:   "Executing bento...",
			Progress: newProgress,
		}
	}
}
