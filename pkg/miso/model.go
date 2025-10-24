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
	themeView
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
	variableForm
)

// Model holds the TUI state
type Model struct {
	currentView        int
	list               list.Model
	settingsList       list.Model
	secretsList        list.Model
	variablesList      list.Model
	themeList          list.Model
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
	reorderMode        bool             // Tracks if we're in bento reorder mode

	// Key bindings for each view
	listKeys      listKeyMap
	settingsKeys  settingsKeyMap
	secretsKeys   secretsKeyMap
	variablesKeys variablesKeyMap
	formKeys      formKeyMap
	executionKeys executionKeyMap
	themeKeys     themeKeyMap
}

// BentoItem represents a bento in the list
type BentoItem struct {
	Name     string
	FilePath string
}

func (i BentoItem) Title() string       { return i.Name }
func (i BentoItem) Description() string { return CompressPath(i.FilePath) }
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

func (i VariableItem) Title() string { return i.Key }
func (i VariableItem) Description() string {
	// Show compressed path if it looks like a path
	if len(i.Value) > 0 && (i.Value[0] == '/' || i.Value[0] == '~' || (len(i.Value) > 1 && i.Value[1] == ':')) {
		return CompressPath(i.Value)
	}
	return i.Value
}
func (i VariableItem) FilterValue() string { return i.Key }

// ThemeItem represents a theme in the list
type ThemeItem struct {
	Variant     Variant
	DisplayName string
	Desc        string
}

func (i ThemeItem) Title() string       { return i.DisplayName }
func (i ThemeItem) Description() string { return i.Desc }
func (i ThemeItem) FilterValue() string { return i.DisplayName }

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
	l.SetShowHelp(false) // Disable built-in help - we provide our own

	// Load current values for settings display
	currentHome := LoadBentoHome()
	currentTheme := LoadSavedTheme()

	// Create settings list
	settingsItems := []list.Item{
		SettingsItem{
			Name:   "Bento Home",
			Desc:   fmt.Sprintf("Current: %s", CompressPath(currentHome)),
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
	sl.SetShowHelp(false) // Disable built-in help - we provide our own

	// Create empty secrets list (loaded on demand)
	secretsl := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	secretsl.Title = "🔐 Secrets"
	secretsl.SetShowStatusBar(false)
	secretsl.SetShowHelp(false) // Disable built-in help - we provide our own

	// Create empty variables list (loaded on demand)
	variablesl := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	variablesl.Title = "📝 Variables"
	variablesl.SetShowStatusBar(false)
	variablesl.SetShowHelp(false) // Disable built-in help - we provide our own

	// Create theme list with all available themes
	themeItems := []list.Item{
		ThemeItem{
			Variant:     VariantNasu,
			DisplayName: "Nasu",
			Desc:        "Purple - eggplant sushi",
		},
		ThemeItem{
			Variant:     VariantWasabi,
			DisplayName: "Wasabi",
			Desc:        "Green - wasabi",
		},
		ThemeItem{
			Variant:     VariantToro,
			DisplayName: "Toro",
			Desc:        "Pink - fatty tuna",
		},
		ThemeItem{
			Variant:     VariantTamago,
			DisplayName: "Tamago",
			Desc:        "Yellow - egg sushi",
		},
		ThemeItem{
			Variant:     VariantTonkotsu,
			DisplayName: "Tonkotsu",
			Desc:        "Red - pork bone broth",
		},
		ThemeItem{
			Variant:     VariantSaba,
			DisplayName: "Saba",
			Desc:        "Cyan - mackerel",
		},
		ThemeItem{
			Variant:     VariantIka,
			DisplayName: "Ika",
			Desc:        "White - squid",
		},
	}
	themel := list.New(themeItems, list.NewDefaultDelegate(), 0, 0)
	themel.Title = "🎨 Themes"
	themel.SetShowStatusBar(false)
	themel.SetShowHelp(false) // Disable built-in help - we provide our own

	return &Model{
		currentView:   listView,
		list:          l,
		settingsList:  sl,
		secretsList:   secretsl,
		variablesList: variablesl,
		themeList:     themel,
		theme:         VariantNasu, // Default theme

		// Initialize key bindings
		listKeys:      newListKeyMap(),
		settingsKeys:  newSettingsKeyMap(),
		secretsKeys:   newSecretsKeyMap(),
		variablesKeys: newVariablesKeyMap(),
		formKeys:      newFormKeyMap(),
		executionKeys: newExecutionKeyMap(),
		themeKeys:     newThemeKeyMap(),
	}, nil
}

// Init initializes the TUI
func (m Model) Init() tea.Cmd {
	return nil
}
