package screens

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestSettingsInModalMode(t *testing.T) {
	tests := []struct {
		name           string
		selectingDir   bool
		selectingTheme bool
		expectedModal  bool
	}{
		{
			name:           "not in modal mode",
			selectingDir:   false,
			selectingTheme: false,
			expectedModal:  false,
		},
		{
			name:           "in directory picker mode",
			selectingDir:   true,
			selectingTheme: false,
			expectedModal:  true,
		},
		{
			name:           "in theme picker mode",
			selectingDir:   false,
			selectingTheme: true,
			expectedModal:  true,
		},
		{
			name:           "both pickers active",
			selectingDir:   true,
			selectingTheme: true,
			expectedModal:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSettings()
			s.selectingDir = tt.selectingDir
			s.selectingTheme = tt.selectingTheme

			result := s.InModalMode()
			if result != tt.expectedModal {
				t.Errorf("InModalMode() = %v, want %v", result, tt.expectedModal)
			}
		})
	}
}

func TestSettingsResetKey(t *testing.T) {
	s := NewSettings()

	// Test lowercase 'r' triggers reset
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
	updated, _ := s.Update(msg)

	// Should have processed the key (not panic or error)
	if updated.cursor < 0 {
		t.Error("Reset key 'r' caused invalid state")
	}
}

func TestSettingsSpaceKey(t *testing.T) {
	s := NewSettings()
	s.cursor = 0 // Position on Theme setting

	// Test space key activates setting
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}}
	updated, _ := s.Update(msg)

	// Should enter theme selection mode
	if !updated.selectingTheme {
		t.Error("Space key should activate theme selection")
	}
}

func TestSettingsEnterKey(t *testing.T) {
	s := NewSettings()
	s.cursor = 0 // Position on Theme setting

	// Test enter key activates setting
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updated, _ := s.Update(msg)

	// Should enter theme selection mode
	if !updated.selectingTheme {
		t.Error("Enter key should activate theme selection")
	}
}

func TestSettingsTabExitsDirPicker(t *testing.T) {
	s := NewSettings()
	s.selectingDir = true

	// Test tab exits directory picker
	msg := tea.KeyMsg{Type: tea.KeyTab}
	updated, _ := s.Update(msg)

	if updated.selectingDir {
		t.Error("Tab key should exit directory picker")
	}
}

func TestSettingsTabExitsThemePicker(t *testing.T) {
	s := NewSettings()
	s.selectingTheme = true

	// Test tab exits theme picker
	msg := tea.KeyMsg{Type: tea.KeyTab}
	updated, _ := s.Update(msg)

	if updated.selectingTheme {
		t.Error("Tab key should exit theme picker")
	}
}

func TestSettingsEscExitsDirPicker(t *testing.T) {
	s := NewSettings()
	s.selectingDir = true

	// Test esc exits directory picker
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	updated, _ := s.Update(msg)

	if updated.selectingDir {
		t.Error("Esc key should exit directory picker")
	}
}

func TestSettingsEscExitsThemePicker(t *testing.T) {
	s := NewSettings()
	s.selectingTheme = true

	// Test esc exits theme picker
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	updated, _ := s.Update(msg)

	if updated.selectingTheme {
		t.Error("Esc key should exit theme picker")
	}
}

func TestSettingsNavigation(t *testing.T) {
	tests := []struct {
		name           string
		key            tea.KeyMsg
		initialCursor  int
		expectedCursor int
	}{
		{
			name:           "down moves cursor",
			key:            tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}},
			initialCursor:  0,
			expectedCursor: 1,
		},
		{
			name:           "up moves cursor",
			key:            tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}},
			initialCursor:  1,
			expectedCursor: 0,
		},
		{
			name:           "down at bottom stays",
			key:            tea.KeyMsg{Type: tea.KeyDown},
			initialCursor:  1,
			expectedCursor: 1,
		},
		{
			name:           "up at top stays",
			key:            tea.KeyMsg{Type: tea.KeyUp},
			initialCursor:  0,
			expectedCursor: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSettings()
			s.cursor = tt.initialCursor

			updated, _ := s.Update(tt.key)

			if updated.cursor != tt.expectedCursor {
				t.Errorf("cursor = %d, want %d", updated.cursor, tt.expectedCursor)
			}
		})
	}
}

func TestSettingsBuildSettings(t *testing.T) {
	s := NewSettings()
	items := s.buildSettings()

	// Should have exactly 2 settings (Theme and Save Directory)
	if len(items) != 2 {
		t.Errorf("Expected 2 settings, got %d", len(items))
	}

	// Verify Theme setting
	if items[0].name != "Theme" {
		t.Errorf("First setting should be Theme, got %s", items[0].name)
	}
	if !items[0].editable {
		t.Error("Theme setting should be editable")
	}

	// Verify Save Directory setting
	if items[1].name != "Save Directory" {
		t.Errorf("Second setting should be Save Directory, got %s", items[1].name)
	}
	if !items[1].editable {
		t.Error("Save Directory setting should be editable")
	}

	// Verify description mentions "all app data"
	if items[1].desc != "Directory for all app data (press Enter/Space to change)" {
		t.Errorf("Save Directory desc = %s, want 'Directory for all app data (press Enter/Space to change)'", items[1].desc)
	}
}
