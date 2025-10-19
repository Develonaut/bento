// Package itamae provides the orchestration engine for executing bentos.
//
// "Itamae" (板前 - "sushi chef") is the skilled chef who coordinates every
// aspect of sushi preparation. Similarly, the itamae orchestrates bento
// execution, managing data flow, concurrency, and error handling.
//
// # Usage
//
//	p := pantry.New()
//	logger := shoyu.New(shoyu.Config{Level: shoyu.LevelInfo})
//	chef := itamae.New(p, logger)
//
//	// Execute a bento
//	result, err := chef.Serve(ctx, bentoDef)
//	if err != nil {
//	    log.Fatalf("Execution failed: %v", err)
//	}
//
// # Context Management
//
// The itamae passes data between neta through an execution context.
// Each neta's output becomes available to downstream neta via template variables.
//
// Learn more about context.Context: https://go.dev/blog/context
package itamae

import (
	"context"
	"time"

	"github.com/Develonaut/bento/pkg/neta"
	"github.com/Develonaut/bento/pkg/pantry"
	"github.com/Develonaut/bento/pkg/shoyu"
)

// Itamae orchestrates bento execution.
type Itamae struct {
	pantry     *pantry.Pantry
	logger     *shoyu.Logger
	onProgress ProgressCallback
}

// ProgressCallback is called when a node starts/completes execution.
type ProgressCallback func(nodeID string, status string)

// Result contains the result of a bento execution.
type Result struct {
	Status        Status                 // Execution status
	NodesExecuted int                    // Number of nodes executed
	NodeOutputs   map[string]interface{} // Output from each node
	Duration      time.Duration          // Total execution time
	Error         error                  // Error if execution failed
}

// Status represents the execution status.
type Status string

const (
	StatusSuccess   Status = "success"
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
)

// New creates a new Itamae orchestrator.
func New(p *pantry.Pantry, logger *shoyu.Logger) *Itamae {
	return &Itamae{
		pantry: p,
		logger: logger,
	}
}

// OnProgress registers a callback for progress updates.
func (i *Itamae) OnProgress(callback ProgressCallback) {
	i.onProgress = callback
}

// Serve executes a bento definition.
//
// Returns:
//   - *Result: Execution result with outputs from all nodes
//   - error: Any error that occurred during execution
func (i *Itamae) Serve(ctx context.Context, def *neta.Definition) (*Result, error) {
	start := time.Now()

	i.logger.Info("🍱 Starting bento execution",
		"bento_id", def.ID,
		"bento_name", def.Name)

	result := &Result{
		NodeOutputs: make(map[string]interface{}),
	}

	// Create execution context
	execCtx := newExecutionContext()

	// Execute the bento
	err := i.executeNode(ctx, def, execCtx, result)

	result.Duration = time.Since(start)

	if err != nil {
		result.Status = StatusFailed
		result.Error = err

		i.logger.Error("🍱 Bento execution failed",
			"bento_id", def.ID,
			"nodes_executed", result.NodesExecuted,
			"duration", result.Duration,
			"error", err)

		return result, err
	}

	result.Status = StatusSuccess

	i.logger.Info("🍱 Bento execution completed",
		"bento_id", def.ID,
		"nodes_executed", result.NodesExecuted,
		"duration", result.Duration)

	return result, nil
}
