// Package group provides group execution nodes.
package group

import (
	"context"
	"sync"

	"bento/pkg/neta"
)

// Parallel executes nodes concurrently.
type Parallel struct {
	executor neta.Executor
}

// NewParallel creates a parallel group executor.
func NewParallel(executor neta.Executor) *Parallel {
	return &Parallel{executor: executor}
}

// Execute runs nodes concurrently.
func (p *Parallel) Execute(ctx context.Context, params map[string]interface{}) (neta.Result, error) {
	nodes := getNodes(params)
	if len(nodes) == 0 {
		return neta.Result{Output: []neta.Result{}}, nil
	}

	return p.executeParallel(ctx, nodes)
}

// executeParallel runs all nodes concurrently.
func (p *Parallel) executeParallel(ctx context.Context, nodes []neta.Definition) (neta.Result, error) {
	results := make([]neta.Result, len(nodes))
	errs := make([]error, len(nodes))
	var wg sync.WaitGroup

	for i, node := range nodes {
		wg.Add(1)
		go func(idx int, n neta.Definition) {
			defer wg.Done()
			result, err := p.executor.Execute(ctx, n)
			results[idx] = result
			errs[idx] = err
		}(i, node)
	}

	wg.Wait()

	// Check for errors
	for _, err := range errs {
		if err != nil {
			return neta.Result{}, err
		}
	}

	return neta.Result{Output: results}, nil
}
