// Package pantry provides the node type registry.
// Pantry stores all available neta types and provides lookup.
package pantry

import (
	"fmt"
	"sync"

	"bento/pkg/neta"
)

// Pantry is a thread-safe registry of node types.
type Pantry struct {
	mu    sync.RWMutex
	nodes map[string]neta.Executable
}

// New creates a new empty Pantry.
func New() *Pantry {
	return &Pantry{
		nodes: make(map[string]neta.Executable),
	}
}


// Register adds a node type to the pantry.
func (p *Pantry) Register(nodeType string, exec neta.Executable) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.nodes[nodeType]; exists {
		return fmt.Errorf("node type already registered: %s", nodeType)
	}

	p.nodes[nodeType] = exec
	return nil
}

// Get retrieves a node type from the pantry.
func (p *Pantry) Get(nodeType string) (neta.Executable, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	exec, exists := p.nodes[nodeType]
	if !exists {
		return nil, fmt.Errorf("node type not found: %s", nodeType)
	}

	return exec, nil
}

// List returns all registered node types.
func (p *Pantry) List() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	types := make([]string, 0, len(p.nodes))
	for t := range p.nodes {
		types = append(types, t)
	}
	return types
}
