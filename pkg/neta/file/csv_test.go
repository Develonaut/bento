package file

import (
	"context"
	"os"
	"testing"
)

func TestCSV_Execute(t *testing.T) {
	tests := []struct {
		name       string
		csvContent string
		params     map[string]interface{}
		wantErr    bool
		wantRows   int
		verify     func(t *testing.T, output interface{})
	}{
		{
			name: "read CSV with headers",
			csvContent: `name,age,city
Alice,30,NYC
Bob,25,LA
Charlie,35,SF`,
			params: map[string]interface{}{
				"has_header": true,
			},
			wantErr:  false,
			wantRows: 3,
			verify: func(t *testing.T, output interface{}) {
				rows, ok := output.([]map[string]interface{})
				if !ok {
					t.Fatal("Output is not []map[string]interface{}")
				}
				if len(rows) != 3 {
					t.Fatalf("Expected 3 rows, got %d", len(rows))
				}
				// Check first row
				if rows[0]["name"] != "Alice" {
					t.Errorf("First row name = %v, want Alice", rows[0]["name"])
				}
				if rows[0]["age"] != "30" {
					t.Errorf("First row age = %v, want 30", rows[0]["age"])
				}
				if rows[0]["city"] != "NYC" {
					t.Errorf("First row city = %v, want NYC", rows[0]["city"])
				}
			},
		},
		{
			name: "read CSV without headers",
			csvContent: `Alice,30,NYC
Bob,25,LA`,
			params: map[string]interface{}{
				"has_header": false,
			},
			wantErr:  false,
			wantRows: 2,
			verify: func(t *testing.T, output interface{}) {
				rows, ok := output.([]map[string]interface{})
				if !ok {
					t.Fatal("Output is not []map[string]interface{}")
				}
				if len(rows) != 2 {
					t.Fatalf("Expected 2 rows, got %d", len(rows))
				}
				// Check columns are named col0, col1, col2
				if rows[0]["col0"] != "Alice" {
					t.Errorf("First row col0 = %v, want Alice", rows[0]["col0"])
				}
				if rows[0]["col1"] != "30" {
					t.Errorf("First row col1 = %v, want 30", rows[0]["col1"])
				}
			},
		},
		{
			name:       "empty CSV file",
			csvContent: "",
			params: map[string]interface{}{
				"has_header": true,
			},
			wantErr:  false,
			wantRows: 0,
			verify: func(t *testing.T, output interface{}) {
				rows, ok := output.([]map[string]interface{})
				if !ok {
					t.Fatal("Output is not []map[string]interface{}")
				}
				if len(rows) != 0 {
					t.Fatalf("Expected 0 rows, got %d", len(rows))
				}
			},
		},
		{
			name: "CSV with only headers",
			csvContent: `name,age,city
`,
			params: map[string]interface{}{
				"has_header": true,
			},
			wantErr:  false,
			wantRows: 0,
			verify: func(t *testing.T, output interface{}) {
				rows, ok := output.([]map[string]interface{})
				if !ok {
					t.Fatal("Output is not []map[string]interface{}")
				}
				if len(rows) != 0 {
					t.Fatalf("Expected 0 rows for headers-only, got %d", len(rows))
				}
			},
		},
		{
			name: "default has_header is true",
			csvContent: `name,age
Alice,30`,
			params:   map[string]interface{}{},
			wantErr:  false,
			wantRows: 1,
			verify: func(t *testing.T, output interface{}) {
				rows, ok := output.([]map[string]interface{})
				if !ok {
					t.Fatal("Output is not []map[string]interface{}")
				}
				if len(rows) != 1 {
					t.Fatalf("Expected 1 row, got %d", len(rows))
				}
				if rows[0]["name"] != "Alice" {
					t.Error("Default should treat first row as header")
				}
			},
		},
		{
			name:       "missing path parameter",
			csvContent: "name,age\n",
			params:     map[string]interface{}{},
			wantErr:    true,
		},
		{
			name:       "file does not exist",
			csvContent: "",
			params: map[string]interface{}{
				"path": "/nonexistent/file.csv",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp CSV file if not testing missing file
			var csvPath string
			if tt.name != "missing path parameter" && tt.name != "file does not exist" {
				tmpFile, err := os.CreateTemp(t.TempDir(), "test-*.csv")
				if err != nil {
					t.Fatalf("Failed to create temp file: %v", err)
				}
				defer tmpFile.Close()

				if tt.csvContent != "" {
					if _, err := tmpFile.WriteString(tt.csvContent); err != nil {
						t.Fatalf("Failed to write CSV content: %v", err)
					}
				}

				csvPath = tmpFile.Name()
				tt.params["path"] = csvPath
			}

			csv := NewCSV()
			result, err := csv.Execute(context.Background(), tt.params)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Execute() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Execute() unexpected error: %v", err)
				return
			}

			if result.Output == nil {
				t.Error("Execute() output is nil")
				return
			}

			// Verify output
			if tt.verify != nil {
				tt.verify(t, result.Output)
			}
		})
	}
}
