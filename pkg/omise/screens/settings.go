package screens

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"bento/pkg/omise/components"
	"bento/pkg/omise/config"
	"bento/pkg/omise/styles"
)

// Settings shows configuration options
type Settings struct {
	cursor          int
	items           []settingItem
	themeManager    *styles.Manager
	selectingTheme  bool
	themeCursor     int
	availableThemes []styles.Variant
	config          config.Config
	dirPicker       components.DirPicker
	selectingDir    bool
}

type settingItem struct {
	name       string
	value      string
	desc       string
	editable   bool
	valueStyle lipgloss.Style // Optional custom style
}

// NewSettings creates a settings screen
func NewSettings() Settings {
	tm := styles.NewManager()
	cfg := config.Load()
	defaultCfg := config.Default()

	// Create DirPicker: start at current directory, but reset goes to app default
	dp := components.NewDirPicker(cfg.SaveDirectory, defaultCfg.SaveDirectory)

	s := Settings{
		cursor:          0,
		themeManager:    tm,
		selectingTheme:  false,
		selectingDir:    false,
		themeCursor:     0,
		availableThemes: styles.AllVariants(),
		config:          cfg,
		dirPicker:       dp,
	}
	s.items = s.buildSettings()
	return s
}

// buildSettings creates setting items with theme manager
func (s *Settings) buildSettings() []settingItem {
	return []settingItem{
		{
			name:     "Theme",
			value:    string(s.themeManager.GetVariant()),
			desc:     "Sushi-themed color variant (press Enter/Space to select)",
			editable: true,
			valueStyle: lipgloss.NewStyle().
				Foreground(styles.Primary).
				Bold(true),
		},
		{
			name:     "Save Directory",
			value:    s.config.GetSaveDirectory(),
			desc:     "Directory for all app data (press Enter/Space to change)",
			editable: true,
		},
	}
}

// Init initializes the settings
func (s Settings) Init() tea.Cmd {
	return s.dirPicker.Init()
}

// InModalMode returns true if settings is in modal picker mode
func (s Settings) InModalMode() bool {
	return s.selectingDir || s.selectingTheme
}

// Update handles settings messages
func (s Settings) Update(msg tea.Msg) (Settings, tea.Cmd) {
	// Handle directory selection
	if msg, ok := msg.(components.DirSelectedMsg); ok {
		return s.handleDirectorySelected(msg)
	}

	// Handle directory picker mode
	if s.selectingDir {
		return s.handleDirectoryPickerMode(msg)
	}

	// Handle theme changes and rebuild styles
	if _, ok := msg.(styles.ThemeChangedMsg); ok {
		return s.handleThemeChanged()
	}

	// Handle keyboard input
	if msg, ok := msg.(tea.KeyMsg); ok {
		return s.handleKeyInput(msg)
	}

	return s, nil
}

// handleThemeChanged rebuilds items and styles when theme changes
func (s Settings) handleThemeChanged() (Settings, tea.Cmd) {
	s.items = s.buildSettings()
	s.dirPicker = s.dirPicker.RebuildStyles()
	return s, nil
}

// handleKeyInput processes keyboard input
func (s Settings) handleKeyInput(msg tea.KeyMsg) (Settings, tea.Cmd) {
	// Handle theme selection mode
	if s.selectingTheme {
		return s.handleThemeSelection(msg)
	}

	// Normal settings navigation
	switch msg.String() {
	case "up", "k":
		return s.moveCursorUp(), nil
	case "down", "j":
		return s.moveCursorDown(), nil
	case "enter", " ":
		return s.activateSetting()
	case "r":
		return s.resetCurrentSetting()
	}

	return s, nil
}

// moveCursorUp moves cursor to previous item
func (s Settings) moveCursorUp() Settings {
	if s.cursor > 0 {
		s.cursor--
	}
	return s
}

// moveCursorDown moves cursor to next item
func (s Settings) moveCursorDown() Settings {
	if s.cursor < len(s.items)-1 {
		s.cursor++
	}
	return s
}

// activateSetting activates the current setting (e.g., opens theme selector)
func (s Settings) activateSetting() (Settings, tea.Cmd) {
	item := s.items[s.cursor]
	if !item.editable {
		return s, nil
	}

	// Check if this is the theme setting
	if item.name == "Theme" {
		s = s.activateThemeSetting()
		return s, nil
	}

	// Check if this is the save directory setting
	if item.name == "Save Directory" {
		s.selectingDir = true
		return s, nil
	}

	return s, nil
}

// resetCurrentSetting resets the current setting to its default value
func (s Settings) resetCurrentSetting() (Settings, tea.Cmd) {
	item := s.items[s.cursor]
	if !item.editable {
		return s, nil
	}

	// Check if this is the theme setting
	if item.name == "Theme" {
		// Reset to default theme (Maguro)
		s.themeManager.SetVariant(styles.VariantMaguro)
		s.items = s.buildSettings()
		return s, func() tea.Msg {
			return styles.ThemeChangedMsg{}
		}
	}

	// Check if this is the save directory setting
	if item.name == "Save Directory" {
		return s.resetDirectorySetting()
	}

	return s, nil
}

// View renders the settings
func (s Settings) View() string {
	title := styles.Title.Render("Settings")

	// Show directory picker if in directory selection mode
	if s.selectingDir {
		return s.renderDirectoryPickerView(title)
	}

	// Show theme selector if in selection mode
	if s.selectingTheme {
		return s.renderThemeSelector(title)
	}

	// Show normal settings list
	return s.renderSettingsListView(title)
}

// renderSettingsListView renders the normal settings list
func (s Settings) renderSettingsListView(title string) string {
	settingsView := s.renderSettings()
	hint := styles.Subtle.Render("↑/↓: Navigate • Enter/Space: Select • r: Reset to Default • Esc: Back")
	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		settingsView,
		"",
		hint,
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

	valueStyle := styles.Subtle
	// Use custom valueStyle if it has been set
	if item.valueStyle.GetBold() || item.valueStyle.GetForeground() != lipgloss.Color("") {
		valueStyle = item.valueStyle
	}

	editIndicator := ""
	if item.editable {
		editIndicator = " [↵]"
	}

	return cursor + nameStyle.Render(item.name) + "\n" +
		"  " + valueStyle.Render("  "+item.value+editIndicator) + "\n" +
		"  " + styles.Subtle.Render("  "+item.desc) + "\n\n"
}
