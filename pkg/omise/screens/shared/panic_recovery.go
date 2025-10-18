package shared

import (
	"fmt"
	"runtime/debug"

	tea "github.com/charmbracelet/bubbletea"
)

// PanicRecoveryMsg is sent when a panic is recovered in a command
type PanicRecoveryMsg struct {
	Err        error
	StackTrace string
	Context    string
}

// RecoverFromPanic wraps a tea.Cmd with panic recovery
// This ensures the terminal is always restored even if a command panics
func RecoverFromPanic(cmd tea.Cmd, context string) tea.Cmd {
	if cmd == nil {
		return nil
	}

	return func() tea.Msg {
		defer func() {
			if r := recover(); r != nil {
				// Capture the panic and stack trace
				err := fmt.Errorf("panic in command (%s): %v", context, r)
				stack := string(debug.Stack())

				// Send a message about the panic instead of crashing
				// This allows the terminal to be restored properly
				tea.Printf("PANIC RECOVERED: %v\n%s\n", err, stack)
			}
		}()

		// Execute the original command
		return cmd()
	}
}

// WrapCmd is a convenience function for wrapping commands with panic recovery
func WrapCmd(context string) func(tea.Cmd) tea.Cmd {
	return func(cmd tea.Cmd) tea.Cmd {
		return RecoverFromPanic(cmd, context)
	}
}
