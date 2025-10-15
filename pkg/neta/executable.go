package neta

import "context"

// Result represents the outcome of executing a node.
//
// Result uses interface{} for Output to support heterogeneous node types
// in the orchestration system. Each node type produces different output types:
// - HTTP nodes return response bodies ([]byte, string, or parsed JSON)
// - File nodes return file contents or paths
// - Transform nodes return modified data structures
// - Shell nodes return command output and exit codes
//
// The orchestrator (itamae) treats these as opaque values and passes them
// between nodes. Type assertions happen at the node implementation level
// where the concrete types are known.
//
// This design enables a plugin-based architecture where node types can be
// registered dynamically without recompiling the core engine.
type Result struct {
	// Output contains the result data. Type varies by node implementation.
	// Common types: string, []byte, map[string]interface{}, int, bool.
	// Nil indicates no output (e.g., write-only operations).
	Output interface{}
}

// Executable is implemented by all node types that can be executed.
// Accept interfaces, return structs.
type Executable interface {
	// Execute runs the node and returns its result.
	//
	// The params map uses interface{} values to support heterogeneous node
	// configurations. Each node type expects different parameter types:
	// - HTTP nodes: "url" (string), "method" (string), "body" ([]byte)
	// - File nodes: "path" (string), "mode" (int)
	// - Shell nodes: "command" (string), "args" ([]string)
	// - Transform nodes: "script" (string), "input" (interface{})
	//
	// Node implementations validate and type-assert parameters, returning
	// errors for invalid or missing required parameters.
	//
	// This design enables dynamic node configuration from YAML/JSON
	// without requiring compile-time knowledge of all node types.
	Execute(ctx context.Context, params map[string]interface{}) (Result, error)
}

// Executor can execute neta definitions.
// Used by nodes that need to orchestrate other nodes (conditional, loop, group).
type Executor interface {
	Execute(ctx context.Context, def Definition) (Result, error)
}
