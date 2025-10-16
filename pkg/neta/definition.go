// Package neta defines the core node types for Bento.
// Neta (ネタ) means "ingredients" or "toppings" in sushi terminology.
package neta

// NodeEdge defines a connection between two nodes in a bento graph.
// Edges determine execution order and data flow between nodes.
type NodeEdge struct {
	// ID is the unique identifier for this edge
	ID string `yaml:"id" json:"id"`

	// Source is the ID of the node this edge originates from
	Source string `yaml:"source" json:"source"`

	// Target is the ID of the node this edge points to
	Target string `yaml:"target" json:"target"`

	// SourceHandle identifies which output port of the source node (optional)
	SourceHandle string `yaml:"sourceHandle,omitempty" json:"sourceHandle,omitempty"`

	// TargetHandle identifies which input port of the target node (optional)
	TargetHandle string `yaml:"targetHandle,omitempty" json:"targetHandle,omitempty"`
}

// Definition describes a node that can be executed by Itamae.
// It may be a single executable node or a group containing other nodes.
//
// Definitions are typically loaded from YAML/JSON configuration files and
// represent the declarative bento specification. The orchestrator (itamae)
// converts definitions into executable operations by looking up node types
// in the registry (pantry).
type Definition struct {
	// ID is the unique identifier for this node
	// Optional for simple bentos, but required for graph-based execution
	ID string `yaml:"id,omitempty" json:"id,omitempty"`

	// ParentID references the parent node (for nested/hierarchical display)
	// Empty for root nodes
	ParentID string `yaml:"parentId,omitempty" json:"parentId,omitempty"`
	// Version specifies the definition schema version (e.g., "1.0")
	// REQUIRED: Must be present in all .bento.yaml files
	// Format: MAJOR.MINOR (semantic versioning)
	Version string `yaml:"version" json:"version"`

	// Type identifies what kind of node this is (http, transform, group, etc)
	Type string `yaml:"type" json:"type"`

	// Name is the human-readable identifier for this node
	Name string `yaml:"name" json:"name"`

	// Parameters contains type-specific configuration.
	//
	// Parameters uses map[string]interface{} to support heterogeneous node
	// types loaded from YAML/JSON configuration files. Each node type expects
	// different parameter schemas:
	// - HTTP nodes: {url: string, method: string, headers: map}
	// - File nodes: {path: string, operation: string}
	// - Shell nodes: {command: string, args: []string}
	// - Transform nodes: {input: any, script: string}
	//
	// This flexibility enables:
	// 1. Loading bentos from external configuration (YAML/JSON)
	// 2. Dynamic node type registration without core recompilation
	// 3. Schema validation at the node implementation level
	//
	// Parameter validation happens when the node is executed, not when the
	// definition is created, allowing late binding and runtime configuration.
	Parameters map[string]interface{} `yaml:"parameters,omitempty" json:"parameters,omitempty"`

	// Nodes contains child nodes (for group types)
	// In graph-based execution, all nodes are stored flat in this array
	// and connected via Edges rather than nested hierarchically
	Nodes []Definition `yaml:"nodes,omitempty" json:"nodes,omitempty"`

	// Edges defines connections between child nodes (for group types)
	// Edges determine execution order and data flow in graph-based execution
	Edges []NodeEdge `yaml:"edges,omitempty" json:"edges,omitempty"`
}

// CurrentVersion is the version of definitions this build supports
const CurrentVersion = "1.0"

// IsGroup returns true if this definition contains child nodes
func (d Definition) IsGroup() bool {
	return len(d.Nodes) > 0
}

// IsVersionCompatible checks if the definition version is compatible
func (d Definition) IsVersionCompatible() bool {
	return isCompatibleVersion(d.Version)
}
