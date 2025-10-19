# Phase 6: Itamae Package (板前 - "Sushi Chef")

**Duration:** 2-3 weeks
**Package:** `pkg/itamae/`
**Dependencies:** `pkg/neta`, `pkg/pantry`, `pkg/shoyu`

---

## TDD Philosophy

> **Write tests FIRST to define contracts**

Orchestration tests should verify:
1. Execute simple linear bentos (A → B → C)
2. Execute group bentos with nested nodes
3. Pass context between neta (output from A becomes input to B)
4. Handle errors gracefully (fail one neta, propagate error)
5. Execute loops correctly (forEach through items)
6. Execute parallel neta concurrently
7. Track progress for long-running processes
8. Handle cancellation (context.Context)

---

## Phase Overview

The itamae ("sushi chef") is the heart of the bento system. Like a skilled sushi chef who orchestrates every aspect of meal preparation, the itamae coordinates the execution of all neta in a bento.

Key responsibilities:
- **Execution:** Run bentos from start to finish
- **Context management:** Pass data between neta
- **Error handling:** Graceful failures with clear error chains
- **Progress tracking:** Report execution status
- **Concurrency:** Handle parallel and loop neta correctly
- **Cancellation:** Respect context.Context for clean shutdowns

### Why "Itamae"?

An itamae (板前) is a highly skilled sushi chef who has mastered the art of preparation. The itamae doesn't just make sushi—they coordinate timing, temperature, presentation, and flow. Similarly, our orchestration engine doesn't just run neta—it coordinates execution, data flow, error handling, and progress.

---

## Success Criteria

**Phase 6 Complete When:**
- [ ] Execute simple linear bentos (sequential neta)
- [ ] Execute group bentos with nested nodes
- [ ] Context passing between neta (data flows correctly)
- [ ] Error handling with clear error chains
- [ ] Loop neta execution (forEach, times, while)
- [ ] Parallel neta execution (goroutines, wait groups)
- [ ] Progress tracking with callbacks
- [ ] Context cancellation support
- [ ] Integration tests for all scenarios
- [ ] Files < 250 lines each
- [ ] File-level documentation complete
- [ ] `/code-review` run with Karen + Colossus approval

---

## Test-First Approach

### Step 1: Test simple linear execution

Create `pkg/itamae/itamae_test.go`:

```go
package itamae_test

import (
	"context"
	"testing"

	"github.com/Develonaut/bento/pkg/itamae"
	"github.com/Develonaut/bento/pkg/neta"
	"github.com/Develonaut/bento/pkg/pantry"
	"github.com/Develonaut/bento/pkg/shoyu"
)

// Test: Execute a simple linear bento (A → B → C)
func TestItamae_LinearExecution(t *testing.T) {
	ctx := context.Background()

	// Create bento: node-1 → node-2 → node-3
	bento := &neta.Definition{
		ID:      "linear-bento",
		Type:    "group",
		Version: "1.0.0",
		Name:    "Linear Bento",
		Nodes: []neta.Definition{
			{
				ID:   "node-1",
				Type: "edit-fields",
				Parameters: map[string]interface{}{
					"values": map[string]interface{}{"step": 1},
				},
			},
			{
				ID:   "node-2",
				Type: "edit-fields",
				Parameters: map[string]interface{}{
					"values": map[string]interface{}{"step": 2},
				},
			},
			{
				ID:   "node-3",
				Type: "edit-fields",
				Parameters: map[string]interface{}{
					"values": map[string]interface{}{"step": 3},
				},
			},
		},
		Edges: []neta.Edge{
			{ID: "edge-1", Source: "node-1", Target: "node-2"},
			{ID: "edge-2", Source: "node-2", Target: "node-3"},
		},
	}

	// Create itamae
	p := pantry.New()
	logger := shoyu.New(shoyu.Config{Level: shoyu.LevelInfo})
	chef := itamae.New(p, logger)

	// Execute
	result, err := chef.Serve(ctx, bento)
	if err != nil {
		t.Fatalf("Serve failed: %v", err)
	}

	// Verify all nodes executed
	if result.NodesExecuted != 3 {
		t.Errorf("NodesExecuted = %d, want 3", result.NodesExecuted)
	}

	if result.Status != itamae.StatusSuccess {
		t.Errorf("Status = %v, want Success", result.Status)
	}
}
```

### Step 2: Test context passing between neta

```go
// Test: Data should flow from one neta to the next
func TestItamae_ContextPassing(t *testing.T) {
	ctx := context.Background()

	// Bento: set name → use name in template
	bento := &neta.Definition{
		ID:   "context-bento",
		Type: "group",
		Nodes: []neta.Definition{
			{
				ID:   "set-name",
				Type: "edit-fields",
				Parameters: map[string]interface{}{
					"values": map[string]interface{}{
						"productName": "Widget",
					},
				},
			},
			{
				ID:   "use-name",
				Type: "edit-fields",
				Parameters: map[string]interface{}{
					"values": map[string]interface{}{
						"title": "{{.set-name.productName}}",
					},
				},
			},
		},
		Edges: []neta.Edge{
			{ID: "edge-1", Source: "set-name", Target: "use-name"},
		},
	}

	p := pantry.New()
	logger := shoyu.New(shoyu.Config{Level: shoyu.LevelInfo})
	chef := itamae.New(p, logger)

	result, err := chef.Serve(ctx, bento)
	if err != nil {
		t.Fatalf("Serve failed: %v", err)
	}

	// Verify template was resolved
	output := result.NodeOutputs["use-name"].(map[string]interface{})
	if output["title"] != "Widget" {
		t.Errorf("title = %v, want Widget (template should be resolved)", output["title"])
	}
}
```

### Step 3: Test error handling

```go
// Test: Error in one neta should stop execution and report clearly
func TestItamae_ErrorHandling(t *testing.T) {
	ctx := context.Background()

	// Bento with an invalid HTTP request (bad URL)
	bento := &neta.Definition{
		ID:   "error-bento",
		Type: "group",
		Nodes: []neta.Definition{
			{
				ID:   "node-1",
				Type: "edit-fields",
				Parameters: map[string]interface{}{
					"values": map[string]interface{}{"step": 1},
				},
			},
			{
				ID:   "bad-request",
				Type: "http-request",
				Parameters: map[string]interface{}{
					"url":    "htp://invalid-url", // Invalid URL
					"method": "GET",
				},
			},
			{
				ID:   "node-3",
				Type: "edit-fields",
				Parameters: map[string]interface{}{
					"values": map[string]interface{}{"step": 3},
				},
			},
		},
		Edges: []neta.Edge{
			{ID: "edge-1", Source: "node-1", Target: "bad-request"},
			{ID: "edge-2", Source: "bad-request", Target: "node-3"},
		},
	}

	p := pantry.New()
	logger := shoyu.New(shoyu.Config{Level: shoyu.LevelInfo})
	chef := itamae.New(p, logger)

	result, err := chef.Serve(ctx, bento)

	// Should return error
	if err == nil {
		t.Fatal("Expected error from invalid HTTP request")
	}

	// Error should mention which node failed
	if !strings.Contains(err.Error(), "bad-request") {
		t.Errorf("Error should mention failing node 'bad-request': %v", err)
	}

	// node-3 should NOT have executed (execution stopped at error)
	if result != nil && result.NodesExecuted > 2 {
		t.Error("Execution should stop after error")
	}
}
```

### Step 4: Test loop execution

```go
// Test: Loop neta should execute body for each item
func TestItamae_LoopExecution(t *testing.T) {
	// CRITICAL FOR PHASE 8: CSV iteration
	ctx := context.Background()

	// Simulate CSV rows
	csvRows := []map[string]interface{}{
		{"sku": "PROD-001", "name": "Product A"},
		{"sku": "PROD-002", "name": "Product B"},
		{"sku": "PROD-003", "name": "Product C"},
	}

	bento := &neta.Definition{
		ID:   "loop-bento",
		Type: "group",
		Nodes: []neta.Definition{
			{
				ID:   "provide-items",
				Type: "edit-fields",
				Parameters: map[string]interface{}{
					"values": map[string]interface{}{
						"rows": csvRows,
					},
				},
			},
			{
				ID:   "loop-rows",
				Type: "loop",
				Parameters: map[string]interface{}{
					"mode":  "forEach",
					"items": "{{.provide-items.rows}}",
					"body": map[string]interface{}{
						// Loop body (simplified for test)
						"type": "edit-fields",
						"parameters": map[string]interface{}{
							"values": map[string]interface{}{
								"folder": "products/{{.item.sku}}",
							},
						},
					},
				},
			},
		},
		Edges: []neta.Edge{
			{ID: "edge-1", Source: "provide-items", Target: "loop-rows"},
		},
	}

	p := pantry.New()
	logger := shoyu.New(shoyu.Config{Level: shoyu.LevelInfo})
	chef := itamae.New(p, logger)

	result, err := chef.Serve(ctx, bento)
	if err != nil {
		t.Fatalf("Serve failed: %v", err)
	}

	// Loop should have executed 3 times
	loopResult := result.NodeOutputs["loop-rows"].(map[string]interface{})
	if loopResult["iterations"] != 3 {
		t.Errorf("iterations = %v, want 3", loopResult["iterations"])
	}
}
```

### Step 5: Test parallel execution

```go
// Test: Parallel neta should execute children concurrently
func TestItamae_ParallelExecution(t *testing.T) {
	ctx := context.Background()

	bento := &neta.Definition{
		ID:   "parallel-bento",
		Type: "group",
		Nodes: []neta.Definition{
			{
				ID:   "parallel-group",
				Type: "parallel",
				Parameters: map[string]interface{}{
					"maxConcurrency": 4,
				},
				Nodes: []neta.Definition{
					{ID: "task-1", Type: "edit-fields", Parameters: map[string]interface{}{"values": map[string]interface{}{"task": 1}}},
					{ID: "task-2", Type: "edit-fields", Parameters: map[string]interface{}{"values": map[string]interface{}{"task": 2}}},
					{ID: "task-3", Type: "edit-fields", Parameters: map[string]interface{}{"values": map[string]interface{}{"task": 3}}},
					{ID: "task-4", Type: "edit-fields", Parameters: map[string]interface{}{"values": map[string]interface{}{"task": 4}}},
				},
			},
		},
	}

	p := pantry.New()
	logger := shoyu.New(shoyu.Config{Level: shoyu.LevelInfo})
	chef := itamae.New(p, logger)

	start := time.Now()
	result, err := chef.Serve(ctx, bento)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Serve failed: %v", err)
	}

	// All 4 tasks should have executed
	if result.NodesExecuted != 5 { // parallel-group + 4 children
		t.Errorf("NodesExecuted = %d, want 5", result.NodesExecuted)
	}

	// Should complete faster than sequential (hard to test reliably, but good to log)
	t.Logf("Parallel execution took %v", duration)
}
```

### Step 6: Test progress tracking

```go
// Test: Progress callback should be called for each node
func TestItamae_ProgressTracking(t *testing.T) {
	ctx := context.Background()

	bento := &neta.Definition{
		ID:   "progress-bento",
		Type: "group",
		Nodes: []neta.Definition{
			{ID: "node-1", Type: "edit-fields"},
			{ID: "node-2", Type: "edit-fields"},
			{ID: "node-3", Type: "edit-fields"},
		},
	}

	p := pantry.New()
	logger := shoyu.New(shoyu.Config{Level: shoyu.LevelInfo})
	chef := itamae.New(p, logger)

	progressCalls := 0
	onProgress := func(nodeID string, status string) {
		progressCalls++
		t.Logf("Progress: %s - %s", nodeID, status)
	}

	chef.OnProgress(onProgress)

	_, err := chef.Serve(ctx, bento)
	if err != nil {
		t.Fatalf("Serve failed: %v", err)
	}

	// Should have called progress for each node (at least starting/completed)
	if progressCalls < 3 {
		t.Errorf("progressCalls = %d, want at least 3", progressCalls)
	}
}
```

### Step 7: Test context cancellation

```go
// Test: Cancelling context should stop execution cleanly
func TestItamae_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// Bento with long-running task
	bento := &neta.Definition{
		ID:   "cancellation-bento",
		Type: "group",
		Nodes: []neta.Definition{
			{
				ID:   "long-task",
				Type: "shell-command",
				Parameters: map[string]interface{}{
					"command": "sleep",
					"args":    []string{"10"},
				},
			},
		},
	}

	p := pantry.New()
	logger := shoyu.New(shoyu.Config{Level: shoyu.LevelInfo})
	chef := itamae.New(p, logger)

	// Cancel after 1 second
	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()

	_, err := chef.Serve(ctx, bento)

	// Should return context cancelled error
	if err == nil {
		t.Fatal("Expected error from context cancellation")
	}

	if !strings.Contains(err.Error(), "context canceled") {
		t.Errorf("Error should mention context cancellation: %v", err)
	}
}
```

---

## File Structure

```
pkg/itamae/
├── itamae.go              # Main orchestrator (~200 lines)
├── executor.go            # Node execution logic (~200 lines)
├── context.go             # Context management (~150 lines)
├── graph.go               # Graph traversal (~150 lines)
├── loop.go                # Loop neta execution (~200 lines)
├── parallel.go            # Parallel neta execution (~200 lines)
├── progress.go            # Progress tracking (~100 lines)
├── errors.go              # Error handling (~100 lines)
└── itamae_test.go         # Integration tests (~500 lines)
```

---

## Implementation Guidance

**File: `pkg/itamae/itamae.go`**

```go
// Package itamae provides the orchestration engine for executing bentos.
//
// "Itamae" (板前 - "sushi chef") is the skilled chef who coordinates every
// aspect of sushi preparation. Similarly, the itamae orchestrates bento
// execution, managing data flow, concurrency, and error handling.
//
// Usage:
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
// Context Management:
// The itamae passes data between neta through an execution context.
// Each neta's output becomes available to downstream neta via template variables.
//
// Learn more about context.Context:
// https://go.dev/blog/context
package itamae

import (
	"context"
	"fmt"

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
	Status        Status                            // Execution status
	NodesExecuted int                               // Number of nodes executed
	NodeOutputs   map[string]interface{}            // Output from each node
	Duration      time.Duration                     // Total execution time
	Error         error                             // Error if execution failed
}

// Status represents the execution status.
type Status string

const (
	StatusSuccess Status = "success"
	StatusFailed  Status = "failed"
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

	i.logger.Info().
		Str("bento_id", def.ID).
		Str("bento_name", def.Name).
		Msg("🍱 Starting bento execution")

	result := &Result{
		NodeOutputs: make(map[string]interface{}),
	}

	// Execute the bento
	err := i.execute(ctx, def, result)

	result.Duration = time.Since(start)

	if err != nil {
		result.Status = StatusFailed
		result.Error = err

		i.logger.Error().
			Err(err).
			Str("bento_id", def.ID).
			Int("nodes_executed", result.NodesExecuted).
			Dur("duration", result.Duration).
			Msg("🍱 Bento execution failed")

		return result, err
	}

	result.Status = StatusSuccess

	i.logger.Info().
		Str("bento_id", def.ID).
		Int("nodes_executed", result.NodesExecuted).
		Dur("duration", result.Duration).
		Msg("🍱 Bento execution completed")

	return result, nil
}
```

**File: `pkg/itamae/executor.go`**

```go
package itamae

import (
	"context"
	"fmt"

	"github.com/Develonaut/bento/pkg/neta"
)

// execute runs a neta definition (handles all neta types).
func (i *Itamae) execute(ctx context.Context, def *neta.Definition, result *Result) error {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Special handling for group/loop/parallel neta
	switch def.Type {
	case "group":
		return i.executeGroup(ctx, def, result)
	case "loop":
		return i.executeLoop(ctx, def, result)
	case "parallel":
		return i.executeParallel(ctx, def, result)
	}

	// Regular neta execution
	return i.executeSingleNeta(ctx, def, result)
}

// executeSingleNeta executes a single (non-group) neta.
func (i *Itamae) executeSingleNeta(ctx context.Context, def *neta.Definition, result *Result) error {
	// Progress callback: starting
	if i.onProgress != nil {
		i.onProgress(def.ID, "starting")
	}

	i.logger.Info().
		Str("neta_id", def.ID).
		Str("neta_type", def.Type).
		Msg("🍙 Executing neta")

	// Get neta implementation from pantry
	netaImpl, err := i.pantry.GetNew(def.Type)
	if err != nil {
		return fmt.Errorf("failed to get neta '%s' (type: %s): %w", def.ID, def.Type, err)
	}

	// Resolve template variables in parameters
	resolvedParams := i.resolveTemplates(def.Parameters, result.NodeOutputs)

	// Execute neta
	output, err := netaImpl.Execute(ctx, resolvedParams)
	if err != nil {
		return fmt.Errorf("neta '%s' (type: %s) failed: %w", def.ID, def.Type, err)
	}

	// Store output
	result.NodeOutputs[def.ID] = output
	result.NodesExecuted++

	// Progress callback: completed
	if i.onProgress != nil {
		i.onProgress(def.ID, "completed")
	}

	i.logger.Info().
		Str("neta_id", def.ID).
		Str("neta_type", def.Type).
		Msg("✓ Neta completed")

	return nil
}
```

---

## Common Go Pitfalls to Avoid

1. **Context propagation**: Always pass ctx through the call chain
   ```go
   // ❌ BAD - creates new context, loses cancellation
   func (i *Itamae) execute(ctx context.Context, def *neta.Definition) error {
       netaImpl.Execute(context.Background(), params)
   }

   // ✅ GOOD - passes context through
   func (i *Itamae) execute(ctx context.Context, def *neta.Definition) error {
       netaImpl.Execute(ctx, params)
   }
   ```

2. **Goroutine leaks**: Always wait for goroutines to finish
   ```go
   // ❌ BAD - goroutines may still be running after return
   for _, child := range children {
       go execute(child)
   }
   return nil

   // ✅ GOOD - wait for all goroutines
   var wg sync.WaitGroup
   for _, child := range children {
       wg.Add(1)
       go func(c neta.Definition) {
           defer wg.Done()
           execute(c)
       }(child)
   }
   wg.Wait()
   ```

3. **Error wrapping**: Use fmt.Errorf with %w to preserve error chain
   ```go
   // ❌ BAD - loses original error
   return fmt.Errorf("neta failed")

   // ✅ GOOD - wraps error with context
   return fmt.Errorf("neta '%s' failed: %w", nodeID, err)
   ```

---

## Critical for Phase 8

**Loop Execution:**
- Must iterate through CSV rows correctly
- Each iteration gets fresh context with `{{.item}}` variable
- Must handle 50+ iterations without memory leaks

**Long-Running Processes:**
- Shell-command neta (Blender) can take 2-5 minutes
- Must stream progress updates
- Context cancellation should kill child processes

**Error Recovery:**
- If Figma API fails on row 23, should we stop or continue?
- Decision: Stop by default, add "continueOnError" flag for future

**Progress Tracking:**
- Must show "Rendering product 3 of 50... 45% complete"
- Progress callback called for each neta start/complete
- Long-running neta should emit intermediate progress

---

## Bento Box Principle Checklist

- [ ] Files < 250 lines each (split into executor.go, context.go, loop.go, etc.)
- [ ] Functions < 20 lines
- [ ] Single responsibility per file
- [ ] Clear error messages (mention node ID, type, what failed)
- [ ] File-level documentation

---

## Phase Completion

**Phase 6 MUST end with:**

1. All tests passing (`go test ./pkg/itamae/...`)
2. Run `/code-review` slash command
3. Address feedback from Karen and Colossus
4. Get explicit approval from both agents
5. Document any decisions in `.claude/strategy/`

**Do not proceed to Phase 7 until code review is approved.**

---

## Claude Prompt Template

```
I need to implement Phase 6: itamae (orchestration engine) following TDD principles.

This is the most complex phase - the "sushi chef" that coordinates everything.

Please read:
- .claude/strategy/phase-6-itamae.md (this file)
- .claude/BENTO_BOX_PRINCIPLE.md
- .claude/COMPLETE_NODE_INVENTORY.md (for loop/parallel/group behavior)

Then:

1. Create `pkg/itamae/itamae_test.go` with integration tests for:
   - Linear execution (A → B → C)
   - Context passing (A's output becomes B's input via templates)
   - Error handling (clear error chains with node IDs)
   - Loop execution (forEach through items) - CRITICAL FOR PHASE 8
   - Parallel execution (goroutines)
   - Progress tracking (callbacks)
   - Context cancellation (clean shutdown)

2. Watch the tests fail

3. Implement to make tests pass (split into multiple files):
   - pkg/itamae/itamae.go (~200 lines) - Main orchestrator
   - pkg/itamae/executor.go (~200 lines) - Node execution
   - pkg/itamae/context.go (~150 lines) - Context management
   - pkg/itamae/graph.go (~150 lines) - Graph traversal
   - pkg/itamae/loop.go (~200 lines) - Loop neta
   - pkg/itamae/parallel.go (~200 lines) - Parallel neta
   - pkg/itamae/progress.go (~100 lines) - Progress tracking
   - pkg/itamae/errors.go (~100 lines) - Error handling

4. Add file-level documentation for each file

Remember:
- Write tests FIRST
- Files < 250 lines (split into focused files)
- Functions < 20 lines
- Clear error chains (fmt.Errorf with %w)
- Context propagation (always pass ctx through)
- No goroutine leaks (always wait with sync.WaitGroup)

CRITICAL FOR PHASE 8:
- Loop must iterate CSV rows correctly
- Progress tracking must show "Rendering 3/50"
- Context cancellation must kill shell commands

When complete, run `/code-review` and get Karen + Colossus approval.
```

---

## Dependencies

No additional dependencies needed - uses stdlib only:
- `context` for cancellation
- `sync` for WaitGroup (parallel execution)
- `time` for duration tracking

---

## Notes

- This is the most complex package - take time to get it right
- Split into multiple files to maintain Bento Box Principle
- Test-driven approach is critical here - tests define the behavior
- Error messages must include node ID for debugging
- Context cancellation must be respected everywhere
- Loop and parallel execution are critical for Phase 8

---

**Status:** Ready for implementation
**Next Phase:** Phase 7 (cmd/bento CLI) - depends on completion of Phase 6
