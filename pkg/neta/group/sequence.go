// Package group provides group execution nodes.
package group

import (
	"context"

	"bento/pkg/neta"
)

// Sequence executes nodes in order.
type Sequence struct {
	executor neta.Executor
}

// NewSequence creates a sequential group executor.
func NewSequence(executor neta.Executor) *Sequence {
	return &Sequence{executor: executor}
}

// Execute runs nodes one after another.
func (s *Sequence) Execute(ctx context.Context, params map[string]interface{}) (neta.Result, error) {
	nodes := getNodes(params)
	results := make([]neta.Result, 0, len(nodes))

	for _, node := range nodes {
		result, err := s.executor.Execute(ctx, node)
		if err != nil {
			return neta.Result{}, err
		}
		results = append(results, result)
	}

	return neta.Result{Output: results}, nil
}

// getNodes extracts child nodes from params.
func getNodes(params map[string]interface{}) []neta.Definition {
	nodes, ok := params["nodes"].([]neta.Definition)
	if !ok {
		return []neta.Definition{}
	}
	return nodes
}
