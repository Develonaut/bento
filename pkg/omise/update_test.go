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

func TestHandleKeySettings(t *testing.T) {
	m := NewModel()
	m.screen = ScreenBrowser

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
	updated, _ := m.handleKey(msg)
	model := updated.(Model)

	if model.screen != ScreenSettings {
		t.Errorf("Expected screen to change to Settings, got %v", model.screen)
	}
}
