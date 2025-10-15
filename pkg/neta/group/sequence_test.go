package group

import (
	"context"
	"errors"
	"sync"
	"testing"

	"bento/pkg/neta"
)

// mockExecutor implements Executor for testing.
type mockExecutor struct {
	executeFunc func(ctx context.Context, def neta.Definition) (neta.Result, error)
	mu          sync.Mutex
	callCount   int
}

func (m *mockExecutor) Execute(ctx context.Context, def neta.Definition) (neta.Result, error) {
	m.mu.Lock()
	m.callCount++
	m.mu.Unlock()

	if m.executeFunc != nil {
		return m.executeFunc(ctx, def)
	}
	return neta.Result{Output: "executed"}, nil
}

func TestSequence_Execute(t *testing.T) {
	tests := []struct {
		name      string
		params    map[string]interface{}
		mockError error
		wantCalls int
		wantErr   bool
	}{
		{
			name: "execute multiple nodes in sequence",
			params: map[string]interface{}{
				"nodes": []neta.Definition{
					{Type: "test1", Name: "node1"},
					{Type: "test2", Name: "node2"},
					{Type: "test3", Name: "node3"},
				},
			},
			wantCalls: 3,
			wantErr:   false,
		},
		{
			name: "execute single node",
			params: map[string]interface{}{
				"nodes": []neta.Definition{
					{Type: "test", Name: "node"},
				},
			},
			wantCalls: 1,
			wantErr:   false,
		},
		{
			name: "empty nodes array",
			params: map[string]interface{}{
				"nodes": []neta.Definition{},
			},
			wantCalls: 0,
			wantErr:   false,
		},
		{
			name:      "missing nodes parameter",
			params:    map[string]interface{}{},
			wantCalls: 0,
			wantErr:   false,
		},
		{
			name: "error stops execution",
			params: map[string]interface{}{
				"nodes": []neta.Definition{
					{Type: "test1", Name: "node1"},
					{Type: "test2", Name: "node2"},
					{Type: "test3", Name: "node3"},
				},
			},
			mockError: errors.New("execution failed"),
			wantCalls: 1,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockExecutor{
				executeFunc: func(ctx context.Context, def neta.Definition) (neta.Result, error) {
					if tt.mockError != nil {
						return neta.Result{}, tt.mockError
					}
					return neta.Result{Output: def.Name}, nil
				},
			}

			seq := NewSequence(mock)
			result, err := seq.Execute(context.Background(), tt.params)

			if tt.wantErr {
				if err == nil {
					t.Error("Execute() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Execute() unexpected error: %v", err)
				return
			}

			if mock.callCount != tt.wantCalls {
				t.Errorf("Execute() call count = %d, want %d", mock.callCount, tt.wantCalls)
			}

			// Verify results array length
			if results, ok := result.Output.([]neta.Result); ok {
				if len(results) != tt.wantCalls {
					t.Errorf("Execute() results length = %d, want %d", len(results), tt.wantCalls)
				}
			}
		})
	}
}

func TestGetNodes(t *testing.T) {
	tests := []struct {
		name   string
		params map[string]interface{}
		want   int
	}{
		{
			name: "valid nodes array",
			params: map[string]interface{}{
				"nodes": []neta.Definition{
					{Type: "test1"},
					{Type: "test2"},
				},
			},
			want: 2,
		},
		{
			name: "empty nodes array",
			params: map[string]interface{}{
				"nodes": []neta.Definition{},
			},
			want: 0,
		},
		{
			name:   "missing nodes",
			params: map[string]interface{}{},
			want:   0,
		},
		{
			name: "wrong type",
			params: map[string]interface{}{
				"nodes": "not an array",
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes := getNodes(tt.params)
			if len(nodes) != tt.want {
				t.Errorf("getNodes() length = %d, want %d", len(nodes), tt.want)
			}
		})
	}
}
