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

		// Validation passed - clear error
		m.validationErr = nil

		// Check if the node we're about to add is a group
		isGroup := isGroupNode(m.currentNode.Type)
		if isGroup {
			// Initialize the group's Nodes array if needed
			if m.currentNode.Nodes == nil {
				m.currentNode.Nodes = []neta.Definition{}
			}
		}

		// Add node to current parent (or root)
		m.addNodeToCurrent(*m.currentNode)

		// If it's a group, we need to get a pointer to the node we just added
		// so we can use it as the parent for children
		var addedNodePtr *neta.Definition
		if isGroup {
			nodes := m.getCurrentNodes()
			addedNodePtr = &(*nodes)[len(*nodes)-1]
		}

		// If group, show group context menu
		if isGroup {
			// Store pointer to the group we just added
			m.currentNode = addedNodePtr

			// Move to group context prompt
			m.stage = guidedStageGroupContext
			m.form = m.createGroupContextForm(m.currentNode.Name)
			return m, m.form.Init()
		}

		// Not a group, proceed to normal continue stage
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

	case guidedStageGroupContext:
		// Handle group context menu choice
		choice := m.form.GetString("group_context")

		switch choice {
		case "add_child":
			// Push current node onto stack and make it the parent
			m.pushParentContext(m.currentNode)

			// Reset current node and go to type selection
			m.currentNode = nil
			m.stage = guidedStageNodeTypeSelect
			m.form = m.createNodeTypeSelectForm()
			return m, m.form.Init()

		case "add_sibling":
			// Add another node at the same level
			m.currentNode = nil
			m.stage = guidedStageNodeTypeSelect
			m.form = m.createNodeTypeSelectForm()
			return m, m.form.Init()

		case "done_level":
			// Pop back to parent level
			if m.popParentContext() {
				// Successfully popped to parent level
				m.stage = guidedStageContinue
				m.form = m.createContinueForm()
				return m, m.form.Init()
			} else {
				// Already at root, save
				m.state = guidedStateCompleted
				return m, m.saveBento()
			}

		case "save":
			// Save and exit
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
