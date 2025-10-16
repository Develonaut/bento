package neta

import (
	"testing"
)

func TestAnalyzeExecutionGraph_SingleNode(t *testing.T) {
	def := Definition{
		ID:      "node-1",
		Type:    "http",
		Name:    "Test Node",
		Version: "1.0",
	}

	graph, err := AnalyzeExecutionGraph(def)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(graph.Nodes) != 1 {
		t.Errorf("Expected 1 node, got %d", len(graph.Nodes))
	}

	if len(graph.ExecutionOrder) != 1 {
		t.Errorf("Expected 1 execution layer, got %d", len(graph.ExecutionOrder))
	}

	if graph.TotalWeight != 1 {
		t.Errorf("Expected total weight 1, got %d", graph.TotalWeight)
	}
}

func TestAnalyzeExecutionGraph_SequentialNodes(t *testing.T) {
	def := Definition{
		ID:      "root",
		Type:    "group.sequence",
		Name:    "Sequential Workflow",
		Version: "1.0",
		Nodes: []Definition{
			{ID: "node-1", Type: "http", Name: "Node 1", Version: "1.0"},
			{ID: "node-2", Type: "transform.jq", Name: "Node 2", Version: "1.0"},
			{ID: "node-3", Type: "http", Name: "Node 3", Version: "1.0"},
		},
		Edges: []NodeEdge{
			{ID: "edge-1-2", Source: "node-1", Target: "node-2"},
			{ID: "edge-2-3", Source: "node-2", Target: "node-3"},
		},
	}

	graph, err := AnalyzeExecutionGraph(def)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(graph.Nodes) != 3 {
		t.Errorf("Expected 3 nodes, got %d", len(graph.Nodes))
	}

	if len(graph.ExecutionOrder) != 3 {
		t.Errorf("Expected 3 execution layers (sequential), got %d", len(graph.ExecutionOrder))
	}

	// Verify sequential order
	if len(graph.ExecutionOrder[0]) != 1 || graph.ExecutionOrder[0][0] != "node-1" {
		t.Errorf("Expected first layer to contain node-1")
	}
	if len(graph.ExecutionOrder[1]) != 1 || graph.ExecutionOrder[1][0] != "node-2" {
		t.Errorf("Expected second layer to contain node-2")
	}
	if len(graph.ExecutionOrder[2]) != 1 || graph.ExecutionOrder[2][0] != "node-3" {
		t.Errorf("Expected third layer to contain node-3")
	}

	if graph.MaxParallelism != 1 {
		t.Errorf("Expected max parallelism 1 (sequential), got %d", graph.MaxParallelism)
	}
}

func TestAnalyzeExecutionGraph_ParallelNodes(t *testing.T) {
	def := Definition{
		ID:      "root",
		Type:    "group.parallel",
		Name:    "Parallel Workflow",
		Version: "1.0",
		Nodes: []Definition{
			{ID: "node-1", Type: "http", Name: "Node 1", Version: "1.0"},
			{ID: "node-2", Type: "http", Name: "Node 2", Version: "1.0"},
			{ID: "node-3", Type: "http", Name: "Node 3", Version: "1.0"},
		},
		Edges: []NodeEdge{}, // No edges = all parallel
	}

	graph, err := AnalyzeExecutionGraph(def)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(graph.Nodes) != 3 {
		t.Errorf("Expected 3 nodes, got %d", len(graph.Nodes))
	}

	if len(graph.ExecutionOrder) != 1 {
		t.Errorf("Expected 1 execution layer (all parallel), got %d", len(graph.ExecutionOrder))
	}

	if len(graph.ExecutionOrder[0]) != 3 {
		t.Errorf("Expected 3 nodes in first layer, got %d", len(graph.ExecutionOrder[0]))
	}

	if graph.MaxParallelism != 3 {
		t.Errorf("Expected max parallelism 3, got %d", graph.MaxParallelism)
	}
}

func TestAnalyzeExecutionGraph_CyclicDependency(t *testing.T) {
	def := Definition{
		ID:      "root",
		Type:    "group.sequence",
		Name:    "Cyclic Workflow",
		Version: "1.0",
		Nodes: []Definition{
			{ID: "node-1", Type: "http", Name: "Node 1", Version: "1.0"},
			{ID: "node-2", Type: "http", Name: "Node 2", Version: "1.0"},
		},
		Edges: []NodeEdge{
			{ID: "edge-1-2", Source: "node-1", Target: "node-2"},
			{ID: "edge-2-1", Source: "node-2", Target: "node-1"}, // Creates cycle
		},
	}

	_, err := AnalyzeExecutionGraph(def)
	if err == nil {
		t.Fatal("Expected error for cyclic dependency, got none")
	}
}

func TestAnalyzeExecutionGraph_DiamondShape(t *testing.T) {
	// Diamond: node-1 -> node-2, node-3 -> node-4
	//                  \-> node-3 /
	def := Definition{
		ID:      "root",
		Type:    "group.sequence",
		Name:    "Diamond Workflow",
		Version: "1.0",
		Nodes: []Definition{
			{ID: "node-1", Type: "http", Name: "Start", Version: "1.0"},
			{ID: "node-2", Type: "http", Name: "Branch A", Version: "1.0"},
			{ID: "node-3", Type: "http", Name: "Branch B", Version: "1.0"},
			{ID: "node-4", Type: "http", Name: "End", Version: "1.0"},
		},
		Edges: []NodeEdge{
			{ID: "edge-1-2", Source: "node-1", Target: "node-2"},
			{ID: "edge-1-3", Source: "node-1", Target: "node-3"},
			{ID: "edge-2-4", Source: "node-2", Target: "node-4"},
			{ID: "edge-3-4", Source: "node-3", Target: "node-4"},
		},
	}

	graph, err := AnalyzeExecutionGraph(def)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(graph.Nodes) != 4 {
		t.Errorf("Expected 4 nodes, got %d", len(graph.Nodes))
	}

	if len(graph.ExecutionOrder) != 3 {
		t.Errorf("Expected 3 execution layers, got %d", len(graph.ExecutionOrder))
	}

	// Layer 1: node-1
	if len(graph.ExecutionOrder[0]) != 1 {
		t.Errorf("Expected 1 node in first layer, got %d", len(graph.ExecutionOrder[0]))
	}

	// Layer 2: node-2 and node-3 (parallel)
	if len(graph.ExecutionOrder[1]) != 2 {
		t.Errorf("Expected 2 nodes in second layer (parallel branches), got %d", len(graph.ExecutionOrder[1]))
	}

	// Layer 3: node-4
	if len(graph.ExecutionOrder[2]) != 1 {
		t.Errorf("Expected 1 node in third layer, got %d", len(graph.ExecutionOrder[2]))
	}

	if graph.MaxParallelism != 2 {
		t.Errorf("Expected max parallelism 2 (diamond middle), got %d", graph.MaxParallelism)
	}
}
