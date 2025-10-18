package guided_creation

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"

	"bento/pkg/neta"
)

// navigationHistory tracks the user's journey through the form stages
type navigationHistory struct {
	stages        []guidedStage       // History of visited stages
	stageData     []map[string]string // Form data at each stage
	nodeSnapshots []*neta.Definition  // Snapshots of nodes being built
}

// newNavigationHistory creates a new navigation history
func newNavigationHistory() navigationHistory {
	return navigationHistory{
		stages:        make([]guidedStage, 0),
		stageData:     make([]map[string]string, 0),
		nodeSnapshots: make([]*neta.Definition, 0),
	}
}

// captureSnapshot saves the current state to history
func (m *GuidedModal) captureSnapshot() {
	// Don't capture if we're in navigation mode
	if m.navigating {
		return
	}

	// Capture current stage
	m.history.stages = append(m.history.stages, m.stage)

	// Capture form data
	formData := m.extractFormData()
	m.history.stageData = append(m.history.stageData, formData)

	// Capture node snapshot if we're building a node
	var nodeSnapshot *neta.Definition
	if m.currentNode != nil {
		nodeCopy := *m.currentNode
		nodeSnapshot = &nodeCopy
	}
	m.history.nodeSnapshots = append(m.history.nodeSnapshots, nodeSnapshot)
}

// extractFormData gets all current form field values
func (m *GuidedModal) extractFormData() map[string]string {
	data := make(map[string]string)

	// Extract based on current stage
	switch m.stage {
	case guidedStageMetadata:
		data["name"] = m.form.GetString("name")
		data["description"] = m.form.GetString("description")

	case guidedStageNodeTypeSelect:
		data["node_type"] = m.form.GetString("node_type")

	case guidedStageNodeParameters:
		// Extract all node-specific fields
		if m.currentNode != nil {
			data["node_name"] = m.form.GetString("node_name")
			switch m.currentNode.Type {
			case "http":
				data["url"] = m.form.GetString("url")
				data["method"] = m.form.GetString("method")
			case "transform.jq", "jq":
				data["query"] = m.form.GetString("query")
			case "file.write":
				data["path"] = m.form.GetString("path")
			}
		}

	case guidedStageContinue:
		data["continue"] = m.form.GetString("continue")

	case guidedStageGroupContext:
		data["group_context"] = m.form.GetString("group_context")
	}

	return data
}

// navigateHistory moves backward or forward in the history
func (m *GuidedModal) navigateHistory(direction int) tea.Cmd {
	newIndex := m.historyIndex + direction

	// Validate bounds
	if newIndex < 0 || newIndex >= len(m.history.stages) {
		return nil
	}

	// Enter navigation mode
	m.navigating = true
	m.historyIndex = newIndex

	// Restore stage
	m.stage = m.history.stages[newIndex]

	// Restore node snapshot if available
	if newIndex < len(m.history.nodeSnapshots) && m.history.nodeSnapshots[newIndex] != nil {
		nodeCopy := *m.history.nodeSnapshots[newIndex]
		m.currentNode = &nodeCopy
	}

	// Recreate form for this stage with saved data
	m.form = m.recreateFormForStage(m.stage, m.history.stageData[newIndex])

	return m.form.Init()
}

// recreateFormForStage creates a form for a given stage with pre-filled data
func (m *GuidedModal) recreateFormForStage(stage guidedStage, data map[string]string) *huh.Form {
	switch stage {
	case guidedStageMetadata:
		return m.createMetadataForm()

	case guidedStageNodeTypeSelect:
		return m.createNodeTypeSelectForm()

	case guidedStageNodeParameters:
		if m.currentNode != nil {
			return m.createNodeFormForType(m.currentNode.Type)
		}
		return m.createNodeTypeSelectForm()

	case guidedStageContinue:
		return m.createContinueForm()

	case guidedStageGroupContext:
		if m.currentNode != nil {
			return m.createGroupContextForm(m.currentNode.Name)
		}
		return m.createContinueForm()

	default:
		return m.createMetadataForm()
	}
}

// deleteCurrentNode removes the current node from the definition
func (m *GuidedModal) deleteCurrentNode() tea.Cmd {
	if m.currentNode == nil {
		return nil
	}

	// Find and remove node from parent or root
	if m.currentParent != nil {
		m.currentParent.Nodes = removeNodeByName(m.currentParent.Nodes, m.currentNode.Name)
	} else {
		m.definition.Nodes = removeNodeByName(m.definition.Nodes, m.currentNode.Name)
	}

	// Clear current node and go back to continue stage
	m.currentNode = nil
	m.stage = guidedStageContinue
	m.form = m.createContinueForm()

	return m.form.Init()
}

// removeNodeByName removes a node from a slice by name
func removeNodeByName(nodes []neta.Definition, name string) []neta.Definition {
	result := make([]neta.Definition, 0, len(nodes))
	for _, node := range nodes {
		if node.Name != name {
			result = append(result, node)
		}
	}
	return result
}

// exitNavigationMode returns to normal flow from navigation mode
func (m *GuidedModal) exitNavigationMode() {
	m.navigating = false
	m.historyIndex = -1
}
