// Package itamae provides the orchestration engine for executing neta definitions.
// Itamae (板前) means "sushi chef" - the one who prepares each piece.
package itamae

import (
	"context"
	"fmt"

	"bento/pkg/neta"
)

// Itamae orchestrates the execution of neta definitions.
type Itamae struct {
	pantry Registry
}

// Registry provides node type lookup.
type Registry interface {
	Get(nodeType string) (neta.Executable, error)
}

// New creates a new Itamae with the provided registry.
func New(registry Registry) *Itamae {
	return &Itamae{
		pantry: registry,
	}
}

// Execute runs a neta definition and returns the result.
func (i *Itamae) Execute(ctx context.Context, def neta.Definition) (neta.Result, error) {
	if def.IsGroup() {
		return i.executeGroup(ctx, def)
	}
	return i.executeSingle(ctx, def)
}

// executeSingle runs a single node.
func (i *Itamae) executeSingle(ctx context.Context, def neta.Definition) (neta.Result, error) {
	exec, err := i.pantry.Get(def.Type)
	if err != nil {
		return neta.Result{}, fmt.Errorf("node type not found: %s: %w", def.Type, err)
	}
	return exec.Execute(ctx, def.Parameters)
}

// executeGroup runs a group of nodes in sequence.
func (i *Itamae) executeGroup(ctx context.Context, def neta.Definition) (neta.Result, error) {
	results := make([]neta.Result, 0, len(def.Nodes))
	for _, child := range def.Nodes {
		result, err := i.Execute(ctx, child)
		if err != nil {
			return neta.Result{}, err
		}
		results = append(results, result)
	}
	return neta.Result{Output: results}, nil
}
