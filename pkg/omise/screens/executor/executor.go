package executor

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"bento/pkg/omise/components"
	"bento/pkg/omise/emoji"

	"github.com/charmbracelet/bubbles/key"
)

// NodeStatus represents node execution state
type NodeStatus int

const (
	NodePending NodeStatus = iota
	NodeRunning
	NodeCompleted
	NodeFailed
)

// NodeState tracks individual node execution
type NodeState struct {
	path      string
	name      string
	nodeType  string
	status    NodeStatus
	startTime time.Time
	duration  time.Duration
	depth     int // Nesting level for indentation
}

// Executor shows bento execution progress
type Executor struct {
	spinner          components.Spinner
	progress         components.Progress
	sequence         components.Sequence
	progressPercent  float64
	status           string
	running          bool
	complete         bool
	success          bool
	bentoName        string
	bentoPath        string
	workDir          string
	errorMsg         string
	result           string
	copyFeedback     string
	startTime        time.Time
	endTime          time.Time
	lifecycleHistory []string
	nodeStates       []NodeState
}

// NewExecutor creates an executor screen
func NewExecutor() Executor {
	return Executor{
		spinner:  components.NewSpinner(),
		progress: components.NewProgress(40),
		sequence: components.NewSequence(),
		status:   "Ready to execute bentos",
		running:  false,
	}
}

// Init initializes the executor
func (e Executor) Init() tea.Cmd {
	return nil
}

// Update handles executor messages
func (e Executor) Update(msg tea.Msg) (Executor, tea.Cmd) {
	// Handle theme changes and keyboard input
	if updated, cmd, handled := e.handleThemeAndInput(msg); handled {
		return updated, cmd
	}

	// Handle execution-related messages
	return e.handleExecutionMessages(msg)
}

// StartBento prepares the executor to run a bento
func (e Executor) StartBento(name, path, workDir string) Executor {
	e.bentoName = name
	e.bentoPath = path
	e.workDir = workDir
	e.running = true
	e.complete = false
	e.success = false
	e.errorMsg = ""
	e.result = ""
	e.copyFeedback = ""
	e.status = "Adding neta..."
	e.progressPercent = 0.0
	e.startTime = time.Now()
	e.endTime = time.Time{} // Zero value
	e.lifecycleHistory = []string{emoji.Bento + " Preparing Bento..."}
	e.nodeStates = []NodeState{} // Reset node states from previous run
	return e
}

// ExecuteCmd returns the command to start execution
func (e Executor) ExecuteCmd(program *tea.Program) tea.Cmd {
	return tea.Batch(
		e.spinner.Tick,
		ExecuteBentoCmd(e.bentoName, e.workDir, program),
	)
}

// copyToClipboard copies the result to clipboard
func (e Executor) copyToClipboard() tea.Cmd {
	return func() tea.Msg {
		return CopyResultCmd(e.result, e.bentoName, e.errorMsg, e.success)
	}
}

// copyEntireView copies the entire rendered view to clipboard
func (e Executor) copyEntireView() tea.Cmd {
	return func() tea.Msg {
		return CopyEntireViewCmd(e.View())
	}
}

// IsRunning returns whether the executor is currently running
func (e Executor) IsRunning() bool {
	return e.running
}

// KeyBindings returns the contextual key bindings for the footer
func (e Executor) KeyBindings() []key.Binding {
	// When execution is complete, show copy and rerun keys
	if e.complete {
		keys := []key.Binding{
			key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "rerun"),
			),
		}
		// Add copy key if we have content to copy
		if (e.success && e.result != "") || (!e.success && e.errorMsg != "") {
			keys = append(keys, key.NewBinding(
				key.WithKeys("c"),
				key.WithHelp("c", "copy output"),
			))
		}
		// Always show shift+c for debugging
		keys = append(keys, key.NewBinding(
			key.WithKeys("C"),
			key.WithHelp("shift+c", "copy view"),
		))
		return keys
	}
	// No contextual keys during execution or when no content
	return []key.Binding{}
}
