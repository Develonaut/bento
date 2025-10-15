package screens

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"

	"bento/pkg/jubako"
	"bento/pkg/neta"
	"bento/pkg/pantry"
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
	ctx       context.Context

	// Bento being edited
	bentoName string
	bentoPath string
	def       neta.Definition

	// Current node being configured
	currentNodeType string

	// Visual navigation
	selectedNodeIndex int
	viewMode          ViewMode

	// UI state
	message string
	width   int
	height  int
}

// NewEditorCreate creates editor in create mode
func NewEditorCreate(store *jubako.Store, registry *pantry.Pantry) Editor {
	return Editor{
		mode:              EditorModeCreate,
		state:             StateNaming,
		store:             store,
		registry:          registry,
		validator:         neta.NewValidator(),
		ctx:               context.Background(),
		selectedNodeIndex: 0,
		viewMode:          ViewModeList,
		def: neta.Definition{
			Version: neta.CurrentVersion,
			Nodes:   []neta.Definition{},
		},
	}
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
		ctx:               context.Background(),
		selectedNodeIndex: 0,
		viewMode:          ViewModeList,
		bentoName:         name,
		bentoPath:         path,
		def:               def,
	}, nil
}

// Init initializes the editor
func (e Editor) Init() tea.Cmd {
	return nil
}

// Update handles editor messages
func (e Editor) Update(msg tea.Msg) (Editor, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return e.handleResize(msg)
	case tea.KeyMsg:
		return e.handleKey(msg)
	case BentoNameEnteredMsg:
		return e.handleNameEntered(msg)
	case NodeTypeSelectedMsg:
		return e.handleTypeSelected(msg)
	case NodeConfiguredMsg:
		return e.handleNodeConfigured(msg)
	}

	return e, nil
}

// GetBentoName returns the bento name
func (e Editor) GetBentoName() string {
	return e.bentoName
}

// GetDefinition returns the bento definition
func (e Editor) GetDefinition() neta.Definition {
	return e.def
}
