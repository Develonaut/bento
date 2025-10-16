package neta

import (
	"fmt"
)

// ExecutionGraph represents the analyzed graph structure
type ExecutionGraph struct {
	Nodes          map[string]*ExecutionGraphNode
	Edges          []ExecutionGraphEdge
	ExecutionOrder [][]string // Layers of nodes that can execute in parallel
	CriticalPath   []string   // Longest path through graph
	TotalWeight    int
	MaxParallelism int
}

// AnalyzeExecutionGraph converts a Definition to an execution graph
func AnalyzeExecutionGraph(def Definition) (*ExecutionGraph, error) {
	// For non-group nodes, create simple single-node graph
	if !def.IsGroup() {
		node := &ExecutionGraphNode{
			ID:           def.ID,
			Name:         def.Name,
			Type:         def.Type,
			Weight:       1,
			Dependencies: []string{},
		}
		return &ExecutionGraph{
			Nodes: map[string]*ExecutionGraphNode{
				def.ID: node,
			},
			Edges:          []ExecutionGraphEdge{},
			ExecutionOrder: [][]string{{def.ID}},
			CriticalPath:   []string{def.ID},
			TotalWeight:    1,
			MaxParallelism: 1,
		}, nil
	}

	// Build graph from edges
	nodes := make(map[string]*ExecutionGraphNode)
	for _, childDef := range def.Nodes {
		nodes[childDef.ID] = &ExecutionGraphNode{
			ID:           childDef.ID,
			Name:         childDef.Name,
			Type:         childDef.Type,
			Weight:       1,
			Dependencies: []string{},
		}
	}

	// Build dependency map from edges
	edges := make([]ExecutionGraphEdge, 0, len(def.Edges))
	for _, edge := range def.Edges {
		edges = append(edges, ExecutionGraphEdge{
			From: edge.Source,
			To:   edge.Target,
		})

		// Add dependency
		if targetNode, ok := nodes[edge.Target]; ok {
			targetNode.Dependencies = append(targetNode.Dependencies, edge.Source)
		}
	}

	// Detect cycles
	if hasCycle(nodes, edges) {
		return nil, fmt.Errorf("circular dependency detected in execution graph")
	}

	// Topological sort to get execution order
	executionOrder, err := topologicalSort(nodes, edges)
	if err != nil {
		return nil, err
	}

	// Calculate critical path (longest path)
	criticalPath := calculateCriticalPath(nodes, executionOrder)

	// Calculate total weight and max parallelism
	totalWeight := len(nodes)
	maxParallelism := 1
	for _, layer := range executionOrder {
		if len(layer) > maxParallelism {
			maxParallelism = len(layer)
		}
	}

	return &ExecutionGraph{
		Nodes:          nodes,
		Edges:          edges,
		ExecutionOrder: executionOrder,
		CriticalPath:   criticalPath,
		TotalWeight:    totalWeight,
		MaxParallelism: maxParallelism,
	}, nil
}

// topologicalSort performs Kahn's algorithm for topological sorting
func topologicalSort(nodes map[string]*ExecutionGraphNode, edges []ExecutionGraphEdge) ([][]string, error) {
	// Calculate in-degree for each node
	inDegree := make(map[string]int)
	for id := range nodes {
		inDegree[id] = 0
	}
	for _, edge := range edges {
		inDegree[edge.To]++
	}

	// Find all nodes with in-degree 0 (no dependencies)
	queue := []string{}
	for id, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, id)
		}
	}

	layers := [][]string{}
	visited := 0

	for len(queue) > 0 {
		// All nodes in queue can execute in parallel
		layer := make([]string, len(queue))
		copy(layer, queue)
		layers = append(layers, layer)

		// Process current layer
		nextQueue := []string{}
		for _, nodeID := range queue {
			visited++

			// Reduce in-degree for all dependent nodes
			for _, edge := range edges {
				if edge.From == nodeID {
					inDegree[edge.To]--
					if inDegree[edge.To] == 0 {
						nextQueue = append(nextQueue, edge.To)
					}
				}
			}
		}

		queue = nextQueue
	}

	// Check if all nodes were visited (no cycles)
	if visited != len(nodes) {
		return nil, fmt.Errorf("graph contains a cycle")
	}

	return layers, nil
}

// hasCycle detects if there's a cycle in the graph using DFS
func hasCycle(nodes map[string]*ExecutionGraphNode, edges []ExecutionGraphEdge) bool {
	// Build adjacency list
	adjacency := make(map[string][]string)
	for _, edge := range edges {
		adjacency[edge.From] = append(adjacency[edge.From], edge.To)
	}

	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	var hasCycleDFS func(string) bool
	hasCycleDFS = func(nodeID string) bool {
		visited[nodeID] = true
		recStack[nodeID] = true

		for _, neighbor := range adjacency[nodeID] {
			if !visited[neighbor] {
				if hasCycleDFS(neighbor) {
					return true
				}
			} else if recStack[neighbor] {
				return true
			}
		}

		recStack[nodeID] = false
		return false
	}

	for nodeID := range nodes {
		if !visited[nodeID] {
			if hasCycleDFS(nodeID) {
				return true
			}
		}
	}

	return false
}

// calculateCriticalPath finds the longest path through the graph
func calculateCriticalPath(nodes map[string]*ExecutionGraphNode, executionOrder [][]string) []string {
	// For now, simply return the first node from each layer
	// A more sophisticated implementation would calculate actual weights
	path := []string{}
	for _, layer := range executionOrder {
		if len(layer) > 0 {
			path = append(path, layer[0])
		}
	}
	return path
}
