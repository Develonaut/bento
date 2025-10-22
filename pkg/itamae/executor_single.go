package itamae

import (
	"context"
	"time"

	"github.com/Develonaut/bento/pkg/neta"
)

// executeSingle executes a single (non-group) neta.
func (i *Itamae) executeSingle(ctx context.Context, def *neta.Definition, execCtx *executionContext, result *Result) error {
	i.logExecutionStart(def, execCtx)
	netaImpl, err := i.loadNetaImplementation(def)
	if err != nil {
		return err
	}
	if err := i.executeAndRecordNeta(ctx, def, netaImpl, execCtx, result); err != nil {
		return err
	}
	return nil
}

// loadNetaImplementation loads neta from pantry with error wrapping.
func (i *Itamae) loadNetaImplementation(def *neta.Definition) (neta.Executable, error) {
	netaImpl, err := i.pantry.GetNew(def.Type)
	if err != nil {
		return nil, newNodeError(def.ID, def.Type, "get neta", err)
	}
	return netaImpl, nil
}

// executeAndRecordNeta executes neta and records results.
func (i *Itamae) executeAndRecordNeta(ctx context.Context, def *neta.Definition, netaImpl neta.Executable,
	execCtx *executionContext, result *Result) error {
	params := i.prepareNetaParams(def, execCtx)
	output, duration, err := i.executeNetaWithTiming(ctx, netaImpl, params)
	i.sendNodeCompleted(def.ID, duration, err)
	if err != nil {
		return newNodeError(def.ID, def.Type, "execute", err)
	}
	i.storeExecutionResult(def.ID, output, execCtx, result)
	i.logExecutionComplete(def, execCtx, duration)
	return nil
}

// logExecutionStart logs and notifies at the start of node execution.
func (i *Itamae) logExecutionStart(def *neta.Definition, execCtx *executionContext) {
	i.notifyProgress(def.ID, "starting")

	// Mark node as executing in execution state
	i.state.setNodeState(def.ID, "executing")
	i.state.setNodeProgress(def.ID, 0, "Starting")

	if i.messenger != nil {
		i.messenger.SendNodeStarted(def.ID, def.Name, def.Type)
	}

	if i.logger != nil {
		msg := msgNetaStarted()
		i.logger.Debug(msg.format(),
			"neta_id", def.ID,
			"neta_type", def.Type)
	}
}

// prepareNetaParams prepares execution parameters with context resolution.
func (i *Itamae) prepareNetaParams(def *neta.Definition, execCtx *executionContext) map[string]interface{} {
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

// executeNetaWithTiming executes a neta and tracks duration.
func (i *Itamae) executeNetaWithTiming(
	ctx context.Context,
	netaImpl neta.Executable,
	params map[string]interface{},
) (interface{}, time.Duration, error) {
	start := time.Now()
	output, err := netaImpl.Execute(ctx, params)
	duration := time.Since(start)

	if i.slowMoDelay > 0 {
		time.Sleep(i.slowMoDelay)
	}

	return output, duration, err
}

// sendNodeCompleted sends messenger event for completed node.
func (i *Itamae) sendNodeCompleted(nodeID string, duration time.Duration, err error) {
	if i.messenger != nil {
		i.messenger.SendNodeCompleted(nodeID, duration, err)
	}
}

// storeExecutionResult stores node output and marks node as completed.
func (i *Itamae) storeExecutionResult(
	nodeID string,
	output interface{},
	execCtx *executionContext,
	result *Result,
) {
	execCtx.set(nodeID, output)
	result.NodeOutputs[nodeID] = output
	result.NodesExecuted++

	// Mark node as completed in execution state
	i.state.setNodeProgress(nodeID, 100, "Completed")
	i.state.setNodeState(nodeID, "completed")
}

// logExecutionComplete logs completion with progress tracking.
func (i *Itamae) logExecutionComplete(def *neta.Definition, execCtx *executionContext, duration time.Duration) {
	i.notifyProgress(def.ID, "completed")

	if i.logger != nil {
		progressPct := i.state.getProgress()
		durationStr := formatDuration(duration)
		msg := msgChildNodeCompleted(execCtx.depth, def.Type, def.Name, durationStr, progressPct)
		i.logger.Info(msg.format())
	}
}
