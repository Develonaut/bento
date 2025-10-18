package guided_creation

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *GuidedModal) View() string {
	s := m.styles

	switch m.state {
	case guidedStateCompleted:
		title := s.Highlight.Render(m.definition.Name)
		var b strings.Builder
		fmt.Fprintf(&b, "✓ Bento created successfully!\n\n")
		fmt.Fprintf(&b, "Name: %s\n", title)
		fmt.Fprintf(&b, "Nodes: %d\n", len(m.definition.Nodes))
		fmt.Fprintf(&b, "\nPress any key to return to browser...")
		return s.Status.Margin(0, 1).Padding(1, 2).Width(60).Render(b.String()) + "\n\n"

	case guidedStateCancelled:
		var b strings.Builder
		fmt.Fprintf(&b, "✗ Bento creation cancelled\n\n")
		fmt.Fprintf(&b, "Press any key to return to browser...")
		return s.Status.Margin(0, 1).Padding(1, 2).Width(60).Render(b.String()) + "\n\n"

	default:
		// Form (left side)
		v := strings.TrimSuffix(m.form.View(), "\n\n")
		form := m.lg.NewStyle().Margin(1, 0).Render(v)

		// Preview (right side)
		preview := m.renderPreview(lipgloss.Height(form))

		errors := m.form.Errors()

		// Check for both form validation errors and node validation errors
		hasError := len(errors) > 0 || m.validationErr != nil

		header := m.appBoundaryView("Create New Bento")
		if hasError {
			if m.validationErr != nil {
				header = m.appErrorBoundaryView("Validation Error: " + m.validationErr.Error())
			} else {
				header = m.appErrorBoundaryView(m.errorView())
			}
		}

		body := lipgloss.JoinHorizontal(lipgloss.Left, form, preview)

		footer := m.appBoundaryView(m.form.Help().ShortHelpView(m.form.KeyBinds()))
		if hasError {
			footer = m.appErrorBoundaryView("")
		}

		return s.Base.Render(header + "\n" + body + "\n\n" + footer)
	}
}

func (m *GuidedModal) errorView() string {
	var s string
	for _, err := range m.form.Errors() {
		s += err.Error()
	}
	return s
}

func (m *GuidedModal) appBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.HeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(guidedIndigo),
	)
}

func (m *GuidedModal) appErrorBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.ErrorHeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(guidedRed),
	)
}
