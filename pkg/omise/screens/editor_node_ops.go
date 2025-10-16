package screens

import (
	"fmt"

	"bento/pkg/neta"
)

// Node operation functions for building and manipulating bento definitions

// shouldSetAsRoot determines if the configured node should become the root node.
// Returns true when the definition has no nodes and no type set (first node).
func (e Editor) shouldSetAsRoot() bool {
	return len(e.def.Nodes) == 0 && e.def.Type == ""
}

// buildNode constructs a neta.Definition from a configured node message.
func buildNode(msg NodeConfiguredMsg) neta.Definition {
	return neta.Definition{
		Version:    neta.CurrentVersion,
		Type:       msg.Type,
		Name:       msg.Name,
		Parameters: msg.Parameters,
	}
}

// setRootNode sets the configured node as the root node of the definition.
// Used when adding the first node to an empty bento.
func setRootNode(def neta.Definition, msg NodeConfiguredMsg) neta.Definition {
	def.Type = msg.Type
	def.Parameters = msg.Parameters
	return def
}

// appendNode adds a node to the definition's Nodes array.
// If the definition has no type, it converts it to a group.sequence.
func appendNode(def neta.Definition, node neta.Definition) neta.Definition {
	if def.Type == "" {
		def.Type = "group.sequence"
	}
	def.Nodes = append(def.Nodes, node)
	return def
}

// defaultNodeConfig returns a default configuration when no schema is available.
func defaultNodeConfig(nodeType string) NodeConfiguredMsg {
	return NodeConfiguredMsg{
		Type:       nodeType,
		Name:       fmt.Sprintf("New %s Node", nodeType),
		Parameters: map[string]interface{}{},
	}
}

// extractNodeName extracts and removes the "name" field from parameters.
// Returns a default name if not found.
func extractNodeName(params map[string]interface{}, nodeType string) string {
	if name, ok := params["name"]; ok {
		if nameStr, ok := name.(*string); ok {
			delete(params, "name")
			return *nameStr
		}
	}
	return fmt.Sprintf("New %s Node", nodeType)
}

// convertParamPointers converts pointer values in parameters to actual values.
// Huh forms use pointers for binding, but we need actual values for storage.
// Uses interface{} because parameters are dynamically typed from schema definitions.
func convertParamPointers(params map[string]interface{}) map[string]interface{} {
	actualParams := make(map[string]interface{})
	for k, v := range params {
		actualParams[k] = derefValue(v)
	}
	return actualParams
}

// derefValue dereferences pointer values to their underlying types.
// Used for converting Huh form pointer bindings to actual values.
// interface{} required for type-agnostic pointer dereferencing.
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
