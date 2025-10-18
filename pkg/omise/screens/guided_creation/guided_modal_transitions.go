package guided_creation

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"bento/pkg/neta"
)

// handleStageTransition manages transitions between guided creation stages
func (m *GuidedModal) handleStageTransition() (*GuidedModal, tea.Cmd) {
	switch m.stage {
	case guidedStageMetadata:
		// Metadata complete - move to node type selection
		m.stage = guidedStageNodeTypeSelect
		m.currentNode = nil // Reset current node
		m.form = m.createNodeTypeSelectForm()
		return m, m.form.Init()

	case guidedStageNodeTypeSelect:
		// Node type selected - create specialized form for that type
		nodeType := m.form.GetString("node_type")
		if nodeType == "" {
			// No type selected, stay on selection form
			return m, nil
		}

		// Initialize current node with selected type
		m.currentNode = &neta.Definition{
			Version:    "1.0",
			Type:       nodeType,
			Parameters: make(map[string]interface{}),
		}

		// Move to node parameters stage with specialized form
		m.stage = guidedStageNodeParameters
		m.form = m.createNodeFormForType(nodeType)
		return m, m.form.Init()

	case guidedStageNodeParameters:
		// Node parameters complete - validate before continuing
		if err := m.validateCurrentNode(); err != nil {
			// Validation failed - show error and stay on node form
			m.validationErr = err
			return m, nil
		}

		// Validation passed - clear error and add node to definition
		m.validationErr = nil
		m.definition.Nodes = append(m.definition.Nodes, *m.currentNode)

		// Move to continue prompt
		m.stage = guidedStageContinue
		m.form = m.createContinueForm()
		return m, m.form.Init()

	case guidedStageContinue:
		// Check user's choice
		choice := m.form.GetString("continue")
		if choice == "add" {
			// Add another node - go back to type selection
			m.stage = guidedStageNodeTypeSelect
			m.currentNode = nil // Reset for new node
			m.form = m.createNodeTypeSelectForm()
			return m, m.form.Init()
		} else {
			// Done - save bento
			m.state = guidedStateCompleted
			return m, m.saveBento()
		}
	}

	return m, nil
}

// validateCurrentNode validates the current node being built
func (m *GuidedModal) validateCurrentNode() error {
	if m.currentNode == nil {
		return fmt.Errorf("no node to validate")
	}

	// Validate using neta validator
	return m.validator.Validate(*m.currentNode)
}
