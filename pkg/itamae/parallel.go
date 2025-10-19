package itamae

import (
	"context"
	"sync"

	"github.com/Develonaut/bento/pkg/neta"
)

// executeParallel executes a parallel neta.
// Runs all child nodes concurrently using goroutines.
func (i *Itamae) executeParallel(
	ctx context.Context,
	def *neta.Definition,
	execCtx *executionContext,
	result *Result,
) error {
	i.notifyProgress(def.ID, "starting")

	childCount := len(def.Nodes)

	i.logger.Info("⚡ Starting parallel execution",
		"parallel_id", def.ID,
		"child_count", childCount)

	// Handle empty parallel
	if childCount == 0 {
		i.notifyProgress(def.ID, "completed")
		result.NodesExecuted++
		return nil
	}

	// Get max concurrency (default: no limit)
	maxConcurrency := 0
	if mc, ok := def.Parameters["maxConcurrency"]; ok {
		if mcFloat, ok := mc.(float64); ok {
			maxConcurrency = int(mcFloat)
		}
	}

	// Execute children in parallel
	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error

	// Create semaphore for concurrency control
	var sem chan struct{}
	if maxConcurrency > 0 {
		sem = make(chan struct{}, maxConcurrency)
	}

	for idx := range def.Nodes {
		child := &def.Nodes[idx]

		wg.Add(1)
		go func(node *neta.Definition) {
			defer wg.Done()

			// Acquire semaphore if concurrency limited
			if sem != nil {
				sem <- struct{}{}
				defer func() { <-sem }()
			}

			// Create child context
			childCtx := execCtx.copy()

			// Execute child
			childResult := &Result{
				NodeOutputs: make(map[string]interface{}),
			}

			err := i.executeNode(ctx, node, childCtx, childResult)

			// Handle error (capture first error only)
			mu.Lock()
			if err != nil && firstErr == nil {
				firstErr = err
			}

			// Merge child outputs
			if err == nil {
				for k, v := range childResult.NodeOutputs {
					execCtx.set(k, v)
					result.NodeOutputs[k] = v
				}
				result.NodesExecuted += childResult.NodesExecuted
			}
			mu.Unlock()
		}(child)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	if firstErr != nil {
		return newNodeError(def.ID, "parallel", "execute", firstErr)
	}

	i.notifyProgress(def.ID, "completed")
	result.NodesExecuted++

	i.logger.Info("✓ Parallel execution completed",
		"parallel_id", def.ID,
		"child_count", childCount)

	return nil
}
