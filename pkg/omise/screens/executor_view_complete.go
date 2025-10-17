package screens

import (
	"github.com/charmbracelet/lipgloss"

	"bento/pkg/omise/styles"
)

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
			styles.Subtle.Render(truncateToOneLine(result, 80)),
			"",
		)
	}
	return lines
}

// buildCompletionCenter builds the center section for completion view
func (e Executor) buildCompletionCenter() string {
	lines := []string{}

	// Show lifecycle status
	if len(e.lifecycleHistory) > 0 {
		// Show only the most recent lifecycle message
		lines = append(lines, e.lifecycleHistory[len(e.lifecycleHistory)-1], "")
	}

	// Show node sequence if available
	if len(e.nodeStates) > 0 {
		lines = append(lines, e.sequence.View(), "")
	}

	// Add fun completion message
	if e.success {
		completionMsg := "✨ Delicious!"
		lines = append(lines, styles.SuccessStyle.Render(completionMsg), "")
	} else {
		failureMsg := "👹 Oh no!"
		lines = append(lines, styles.ErrorStyle.Render(failureMsg), "")
	}

	// Show errors if execution failed
	lines = appendErrorLines(lines, e.success, e.errorMsg)

	// Show output if execution succeeded
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
