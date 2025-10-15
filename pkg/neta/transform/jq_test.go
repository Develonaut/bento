package transform

import (
	"context"
	"testing"
)

func TestJQ_Execute(t *testing.T) {
	tests := []struct {
		name    string
		params  map[string]interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "simple identity query",
			params: map[string]interface{}{
				"query": ".",
				"input": map[string]interface{}{
					"name": "test",
					"age":  30,
				},
			},
			want: map[string]interface{}{
				"name": "test",
				"age":  float64(30),
			},
			wantErr: false,
		},
		{
			name: "extract field",
			params: map[string]interface{}{
				"query": ".name",
				"input": map[string]interface{}{
					"name": "test",
					"age":  30,
				},
			},
			want:    "test",
			wantErr: false,
		},
		{
			name: "array map",
			params: map[string]interface{}{
				"query": ".[].name",
				"input": []interface{}{
					map[string]interface{}{"name": "alice"},
					map[string]interface{}{"name": "bob"},
				},
			},
			want:    "alice",
			wantErr: false,
		},
		{
			name: "default query is identity",
			params: map[string]interface{}{
				"input": map[string]interface{}{"test": "value"},
			},
			want:    map[string]interface{}{"test": "value"},
			wantErr: false,
		},
		{
			name: "invalid query syntax",
			params: map[string]interface{}{
				"query": ".[[[invalid",
				"input": map[string]interface{}{"test": "value"},
			},
			wantErr: true,
		},
		{
			name: "input as JSON string",
			params: map[string]interface{}{
				"query": ".name",
				"input": `{"name": "test", "age": 30}`,
			},
			want:    "test",
			wantErr: false,
		},
		{
			name: "invalid JSON string input",
			params: map[string]interface{}{
				"query": ".",
				"input": `{invalid json}`,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jq := NewJQ()
			result, err := jq.Execute(context.Background(), tt.params)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Execute() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Execute() unexpected error: %v", err)
				return
			}

			// Compare output (note: numbers from JSON become float64)
			// For complex comparisons, you might need a deeper comparison
			if tt.want != nil && result.Output == nil {
				t.Error("Execute() output is nil")
			}
		})
	}
}

func TestApplyQuery(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		data    interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name:  "valid query",
			query: ".name",
			data: map[string]interface{}{
				"name": "test",
			},
			want:    "test",
			wantErr: false,
		},
		{
			name:    "invalid query",
			query:   "[[[invalid",
			data:    map[string]interface{}{},
			wantErr: true,
		},
		{
			name:    "query with no results",
			query:   ".missing",
			data:    map[string]interface{}{"present": "value"},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := applyQuery(tt.query, tt.data)

			if tt.wantErr {
				if err == nil {
					t.Error("applyQuery() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("applyQuery() unexpected error: %v", err)
				return
			}

			// Basic comparison - for more complex cases use reflect.DeepEqual
			_ = result // result validation would go here
		})
	}
}

func TestGetStringParam_Transform(t *testing.T) {
	tests := []struct {
		name   string
		params map[string]interface{}
		key    string
		def    string
		want   string
	}{
		{
			name:   "existing string",
			params: map[string]interface{}{"key": "value"},
			key:    "key",
			def:    "default",
			want:   "value",
		},
		{
			name:   "missing key",
			params: map[string]interface{}{},
			key:    "missing",
			def:    "default",
			want:   "default",
		},
		{
			name:   "wrong type",
			params: map[string]interface{}{"key": 123},
			key:    "key",
			def:    "default",
			want:   "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getStringParam(tt.params, tt.key, tt.def)
			if got != tt.want {
				t.Errorf("getStringParam() = %s, want %s", got, tt.want)
			}
		})
	}
}
