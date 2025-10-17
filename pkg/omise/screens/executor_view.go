package screens

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"

	"bento/pkg/omise/emoji"
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
			lines = append(lines, emoji.Success+" Success")
		} else {
			lines = append(lines, emoji.Failure+" Failed")
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
		styles.Subtle.Render(emoji.Bento+" Ready to execute bentos"),
		"",
		styles.Subtle.Render("Select a bento from the Browser and press 'r' to run."),
	)
}
