package screens

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"bento/pkg/neta"
)

// editNode launches wizard to edit node
func (e Editor) editNode(index int) tea.Cmd {
	node := e.getNode(index)
	if node == nil {
		return nil
	}

	// Launch wizard with existing node type
	return e.launchWizard(node.Type)
}

// moveNode swaps node with next
func (e Editor) moveNode(index int) (Editor, tea.Cmd) {
	if index >= len(e.def.Nodes)-1 {
		e.message = "Cannot move last node down"
		return e, nil
	}

	// Swap with next
	e.def.Nodes[index], e.def.Nodes[index+1] = e.def.Nodes[index+1], e.def.Nodes[index]

	e.message = fmt.Sprintf("Moved node %d down", index+1)
	return e, nil
}

// deleteNode removes node
func (e Editor) deleteNode(index int) (Editor, tea.Cmd) {
	if index >= len(e.def.Nodes) {
		return e, nil
	}

	// Remove node
	e.def.Nodes = append(e.def.Nodes[:index], e.def.Nodes[index+1:]...)

	// Adjust selection
	if e.selectedNodeIndex >= len(e.def.Nodes) && e.selectedNodeIndex > 0 {
		e.selectedNodeIndex--
	}

	e.message = "Node deleted"
	return e, nil
}

// runBento executes the bento
func (e Editor) runBento() tea.Cmd {
	return func() tea.Msg {
		return RunBentoFromEditorMsg{
			Def: e.def,
		}
	}
}

// getNode returns node at index
func (e Editor) getNode(index int) *neta.Definition {
	nodes := e.getNodes()
	if index < 0 || index >= len(nodes) {
		return nil
	}
	return &nodes[index]
}

// getNodes returns all nodes
func (e Editor) getNodes() []neta.Definition {
	// If single-node bento (root has type that's not a group)
	if e.def.Type != "" && e.def.Type != "group.sequence" && e.def.Type != "group.parallel" {
		return []neta.Definition{e.def}
	}
	// Multi-node bento
	return e.def.Nodes
}

// navigateUp moves selection up
func (e Editor) navigateUp() Editor {
	if e.selectedNodeIndex > 0 {
		e.selectedNodeIndex--
	}
	return e
}

// navigateDown moves selection down
func (e Editor) navigateDown() Editor {
	nodes := e.getNodes()
	if e.selectedNodeIndex < len(nodes)-1 {
		e.selectedNodeIndex++
	}
	return e
}

// toggleViewMode toggles between list and visual
func (e Editor) toggleViewMode() Editor {
	if e.viewMode == ViewModeList {
		e.viewMode = ViewModeVisual
	} else {
		e.viewMode = ViewModeList
	}
	return e
}
