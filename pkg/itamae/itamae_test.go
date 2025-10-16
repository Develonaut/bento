package itamae

import (
	"context"
	"errors"
	"testing"

	"bento/pkg/neta"
)

// mockRegistry is a test implementation of Registry
type mockRegistry struct {
	nodes map[string]neta.Executable
}

func newMockRegistry() *mockRegistry {
	return &mockRegistry{
		nodes: make(map[string]neta.Executable),
	}
}

func (m *mockRegistry) register(nodeType string, exec neta.Executable) {
	m.nodes[nodeType] = exec
}

func (m *mockRegistry) Get(nodeType string) (neta.Executable, error) {
	exec, exists := m.nodes[nodeType]
	if !exists {
		return nil, errors.New("node type not found")
	}
	return exec, nil
}

// mockExecutable is a test implementation of Executable
type mockExecutable struct {
	output interface{}
	err    error
}

func (m *mockExecutable) Execute(ctx context.Context, params map[string]interface{}) (neta.Result, error) {
	if m.err != nil {
		return neta.Result{}, m.err
	}
	return neta.Result{Output: m.output}, nil
}

// mockInputPassthrough is a test executable that passes through its input parameter
type mockInputPassthrough struct{}

func (m *mockInputPassthrough) Execute(ctx context.Context, params map[string]interface{}) (neta.Result, error) {
	input, ok := params["input"]
	if !ok {
		return neta.Result{Output: "no input"}, nil
	}
	return neta.Result{Output: input}, nil
}

func TestItamae_ExecuteSingle(t *testing.T) {
	tests := []struct {
		name    string
		def     neta.Definition
		setup   func(*mockRegistry)
		want    interface{}
		wantErr bool
	}{
		{
			name: "successful single node execution",
			def: neta.Definition{
				Type: "test",
				Name: "Test Node",
				Parameters: map[string]interface{}{
					"key": "value",
				},
			},
			setup: func(r *mockRegistry) {
				r.register("test", &mockExecutable{output: "success"})
			},
			want:    "success",
			wantErr: false,
		},
		{
			name: "node type not found",
			def: neta.Definition{
				Type: "unknown",
				Name: "Unknown Node",
			},
			setup:   func(r *mockRegistry) {},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := newMockRegistry()
			tt.setup(registry)

			itamae := New(registry)
			ctx := context.Background()

			result, err := itamae.Execute(ctx, tt.def)

			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result.Output != tt.want {
				t.Errorf("Execute() Output = %v, want %v", result.Output, tt.want)
			}
		})
	}
}

func TestItamae_ExecuteGroup(t *testing.T) {
	registry := newMockRegistry()
	registry.register("test", &mockExecutable{output: "success"})

	itamae := New(registry)
	ctx := context.Background()

	groupDef := neta.Definition{
		Type: "group",
		Name: "Test Group",
		Nodes: []neta.Definition{
			{Type: "test", Name: "Step 1"},
			{Type: "test", Name: "Step 2"},
		},
	}

	result, err := itamae.Execute(ctx, groupDef)
	if err != nil {
		t.Errorf("Execute() error = %v", err)
		return
	}

	results, ok := result.Output.([]neta.Result)
	if !ok {
		t.Errorf("Execute() Output type = %T, want []neta.Result", result.Output)
		return
	}

	if len(results) != 2 {
		t.Errorf("Execute() result count = %d, want 2", len(results))
	}
}

func TestItamae_ExecuteGraphWithDataFlow(t *testing.T) {
	registry := newMockRegistry()
	// Node 1 outputs "hello"
	registry.register("source", &mockExecutable{output: "hello"})
	// Node 2 passes through its input
	registry.register("passthrough", &mockInputPassthrough{})

	store := neta.NewExecutionGraphStore()
	itamae := New(registry)
	itamae.SetStore(store)

	ctx := context.Background()

	graphDef := neta.Definition{
		Type: "group.sequence",
		Name: "Data Flow Test",
		Nodes: []neta.Definition{
			{
				ID:   "node-1",
				Type: "source",
				Name: "Source Node",
			},
			{
				ID:   "node-2",
				Type: "passthrough",
				Name: "Passthrough Node",
			},
		},
		Edges: []neta.NodeEdge{
			{
				ID:     "edge-1",
				Source: "node-1",
				Target: "node-2",
			},
		},
	}

	result, err := itamae.Execute(ctx, graphDef)
	if err != nil {
		t.Errorf("Execute() error = %v", err)
		return
	}

	// The result should be "hello" from node-1, passed through node-2
	if result.Output != "hello" {
		t.Errorf("Execute() Output = %v, want 'hello'", result.Output)
	}
}
