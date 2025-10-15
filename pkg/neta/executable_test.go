package neta

import (
	"context"
	"errors"
	"testing"
)

// mockExecutable is a test implementation of Executable
type mockExecutable struct {
	output interface{}
	err    error
}

func (m *mockExecutable) Execute(ctx context.Context, params map[string]interface{}) (Result, error) {
	if m.err != nil {
		return Result{}, m.err
	}
	return Result{Output: m.output}, nil
}

func TestExecutable_Interface(t *testing.T) {
	tests := []struct {
		name    string
		exec    Executable
		params  map[string]interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "successful execution",
			exec: &mockExecutable{
				output: "success",
				err:    nil,
			},
			params:  map[string]interface{}{"key": "value"},
			want:    "success",
			wantErr: false,
		},
		{
			name: "execution with error",
			exec: &mockExecutable{
				output: nil,
				err:    errors.New("execution failed"),
			},
			params:  map[string]interface{}{},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			result, err := tt.exec.Execute(ctx, tt.params)

			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if result.Output != tt.want {
				t.Errorf("Execute() Output = %v, want %v", result.Output, tt.want)
			}
		})
	}
}
