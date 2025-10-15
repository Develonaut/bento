package jubako

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewHistory(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("creates history with valid directory", func(t *testing.T) {
		histDir := filepath.Join(tmpDir, "history")
		hist, err := NewHistory(histDir)
		if err != nil {
			t.Errorf("NewHistory() error = %v", err)
			return
		}
		if hist == nil {
			t.Error("NewHistory() returned nil history")
		}

		// Verify directory was created
		if _, err := os.Stat(histDir); os.IsNotExist(err) {
			t.Error("NewHistory() did not create directory")
		}
	})

	t.Run("creates nested directories", func(t *testing.T) {
		histDir := filepath.Join(tmpDir, "nested", "history")
		_, err := NewHistory(histDir)
		if err != nil {
			t.Errorf("NewHistory() error = %v", err)
		}

		// Verify nested directory was created
		if _, err := os.Stat(histDir); os.IsNotExist(err) {
			t.Error("NewHistory() did not create nested directory")
		}
	})
}

func TestHistory_RecordAndGet(t *testing.T) {
	tmpDir := t.TempDir()
	hist, err := NewHistory(tmpDir)
	if err != nil {
		t.Fatalf("NewHistory() error = %v", err)
	}

	rec := ExecutionRecord{
		ID:        "test-id-123",
		Bento:     "test-bento",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(5 * time.Second),
		Success:   true,
		Result:    map[string]interface{}{"status": "ok"},
	}

	t.Run("record and get execution", func(t *testing.T) {
		err := hist.Record(rec)
		if err != nil {
			t.Errorf("Record() error = %v", err)
			return
		}

		retrieved, err := hist.Get(rec.ID)
		if err != nil {
			t.Errorf("Get() error = %v", err)
			return
		}

		if retrieved.ID != rec.ID {
			t.Errorf("Get() got ID = %v, want %v", retrieved.ID, rec.ID)
		}
		if retrieved.Bento != rec.Bento {
			t.Errorf("Get() got Bento = %v, want %v", retrieved.Bento, rec.Bento)
		}
		if retrieved.Success != rec.Success {
			t.Errorf("Get() got Success = %v, want %v", retrieved.Success, rec.Success)
		}
	})

	t.Run("record without ID generates UUID", func(t *testing.T) {
		rec := ExecutionRecord{
			Bento:   "test-bento-2",
			Success: true,
		}

		err := hist.Record(rec)
		if err != nil {
			t.Errorf("Record() error = %v", err)
			return
		}

		// Verify a file was created (UUID was generated)
		files, err := hist.listFiles()
		if err != nil {
			t.Errorf("listFiles() error = %v", err)
			return
		}
		if len(files) < 2 {
			t.Error("Record() did not generate UUID for record")
		}
	})

	t.Run("get non-existent record", func(t *testing.T) {
		_, err := hist.Get("nonexistent-id")
		if err == nil {
			t.Error("Get() expected error for non-existent record")
		}
	})
}
