package guided_creation

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *GuidedModal) renderPreview(formHeight int) string {
	s := m.styles

	var (
		bentoInfo string
		nodeInfo  string
	)

	// Bento metadata
	if m.definition.Name != "" {
		bentoInfo = fmt.Sprintf("%s %s\n%s",
			m.definition.Icon,
			m.definition.Name,
			m.definition.Description,
		)
	} else {
		bentoInfo = "(Bento metadata pending...)"
	}

	// Node information
	if len(m.definition.Nodes) > 0 {
		var nodeList strings.Builder
		for i, node := range m.definition.Nodes {
			fmt.Fprintf(&nodeList, "%d. %s (%s)\n", i+1, node.Name, node.Type)
		}
		nodeInfo = "\n\n" + s.StatusHeader.Render("Nodes") + "\n" + nodeList.String()
	} else {
		nodeInfo = "\n\n" + s.StatusHeader.Render("Nodes") + "\n(No nodes yet)"
	}

	const previewWidth = 40
	previewMarginLeft := m.width - previewWidth - lipgloss.Width(m.form.View()) - s.Status.GetMarginRight()

	return s.Status.
		Height(formHeight).
		Width(previewWidth).
		MarginLeft(previewMarginLeft).
		Render(s.StatusHeader.Render("Bento Preview") + "\n\n" +
			bentoInfo +
			nodeInfo)
}
