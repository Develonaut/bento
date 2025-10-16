package screens

import (
	"context"
	"fmt"
	"time"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"

	"bento/pkg/itamae"
	"bento/pkg/jubako"
	"bento/pkg/neta"
	"bento/pkg/neta/conditional"
	"bento/pkg/neta/group"
	"bento/pkg/neta/http"
	"bento/pkg/neta/loop"
	"bento/pkg/neta/transform"
	"bento/pkg/pantry"
)

// executorMessenger sends progress messages to TUI
type executorMessenger struct {
	program *tea.Program
}

// SendNodeStarted sends node start message
func (m *executorMessenger) SendNodeStarted(path, name, nodeType string) {
	if m.program != nil {
		m.program.Send(NodeStartedMsg{
			Path:     path,
			Name:     name,
			NodeType: nodeType,
		})
	}
}

// SendNodeCompleted sends node completion message
func (m *executorMessenger) SendNodeCompleted(path string, duration time.Duration, err error) {
	if m.program != nil {
		m.program.Send(NodeCompletedMsg{
			Path:     path,
			Duration: duration,
			Error:    err,
		})
	}
}

// CopyResultCmd copies result to clipboard and returns feedback message
func CopyResultCmd(result, bentoName, errorMsg string, success bool) tea.Msg {
	content := buildClipboardContent(result, bentoName, errorMsg, success)
	if content == "" {
		return CopyResultMsg("No output or error to copy")
	}

	if err := clipboard.WriteAll(content); err != nil {
		return CopyResultMsg(fmt.Sprintf("Failed to copy: %s", err.Error()))
	}

	return CopyResultMsg("✓ Copied to clipboard!")
}

// buildClipboardContent formats content for clipboard
func buildClipboardContent(result, bentoName, errorMsg string, success bool) string {
	if success && result != "" {
		return fmt.Sprintf("Bento: %s\n\nStatus: Success\n\nOutput:\n%s", bentoName, result)
	}
	if !success && errorMsg != "" {
		return fmt.Sprintf("Bento: %s\n\nStatus: Failed\n\nError:\n%s", bentoName, errorMsg)
	}
	if result != "" {
		return fmt.Sprintf("Bento: %s\n\nOutput:\n%s", bentoName, result)
	}
	if errorMsg != "" {
		return fmt.Sprintf("Bento: %s\n\nError:\n%s", bentoName, errorMsg)
	}
	return ""
}

// executionState tracks ongoing execution
var executionState struct {
	running bool
	done    chan ExecutionCompleteMsg
}

// ExecuteBentoCmd creates a command that executes a bento by name
func ExecuteBentoCmd(bentoName string, workDir string, program *tea.Program) tea.Cmd {
	executionState.running = true
	executionState.done = make(chan ExecutionCompleteMsg, 1)

	go executeBentoInBackground(bentoName, workDir, program)

	return initialProgressMsg
}

// executeBentoInBackground runs bento execution in background
func executeBentoInBackground(bentoName, workDir string, program *tea.Program) {
	def, err := loadBentoDefinition(bentoName, workDir)
	if err != nil {
		sendExecutionError(err)
		return
	}

	sendInitMessage(program, def)

	result, err := runBentoExecution(def, program)
	executionState.running = false

	if err != nil {
		executionState.done <- ExecutionCompleteMsg{Success: false, Error: err}
		return
	}

	executionState.done <- ExecutionCompleteMsg{Success: true, Result: result}
}

// loadBentoDefinition loads bento from store
func loadBentoDefinition(bentoName, workDir string) (neta.Definition, error) {
	store, err := jubako.NewStore(workDir)
	if err != nil {
		return neta.Definition{}, err
	}
	return store.Load(bentoName)
}

// sendExecutionError sends error and marks execution as not running
func sendExecutionError(err error) {
	executionState.done <- ExecutionCompleteMsg{Success: false, Error: err}
	executionState.running = false
}

// sendInitMessage sends initialization message to TUI
func sendInitMessage(program *tea.Program, def neta.Definition) {
	if program != nil {
		program.Send(ExecutionInitMsg{Definition: def})
	}
}

// runBentoExecution creates chef and executes bento
func runBentoExecution(def neta.Definition, program *tea.Program) (neta.Result, error) {
	registry := pantry.New()
	messenger := &executorMessenger{program: program}
	chef := itamae.NewWithMessenger(registry, messenger)

	registerStandardNetas(registry, chef)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	return chef.Execute(ctx, def)
}

// registerStandardNetas registers all built-in neta types
func registerStandardNetas(registry *pantry.Pantry, chef *itamae.Itamae) {
	_ = registry.Register("http", http.New())
	_ = registry.Register("transform.jq", transform.NewJQ())
	_ = registry.Register("group.sequence", group.NewSequence(chef))
	_ = registry.Register("group.parallel", group.NewParallel(chef))
	_ = registry.Register("conditional.if", conditional.NewIf(chef))
	_ = registry.Register("loop.for", loop.NewFor(chef))
}

// initialProgressMsg returns initial progress message
func initialProgressMsg() tea.Msg {
	return ExecutionProgressMsg{
		Status:   "Loading bento definition...",
		Progress: 0.1,
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
