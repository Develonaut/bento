// Package loop provides iteration nodes.
package loop

import (
	"context"
	"fmt"

	"bento/pkg/neta"
)

// For iterates over a collection.
type For struct {
	executor neta.Executor
}

// NewFor creates a new For loop node.
func NewFor(executor neta.Executor) *For {
	return &For{executor: executor}
}

// Execute iterates over items and executes body for each.
func (f *For) Execute(ctx context.Context, params map[string]interface{}) (neta.Result, error) {
	items, err := getItems(params)
	if err != nil {
		return neta.Result{}, err
	}

	body, ok := params["body"].(neta.Definition)
	if !ok {
		return neta.Result{}, fmt.Errorf("body parameter required")
	}

	return f.iterate(ctx, items, body)
}

// iterate executes body for each item.
func (f *For) iterate(ctx context.Context, items []interface{}, body neta.Definition) (neta.Result, error) {
	results := make([]neta.Result, 0, len(items))

	for _, item := range items {
		// Create body with injected item
		bodyWithItem := body
		if bodyWithItem.Parameters == nil {
			bodyWithItem.Parameters = make(map[string]interface{})
		}
		bodyWithItem.Parameters["item"] = item

		result, err := f.executor.Execute(ctx, bodyWithItem)
		if err != nil {
			return neta.Result{}, err
		}
		results = append(results, result)
	}

	return neta.Result{Output: results}, nil
}

// getItems extracts the items array from params.
func getItems(params map[string]interface{}) ([]interface{}, error) {
	items, ok := params["items"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("items parameter required and must be array")
	}
	return items, nil
}
