package jubako

import (
	"os"
	"testing"
)

func TestHistory_Clear(t *testing.T) {
	tmpDir := t.TempDir()
	hist, err := NewHistory(tmpDir)
	if err != nil {
		t.Fatalf("NewHistory() error = %v", err)
	}

	// Create some records
	for i := 0; i < 3; i++ {
		rec := ExecutionRecord{
			Bento:   "test",
			Success: true,
		}
		if err := hist.Record(rec); err != nil {
			t.Fatalf("Record() error = %v", err)
		}
	}

	t.Run("clear all records", func(t *testing.T) {
		err := hist.Clear()
		if err != nil {
			t.Errorf("Clear() error = %v", err)
			return
		}

		// Verify directory no longer exists
		if _, err := os.Stat(tmpDir); !os.IsNotExist(err) {
			// Directory might exist but should be empty
			files, _ := hist.listFiles()
			if len(files) > 0 {
				t.Error("Clear() did not remove all records")
			}
		}
	})
}
