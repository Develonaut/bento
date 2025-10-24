package itamae

import (
	"context"
	"time"

	"github.com/Develonaut/bento/pkg/neta"
)

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
		msg := msgChildNodeStarted(execCtx.getBreadcrumb(), def.Type, def.Name)
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
			// Stream output with breadcrumb context
			breadcrumb := execCtx.getBreadcrumb()
			if breadcrumb != "" {
				formattedLine := formatStreamingOutput(breadcrumb, line)
				i.logger.Info(formattedLine)
			} else {
				i.logger.Info(line)
			}
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
		msg := msgChildNodeCompleted(execCtx.getBreadcrumb(), def.Type, def.Name, durationStr, progressPct)
		i.logger.Info(msg.format())
	}
}
