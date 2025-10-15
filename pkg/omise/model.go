package omise

import (
	tea "github.com/charmbracelet/bubbletea"

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
	screenCount // Editor is excluded from tab cycle

	// Modal screens (not in tab cycle)
	ScreenEditor
)

// String returns the screen name
func (s Screen) String() string {
	return [...]string{"Browser", "Executor", "Pantry", "Settings", "Help", "Editor"}[s]
}

// Model is the root Bubble Tea model for Omise
type Model struct {
	screen Screen
	width  int
	height int

	// Screen models
	browser  screens.Browser
	executor screens.Executor
	pantry   screens.Pantry
	settings screens.Settings
	help     screens.Help
	editor   screens.Editor

	// Application state
	quitting bool
	workDir  string
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

	return Model{
		screen:   ScreenBrowser,
		browser:  browser,
		executor: screens.NewExecutor(),
		pantry:   screens.NewPantry(),
		settings: screens.NewSettings(),
		help:     screens.NewHelp(),
		editor:   screens.Editor{},
		workDir:  workDir,
	}
}

// NewModelWithWorkDir creates model with configured work directory
func NewModelWithWorkDir(workDir string) (Model, error) {
	browser, err := screens.NewBrowser(workDir)
	if err != nil {
		return Model{}, err
	}

	return Model{
		screen:   ScreenBrowser,
		browser:  browser,
		executor: screens.NewExecutor(),
		pantry:   screens.NewPantry(),
		settings: screens.NewSettings(),
		help:     screens.NewHelp(),
		editor:   screens.Editor{},
		workDir:  workDir,
	}, nil
}

// getDefaultWorkDir returns the default bento work directory
func getDefaultWorkDir() string {
	// Default to current directory if home unavailable
	return "."
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
