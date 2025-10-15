package jubako

import (
	"os"
	"path/filepath"
	"testing"

	"bento/pkg/neta"
)

func TestStore_List(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewStore(tmpDir)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	t.Run("list empty store", func(t *testing.T) {
		infos, err := store.List()
		if err != nil {
			t.Errorf("List() error = %v", err)
			return
		}
		if len(infos) != 0 {
			t.Errorf("List() got %d bentos, want 0", len(infos))
		}
	})

	t.Run("list multiple bentos", func(t *testing.T) {
		// Save multiple bentos with valid parameters
		bentos := []struct {
			name string
			def  neta.Definition
		}{
			{"flow1", neta.Definition{
				Version:    "1.0",
				Type:       "http",
				Name:       "flow1",
				Parameters: map[string]interface{}{"url": "https://example.com"},
			}},
			{"flow2", neta.Definition{
				Version:    "1.0",
				Type:       "jq",
				Name:       "flow2",
				Parameters: map[string]interface{}{"query": "."},
			}},
			{"flow3", neta.Definition{
				Version: "1.0",
				Type:    "sequence",
				Name:    "flow3",
			}},
		}

		for _, b := range bentos {
			if err := store.Save(b.name, b.def); err != nil {
				t.Fatalf("Save() error = %v", err)
			}
		}

		infos, err := store.List()
		if err != nil {
			t.Errorf("List() error = %v", err)
			return
		}

		if len(infos) != len(bentos) {
			t.Errorf("List() got %d bentos, want %d", len(infos), len(bentos))
		}

		// Verify each bento info
		for _, info := range infos {
			if info.Name == "" {
				t.Error("List() returned bento with empty name")
			}
			if info.Path == "" {
				t.Error("List() returned bento with empty path")
			}
			if info.Type == "" {
				t.Error("List() returned bento with empty type")
			}
		}
	})

	t.Run("list skips invalid files", func(t *testing.T) {
		// Create an invalid .bento.yaml file
		invalidPath := filepath.Join(tmpDir, "invalid.bento.yaml")
		if err := os.WriteFile(invalidPath, []byte("invalid: yaml: content:"), 0644); err != nil {
			t.Fatalf("Failed to write invalid file: %v", err)
		}

		infos, err := store.List()
		if err != nil {
			t.Errorf("List() error = %v", err)
			return
		}

		// Should skip invalid file and return only valid ones
		// From previous test we had 3 valid files
		if len(infos) < 1 {
			t.Errorf("List() got %d bentos, expected at least 1", len(infos))
		}
	})
}
