package omise

import (
	"path/filepath"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"bento/pkg/omise/components"
	"bento/pkg/omise/config"
	"bento/pkg/omise/screens"
)

// Screen identifies which screen is active
type Screen int

const (
	ScreenBrowser Screen = iota
	ScreenExecutor
	ScreenPantry
	ScreenSettings
	ScreenHelp
	screenCount // Marker for tab cycle end

	// Modal screens (not in tab cycle)
	ScreenEditor // TODO: Editor not yet implemented
)

// String returns the screen name
func (s Screen) String() string {
	names := [...]string{"Browser", "Executor", "Pantry", "Settings", "Help", "screenCount", "Editor"}
	if int(s) >= len(names) {
		return "Unknown"
	}
	return names[s]
}

// Model is the root Bubble Tea model for Omise
type Model struct {
	screen   Screen
	width    int
	height   int
	viewport viewport.Model
	tabView  components.TabView

	// Screen models
	browser  screens.Browser
	executor screens.Executor
	pantry   screens.Pantry
	settings screens.Settings
	help     screens.Help

	// Application state
	quitting bool
	workDir  string
	program  *tea.Program // For executor messaging
}

// NewModel creates the initial application model
func NewModel() Model {
	// Use default work directory
	workDir := getDefaultWorkDir()
	browser, err := screens.NewBrowser(workDir)
	if err != nil {
		// Fall back to empty browser on error
		browser = screens.Browser{}
	}

	// Create viewport with default size (will be updated on first resize)
	vp := viewport.New(80, 20)

	return Model{
		screen:   ScreenBrowser,
		viewport: vp,
		tabView:  components.NewTabView(),
		browser:  browser,
		executor: screens.NewExecutor(),
		pantry:   screens.NewPantry(),
		settings: screens.NewSettings(),
		help:     screens.NewHelp(),
		workDir:  workDir,
	}
}

// NewModelWithWorkDir creates model with configured work directory
func NewModelWithWorkDir(workDir string) (Model, error) {
	browser, err := screens.NewBrowser(workDir)
	if err != nil {
		return Model{}, err
	}

	// Create viewport with default size (will be updated on first resize)
	vp := viewport.New(80, 20)

	return Model{
		screen:   ScreenBrowser,
		viewport: vp,
		tabView:  components.NewTabView(),
		browser:  browser,
		executor: screens.NewExecutor(),
		pantry:   screens.NewPantry(),
		settings: screens.NewSettings(),
		help:     screens.NewHelp(),
		workDir:  workDir,
	}, nil
}

// getDefaultWorkDir returns the default bento work directory
func getDefaultWorkDir() string {
	cfg := config.Load()
	// Append /bentos subdirectory to the save directory
	return filepath.Join(cfg.SaveDirectory, "bentos")
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// NextScreen cycles to the next screen
func (m Model) NextScreen() Screen {
	return (m.screen + 1) % screenCount
}

// PrevScreen cycles to the previous screen
func (m Model) PrevScreen() Screen {
	if m.screen == 0 {
		return screenCount - 1
	}
	return m.screen - 1
}

// SetProgram stores program reference for messaging
func (m *Model) SetProgram(p *tea.Program) {
	m.program = p
}

// ScreenToTab maps a screen to its corresponding tab
func (m Model) ScreenToTab() components.TabID {
	switch m.screen {
	case ScreenBrowser:
		return components.TabBentos
	case ScreenPantry:
		return components.TabRecipes
	case ScreenSettings:
		return components.TabMise
	case ScreenHelp:
		return components.TabSensei
	default:
		return components.TabBentos
	}
}

// TabToScreen maps a tab to its corresponding screen
func (m Model) TabToScreen(tab components.TabID) Screen {
	switch tab {
	case components.TabBentos:
		return ScreenBrowser
	case components.TabRecipes:
		return ScreenPantry
	case components.TabMise:
		return ScreenSettings
	case components.TabSensei:
		return ScreenHelp
	default:
		return ScreenBrowser
	}
}

// SwitchToTab switches to the screen for the given tab
func (m Model) SwitchToTab(tab components.TabID) Model {
	m.screen = m.TabToScreen(tab)
	m.tabView = m.tabView.SetActiveTab(tab)
	return m
}
