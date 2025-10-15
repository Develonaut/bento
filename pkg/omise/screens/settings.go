package screens

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"bento/pkg/omise/styles"
)

// Settings shows configuration options
type Settings struct {
	cursor int
	items  []settingItem
}

type settingItem struct {
	name  string
	value string
	desc  string
}

// NewSettings creates a settings screen
func NewSettings() Settings {
	return Settings{
		cursor: 0,
		items:  defaultSettings(),
	}
}

// defaultSettings returns default setting items
func defaultSettings() []settingItem {
	return []settingItem{
		{
			name:  "Workflow Directory",
			value: "./workflows",
			desc:  "Default location for .bento.yaml files",
		},
		{
			name:  "Execution Timeout",
			value: "30s",
			desc:  "Default timeout for workflow execution",
		},
		{
			name:  "Log Level",
			value: "info",
			desc:  "Logging verbosity (debug, info, warn, error)",
		},
		{
			name:  "Theme",
			value: "bento",
			desc:  "UI color theme",
		},
	}
}

// Init initializes the settings
func (s Settings) Init() tea.Cmd {
	return nil
}

// Update handles settings messages
func (s Settings) Update(msg tea.Msg) (Settings, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if s.cursor > 0 {
				s.cursor--
			}
		case "down", "j":
			if s.cursor < len(s.items)-1 {
				s.cursor++
			}
		}
	}
	return s, nil
}

// View renders the settings
func (s Settings) View() string {
	title := styles.Title.Render("Settings")
	settingsView := s.renderSettings()
	note := styles.Subtle.Render("Note: Settings functionality coming in Phase 5")
	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		settingsView,
		note,
	)
}

// renderSettings renders the settings list
func (s Settings) renderSettings() string {
	var view string
	for i, item := range s.items {
		view += s.renderSetting(i, item)
	}
	return view
}

// renderSetting renders a single setting item
func (s Settings) renderSetting(index int, item settingItem) string {
	cursor := "  "
	nameStyle := styles.Normal
	if index == s.cursor {
		cursor = "> "
		nameStyle = styles.Selected
	}
	return cursor + nameStyle.Render(item.name) + "\n" +
		"  " + styles.Subtle.Render("  "+item.value) + "\n" +
		"  " + styles.Subtle.Render("  "+item.desc) + "\n\n"
}
