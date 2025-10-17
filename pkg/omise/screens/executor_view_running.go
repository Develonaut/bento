package screens

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"

	"bento/pkg/omise/styles"
)

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

	// Show lifecycle status
	if len(e.lifecycleHistory) > 0 {
		// Show only the most recent lifecycle message
		lines = append(lines, e.lifecycleHistory[len(e.lifecycleHistory)-1], "")
	}

	// Show per-node progress using sequence component
	if len(e.nodeStates) > 0 {
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
