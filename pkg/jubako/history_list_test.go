package jubako

import (
	"testing"
)

func TestHistory_List(t *testing.T) {
	tmpDir := t.TempDir()
	hist, err := NewHistory(tmpDir)
	if err != nil {
		t.Fatalf("NewHistory() error = %v", err)
	}

	// Create test records
	records := []ExecutionRecord{
		{
			ID:      "rec1",
			Bento:   "bento1",
			Success: true,
		},
		{
			ID:      "rec2",
			Bento:   "bento1",
			Success: false,
			Error:   "test error",
		},
		{
			ID:      "rec3",
			Bento:   "bento2",
			Success: true,
		},
	}

	for _, rec := range records {
		if err := hist.Record(rec); err != nil {
			t.Fatalf("Record() error = %v", err)
		}
	}

	t.Run("list all records", func(t *testing.T) {
		filter := HistoryFilter{}
		list, err := hist.List(filter)
		if err != nil {
			t.Errorf("List() error = %v", err)
			return
		}

		if len(list) != len(records) {
			t.Errorf("List() got %d records, want %d", len(list), len(records))
		}
	})

	t.Run("list with bento filter", func(t *testing.T) {
		filter := HistoryFilter{
			Bento: "bento1",
		}
		list, err := hist.List(filter)
		if err != nil {
			t.Errorf("List() error = %v", err)
			return
		}

		expectedCount := 2
		if len(list) != expectedCount {
			t.Errorf("List() got %d records, want %d", len(list), expectedCount)
		}

		// Verify all records match the bento
		for _, rec := range list {
			if rec.Bento != "bento1" {
				t.Errorf("List() returned record with bento = %v, want bento1", rec.Bento)
			}
		}
	})

	t.Run("list success only", func(t *testing.T) {
		filter := HistoryFilter{
			SuccessOnly: true,
		}
		list, err := hist.List(filter)
		if err != nil {
			t.Errorf("List() error = %v", err)
			return
		}

		expectedCount := 2
		if len(list) != expectedCount {
			t.Errorf("List() got %d records, want %d", len(list), expectedCount)
		}

		// Verify all records are successful
		for _, rec := range list {
			if !rec.Success {
				t.Error("List() returned failed record with SuccessOnly filter")
			}
		}
	})

	t.Run("list with limit", func(t *testing.T) {
		filter := HistoryFilter{
			Limit: 2,
		}
		list, err := hist.List(filter)
		if err != nil {
			t.Errorf("List() error = %v", err)
			return
		}

		if len(list) != 2 {
			t.Errorf("List() got %d records, want 2", len(list))
		}
	})

	t.Run("list with combined filters", func(t *testing.T) {
		filter := HistoryFilter{
			Bento:       "bento1",
			SuccessOnly: true,
			Limit:       1,
		}
		list, err := hist.List(filter)
		if err != nil {
			t.Errorf("List() error = %v", err)
			return
		}

		if len(list) != 1 {
			t.Errorf("List() got %d records, want 1", len(list))
		}

		if len(list) > 0 {
			if list[0].Bento != "bento1" {
				t.Errorf("List() returned record with bento = %v, want bento1", list[0].Bento)
			}
			if !list[0].Success {
				t.Error("List() returned failed record with SuccessOnly filter")
			}
		}
	})
}

func TestMatchesFilter(t *testing.T) {
	tests := []struct {
		name   string
		rec    ExecutionRecord
		filter HistoryFilter
		want   bool
	}{
		{
			name: "empty filter matches all",
			rec: ExecutionRecord{
				Bento:   "test",
				Success: true,
			},
			filter: HistoryFilter{},
			want:   true,
		},
		{
			name: "bento filter matches",
			rec: ExecutionRecord{
				Bento:   "test",
				Success: true,
			},
			filter: HistoryFilter{
				Bento: "test",
			},
			want: true,
		},
		{
			name: "bento filter does not match",
			rec: ExecutionRecord{
				Bento:   "test",
				Success: true,
			},
			filter: HistoryFilter{
				Bento: "other",
			},
			want: false,
		},
		{
			name: "success only matches successful",
			rec: ExecutionRecord{
				Bento:   "test",
				Success: true,
			},
			filter: HistoryFilter{
				SuccessOnly: true,
			},
			want: true,
		},
		{
			name: "success only does not match failed",
			rec: ExecutionRecord{
				Bento:   "test",
				Success: false,
			},
			filter: HistoryFilter{
				SuccessOnly: true,
			},
			want: false,
		},
		{
			name: "combined filter matches",
			rec: ExecutionRecord{
				Bento:   "test",
				Success: true,
			},
			filter: HistoryFilter{
				Bento:       "test",
				SuccessOnly: true,
			},
			want: true,
		},
		{
			name: "combined filter does not match",
			rec: ExecutionRecord{
				Bento:   "test",
				Success: false,
			},
			filter: HistoryFilter{
				Bento:       "test",
				SuccessOnly: true,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchesFilter(tt.rec, tt.filter)
			if got != tt.want {
				t.Errorf("matchesFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}
