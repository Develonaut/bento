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
		workDir := filepath.Join(tmpDir, "bentos")
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
		workDir := filepath.Join(tmpDir, "nested", "path", "bentos")
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
		Version: "1.0",
		Type:    "http",
		Name:    "Test Bento",
		Parameters: map[string]interface{}{
			"url":    "https://example.com",
			"method": "GET",
		},
	}

	t.Run("save and load bento", func(t *testing.T) {
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

	t.Run("save with .bento.json extension", func(t *testing.T) {
		err := store.Save("test.bento.json", def)
		if err != nil {
			t.Errorf("Save() error = %v", err)
		}

		// Verify file was created
		path := filepath.Join(tmpDir, "test.bento.json")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Error("Save() did not create file")
		}
	})

	t.Run("load non-existent bento", func(t *testing.T) {
		_, err := store.Load("nonexistent")
		if err == nil {
			t.Error("Load() expected error for non-existent bento")
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
		Version: "1.0",
		Type:    "http",
		Name:    "Test",
	}

	// Save a bento first
	if err := store.Save("test", def); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	t.Run("delete existing bento", func(t *testing.T) {
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

	t.Run("delete non-existent bento", func(t *testing.T) {
		err := store.Delete("nonexistent")
		if err == nil {
			t.Error("Delete() expected error for non-existent bento")
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
			want:  "test.bento.json",
		},
		{
			name:  "name with extension",
			input: "test.bento.json",
			want:  "test.bento.json",
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
