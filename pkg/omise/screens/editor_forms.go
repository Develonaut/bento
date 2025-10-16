package screens

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Form launching functions for editor wizards and save operations

// launchWizard starts the node configuration wizard for the given node type.
// Uses schemas from validator to build type-specific Huh forms.
func (e Editor) launchWizard(nodeType string) tea.Cmd {
	return func() tea.Msg {
		schema, ok := e.validator.GetSchema(nodeType)
		if !ok {
			return defaultNodeConfig(nodeType)
		}

		wizard := NewNodeWizard(nodeType, schema)
		params, err := wizard.Run()
		if err != nil {
			return EditorCancelledMsg{}
		}

		nodeName := extractNodeName(params, nodeType)
		actualParams := convertParamPointers(params)

		return NodeConfiguredMsg{
			Type:       nodeType,
			Name:       nodeName,
			Parameters: actualParams,
		}
	}
}

// launchNameForm prompts the user to enter a bento name using Huh.
// Returns BentoNameEnteredMsg on success, EditorCancelledMsg on cancel.
func (e Editor) launchNameForm() tea.Cmd {
	return func() tea.Msg {
		name, err := promptBentoName()
		if err != nil {
			return EditorCancelledMsg{}
		}
		return BentoNameEnteredMsg{Name: name}
	}
}

// launchTypeForm prompts the user to select a node type from pantry.
// Returns NodeTypeSelectedMsg on success, EditorCancelledMsg on cancel.
func (e Editor) launchTypeForm() tea.Cmd {
	return func() tea.Msg {
		nodeTypes := e.registry.List()
		if len(nodeTypes) == 0 {
			return EditorCancelledMsg{}
		}

		nodeType, err := promptNodeType(nodeTypes)
		if err != nil {
			return EditorCancelledMsg{}
		}
		return NodeTypeSelectedMsg{Type: nodeType}
	}
}

// saveBento saves the current bento definition to Jubako storage.
// Returns EditorSavedMsg on success, EditorSaveErrorMsg on failure.
// Context cancellation support will be added when Store.Save accepts context.
func (e Editor) saveBento() tea.Cmd {
	return func() tea.Msg {
		if err := e.store.Save(e.bentoName, e.def); err != nil {
			return EditorSaveErrorMsg{Error: err}
		}
		return EditorSavedMsg{Name: e.bentoName}
	}
}

// cancelEditor creates a command that cancels the editor and returns to browser.
func cancelEditor() tea.Cmd {
	return func() tea.Msg {
		return EditorCancelledMsg{}
	}
}
