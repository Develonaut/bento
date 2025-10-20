package template_test

import (
	"context"
	"os"
	"strings"
	"testing"

	templateneta "github.com/Develonaut/bento/pkg/neta/library/template"
)

// createTestSVG creates a simple test SVG file with ID-based text elements
func createTestSVG(t *testing.T) string {
	t.Helper()

	svgContent := `<?xml version="1.0" encoding="UTF-8"?>
<svg width="200" height="200" xmlns="http://www.w3.org/2000/svg">
  <text id="product_name" x="10" y="50">Default Name</text>
  <text id="product_scale" x="10" y="100">Default Scale</text>
  <text id="studio_name" x="10" y="150">Default Studio</text>
</svg>`

	tmpfile, err := os.CreateTemp("", "test-*.svg")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tmpfile.WriteString(svgContent); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	return tmpfile.Name()
}

// createTestXML creates a simple test XML file
func createTestXML(t *testing.T) string {
	t.Helper()

	xmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<root>
  <item id="name">John Doe</item>
  <item id="email">john@example.com</item>
</root>`

	tmpfile, err := os.CreateTemp("", "test-*.xml")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tmpfile.WriteString(xmlContent); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	return tmpfile.Name()
}

// TestTemplate_ReplaceByID tests XML element replacement by ID attribute
func TestTemplate_ReplaceByID(t *testing.T) {
	ctx := context.Background()

	inputPath := createTestSVG(t)
	defer os.Remove(inputPath)

	outputPath := strings.Replace(inputPath, ".svg", "-output.svg", 1)
	defer os.Remove(outputPath)

	tmplNeta := templateneta.New()

	params := map[string]interface{}{
		"operation": "replace",
		"input":     inputPath,
		"output":    outputPath,
		"replacements": map[string]interface{}{
			"product_name":  "Test Product",
			"product_scale": "32mm - Heroic",
			"studio_name":   "Test Studio",
		},
		"mode": "id",
	}

	result, err := tmplNeta.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	output, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected map[string]interface{} result")
	}

	if output["path"] != outputPath {
		t.Errorf("path = %v, want %v", output["path"], outputPath)
	}

	// Verify output file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("Output file not created")
	}

	// Verify replacements were made
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "Test Product") {
		t.Error("Expected 'Test Product' in output")
	}
	if !strings.Contains(contentStr, "32mm - Heroic") {
		t.Error("Expected '32mm - Heroic' in output")
	}
	if !strings.Contains(contentStr, "Test Studio") {
		t.Error("Expected 'Test Studio' in output")
	}

	// Verify old content was replaced
	if strings.Contains(contentStr, "Default Name") {
		t.Error("Old content 'Default Name' should be replaced")
	}
}

// TestTemplate_ReplacePlaceholders tests simple placeholder replacement
func TestTemplate_ReplacePlaceholders(t *testing.T) {
	ctx := context.Background()

	// Create a test text file with placeholders
	tmpfile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	inputPath := tmpfile.Name()
	defer os.Remove(inputPath)

	content := "Hello {{name}}! Your email is {{email}}."
	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	outputPath := strings.Replace(inputPath, ".txt", "-output.txt", 1)
	defer os.Remove(outputPath)

	tmplNeta := templateneta.New()

	params := map[string]interface{}{
		"operation": "replace",
		"input":     inputPath,
		"output":    outputPath,
		"replacements": map[string]interface{}{
			"{{name}}":  "Alice",
			"{{email}}": "alice@example.com",
		},
		"mode": "placeholder",
	}

	result, err := tmplNeta.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	output, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected map[string]interface{} result")
	}

	if output["path"] != outputPath {
		t.Errorf("path = %v, want %v", output["path"], outputPath)
	}

	// Verify replacements were made
	resultContent, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}

	expected := "Hello Alice! Your email is alice@example.com."
	if string(resultContent) != expected {
		t.Errorf("Content = %q, want %q", string(resultContent), expected)
	}
}

// TestTemplate_MissingParameters tests error handling
func TestTemplate_MissingParameters(t *testing.T) {
	ctx := context.Background()
	tmplNeta := templateneta.New()

	tests := []struct {
		name   string
		params map[string]interface{}
	}{
		{
			name:   "missing operation",
			params: map[string]interface{}{"input": "test.svg", "output": "out.svg"},
		},
		{
			name:   "missing input",
			params: map[string]interface{}{"operation": "replace", "output": "out.svg"},
		},
		{
			name:   "missing output",
			params: map[string]interface{}{"operation": "replace", "input": "test.svg"},
		},
		{
			name:   "missing replacements",
			params: map[string]interface{}{"operation": "replace", "input": "test.svg", "output": "out.svg"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tmplNeta.Execute(ctx, tt.params)
			if err == nil {
				t.Error("Expected error for missing parameter, got nil")
			}
		})
	}
}

// TestTemplate_InvalidInput tests error handling for non-existent files
func TestTemplate_InvalidInput(t *testing.T) {
	ctx := context.Background()
	tmplNeta := templateneta.New()

	params := map[string]interface{}{
		"operation": "replace",
		"input":     "/nonexistent/file.svg",
		"output":    "/tmp/output.svg",
		"replacements": map[string]interface{}{
			"test": "value",
		},
	}

	_, err := tmplNeta.Execute(ctx, params)
	if err == nil {
		t.Error("Expected error for non-existent input file, got nil")
	}
}

// TestTemplate_RealWorldSVG tests with Figma-style IDs
func TestTemplate_RealWorldSVG(t *testing.T) {
	ctx := context.Background()

	// Create SVG with Figma-style IDs (with special characters)
	svgContent := `<?xml version="1.0" encoding="UTF-8"?>
<svg width="2000" height="2000" xmlns="http://www.w3.org/2000/svg">
  <text id="${{scale}}" fill="white">32mm - Heroic</text>
  <text id="${{name}}" fill="yellow">Yeoman (Shotgun)</text>
  <text id="{{studio}}" fill="white">Bite The Bullet Studio</text>
</svg>`

	tmpfile, err := os.CreateTemp("", "test-figma-*.svg")
	if err != nil {
		t.Fatal(err)
	}
	inputPath := tmpfile.Name()
	defer os.Remove(inputPath)

	if _, err := tmpfile.WriteString(svgContent); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	outputPath := strings.Replace(inputPath, ".svg", "-output.svg", 1)
	defer os.Remove(outputPath)

	tmplNeta := templateneta.New()

	params := map[string]interface{}{
		"operation": "replace",
		"input":     inputPath,
		"output":    outputPath,
		"replacements": map[string]interface{}{
			"${{scale}}": "28mm - Standard",
			"${{name}}":  "Combat Dog",
			"{{studio}}": "Heavy Handed Studio",
		},
		"mode": "id",
	}

	result, err := tmplNeta.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	output, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected map[string]interface{} result")
	}

	if output["path"] != outputPath {
		t.Errorf("path = %v, want %v", output["path"], outputPath)
	}

	// Verify replacements
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "28mm - Standard") {
		t.Error("Expected '28mm - Standard' in output")
	}
	if !strings.Contains(contentStr, "Combat Dog") {
		t.Error("Expected 'Combat Dog' in output")
	}
	if !strings.Contains(contentStr, "Heavy Handed Studio") {
		t.Error("Expected 'Heavy Handed Studio' in output")
	}
}
