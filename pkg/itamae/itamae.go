// Package itamae provides the orchestration engine for executing neta definitions.
// Itamae (板前) means "sushi chef" - the one who prepares each piece.
package itamae

import (
	"context"
	"fmt"
	"time"

	"bento/pkg/neta"
)

// ProgressMessenger receives execution progress events.
// Used by TUI to display real-time progress.
// Optional - nil check before use for non-TUI execution.
type ProgressMessenger interface {
	// SendNodeStarted notifies that a node has started execution.
	// path: node path in tree (e.g., "0", "1.2", "0.1.3")
	// name: human-readable node name
	// nodeType: node type (e.g., "http", "transform.jq")
	SendNodeStarted(path, name, nodeType string)

	// SendNodeCompleted notifies that a node has finished execution.
	// path: node path in tree
	// duration: how long the node took to execute
	// err: error if node failed, nil if successful
	SendNodeCompleted(path string, duration time.Duration, err error)
}

// Itamae orchestrates the execution of neta definitions.
type Itamae struct {
	pantry        Registry
	messenger     ProgressMessenger         // Optional - can be nil
	slowMoDelayMs int                       // Delay in milliseconds for slow-mo mode (0 = off)
	store         *neta.ExecutionGraphStore // Optional - for graph-based execution tracking
}

// Registry provides node type lookup.
type Registry interface {
	Get(nodeType string) (neta.Executable, error)
}

// New creates a new Itamae with the provided registry.
func New(registry Registry) *Itamae {
	return &Itamae{
		pantry:        registry,
		messenger:     nil,
		slowMoDelayMs: 0,
	}
}

// NewWithMessenger creates an Itamae with progress messaging.
func NewWithMessenger(registry Registry, messenger ProgressMessenger) *Itamae {
	return &Itamae{
		pantry:        registry,
		messenger:     messenger,
		slowMoDelayMs: 0,
	}
}

// SetSlowMoDelay sets the slow-mo delay in milliseconds (0 = off)
func (i *Itamae) SetSlowMoDelay(delayMs int) {
	i.slowMoDelayMs = delayMs
}

// SetStore sets the execution graph store for tracking node states
func (i *Itamae) SetStore(store *neta.ExecutionGraphStore) {
	i.store = store
}

// Execute runs a neta definition and returns the result.
func (i *Itamae) Execute(ctx context.Context, def neta.Definition) (neta.Result, error) {
	// Use graph-based execution if store is available and node has edges
	if i.store != nil && def.IsGroup() && len(def.Edges) > 0 {
		return i.executeGraph(ctx, def)
	}

	// Fall back to hierarchical execution
	if def.IsGroup() {
		return i.executeGroup(ctx, def, "")
	}
	return i.executeSingle(ctx, def, "")
}

// executeSingle runs a single node.
func (i *Itamae) executeSingle(ctx context.Context, def neta.Definition, path string) (neta.Result, error) {
	i.notifyNodeStarted(path, def.Name, def.Type)

	exec, err := i.pantry.Get(def.Type)
	if err != nil {
		i.notifyNodeCompleted(path, 0, err)
		return neta.Result{}, fmt.Errorf("node type not found: %s: %w", def.Type, err)
	}

	return i.executeWithTiming(ctx, exec, def.Parameters, path)
}

// notifyNodeStarted sends node start notification if messenger available
func (i *Itamae) notifyNodeStarted(path, name, nodeType string) {
	if i.messenger != nil {
		i.messenger.SendNodeStarted(path, name, nodeType)
	}

	// Apply slow-mo delay if configured
	if i.slowMoDelayMs > 0 {
		time.Sleep(time.Duration(i.slowMoDelayMs) * time.Millisecond)
	}
}

// notifyNodeCompleted sends node completion notification if messenger available
func (i *Itamae) notifyNodeCompleted(path string, duration time.Duration, err error) {
	if i.messenger != nil {
		i.messenger.SendNodeCompleted(path, duration, err)
	}
}

// executeWithTiming executes a node and tracks timing
func (i *Itamae) executeWithTiming(ctx context.Context, exec neta.Executable, params map[string]interface{}, path string) (neta.Result, error) {
	start := time.Now()
	result, err := exec.Execute(ctx, params)
	duration := time.Since(start)

	i.notifyNodeCompleted(path, duration, err)
	return result, err
}

// executeGroup runs a group of nodes in sequence.
func (i *Itamae) executeGroup(ctx context.Context, def neta.Definition, basePath string) (neta.Result, error) {
	results := make([]neta.Result, 0, len(def.Nodes))

	for idx, child := range def.Nodes {
		nodePath := buildPath(basePath, idx)
		result, err := i.executeChildNode(ctx, child, nodePath)
		if err != nil {
			return neta.Result{}, err
		}
		results = append(results, result)
	}

	return neta.Result{Output: results}, nil
}

// executeChildNode executes a single child node with progress tracking
func (i *Itamae) executeChildNode(ctx context.Context, child neta.Definition, nodePath string) (neta.Result, error) {
	i.notifyNodeStarted(nodePath, child.Name, child.Type)

	start := time.Now()
	result, err := i.executeNode(ctx, child, nodePath)
	duration := time.Since(start)

	i.notifyNodeCompleted(nodePath, duration, err)
	return result, err
}

// executeNode routes to group or single execution
func (i *Itamae) executeNode(ctx context.Context, def neta.Definition, path string) (neta.Result, error) {
	if def.IsGroup() {
		return i.executeGroup(ctx, def, path)
	}
	return i.executeSingle(ctx, def, path)
}

// buildPath constructs node path for tracking
func buildPath(basePath string, index int) string {
	if basePath == "" {
		return fmt.Sprintf("%d", index)
	}
	return fmt.Sprintf("%s.%d", basePath, index)
}
