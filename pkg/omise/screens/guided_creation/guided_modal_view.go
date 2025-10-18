package guided_creation

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *GuidedModal) View() string {
	switch m.state {
	case guidedStateCompleted:
		return m.renderCompleted()
	case guidedStateCancelled:
		return m.renderCancelled()
	default:
		return m.renderActiveForm()
	}
}

func (m *GuidedModal) renderCompleted() string {
	title := m.styles.Highlight.Render(m.definition.Name)
	var b strings.Builder
	fmt.Fprintf(&b, "✓ Bento created successfully!\n\n")
	fmt.Fprintf(&b, "Name: %s\n", title)
	fmt.Fprintf(&b, "Nodes: %d\n", len(m.definition.Nodes))
	fmt.Fprintf(&b, "\nPress any key to return to browser...")
	return m.styles.Status.Margin(0, 1).Padding(1, 2).Width(60).Render(b.String()) + "\n\n"
}

func (m *GuidedModal) renderCancelled() string {
	var b strings.Builder
	fmt.Fprintf(&b, "✗ Bento creation cancelled\n\n")
	fmt.Fprintf(&b, "Press any key to return to browser...")
	return m.styles.Status.Margin(0, 1).Padding(1, 2).Width(60).Render(b.String()) + "\n\n"
}

func (m *GuidedModal) renderActiveForm() string {
	form := m.renderForm()
	preview := m.renderPreview(lipgloss.Height(form))
	header := m.renderHeader()
	breadcrumb := m.renderBreadcrumb()
	footer := m.renderFooter()
	body := lipgloss.JoinHorizontal(lipgloss.Left, form, preview)
	return m.styles.Base.Render(header + "\n" + breadcrumb + "\n" + body + "\n\n" + footer)
}

func (m *GuidedModal) renderForm() string {
	v := strings.TrimSuffix(m.form.View(), "\n\n")
	return m.lg.NewStyle().Margin(1, 0).Render(v)
}

func (m *GuidedModal) renderHeader() string {
	hasError := len(m.form.Errors()) > 0 || m.validationErr != nil
	if !hasError {
		return m.appBoundaryView("Create New Bento")
	}
	if m.validationErr != nil {
		return m.appErrorBoundaryView("Validation Error: " + m.validationErr.Error())
	}
	return m.appErrorBoundaryView(m.errorView())
}

func (m *GuidedModal) renderFooter() string {
	hasError := len(m.form.Errors()) > 0 || m.validationErr != nil
	if hasError {
		return m.appErrorBoundaryView("")
	}

	// Build help text with navigation shortcuts
	formHelp := m.form.Help().ShortHelpView(m.form.KeyBinds())
	navHelp := m.renderNavigationHelp()
	helpText := formHelp
	if navHelp != "" {
		if formHelp != "" {
			helpText += " • " + navHelp
		} else {
			helpText = navHelp
		}
	}

	return m.appBoundaryView(helpText)
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

// renderBreadcrumb shows the current hierarchy context
func (m *GuidedModal) renderBreadcrumb() string {
	if m.currentParent == nil {
		// At root level
		return m.styles.Breadcrumb.Render("Context: Root")
	}

	// Build breadcrumb from stack
	parts := []string{"Root"}
	for _, parent := range m.nodeStack {
		if parent != nil {
			parts = append(parts, parent.Name)
		}
	}
	parts = append(parts, m.currentParent.Name)

	breadcrumbText := "Context: " + strings.Join(parts, " > ")
	return m.styles.Breadcrumb.Render(breadcrumbText)
}

// renderNavigationHelp shows available navigation shortcuts
func (m *GuidedModal) renderNavigationHelp() string {
	var parts []string

	// Show navigation if we have history
	if len(m.history.stages) > 0 {
		parts = append(parts, "ctrl+↑/↓: navigate")
	}

	// Show delete if on node stage
	if m.stage == guidedStageNodeParameters && m.currentNode != nil {
		parts = append(parts, "ctrl+d: delete")
	}

	if len(parts) == 0 {
		return ""
	}

	return strings.Join(parts, " • ")
}
