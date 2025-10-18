// Package bento provides bento composition - bentos calling other bentos.
package bento

import (
	"context"
	"fmt"

	"bento/pkg/neta"
)

// Store interface for loading bentos.
type Store interface {
	Load(name string) (neta.Definition, error)
}

// Chef interface for executing bentos.
type Chef interface {
	Execute(ctx context.Context, def neta.Definition) (neta.Result, error)
}

// Execute loads and executes another bento as a node.
// This enables composable bento architecture - "a node is a node is a node".
type Execute struct {
	store Store
	chef  Chef
}

// NewExecute creates a bento executor node.
func NewExecute(store Store, chef Chef) *Execute {
	return &Execute{
		store: store,
		chef:  chef,
	}
}

// Execute runs another bento as a node.
func (e *Execute) Execute(ctx context.Context, params map[string]interface{}) (neta.Result, error) {
	bentoName := neta.GetStringParam(params, "bento", "")
	if bentoName == "" {
		return neta.Result{}, fmt.Errorf("bento parameter required")
	}

	// Load the bento definition
	def, err := e.store.Load(bentoName)
	if err != nil {
		return neta.Result{}, fmt.Errorf("failed to load bento %s: %w", bentoName, err)
	}

	// Extract inputs for the sub-bento
	if inputs, ok := params["inputs"].(map[string]interface{}); ok {
		// Merge inputs into the bento's parameters
		if def.Parameters == nil {
			def.Parameters = make(map[string]interface{})
		}
		for key, value := range inputs {
			def.Parameters[key] = value
		}
	}

	// Execute the bento (potentially recursive!)
	result, err := e.chef.Execute(ctx, def)
	if err != nil {
		return neta.Result{}, fmt.Errorf("failed to execute bento %s: %w", bentoName, err)
	}

	return result, nil
}
