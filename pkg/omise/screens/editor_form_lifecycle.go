package screens

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// Form lifecycle methods for Huh integration

// createNameForm creates the bento name input form
func (e Editor) createNameForm() *huh.Form {
	var bentoName string
	e.formValues["bentoName"] = &bentoName
	return createBentoNameForm(&bentoName)
}

// extractNodeName extracts the node name from formValues
func (e Editor) extractNodeName() string {
	return extractNodeName(e.formValues, e.currentNodeType)
}

// updateForm handles form updates and checks for completion
func (e Editor) updateForm(msg tea.Msg) (Editor, tea.Cmd) {
	form, cmd := e.currentForm.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		e.currentForm = f
	}

	// Check if form completed
	if e.currentForm.State == huh.StateCompleted {
		return e.handleFormCompletion()
	}

	return e, cmd
}

// handleFormCompletion processes completed forms based on current state
func (e Editor) handleFormCompletion() (Editor, tea.Cmd) {
	e.currentForm = nil

	switch e.state {
	case StateNaming:
		return e.handleNamingCompletion()
	case StateSelectingType:
		return e.handleTypeSelectionCompletion()
	case StateConfiguringNode:
		return e.handleConfigurationCompletion()
	}

	return e, nil
}

// handleNamingCompletion processes the bento name form completion
func (e Editor) handleNamingCompletion() (Editor, tea.Cmd) {
	namePtr, ok := e.formValues["bentoName"].(*string)
	if !ok {
		return e, cancelEditor()
	}

	e.bentoName = *namePtr
	e.def.Name = *namePtr
	e.state = StateSelectingType

	var nodeType string
	nodeTypes := e.validator.ListTypes()
	e.currentForm = createNodeTypeForm(nodeTypes, &nodeType)
	e.formValues = map[string]interface{}{"nodeType": &nodeType}
	return e, e.currentForm.Init()
}

// handleTypeSelectionCompletion processes the node type selection form completion
func (e Editor) handleTypeSelectionCompletion() (Editor, tea.Cmd) {
	typePtr, ok := e.formValues["nodeType"].(*string)
	if !ok {
		return e, cancelEditor()
	}

	e.currentNodeType = *typePtr
	e.state = StateConfiguringNode

	schema, ok := e.validator.GetSchema(*typePtr)
	if !ok {
		e.state = StateReview
		return e, nil
	}

	e.formValues = make(map[string]interface{})
	wizard := NewNodeWizard(*typePtr, schema, e.formValues)
	e.currentForm = wizard.Form()
	return e, e.currentForm.Init()
}

// handleConfigurationCompletion processes the node configuration form completion
func (e Editor) handleConfigurationCompletion() (Editor, tea.Cmd) {
	nodeName := e.extractNodeName()
	actualParams := convertParamPointers(e.formValues)

	msg := NodeConfiguredMsg{
		Type:       e.currentNodeType,
		Name:       nodeName,
		Parameters: actualParams,
	}
	return e.handleNodeConfigured(msg)
}
