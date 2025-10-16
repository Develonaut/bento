package screens

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"bento/pkg/omise/styles"
)

// handleThemeAndInput handles theme changes and keyboard input
func (e Executor) handleThemeAndInput(msg tea.Msg) (Executor, tea.Cmd, bool) {
	// Handle theme changes
	if _, ok := msg.(styles.ThemeChangedMsg); ok {
		e.spinner = e.spinner.RebuildStyles()
		e.progress = e.progress.RebuildStyles()
		return e, nil, true
	}

	// Handle keyboard input when complete
	if e.complete {
		return e.handleCompleteKeys(msg)
	}

	return e, nil, false
}

// handleCompleteKeys handles keyboard input when execution is complete
func (e Executor) handleCompleteKeys(msg tea.Msg) (Executor, tea.Cmd, bool) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return e, nil, false
	}

	switch keyMsg.String() {
	case "c":
		return e, e.copyToClipboard(), true
	case "r":
		return e, nil, true
	}

	return e, nil, false
}

// handleExecutionMessages handles execution-related messages
func (e Executor) handleExecutionMessages(msg tea.Msg) (Executor, tea.Cmd) {
	if updated, cmd, handled := e.handleNodeMessages(msg); handled {
		return updated, cmd
	}
	switch msg := msg.(type) {
	case ExecutionProgressMsg:
		return e.handleProgressMsg(msg)
	case executionStillRunningMsg:
		return e.handleStillRunningMsg()
	case ExecutionCompleteMsg:
		return e.handleCompleteMsg(msg)
	case CopyResultMsg:
		return e.handleCopyMsg(msg)
	case ExecutionErrorMsg:
		return e.handleErrorMsg(msg)
	}
	return e.updateSpinner(msg)
}

// handleNodeMessages handles node-specific execution messages
func (e Executor) handleNodeMessages(msg tea.Msg) (Executor, tea.Cmd, bool) {
	switch msg := msg.(type) {
	case ExecutionInitMsg:
		updated, cmd := e.handleInitMsg(msg)
		return updated, cmd, true
	case NodeStartedMsg:
		updated, cmd := e.handleNodeStarted(msg)
		return updated, cmd, true
	case NodeCompletedMsg:
		updated, cmd := e.handleNodeCompleted(msg)
		return updated, cmd, true
	}
	return e, nil, false
}

// handleCopyMsg handles copy result messages
func (e Executor) handleCopyMsg(msg CopyResultMsg) (Executor, tea.Cmd) {
	e.copyFeedback = string(msg)
	return e, nil
}

// handleProgressMsg handles execution progress updates
func (e Executor) handleProgressMsg(msg ExecutionProgressMsg) (Executor, tea.Cmd) {
	if !e.running {
		return e, nil
	}
	e.status = msg.Status
	e.progressPercent = msg.Progress

	return e, tea.Batch(
		WaitForExecutionCmd(),
		ProgressTickCmd(msg.Progress),
	)
}

// handleStillRunningMsg continues polling for completion
func (e Executor) handleStillRunningMsg() (Executor, tea.Cmd) {
	if !e.running {
		return e, nil
	}
	return e, WaitForExecutionCmd()
}

// handleCompleteMsg handles execution completion
func (e Executor) handleCompleteMsg(msg ExecutionCompleteMsg) (Executor, tea.Cmd) {
	e.endTime = time.Now()
	e.running = false
	e.complete = true
	e.success = msg.Success
	e.result = formatResult(msg.Result)
	e.progressPercent = 1.0
	e.status = "Execution completed successfully!"

	// Add completion message to history
	e.lifecycleHistory = append(e.lifecycleHistory, emojiBento+" Bento packed!")

	if !msg.Success {
		e.status = "Execution failed"
		if msg.Error != nil {
			e.errorMsg = msg.Error.Error()
		}
	}
	return e, nil
}

// handleErrorMsg handles execution errors
func (e Executor) handleErrorMsg(msg ExecutionErrorMsg) (Executor, tea.Cmd) {
	e.running = false
	e.complete = true
	e.success = false
	e.status = "Execution error"
	e.errorMsg = msg.Error.Error()
	e.progressPercent = 0.0
	return e, nil
}

// updateSpinner updates the spinner during execution
func (e Executor) updateSpinner(msg tea.Msg) (Executor, tea.Cmd) {
	if !e.running {
		return e, nil
	}
	var spinnerCmd tea.Cmd
	e.spinner, spinnerCmd = e.spinner.Update(msg)
	return e, spinnerCmd
}
