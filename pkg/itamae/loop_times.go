package itamae

import (
	"context"
	"fmt"

	"github.com/Develonaut/bento/pkg/neta"
)

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

	if i.logger != nil {
		msg := msgLoopStarted("times")
		i.logger.Info(msg.format(),
			"loop_id", def.ID,
			"count", count)
	}

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
			if i.logger != nil {
				i.logger.Error("Times loop iteration failed",
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
