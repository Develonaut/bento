package itamae

import (
	"context"
	"fmt"

	"github.com/Develonaut/bento/pkg/neta"
)

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
