package screens

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// Message and keyboard event handlers for the editor

func (e Editor) handleResize(msg tea.WindowSizeMsg) (Editor, tea.Cmd) {
	e.width = msg.Width
	e.height = msg.Height
	return e, nil
}

func (e Editor) handleKey(msg tea.KeyMsg) (Editor, tea.Cmd) {
	switch msg.String() {
	case "esc":
		return e, cancelEditor()
	case "ctrl+s":
		return e, e.saveBento()
	default:
		return e.handleStateKey(msg)
	}
}

func (e Editor) handleStateKey(msg tea.KeyMsg) (Editor, tea.Cmd) {
	switch e.state {
	case StateNaming:
		return e.handleNamingKey(msg)
	case StateSelectingType:
		return e.handleTypeSelectionKey(msg)
	case StateReview:
		return e.handleReviewKey(msg)
	}
	return e, nil
}

func (e Editor) handleNamingKey(msg tea.KeyMsg) (Editor, tea.Cmd) {
	// Name entry handled by launchNameForm via message
	// This handler only processes ESC (handled globally)
	return e, nil
}

func (e Editor) handleTypeSelectionKey(msg tea.KeyMsg) (Editor, tea.Cmd) {
	// Type selection handled by launchTypeForm via message
	// This handler only processes ESC (handled globally)
	return e, nil
}

func (e Editor) handleReviewKey(msg tea.KeyMsg) (Editor, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		return e.navigateUp(), nil
	case "down", "j":
		return e.navigateDown(), nil
	case "e":
		return e, e.editNode(e.selectedNodeIndex)
	case "m":
		return e.moveNode(e.selectedNodeIndex)
	case "d":
		return e.deleteNode(e.selectedNodeIndex)
	case "r":
		return e, e.runBento()
	case "v":
		return e.toggleViewMode(), nil
	case "a":
		e.state = StateSelectingType
		return e, e.launchTypeForm()
	case "s", "enter":
		return e, e.saveBento()
	}
	return e, nil
}

func (e Editor) handleNameEntered(msg BentoNameEnteredMsg) (Editor, tea.Cmd) {
	e.bentoName = msg.Name
	e.def.Name = msg.Name
	e.state = StateSelectingType
	return e, e.launchTypeForm()
}

func (e Editor) handleTypeSelected(msg NodeTypeSelectedMsg) (Editor, tea.Cmd) {
	e.currentNodeType = msg.Type
	e.state = StateConfiguringNode
	return e, e.launchWizard(msg.Type)
}

func (e Editor) handleNodeConfigured(msg NodeConfiguredMsg) (Editor, tea.Cmd) {
	node := buildNode(msg)

	if e.shouldSetAsRoot() {
		e.def = setRootNode(e.def, msg)
	} else {
		e.def = appendNode(e.def, node)
	}

	e.state = StateReview
	e.message = fmt.Sprintf("Added node: %s", msg.Name)
	return e, nil
}
