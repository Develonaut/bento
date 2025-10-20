package itamae

import (
	"context"
	"fmt"

	"github.com/Develonaut/bento/pkg/neta"
)

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

	if i.logger != nil {
		msg := msgLoopStarted("forEach")
		i.logger.Info(msg.format(),
			"loop_id", def.ID,
			"loop_name", def.Name,
			"item_count", len(items))
	}

	iterations := i.executeLoopIterations(ctx, def, items, execCtx, result)

	i.storeLoopResult(def, execCtx, result, iterations)

	if i.logger != nil {
		msg := msgLoopCompleted("forEach")
		i.logger.Info(msg.format(),
			"loop_id", def.ID,
			"loop_name", def.Name,
			"iterations", iterations)
	}

	return nil
}

// extractLoopItems extracts and validates items for forEach loop.
func (i *Itamae) extractLoopItems(
	def *neta.Definition,
	execCtx *executionContext,
) ([]interface{}, error) {
	itemsParam := def.Parameters["items"]

	if i.logger != nil {
		i.logger.Debug("Loop items parameter",
			"loop_id", def.ID,
			"itemsParam", itemsParam,
			"itemsParam_type", fmt.Sprintf("%T", itemsParam))
	}

	resolved := execCtx.resolveValue(itemsParam)

	if i.logger != nil {
		i.logger.Debug("Loop items resolved",
			"loop_id", def.ID,
			"resolved", resolved,
			"resolved_type", fmt.Sprintf("%T", resolved))
	}

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
		if i.logger != nil {
			i.logger.Error("Loop items not an array",
				"loop_id", def.ID,
				"resolved_type", fmt.Sprintf("%T", resolved),
				"resolved_value", resolved)
		}
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
			if i.logger != nil {
				i.logger.Error("Loop iteration failed",
					"loop_id", def.ID,
					"iteration", idx,
					"error", err)
			}
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

	if i.logger != nil {
		i.logger.Debug("Loop iteration",
			"loop_id", def.ID,
			"iteration", idx+1,
			"total", len(def.Nodes))
	}

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
