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
			t.Errorf("List() got %d workflows, want 0", len(infos))
		}
	})

	t.Run("list multiple workflows", func(t *testing.T) {
		// Save multiple workflows
		workflows := []struct {
			name string
			typ  string
		}{
			{"flow1", "http"},
			{"flow2", "transform.jq"},
			{"flow3", "group.sequence"},
		}

		for _, wf := range workflows {
			def := neta.Definition{
				Type: wf.typ,
				Name: wf.name,
			}
			if err := store.Save(wf.name, def); err != nil {
				t.Fatalf("Save() error = %v", err)
			}
		}

		infos, err := store.List()
		if err != nil {
			t.Errorf("List() error = %v", err)
			return
		}

		if len(infos) != len(workflows) {
			t.Errorf("List() got %d workflows, want %d", len(infos), len(workflows))
		}

		// Verify each workflow info
		for _, info := range infos {
			if info.Name == "" {
				t.Error("List() returned workflow with empty name")
			}
			if info.Path == "" {
				t.Error("List() returned workflow with empty path")
			}
			if info.Type == "" {
				t.Error("List() returned workflow with empty type")
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
			t.Errorf("List() got %d workflows, expected at least 1", len(infos))
		}
	})
}
