package itamae

import (
	"context"
	"fmt"
	"time"

	"github.com/Develonaut/bento/pkg/neta"
)

// executeForEach executes a forEach loop AS A LEAF NODE.
// The loop is a single unit in the progress graph, with children executing internally.
func (i *Itamae) executeForEach(
	ctx context.Context,
	def *neta.Definition,
	execCtx *executionContext,
	result *Result,
) error {
	items, err := i.extractLoopItems(def, execCtx)
	if err != nil {
		i.state.setNodeState(def.ID, "error")
		return err
	}

	i.initializeLoopExecution(def, execCtx)

	start := time.Now()
	loopResults, err := i.executeForEachIterations(ctx, def, items, execCtx)
	if err != nil {
		i.state.setNodeState(def.ID, "error")
		return err
	}

	i.finalizeLoopExecution(def, loopResults, execCtx, result, time.Since(start))
	return nil
}

// initializeLoopExecution sets up loop state and logging.
func (i *Itamae) initializeLoopExecution(def *neta.Definition, execCtx *executionContext) {
	i.state.setNodeState(def.ID, "executing")
	i.state.setNodeProgress(def.ID, 0, "Starting loop")

	if i.logger != nil {
		msg := msgLoopStarted(execCtx.depth, def.Name)
		i.logger.Info(msg.format())
	}
}

// executeForEachIterations executes all loop iterations and returns results.
func (i *Itamae) executeForEachIterations(
	ctx context.Context,
	def *neta.Definition,
	items []interface{},
	execCtx *executionContext,
) ([]interface{}, error) {
	loopResults := make([]interface{}, 0, len(items))
	totalItems := len(items)

	for idx, item := range items {
		if err := i.checkLoopCancellation(ctx, def); err != nil {
			return nil, err
		}

		i.reportLoopProgress(def, idx, totalItems)

		iterResult, err := i.executeLoopIteration(ctx, def, item, idx, totalItems, execCtx)
		if err != nil {
			i.logLoopIterationError(def, idx, err)
			return nil, fmt.Errorf("iteration %d failed: %w", idx, err)
		}

		loopResults = append(loopResults, iterResult)
	}

	return loopResults, nil
}

// checkLoopCancellation checks if context is cancelled.
func (i *Itamae) checkLoopCancellation(ctx context.Context, def *neta.Definition) error {
	select {
	case <-ctx.Done():
		i.state.setNodeState(def.ID, "error")
		return ctx.Err()
	default:
		return nil
	}
}

// reportLoopProgress reports partial progress for current iteration.
func (i *Itamae) reportLoopProgress(def *neta.Definition, idx, total int) {
	progress := (idx * 100) / total
	message := fmt.Sprintf("Iteration %d/%d", idx+1, total)
	i.state.setNodeProgress(def.ID, progress, message)

	if i.logger != nil {
		i.logger.Debug("Loop iteration",
			"loop_id", def.ID,
			"iteration", idx+1,
			"total", total)
	}
}

// logLoopIterationError logs an error during loop iteration.
func (i *Itamae) logLoopIterationError(def *neta.Definition, idx int, err error) {
	if i.logger != nil {
		i.logger.Error("│  │   ✗ Loop iteration failed",
			"loop_id", def.ID,
			"iteration", idx,
			"error", err)
	}
}

// finalizeLoopExecution completes loop execution and stores results.
func (i *Itamae) finalizeLoopExecution(
	def *neta.Definition,
	loopResults []interface{},
	execCtx *executionContext,
	result *Result,
	duration time.Duration,
) {
	i.state.setNodeProgress(def.ID, 100, "Completed")
	i.state.setNodeState(def.ID, "completed")

	execCtx.set(def.ID, loopResults)
	result.NodeOutputs[def.ID] = loopResults
	result.NodesExecuted++

	if i.logger != nil {
		durationStr := formatDuration(duration)
		progressPct := i.state.getProgress()
		msg := msgLoopCompleted(execCtx.depth, def.Name, durationStr, progressPct)
		i.logger.Info(msg.format())
	}

	i.notifyProgress(def.ID, "completed")
}

// executeLoopIteration executes one iteration (INTERNAL - not tracked in graph).
func (i *Itamae) executeLoopIteration(
	ctx context.Context,
	def *neta.Definition,
	item interface{},
	idx int,
	total int,
	execCtx *executionContext,
) (map[string]interface{}, error) {
	iterCtx := execCtx.withDepth(1)
	iterCtx.set("item", item)
	iterCtx.set("index", idx)

	return i.executeIterationChildren(ctx, def, idx, total, iterCtx)
}

// executeIterationChildren executes all child nodes for one iteration.
func (i *Itamae) executeIterationChildren(
	ctx context.Context,
	def *neta.Definition,
	idx int,
	total int,
	iterCtx *executionContext,
) (map[string]interface{}, error) {
	iterResult := make(map[string]interface{})

	for j := range def.Nodes {
		childDef := &def.Nodes[j]

		// Notify messenger which child is currently executing
		if i.messenger != nil {
			i.messenger.SendLoopChild(def.ID, childDef.Name, idx, total)
		}

		output, err := i.executeNodeInternal(ctx, childDef, iterCtx)
		if err != nil {
			return nil, err
		}

		iterResult[childDef.ID] = output
		iterCtx.set(childDef.ID, output)
	}

	return iterResult, nil
}

// executeNodeInternal executes a node without state tracking (for loop children).
func (i *Itamae) executeNodeInternal(
	ctx context.Context,
	def *neta.Definition,
	execCtx *executionContext,
) (interface{}, error) {
	i.logInternalNodeStart(def, execCtx)

	netaImpl, err := i.loadNetaForInternal(def)
	if err != nil {
		return nil, err
	}

	params := i.prepareInternalNodeParams(def, execCtx)

	output, duration, err := i.executeInternalNode(ctx, netaImpl, params)
	if err != nil {
		return nil, newNodeError(def.ID, def.Type, "execute", err)
	}

	i.logInternalNodeComplete(def, execCtx, duration)
	return output, nil
}

// logInternalNodeStart logs execution start for internal node.
func (i *Itamae) logInternalNodeStart(def *neta.Definition, execCtx *executionContext) {
	if i.logger != nil {
		msg := msgChildNodeStarted(execCtx.depth, def.Type, def.Name)
		i.logger.Info(msg.format())
	}
}

// loadNetaForInternal loads neta implementation for internal execution.
func (i *Itamae) loadNetaForInternal(def *neta.Definition) (neta.Executable, error) {
	netaImpl, err := i.pantry.GetNew(def.Type)
	if err != nil {
		return nil, newNodeError(def.ID, def.Type, "get neta", err)
	}
	return netaImpl, nil
}

// prepareInternalNodeParams prepares parameters for internal node execution.
func (i *Itamae) prepareInternalNodeParams(
	def *neta.Definition,
	execCtx *executionContext,
) map[string]interface{} {
	params := make(map[string]interface{})
	for k, v := range def.Parameters {
		params[k] = execCtx.resolveValue(v)
	}
	params["_context"] = execCtx.toMap()
	params["_onOutput"] = func(line string) {
		if i.logger != nil {
			i.logger.Stream(line)
		}
	}
	return params
}

// executeInternalNode executes neta and tracks duration.
func (i *Itamae) executeInternalNode(
	ctx context.Context,
	netaImpl neta.Executable,
	params map[string]interface{},
) (interface{}, time.Duration, error) {
	start := time.Now()
	output, err := netaImpl.Execute(ctx, params)
	return output, time.Since(start), err
}

// logInternalNodeComplete logs completion for internal node.
func (i *Itamae) logInternalNodeComplete(
	def *neta.Definition,
	execCtx *executionContext,
	duration time.Duration,
) {
	if i.logger != nil {
		durationStr := formatDuration(duration)
		progressPct := i.state.getProgress()
		msg := msgChildNodeCompleted(execCtx.depth, def.Type, def.Name, durationStr, progressPct)
		i.logger.Info(msg.format())
	}
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
