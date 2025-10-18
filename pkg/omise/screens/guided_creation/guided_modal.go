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
	definition *neta.Definition
	editing    bool

	// Current node being built
	currentNode *neta.Definition

	// Validation error to display
	validationErr error
}

// NewGuidedModal creates a new guided creation modal
func NewGuidedModal(store *jubako.Store, workDir string, width, height int) *GuidedModal {
	m := &GuidedModal{
		state:     guidedStateActive,
		stage:     guidedStageMetadata,
		width:     min(width, guidedMaxWidth),
		height:    height,
		store:     store,
		workDir:   workDir,
		validator: neta.NewValidator(),
		editing:   false,
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
