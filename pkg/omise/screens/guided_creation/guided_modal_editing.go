package guided_creation

import (
	"fmt"

	"github.com/charmbracelet/huh"

	"bento/pkg/neta"
)

// createEditMenuForm presents top-level editing options
func (m *GuidedModal) createEditMenuForm() *huh.Form {
	var choice string

	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("edit_choice").
				Title("What would you like to edit?").
				Description("Choose an action to perform on this bento").
				Options(m.editMenuOptions()...).
				Value(&choice),
		).Title(fmt.Sprintf("Editing: %s", m.definition.Name)),
	).
		WithWidth(m.width - 48).
		WithShowHelp(false).
		WithShowErrors(false)
}

// editMenuOptions returns the edit menu options
func (m *GuidedModal) editMenuOptions() []huh.Option[string] {
	return []huh.Option[string]{
		huh.NewOption("Edit metadata (name, description)", "metadata"),
		huh.NewOption("Edit an existing node", "edit_node"),
		huh.NewOption("Add a new node", "add_node"),
		huh.NewOption("Delete a node", "delete_node"),
		huh.NewOption("Save and exit", "save"),
		huh.NewOption("Cancel without saving", "cancel"),
	}
}

// createNodeListForm shows all nodes for selection
func (m *GuidedModal) createNodeListForm(forDelete bool) *huh.Form {
	title := "Select a node to edit"
	if forDelete {
		title = "Select a node to delete"
	}

	// Build options from existing nodes
	options := m.buildNodeListOptions(m.definition.Nodes, "")

	// Initialize selection - start with first node if available
	if len(m.definition.Nodes) > 0 {
		m.tempSelectedNode = m.definition.Nodes[0].Name
	} else {
		m.tempSelectedNode = ""
	}

	if len(options) == 0 {
		options = append(options, huh.NewOption("(No nodes to select)", ""))
	}

	// Add a "Back" option
	options = append(options, huh.NewOption("← Back to menu", ""))

	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("selected_node").
				Title(title).
				Options(options...).
				Value(&m.tempSelectedNode).
				Height(10), // Set explicit height to show options
		).Title("Nodes:"),
	).
		WithWidth(m.width - 48).
		WithShowHelp(false).
		WithShowErrors(false)
}

// buildNodeListOptions recursively builds node options with hierarchy
func (m *GuidedModal) buildNodeListOptions(nodes []neta.Definition, prefix string) []huh.Option[string] {
	options := make([]huh.Option[string], 0)

	for _, node := range nodes {
		label := fmt.Sprintf("%s%s (%s)", prefix, node.Name, node.Type)
		options = append(options, huh.NewOption(label, node.Name))

		// If it's a group with children, add them indented
		if len(node.Nodes) > 0 {
			childOptions := m.buildNodeListOptions(node.Nodes, prefix+"  ")
			options = append(options, childOptions...)
		}
	}

	return options
}

// createMetadataEditForm creates a metadata form pre-populated with values
func (m *GuidedModal) createMetadataEditForm() *huh.Form {
	// Store current values in the modal for form binding
	// This is a workaround since huh forms need variables to bind to
	// We'll use fields on the modal struct for this
	m.tempName = m.definition.Name
	m.tempDescription = m.definition.Description

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("Name").
				Description("A short, descriptive name for this workflow").
				Value(&m.tempName).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("name is required")
					}
					return nil
				}),

			huh.NewInput().
				Key("description").
				Title("Description").
				Description("What does this bento do?").
				Value(&m.tempDescription).
				CharLimit(200),
		).Title("Edit Metadata:"),
	).
		WithWidth(m.width - 48).
		WithShowHelp(false).
		WithShowErrors(false)
}

// updateDefinitionFromEditForm updates the definition from temp fields after editing
func (m *GuidedModal) updateDefinitionFromEditForm() {
	m.definition.Name = m.tempName
	m.definition.Description = m.tempDescription
}

// findNodeByName recursively searches for a node by name
func (m *GuidedModal) findNodeByName(name string) *neta.Definition {
	return findNodeInTree(m.definition.Nodes, name)
}

// findNodeInTree recursively searches for a node in a tree
func findNodeInTree(nodes []neta.Definition, name string) *neta.Definition {
	for i := range nodes {
		if nodes[i].Name == name {
			return &nodes[i]
		}

		// Search children if this is a group
		if len(nodes[i].Nodes) > 0 {
			if found := findNodeInTree(nodes[i].Nodes, name); found != nil {
				return found
			}
		}
	}
	return nil
}

// deleteNodeByName removes a node from the definition by name
func (m *GuidedModal) deleteNodeByName(name string) bool {
	return deleteNodeFromTree(&m.definition.Nodes, name)
}

// deleteNodeFromTree recursively deletes a node from a tree
func deleteNodeFromTree(nodes *[]neta.Definition, name string) bool {
	for i := range *nodes {
		if (*nodes)[i].Name == name {
			// Found it - delete it
			*nodes = append((*nodes)[:i], (*nodes)[i+1:]...)
			return true
		}

		// Search children if this is a group
		if len((*nodes)[i].Nodes) > 0 {
			if deleteNodeFromTree(&(*nodes)[i].Nodes, name) {
				return true
			}
		}
	}
	return false
}

// updateNodeInPlace updates a node's parameters in the definition
func (m *GuidedModal) updateNodeInPlace(name string, updatedNode *neta.Definition) bool {
	return updateNodeInTree(&m.definition.Nodes, name, updatedNode)
}

// updateNodeInTree recursively updates a node in a tree
func updateNodeInTree(nodes *[]neta.Definition, name string, updatedNode *neta.Definition) bool {
	for i := range *nodes {
		if (*nodes)[i].Name == name {
			// Found it - update it (preserve children if it's a group)
			children := (*nodes)[i].Nodes
			(*nodes)[i] = *updatedNode
			if len(children) > 0 {
				(*nodes)[i].Nodes = children
			}
			return true
		}

		// Search children if this is a group
		if len((*nodes)[i].Nodes) > 0 {
			if updateNodeInTree(&(*nodes)[i].Nodes, name, updatedNode) {
				return true
			}
		}
	}
	return false
}

// createHTTPNodeFormWithValues creates HTTP node form with pre-populated values
func (m *GuidedModal) createHTTPNodeFormWithValues() *huh.Form {
	formWidth := m.width - 40 - 8

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("node_name").
				Title("Name").
				Description("Name for this HTTP node").
				Value(&m.tempNodeName).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("name is required")
					}
					return nil
				}),

			huh.NewInput().
				Key("url").
				Title("URL").
				Description("HTTP(S) URL to request").
				Value(&m.tempNodeURL).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("url is required")
					}
					return nil
				}),

			huh.NewSelect[string]().
				Key("method").
				Title("Method").
				Description("HTTP method").
				Options(
					huh.NewOption("GET", "GET"),
					huh.NewOption("POST", "POST"),
					huh.NewOption("PUT", "PUT"),
					huh.NewOption("DELETE", "DELETE"),
					huh.NewOption("PATCH", "PATCH"),
					huh.NewOption("HEAD", "HEAD"),
					huh.NewOption("OPTIONS", "OPTIONS"),
				).
				Value(&m.tempNodeMethod),

			huh.NewInput().
				Key("headers").
				Title("Headers").
				Description("Request headers (optional, JSON object)").
				Value(&m.tempNodeHeaders).
				CharLimit(500),

			huh.NewInput().
				Key("body").
				Title("Body").
				Description("Request body (optional)").
				Value(&m.tempNodeBody).
				CharLimit(500),
		).Title("HTTP Node:"),
	).WithWidth(formWidth).WithShowHelp(false).WithShowErrors(false)
}

// createJQNodeFormWithValues creates JQ node form with pre-populated values
func (m *GuidedModal) createJQNodeFormWithValues() *huh.Form {
	formWidth := m.width - 40 - 8

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("node_name").
				Title("Name").
				Description("Name for this JQ transform node").
				Value(&m.tempNodeName).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("name is required")
					}
					return nil
				}),

			huh.NewInput().
				Key("query").
				Title("JQ Query").
				Description("JQ transformation query").
				Value(&m.tempNodeQuery).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("query is required")
					}
					return nil
				}).
				CharLimit(500),
		).Title("JQ Transform Node:"),
	).WithWidth(formWidth).WithShowHelp(false).WithShowErrors(false)
}

// createFileWriteNodeFormWithValues creates file.write node form with pre-populated values
func (m *GuidedModal) createFileWriteNodeFormWithValues() *huh.Form {
	formWidth := m.width - 40 - 8

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("node_name").
				Title("Name").
				Description("Name for this file write node").
				Value(&m.tempNodeName).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("name is required")
					}
					return nil
				}),

			huh.NewInput().
				Key("path").
				Title("Path").
				Description("File path to write to").
				Value(&m.tempNodePath).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("path is required")
					}
					return nil
				}),

			huh.NewInput().
				Key("content").
				Title("Content").
				Description("Content to write to file").
				Value(&m.tempNodeContent).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("content is required")
					}
					return nil
				}).
				CharLimit(500),
		).Title("File Write Node:"),
	).WithWidth(formWidth).WithShowHelp(false).WithShowErrors(false)
}

// createSequenceNodeFormWithValues creates sequence group form with pre-populated values
func (m *GuidedModal) createSequenceNodeFormWithValues() *huh.Form {
	formWidth := m.width - 40 - 8

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("node_name").
				Title("Name").
				Description("Name for this sequence group").
				Value(&m.tempNodeName).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("name is required")
					}
					return nil
				}),
		).Title("Sequence Group:"),
	).WithWidth(formWidth).WithShowHelp(false).WithShowErrors(false)
}

// createParallelNodeFormWithValues creates parallel group form with pre-populated values
func (m *GuidedModal) createParallelNodeFormWithValues() *huh.Form {
	formWidth := m.width - 40 - 8

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("node_name").
				Title("Name").
				Description("Name for this parallel group").
				Value(&m.tempNodeName).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("name is required")
					}
					return nil
				}),
		).Title("Parallel Group:"),
	).WithWidth(formWidth).WithShowHelp(false).WithShowErrors(false)
}

// updateCurrentNodeFromTempFields updates current node from temp fields (for edit mode)
func (m *GuidedModal) updateCurrentNodeFromTempFields(nodeType string) {
	if m.currentNode == nil {
		m.currentNode = &neta.Definition{
			Version:    "1.0",
			Parameters: make(map[string]interface{}),
		}
	}

	// Set node type and name from temp fields
	m.currentNode.Type = nodeType
	m.currentNode.Name = m.tempNodeName

	// Set parameters based on node type from temp fields
	switch nodeType {
	case "http":
		if m.tempNodeURL != "" {
			m.currentNode.Parameters["url"] = m.tempNodeURL
		}
		if m.tempNodeMethod != "" {
			m.currentNode.Parameters["method"] = m.tempNodeMethod
		}
		if m.tempNodeHeaders != "" {
			m.currentNode.Parameters["headers"] = m.tempNodeHeaders
		}
		if m.tempNodeBody != "" {
			m.currentNode.Parameters["body"] = m.tempNodeBody
		}

	case "transform.jq", "jq":
		if m.tempNodeQuery != "" {
			m.currentNode.Parameters["query"] = m.tempNodeQuery
		}

	case "file.write":
		if m.tempNodePath != "" {
			m.currentNode.Parameters["path"] = m.tempNodePath
		}
		if m.tempNodeContent != "" {
			m.currentNode.Parameters["content"] = m.tempNodeContent
		}

	case "group.sequence", "sequence":
		// Sequence nodes don't have parameters beyond name

	case "group.parallel", "parallel":
		// Parallel nodes don't have parameters beyond name
	}
}

// createNodeFormForTypeWithValues creates a node form pre-populated with values
func (m *GuidedModal) createNodeFormForTypeWithValues(nodeType string, node *neta.Definition) *huh.Form {
	// Pre-populate temp fields with existing node values
	m.tempNodeName = node.Name

	// Extract parameters based on node type
	switch nodeType {
	case "http":
		if url, ok := node.Parameters["url"].(string); ok {
			m.tempNodeURL = url
		}
		if method, ok := node.Parameters["method"].(string); ok {
			m.tempNodeMethod = method
		} else {
			m.tempNodeMethod = "GET" // Default
		}
		if headers, ok := node.Parameters["headers"].(string); ok {
			m.tempNodeHeaders = headers
		}
		if body, ok := node.Parameters["body"].(string); ok {
			m.tempNodeBody = body
		}
		return m.createHTTPNodeFormWithValues()

	case "transform.jq", "jq":
		if query, ok := node.Parameters["query"].(string); ok {
			m.tempNodeQuery = query
		}
		return m.createJQNodeFormWithValues()

	case "file.write":
		if path, ok := node.Parameters["path"].(string); ok {
			m.tempNodePath = path
		}
		if content, ok := node.Parameters["content"].(string); ok {
			m.tempNodeContent = content
		}
		return m.createFileWriteNodeFormWithValues()

	case "group.sequence", "sequence":
		return m.createSequenceNodeFormWithValues()

	case "group.parallel", "parallel":
		return m.createParallelNodeFormWithValues()

	default:
		// Fallback to empty form
		return m.createNodeFormForType(nodeType)
	}
}
