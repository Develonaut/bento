package omise

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestHandleResize(t *testing.T) {
	m := NewModel()

	msg := tea.WindowSizeMsg{Width: 100, Height: 50}
	updated, _ := m.handleResize(msg)
	model := updated.(Model)

	if model.width != 100 {
		t.Errorf("Expected width 100, got %d", model.width)
	}
	if model.height != 50 {
		t.Errorf("Expected height 50, got %d", model.height)
	}
}

func TestHandleKeyQuit(t *testing.T) {
	m := NewModel()

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	updated, cmd := m.handleKey(msg)
	model := updated.(Model)

	if !model.quitting {
		t.Error("Expected quitting to be true after 'q' key")
	}

	if cmd == nil {
		t.Error("Expected tea.Quit command")
	}
}

func TestHandleKeyTab(t *testing.T) {
	m := NewModel()
	m.screen = ScreenBrowser

	msg := tea.KeyMsg{Type: tea.KeyTab}
	updated, _ := m.handleKey(msg)
	model := updated.(Model)

	if model.screen != ScreenExecutor {
		t.Errorf("Expected screen to change to Executor, got %v", model.screen)
	}
}

// Note: Modal mode blocking tests are in screens/settings_test.go
// where we can access the internal state of Settings

func TestHandleKeyTabWorksWhenNotInModalMode(t *testing.T) {
	m := NewModel()
	m.screen = ScreenSettings
	// NOT in modal mode (both false by default)

	msg := tea.KeyMsg{Type: tea.KeyTab}
	updated, _ := m.handleKey(msg)
	model := updated.(Model)

	// Tab SHOULD switch screens when not in modal mode
	if model.screen != ScreenHelp {
		t.Errorf("Expected screen to switch to Help, got %v", model.screen)
	}
}
