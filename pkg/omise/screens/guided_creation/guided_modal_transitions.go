package guided_creation

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"bento/pkg/neta"
)

// handleStageTransition manages transitions between guided creation stages
func (m *GuidedModal) handleStageTransition() (*GuidedModal, tea.Cmd) {
	// Exit navigation mode if we were navigating
	if m.navigating {
		m.exitNavigationMode()
	}

	switch m.stage {
	case guidedStageMetadata:
		// Capture snapshot before transition
		m.captureSnapshot()

		// Update definition with metadata
		if m.editing {
			// In edit mode, use temp fields
			m.updateDefinitionFromEditForm()
		} else {
			// In create mode, use regular form fields
			m.updateDefinitionFromForm()
		}

		// If in edit mode, return to edit menu
		if m.editing {
			m.stage = guidedStageEditMenu
			m.form = m.createEditMenuForm()
			return m, m.form.Init()
		}

		// Otherwise, move to node type selection
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

		// Capture snapshot before transition
		m.captureSnapshot()

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

		// Capture snapshot before transition
		m.captureSnapshot()

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

		// If in edit mode, return to edit menu
		if m.editing {
			m.currentNode = nil
			m.stage = guidedStageEditMenu
			m.form = m.createEditMenuForm()
			return m, m.form.Init()
		}

		// Not a group, proceed to normal continue stage
		m.stage = guidedStageContinue
		m.form = m.createContinueForm()
		return m, m.form.Init()

	case guidedStageContinue:
		// Capture snapshot before transition
		m.captureSnapshot()

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
		// Capture snapshot before transition
		m.captureSnapshot()

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

	case guidedStageEditMenu:
		// Capture snapshot before transition
		m.captureSnapshot()

		// Handle edit menu choice
		choice := m.form.GetString("edit_choice")

		switch choice {
		case "metadata":
			// Edit metadata
			m.stage = guidedStageMetadata
			m.form = m.createMetadataEditForm()
			return m, m.form.Init()

		case "add_node":
			// Go to standard node creation flow
			m.stage = guidedStageNodeTypeSelect
			m.form = m.createNodeTypeSelectForm()
			return m, m.form.Init()

		case "edit_node":
			// Show list of nodes to select for editing
			m.deletingNode = false
			m.stage = guidedStageNodeList
			m.form = m.createNodeListForm(false)
			return m, m.form.Init()

		case "delete_node":
			// Show list of nodes to select for deletion
			m.deletingNode = true
			m.stage = guidedStageNodeList
			m.form = m.createNodeListForm(true)
			return m, m.form.Init()

		case "save":
			// Save and exit
			m.state = guidedStateCompleted
			return m, m.saveBento()

		case "cancel":
			// Cancel without saving
			m.state = guidedStateCancelled
			return m, func() tea.Msg {
				return GuidedCompleteMsg{Cancelled: true}
			}
		}

	case guidedStageNodeList:
		// Capture snapshot before transition
		m.captureSnapshot()

		// Load selected node from temp field
		nodeName := m.tempSelectedNode
		if nodeName == "" {
			// Back to menu
			m.stage = guidedStageEditMenu
			m.form = m.createEditMenuForm()
			return m, m.form.Init()
		}

		// Find node in definition
		node := m.findNodeByName(nodeName)
		if node == nil {
			// Node not found, return to menu
			m.stage = guidedStageEditMenu
			m.form = m.createEditMenuForm()
			return m, m.form.Init()
		}

		// Check if we're deleting or editing based on the flag
		if m.deletingNode {
			// Delete the node
			if m.deleteNodeByName(nodeName) {
				// Successfully deleted
			}
			// Clear the flag and return to edit menu
			m.deletingNode = false
			m.stage = guidedStageEditMenu
			m.form = m.createEditMenuForm()
			return m, m.form.Init()
		}

		// Editing mode - set as current node and load appropriate form
		m.editingNodeName = nodeName
		m.currentNode = node
		m.stage = guidedStageNodeEdit
		m.form = m.createNodeFormForTypeWithValues(node.Type, node)
		return m, m.form.Init()

	case guidedStageNodeEdit:
		// Capture snapshot before transition
		m.captureSnapshot()

		// Update node parameters from temp fields (not from form GetString)
		m.updateCurrentNodeFromTempFields(m.currentNode.Type)

		// Update the node in place in the definition
		m.updateNodeInPlace(m.editingNodeName, m.currentNode)

		// Clear editing state
		m.editingNodeName = ""
		m.currentNode = nil

		// Return to edit menu
		m.stage = guidedStageEditMenu
		m.form = m.createEditMenuForm()
		return m, m.form.Init()
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
