package itamae

import (
	"context"
	"fmt"
	"time"

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
		msg := msgLoopStarted(execCtx.depth, def.Name)
		i.logger.Info(msg.format())
	}

	start := time.Now()
	iterations, err := i.executeTimesIterations(ctx, def, count, execCtx, result)
	duration := time.Since(start)

	// Store result even if there was an error (to track partial progress)
	i.storeLoopResult(def, execCtx, result, iterations)

	if i.logger != nil {
		durationStr := formatDuration(duration)
		msg := msgLoopCompleted(execCtx.depth, def.Name, durationStr)
		i.logger.Info(msg.format())
	}

	// Return error if iterations failed
	if err != nil {
		return err
	}

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
) (int, error) {
	iterations := 0
	for idx := 0; idx < count; idx++ {
		if err := i.executeSingleIteration(ctx, def, idx, nil, execCtx, result); err != nil {
			if i.logger != nil {
				// Use proper indentation for error messages inside loops (depth 1)
				i.logger.Error("│  │   ✗ Times loop iteration failed",
					"loop_id", def.ID,
					"iteration", idx,
					"error", err)
			}
			return iterations, err
		}
		iterations++
	}
	return iterations, nil
}
