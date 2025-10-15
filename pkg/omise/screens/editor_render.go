package screens

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"bento/pkg/omise/styles"
)

// View renders the editor
func (e Editor) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		e.renderTitle(),
		"",
		e.renderContent(),
		"",
		e.renderFooter(),
	)
}

// renderTitle renders the editor title
func (e Editor) renderTitle() string {
	mode := "Create New Bento"
	if e.mode == EditorModeEdit {
		mode = fmt.Sprintf("Edit Bento: %s", e.bentoName)
	}
	return styles.Title.Render(mode)
}

// renderContent renders state-specific content
func (e Editor) renderContent() string {
	switch e.state {
	case StateNaming:
		return e.renderNaming()
	case StateSelectingType:
		return e.renderTypeSelection()
	case StateConfiguringNode:
		return e.renderConfiguration()
	case StateReview:
		return e.renderReview()
	}
	return ""
}

// renderNaming renders name entry
func (e Editor) renderNaming() string {
	return styles.Subtle.Render("Enter bento name:\n\n[Name entry form here]")
}

// renderTypeSelection renders node type selection
func (e Editor) renderTypeSelection() string {
	types := e.registry.List()
	content := "Select node type:\n\n"
	for _, t := range types {
		content += fmt.Sprintf("  • %s\n", t)
	}
	return styles.Subtle.Render(content)
}

// renderConfiguration renders parameter configuration
func (e Editor) renderConfiguration() string {
	return styles.Subtle.Render(
		fmt.Sprintf("Configure %s node:\n\n[Wizard form here]", e.currentNodeType),
	)
}

// renderReview renders bento review
func (e Editor) renderReview() string {
	content := fmt.Sprintf("Bento: %s (v%s)\n", e.def.Name, e.def.Version)
	content += fmt.Sprintf("Type: %s\n\n", e.def.Type)

	if len(e.def.Nodes) > 0 {
		content += "Nodes:\n"
		content += e.renderNodeList()
	}

	return styles.Subtle.Render(content)
}

// renderNodeList renders the list of nodes
func (e Editor) renderNodeList() string {
	content := ""
	for i, node := range e.def.Nodes {
		content += fmt.Sprintf("  %d. %s (%s)\n", i+1, node.Name, node.Type)
	}
	return content
}

// renderFooter renders keyboard shortcuts
func (e Editor) renderFooter() string {
	shortcuts := e.getShortcuts()
	if e.message != "" {
		shortcuts = e.message + " • " + shortcuts
	}
	return styles.Subtle.Render(shortcuts)
}

// getShortcuts returns state-specific shortcuts
func (e Editor) getShortcuts() string {
	if e.state == StateReview {
		return "a: Add node • s: Save • esc: Cancel"
	}
	return "esc: Cancel • ctrl+s: Save"
}
