package jubako

import (
	"os"
	"path/filepath"
	"testing"

	"bento/pkg/neta"
)

func TestNewStore(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("creates store with valid directory", func(t *testing.T) {
		workDir := filepath.Join(tmpDir, "workflows")
		store, err := NewStore(workDir)
		if err != nil {
			t.Errorf("NewStore() error = %v", err)
			return
		}
		if store == nil {
			t.Error("NewStore() returned nil store")
		}

		// Verify directory was created
		if _, err := os.Stat(workDir); os.IsNotExist(err) {
			t.Error("NewStore() did not create directory")
		}
	})

	t.Run("creates nested directories", func(t *testing.T) {
		workDir := filepath.Join(tmpDir, "nested", "path", "workflows")
		_, err := NewStore(workDir)
		if err != nil {
			t.Errorf("NewStore() error = %v", err)
		}

		// Verify nested directory was created
		if _, err := os.Stat(workDir); os.IsNotExist(err) {
			t.Error("NewStore() did not create nested directory")
		}
	})
}

func TestStore_SaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewStore(tmpDir)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	def := neta.Definition{
		Type: "http",
		Name: "Test Workflow",
		Parameters: map[string]interface{}{
			"url":    "https://example.com",
			"method": "GET",
		},
	}

	t.Run("save and load workflow", func(t *testing.T) {
		err := store.Save("test", def)
		if err != nil {
			t.Errorf("Save() error = %v", err)
			return
		}

		loaded, err := store.Load("test")
		if err != nil {
			t.Errorf("Load() error = %v", err)
			return
		}

		if loaded.Type != def.Type {
			t.Errorf("Load() got type = %v, want %v", loaded.Type, def.Type)
		}
		if loaded.Name != def.Name {
			t.Errorf("Load() got name = %v, want %v", loaded.Name, def.Name)
		}
	})

	t.Run("save with .bento.yaml extension", func(t *testing.T) {
		err := store.Save("test.bento.yaml", def)
		if err != nil {
			t.Errorf("Save() error = %v", err)
		}

		// Verify file was created
		path := filepath.Join(tmpDir, "test.bento.yaml")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Error("Save() did not create file")
		}
	})

	t.Run("load non-existent workflow", func(t *testing.T) {
		_, err := store.Load("nonexistent")
		if err == nil {
			t.Error("Load() expected error for non-existent workflow")
		}
	})
}

func TestStore_Delete(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewStore(tmpDir)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	def := neta.Definition{
		Type: "http",
		Name: "Test",
	}

	// Save a workflow first
	if err := store.Save("test", def); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	t.Run("delete existing workflow", func(t *testing.T) {
		err := store.Delete("test")
		if err != nil {
			t.Errorf("Delete() error = %v", err)
		}

		// Verify file was deleted
		_, err = store.Load("test")
		if err == nil {
			t.Error("Delete() did not remove file")
		}
	})

	t.Run("delete non-existent workflow", func(t *testing.T) {
		err := store.Delete("nonexistent")
		if err == nil {
			t.Error("Delete() expected error for non-existent workflow")
		}
	})
}

func TestStore_PathFor(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewStore(tmpDir)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "name without extension",
			input: "test",
			want:  "test.bento.yaml",
		},
		{
			name:  "name with extension",
			input: "test.bento.yaml",
			want:  "test.bento.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := store.pathFor(tt.input)
			wantPath := filepath.Join(tmpDir, tt.want)
			if got != wantPath {
				t.Errorf("pathFor() = %v, want %v", got, wantPath)
			}
		})
	}
}
