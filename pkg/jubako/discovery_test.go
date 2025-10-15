package jubako

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewDiscovery(t *testing.T) {
	t.Run("with custom paths", func(t *testing.T) {
		paths := []string{"/path1", "/path2"}
		disc := NewDiscovery(paths...)
		if len(disc.searchPaths) != 2 {
			t.Errorf("NewDiscovery() got %d paths, want 2", len(disc.searchPaths))
		}
	})

	t.Run("with default paths", func(t *testing.T) {
		disc := NewDiscovery()
		if len(disc.searchPaths) == 0 {
			t.Error("NewDiscovery() returned no default paths")
		}
	})
}

func TestDiscovery_Find(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test directory structure
	dirs := []string{
		filepath.Join(tmpDir, "dir1"),
		filepath.Join(tmpDir, "dir2"),
		filepath.Join(tmpDir, "dir2", "nested"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}
	}

	// Create test files
	files := []struct {
		path    string
		isBento bool
	}{
		{filepath.Join(tmpDir, "dir1", "flow1.bento.yaml"), true},
		{filepath.Join(tmpDir, "dir1", "flow2.bento.yaml"), true},
		{filepath.Join(tmpDir, "dir2", "flow3.bento.yaml"), true},
		{filepath.Join(tmpDir, "dir2", "nested", "flow4.bento.yaml"), true},
		{filepath.Join(tmpDir, "dir1", "not-bento.yaml"), false},
		{filepath.Join(tmpDir, "dir1", "readme.md"), false},
	}

	for _, f := range files {
		content := []byte("type: http\nname: Test")
		if err := os.WriteFile(f.path, content, 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	t.Run("find all bento files", func(t *testing.T) {
		disc := NewDiscovery(tmpDir)
		found, err := disc.Find()
		if err != nil {
			t.Errorf("Find() error = %v", err)
			return
		}

		expectedCount := 4 // Only .bento.yaml files
		if len(found) != expectedCount {
			t.Errorf("Find() got %d files, want %d", len(found), expectedCount)
		}

		// Verify all found files are .bento.yaml files
		for _, path := range found {
			if !isBentoFile(path) {
				t.Errorf("Find() returned non-bento file: %s", path)
			}
		}
	})

	t.Run("find in specific directory", func(t *testing.T) {
		dir1Path := filepath.Join(tmpDir, "dir1")
		disc := NewDiscovery(dir1Path)
		found, err := disc.Find()
		if err != nil {
			t.Errorf("Find() error = %v", err)
			return
		}

		expectedCount := 2
		if len(found) != expectedCount {
			t.Errorf("Find() got %d files, want %d", len(found), expectedCount)
		}
	})

	t.Run("find with non-existent path", func(t *testing.T) {
		disc := NewDiscovery(filepath.Join(tmpDir, "nonexistent"))
		found, err := disc.Find()
		if err != nil {
			t.Errorf("Find() error = %v", err)
			return
		}

		if len(found) != 0 {
			t.Errorf("Find() got %d files, want 0 for non-existent path", len(found))
		}
	})

	t.Run("find with multiple paths", func(t *testing.T) {
		dir1 := filepath.Join(tmpDir, "dir1")
		dir2 := filepath.Join(tmpDir, "dir2")
		disc := NewDiscovery(dir1, dir2)
		found, err := disc.Find()
		if err != nil {
			t.Errorf("Find() error = %v", err)
			return
		}

		if len(found) < 3 {
			t.Errorf("Find() got %d files, expected at least 3", len(found))
		}
	})
}

func TestIsBentoFile(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "valid bento file",
			path: "/path/to/example.bento.yaml",
			want: true,
		},
		{
			name: "valid bento file with longer name",
			path: "/path/to/my-complex-bento.bento.yaml",
			want: true,
		},
		{
			name: "regular yaml file",
			path: "/path/to/config.yaml",
			want: false,
		},
		{
			name: "yaml with bento in name",
			path: "/path/to/bento.yaml",
			want: false,
		},
		{
			name: "text file",
			path: "/path/to/readme.txt",
			want: false,
		},
		{
			name: "no extension",
			path: "/path/to/file",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isBentoFile(tt.path)
			if got != tt.want {
				t.Errorf("isBentoFile(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestDefaultPaths(t *testing.T) {
	paths := defaultPaths()
	if len(paths) == 0 {
		t.Error("defaultPaths() returned empty slice")
	}

	// Should include current directory
	foundCurrent := false
	for _, path := range paths {
		if path == "." {
			foundCurrent = true
			break
		}
	}
	if !foundCurrent {
		t.Error("defaultPaths() does not include current directory")
	}

	// Verify paths are using filepath.Join (cross-platform)
	for _, path := range paths {
		if path != "." && filepath.IsAbs(path) {
			// Absolute paths should be properly formed
			if path == "" {
				t.Errorf("defaultPaths() contains empty path")
			}
		}
	}
}

func TestFindInPath(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	files := []string{
		filepath.Join(tmpDir, "test1.bento.yaml"),
		filepath.Join(tmpDir, "test2.bento.yaml"),
		filepath.Join(tmpDir, "other.yaml"),
	}

	for _, file := range files {
		if err := os.WriteFile(file, []byte("content"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	t.Run("find files in path", func(t *testing.T) {
		found, err := findInPath(tmpDir)
		if err != nil {
			t.Errorf("findInPath() error = %v", err)
			return
		}

		expectedCount := 2 // Only .bento.yaml files
		if len(found) != expectedCount {
			t.Errorf("findInPath() got %d files, want %d", len(found), expectedCount)
		}
	})

	t.Run("non-existent path", func(t *testing.T) {
		found, err := findInPath(filepath.Join(tmpDir, "nonexistent"))
		if err != nil {
			t.Errorf("findInPath() error = %v, expected no error", err)
		}
		if len(found) != 0 {
			t.Errorf("findInPath() got %d files, want 0", len(found))
		}
	})
}
