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

// ExecuteBentoCmd creates a command that executes a bento by name
func ExecuteBentoCmd(bentoName string, workDir string) tea.Cmd {
	return func() tea.Msg {
		// Load the bento definition
		store, err := jubako.NewStore(workDir)
		if err != nil {
			return ExecutionErrorMsg{Error: err}
		}

		def, err := store.Load(bentoName)
		if err != nil {
			return ExecutionErrorMsg{Error: err}
		}

		// Create pantry and register all standard node types
		registry := pantry.New()
		chef := itamae.New(registry)

		// Register all standard neta types
		_ = registry.Register("http", http.New())
		_ = registry.Register("jq", transform.NewJQ())
		_ = registry.Register("sequence", group.NewSequence(chef))
		_ = registry.Register("parallel", group.NewParallel(chef))
		_ = registry.Register("if", conditional.NewIf(chef))
		_ = registry.Register("for", loop.NewFor(chef))

		// Execute with context
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		// Send start message
		time.Sleep(100 * time.Millisecond) // Brief delay for UI

		result, err := chef.Execute(ctx, def)
		if err != nil {
			return ExecutionCompleteMsg{
				Success: false,
				Error:   err,
			}
		}

		return ExecutionCompleteMsg{
			Success: true,
			Result:  result,
		}
	}
}
