package loop

import (
	"context"
	"errors"
	"testing"

	"bento/pkg/neta"
)

// mockExecutor implements Executor for testing.
type mockExecutor struct {
	executeFunc func(ctx context.Context, def neta.Definition) (neta.Result, error)
	callCount   int
}

func (m *mockExecutor) Execute(ctx context.Context, def neta.Definition) (neta.Result, error) {
	m.callCount++
	if m.executeFunc != nil {
		return m.executeFunc(ctx, def)
	}
	return neta.Result{Output: "executed"}, nil
}

func TestFor_Execute(t *testing.T) {
	tests := []struct {
		name      string
		params    map[string]interface{}
		mockError error
		wantCalls int
		wantErr   bool
	}{
		{
			name: "iterate over items",
			params: map[string]interface{}{
				"items": []interface{}{1, 2, 3},
				"body": neta.Definition{
					Type: "test",
					Name: "body",
				},
			},
			wantCalls: 3,
			wantErr:   false,
		},
		{
			name: "empty items array",
			params: map[string]interface{}{
				"items": []interface{}{},
				"body": neta.Definition{
					Type: "test",
					Name: "body",
				},
			},
			wantCalls: 0,
			wantErr:   false,
		},
		{
			name: "missing items parameter",
			params: map[string]interface{}{
				"body": neta.Definition{
					Type: "test",
					Name: "body",
				},
			},
			wantErr: true,
		},
		{
			name: "missing body parameter",
			params: map[string]interface{}{
				"items": []interface{}{1, 2, 3},
			},
			wantErr: true,
		},
		{
			name: "executor error stops iteration",
			params: map[string]interface{}{
				"items": []interface{}{1, 2, 3},
				"body": neta.Definition{
					Type: "test",
					Name: "body",
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
					// Verify item was injected
					if item, ok := def.Parameters["item"]; ok {
						return neta.Result{Output: item}, nil
					}
					return neta.Result{Output: "executed"}, nil
				},
			}

			forNode := NewFor(mock)
			result, err := forNode.Execute(context.Background(), tt.params)

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

func TestGetItems(t *testing.T) {
	tests := []struct {
		name    string
		params  map[string]interface{}
		want    int
		wantErr bool
	}{
		{
			name:    "valid items array",
			params:  map[string]interface{}{"items": []interface{}{1, 2, 3}},
			want:    3,
			wantErr: false,
		},
		{
			name:    "empty items array",
			params:  map[string]interface{}{"items": []interface{}{}},
			want:    0,
			wantErr: false,
		},
		{
			name:    "missing items",
			params:  map[string]interface{}{},
			wantErr: true,
		},
		{
			name:    "wrong type",
			params:  map[string]interface{}{"items": "not an array"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, err := getItems(tt.params)

			if tt.wantErr {
				if err == nil {
					t.Error("getItems() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("getItems() unexpected error: %v", err)
				return
			}

			if len(items) != tt.want {
				t.Errorf("getItems() length = %d, want %d", len(items), tt.want)
			}
		})
	}
}
