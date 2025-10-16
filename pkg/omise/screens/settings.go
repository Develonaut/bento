package screens

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"bento/pkg/omise/components"
	"bento/pkg/omise/config"
	"bento/pkg/omise/styles"
)

// Settings shows configuration options
type Settings struct {
	list            components.StyledList
	themeManager    *styles.Manager
	selectingTheme  bool
	themeCursor     int
	availableThemes []styles.Variant
	config          config.Config
	dirPicker       components.DirPicker
	selectingDir    bool
	helpView        components.HelpView
	keys            components.SettingsKeyMap
	pickerKeys      components.PickerKeyMap
}

type settingItem struct {
	name       string
	value      string
	desc       string
	editable   bool
	valueStyle lipgloss.Style // Optional custom style
}

// Implement list.Item interface
func (s settingItem) FilterValue() string { return s.name }
func (s settingItem) Title() string {
	editIndicator := ""
	if s.editable {
		editIndicator = " [↵]"
	}
	return s.name + editIndicator
}
func (s settingItem) Description() string {
	valueStyle := styles.Subtle
	if s.valueStyle.GetBold() || s.valueStyle.GetForeground() != lipgloss.Color("") {
		valueStyle = s.valueStyle
	}
	return valueStyle.Render(s.value) + " • " + s.desc
}

// NewSettings creates a settings screen
func NewSettings() Settings {
	tm := styles.NewManager()
	cfg := config.Load()
	defaultCfg := config.Default()

	// Create DirPicker: start at current directory, but reset goes to app default
	dp := components.NewDirPicker(cfg.SaveDirectory, defaultCfg.SaveDirectory)

	s := Settings{
		themeManager:    tm,
		selectingTheme:  false,
		selectingDir:    false,
		themeCursor:     0,
		availableThemes: styles.AllVariants(),
		config:          cfg,
		dirPicker:       dp,
		helpView:        components.NewHelpView(),
		keys:            components.NewSettingsKeyMap(),
		pickerKeys:      components.NewPickerKeyMap(),
	}

	items := s.buildSettings()
	s.list = components.NewStyledList(items, "")
	s.list.SetSize(80, 20) // Set default size, will be updated on window resize

	return s
}

// buildSettings creates setting items with theme manager
func (s *Settings) buildSettings() []list.Item {
	return []list.Item{
		settingItem{
			name:     "Theme",
			value:    string(s.themeManager.GetVariant()),
			desc:     "Sushi-themed color variant (press Enter/Space to select)",
			editable: true,
			valueStyle: lipgloss.NewStyle().
				Foreground(styles.Primary).
				Bold(true),
		},
		settingItem{
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

	// Handle window resize
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		h, v := lipgloss.NewStyle().Margin(2, 2).GetFrameSize()
		s.list.SetSize(msg.Width-h, msg.Height-v-4)
		return s, nil
	}

	// Handle keyboard input
	if msg, ok := msg.(tea.KeyMsg); ok {
		return s.handleKeyInput(msg)
	}

	return s, nil
}

// handleThemeChanged rebuilds items and styles when theme changes
func (s Settings) handleThemeChanged() (Settings, tea.Cmd) {
	items := s.buildSettings()
	s.list.SetItems(items)
	s.list = s.list.RebuildStyles()
	s.dirPicker = s.dirPicker.RebuildStyles()
	return s, nil
}

// handleKeyInput processes keyboard input
func (s Settings) handleKeyInput(msg tea.KeyMsg) (Settings, tea.Cmd) {
	// Handle theme selection mode
	if s.selectingTheme {
		return s.handleThemeSelection(msg)
	}

	// Handle setting-specific actions
	switch msg.String() {
	case "enter", " ":
		return s.activateSetting()
	case "r":
		return s.resetCurrentSetting()
	}

	// Delegate navigation to list
	var cmd tea.Cmd
	s.list.Model, cmd = s.list.Model.Update(msg)
	return s, cmd
}

// activateSetting activates the current setting (e.g., opens theme selector)
func (s Settings) activateSetting() (Settings, tea.Cmd) {
	selected := s.list.SelectedItem()
	if selected == nil {
		return s, nil
	}

	item, ok := selected.(settingItem)
	if !ok || !item.editable {
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
	selected := s.list.SelectedItem()
	if selected == nil {
		return s, nil
	}

	item, ok := selected.(settingItem)
	if !ok || !item.editable {
		return s, nil
	}

	// Check if this is the theme setting
	if item.name == "Theme" {
		// Reset to default theme (Maguro)
		s.themeManager.SetVariant(styles.VariantMaguro)
		items := s.buildSettings()
		s.list.SetItems(items)
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

// ContextualKeys returns the most important contextual keys for the footer
func (s Settings) ContextualKeys() []components.KeyHelp {
	// When in theme picker or directory picker mode, don't show main settings keys
	if s.selectingTheme || s.selectingDir {
		return []components.KeyHelp{}
	}

	// Main settings keys
	return []components.KeyHelp{
		{Key: "enter", Desc: "edit"},
		{Key: "r", Desc: "reset"},
	}
}
