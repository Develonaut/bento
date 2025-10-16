package jubako

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"bento/pkg/neta"
)

// Parser handles .bento.yaml file parsing.
type Parser struct{}

// NewParser creates a new parser.
func NewParser() *Parser {
	return &Parser{}
}

// Parse reads and parses a .bento.yaml file.
func (p *Parser) Parse(path string) (neta.Definition, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return neta.Definition{}, fmt.Errorf("read failed: %w", err)
	}

	return p.ParseBytes(data)
}

// ParseBytes parses .bento.yaml from bytes.
func (p *Parser) ParseBytes(data []byte) (neta.Definition, error) {
	var def neta.Definition
	if err := yaml.Unmarshal(data, &def); err != nil {
		return neta.Definition{}, fmt.Errorf("invalid YAML: %w", err)
	}

	// Normalize group definitions (extract child nodes from parameters)
	def = normalizeDefinition(def)

	// Assign IDs to nodes that don't have them
	def = assignNodeIDs(def)

	// Auto-generate edges if missing (backward compatibility)
	def = autoGenerateEdges(def)

	if err := validateDefinition(def); err != nil {
		return neta.Definition{}, fmt.Errorf("validation failed: %w", err)
	}

	return def, nil
}

// Format converts a definition to YAML.
func (p *Parser) Format(def neta.Definition) ([]byte, error) {
	data, err := yaml.Marshal(def)
	if err != nil {
		return nil, fmt.Errorf("marshal failed: %w", err)
	}
	return data, nil
}

// validateDefinition ensures a definition is well-formed.
func validateDefinition(def neta.Definition) error {
	// Validate version and type first
	if err := validateStructure(def); err != nil {
		return err
	}

	// Validate node parameters using validation framework
	validator := neta.NewValidator()
	if err := validator.ValidateRecursive(def); err != nil {
		return err
	}

	return nil
}

// normalizeDefinition extracts child nodes from parameters for group types
func normalizeDefinition(def neta.Definition) neta.Definition {
	// Check if this is a group type with nodes in parameters
	if isGroupType(def.Type) && len(def.Nodes) == 0 {
		if nodes, ok := extractNodesFromParams(def.Parameters); ok {
			def.Nodes = nodes
		}
	}

	// Recursively normalize child nodes
	for i := range def.Nodes {
		def.Nodes[i] = normalizeDefinition(def.Nodes[i])
	}

	return def
}

// isGroupType checks if a type is a group orchestration type
func isGroupType(nodeType string) bool {
	return nodeType == "group.sequence" ||
		nodeType == "group.parallel" ||
		nodeType == "loop.for" ||
		nodeType == "conditional.if"
}

// extractNodesFromParams attempts to extract child nodes from parameters
func extractNodesFromParams(params map[string]interface{}) ([]neta.Definition, bool) {
	nodesParam, ok := params["nodes"]
	if !ok {
		return nil, false
	}

	// Handle []interface{} from YAML unmarshaling
	nodesSlice, ok := nodesParam.([]interface{})
	if !ok {
		return nil, false
	}

	nodes := make([]neta.Definition, 0, len(nodesSlice))
	for _, item := range nodesSlice {
		// Convert map to Definition via re-marshaling
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		// Re-marshal and unmarshal to get proper Definition
		data, err := yaml.Marshal(itemMap)
		if err != nil {
			continue
		}

		var node neta.Definition
		if err := yaml.Unmarshal(data, &node); err != nil {
			continue
		}

		nodes = append(nodes, node)
	}

	return nodes, len(nodes) > 0
}

// validateStructure recursively validates version and type of a definition and its children.
func validateStructure(def neta.Definition) error {
	// Validate version
	if err := neta.ValidateVersion(def.Version); err != nil {
		return fmt.Errorf("version error: %w", err)
	}

	// Validate type is present
	if def.Type == "" {
		return fmt.Errorf("type is required")
	}

	// Recursively validate child nodes
	for i, child := range def.Nodes {
		if err := validateStructure(child); err != nil {
			return fmt.Errorf("node %d: %w", i, err)
		}
	}

	return nil
}

// assignNodeIDs recursively assigns unique IDs to nodes that don't have them
func assignNodeIDs(def neta.Definition) neta.Definition {
	return assignNodeIDsWithCounter(def, 1)
}

// assignNodeIDsWithCounter assigns IDs using a counter
func assignNodeIDsWithCounter(def neta.Definition, counter int) neta.Definition {
	// Assign ID to root if missing
	if def.ID == "" {
		def.ID = fmt.Sprintf("node-%d", counter)
		counter++
	}

	// Recursively assign IDs to children
	for i := range def.Nodes {
		if def.Nodes[i].ID == "" {
			def.Nodes[i].ID = fmt.Sprintf("node-%d", counter)
			counter++
		}
		// Recursively process nested nodes
		def.Nodes[i] = assignNodeIDsWithCounter(def.Nodes[i], counter)
	}

	return def
}

// autoGenerateEdges generates edges for group nodes that don't have explicit edges
func autoGenerateEdges(def neta.Definition) neta.Definition {
	// Only generate edges for group types with nodes but no edges
	if isGroupType(def.Type) && len(def.Nodes) > 0 && len(def.Edges) == 0 {
		// Generate sequential edges for sequence and loop types
		if def.Type == "group.sequence" || def.Type == "loop.for" {
			def.Edges = generateSequentialEdges(def.Nodes)
		}
		// For parallel groups, no edges needed (all run in parallel)
		// For conditional.if, edges would be generated based on condition logic
	}

	// Recursively process child nodes
	for i := range def.Nodes {
		def.Nodes[i] = autoGenerateEdges(def.Nodes[i])
	}

	return def
}

// generateSequentialEdges creates edges connecting nodes in sequence
func generateSequentialEdges(nodes []neta.Definition) []neta.NodeEdge {
	if len(nodes) <= 1 {
		return nil
	}

	edges := make([]neta.NodeEdge, 0, len(nodes)-1)
	for i := 0; i < len(nodes)-1; i++ {
		edge := neta.NodeEdge{
			ID:     fmt.Sprintf("edge-%d-%d", i+1, i+2),
			Source: nodes[i].ID,
			Target: nodes[i+1].ID,
		}
		edges = append(edges, edge)
	}
	return edges
}
