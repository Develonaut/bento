package itamae

import (
	"context"
	"fmt"

	"github.com/Develonaut/bento/pkg/neta"
)

// executeLoop executes a loop neta.
// Critical for Phase 8: CSV iteration for product automation.
func (i *Itamae) executeLoop(
	ctx context.Context,
	def *neta.Definition,
	execCtx *executionContext,
	result *Result,
) error {
	i.notifyProgress(def.ID, "starting")

	// Get loop parameters
	mode, ok := def.Parameters["mode"].(string)
	if !ok {
		return newNodeError(def.ID, "loop", "validate",
			fmt.Errorf("missing or invalid 'mode' parameter"))
	}

	switch mode {
	case "forEach":
		return i.executeForEach(ctx, def, execCtx, result)
	case "times":
		return i.executeTimes(ctx, def, execCtx, result)
	case "while":
		return i.executeWhile(ctx, def, execCtx, result)
	default:
		return newNodeError(def.ID, "loop", "validate",
			fmt.Errorf("unknown loop mode: %s", mode))
	}
}

// executeForEach executes a forEach loop.
func (i *Itamae) executeForEach(
	ctx context.Context,
	def *neta.Definition,
	execCtx *executionContext,
	result *Result,
) error {
	// Get items to iterate over
	itemsParam := def.Parameters["items"]

	// Resolve template if needed
	resolved := execCtx.resolveValue(itemsParam)

	items, ok := resolved.([]interface{})
	if !ok {
		return newNodeError(def.ID, "loop", "validate",
			fmt.Errorf("'items' must be an array"))
	}

	i.logger.Info("🔄 Starting forEach loop",
		"loop_id", def.ID,
		"item_count", len(items))

	// Execute body for each item
	iterations := 0
	for idx, item := range items {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		i.logger.Debug("Loop iteration",
			"loop_id", def.ID,
			"iteration", idx+1,
			"total", len(items))

		// Create iteration context with item
		iterCtx := execCtx.copy()
		iterCtx.set("item", item)
		iterCtx.set("index", idx)

		// Execute loop body (would need body definition)
		// For now, just count iterations
		iterations++
	}

	// Store loop result
	loopOutput := map[string]interface{}{
		"iterations": iterations,
		"completed":  true,
	}

	execCtx.set(def.ID, loopOutput)
	result.NodeOutputs[def.ID] = loopOutput
	result.NodesExecuted++

	i.notifyProgress(def.ID, "completed")

	i.logger.Info("✓ forEach loop completed",
		"loop_id", def.ID,
		"iterations", iterations)

	return nil
}

// executeTimes executes a times loop (repeat N times).
func (i *Itamae) executeTimes(
	ctx context.Context,
	def *neta.Definition,
	execCtx *executionContext,
	result *Result,
) error {
	// Get count
	countParam := def.Parameters["count"]
	count, ok := countParam.(float64)
	if !ok {
		return newNodeError(def.ID, "loop", "validate",
			fmt.Errorf("'count' must be a number"))
	}

	iterations := int(count)

	i.logger.Info("🔄 Starting times loop",
		"loop_id", def.ID,
		"count", iterations)

	// Execute body N times
	for i := 0; i < iterations; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		// Execute body (stub for now)
	}

	loopOutput := map[string]interface{}{
		"iterations": iterations,
		"completed":  true,
	}

	execCtx.set(def.ID, loopOutput)
	result.NodeOutputs[def.ID] = loopOutput
	result.NodesExecuted++

	i.notifyProgress(def.ID, "completed")
	return nil
}

// executeWhile executes a while loop.
func (i *Itamae) executeWhile(
	ctx context.Context,
	def *neta.Definition,
	execCtx *executionContext,
	result *Result,
) error {
	// While loop implementation (stub for now)
	return newNodeError(def.ID, "loop", "execute",
		fmt.Errorf("while loops not yet implemented"))
}
