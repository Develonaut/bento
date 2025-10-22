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
	"fmt"
	"time"

	"github.com/Develonaut/bento/pkg/neta"
	"github.com/Develonaut/bento/pkg/pantry"
	"github.com/Develonaut/bento/pkg/shoyu"
)

// ProgressMessenger receives execution progress events.
// Used by miso package for TUI/CLI progress display.
// Optional - nil check before use.
type ProgressMessenger interface {
	SendNodeStarted(path, name, nodeType string)
	SendNodeCompleted(path string, duration time.Duration, err error)
	SendLoopChild(loopPath, childName string, index, total int)
}

// Itamae orchestrates bento execution.
type Itamae struct {
	pantry      *pantry.Pantry
	logger      *shoyu.Logger     // Optional - can be nil
	messenger   ProgressMessenger // Optional - for TUI progress updates
	onProgress  ProgressCallback
	slowMoDelay time.Duration   // Delay between node completions for animations
	state       *executionState // Progress tracking state
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

// NewWithMessenger creates an Itamae with progress messaging.
// Messenger is used for TUI/CLI progress updates.
// Both logger and messenger are optional - can be nil.
// Automatically loads slowMo delay from config for TUI animations.
func NewWithMessenger(p *pantry.Pantry, logger *shoyu.Logger, messenger ProgressMessenger) *Itamae {
	// Note: Import miso here would create circular dependency, so we'll set slowMo from outside
	return &Itamae{
		pantry:      p,
		logger:      logger,
		messenger:   messenger,
		slowMoDelay: 0, // Will be set by caller to avoid circular dependency
	}
}

// SetSlowMoDelay sets the delay between node completions for animations.
func (i *Itamae) SetSlowMoDelay(delay time.Duration) {
	i.slowMoDelay = delay
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

	// Analyze graph structure and initialize execution state
	graph := analyzeGraph(def)
	i.state = newExecutionState(graph)

	if i.logger != nil {
		msg := msgBentoStarted(def.Name)
		i.logger.Info(msg.format())
	}

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

		if i.logger != nil {
			durationStr := formatDuration(result.Duration)
			msg := msgBentoFailed(durationStr)
			i.logger.Error(msg.format())
			i.logger.Error("Error: " + err.Error())
		}

		return result, err
	}

	result.Status = StatusSuccess

	if i.logger != nil {
		durationStr := formatDuration(result.Duration)
		msg := msgBentoCompleted(durationStr)
		i.logger.Info(msg.format())
	}

	return result, nil
}

// formatDuration formats a duration to match CLI output (e.g., "6ms", "1.2s")
func formatDuration(d time.Duration) string {
	if d < time.Second {
		// Less than 1 second - show milliseconds
		ms := d.Milliseconds()
		return fmt.Sprintf("%dms", ms)
	}
	// 1 second or more - show with decimal
	s := d.Seconds()
	return fmt.Sprintf("%.1fs", s)
}
