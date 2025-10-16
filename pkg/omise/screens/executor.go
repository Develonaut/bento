package screens

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"bento/pkg/omise/components"
	"bento/pkg/omise/styles"
)

// Executor shows bento execution progress
type Executor struct {
	spinner      components.Spinner
	progress     components.Progress
	status       string
	running      bool
	complete     bool
	success      bool
	bentoName    string
	bentoPath    string
	workDir      string
	errorMsg     string
	result       string
	copyFeedback string
}

// NewExecutor creates an executor screen
func NewExecutor() Executor {
	return Executor{
		spinner:  components.NewSpinner(),
		progress: components.NewProgress(40),
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
	// Handle theme changes
	if _, ok := msg.(styles.ThemeChangedMsg); ok {
		e.spinner = e.spinner.RebuildStyles()
	}

	// Handle keyboard input for complete state
	if e.complete {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			if keyMsg.String() == "c" {
				return e, e.copyToClipboard()
			}
		}
	}

	// Handle execution messages
	switch msg := msg.(type) {
	case ExecutionProgressMsg:
		e.status = msg.Status
		e.progress.SetPercent(msg.Progress)
		return e, nil
	case ExecutionCompleteMsg:
		e.running = false
		e.complete = true
		e.success = msg.Success
		e.result = formatResult(msg.Result)
		if msg.Success {
			e.status = "Execution completed successfully!"
			e.progress.SetPercent(1.0)
		} else {
			e.status = "Execution failed"
			if msg.Error != nil {
				e.errorMsg = msg.Error.Error()
			}
		}
		return e, nil
	case CopyResultMsg:
		e.copyFeedback = string(msg)
		return e, nil
	case ExecutionErrorMsg:
		e.running = false
		e.complete = true
		e.success = false
		e.status = "Execution error"
		e.errorMsg = msg.Error.Error()
		return e, nil
	}

	if !e.running {
		return e, nil
	}

	var cmd tea.Cmd
	e.spinner, cmd = e.spinner.Update(msg)
	return e, cmd
}

// View renders the executor
func (e Executor) View() string {
	title := styles.Title.Render("Bento Executor")
	if e.complete {
		return e.completeView(title)
	}
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
		styles.Subtle.Render("Select a bento from the Browser and press Enter/Space to execute."),
	)
}

// runningView renders the executor when running
func (e Executor) runningView(title string) string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		styles.Subtle.Render("Bento: "+e.bentoName),
		styles.Subtle.Render("Path: "+e.bentoPath),
		"",
		e.spinner.View()+" "+e.status,
		"",
		e.progress.View(),
		"",
		styles.Subtle.Render("Execution in progress..."),
	)
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
	e.status = "Starting bento..."
	e.progress.SetPercent(0.0)
	return e
}

// ExecuteCmd returns the command to start execution
func (e Executor) ExecuteCmd() tea.Cmd {
	return tea.Batch(
		e.spinner.Tick,
		ExecuteBentoCmd(e.bentoName, e.workDir),
	)
}

// completeView renders the executor when execution is complete
func (e Executor) completeView(title string) string {
	statusStyle := styles.SuccessStyle
	statusText := "✓ Success"
	if !e.success {
		statusStyle = styles.ErrorStyle
		statusText = "✗ Failed"
	}

	content := []string{
		title,
		"",
		styles.Subtle.Render("Bento: " + e.bentoName),
		styles.Subtle.Render("Path: " + e.bentoPath),
		"",
		statusStyle.Render(statusText),
		styles.Subtle.Render(e.status),
		"",
		e.progress.View(),
	}

	if e.errorMsg != "" {
		content = append(content, "", styles.ErrorStyle.Render("Error: "+e.errorMsg))
	}

	if e.result != "" && e.success {
		content = append(content, "", styles.Subtle.Render("Output:"), styles.Subtle.Render(e.result))
	}

	// Show copy feedback if present
	if e.copyFeedback != "" {
		content = append(content, "", styles.SuccessStyle.Render(e.copyFeedback))
	}

	return lipgloss.JoinVertical(lipgloss.Left, content...)
}

// copyToClipboard copies the result to clipboard
func (e Executor) copyToClipboard() tea.Cmd {
	return func() tea.Msg {
		return CopyResultCmd(e.result, e.bentoName, e.errorMsg, e.success)
	}
}

// formatResult formats the execution result for display/copying
func formatResult(result interface{}) string {
	if result == nil {
		return "No output"
	}
	return fmt.Sprintf("%v", result)
}

// IsRunning returns whether the executor is currently running
func (e Executor) IsRunning() bool {
	return e.running
}

// ContextualKeys returns the most important contextual keys for the footer
func (e Executor) ContextualKeys() []components.KeyHelp {
	// When execution is complete, show copy key
	if e.complete {
		// Only show copy if we have content to copy
		if (e.success && e.result != "") || (!e.success && e.errorMsg != "") {
			return []components.KeyHelp{
				{Key: "c", Desc: "copy output"},
			}
		}
	}
	// No contextual keys during execution or when no content
	return []components.KeyHelp{}
}
