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
	if !def.IsGroup() {
		return createSingleNodeGraph(def), nil
	}

	nodes := buildGraphNodes(def.Nodes)
	edges := buildGraphEdges(def.Edges, nodes)

	if hasCycle(nodes, edges) {
		return nil, fmt.Errorf("circular dependency detected in execution graph")
	}

	executionOrder, err := topologicalSort(nodes, edges)
	if err != nil {
		return nil, err
	}

	return buildExecutionGraph(nodes, edges, executionOrder), nil
}

// createSingleNodeGraph creates a graph for a single non-group node
func createSingleNodeGraph(def Definition) *ExecutionGraph {
	node := &ExecutionGraphNode{
		ID:           def.ID,
		Name:         def.Name,
		Type:         def.Type,
		Weight:       1,
		Dependencies: []string{},
	}
	return &ExecutionGraph{
		Nodes:          map[string]*ExecutionGraphNode{def.ID: node},
		Edges:          []ExecutionGraphEdge{},
		ExecutionOrder: [][]string{{def.ID}},
		CriticalPath:   []string{def.ID},
		TotalWeight:    1,
		MaxParallelism: 1,
	}
}

// buildGraphNodes creates nodes map from child definitions
func buildGraphNodes(children []Definition) map[string]*ExecutionGraphNode {
	nodes := make(map[string]*ExecutionGraphNode)
	for _, childDef := range children {
		nodes[childDef.ID] = &ExecutionGraphNode{
			ID:           childDef.ID,
			Name:         childDef.Name,
			Type:         childDef.Type,
			Weight:       1,
			Dependencies: []string{},
		}
	}
	return nodes
}

// buildGraphEdges creates edges and populates dependencies
func buildGraphEdges(defEdges []NodeEdge, nodes map[string]*ExecutionGraphNode) []ExecutionGraphEdge {
	edges := make([]ExecutionGraphEdge, 0, len(defEdges))
	for _, edge := range defEdges {
		edges = append(edges, ExecutionGraphEdge{
			From: edge.Source,
			To:   edge.Target,
		})
		if targetNode, ok := nodes[edge.Target]; ok {
			targetNode.Dependencies = append(targetNode.Dependencies, edge.Source)
		}
	}
	return edges
}

// buildExecutionGraph assembles final graph with metrics
func buildExecutionGraph(nodes map[string]*ExecutionGraphNode, edges []ExecutionGraphEdge, executionOrder [][]string) *ExecutionGraph {
	criticalPath := calculateCriticalPath(nodes, executionOrder)
	totalWeight, maxParallelism := calculateGraphMetrics(nodes, executionOrder)

	return &ExecutionGraph{
		Nodes:          nodes,
		Edges:          edges,
		ExecutionOrder: executionOrder,
		CriticalPath:   criticalPath,
		TotalWeight:    totalWeight,
		MaxParallelism: maxParallelism,
	}
}

// calculateGraphMetrics computes total weight and max parallelism
func calculateGraphMetrics(nodes map[string]*ExecutionGraphNode, executionOrder [][]string) (int, int) {
	totalWeight := len(nodes)
	maxParallelism := 1
	for _, layer := range executionOrder {
		if len(layer) > maxParallelism {
			maxParallelism = len(layer)
		}
	}
	return totalWeight, maxParallelism
}

// topologicalSort performs Kahn's algorithm for topological sorting
func topologicalSort(nodes map[string]*ExecutionGraphNode, edges []ExecutionGraphEdge) ([][]string, error) {
	inDegree := calculateInDegree(nodes, edges)
	queue := findNodesWithNoDependencies(inDegree)
	layers, visited := processTopologicalLayers(queue, edges, inDegree)

	if visited != len(nodes) {
		return nil, fmt.Errorf("graph contains a cycle")
	}
	return layers, nil
}

// calculateInDegree computes in-degree for all nodes
func calculateInDegree(nodes map[string]*ExecutionGraphNode, edges []ExecutionGraphEdge) map[string]int {
	inDegree := make(map[string]int)
	for id := range nodes {
		inDegree[id] = 0
	}
	for _, edge := range edges {
		inDegree[edge.To]++
	}
	return inDegree
}

// findNodesWithNoDependencies returns nodes with in-degree 0
func findNodesWithNoDependencies(inDegree map[string]int) []string {
	queue := []string{}
	for id, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, id)
		}
	}
	return queue
}

// processTopologicalLayers builds execution layers using Kahn's algorithm
func processTopologicalLayers(queue []string, edges []ExecutionGraphEdge, inDegree map[string]int) ([][]string, int) {
	layers := [][]string{}
	visited := 0

	for len(queue) > 0 {
		layer := make([]string, len(queue))
		copy(layer, queue)
		layers = append(layers, layer)

		nextQueue := processLayer(queue, edges, inDegree, &visited)
		queue = nextQueue
	}
	return layers, visited
}

// processLayer processes one layer and returns next queue
func processLayer(queue []string, edges []ExecutionGraphEdge, inDegree map[string]int, visited *int) []string {
	nextQueue := []string{}
	for _, nodeID := range queue {
		*visited++
		for _, edge := range edges {
			if edge.From == nodeID {
				inDegree[edge.To]--
				if inDegree[edge.To] == 0 {
					nextQueue = append(nextQueue, edge.To)
				}
			}
		}
	}
	return nextQueue
}

// hasCycle detects if there's a cycle in the graph using DFS
func hasCycle(nodes map[string]*ExecutionGraphNode, edges []ExecutionGraphEdge) bool {
	adjacency := buildAdjacencyList(edges)
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for nodeID := range nodes {
		if !visited[nodeID] && detectCycleDFS(nodeID, adjacency, visited, recStack) {
			return true
		}
	}
	return false
}

// buildAdjacencyList creates adjacency list from edges
func buildAdjacencyList(edges []ExecutionGraphEdge) map[string][]string {
	adjacency := make(map[string][]string)
	for _, edge := range edges {
		adjacency[edge.From] = append(adjacency[edge.From], edge.To)
	}
	return adjacency
}

// detectCycleDFS performs DFS to detect cycles
func detectCycleDFS(nodeID string, adjacency map[string][]string, visited, recStack map[string]bool) bool {
	visited[nodeID] = true
	recStack[nodeID] = true

	for _, neighbor := range adjacency[nodeID] {
		if !visited[neighbor] {
			if detectCycleDFS(neighbor, adjacency, visited, recStack) {
				return true
			}
		} else if recStack[neighbor] {
			return true
		}
	}

	recStack[nodeID] = false
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
