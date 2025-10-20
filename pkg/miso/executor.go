package miso

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Develonaut/bento/pkg/neta"
)

// NodeStatus represents node execution state.
type NodeStatus int

const (
	NodePending NodeStatus = iota
	NodeRunning
	NodeCompleted
	NodeFailed
)

// NodeState tracks individual node execution.
type NodeState struct {
	path      string
	name      string
	nodeType  string
	status    NodeStatus
	startTime time.Time
	duration  time.Duration
	depth     int // Nesting level for indentation
}

// Executor displays bento execution progress using Bubbletea.
// This is a lightweight executor view that shows real-time progress
// and exits automatically when execution completes.
type Executor struct {
	theme      *Theme
	palette    Palette
	sequence   *Sequence
	nodeStates []NodeState
	bentoName  string
	running    bool
	complete   bool
	success    bool
	errorMsg   string
	spinner    Spinner
}

// Message types for Bubbletea

// NodeStartedMsg signals that a node has started execution.
type NodeStartedMsg struct {
	Path     string
	Name     string
	NodeType string
}

// NodeCompletedMsg signals that a node has finished execution.
type NodeCompletedMsg struct {
	Path     string
	Duration time.Duration
	Error    error
}

// ExecutionInitMsg initializes the executor with bento definition.
type ExecutionInitMsg struct {
	Definition *neta.Definition
}

// ExecutionCompleteMsg signals that bento execution is complete.
type ExecutionCompleteMsg struct {
	Success bool
	Error   error
}

// NewExecutor creates an executor for the given bento definition.
func NewExecutor(def *neta.Definition, theme *Theme, palette Palette) Executor {
	sequence := NewSequenceWithTheme(theme, palette)
	spinner := NewSpinner(palette)

	return Executor{
		theme:      theme,
		palette:    palette,
		sequence:   sequence,
		nodeStates: []NodeState{},
		bentoName:  def.Name,
		running:    true,
		complete:   false,
		success:    false,
		spinner:    spinner,
	}
}

// Init initializes the Bubbletea model.
func (e Executor) Init() tea.Cmd {
	return tea.Batch(
		e.spinner.Model.Tick,
	)
}

// Update handles Bubbletea messages.
func (e Executor) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Allow quit with Ctrl+C
		if msg.Type == tea.KeyCtrlC {
			return e, tea.Quit
		}

	case ExecutionInitMsg:
		// Flatten definition to get all nodes
		e.nodeStates = flattenDefinition(*msg.Definition, "")
		e.updateSequence()
		return e, nil

	case NodeStartedMsg:
		e.handleNodeStarted(msg)
		e.updateSequence()
		return e, nil

	case NodeCompletedMsg:
		e.handleNodeCompleted(msg)
		e.updateSequence()
		return e, nil

	case ExecutionCompleteMsg:
		e.complete = true
		e.running = false
		e.success = msg.Success
		if msg.Error != nil {
			e.errorMsg = msg.Error.Error()
		}
		// Exit after showing completion message briefly
		return e, tea.Sequence(
			tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
				return tea.Quit()
			}),
		)

	case spinner.TickMsg:
		// Update spinner animation
		var cmd tea.Cmd
		e.spinner, cmd = e.spinner.Update(msg)
		e.sequence.UpdateSpinner(e.spinner)
		return e, cmd
	}

	return e, nil
}

// handleNodeStarted updates node to running state.
func (e *Executor) handleNodeStarted(msg NodeStartedMsg) {
	for i := range e.nodeStates {
		if e.nodeStates[i].path == msg.Path {
			e.nodeStates[i].status = NodeRunning
			e.nodeStates[i].startTime = time.Now()
			return
		}
	}
}

// handleNodeCompleted updates node to completed/failed state.
func (e *Executor) handleNodeCompleted(msg NodeCompletedMsg) {
	for i := range e.nodeStates {
		if e.nodeStates[i].path == msg.Path {
			e.nodeStates[i].duration = msg.Duration
			if msg.Error != nil {
				e.nodeStates[i].status = NodeFailed
			} else {
				e.nodeStates[i].status = NodeCompleted
			}
			return
		}
	}
}

// updateSequence converts node states to sequence steps.
func (e *Executor) updateSequence() {
	steps := make([]Step, len(e.nodeStates))
	for i, node := range e.nodeStates {
		steps[i] = Step{
			Name:     node.name,
			Type:     node.nodeType,
			Status:   convertNodeStatusToStepStatus(node.status),
			Duration: node.duration,
			Depth:    node.depth,
		}
	}
	e.sequence.SetSteps(steps)
}

// convertNodeStatusToStepStatus converts NodeStatus to StepStatus.
func convertNodeStatusToStepStatus(status NodeStatus) StepStatus {
	switch status {
	case NodePending:
		return StepPending
	case NodeRunning:
		return StepRunning
	case NodeCompleted:
		return StepCompleted
	case NodeFailed:
		return StepFailed
	default:
		return StepPending
	}
}

// Success returns whether execution was successful.
func (e Executor) Success() bool {
	return e.success
}
