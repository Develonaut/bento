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
	screenCount
)

// String returns the screen name
func (s Screen) String() string {
	return [...]string{"Browser", "Executor", "Pantry", "Settings", "Help"}[s]
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

	// Application state
	quitting bool
}

// NewModel creates the initial application model
func NewModel() Model {
	return Model{
		screen:   ScreenBrowser,
		browser:  screens.NewBrowser(),
		executor: screens.NewExecutor(),
		pantry:   screens.NewPantry(),
		settings: screens.NewSettings(),
		help:     screens.NewHelp(),
	}
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
