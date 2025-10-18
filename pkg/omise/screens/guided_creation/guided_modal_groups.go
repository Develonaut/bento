package guided_creation

import (
	"fmt"

	"github.com/charmbracelet/huh"

	"bento/pkg/neta"
)

// createGroupContextForm prompts user what to do after creating a group node
func (m *GuidedModal) createGroupContextForm(groupName string) *huh.Form {
	var choice string
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("group_context").
				Title(fmt.Sprintf("Group '%s' created. What next?", groupName)).
				Description("Choose whether to add children to this group or continue at the current level").
				Options(m.groupContextOptions(groupName)...).
				Value(&choice),
		).Title("Group Context:"),
	).WithWidth(m.width - 48).WithShowHelp(false).WithShowErrors(false)
}

func (m *GuidedModal) groupContextOptions(groupName string) []huh.Option[string] {
	return []huh.Option[string]{
		huh.NewOption(fmt.Sprintf("Add child to '%s'", groupName), "add_child"),
		huh.NewOption("Add another node at current level", "add_sibling"),
		huh.NewOption("Done with current level", "done_level"),
		huh.NewOption("Save bento", "save"),
	}
}

// pushParentContext pushes the current parent onto the stack and sets a new parent
func (m *GuidedModal) pushParentContext(newParent *neta.Definition) {
	m.nodeStack = append(m.nodeStack, m.currentParent)
	m.currentParent = newParent
}

// popParentContext pops the parent stack and restores the previous parent
func (m *GuidedModal) popParentContext() bool {
	if len(m.nodeStack) == 0 {
		// Already at root
		return false
	}

	m.currentParent = m.nodeStack[len(m.nodeStack)-1]
	m.nodeStack = m.nodeStack[:len(m.nodeStack)-1]
	return true
}

// getCurrentNodes returns the nodes array we're currently adding to
func (m *GuidedModal) getCurrentNodes() *[]neta.Definition {
	if m.currentParent != nil {
		return &m.currentParent.Nodes
	}
	return &m.definition.Nodes
}

// addNodeToCurrent adds a node to the current parent (or root if no parent)
func (m *GuidedModal) addNodeToCurrent(node neta.Definition) {
	nodes := m.getCurrentNodes()
	*nodes = append(*nodes, node)
}

// isGroupNode checks if a node type is a group type
func isGroupNode(nodeType string) bool {
	return nodeType == "group.sequence" || nodeType == "group.parallel"
}
