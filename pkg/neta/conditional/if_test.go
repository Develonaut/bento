package conditional

import (
	"context"
	"errors"
	"testing"

	"bento/pkg/neta"
)

// mockExecutor implements Executor for testing.
type mockExecutor struct {
	executeFunc func(ctx context.Context, def neta.Definition) (neta.Result, error)
}

func (m *mockExecutor) Execute(ctx context.Context, def neta.Definition) (neta.Result, error) {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, def)
	}
	return neta.Result{Output: "executed"}, nil
}

func TestIf_Execute(t *testing.T) {
	tests := []struct {
		name       string
		params     map[string]interface{}
		mockResult neta.Result
		mockError  error
		wantOutput interface{}
		wantErr    bool
	}{
		{
			name: "condition true executes then branch",
			params: map[string]interface{}{
				"condition": true,
				"then": neta.Definition{
					Type: "test",
					Name: "then-branch",
				},
			},
			mockResult: neta.Result{Output: "then executed"},
			wantOutput: "then executed",
			wantErr:    false,
		},
		{
			name: "condition false executes else branch",
			params: map[string]interface{}{
				"condition": false,
				"else": neta.Definition{
					Type: "test",
					Name: "else-branch",
				},
			},
			mockResult: neta.Result{Output: "else executed"},
			wantOutput: "else executed",
			wantErr:    false,
		},
		{
			name: "missing condition defaults to false",
			params: map[string]interface{}{
				"else": neta.Definition{
					Type: "test",
					Name: "else-branch",
				},
			},
			mockResult: neta.Result{Output: "else executed"},
			wantOutput: "else executed",
			wantErr:    false,
		},
		{
			name: "missing branch returns empty result",
			params: map[string]interface{}{
				"condition": true,
			},
			mockResult: neta.Result{},
			wantOutput: nil,
			wantErr:    false,
		},
		{
			name: "executor error propagates",
			params: map[string]interface{}{
				"condition": true,
				"then": neta.Definition{
					Type: "test",
					Name: "then-branch",
				},
			},
			mockError: errors.New("execution failed"),
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
					return tt.mockResult, nil
				},
			}

			ifNode := NewIf(mock)
			result, err := ifNode.Execute(context.Background(), tt.params)

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

			if tt.wantOutput != nil && result.Output != tt.wantOutput {
				t.Errorf("Execute() output = %v, want %v", result.Output, tt.wantOutput)
			}
		})
	}
}

func TestGetBoolParam(t *testing.T) {
	tests := []struct {
		name   string
		params map[string]interface{}
		key    string
		def    bool
		want   bool
	}{
		{
			name:   "existing bool true",
			params: map[string]interface{}{"key": true},
			key:    "key",
			def:    false,
			want:   true,
		},
		{
			name:   "existing bool false",
			params: map[string]interface{}{"key": false},
			key:    "key",
			def:    true,
			want:   false,
		},
		{
			name:   "missing key returns default",
			params: map[string]interface{}{},
			key:    "missing",
			def:    true,
			want:   true,
		},
		{
			name:   "wrong type returns default",
			params: map[string]interface{}{"key": "not a bool"},
			key:    "key",
			def:    true,
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getBoolParam(tt.params, tt.key, tt.def)
			if got != tt.want {
				t.Errorf("getBoolParam() = %v, want %v", got, tt.want)
			}
		})
	}
}
