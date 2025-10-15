package screens

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"bento/pkg/omise/styles"
)

// Executor shows workflow execution progress
type Executor struct {
	spinner      spinner.Model
	progress     progress.Model
	status       string
	running      bool
	workflowName string
	workflowPath string
}

// NewExecutor creates an executor screen
func NewExecutor() Executor {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(styles.Primary)

	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
	)

	return Executor{
		spinner:  s,
		progress: p,
		status:   "Ready to execute workflows",
		running:  false,
	}
}

// Init initializes the executor
func (e Executor) Init() tea.Cmd {
	return nil
}

// Update handles executor messages
func (e Executor) Update(msg tea.Msg) (Executor, tea.Cmd) {
	if !e.running {
		return e, nil
	}

	var cmd tea.Cmd
	e.spinner, cmd = e.spinner.Update(msg)
	return e, cmd
}

// View renders the executor
func (e Executor) View() string {
	title := styles.Title.Render("Workflow Executor")
	if !e.running {
		return e.idleView(title)
	}
	return e.runningView(title)
}

// idleView renders the executor when idle
func (e Executor) idleView(title string) string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		styles.Subtle.Render(e.status),
		"",
		styles.Subtle.Render("Select a workflow from the Browser and press Enter to execute."),
	)
}

// runningView renders the executor when running
func (e Executor) runningView(title string) string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		styles.Subtle.Render("Workflow: "+e.workflowName),
		styles.Subtle.Render("Path: "+e.workflowPath),
		"",
		e.spinner.View()+" "+e.status,
		"",
		e.progress.View(),
		"",
		styles.Subtle.Render("Execution in progress..."),
	)
}

// StartWorkflow prepares the executor to run a workflow
func (e Executor) StartWorkflow(name, path string) Executor {
	e.workflowName = name
	e.workflowPath = path
	e.running = true
	e.status = "Starting workflow..."
	return e
}

// ExecuteCmd returns the command to start execution
func (e Executor) ExecuteCmd() tea.Cmd {
	return e.spinner.Tick
}

// IsRunning returns whether the executor is currently running
func (e Executor) IsRunning() bool {
	return e.running
}
