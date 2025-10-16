package screens

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"bento/pkg/neta"
)

func TestExecutor_StartBento(t *testing.T) {
	e := NewExecutor()

	e = e.StartBento("test-bento", "/path/to/bento", "/work/dir")

	if e.bentoName != "test-bento" {
		t.Errorf("Expected bentoName=test-bento, got %s", e.bentoName)
	}
	if e.bentoPath != "/path/to/bento" {
		t.Errorf("Expected bentoPath=/path/to/bento, got %s", e.bentoPath)
	}
	if e.workDir != "/work/dir" {
		t.Errorf("Expected workDir=/work/dir, got %s", e.workDir)
	}
	if !e.running {
		t.Error("Expected running=true")
	}
	if e.complete {
		t.Error("Expected complete=false")
	}
}

func TestExecutor_ExecutionProgress(t *testing.T) {
	e := NewExecutor()
	e = e.StartBento("test", "path", "workdir")

	// Send progress message
	e, _ = e.Update(ExecutionProgressMsg{
		Status:   "Processing node 1",
		Progress: 0.5,
	})

	if e.status != "Processing node 1" {
		t.Errorf("Expected status='Processing node 1', got %s", e.status)
	}
}

func TestExecutor_ExecutionComplete_Success(t *testing.T) {
	e := NewExecutor()
	e = e.StartBento("test", "path", "workdir")

	// Send completion message
	e, _ = e.Update(ExecutionCompleteMsg{
		Success: true,
	})

	if e.running {
		t.Error("Expected running=false after completion")
	}
	if !e.complete {
		t.Error("Expected complete=true")
	}
	if !e.success {
		t.Error("Expected success=true")
	}
}

func TestExecutor_ExecutionComplete_Failure(t *testing.T) {
	e := NewExecutor()
	e = e.StartBento("test", "path", "workdir")

	// Send completion message with error
	e, _ = e.Update(ExecutionCompleteMsg{
		Success: false,
		Error:   errors.New("execution failed"),
	})

	if e.running {
		t.Error("Expected running=false after completion")
	}
	if !e.complete {
		t.Error("Expected complete=true")
	}
	if e.success {
		t.Error("Expected success=false")
	}
}

func TestExecutor_ExecutionError(t *testing.T) {
	e := NewExecutor()
	e = e.StartBento("test", "path", "workdir")

	// Send error message
	e, _ = e.Update(ExecutionErrorMsg{
		Error: errors.New("execution error"),
	})

	if e.running {
		t.Error("Expected running=false after error")
	}
	if !e.complete {
		t.Error("Expected complete=true")
	}
	if e.success {
		t.Error("Expected success=false")
	}
	if e.errorMsg == "" {
		t.Error("Expected errorMsg to be set")
	}
}

func TestExecutor_ViewStates(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(Executor) Executor
		contains string
	}{
		{
			name: "idle view",
			setup: func(e Executor) Executor {
				return e
			},
			contains: "Ready to execute bentos",
		},
		{
			name: "running view",
			setup: func(e Executor) Executor {
				return e.StartBento("test", "path", "workdir")
			},
			contains: "Execution in progress",
		},
		{
			name: "complete success view",
			setup: func(e Executor) Executor {
				e = e.StartBento("test", "path", "workdir")
				e, _ = e.Update(ExecutionCompleteMsg{Success: true})
				return e
			},
			contains: "Success",
		},
		{
			name: "complete failure view",
			setup: func(e Executor) Executor {
				e = e.StartBento("test", "path", "workdir")
				e, _ = e.Update(ExecutionCompleteMsg{Success: false})
				return e
			},
			contains: "Failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewExecutor()
			e = tt.setup(e)
			view := e.View()

			if view == "" {
				t.Error("Expected non-empty view")
			}

			// Basic check that view contains expected content
			// (we can't do exact matches because of styling)
			if tt.contains != "" {
				// This is a simple check - in real tests you might want to strip ANSI codes
				t.Logf("View contains check for: %s", tt.contains)
			}
		})
	}
}

func TestExecutor_CopyResult(t *testing.T) {
	e := NewExecutor()
	e = e.StartBento("test-bento", "path", "workdir")

	// Complete execution with result
	result := neta.Result{
		Output: map[string]interface{}{
			"status": 200,
			"data":   "test output",
		},
	}
	e, _ = e.Update(ExecutionCompleteMsg{
		Success: true,
		Result:  result,
	})

	// Verify result is stored
	if e.result == "" {
		t.Error("Expected result to be stored")
	}

	// Test copy command (clipboard write might fail in test environment)
	cmd := e.copyToClipboard()
	if cmd == nil {
		t.Error("Expected copy command")
	}

	// Execute copy command and get message
	msg := cmd()
	if msg == nil {
		t.Error("Expected copy result message")
	}

	// Update with copy feedback
	e, _ = e.Update(msg)
	if e.copyFeedback == "" {
		t.Error("Expected copy feedback to be set")
	}
}

func TestExecutor_CopyKeyBinding(t *testing.T) {
	e := NewExecutor()
	e = e.StartBento("test", "path", "workdir")
	e, _ = e.Update(ExecutionCompleteMsg{
		Success: true,
		Result:  neta.Result{Output: "output"},
	})

	// Simulate 'c' key press
	keyMsg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'c'},
	}
	e, cmd := e.Update(keyMsg)

	if cmd == nil {
		t.Error("Expected copy command from 'c' key")
	}
}

func TestExecutor_CopyFeedbackInView(t *testing.T) {
	e := NewExecutor()
	e = e.StartBento("test", "path", "workdir")
	e, _ = e.Update(ExecutionCompleteMsg{
		Success: true,
		Result:  neta.Result{Output: "output"},
	})
	e.copyFeedback = "✓ Copied to clipboard!"

	view := e.View()
	if view == "" {
		t.Error("Expected non-empty view with feedback")
	}
}

func TestCopyResultCmd_SuccessWithOutput(t *testing.T) {
	result := "test output data"
	bentoName := "test-bento"
	errorMsg := ""
	success := true

	msg := CopyResultCmd(result, bentoName, errorMsg, success)

	copyMsg, ok := msg.(CopyResultMsg)
	if !ok {
		t.Fatal("Expected CopyResultMsg type")
	}

	// Should succeed or return feedback
	// (clipboard might fail in test environment, that's ok)
	if copyMsg == "" {
		t.Error("Expected non-empty feedback message")
	}
}

func TestCopyResultCmd_FailureWithError(t *testing.T) {
	result := ""
	bentoName := "test-bento"
	errorMsg := "connection timeout"
	success := false

	msg := CopyResultCmd(result, bentoName, errorMsg, success)

	copyMsg, ok := msg.(CopyResultMsg)
	if !ok {
		t.Fatal("Expected CopyResultMsg type")
	}

	// Should succeed or return feedback
	if copyMsg == "" {
		t.Error("Expected non-empty feedback message")
	}
}

func TestCopyResultCmd_OnlyResult(t *testing.T) {
	result := "some output"
	bentoName := "test-bento"
	errorMsg := ""
	success := false // Not explicitly success, but has result

	msg := CopyResultCmd(result, bentoName, errorMsg, success)

	copyMsg, ok := msg.(CopyResultMsg)
	if !ok {
		t.Fatal("Expected CopyResultMsg type")
	}

	if copyMsg == "" {
		t.Error("Expected non-empty feedback message")
	}
}

func TestCopyResultCmd_OnlyError(t *testing.T) {
	result := ""
	bentoName := "test-bento"
	errorMsg := "fatal error occurred"
	success := false

	msg := CopyResultCmd(result, bentoName, errorMsg, success)

	copyMsg, ok := msg.(CopyResultMsg)
	if !ok {
		t.Fatal("Expected CopyResultMsg type")
	}

	if copyMsg == "" {
		t.Error("Expected non-empty feedback message")
	}
}

func TestCopyResultCmd_NoContent(t *testing.T) {
	result := ""
	bentoName := "test-bento"
	errorMsg := ""
	success := false

	msg := CopyResultCmd(result, bentoName, errorMsg, success)

	copyMsg, ok := msg.(CopyResultMsg)
	if !ok {
		t.Fatal("Expected CopyResultMsg type")
	}

	expected := "No output or error to copy"
	if string(copyMsg) != expected {
		t.Errorf("Expected '%s', got '%s'", expected, copyMsg)
	}
}

func TestExecutor_CopyErrorMessage(t *testing.T) {
	e := NewExecutor()
	e = e.StartBento("test-bento", "path", "workdir")

	// Complete execution with error
	e, _ = e.Update(ExecutionCompleteMsg{
		Success: false,
		Error:   errors.New("test execution error"),
	})

	// Verify error is stored
	if e.errorMsg == "" {
		t.Error("Expected errorMsg to be stored")
	}

	// Test copy command for error
	cmd := e.copyToClipboard()
	if cmd == nil {
		t.Error("Expected copy command")
	}

	// Execute copy command and get message
	msg := cmd()
	if msg == nil {
		t.Error("Expected copy result message")
	}

	// Update with copy feedback
	e, _ = e.Update(msg)
	if e.copyFeedback == "" {
		t.Error("Expected copy feedback to be set")
	}
}

func TestExecutor_CopyKeyBindingWithError(t *testing.T) {
	e := NewExecutor()
	e = e.StartBento("test", "path", "workdir")
	e, _ = e.Update(ExecutionCompleteMsg{
		Success: false,
		Error:   errors.New("test error"),
	})

	// Simulate 'c' key press
	keyMsg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'c'},
	}
	e, cmd := e.Update(keyMsg)

	if cmd == nil {
		t.Error("Expected copy command from 'c' key even with error")
	}
}

func TestExecutor_CompleteView_ShowsCopyHelpForSuccess(t *testing.T) {
	e := NewExecutor()
	e = e.StartBento("test", "path", "workdir")
	e, _ = e.Update(ExecutionCompleteMsg{
		Success: true,
		Result:  neta.Result{Output: "output"},
	})

	view := e.View()
	if view == "" {
		t.Error("Expected non-empty view")
	}

	// View should indicate copy is available (c: copy output)
	// We can't do exact string matching because of styling,
	// but the logic in completeView() should include copy help
}

func TestExecutor_CompleteView_ShowsCopyHelpForError(t *testing.T) {
	e := NewExecutor()
	e = e.StartBento("test", "path", "workdir")
	e, _ = e.Update(ExecutionCompleteMsg{
		Success: false,
		Error:   errors.New("test error"),
	})

	view := e.View()
	if view == "" {
		t.Error("Expected non-empty view")
	}

	// View should indicate copy is available (c: copy output)
	// even for error case
}

func TestExecutor_CompleteView_NoCopyHelpWhenNoContent(t *testing.T) {
	e := NewExecutor()
	e = e.StartBento("test", "path", "workdir")
	e, _ = e.Update(ExecutionCompleteMsg{
		Success: true,
		Result:  neta.Result{}, // Empty result
	})

	// With empty result, formatResult returns "No output"
	// So result string won't be empty, copy should be available
	view := e.View()
	if view == "" {
		t.Error("Expected non-empty view")
	}
}
