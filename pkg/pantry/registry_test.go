package pantry

import (
	"context"
	"sync"
	"testing"

	"bento/pkg/neta"
)

// mockExecutable is a test implementation
type mockExecutable struct {
	name string
}

func (m *mockExecutable) Execute(ctx context.Context, params map[string]interface{}) (neta.Result, error) {
	return neta.Result{Output: m.name}, nil
}

func TestPantry_Register(t *testing.T) {
	tests := []struct {
		name      string
		nodeType  string
		exec      neta.Executable
		wantErr   bool
		duplicate bool
	}{
		{
			name:     "successful registration",
			nodeType: "http",
			exec:     &mockExecutable{name: "http"},
			wantErr:  false,
		},
		{
			name:      "duplicate registration",
			nodeType:  "http",
			exec:      &mockExecutable{name: "http"},
			wantErr:   true,
			duplicate: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			if tt.duplicate {
				_ = p.Register(tt.nodeType, tt.exec)
			}

			err := p.Register(tt.nodeType, tt.exec)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPantry_Get(t *testing.T) {
	p := New()
	exec := &mockExecutable{name: "test"}
	_ = p.Register("test", exec)

	tests := []struct {
		name     string
		nodeType string
		wantErr  bool
	}{
		{
			name:     "get existing node",
			nodeType: "test",
			wantErr:  false,
		},
		{
			name:     "get non-existent node",
			nodeType: "unknown",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := p.Get(tt.nodeType)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPantry_List(t *testing.T) {
	p := New()
	_ = p.Register("http", &mockExecutable{name: "http"})
	_ = p.Register("transform", &mockExecutable{name: "transform"})

	types := p.List()
	if len(types) != 2 {
		t.Errorf("List() count = %d, want 2", len(types))
	}
}

func TestPantry_ThreadSafety(t *testing.T) {
	p := New()
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			nodeType := string(rune('a' + id))
			p.Register(nodeType, &mockExecutable{name: nodeType})
		}(i)
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p.List()
		}()
	}

	wg.Wait()

	types := p.List()
	if len(types) != 10 {
		t.Errorf("Thread safety test: expected 10 types, got %d", len(types))
	}
}
