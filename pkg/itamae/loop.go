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
	items, err := i.extractLoopItems(def, execCtx)
	if err != nil {
		return err
	}

	i.logger.Info("🔄 Starting forEach loop",
		"loop_id", def.ID,
		"item_count", len(items))

	iterations := i.executeLoopIterations(ctx, def, items, execCtx, result)

	i.storeLoopResult(def, execCtx, result, iterations)

	i.logger.Info("✓ forEach loop completed",
		"loop_id", def.ID,
		"iterations", iterations)

	return nil
}

// extractLoopItems extracts and validates items for forEach loop.
func (i *Itamae) extractLoopItems(
	def *neta.Definition,
	execCtx *executionContext,
) ([]interface{}, error) {
	itemsParam := def.Parameters["items"]

	i.logger.Debug("Loop items parameter",
		"loop_id", def.ID,
		"itemsParam", itemsParam,
		"itemsParam_type", fmt.Sprintf("%T", itemsParam))

	resolved := execCtx.resolveValue(itemsParam)

	i.logger.Debug("Loop items resolved",
		"loop_id", def.ID,
		"resolved", resolved,
		"resolved_type", fmt.Sprintf("%T", resolved))

	return i.convertToInterfaceArray(def, resolved)
}

// convertToInterfaceArray converts resolved value to []interface{}.
func (i *Itamae) convertToInterfaceArray(
	def *neta.Definition,
	resolved interface{},
) ([]interface{}, error) {
	switch v := resolved.(type) {
	case []interface{}:
		return v, nil
	case []map[string]interface{}:
		items := make([]interface{}, len(v))
		for idx, item := range v {
			items[idx] = item
		}
		return items, nil
	default:
		i.logger.Error("Loop items not an array",
			"loop_id", def.ID,
			"resolved_type", fmt.Sprintf("%T", resolved),
			"resolved_value", resolved)
		return nil, newNodeError(def.ID, "loop", "validate",
			fmt.Errorf("'items' must be an array, got %T", resolved))
	}
}

// executeLoopIterations executes loop body for each item.
func (i *Itamae) executeLoopIterations(
	ctx context.Context,
	def *neta.Definition,
	items []interface{},
	execCtx *executionContext,
	result *Result,
) int {
	iterations := 0
	for idx, item := range items {
		if err := i.executeSingleIteration(ctx, def, idx, item, execCtx, result); err != nil {
			i.logger.Error("Loop iteration failed",
				"loop_id", def.ID,
				"iteration", idx,
				"error", err)
			return iterations
		}
		iterations++
	}
	return iterations
}

// executeSingleIteration executes loop body for one iteration.
func (i *Itamae) executeSingleIteration(
	ctx context.Context,
	def *neta.Definition,
	idx int,
	item interface{},
	execCtx *executionContext,
	result *Result,
) error {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	i.logger.Debug("Loop iteration",
		"loop_id", def.ID,
		"iteration", idx+1,
		"total", len(def.Nodes))

	iterCtx := i.createIterationContext(execCtx, item, idx)

	return i.executeIterationBody(ctx, def, iterCtx, result, idx)
}

// createIterationContext creates context for a single iteration.
func (i *Itamae) createIterationContext(
	execCtx *executionContext,
	item interface{},
	idx int,
) *executionContext {
	iterCtx := execCtx.copy()
	iterCtx.set("item", item)
	iterCtx.set("index", idx)
	return iterCtx
}

// executeIterationBody executes nested nodes for one iteration.
func (i *Itamae) executeIterationBody(
	ctx context.Context,
	def *neta.Definition,
	iterCtx *executionContext,
	result *Result,
	idx int,
) error {
	for j := range def.Nodes {
		childDef := &def.Nodes[j]
		if err := i.executeNode(ctx, childDef, iterCtx, result); err != nil {
			return fmt.Errorf("iteration %d failed: %w", idx, err)
		}
	}
	return nil
}

// storeLoopResult stores the loop execution result.
func (i *Itamae) storeLoopResult(
	def *neta.Definition,
	execCtx *executionContext,
	result *Result,
	iterations int,
) {
	loopOutput := map[string]interface{}{
		"iterations": iterations,
		"completed":  true,
	}

	execCtx.set(def.ID, loopOutput)
	result.NodeOutputs[def.ID] = loopOutput
	result.NodesExecuted++

	i.notifyProgress(def.ID, "completed")
}

// executeTimes executes a times loop (repeat N times).
func (i *Itamae) executeTimes(
	ctx context.Context,
	def *neta.Definition,
	execCtx *executionContext,
	result *Result,
) error {
	count, err := i.extractTimesCount(def)
	if err != nil {
		return err
	}

	i.logger.Info("🔄 Starting times loop",
		"loop_id", def.ID,
		"count", count)

	iterations := i.executeTimesIterations(ctx, def, count, execCtx, result)

	i.storeLoopResult(def, execCtx, result, iterations)

	return nil
}

// extractTimesCount extracts and validates count parameter for times loop.
func (i *Itamae) extractTimesCount(def *neta.Definition) (int, error) {
	countParam := def.Parameters["count"]
	count, ok := countParam.(float64)
	if !ok {
		return 0, newNodeError(def.ID, "loop", "validate",
			fmt.Errorf("'count' must be a number"))
	}
	return int(count), nil
}

// executeTimesIterations executes N iterations for times loop.
func (i *Itamae) executeTimesIterations(
	ctx context.Context,
	def *neta.Definition,
	count int,
	execCtx *executionContext,
	result *Result,
) int {
	iterations := 0
	for idx := 0; idx < count; idx++ {
		if err := i.executeSingleIteration(ctx, def, idx, nil, execCtx, result); err != nil {
			i.logger.Error("Times loop iteration failed",
				"loop_id", def.ID,
				"iteration", idx,
				"error", err)
			return iterations
		}
		iterations++
	}
	return iterations
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
