package screens

import (
	"encoding/json"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"bento/pkg/neta"
	"bento/pkg/omise/components"
	"bento/pkg/omise/styles"
)

// Emoji constants for lifecycle states
const (
	emojiBento     = "🍱"
	emojiExecuting = "⏳"
	emojiSuccess   = "✓"
	emojiFailure   = "✗"
)

// Executor shows bento execution progress
type Executor struct {
	spinner          components.Spinner
	progress         components.Progress
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
	// Handle theme changes and keyboard input
	if updated, cmd, handled := e.handleThemeAndInput(msg); handled {
		return updated, cmd
	}

	// Handle execution-related messages
	return e.handleExecutionMessages(msg)
}

// handleThemeAndInput handles theme changes and keyboard input
func (e Executor) handleThemeAndInput(msg tea.Msg) (Executor, tea.Cmd, bool) {
	// Handle theme changes
	if _, ok := msg.(styles.ThemeChangedMsg); ok {
		e.spinner = e.spinner.RebuildStyles()
		return e, nil, true
	}

	// Handle commands when execution is complete
	if e.complete {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "c":
				return e, e.copyToClipboard(), true
			case "r":
				// Rerun the bento
				e = e.StartBento(e.bentoName, e.bentoPath, e.workDir)
				return e, e.ExecuteCmd(), true
			}
		}
	}

	return e, nil, false
}

// handleExecutionMessages handles execution-related messages
func (e Executor) handleExecutionMessages(msg tea.Msg) (Executor, tea.Cmd) {
	switch msg := msg.(type) {
	case ExecutionProgressMsg:
		return e.handleProgressMsg(msg)
	case executionStillRunningMsg:
		return e.handleStillRunningMsg()
	case ExecutionCompleteMsg:
		return e.handleCompleteMsg(msg)
	case CopyResultMsg:
		e.copyFeedback = string(msg)
		return e, nil
	case ExecutionErrorMsg:
		return e.handleErrorMsg(msg)
	}

	return e.updateSpinner(msg)
}

// handleProgressMsg handles execution progress updates
func (e Executor) handleProgressMsg(msg ExecutionProgressMsg) (Executor, tea.Cmd) {
	if !e.running {
		return e, nil
	}
	e.status = msg.Status
	e.progressPercent = msg.Progress

	return e, tea.Batch(
		WaitForExecutionCmd(),
		ProgressTickCmd(msg.Progress),
	)
}

// handleStillRunningMsg continues polling for completion
func (e Executor) handleStillRunningMsg() (Executor, tea.Cmd) {
	if !e.running {
		return e, nil
	}
	return e, WaitForExecutionCmd()
}

// handleCompleteMsg handles execution completion
func (e Executor) handleCompleteMsg(msg ExecutionCompleteMsg) (Executor, tea.Cmd) {
	e.endTime = time.Now()
	e.running = false
	e.complete = true
	e.success = msg.Success
	e.result = formatResult(msg.Result)
	e.progressPercent = 1.0
	e.status = "Execution completed successfully!"

	// Add completion message to history
	e.lifecycleHistory = append(e.lifecycleHistory, emojiBento+" Bento packed!")

	if !msg.Success {
		e.status = "Execution failed"
		if msg.Error != nil {
			e.errorMsg = msg.Error.Error()
		}
	}
	return e, nil
}

// handleErrorMsg handles execution errors
func (e Executor) handleErrorMsg(msg ExecutionErrorMsg) (Executor, tea.Cmd) {
	e.running = false
	e.complete = true
	e.success = false
	e.status = "Execution error"
	e.errorMsg = msg.Error.Error()
	e.progressPercent = 0.0
	return e, nil
}

// updateSpinner updates the spinner during execution
func (e Executor) updateSpinner(msg tea.Msg) (Executor, tea.Cmd) {
	if !e.running {
		return e, nil
	}
	var spinnerCmd tea.Cmd
	e.spinner, spinnerCmd = e.spinner.Update(msg)
	return e, spinnerCmd
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

// calculateElapsedTime returns the elapsed execution time
func (e Executor) calculateElapsedTime() time.Duration {
	if !e.endTime.IsZero() {
		return e.endTime.Sub(e.startTime)
	}
	if !e.startTime.IsZero() {
		return time.Since(e.startTime)
	}
	return 0
}

// buildCompletionStatus returns the completion status line
func (e Executor) buildCompletionStatus() string {
	if !e.complete {
		return ""
	}
	if e.success {
		return emojiSuccess + " Success"
	}
	return emojiFailure + " Failed"
}

// renderFooter renders the footer section with status and timer
func (e Executor) renderFooter() string {
	elapsed := e.calculateElapsedTime()
	lines := []string{}

	// Status line
	if status := e.buildCompletionStatus(); status != "" {
		lines = append(lines, status)
	}

	// Timer line
	if elapsed > 0 {
		timerText := fmt.Sprintf("Execution time: %s", elapsed.Round(time.Millisecond))
		lines = append(lines, styles.Subtle.Render(timerText))
	}

	// Progress bar
	lines = append(lines, e.progress.ViewAs(e.progressPercent))

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// idleView renders the executor when idle
func (e Executor) idleView(title string) string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		styles.Subtle.Render(emojiBento+" Ready to execute bentos"),
		"",
		styles.Subtle.Render("Select a bento from the Browser and press 'r' to run."),
	)
}

// runningView renders the executor when running
func (e Executor) runningView(title string) string {
	header := e.buildRunningHeader(title)
	center := e.buildRunningCenter()
	footer := e.buildRunningFooter()
	return lipgloss.JoinVertical(lipgloss.Left, header, center, footer)
}

// buildRunningHeader builds the header section for running view
func (e Executor) buildRunningHeader(title string) string {
	lines := []string{
		title,
		"",
		styles.Subtle.Render("Bento: " + e.bentoName),
		styles.Subtle.Render("Path: " + e.bentoPath),
		"",
	}
	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// buildRunningCenter builds the center section for running view
func (e Executor) buildRunningCenter() string {
	lines := []string{}

	// Add all lifecycle history
	lines = append(lines, e.lifecycleHistory...)

	// Add current status with spinner
	lines = append(lines, e.spinner.View()+" "+e.status)
	lines = append(lines, "")

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// buildRunningFooter builds the footer section for running view
func (e Executor) buildRunningFooter() string {
	elapsed := time.Since(e.startTime)
	lines := []string{
		"",
		styles.Subtle.Render(fmt.Sprintf("Elapsed: %s", elapsed.Round(time.Millisecond))),
		"",
		e.progress.ViewAs(e.progressPercent),
	}
	return lipgloss.JoinVertical(lipgloss.Left, lines...)
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
	e.status = "Adding neta..."
	e.progressPercent = 0.0
	e.startTime = time.Now()
	e.endTime = time.Time{} // Zero value
	e.lifecycleHistory = []string{emojiBento + " Preparing bento box..."}
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
	header := e.buildCompletionHeader(title)
	center := e.buildCompletionCenter()
	footer := e.buildCompletionFooter()
	return lipgloss.JoinVertical(lipgloss.Left, header, center, footer)
}

// buildCompletionHeader builds the header section for completion view
func (e Executor) buildCompletionHeader(title string) string {
	lines := []string{
		title,
		"",
		styles.Subtle.Render("Bento: " + e.bentoName),
		styles.Subtle.Render("Path: " + e.bentoPath),
		"",
	}
	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// buildCompletionCenter builds the center section for completion view
func (e Executor) buildCompletionCenter() string {
	lines := []string{}

	// Add all lifecycle history
	lines = append(lines, e.lifecycleHistory...)
	lines = append(lines, "")

	// Add error message if failed
	if !e.success && e.errorMsg != "" {
		lines = append(lines,
			styles.ErrorStyle.Render("Error:"),
			styles.ErrorStyle.Render(truncateToOneLine(e.errorMsg, 80)),
			"",
		)
	}

	// Add output if successful (truncated to one line)
	if e.success && e.result != "" {
		lines = append(lines,
			styles.Subtle.Render("Output: "+truncateToOneLine(e.result, 80)),
			"",
		)
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// buildCompletionFooter builds the footer section for completion view
func (e Executor) buildCompletionFooter() string {
	lines := []string{e.renderFooter()}

	// Show copy feedback if present
	if e.copyFeedback != "" {
		lines = append(lines, "",
			styles.SuccessStyle.Render(e.copyFeedback))
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// copyToClipboard copies the result to clipboard
func (e Executor) copyToClipboard() tea.Cmd {
	return func() tea.Msg {
		return CopyResultCmd(e.result, e.bentoName, e.errorMsg, e.success)
	}
}

// formatResult formats the execution result with syntax highlighting
func formatResult(result interface{}) string {
	jsonStr, ok := extractJSON(result)
	if !ok {
		return jsonStr // Already a fallback message
	}
	return renderWithGlamour(jsonStr)
}

// extractJSON extracts and formats JSON from result
func extractJSON(result interface{}) (string, bool) {
	if result == nil {
		return "No output", false
	}

	// Type assert to neta.Result
	netaResult, ok := result.(neta.Result)
	if !ok {
		return fmt.Sprintf("%v", result), false
	}

	// Handle nil output
	if netaResult.Output == nil {
		return "No output", false
	}

	// Marshal to pretty JSON
	jsonBytes, err := json.MarshalIndent(netaResult.Output, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", netaResult.Output), false
	}

	return string(jsonBytes), true
}

// renderWithGlamour renders JSON with Glamour syntax highlighting
func renderWithGlamour(jsonStr string) string {
	// For now, skip Glamour to keep output on one line
	// Future: implement collapsible container for full output
	return jsonStr
}

// truncateToOneLine truncates a string to fit on one line
func truncateToOneLine(s string, maxLen int) string {
	// Remove all newlines and extra whitespace
	s = lipgloss.NewStyle().Inline(true).Render(s)

	// Find first newline if any
	for i, r := range s {
		if r == '\n' || r == '\r' {
			s = s[:i]
			break
		}
	}

	// Truncate to max length
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	return s
}

// IsRunning returns whether the executor is currently running
func (e Executor) IsRunning() bool {
	return e.running
}

// ContextualKeys returns the most important contextual keys for the footer
func (e Executor) ContextualKeys() []components.KeyHelp {
	// When execution is complete, show copy and rerun keys
	if e.complete {
		keys := []components.KeyHelp{
			{Key: "r", Desc: "rerun"},
		}
		// Add copy key if we have content to copy
		if (e.success && e.result != "") || (!e.success && e.errorMsg != "") {
			keys = append(keys, components.KeyHelp{Key: "c", Desc: "copy output"})
		}
		return keys
	}
	// No contextual keys during execution or when no content
	return []components.KeyHelp{}
}
