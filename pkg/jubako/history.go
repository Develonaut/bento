package jubako

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// History manages execution history.
type History struct {
	historyDir string
}

// NewHistory creates a new history manager.
func NewHistory(dir string) (*History, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	return &History{historyDir: dir}, nil
}

// Record saves an execution record.
func (h *History) Record(rec ExecutionRecord) error {
	if rec.ID == "" {
		rec.ID = uuid.New().String()
	}

	path := h.recordPath(rec.ID)
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// Get retrieves an execution record by ID.
func (h *History) Get(id string) (ExecutionRecord, error) {
	return h.loadRecord(h.recordPath(id))
}

// List returns execution history with optional filtering.
func (h *History) List(filter HistoryFilter) ([]ExecutionRecord, error) {
	files, err := h.listFiles()
	if err != nil {
		return nil, err
	}

	return h.filterRecords(files, filter), nil
}

// filterRecords loads and filters records from files.
func (h *History) filterRecords(files []string, filter HistoryFilter) []ExecutionRecord {
	records := []ExecutionRecord{}
	for _, file := range files {
		if filter.Limit > 0 && len(records) >= filter.Limit {
			break
		}
		if rec, err := h.loadRecord(file); err == nil && matchesFilter(rec, filter) {
			records = append(records, rec)
		}
	}
	return records
}

// Clear removes all history records.
func (h *History) Clear() error {
	return os.RemoveAll(h.historyDir)
}

// recordPath returns the file path for a record.
func (h *History) recordPath(id string) string {
	return filepath.Join(h.historyDir, id+".json")
}

// listFiles returns all history files sorted by modification time.
func (h *History) listFiles() ([]string, error) {
	pattern := filepath.Join(h.historyDir, "*.json")
	return filepath.Glob(pattern)
}

// loadRecord loads a record from a file.
func (h *History) loadRecord(path string) (ExecutionRecord, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ExecutionRecord{}, err
	}

	var rec ExecutionRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return ExecutionRecord{}, err
	}

	return rec, nil
}

// matchesFilter checks if a record matches the filter.
func matchesFilter(rec ExecutionRecord, filter HistoryFilter) bool {
	if filter.Workflow != "" && rec.Workflow != filter.Workflow {
		return false
	}

	if filter.SuccessOnly && !rec.Success {
		return false
	}

	return true
}
