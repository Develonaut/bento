package screens

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
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
