package miso

import (
	"fmt"
	"strings"

	"github.com/Develonaut/bento/pkg/neta"
)

// flattenDefinition converts a bento definition tree into a flat list of node states.
// This recursively processes groups, loops, and parallel nodes to create a linear
// sequence suitable for display in the TUI.
// Note: Loop children are not shown individually since they execute multiple times.
func flattenDefinition(def neta.Definition, basePath string) []NodeState {
	if def.Type != "group" && def.Type != "loop" && def.Type != "parallel" {
		return flattenSingleNode(def, basePath)
	}
	return flattenGroupNodes(def, basePath)
}

// flattenSingleNode creates state for a single node bento.
func flattenSingleNode(def neta.Definition, basePath string) []NodeState {
	// Use node ID if present, otherwise use basePath (or "0" as fallback)
	path := basePath
	if def.ID != "" {
		path = def.ID
	} else if path == "" {
		path = "0" // Fallback for single nodes without ID
	}

	return []NodeState{{
		path:     path,
		name:     def.Name,
		nodeType: def.Type,
		status:   NodePending,
		depth:    0,
	}}
}

// flattenGroupNodes flattens all nodes in a group recursively.
func flattenGroupNodes(def neta.Definition, basePath string) []NodeState {
	states := []NodeState{}

	for idx, child := range def.Nodes {
		// Use node ID if present (graph-based execution), otherwise use hierarchical path
		path := getNodePath(child, basePath, idx)

		// Recursively flatten child groups/loops/parallel (don't track containers, only leaf nodes)
		if child.Type == "group" || child.Type == "loop" || child.Type == "parallel" {
			childStates := flattenDefinition(child, path)
			states = append(states, childStates...)
		} else {
			// Only track actual execution nodes (not containers)
			states = append(states, createNodeState(child, path))
		}
	}

	return states
}

// getNodePath returns the appropriate path for a node.
// Uses node ID for graph-based execution, falls back to hierarchical index.
func getNodePath(child neta.Definition, basePath string, idx int) string {
	if child.ID != "" {
		return child.ID // Use node ID for graph-based execution
	}
	return buildPathForNode(basePath, idx) // Use index-based path for hierarchical execution
}

// createNodeState builds a NodeState from definition and path.
func createNodeState(def neta.Definition, path string) NodeState {
	return NodeState{
		path:     path,
		name:     def.Name,
		nodeType: def.Type,
		status:   NodePending,
		depth:    parseDepth(path),
	}
}

// buildPathForNode constructs node path for tracking.
func buildPathForNode(basePath string, index int) string {
	if basePath == "" {
		return fmt.Sprintf("%d", index)
	}
	return fmt.Sprintf("%s.%d", basePath, index)
}

// parseDepth calculates nesting level from path.
func parseDepth(path string) int {
	if path == "" {
		return 0
	}
	return strings.Count(path, ".") + 1
}
