package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"bento/pkg/neta"
	"bento/pkg/omise/styles"
)

// renderVisualBentoBox renders visual bento box
func (e Editor) renderVisualBentoBox() string {
	nodes := e.getNodes()

	// Header
	header := e.renderBentoHeader()

	// Compartments
	compartments := []string{}
	for i, node := range nodes {
		compartment := e.renderCompartment(node, i == e.selectedNodeIndex, i+1)
		compartments = append(compartments, compartment)
	}

	// Join vertically
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		strings.Join(compartments, "\n\n"),
	)

	return content
}

// renderBentoHeader renders the bento box header
func (e Editor) renderBentoHeader() string {
	title := fmt.Sprintf("🍱 %s (v%s)", e.def.Name, e.def.Version)

	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, true, false, true).
		BorderForeground(styles.Primary).
		Padding(0, 1).
		Width(60).
		Render(title)
}

// renderCompartment renders a node as a compartment
func (e Editor) renderCompartment(node neta.Definition, selected bool, number int) string {
	border, borderColor := e.getCompartmentStyle(selected)
	content := e.buildCompartmentContent(node, number)
	return e.renderCompartmentBox(content, border, borderColor)
}

// getCompartmentStyle returns border style based on selection state
func (e Editor) getCompartmentStyle(selected bool) (lipgloss.Border, lipgloss.Color) {
	if selected {
		return lipgloss.ThickBorder(), styles.Primary
	}
	return lipgloss.RoundedBorder(), styles.Muted
}

// buildCompartmentContent builds the content for a compartment
func (e Editor) buildCompartmentContent(node neta.Definition, number int) string {
	title := fmt.Sprintf("%d. %s [%s]", number, node.Name, node.Type)
	params := e.formatParameters(node)

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.Text)
	return lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render(title),
		styles.Subtle.Render(params),
	)
}

// renderCompartmentBox renders the final compartment box with styling
func (e Editor) renderCompartmentBox(content string, border lipgloss.Border, borderColor lipgloss.Color) string {
	return lipgloss.NewStyle().
		Border(border).
		BorderForeground(borderColor).
		Padding(1, 2).
		Width(56).
		Render(content)
}

// formatParameters formats node parameters for display
func (e Editor) formatParameters(node neta.Definition) string {
	switch node.Type {
	case "http":
		return e.formatHTTPParams(node.Parameters)
	case "transform.jq":
		return e.formatJQParams(node.Parameters)
	case "conditional.if":
		return e.formatConditionalParams(node.Parameters)
	default:
		return e.formatGenericParams(node.Parameters)
	}
}

// formatHTTPParams formats HTTP parameters
func (e Editor) formatHTTPParams(params map[string]interface{}) string {
	method := "GET"
	url := ""

	if m, ok := params["method"].(string); ok {
		method = m
	}
	if u, ok := params["url"].(string); ok {
		url = u
	}

	return fmt.Sprintf("%s %s", method, url)
}

// formatJQParams formats JQ parameters
func (e Editor) formatJQParams(params map[string]interface{}) string {
	if filter, ok := params["filter"].(string); ok {
		return filter
	}
	return "No filter"
}

// formatConditionalParams formats conditional parameters
func (e Editor) formatConditionalParams(params map[string]interface{}) string {
	if cond, ok := params["condition"].(string); ok {
		return cond
	}
	return "No condition"
}

// formatGenericParams formats generic parameters
func (e Editor) formatGenericParams(params map[string]interface{}) string {
	if len(params) == 0 {
		return "No parameters"
	}

	lines := []string{}
	for key, val := range params {
		lines = append(lines, fmt.Sprintf("%s: %v", key, val))
	}

	return strings.Join(lines, "\n")
}

// renderListView renders simple list (Phase 7 view)
func (e Editor) renderListView() string {
	content := fmt.Sprintf("Bento: %s (v%s)\n", e.def.Name, e.def.Version)
	content += fmt.Sprintf("Type: %s\n\n", e.def.Type)

	nodes := e.getNodes()
	if len(nodes) > 0 {
		content += "Nodes:\n"
		for i, node := range nodes {
			selected := "  "
			if i == e.selectedNodeIndex {
				selected = "→ "
			}
			content += fmt.Sprintf("%s%d. %s (%s)\n", selected, i+1, node.Name, node.Type)
		}
	}

	return styles.Subtle.Render(content)
}
