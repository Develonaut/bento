package screens

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"

	"bento/pkg/neta"
	"bento/pkg/omise/styles"
)

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

// renderFooter renders the footer section with status and timer
func (e Executor) renderFooter() string {
	lines := []string{}

	// Status line
	if e.complete {
		if e.success {
			lines = append(lines, emojiSuccess+" Success")
		} else {
			lines = append(lines, emojiFailure+" Failed")
		}
	}

	// Timer line - calculate elapsed time
	var elapsed time.Duration
	if !e.endTime.IsZero() {
		elapsed = e.endTime.Sub(e.startTime)
	} else if !e.startTime.IsZero() {
		elapsed = time.Since(e.startTime)
	}
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

	// Show per-node progress using sequence component
	if len(e.nodeStates) > 0 {
		lines = append(lines, "")
		lines = append(lines, e.sequence.View())
	} else {
		// Fallback to old status display if no nodes yet
		lines = append(lines, e.spinner.View()+" "+e.status)
	}

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

// appendErrorLines adds error message lines if execution failed
func appendErrorLines(lines []string, success bool, errorMsg string) []string {
	if !success && errorMsg != "" {
		return append(lines,
			styles.ErrorStyle.Render("Error:"),
			styles.ErrorStyle.Render(truncateToOneLine(errorMsg, 80)),
			"",
		)
	}
	return lines
}

// appendOutputLines adds output lines if execution succeeded
func appendOutputLines(lines []string, success bool, result string) []string {
	if success && result != "" {
		return append(lines,
			styles.Subtle.Render("Output: "+truncateToOneLine(result, 80)),
			"",
		)
	}
	return lines
}

// buildCompletionCenter builds the center section for completion view
func (e Executor) buildCompletionCenter() string {
	lines := []string{}
	lines = append(lines, e.lifecycleHistory...)

	if len(e.nodeStates) > 0 {
		lines = append(lines, "", e.sequence.View())
	}

	lines = append(lines, "")
	lines = appendErrorLines(lines, e.success, e.errorMsg)
	lines = appendOutputLines(lines, e.success, e.result)

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

// formatResult formats the execution result with syntax highlighting
func formatResult(result interface{}) string {
	jsonStr, ok := extractJSON(result)
	if !ok {
		return jsonStr // Already a fallback message
	}
	return jsonStr
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

	return marshalNetaResult(netaResult)
}

// marshalNetaResult converts neta.Result to JSON string
func marshalNetaResult(result neta.Result) (string, bool) {
	if result.Output == nil {
		return "No output", false
	}

	jsonBytes, err := json.MarshalIndent(result.Output, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", result.Output), false
	}

	return string(jsonBytes), true
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
