package executor

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"bento/pkg/itamae"
	"bento/pkg/jubako"
	"bento/pkg/neta"
	"bento/pkg/neta/conditional"
	"bento/pkg/neta/file"
	"bento/pkg/neta/group"
	"bento/pkg/neta/http"
	"bento/pkg/neta/loop"
	"bento/pkg/neta/transform"
	"bento/pkg/omise/config"
	"bento/pkg/pantry"
)

// executionState tracks ongoing execution
var executionState struct {
	running bool
	done    chan ExecutionCompleteMsg
	ready   chan struct{} // Signals when init message is processed
}

// ExecuteBentoCmd creates a command that executes a bento by name
func ExecuteBentoCmd(bentoName string, workDir string, program *tea.Program) tea.Cmd {
	executionState.running = true
	executionState.done = make(chan ExecutionCompleteMsg, 1)
	executionState.ready = make(chan struct{})

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
	waitForInitialization(program)

	result, err := runBentoExecution(def, program)
	executionState.running = false

	sendExecutionResult(result, err)
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

// waitForInitialization waits for UI to be ready for execution
func waitForInitialization(program *tea.Program) {
	select {
	case <-executionState.ready:
		// Ready to proceed with execution
	case <-time.After(5 * time.Second):
		// Timeout - proceed anyway to avoid deadlock
		if program != nil {
			program.Send(ExecutionErrorMsg{
				Error: fmt.Errorf("initialization timeout - proceeding with execution"),
			})
		}
	}
}

// sendExecutionResult sends the completion message
func sendExecutionResult(result neta.Result, err error) {
	if err != nil {
		executionState.done <- ExecutionCompleteMsg{Success: false, Error: err}
		return
	}
	executionState.done <- ExecutionCompleteMsg{Success: true, Result: result}
}

// runBentoExecution creates chef and executes bento
func runBentoExecution(def neta.Definition, program *tea.Program) (neta.Result, error) {
	registry := pantry.New()
	messenger := &executorMessenger{program: program}
	chef := itamae.NewWithMessenger(registry, messenger)

	// Create and attach execution graph store for graph-based execution tracking
	store := neta.NewExecutionGraphStore()
	chef.SetStore(store)

	// Subscribe to store changes to update UI
	store.Subscribe(func(state neta.ExecutionGraphState) {
		if program != nil {
			program.Send(GraphStateUpdateMsg{State: state})
		}
	})

	// Load config and apply slow-mo delay if configured
	cfg := config.Load()
	if cfg.SlowMoDelayMs > 0 {
		chef.SetSlowMoDelay(cfg.SlowMoDelayMs)
	}

	registerStandardNetas(registry, chef)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	return chef.Execute(ctx, def)
}

// registerStandardNetas registers all built-in neta types
func registerStandardNetas(registry *pantry.Pantry, chef *itamae.Itamae) {
	_ = registry.Register("http", http.New())
	_ = registry.Register("transform.jq", transform.NewJQ())
	_ = registry.Register("file.write", file.NewWriter())
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

// signalInitReadyCmd signals that initialization is complete
func signalInitReadyCmd() tea.Cmd {
	return func() tea.Msg {
		// Signal the background goroutine that init is complete
		// Non-blocking send in case execution already failed
		select {
		case executionState.ready <- struct{}{}:
		default:
		}
		return nil
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
