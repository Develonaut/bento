package group

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"bento/pkg/neta"
)

func TestParallel_Execute(t *testing.T) {
	tests := []struct {
		name      string
		params    map[string]interface{}
		mockError error
		wantCalls int
		wantErr   bool
	}{
		{
			name: "execute multiple nodes in parallel",
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
			name: "error in one node fails all",
			params: map[string]interface{}{
				"nodes": []neta.Definition{
					{Type: "test1", Name: "node1"},
					{Type: "test2", Name: "node2"},
					{Type: "test3", Name: "node3"},
				},
			},
			mockError: errors.New("execution failed"),
			wantCalls: 3,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mu sync.Mutex
			callCount := 0

			mock := &mockExecutor{
				executeFunc: func(ctx context.Context, def neta.Definition) (neta.Result, error) {
					mu.Lock()
					callCount++
					mu.Unlock()

					// Small delay to simulate work
					time.Sleep(10 * time.Millisecond)

					if tt.mockError != nil {
						return neta.Result{}, tt.mockError
					}
					return neta.Result{Output: def.Name}, nil
				},
			}

			par := NewParallel(mock)
			result, err := par.Execute(context.Background(), tt.params)

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

			if callCount != tt.wantCalls {
				t.Errorf("Execute() call count = %d, want %d", callCount, tt.wantCalls)
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

func TestParallel_ExecutesConcurrently(t *testing.T) {
	var mu sync.Mutex
	executing := 0
	maxConcurrent := 0

	mock := &mockExecutor{
		executeFunc: func(ctx context.Context, def neta.Definition) (neta.Result, error) {
			mu.Lock()
			executing++
			if executing > maxConcurrent {
				maxConcurrent = executing
			}
			mu.Unlock()

			// Simulate work
			time.Sleep(50 * time.Millisecond)

			mu.Lock()
			executing--
			mu.Unlock()

			return neta.Result{Output: def.Name}, nil
		},
	}

	par := NewParallel(mock)
	params := map[string]interface{}{
		"nodes": []neta.Definition{
			{Type: "test1", Name: "node1"},
			{Type: "test2", Name: "node2"},
			{Type: "test3", Name: "node3"},
		},
	}

	_, err := par.Execute(context.Background(), params)
	if err != nil {
		t.Errorf("Execute() unexpected error: %v", err)
	}

	// Verify that multiple nodes were executing concurrently
	if maxConcurrent < 2 {
		t.Errorf("Execute() max concurrent = %d, want at least 2 (indicating parallel execution)", maxConcurrent)
	}
}
