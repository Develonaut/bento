package screens

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"bento/pkg/neta"
)

func (e Editor) editNode(index int) (Editor, tea.Cmd) {
	node := e.getNode(index)
	if node == nil {
		return e, nil
	}

	// Start wizard with existing node type
	return e.startWizard(node.Type)
}

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

func (e Editor) runBento() tea.Cmd {
	return func() tea.Msg {
		return RunBentoFromEditorMsg{
			Def: e.def,
		}
	}
}

func (e Editor) getNode(index int) *neta.Definition {
	nodes := e.getNodes()
	if index < 0 || index >= len(nodes) {
		return nil
	}
	return &nodes[index]
}

// getNodes returns all navigable nodes in the bento.
// For single-node bentos (root has non-group type), returns the root as single element.
// For multi-node bentos (group types), returns the nodes array.
func (e Editor) getNodes() []neta.Definition {
	if e.def.Type != "" && e.def.Type != "group.sequence" && e.def.Type != "group.parallel" {
		return []neta.Definition{e.def}
	}
	return e.def.Nodes
}

func (e Editor) navigateUp() Editor {
	if e.selectedNodeIndex > 0 {
		e.selectedNodeIndex--
	}
	return e
}

func (e Editor) navigateDown() Editor {
	nodes := e.getNodes()
	if e.selectedNodeIndex < len(nodes)-1 {
		e.selectedNodeIndex++
	}
	return e
}

func (e Editor) toggleViewMode() Editor {
	if e.viewMode == ViewModeList {
		e.viewMode = ViewModeVisual
	} else {
		e.viewMode = ViewModeList
	}
	return e
}
