package neta

import (
	"sync"
	"time"
)

// NodeExecutionState represents the execution state of a node
type NodeExecutionState string

const (
	NodeStatePending   NodeExecutionState = "pending"
	NodeStateExecuting NodeExecutionState = "executing"
	NodeStateCompleted NodeExecutionState = "completed"
	NodeStateError     NodeExecutionState = "error"
	NodeStateSkipped   NodeExecutionState = "skipped"
)

// ExecutionGraphNode represents a node in the execution graph with state
type ExecutionGraphNode struct {
	ID           string
	Name         string
	Type         string
	State        NodeExecutionState
	Progress     int       // 0-100
	Message      string    // Optional progress message
	StartTime    time.Time // When node started executing
	EndTime      time.Time // When node finished executing
	Error        string    // Error message if state is Error
	Weight       int       // Node weight for progress calculation
	Dependencies []string  // IDs of nodes this depends on
}

// ExecutionGraphEdge represents a connection between nodes
type ExecutionGraphEdge struct {
	From string
	To   string
}

// ExecutionGraphState holds the full state of graph execution
type ExecutionGraphState struct {
	Nodes          map[string]*ExecutionGraphNode
	Edges          []ExecutionGraphEdge
	ExecutionOrder [][]string // Topologically sorted layers
	CriticalPath   []string   // Longest path through the graph
	TotalWeight    int
	MaxParallelism int
	IsExecuting    bool
	StartTime      time.Time
	EndTime        time.Time
	CachedProgress float64 // Weighted progress (0-100)
}

// ExecutionGraphStore manages execution state with thread-safe operations
type ExecutionGraphStore struct {
	state     ExecutionGraphState
	mu        sync.RWMutex
	listeners []func(ExecutionGraphState)
}

// NewExecutionGraphStore creates a new execution graph store
func NewExecutionGraphStore() *ExecutionGraphStore {
	return &ExecutionGraphStore{
		state: ExecutionGraphState{
			Nodes:          make(map[string]*ExecutionGraphNode),
			Edges:          []ExecutionGraphEdge{},
			ExecutionOrder: [][]string{},
			CriticalPath:   []string{},
		},
		listeners: []func(ExecutionGraphState){},
	}
}

// InitializeGraph sets up the graph from analyzed execution graph
func (s *ExecutionGraphStore) InitializeGraph(
	nodes map[string]*ExecutionGraphNode,
	edges []ExecutionGraphEdge,
	executionOrder [][]string,
	criticalPath []string,
	totalWeight int,
	maxParallelism int,
) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Initialize all nodes with pending state
	for _, node := range nodes {
		node.State = NodeStatePending
		node.Progress = 0
	}

	s.state = ExecutionGraphState{
		Nodes:          nodes,
		Edges:          edges,
		ExecutionOrder: executionOrder,
		CriticalPath:   criticalPath,
		TotalWeight:    totalWeight,
		MaxParallelism: maxParallelism,
		IsExecuting:    true,
		StartTime:      time.Now(),
		CachedProgress: 0,
	}

	s.notifyListeners()
}

// SetNodeState updates the execution state of a node
func (s *ExecutionGraphStore) SetNodeState(nodeID string, state NodeExecutionState, errorMsg string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	node, exists := s.state.Nodes[nodeID]
	if !exists {
		return
	}

	node.State = state

	switch state {
	case NodeStateExecuting:
		node.StartTime = time.Now()
		node.Progress = 0

	case NodeStateCompleted, NodeStateSkipped:
		node.EndTime = time.Now()
		node.Progress = 100

	case NodeStateError:
		node.EndTime = time.Now()
		node.Error = errorMsg
		// Keep current progress, don't set to 100
	}

	s.state.CachedProgress = s.calculateProgress()
	s.notifyListeners()
}

// SetNodeProgress updates the progress of a node (0-100)
func (s *ExecutionGraphStore) SetNodeProgress(nodeID string, progress int, message string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	node, exists := s.state.Nodes[nodeID]
	if !exists {
		return
	}

	// Clamp progress to 0-100
	if progress < 0 {
		progress = 0
	}
	if progress > 100 {
		progress = 100
	}

	node.Progress = progress
	if message != "" {
		node.Message = message
	}

	s.state.CachedProgress = s.calculateProgress()
	s.notifyListeners()
}

// CompleteExecution marks the execution as complete
func (s *ExecutionGraphStore) CompleteExecution() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.state.IsExecuting = false
	s.state.EndTime = time.Now()
	s.notifyListeners()
}

// Reset clears the store state
func (s *ExecutionGraphStore) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.state = ExecutionGraphState{
		Nodes:          make(map[string]*ExecutionGraphNode),
		Edges:          []ExecutionGraphEdge{},
		ExecutionOrder: [][]string{},
		CriticalPath:   []string{},
	}
	s.notifyListeners()
}

// GetState returns a copy of the current state (thread-safe read)
func (s *ExecutionGraphStore) GetState() ExecutionGraphState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state
}

// GetNodeState returns the state of a specific node
func (s *ExecutionGraphStore) GetNodeState(nodeID string) (NodeExecutionState, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	node, exists := s.state.Nodes[nodeID]
	if !exists {
		return "", false
	}
	return node.State, true
}

// Subscribe adds a listener that will be called on state changes
func (s *ExecutionGraphStore) Subscribe(listener func(ExecutionGraphState)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.listeners = append(s.listeners, listener)
}

// notifyListeners calls all subscribed listeners with current state
func (s *ExecutionGraphStore) notifyListeners() {
	// Call listeners without holding the lock to avoid deadlocks
	state := s.state
	for _, listener := range s.listeners {
		listener(state)
	}
}

// calculateProgress computes weighted progress across all nodes
func (s *ExecutionGraphStore) calculateProgress() float64 {
	if s.state.TotalWeight == 0 {
		return 0
	}

	var completedWeight float64
	for _, node := range s.state.Nodes {
		nodeWeight := float64(node.Weight)
		if nodeWeight == 0 {
			nodeWeight = 1 // Default weight
		}

		switch node.State {
		case NodeStateCompleted, NodeStateSkipped:
			completedWeight += nodeWeight
		case NodeStateExecuting:
			// Partial credit based on progress
			completedWeight += nodeWeight * (float64(node.Progress) / 100.0)
		case NodeStateError:
			// Count errors as completed for progress calculation
			completedWeight += nodeWeight
		}
	}

	return (completedWeight / float64(s.state.TotalWeight)) * 100.0
}

// GetNodesByState returns all nodes in a given state
func (s *ExecutionGraphStore) GetNodesByState(state NodeExecutionState) []*ExecutionGraphNode {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var nodes []*ExecutionGraphNode
	for _, node := range s.state.Nodes {
		if node.State == state {
			nodes = append(nodes, node)
		}
	}
	return nodes
}
