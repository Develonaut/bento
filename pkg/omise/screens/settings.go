package screens

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
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
	selectingSlowMo bool
	themeForm       components.FormSelect
	slowMoForm      components.FormSelect
	availableThemes []styles.Variant
	config          config.Config
	dirPicker       components.DirPicker
	selectingDir    bool
	helpView        components.HelpView
	keys            components.SettingsKeyMap
	pickerKeys      components.PickerKeyMap
	selectedTheme   string
	selectedSlowMo  string
	statusMessage   string
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

// initializeSettings creates the initial Settings state
func initializeSettings(tm *styles.Manager, cfg config.Config, dp components.DirPicker) Settings {
	return Settings{
		themeManager:    tm,
		selectingTheme:  false,
		selectingSlowMo: false,
		selectingDir:    false,
		availableThemes: styles.AllVariants(),
		config:          cfg,
		dirPicker:       dp,
		helpView:        components.NewHelpView(),
		keys:            components.NewSettingsKeyMap(),
		pickerKeys:      components.NewPickerKeyMap(),
		selectedTheme:   string(tm.GetVariant()),
		selectedSlowMo:  formatSlowMoValue(cfg.SlowMoDelayMs),
		statusMessage:   "",
	}
}

// NewSettings creates a settings screen
func NewSettings() Settings {
	tm := styles.NewManager()
	cfg := config.Load()
	defaultCfg := config.Default()
	dp := components.NewDirPicker(cfg.SaveDirectory, defaultCfg.SaveDirectory)

	s := initializeSettings(tm, cfg, dp)
	items := s.buildSettings()
	s.list = components.NewStyledList(items, "")
	s.list.SetSize(80, 20)

	return s
}

// buildSettings creates setting items with theme manager
func (s *Settings) buildSettings() []list.Item {
	slowMoValue := formatSlowMoValue(s.config.SlowMoDelayMs)

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
			name:     "Slow-Mo Execution",
			value:    slowMoValue,
			desc:     "Slow down execution to watch node progress (press Enter/Space to cycle)",
			editable: true,
		},
		settingItem{
			name:     "Save Directory",
			value:    s.config.GetSaveDirectory(),
			desc:     "Directory for all app data (press Enter/Space to change)",
			editable: true,
		},
	}
}

// formatSlowMoValue formats the slow-mo delay value for display
func formatSlowMoValue(delayMs int) string {
	if delayMs == 0 {
		return "Off"
	}
	return fmt.Sprintf("%dms", delayMs)
}

// Init initializes the settings
func (s Settings) Init() tea.Cmd {
	return s.dirPicker.Init()
}

// InModalMode returns true if settings is in modal picker mode
func (s Settings) InModalMode() bool {
	return s.selectingDir || s.selectingTheme || s.selectingSlowMo
}

// Update handles settings messages
func (s Settings) Update(msg tea.Msg) (Settings, tea.Cmd) {
	if msg, ok := msg.(components.DirSelectedMsg); ok {
		return s.handleDirectorySelected(msg)
	}

	if newS, cmd, handled := s.routeModalModeUpdate(msg); handled {
		return newS, cmd
	}

	if _, ok := msg.(styles.ThemeChangedMsg); ok {
		return s.handleThemeChanged()
	}

	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		return s.handleWindowResize(msg)
	}

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

// KeyBindings returns the contextual key bindings for the footer
func (s Settings) KeyBindings() []key.Binding {
	// When in theme picker, directory picker, or slow-mo picker mode, don't show main settings keys
	if s.selectingTheme || s.selectingDir || s.selectingSlowMo {
		return []key.Binding{}
	}

	// Main settings keys
	return []key.Binding{
		s.keys.Select,
		s.keys.Reset,
	}
}
