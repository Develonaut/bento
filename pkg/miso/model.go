package miso

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// View states
const (
	listView = iota
	settingsView
	secretsView
	variablesView
	formView
	executionView
)

// Messages for async execution
type executionOutputMsg string
type executionCompleteMsg struct {
	err      error
	duration time.Duration
}
type executionStartMsg struct{}

// settingsFormType identifies which settings form is active
type settingsFormType int

const (
	noSettingsForm settingsFormType = iota
	bentoHomeForm
	themeForm
)

// Model holds the TUI state
type Model struct {
	currentView        int
	list               list.Model
	settingsList       list.Model
	secretsList        list.Model
	variablesList      list.Model
	form               *huh.Form
	selectedBento      string
	bentoVars          []Variable
	varHolders         map[string]*string // Pointers to form values
	logs               string
	logViewport        viewport.Model // Viewport for scrollable log display
	logChan            chan string    // Channel for streaming execution logs
	executing          bool
	width              int
	height             int
	theme              Variant
	quitting           bool
	activeSettingsForm settingsFormType // Tracks which settings form is active

	// Key bindings for each view
	listKeys      listKeyMap
	settingsKeys  settingsKeyMap
	secretsKeys   secretsKeyMap
	variablesKeys variablesKeyMap
	formKeys      formKeyMap
	executionKeys executionKeyMap
}

// BentoItem represents a bento in the list
type BentoItem struct {
	Name     string
	FilePath string
}

func (i BentoItem) Title() string       { return i.Name }
func (i BentoItem) Description() string { return i.FilePath }
func (i BentoItem) FilterValue() string { return i.Name }

// SettingsItem represents a settings option
type SettingsItem struct {
	Name   string
	Desc   string
	Action string
}

func (i SettingsItem) Title() string       { return i.Name }
func (i SettingsItem) Description() string { return i.Desc }
func (i SettingsItem) FilterValue() string { return i.Name }

// SecretItem represents a secret in the list
type SecretItem struct {
	Key string
}

func (i SecretItem) Title() string       { return i.Key }
func (i SecretItem) Description() string { return "Use {{SECRETS." + i.Key + "}} in bentos" }
func (i SecretItem) FilterValue() string { return i.Key }

// VariableItem represents a variable in the list
type VariableItem struct {
	Key   string
	Value string
}

func (i VariableItem) Title() string       { return i.Key }
func (i VariableItem) Description() string { return i.Value }
func (i VariableItem) FilterValue() string { return i.Key }

// NewTUI creates a new TUI model
func NewTUI() (*Model, error) {
	// Load bentos from ~/.bento/bentos
	items, err := loadBentos()
	if err != nil {
		return nil, fmt.Errorf("failed to load bentos: %w", err)
	}

	// Create bento list
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "🍱 Bentos"
	l.SetShowStatusBar(false)

	// Load current values for settings display
	currentHome := LoadBentoHome()
	currentTheme := LoadSavedTheme()

	// Create settings list
	settingsItems := []list.Item{
		SettingsItem{
			Name:   "Configure Bento Home",
			Desc:   fmt.Sprintf("Current: %s", currentHome),
			Action: "bentohome",
		},
		SettingsItem{
			Name:   "Manage Secrets",
			Desc:   "Add, view, or delete secrets",
			Action: "secrets",
		},
		SettingsItem{
			Name:   "Manage Variables",
			Desc:   "Add, view, or delete configuration variables",
			Action: "variables",
		},
		SettingsItem{
			Name:   "Change Theme",
			Desc:   fmt.Sprintf("Current: %s", currentTheme),
			Action: "theme",
		},
	}
	sl := list.New(settingsItems, list.NewDefaultDelegate(), 0, 0)
	sl.Title = "⚙️  Settings"
	sl.SetShowStatusBar(false)

	// Create empty secrets list (loaded on demand)
	secretsl := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	secretsl.Title = "🔐 Secrets"
	secretsl.SetShowStatusBar(false)

	// Create empty variables list (loaded on demand)
	variablesl := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	variablesl.Title = "📝 Variables"
	variablesl.SetShowStatusBar(false)

	return &Model{
		currentView:   listView,
		list:          l,
		settingsList:  sl,
		secretsList:   secretsl,
		variablesList: variablesl,
		theme:         VariantNasu, // Default theme

		// Initialize key bindings
		listKeys:      newListKeyMap(),
		settingsKeys:  newSettingsKeyMap(),
		secretsKeys:   newSecretsKeyMap(),
		variablesKeys: newVariablesKeyMap(),
		formKeys:      newFormKeyMap(),
		executionKeys: newExecutionKeyMap(),
	}, nil
}

// Init initializes the TUI
func (m Model) Init() tea.Cmd {
	return nil
}
