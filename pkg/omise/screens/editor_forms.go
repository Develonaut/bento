package screens

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Form setup functions for editor wizards and save operations

// startTypeSelection creates and initializes the node type selection form
func (e Editor) startTypeSelection() (Editor, tea.Cmd) {
	var nodeType string
	nodeTypes := e.validator.ListTypes() // Use validator, not registry
	if len(nodeTypes) == 0 {
		e.state = StateReview
		return e, nil
	}
	e.currentForm = createNodeTypeForm(nodeTypes, &nodeType)
	e.formValues = map[string]interface{}{"nodeType": &nodeType}
	e.state = StateSelectingType
	return e, e.currentForm.Init()
}

// startWizard creates and initializes the node configuration wizard for the given type
func (e Editor) startWizard(nodeType string) (Editor, tea.Cmd) {
	schema, ok := e.validator.GetSchema(nodeType)
	if !ok {
		e.state = StateReview
		e.message = "No schema found for node type"
		return e, nil
	}
	e.formValues = make(map[string]interface{})
	wizard := NewNodeWizard(nodeType, schema, e.formValues)
	e.currentForm = wizard.Form()
	e.currentNodeType = nodeType
	e.state = StateConfiguringNode
	return e, e.currentForm.Init()
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
