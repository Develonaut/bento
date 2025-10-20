package miso

import (
	"fmt"
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// ProgressMessenger receives execution progress events.
// Used to display real-time progress in TUI or CLI mode.
type ProgressMessenger interface {
	// SendNodeStarted notifies that a node has started execution.
	// path: node path in tree (e.g., "0", "1.2", "0.1.3" or node ID)
	// name: human-readable node name
	// nodeType: node type (e.g., "http-request", "file-system")
	SendNodeStarted(path, name, nodeType string)

	// SendNodeCompleted notifies that a node has finished execution.
	// path: node path in tree
	// duration: how long the node took to execute
	// err: error if node failed, nil if successful
	SendNodeCompleted(path string, duration time.Duration, err error)
}

// BubbletMessenger sends progress messages to a Bubbletea program.
// Used for TUI mode with real-time animated display.
type BubbletMessenger struct {
	program *tea.Program
}

// NewBubbletMessenger creates a messenger that sends to a Bubbletea program.
func NewBubbletMessenger(program *tea.Program) *BubbletMessenger {
	return &BubbletMessenger{
		program: program,
	}
}

// SendNodeStarted sends node start message to Bubbletea.
func (m *BubbletMessenger) SendNodeStarted(path, name, nodeType string) {
	if m.program != nil {
		m.program.Send(NodeStartedMsg{
			Path:     path,
			Name:     name,
			NodeType: nodeType,
		})
	}
}

// SendNodeCompleted sends node completion message to Bubbletea.
func (m *BubbletMessenger) SendNodeCompleted(path string, duration time.Duration, err error) {
	if m.program != nil {
		m.program.Send(NodeCompletedMsg{
			Path:     path,
			Duration: duration,
			Error:    err,
		})
	}
}

// SimpleMessenger prints simple progress updates for non-TTY mode.
// Used for CI/CD, pipes, and redirects where Bubbletea cannot run.
type SimpleMessenger struct {
	theme    *Theme
	palette  Palette
	nodeInfo map[string]nodeStartInfo // stores node info from start to completion
}

// nodeStartInfo stores information about a started node.
type nodeStartInfo struct {
	name     string
	nodeType string
	emoji    string
}

// NewSimpleMessenger creates a messenger for simple progress output.
func NewSimpleMessenger(theme *Theme, palette Palette) *SimpleMessenger {
	return &SimpleMessenger{
		theme:    theme,
		palette:  palette,
		nodeInfo: make(map[string]nodeStartInfo),
	}
}

// SendNodeStarted stores node start information (doesn't print yet).
func (m *SimpleMessenger) SendNodeStarted(path, name, nodeType string) {
	emoji := getStepEmoji(name)

	// Store node info for when it completes
	m.nodeInfo[path] = nodeStartInfo{
		name:     name,
		nodeType: nodeType,
		emoji:    emoji,
	}
}

// SendNodeCompleted prints the complete node execution line.
func (m *SimpleMessenger) SendNodeCompleted(path string, duration time.Duration, err error) {
	// Get stored node info
	info, ok := m.nodeInfo[path]
	if !ok {
		// Fallback if we don't have the info
		info = nodeStartInfo{
			name:  path,
			emoji: getStepEmoji(path),
		}
	}

	// Clean up stored info
	delete(m.nodeInfo, path)

	if err != nil {
		// Print: 💥 Failed Create Product Directory… (error message)
		statusWord := getStatusLabel(StepFailed, info.name)
		errorEmoji := getErrorEmoji()
		fmt.Printf("  %s %s %s… %s\n",
			errorEmoji,
			m.theme.Error.Render(statusWord),
			info.name,
			m.theme.Error.Render(err.Error()))
	} else {
		// Print: 🍥 Perfected Read Products CSV… (2ms)
		statusWord := getStatusLabel(StepCompleted, info.name)
		durationStr := formatSimpleDuration(duration)
		fmt.Printf("  %s %s %s… %s\n",
			info.emoji,
			m.theme.Success.Render(statusWord),
			info.name,
			m.theme.Subtle.Render(fmt.Sprintf("(%s)", durationStr)))
	}
}

// formatSimpleDuration formats duration for simple display.
func formatSimpleDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	mins := int(d.Minutes())
	secs := int(d.Seconds()) % 60
	return fmt.Sprintf("%dm %ds", mins, secs)
}

// getErrorEmoji returns a random error emoji for failed operations.
func getErrorEmoji() string {
	errorEmojis := []string{
		"👹",  // oni mask (Japanese demon)
		"👺",  // tengu/goblin mask
		"💀",  // skull
		"☠️", // skull and crossbones
		"💥",  // collision/explosion
		"🔥",  // fire
		"⚠️", // warning
		"❌",  // cross mark
		"🚫",  // no entry
	}
	return errorEmojis[rand.Intn(len(errorEmojis))]
}
