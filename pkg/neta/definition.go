// Package neta defines the core node types for Bento.
// Neta (ネタ) means "ingredients" or "toppings" in sushi terminology.
package neta

// Definition describes a node that can be executed by Itamae.
// It may be a single executable node or a group containing other nodes.
//
// Definitions are typically loaded from YAML/JSON configuration files and
// represent the declarative bento specification. The orchestrator (itamae)
// converts definitions into executable operations by looking up node types
// in the registry (pantry).
type Definition struct {
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
	// If empty/nil, this is a leaf node
	Nodes []Definition `yaml:"nodes,omitempty" json:"nodes,omitempty"`
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
