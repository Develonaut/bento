# Phase 4: Hangiri Package (はんぎり - "Wooden Rice Tub")

**Duration:** 1 week
**Package:** `pkg/hangiri/`
**Dependencies:** `pkg/neta`, `pkg/omakase`

---

## TDD Philosophy

> **Write tests FIRST to define contracts**

Storage tests should verify:
1. Load .bento.json files and deserialize to neta.Definition
2. Save neta.Definition to .bento.json with proper formatting
3. Handle malformed JSON gracefully
4. Discover bento files in a directory
5. Parse complex nested structures (group neta with children)

---

## Phase Overview

The hangiri ("wooden rice tub") package manages storage of bento files. Like the traditional wooden tub used to store and cool sushi rice, hangiri provides a clean container for your bento definitions.

Key responsibilities:
- **Load/Save:** Serialize and deserialize .bento.json files
- **Discovery:** Find all bento files in a directory
- **Validation:** Ensure JSON is well-formed before parsing
- **Pretty printing:** Save with 2-space indentation for human readability

### Why "Hangiri"?

A hangiri (はんぎり) is the traditional wooden tub used by sushi chefs to store and cool freshly cooked sushi rice. It's a specialized container for a specific purpose—just as our storage package is a specialized container for bento definitions.

---

## Success Criteria

**Phase 4 Complete When:**
- [ ] LoadBento() loads and parses .bento.json files
- [ ] SaveBento() saves with pretty-printed JSON (2-space indent)
- [ ] DiscoverBentos() finds all .bento.json files in directory
- [ ] Handles malformed JSON with clear errors
- [ ] Handles nested group structures
- [ ] Integration tests for all operations
- [ ] Files < 250 lines each
- [ ] File-level documentation complete
- [ ] `/code-review` run with Karen + Colossus approval

---

## Test-First Approach

### Step 1: Define storage interface via tests

Create `pkg/hangiri/hangiri_test.go`:

```go
package hangiri_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/bento/pkg/hangiri"
	"github.com/yourusername/bento/pkg/neta"
)

// Test: Load a simple bento file
func TestHangiri_LoadSimpleBento(t *testing.T) {
	// Create a temporary .bento.json file
	tmpDir := t.TempDir()
	bentoPath := filepath.Join(tmpDir, "test.bento.json")

	bentoJSON := `{
  "id": "test-bento",
  "type": "group",
  "version": "1.0.0",
  "name": "Test Bento",
  "nodes": [
    {
      "id": "node-1",
      "type": "edit-fields",
      "version": "1.0.0",
      "name": "Set Fields",
      "position": {"x": 0, "y": 0},
      "metadata": {},
      "parameters": {
        "values": {"foo": "bar"}
      },
      "inputPorts": [],
      "outputPorts": []
    }
  ],
  "edges": [],
  "position": {"x": 0, "y": 0},
  "metadata": {},
  "parameters": {},
  "inputPorts": [],
  "outputPorts": []
}`

	os.WriteFile(bentoPath, []byte(bentoJSON), 0644)

	// Load the bento
	h := hangiri.New()
	def, err := h.LoadBento(bentoPath)
	if err != nil {
		t.Fatalf("LoadBento failed: %v", err)
	}

	// Verify structure
	if def.ID != "test-bento" {
		t.Errorf("ID = %s, want test-bento", def.ID)
	}

	if def.Type != "group" {
		t.Errorf("Type = %s, want group", def.Type)
	}

	if len(def.Nodes) != 1 {
		t.Errorf("Nodes = %d, want 1", len(def.Nodes))
	}

	if def.Nodes[0].ID != "node-1" {
		t.Errorf("Nodes[0].ID = %s, want node-1", def.Nodes[0].ID)
	}
}

// Test: Save a bento with pretty-printed JSON
func TestHangiri_SaveBento(t *testing.T) {
	tmpDir := t.TempDir()
	bentoPath := filepath.Join(tmpDir, "output.bento.json")

	def := &neta.Definition{
		ID:      "my-bento",
		Type:    "group",
		Version: "1.0.0",
		Name:    "My Bento",
		Nodes: []neta.Definition{
			{
				ID:      "node-1",
				Type:    "http-request",
				Version: "1.0.0",
				Name:    "Fetch Data",
				Parameters: map[string]interface{}{
					"url":    "https://api.example.com",
					"method": "GET",
				},
			},
		},
	}

	h := hangiri.New()
	err := h.SaveBento(bentoPath, def)
	if err != nil {
		t.Fatalf("SaveBento failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(bentoPath); os.IsNotExist(err) {
		t.Fatal("Bento file was not created")
	}

	// Read file and verify it's pretty-printed
	content, _ := os.ReadFile(bentoPath)
	contentStr := string(content)

	// Should have 2-space indentation
	if !strings.Contains(contentStr, "  \"id\": \"my-bento\"") {
		t.Error("JSON should be pretty-printed with 2-space indentation")
	}

	// Should be valid JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal(content, &parsed); err != nil {
		t.Errorf("Saved JSON is invalid: %v", err)
	}
}

// Test: Load and save should round-trip correctly
func TestHangiri_RoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	originalPath := filepath.Join(tmpDir, "original.bento.json")
	roundTripPath := filepath.Join(tmpDir, "roundtrip.bento.json")

	original := &neta.Definition{
		ID:      "round-trip-test",
		Type:    "group",
		Version: "1.0.0",
		Name:    "Round Trip Test",
		Parameters: map[string]interface{}{
			"nested": map[string]interface{}{
				"foo": "bar",
				"num": 42,
			},
		},
	}

	h := hangiri.New()

	// Save original
	if err := h.SaveBento(originalPath, original); err != nil {
		t.Fatalf("SaveBento failed: %v", err)
	}

	// Load it back
	loaded, err := h.LoadBento(originalPath)
	if err != nil {
		t.Fatalf("LoadBento failed: %v", err)
	}

	// Save it again
	if err := h.SaveBento(roundTripPath, loaded); err != nil {
		t.Fatalf("Second SaveBento failed: %v", err)
	}

	// Verify they're the same
	originalBytes, _ := os.ReadFile(originalPath)
	roundTripBytes, _ := os.ReadFile(roundTripPath)

	if !jsonEqual(originalBytes, roundTripBytes) {
		t.Error("Round-trip resulted in different JSON")
	}
}

// Test: Malformed JSON should return clear error
func TestHangiri_MalformedJSON(t *testing.T) {
	tmpDir := t.TempDir()
	bentoPath := filepath.Join(tmpDir, "bad.bento.json")

	// Write invalid JSON
	os.WriteFile(bentoPath, []byte(`{"id": "bad", "type": "group",`), 0644)

	h := hangiri.New()
	_, err := h.LoadBento(bentoPath)

	if err == nil {
		t.Fatal("Expected error for malformed JSON")
	}

	// Error should mention JSON parsing
	if !strings.Contains(err.Error(), "JSON") {
		t.Errorf("Error should mention JSON parsing: %v", err)
	}
}

// Test: Missing file should return clear error
func TestHangiri_MissingFile(t *testing.T) {
	h := hangiri.New()
	_, err := h.LoadBento("/nonexistent/file.bento.json")

	if err == nil {
		t.Fatal("Expected error for missing file")
	}

	if !strings.Contains(err.Error(), "not found") && !strings.Contains(err.Error(), "no such file") {
		t.Errorf("Error should mention file not found: %v", err)
	}
}
```

### Step 2: Test bento discovery

```go
// Test: Discover all .bento.json files in directory
func TestHangiri_DiscoverBentos(t *testing.T) {
	tmpDir := t.TempDir()

	// Create multiple .bento.json files
	createBentoFile(t, filepath.Join(tmpDir, "workflow1.bento.json"), "bento-1")
	createBentoFile(t, filepath.Join(tmpDir, "workflow2.bento.json"), "bento-2")

	// Create subdirectory with more bentos
	subDir := filepath.Join(tmpDir, "subdir")
	os.Mkdir(subDir, 0755)
	createBentoFile(t, filepath.Join(subDir, "workflow3.bento.json"), "bento-3")

	// Create non-bento file (should be ignored)
	os.WriteFile(filepath.Join(tmpDir, "readme.md"), []byte("# README"), 0644)

	h := hangiri.New()
	bentos, err := h.DiscoverBentos(tmpDir, true) // recursive=true

	if err != nil {
		t.Fatalf("DiscoverBentos failed: %v", err)
	}

	if len(bentos) != 3 {
		t.Errorf("Expected 3 bentos, found %d", len(bentos))
	}

	// Verify IDs
	ids := make(map[string]bool)
	for _, b := range bentos {
		ids[b.ID] = true
	}

	if !ids["bento-1"] || !ids["bento-2"] || !ids["bento-3"] {
		t.Error("Not all bentos were discovered")
	}
}

// Test: Non-recursive discovery should only find top-level bentos
func TestHangiri_DiscoverBentos_NonRecursive(t *testing.T) {
	tmpDir := t.TempDir()

	createBentoFile(t, filepath.Join(tmpDir, "workflow1.bento.json"), "bento-1")

	subDir := filepath.Join(tmpDir, "subdir")
	os.Mkdir(subDir, 0755)
	createBentoFile(t, filepath.Join(subDir, "workflow2.bento.json"), "bento-2")

	h := hangiri.New()
	bentos, err := h.DiscoverBentos(tmpDir, false) // recursive=false

	if err != nil {
		t.Fatalf("DiscoverBentos failed: %v", err)
	}

	if len(bentos) != 1 {
		t.Errorf("Expected 1 bento, found %d (should not be recursive)", len(bentos))
	}

	if bentos[0].ID != "bento-1" {
		t.Errorf("Found wrong bento: %s", bentos[0].ID)
	}
}

// Helper: Create a simple bento file
func createBentoFile(t *testing.T, path string, id string) {
	def := &neta.Definition{
		ID:      id,
		Type:    "group",
		Version: "1.0.0",
		Name:    "Test Bento",
	}

	h := hangiri.New()
	if err := h.SaveBento(path, def); err != nil {
		t.Fatalf("Failed to create bento file: %v", err)
	}
}

// Helper: Compare JSON for equality (ignoring whitespace)
func jsonEqual(a, b []byte) bool {
	var aMap, bMap map[string]interface{}
	json.Unmarshal(a, &aMap)
	json.Unmarshal(b, &bMap)

	return reflect.DeepEqual(aMap, bMap)
}
```

### Step 3: Test complex nested structures

```go
// Test: Load bento with deeply nested groups
func TestHangiri_NestedGroups(t *testing.T) {
	// CRITICAL FOR PHASE 8: Product automation has nested structures
	tmpDir := t.TempDir()
	bentoPath := filepath.Join(tmpDir, "nested.bento.json")

	bentoJSON := `{
  "id": "main-group",
  "type": "group",
  "version": "1.0.0",
  "name": "Main Group",
  "nodes": [
    {
      "id": "sub-group",
      "type": "group",
      "version": "1.0.0",
      "name": "Sub Group",
      "nodes": [
        {
          "id": "nested-node",
          "type": "edit-fields",
          "version": "1.0.0",
          "name": "Nested Node",
          "parameters": {"values": {"foo": "bar"}}
        }
      ],
      "edges": []
    }
  ],
  "edges": []
}`

	os.WriteFile(bentoPath, []byte(bentoJSON), 0644)

	h := hangiri.New()
	def, err := h.LoadBento(bentoPath)

	if err != nil {
		t.Fatalf("LoadBento failed: %v", err)
	}

	// Verify nested structure
	if len(def.Nodes) != 1 {
		t.Fatalf("Expected 1 top-level node, got %d", len(def.Nodes))
	}

	subGroup := def.Nodes[0]
	if subGroup.Type != "group" {
		t.Errorf("Sub-node should be group, got %s", subGroup.Type)
	}

	if len(subGroup.Nodes) != 1 {
		t.Fatalf("Expected 1 nested node, got %d", len(subGroup.Nodes))
	}

	if subGroup.Nodes[0].ID != "nested-node" {
		t.Errorf("Nested node ID = %s, want nested-node", subGroup.Nodes[0].ID)
	}
}
```

---

## File Structure

```
pkg/hangiri/
├── hangiri.go           # Main storage implementation (~200 lines)
├── load.go              # LoadBento implementation (~150 lines)
├── save.go              # SaveBento implementation (~100 lines)
├── discover.go          # DiscoverBentos implementation (~150 lines)
└── hangiri_test.go      # Integration tests (~400 lines)
```

---

## Implementation Guidance

**File: `pkg/hangiri/hangiri.go`**

```go
// Package hangiri provides storage and retrieval of bento files.
//
// "Hangiri" (はんぎり - "wooden rice tub") is the traditional wooden container
// used by sushi chefs to store and cool sushi rice. Similarly, this package
// provides a clean container for bento definitions.
//
// Usage:
//
//	h := hangiri.New()
//
//	// Load a bento
//	def, err := h.LoadBento("my-workflow.bento.json")
//
//	// Save a bento
//	err = h.SaveBento("output.bento.json", def)
//
//	// Discover all bentos in a directory
//	bentos, err := h.DiscoverBentos("/path/to/bentos", true)
//
// Learn more about JSON encoding in Go:
// https://go.dev/blog/json
package hangiri

import (
	"github.com/yourusername/bento/pkg/neta"
	"github.com/yourusername/bento/pkg/omakase"
)

// Hangiri manages bento file storage and retrieval.
type Hangiri struct {
	validator *omakase.Validator
}

// New creates a new Hangiri storage manager.
func New() *Hangiri {
	return &Hangiri{
		validator: omakase.New(),
	}
}

// LoadBento loads and parses a .bento.json file.
//
// Returns:
//   - *neta.Definition: The parsed bento definition
//   - error: Any error that occurred (file not found, invalid JSON, etc.)
func (h *Hangiri) LoadBento(path string) (*neta.Definition, error) {
	// Implementation in load.go
	return h.loadBento(path)
}

// SaveBento saves a bento definition to a .bento.json file.
//
// The JSON is pretty-printed with 2-space indentation for human readability.
//
// Returns:
//   - error: Any error that occurred during serialization or file writing
func (h *Hangiri) SaveBento(path string, def *neta.Definition) error {
	// Implementation in save.go
	return h.saveBento(path, def)
}

// DiscoverBentos finds all .bento.json files in a directory.
//
// Parameters:
//   - dir: Directory to search
//   - recursive: If true, searches subdirectories
//
// Returns:
//   - []*neta.Definition: All discovered bentos
//   - error: Any error that occurred
func (h *Hangiri) DiscoverBentos(dir string, recursive bool) ([]*neta.Definition, error) {
	// Implementation in discover.go
	return h.discoverBentos(dir, recursive)
}
```

**File: `pkg/hangiri/load.go`**

```go
package hangiri

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/yourusername/bento/pkg/neta"
)

// loadBento loads and parses a .bento.json file.
func (h *Hangiri) loadBento(path string) (*neta.Definition, error) {
	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("bento file not found: %s", path)
		}
		return nil, fmt.Errorf("failed to read bento file: %w", err)
	}

	// Parse JSON
	var def neta.Definition
	if err := json.Unmarshal(data, &def); err != nil {
		return nil, fmt.Errorf("invalid JSON in bento file %s: %w", path, err)
	}

	// Validate structure
	if err := h.validator.Validate(&def); err != nil {
		return nil, fmt.Errorf("invalid bento file %s: %w", path, err)
	}

	return &def, nil
}
```

**File: `pkg/hangiri/save.go`**

```go
package hangiri

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/yourusername/bento/pkg/neta"
)

// saveBento saves a bento definition to a .bento.json file.
func (h *Hangiri) saveBento(path string, def *neta.Definition) error {
	// Validate before saving
	if err := h.validator.Validate(def); err != nil {
		return fmt.Errorf("cannot save invalid bento: %w", err)
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(def, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize bento: %w", err)
	}

	// Write to file
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write bento file: %w", err)
	}

	return nil
}
```

**File: `pkg/hangiri/discover.go`**

```go
package hangiri

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yourusername/bento/pkg/neta"
)

// discoverBentos finds all .bento.json files in a directory.
func (h *Hangiri) discoverBentos(dir string, recursive bool) ([]*neta.Definition, error) {
	var bentos []*neta.Definition

	// Walk directory
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip if not recursive and not in root dir
		if !recursive && filepath.Dir(path) != dir {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check for .bento.json extension
		if !strings.HasSuffix(path, ".bento.json") {
			return nil
		}

		// Load the bento
		def, err := h.LoadBento(path)
		if err != nil {
			// Log warning but continue (don't fail entire discovery)
			fmt.Printf("Warning: Failed to load %s: %v\n", path, err)
			return nil
		}

		bentos = append(bentos, def)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to discover bentos: %w", err)
	}

	return bentos, nil
}
```

---

## Common Go Pitfalls to Avoid

1. **File paths:** Use filepath.Join, not string concatenation
   ```go
   // ❌ BAD - breaks on Windows
   path := dir + "/" + filename

   // ✅ GOOD
   path := filepath.Join(dir, filename)
   ```

2. **JSON indentation:** MarshalIndent requires empty prefix and 2-space indent
   ```go
   // ❌ BAD - no indentation
   data, _ := json.Marshal(def)

   // ✅ GOOD - pretty-printed
   data, _ := json.MarshalIndent(def, "", "  ")
   ```

3. **Recursive filepath.Walk:** Use filepath.SkipDir to skip directories
   ```go
   // ✅ GOOD - skip subdirectories when not recursive
   if !recursive && info.IsDir() && path != dir {
       return filepath.SkipDir
   }
   ```

---

## Critical for Phase 8

**Nested Structures:**
- Product automation bento will have nested groups (main -> loop -> process-product)
- Must handle deep nesting correctly
- Validation must be recursive

**Discovery:**
- `bento menu` command will use DiscoverBentos to list available bentos
- Should handle errors gracefully (one bad bento shouldn't break entire discovery)

---

## Bento Box Principle Checklist

- [ ] Files < 250 lines (hangiri.go ~200, load.go ~150, save.go ~100, discover.go ~150)
- [ ] Functions < 20 lines
- [ ] Single responsibility (storage only, no business logic)
- [ ] Clear error messages (mention file path, what went wrong)
- [ ] File-level documentation

---

## Phase Completion

**Phase 4 MUST end with:**

1. All tests passing (`go test ./pkg/hangiri/...`)
2. Run `/code-review` slash command
3. Address feedback from Karen and Colossus
4. Get explicit approval from both agents
5. Document any decisions in `.claude/strategy/`

**Do not proceed to Phase 5 until code review is approved.**

---

## Claude Prompt Template

```
I need to implement Phase 4: hangiri (storage package) following TDD principles.

Please read:
- .claude/strategy/phase-4-hangiri.md (this file)
- .claude/BENTO_BOX_PRINCIPLE.md

Then:

1. Create `pkg/hangiri/hangiri_test.go` with integration tests for:
   - LoadBento (simple and nested structures)
   - SaveBento (pretty-printed JSON)
   - Round-trip (load -> save -> load should be identical)
   - Malformed JSON (clear error messages)
   - Missing files (clear error messages)
   - DiscoverBentos (recursive and non-recursive)
   - Nested groups (Phase 8 needs this)

2. Watch the tests fail

3. Implement to make tests pass:
   - pkg/hangiri/hangiri.go (~200 lines)
   - pkg/hangiri/load.go (~150 lines)
   - pkg/hangiri/save.go (~100 lines)
   - pkg/hangiri/discover.go (~150 lines)

4. Add file-level documentation explaining:
   - What a hangiri is (wooden rice tub)
   - How to use LoadBento, SaveBento, DiscoverBentos
   - Common Go pitfalls (filepath.Join, MarshalIndent)

Remember:
- Write tests FIRST
- Files < 250 lines
- Functions < 20 lines
- Integration tests over unit tests
- Clear error messages (mention file path)

When complete, run `/code-review` and get Karen + Colossus approval.
```

---

## Dependencies

No additional dependencies needed - uses stdlib only:
- `encoding/json` for JSON parsing
- `os` for file I/O
- `path/filepath` for path operations

---

## Notes

- Pretty-printing (2-space indent) makes .bento.json files human-readable and git-friendly
- Validation on load prevents bad bentos from propagating
- Validation on save prevents creating invalid bentos
- Discovery errors should warn but not fail (one bad bento shouldn't break `bento menu`)
- Nested structures are critical for complex bentos (like Phase 8 product automation)

---

**Status:** Ready for implementation
**Next Phase:** Phase 5 (pantry registry) - depends on completion of Phase 4
