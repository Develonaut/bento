package guided_creation

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"bento/pkg/jubako"
	"bento/pkg/neta"
)

type guidedState int

const (
	guidedStateActive guidedState = iota
	guidedStateCompleted
	guidedStateCancelled
)

type guidedStage int

const (
	guidedStageMetadata guidedStage = iota
	guidedStageNodeTypeSelect
	guidedStageNodeParameters
	guidedStageContinue
	guidedStageGroupContext
	guidedStageEditMenu
	guidedStageNodeList
	guidedStageNodeEdit
)

// GuidedModal is a bubbletea-integrated modal for creating/editing bentos
// It shows a huh form on the left and a live preview of the bento on the right
type GuidedModal struct {
	state  guidedState
	stage  guidedStage
	lg     *lipgloss.Renderer
	styles *GuidedStyles
	form   *huh.Form
	width  int
	height int

	store     *jubako.Store
	workDir   string
	validator *neta.Validator

	// The bento being created/edited
	definition       *neta.Definition
	editing          bool
	originalFilename string // Original filename when editing (to preserve on save)

	// Current node being built
	currentNode *neta.Definition

	// Node being edited (for edit mode)
	editingNodeName string // Name of node being edited
	deletingNode    bool   // True if in delete mode, false if in edit mode

	// Temporary fields for edit forms (to hold pre-populated values)
	tempName         string
	tempDescription  string
	tempSelectedNode string // For node list selection

	// Temporary fields for node parameter editing
	tempNodeName    string
	tempNodeURL     string
	tempNodeMethod  string
	tempNodeHeaders string
	tempNodeBody    string
	tempNodePath    string
	tempNodeContent string
	tempNodeQuery   string

	// Group hierarchy tracking
	nodeStack     []*neta.Definition // Stack of parent nodes
	currentParent *neta.Definition   // Current parent being edited (nil = root)

	// Navigation and history tracking
	history      navigationHistory
	historyIndex int  // Current position in history (-1 = not navigating)
	navigating   bool // Are we in navigation mode?

	// Validation error to display
	validationErr error

	// Save state tracking
	saveInProgress bool // Prevent multiple concurrent saves
}

// NewGuidedModal creates a new guided creation modal
func NewGuidedModal(store *jubako.Store, workDir string, width, height int) *GuidedModal {
	m := &GuidedModal{
		state:        guidedStateActive,
		stage:        guidedStageMetadata,
		width:        min(width, guidedMaxWidth),
		height:       height,
		store:        store,
		workDir:      workDir,
		validator:    neta.NewValidator(),
		editing:      false,
		history:      newNavigationHistory(),
		historyIndex: -1,
		navigating:   false,
		definition: &neta.Definition{
			Version: "1.0",
			Type:    "group.sequence",
			Nodes:   []neta.Definition{},
			Edges:   []neta.NodeEdge{},
		},
	}

	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewGuidedStyles(m.lg)

	// Create the form with metadata fields
	m.form = m.createMetadataForm()

	return m
}

// NewGuidedModalForEdit creates a modal for editing an existing bento
func NewGuidedModalForEdit(store *jubako.Store, workDir string, width, height int, bentoName string) (*GuidedModal, error) {
	// Load the existing bento
	def, err := store.Load(bentoName)
	if err != nil {
		return nil, err
	}

	m := &GuidedModal{
		state:            guidedStateActive,
		stage:            guidedStageEditMenu, // Start at edit menu
		width:            min(width, guidedMaxWidth),
		height:           height,
		store:            store,
		workDir:          workDir,
		validator:        neta.NewValidator(),
		editing:          true,
		originalFilename: bentoName, // Store original filename
		history:          newNavigationHistory(),
		historyIndex:     -1,
		navigating:       false,
		definition:       &def,
	}

	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewGuidedStyles(m.lg)

	// Start with edit menu
	m.form = m.createEditMenuForm()

	return m, nil
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func (m *GuidedModal) Init() tea.Cmd {
	return m.form.Init()
}

func (m *GuidedModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = min(msg.Width, guidedMaxWidth) - m.styles.Base.GetHorizontalFrameSize()
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Interrupt
		case "esc":
			// Cancel the guided flow
			m.state = guidedStateCancelled
			return m, func() tea.Msg {
				return GuidedCompleteMsg{
					Cancelled: true,
				}
			}
		case "ctrl+up":
			// Navigate backward in history
			if cmd := m.navigateHistory(-1); cmd != nil {
				return m, cmd
			}
			return m, nil
		case "ctrl+down":
			// Navigate forward in history
			if cmd := m.navigateHistory(1); cmd != nil {
				return m, cmd
			}
			return m, nil
		case "ctrl+d":
			// Delete current node if we're on a node stage
			if m.stage == guidedStageNodeParameters && m.currentNode != nil {
				if cmd := m.deleteCurrentNode(); cmd != nil {
					return m, cmd
				}
			}
			return m, nil
		}
	}

	var cmds []tea.Cmd

	// Process the form
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	// Update definition with current form values based on stage
	if m.stage == guidedStageMetadata {
		m.updateDefinitionFromForm()
	} else if m.stage == guidedStageNodeParameters {
		// Node type is stored in currentNode.Type
		if m.currentNode != nil {
			m.updateCurrentNodeFromNodeForm(m.currentNode.Type)
		}
	}

	if m.form.State == huh.StateCompleted {
		// Form completed - handle stage transition
		newM, cmd := m.handleStageTransition()
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		return newM, tea.Batch(cmds...)
	}

	return m, tea.Batch(cmds...)
}
