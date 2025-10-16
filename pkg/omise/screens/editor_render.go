package screens

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"bento/pkg/omise/styles"
)

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

func (e Editor) renderTitle() string {
	mode := "Create New Bento"
	if e.mode == EditorModeEdit {
		mode = fmt.Sprintf("Edit Bento: %s", e.bentoName)
	}
	return styles.Title.Render(mode)
}

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

func (e Editor) renderNaming() string {
	if e.currentForm != nil {
		return e.currentForm.View()
	}
	return styles.Subtle.Render("Enter bento name:\n\n[Name entry form here]")
}

func (e Editor) renderTypeSelection() string {
	if e.currentForm != nil {
		return e.currentForm.View()
	}
	types := e.validator.ListTypes() // Use validator, not registry
	content := "Select node type:\n\n"
	for _, t := range types {
		content += fmt.Sprintf("  • %s\n", t)
	}
	return styles.Subtle.Render(content)
}

func (e Editor) renderConfiguration() string {
	if e.currentForm != nil {
		return e.currentForm.View()
	}
	return styles.Subtle.Render(
		fmt.Sprintf("Configure %s node:\n\n[Wizard form here]", e.currentNodeType),
	)
}

func (e Editor) renderReview() string {
	if e.viewMode == ViewModeVisual {
		return e.renderVisualBentoBox()
	}
	return e.renderListView()
}

func (e Editor) renderFooter() string {
	return e.helpView.RenderFooter(e.message, e.keys)
}
