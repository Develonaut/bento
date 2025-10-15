// Package neta defines the core node types for Bento.
// Neta (ネタ) means "ingredients" or "toppings" in sushi terminology.
package neta

// Definition describes a node that can be executed by Itamae.
// It may be a single executable node or a group containing other nodes.
//
// Definitions are typically loaded from YAML/JSON configuration files and
// represent the declarative workflow specification. The orchestrator (itamae)
// converts definitions into executable operations by looking up node types
// in the registry (pantry).
type Definition struct {
	// Type identifies what kind of node this is (http, transform, group, etc)
	Type string

	// Name is the human-readable identifier for this node
	Name string

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
	// 1. Loading workflows from external configuration (YAML/JSON)
	// 2. Dynamic node type registration without core recompilation
	// 3. Schema validation at the node implementation level
	//
	// Parameter validation happens when the node is executed, not when the
	// definition is created, allowing late binding and runtime configuration.
	Parameters map[string]interface{}

	// Nodes contains child nodes (for group types)
	// If empty/nil, this is a leaf node
	Nodes []Definition
}

func (d Definition) IsGroup() bool {
	return len(d.Nodes) > 0
}
