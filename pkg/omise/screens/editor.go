package screens

import (
	tea "github.com/charmbracelet/bubbletea"

	"bento/pkg/jubako"
	"bento/pkg/neta"
	"bento/pkg/omise/components"
	"bento/pkg/pantry"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
)

// EditorMode defines the editor mode
type EditorMode int

const (
	EditorModeCreate EditorMode = iota
	EditorModeEdit
)

// EditorState defines the current editor state
type EditorState int

const (
	StateNaming EditorState = iota
	StateSelectingType
	StateConfiguringNode
	StateReview
)

// ViewMode defines the view mode
type ViewMode int

const (
	ViewModeList   ViewMode = iota // List view (Phase 7)
	ViewModeVisual                 // Visual bento box (Phase 8)
)

// Editor screen for creating and editing bentos
type Editor struct {
	mode      EditorMode
	state     EditorState
	store     *jubako.Store
	registry  *pantry.Pantry
	validator *neta.Validator

	// Bento being edited
	bentoName string
	bentoPath string
	def       neta.Definition

	// Current node being configured
	currentNodeType string

	// Huh form integration
	currentForm *huh.Form
	formValues  map[string]interface{}

	// Visual navigation
	selectedNodeIndex int
	viewMode          ViewMode

	// UI state
	message  string
	width    int
	height   int
	helpView components.HelpView
	keys     components.EditorKeyMap
}

// NewEditorCreate creates editor in create mode
func NewEditorCreate(store *jubako.Store, registry *pantry.Pantry) Editor {
	editor := Editor{
		mode:              EditorModeCreate,
		state:             StateNaming,
		store:             store,
		registry:          registry,
		validator:         neta.NewValidator(),
		formValues:        make(map[string]interface{}),
		selectedNodeIndex: 0,
		viewMode:          ViewModeList,
		helpView:          components.NewHelpView(),
		keys:              components.NewEditorKeyMap(),
		def: neta.Definition{
			Version: neta.CurrentVersion,
			Nodes:   []neta.Definition{},
		},
	}
	// Create initial form for bento naming
	editor.currentForm = editor.createNameForm()
	return editor
}

// NewEditorEdit creates editor in edit mode
func NewEditorEdit(store *jubako.Store, registry *pantry.Pantry, name, path string) (Editor, error) {
	def, err := store.Load(name)
	if err != nil {
		return Editor{}, err
	}

	return Editor{
		mode:              EditorModeEdit,
		state:             StateReview,
		store:             store,
		registry:          registry,
		validator:         neta.NewValidator(),
		formValues:        make(map[string]interface{}),
		selectedNodeIndex: 0,
		viewMode:          ViewModeList,
		helpView:          components.NewHelpView(),
		keys:              components.NewEditorKeyMap(),
		bentoName:         name,
		bentoPath:         path,
		def:               def,
	}, nil
}

// Init initializes the editor
func (e Editor) Init() tea.Cmd {
	// Form is already created in constructor if needed
	if e.currentForm != nil {
		return e.currentForm.Init()
	}
	return nil
}

// Update handles editor messages
func (e Editor) Update(msg tea.Msg) (Editor, tea.Cmd) {
	// Handle window resize
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		return e.handleResize(msg)
	}

	// Handle editor-specific messages before form processing
	// This allows tests and other code to send messages directly
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// If in modal mode, delegate to form first
		if e.currentForm != nil && e.InModalMode() {
			return e.updateForm(msg)
		}
		return e.handleKey(msg)
	case BentoNameEnteredMsg:
		return e.handleNameEntered(msg)
	case NodeTypeSelectedMsg:
		return e.handleTypeSelected(msg)
	case NodeConfiguredMsg:
		return e.handleNodeConfigured(msg)
	}

	// If we have an active form and message wasn't handled above, delegate to it
	if e.currentForm != nil && e.InModalMode() {
		return e.updateForm(msg)
	}

	return e, nil
}

// InModalMode returns true if editor is showing a form
// This prevents tab navigation to other screens during form input
func (e Editor) InModalMode() bool {
	// Modal mode is active during form input states
	return e.state == StateNaming ||
		e.state == StateSelectingType ||
		e.state == StateConfiguringNode
}

// GetBentoName returns the bento name
func (e Editor) GetBentoName() string {
	return e.bentoName
}

// GetDefinition returns the bento definition
func (e Editor) GetDefinition() neta.Definition {
	return e.def
}

// KeyBindings returns the contextual key bindings for the footer
func (e Editor) KeyBindings() []key.Binding {
	// When in modal mode (form), don't show editor keys
	if e.InModalMode() {
		return []key.Binding{}
	}

	// In review state, show main editor keys
	if e.state == StateReview {
		if len(e.def.Nodes) > 0 {
			return []key.Binding{
				e.keys.Add,
				e.keys.Edit,
				e.keys.Delete,
				e.keys.Save,
			}
		}
		// No nodes yet, just show add and save
		return []key.Binding{
			e.keys.Add,
			e.keys.Save,
		}
	}

	return []key.Binding{}
}
