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
	// Verify the settings object is still valid
	if updated.list.Model.Items() == nil {
		t.Error("Reset key 'r' caused invalid state")
	}
}

func TestSettingsSpaceKey(t *testing.T) {
	s := NewSettings()
	s.list.Select(0) // Position on Theme setting

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
	s.list.Select(0) // Position on Theme setting

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
		name          string
		key           tea.KeyMsg
		initialIndex  int
		expectedIndex int
	}{
		{
			name:          "down moves selection",
			key:           tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}},
			initialIndex:  0,
			expectedIndex: 1,
		},
		{
			name:          "up moves selection",
			key:           tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}},
			initialIndex:  1,
			expectedIndex: 0,
		},
		{
			name:          "down at bottom stays",
			key:           tea.KeyMsg{Type: tea.KeyDown},
			initialIndex:  2,
			expectedIndex: 2,
		},
		{
			name:          "up at top stays",
			key:           tea.KeyMsg{Type: tea.KeyUp},
			initialIndex:  0,
			expectedIndex: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSettings()
			s.list.Select(tt.initialIndex)

			updated, _ := s.Update(tt.key)

			if updated.list.Index() != tt.expectedIndex {
				t.Errorf("list index = %d, want %d", updated.list.Index(), tt.expectedIndex)
			}
		})
	}
}

func TestSettingsBuildSettings(t *testing.T) {
	s := NewSettings()
	items := s.buildSettings()

	// Should have exactly 3 settings (Theme, Slow-Mo, and Save Directory)
	if len(items) != 3 {
		t.Errorf("Expected 3 settings, got %d", len(items))
	}

	// Verify Theme setting
	themeItem, ok := items[0].(settingItem)
	if !ok {
		t.Fatal("First item is not a settingItem")
	}
	if themeItem.name != "Theme" {
		t.Errorf("First setting should be Theme, got %s", themeItem.name)
	}
	if !themeItem.editable {
		t.Error("Theme setting should be editable")
	}

	// Verify Slow-Mo Execution setting
	slowMoItem, ok := items[1].(settingItem)
	if !ok {
		t.Fatal("Second item is not a settingItem")
	}
	if slowMoItem.name != "Slow-Mo Execution" {
		t.Errorf("Second setting should be Slow-Mo Execution, got %s", slowMoItem.name)
	}
	if !slowMoItem.editable {
		t.Error("Slow-Mo Execution setting should be editable")
	}

	// Verify Save Directory setting
	dirItem, ok := items[2].(settingItem)
	if !ok {
		t.Fatal("Third item is not a settingItem")
	}
	if dirItem.name != "Save Directory" {
		t.Errorf("Third setting should be Save Directory, got %s", dirItem.name)
	}
	if !dirItem.editable {
		t.Error("Save Directory setting should be editable")
	}

	// Verify description mentions "all app data"
	if dirItem.desc != "Directory for all app data (press Enter/Space to change)" {
		t.Errorf("Save Directory desc = %s, want 'Directory for all app data (press Enter/Space to change)'", dirItem.desc)
	}
}
