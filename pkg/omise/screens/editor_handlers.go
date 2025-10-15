package screens

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"bento/pkg/neta"
)

// handleResize processes window size changes
func (e Editor) handleResize(msg tea.WindowSizeMsg) (Editor, tea.Cmd) {
	e.width = msg.Width
	e.height = msg.Height
	return e, nil
}

// handleKey processes keyboard input
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

// handleStateKey processes state-specific keys
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

// handleNamingKey handles name entry state
func (e Editor) handleNamingKey(msg tea.KeyMsg) (Editor, tea.Cmd) {
	// TODO: Wire up Huh form for name entry
	if msg.String() == "enter" {
		e.bentoName = "new-bento"
		e.def.Name = e.bentoName
		e.state = StateSelectingType
	}
	return e, nil
}

// handleTypeSelectionKey handles type selection
func (e Editor) handleTypeSelectionKey(msg tea.KeyMsg) (Editor, tea.Cmd) {
	// TODO: Wire up type selection from pantry
	if msg.String() == "enter" {
		e.currentNodeType = "http"
		e.state = StateConfiguringNode
		return e, e.launchWizard(e.currentNodeType)
	}
	return e, nil
}

// handleReviewKey handles review state
func (e Editor) handleReviewKey(msg tea.KeyMsg) (Editor, tea.Cmd) {
	switch msg.String() {
	case "a":
		e.state = StateSelectingType
		return e, nil
	case "s", "enter":
		return e, e.saveBento()
	}
	return e, nil
}

// handleNameEntered processes name entry
func (e Editor) handleNameEntered(msg BentoNameEnteredMsg) (Editor, tea.Cmd) {
	e.bentoName = msg.Name
	e.def.Name = msg.Name
	e.state = StateSelectingType
	return e, nil
}

// handleTypeSelected processes type selection
func (e Editor) handleTypeSelected(msg NodeTypeSelectedMsg) (Editor, tea.Cmd) {
	e.currentNodeType = msg.Type
	e.state = StateConfiguringNode
	return e, e.launchWizard(msg.Type)
}

// handleNodeConfigured processes configured node
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

// shouldSetAsRoot checks if node should be root
func (e Editor) shouldSetAsRoot() bool {
	return len(e.def.Nodes) == 0 && e.def.Type == ""
}

// buildNode creates a node from configured message
func buildNode(msg NodeConfiguredMsg) neta.Definition {
	return neta.Definition{
		Version:    neta.CurrentVersion,
		Type:       msg.Type,
		Name:       msg.Name,
		Parameters: msg.Parameters,
	}
}

// setRootNode sets the root node of definition
func setRootNode(def neta.Definition, msg NodeConfiguredMsg) neta.Definition {
	def.Type = msg.Type
	def.Parameters = msg.Parameters
	return def
}

// appendNode adds a node to definition
func appendNode(def neta.Definition, node neta.Definition) neta.Definition {
	if def.Type == "" {
		def.Type = "group.sequence"
	}
	def.Nodes = append(def.Nodes, node)
	return def
}

// launchWizard starts configuration wizard
func (e Editor) launchWizard(nodeType string) tea.Cmd {
	return func() tea.Msg {
		// Get schema for this node type
		schema, ok := e.validator.GetSchema(nodeType)
		if !ok {
			// No schema available, return with default parameters
			return NodeConfiguredMsg{
				Type:       nodeType,
				Name:       fmt.Sprintf("New %s Node", nodeType),
				Parameters: map[string]interface{}{},
			}
		}

		// Create and run wizard
		wizard := NewNodeWizard(nodeType, schema)
		params, err := wizard.Run()
		if err != nil {
			// Wizard cancelled or error occurred
			return EditorCancelledMsg{}
		}

		// Extract node name from params
		var nodeName string
		if name, ok := params["name"]; ok {
			if nameStr, ok := name.(*string); ok {
				nodeName = *nameStr
				delete(params, "name") // Remove name from parameters
			}
		}
		if nodeName == "" {
			nodeName = fmt.Sprintf("New %s Node", nodeType)
		}

		// Convert pointer values to actual values
		actualParams := make(map[string]interface{})
		for k, v := range params {
			actualParams[k] = derefValue(v)
		}

		return NodeConfiguredMsg{
			Type:       nodeType,
			Name:       nodeName,
			Parameters: actualParams,
		}
	}
}

// derefValue dereferences pointer values
func derefValue(v interface{}) interface{} {
	switch val := v.(type) {
	case *string:
		return *val
	case *int:
		return *val
	case *bool:
		return *val
	default:
		return v
	}
}

// saveBento saves the bento to Jubako with context
func (e Editor) saveBento() tea.Cmd {
	return func() tea.Msg {
		// Check if context was cancelled
		select {
		case <-e.ctx.Done():
			return EditorCancelledMsg{}
		default:
		}

		if err := e.store.Save(e.bentoName, e.def); err != nil {
			return EditorSaveErrorMsg{Error: err}
		}
		return EditorSavedMsg{Name: e.bentoName}
	}
}

// cancelEditor cancels the editor
func cancelEditor() tea.Cmd {
	return func() tea.Msg {
		return EditorCancelledMsg{}
	}
}
