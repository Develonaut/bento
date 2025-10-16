package examples

import (
	"testing"
)

func TestGetAll(t *testing.T) {
	categories := GetAll()

	if len(categories) == 0 {
		t.Error("should have categories")
	}

	totalExamples := 0
	for _, cat := range categories {
		totalExamples += len(cat.Examples)
	}

	if totalExamples < 7 {
		t.Errorf("should have at least 7 examples, got %d", totalExamples)
	}

	// Verify each category has a name
	for _, cat := range categories {
		if cat.Name == "" {
			t.Error("category missing name")
		}
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		id       string
		wantName string
	}{
		{"http-get", "Simple GET Request"},
		{"http-post", "POST with JSON Body"},
		{"transform-jq", "Extract User IDs"},
		{"conditional-if", "Check Success Status"},
		{"loop-for", "Process Multiple Users"},
		{"sequence", "Multi-Step Workflow"},
		{"complete-api-workflow", "Complete API Workflow"},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			ex, err := Get(tt.id)
			if err != nil {
				t.Fatalf("should find %s example: %v", tt.id, err)
			}

			if ex.Name != tt.wantName {
				t.Errorf("wrong example name: got %s, want %s", ex.Name, tt.wantName)
			}
		})
	}
}

func TestGetNotFound(t *testing.T) {
	_, err := Get("nonexistent")
	if err == nil {
		t.Error("should return error for nonexistent example")
	}
}

func TestParse(t *testing.T) {
	ex, err := Get("http-get")
	if err != nil {
		t.Fatal(err)
	}

	def, err := Parse(*ex)
	if err != nil {
		t.Fatalf("should parse example: %v", err)
	}

	if def.Version != "1.0" {
		t.Error("example should have version 1.0")
	}

	if def.Type != "http" {
		t.Errorf("wrong type: %s", def.Type)
	}
}

func TestParseAllExamples(t *testing.T) {
	examples := List()

	for _, ex := range examples {
		t.Run(ex.ID, func(t *testing.T) {
			def, err := Parse(ex)
			if err != nil {
				t.Fatalf("failed to parse example %s: %v", ex.ID, err)
			}

			if def.Version == "" {
				t.Error("example missing version")
			}

			if def.Type == "" {
				t.Error("example missing type")
			}
		})
	}
}

func TestList(t *testing.T) {
	examples := List()

	if len(examples) < 7 {
		t.Errorf("should have at least 7 examples, got %d", len(examples))
	}

	// Check each example has required fields
	for _, ex := range examples {
		if ex.ID == "" {
			t.Error("example missing ID")
		}
		if ex.Name == "" {
			t.Error("example missing name")
		}
		if ex.Description == "" {
			t.Error("example missing description")
		}
		if ex.Category == "" {
			t.Error("example missing category")
		}
		if ex.Content == "" {
			t.Error("example missing content")
		}
	}
}

func TestExampleCategories(t *testing.T) {
	categories := GetAll()

	expectedCategories := map[string]bool{
		"HTTP Requests":       false,
		"Data Transformation": false,
		"Control Flow":        false,
		"Complete Workflows":  false,
	}

	for _, cat := range categories {
		if _, ok := expectedCategories[cat.Name]; ok {
			expectedCategories[cat.Name] = true
		}
	}

	for name, found := range expectedCategories {
		if !found {
			t.Errorf("missing expected category: %s", name)
		}
	}
}
