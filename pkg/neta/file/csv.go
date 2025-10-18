package file

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"

	"bento/pkg/neta"
)

// CSV reads CSV files and returns rows as maps.
type CSV struct{}

// NewCSV creates a new CSV reader node.
func NewCSV() *CSV {
	return &CSV{}
}

// Execute reads a CSV file and returns rows.
func (c *CSV) Execute(ctx context.Context, params map[string]interface{}) (neta.Result, error) {
	path := neta.GetStringParam(params, "path", "")
	if path == "" {
		return neta.Result{}, fmt.Errorf("path parameter required")
	}

	// Get has_header flag (default: true)
	hasHeader := true
	if val, ok := params["has_header"].(bool); ok {
		hasHeader = val
	}

	// Read and parse CSV
	rows, err := readCSV(path, hasHeader)
	if err != nil {
		return neta.Result{}, err
	}

	return neta.Result{Output: rows}, nil
}

// readCSV reads and parses a CSV file.
func readCSV(path string, hasHeader bool) ([]map[string]interface{}, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	// Empty file
	if len(records) == 0 {
		return []map[string]interface{}{}, nil
	}

	// Extract headers
	var headers []string
	startRow := 0

	if hasHeader {
		headers = records[0]
		startRow = 1
	} else {
		// Generate column names: col0, col1, etc.
		headers = make([]string, len(records[0]))
		for i := range headers {
			headers[i] = fmt.Sprintf("col%d", i)
		}
	}

	// Convert rows to maps
	rows := make([]map[string]interface{}, 0, len(records)-startRow)
	for _, record := range records[startRow:] {
		row := make(map[string]interface{})
		for i, value := range record {
			if i < len(headers) {
				row[headers[i]] = value
			}
		}
		rows = append(rows, row)
	}

	return rows, nil
}
