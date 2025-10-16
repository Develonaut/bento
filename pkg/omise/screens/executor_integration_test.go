package screens

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"

	"bento/pkg/neta"
	"bento/pkg/omise/components"
)

// executorTestModel wraps Executor to implement tea.Model for testing
type executorTestModel struct {
	executor Executor
}

func (m *executorTestModel) Init() tea.Cmd {
	return m.executor.Init()
}

func (m *executorTestModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle quit for testing
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}

	newExec, cmd := m.executor.Update(msg)
	m.executor = newExec
	return m, cmd
}

func (m *executorTestModel) View() string {
	return m.executor.View()
}

// TestExecutorNodeDisplay_GraphBased tests that graph-based bentos display all nodes
func TestExecutorNodeDisplay_GraphBased(t *testing.T) {
	// Create a graph-based bento definition (like hello-world.bento.yaml)
	def := neta.Definition{
		Version: "1.0",
		Type:    "group.sequence",
		Name:    "Test Workflow",
		Nodes: []neta.Definition{
			{
				ID:      "node-1",
				Version: "1.0",
				Type:    "http",
				Name:    "First Node",
				Parameters: map[string]interface{}{
					"method": "GET",
					"url":    "https://httpbin.org/json",
				},
			},
			{
				ID:      "node-2",
				Version: "1.0",
				Type:    "transform.jq",
				Name:    "Second Node",
				Parameters: map[string]interface{}{
					"query": ".slideshow.title",
				},
			},
			{
				ID:      "node-3",
				Version: "1.0",
				Type:    "http",
				Name:    "Third Node",
				Parameters: map[string]interface{}{
					"method": "POST",
					"url":    "https://httpbin.org/post",
				},
			},
		},
		Edges: []neta.NodeEdge{
			{ID: "edge-1-2", Source: "node-1", Target: "node-2"},
			{ID: "edge-2-3", Source: "node-2", Target: "node-3"},
		},
	}

	// Create executor and start execution (simulates StartBento call)
	executor := NewExecutor()
	executor = executor.StartBento("test-workflow", "/test/path.yaml", "/test/workdir")
	tm := teatest.NewTestModel(t, &executorTestModel{executor: executor}, teatest.WithInitialTermSize(120, 40))

	// Send ExecutionInitMsg (simulates message from background goroutine)
	tm.Send(ExecutionInitMsg{Definition: def})
	time.Sleep(50 * time.Millisecond)

	// Send node started messages
	tm.Send(NodeStartedMsg{Path: "node-1", Name: "First Node", NodeType: "http"})
	time.Sleep(50 * time.Millisecond)

	tm.Send(NodeCompletedMsg{Path: "node-1", Duration: 100 * time.Millisecond, Error: nil})
	time.Sleep(50 * time.Millisecond)

	tm.Send(NodeStartedMsg{Path: "node-2", Name: "Second Node", NodeType: "transform.jq"})
	time.Sleep(50 * time.Millisecond)

	tm.Send(NodeCompletedMsg{Path: "node-2", Duration: 50 * time.Millisecond, Error: nil})
	time.Sleep(50 * time.Millisecond)

	tm.Send(NodeStartedMsg{Path: "node-3", Name: "Third Node", NodeType: "http"})
	time.Sleep(50 * time.Millisecond)

	tm.Send(NodeCompletedMsg{Path: "node-3", Duration: 200 * time.Millisecond, Error: nil})
	time.Sleep(50 * time.Millisecond)

	// Send quit to finish
	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})

	// Wait for test to finish
	tm.WaitFinished(t, teatest.WithFinalTimeout(2*time.Second))

	// Get final model
	finalModel := tm.FinalModel(t)
	finalTestModel, ok := finalModel.(*executorTestModel)
	if !ok {
		t.Fatal("Final model is not *executorTestModel type")
	}
	finalExec := finalTestModel.executor

	// Verify that nodes were initialized
	if len(finalExec.nodeStates) != 3 {
		t.Errorf("Expected 3 nodes in state, got %d", len(finalExec.nodeStates))
	}

	// Verify all nodes have correct paths (should use IDs, not indices)
	expectedPaths := []string{"node-1", "node-2", "node-3"}
	for i, expectedPath := range expectedPaths {
		if i >= len(finalExec.nodeStates) {
			t.Errorf("Missing node state for index %d", i)
			continue
		}
		if finalExec.nodeStates[i].path != expectedPath {
			t.Errorf("Node %d: expected path %s, got %s",
				i, expectedPath, finalExec.nodeStates[i].path)
		}
	}

	// Verify all nodes reached completed status
	for i, node := range finalExec.nodeStates {
		if node.status != NodeCompleted {
			t.Errorf("Node %d (%s): expected status Completed, got %v",
				i, node.name, node.status)
		}
	}

	// Verify node names
	expectedNames := []string{"First Node", "Second Node", "Third Node"}
	for i, expectedName := range expectedNames {
		if i >= len(finalExec.nodeStates) {
			continue
		}
		if finalExec.nodeStates[i].name != expectedName {
			t.Errorf("Node %d: expected name %s, got %s",
				i, expectedName, finalExec.nodeStates[i].name)
		}
	}

	// CRITICAL: Verify View() output contains node names (like Playwright would)
	viewOutput := finalTestModel.View()
	for _, expectedName := range expectedNames {
		if !strings.Contains(viewOutput, expectedName) {
			t.Errorf("View output missing node name: %s\nActual view:\n%s",
				expectedName, viewOutput)
		}
	}
}

// TestExecutorNodeDisplay_Hierarchical tests that hierarchical bentos display nodes
func TestExecutorNodeDisplay_Hierarchical(t *testing.T) {
	// Create a hierarchical bento definition (without explicit IDs)
	def := neta.Definition{
		Version: "1.0",
		Type:    "group.sequence",
		Name:    "Hierarchical Workflow",
		Nodes: []neta.Definition{
			{
				Version: "1.0",
				Type:    "http",
				Name:    "Node Zero",
				Parameters: map[string]interface{}{
					"method": "GET",
					"url":    "https://httpbin.org/json",
				},
			},
			{
				Version: "1.0",
				Type:    "http",
				Name:    "Node One",
				Parameters: map[string]interface{}{
					"method": "GET",
					"url":    "https://httpbin.org/json",
				},
			},
		},
	}

	// Create executor and start execution (simulates StartBento call)
	executor := NewExecutor()
	executor = executor.StartBento("test-workflow", "/test/path.yaml", "/test/workdir")
	tm := teatest.NewTestModel(t, &executorTestModel{executor: executor}, teatest.WithInitialTermSize(120, 40))

	// Send ExecutionInitMsg (simulates message from background goroutine)
	tm.Send(ExecutionInitMsg{Definition: def})
	time.Sleep(50 * time.Millisecond)

	// Send node messages using index-based paths
	tm.Send(NodeStartedMsg{Path: "0", Name: "Node Zero", NodeType: "http"})
	time.Sleep(50 * time.Millisecond)

	tm.Send(NodeCompletedMsg{Path: "0", Duration: 100 * time.Millisecond, Error: nil})
	time.Sleep(50 * time.Millisecond)

	tm.Send(NodeStartedMsg{Path: "1", Name: "Node One", NodeType: "http"})
	time.Sleep(50 * time.Millisecond)

	tm.Send(NodeCompletedMsg{Path: "1", Duration: 50 * time.Millisecond, Error: nil})
	time.Sleep(50 * time.Millisecond)

	// Send quit to finish
	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})

	// Wait for test to finish
	tm.WaitFinished(t, teatest.WithFinalTimeout(2*time.Second))

	// Get final model
	finalModel := tm.FinalModel(t)
	finalTestModel, ok := finalModel.(*executorTestModel)
	if !ok {
		t.Fatal("Final model is not *executorTestModel type")
	}
	finalExec := finalTestModel.executor

	// Verify nodes were initialized
	if len(finalExec.nodeStates) != 2 {
		t.Errorf("Expected 2 nodes in state, got %d", len(finalExec.nodeStates))
	}

	// Verify nodes use index-based paths when no ID present
	expectedPaths := []string{"0", "1"}
	for i, expectedPath := range expectedPaths {
		if i >= len(finalExec.nodeStates) {
			continue
		}
		if finalExec.nodeStates[i].path != expectedPath {
			t.Errorf("Node %d: expected path %s, got %s",
				i, expectedPath, finalExec.nodeStates[i].path)
		}
	}

	// Verify all nodes completed
	for i, node := range finalExec.nodeStates {
		if node.status != NodeCompleted {
			t.Errorf("Node %d: expected Completed, got %v", i, node.status)
		}
	}

	// CRITICAL: Verify View() output contains node names
	viewOutput := finalTestModel.View()
	expectedNames := []string{"Node Zero", "Node One"}
	for _, expectedName := range expectedNames {
		if !strings.Contains(viewOutput, expectedName) {
			t.Errorf("View output missing node name: %s\nActual view:\n%s",
				expectedName, viewOutput)
		}
	}
}

// TestExecutorSequenceRendering tests that the sequence component renders nodes
func TestExecutorSequenceRendering(t *testing.T) {
	def := neta.Definition{
		Version: "1.0",
		Type:    "group.sequence",
		Name:    "Test",
		Nodes: []neta.Definition{
			{
				ID:      "test-node",
				Version: "1.0",
				Type:    "http",
				Name:    "Test Node",
				Parameters: map[string]interface{}{
					"method": "GET",
					"url":    "https://httpbin.org/json",
				},
			},
		},
		Edges: []neta.NodeEdge{},
	}

	executor := NewExecutor()
	executor = executor.StartBento("test", "/test/path.yaml", "/test/workdir")
	tm := teatest.NewTestModel(t, &executorTestModel{executor: executor}, teatest.WithInitialTermSize(120, 40))

	// Initialize (simulates message from background goroutine)
	tm.Send(ExecutionInitMsg{Definition: def})
	time.Sleep(50 * time.Millisecond)

	// Mark as running (simulates chef execution)
	tm.Send(NodeStartedMsg{Path: "test-node", Name: "Test Node", NodeType: "http"})
	time.Sleep(50 * time.Millisecond)

	// Send quit
	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	tm.WaitFinished(t, teatest.WithFinalTimeout(2*time.Second))

	finalModel := tm.FinalModel(t)
	finalTestModel, ok := finalModel.(*executorTestModel)
	if !ok {
		t.Fatal("Final model is not *executorTestModel type")
	}
	finalExec := finalTestModel.executor

	// Verify sequence has steps
	steps := finalExec.convertNodesToSteps()
	if len(steps) != 1 {
		t.Errorf("Expected 1 step in sequence, got %d", len(steps))
	}

	if len(steps) > 0 {
		if steps[0].Name != "Test Node" {
			t.Errorf("Expected step name 'Test Node', got '%s'", steps[0].Name)
		}
		if steps[0].Status != components.StepRunning {
			t.Errorf("Expected step status Running, got %v", steps[0].Status)
		}
	}

	// CRITICAL: Verify View() output contains node name
	viewOutput := finalTestModel.View()
	if !strings.Contains(viewOutput, "Test Node") {
		t.Errorf("View output missing node name 'Test Node'\nActual view:\n%s", viewOutput)
	}
}
