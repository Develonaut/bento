package neta

import "testing"

func TestDefinition_IsGroup(t *testing.T) {
	tests := []struct {
		name string
		def  Definition
		want bool
	}{
		{
			name: "leaf node has no children",
			def: Definition{
				Type:  "http",
				Name:  "Get User",
				Nodes: nil,
			},
			want: false,
		},
		{
			name: "empty nodes slice is not a group",
			def: Definition{
				Type:  "http",
				Name:  "Get User",
				Nodes: []Definition{},
			},
			want: false,
		},
		{
			name: "group node has children",
			def: Definition{
				Type: "group",
				Name: "Workflow",
				Nodes: []Definition{
					{Type: "http", Name: "Step 1"},
				},
			},
			want: true,
		},
		{
			name: "group with multiple children",
			def: Definition{
				Type: "group",
				Name: "Complex Workflow",
				Nodes: []Definition{
					{Type: "http", Name: "Step 1"},
					{Type: "transform", Name: "Step 2"},
					{Type: "http", Name: "Step 3"},
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.def.IsGroup(); got != tt.want {
				t.Errorf("Definition.IsGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}
