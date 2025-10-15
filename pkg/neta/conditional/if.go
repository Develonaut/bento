// Package conditional provides conditional execution nodes.
package conditional

import (
	"context"

	"bento/pkg/neta"
)

// If executes nodes based on conditions.
type If struct {
	executor neta.Executor
}

// NewIf creates a new If node.
func NewIf(executor neta.Executor) *If {
	return &If{executor: executor}
}

// Execute evaluates condition and runs appropriate branch.
func (i *If) Execute(ctx context.Context, params map[string]interface{}) (neta.Result, error) {
	condition := getBoolParam(params, "condition", false)

	if condition {
		return i.executeBranch(ctx, params, "then")
	}
	return i.executeBranch(ctx, params, "else")
}

// executeBranch runs the specified branch.
func (i *If) executeBranch(ctx context.Context, params map[string]interface{}, branch string) (neta.Result, error) {
	def, ok := params[branch].(neta.Definition)
	if !ok {
		// Branch is optional - return empty result if not provided
		return neta.Result{}, nil
	}
	return i.executor.Execute(ctx, def)
}

// getBoolParam extracts a bool parameter.
func getBoolParam(params map[string]interface{}, key string, defaultVal bool) bool {
	if val, ok := params[key].(bool); ok {
		return val
	}
	return defaultVal
}
