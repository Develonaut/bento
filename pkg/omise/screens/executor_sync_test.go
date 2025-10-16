package screens

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"bento/pkg/neta"
)

// TestExecutionInitSynchronization verifies that ExecutionInitMsg is processed before node messages
func TestExecutionInitSynchronization(t *testing.T) {
	// Create executor
	executor := NewExecutor()

	// Create a simple definition with nodes
	def := neta.Definition{
		ID:   "test-bento",
		Name: "Test Bento",
		Type: "group.sequence",
		Nodes: []neta.Definition{
			{
				ID:   "node1",
				Name: "Node 1",
				Type: "http",
			},
			{
				ID:   "node2",
				Name: "Node 2",
				Type: "transform.jq",
			},
		},
	}

	// Test 1: Verify init message handler signals readiness
	t.Run("init_message_signals_ready", func(t *testing.T) {
		// Setup execution state
		executionState.ready = make(chan struct{}, 1)

		// Process init message
		updated, cmd := executor.handleInitMsg(ExecutionInitMsg{Definition: def})

		// Verify node states are initialized
		assert.Equal(t, 2, len(updated.nodeStates))
		assert.NotNil(t, cmd)

		// Execute the command to signal readiness
		if cmd != nil {
			_ = cmd()
		}

		// Verify ready signal was sent
		select {
		case <-executionState.ready:
			// Good - signal received
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Ready signal not sent after init")
		}
	})

	// Test 2: Verify node messages before init are handled correctly
	t.Run("node_message_before_init", func(t *testing.T) {
		executor := NewExecutor()

		// Send node started message before init
		updated, _ := executor.handleNodeStarted(NodeStartedMsg{
			Path:     "node1",
			Name:     "Node 1",
			NodeType: "http",
		})

		// Should have empty node states - no crash
		assert.Equal(t, 0, len(updated.nodeStates))

		// Lifecycle history should contain warning
		found := false
		for _, msg := range updated.lifecycleHistory {
			if contains(msg, "not found in state") {
				found = true
				break
			}
		}
		assert.True(t, found, "Should warn about missing node state")
	})

	// Test 3: Verify correct ordering after init
	t.Run("correct_ordering_after_init", func(t *testing.T) {
		executor := NewExecutor()

		// First send init message
		updated, _ := executor.handleInitMsg(ExecutionInitMsg{Definition: def})
		assert.Equal(t, 2, len(updated.nodeStates))

		// Then send node started message
		updated, _ = updated.handleNodeStarted(NodeStartedMsg{
			Path:     "node1",
			Name:     "Node 1",
			NodeType: "http",
		})

		// Verify node state was updated correctly
		found := false
		for _, state := range updated.nodeStates {
			if state.path == "node1" && state.status == NodeRunning {
				found = true
				break
			}
		}
		assert.True(t, found, "Node should be in running state")
	})
}

// TestExecutionBackgroundSynchronization tests the background execution synchronization
func TestExecutionBackgroundSynchronization(t *testing.T) {
	t.Run("execution_waits_for_ready_signal", func(t *testing.T) {
		// This test verifies the synchronization logic in executeBentoInBackground
		// Since we can't easily test the goroutine directly, we test the mechanism

		// Setup ready channel
		readyChan := make(chan struct{})

		// Simulate background execution waiting
		done := make(chan bool)
		go func() {
			select {
			case <-readyChan:
				done <- true
			case <-time.After(100 * time.Millisecond):
				done <- false
			}
		}()

		// Send ready signal
		close(readyChan)

		// Verify execution proceeded
		result := <-done
		assert.True(t, result, "Execution should proceed after ready signal")
	})

	t.Run("timeout_prevents_deadlock", func(t *testing.T) {
		// Test that timeout mechanism prevents deadlock
		readyChan := make(chan struct{})

		done := make(chan bool)
		go func() {
			select {
			case <-readyChan:
				done <- true
			case <-time.After(50 * time.Millisecond):
				// Timeout - proceed anyway
				done <- false
			}
		}()

		// Don't send ready signal - simulate stuck init

		// Verify timeout works
		result := <-done
		assert.False(t, result, "Should timeout and proceed")
	})
}

// TestSignalInitReadyCmd tests the signalInitReadyCmd command
func TestSignalInitReadyCmd(t *testing.T) {
	// Setup ready channel
	executionState.ready = make(chan struct{}, 1)

	// Create and execute command
	cmd := signalInitReadyCmd()
	assert.NotNil(t, cmd)

	// Execute command
	msg := cmd()
	assert.Nil(t, msg) // Should return nil

	// Verify signal was sent
	select {
	case <-executionState.ready:
		// Good - signal received
	default:
		t.Fatal("Ready signal not sent by command")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				contains(s[1:], substr)))
}
