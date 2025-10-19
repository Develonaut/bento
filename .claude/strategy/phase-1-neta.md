# Phase 1: Neta Package (ネタ - "Ingredients")

**Duration:** 5-6 weeks (broken into 3 sub-phases)
**Package:** `pkg/neta/`
**Dependencies:** None (foundation layer)

---

## TDD Philosophy

> **Write tests FIRST to define contracts**

In this phase, tests are not validation—they are **design tools**. Before writing a single line of implementation code:

1. Write integration tests that describe HOW users will use each neta type
2. Let the tests fail (red)
3. Implement just enough to make tests pass (green)
4. Refactor while keeping tests green

**Integration over Unit:** Test like users would use the neta, not implementation details.

---

## Phase Overview

The neta package is the **foundation** of the entire bento system. It contains:
- Core type definitions (`Definition`, `Executable` interface, `Port`, `Edge`)
- All 10 neta type implementations
- Graph structures for workflow representation

This phase is broken into **3 sub-phases** to make it manageable:

### Phase 1a: Core + Simple Neta (2 weeks)
- Core types and interfaces
- `edit-fields` neta (field editor using text/template)
- `group` neta (sequential/parallel execution container)

### Phase 1b: I/O Neta (2 weeks)
- `http-request` neta (HTTP client using net/http)
- `file-system` neta (file operations)
- `shell-command` neta (execute shell commands with configurable timeouts)

### Phase 1c: Advanced Neta (2 weeks)
- `loop` neta (iteration: forEach, times, while)
- `parallel` neta (advanced parallelism with worker pools)
- `spreadsheet` neta (Excel/CSV using excelize)
- `image` neta (image processing using govips)
- `transform` neta (data transformation using expr)

---

## Success Criteria

**Phase 1a Complete When:**
- [ ] Core types (`Definition`, `Executable`, `Port`, `Edge`) defined and tested
- [ ] `edit-fields` neta fully implemented with integration tests
- [ ] `group` neta fully implemented with integration tests
- [ ] All tests passing
- [ ] Files < 250 lines each (Bento Box Principle)
- [ ] File-level documentation complete
- [ ] `/code-review` run with Karen + Colossus approval

**Phase 1b Complete When:**
- [ ] `http-request` neta with integration tests (GET, POST, headers, auth)
- [ ] `file-system` neta with integration tests (read, write, copy, move, delete)
- [ ] `shell-command` neta with tests for timeouts, streaming output, exit codes
- [ ] All tests passing
- [ ] `/code-review` run with Karen + Colossus approval

**Phase 1c Complete When:**
- [ ] `loop` neta (forEach, times, while) with integration tests
- [ ] `parallel` neta with worker pool tests
- [ ] `spreadsheet` neta (CSV/Excel read/write) with integration tests
- [ ] `image` neta (resize, format conversion, optimization) with integration tests
- [ ] `transform` neta (expr execution) with integration tests
- [ ] All 10 neta types complete
- [ ] `/code-review` run with Karen + Colossus approval

---

## Phase 1a: Core + Simple Neta

### Test-First Approach

**Step 1: Define core interface via tests**

Create `pkg/neta/neta_test.go`:

```go
package neta_test

import (
	"context"
	"testing"

	"github.com/yourusername/bento/pkg/neta"
)

// Test: Every neta must implement Executable interface
func TestExecutableInterface(t *testing.T) {
	ctx := context.Background()

	// This test ensures our interface contract is clear
	var _ neta.Executable = &mockNeta{}

	mock := &mockNeta{
		result: map[string]interface{}{"foo": "bar"},
	}

	result, err := mock.Execute(ctx, map[string]interface{}{})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

type mockNeta struct {
	result map[string]interface{}
}

func (m *mockNeta) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	return m.result, nil
}
```

**Step 2: Define Definition structure via test**

```go
// Test: Definition should serialize to/from JSON
func TestDefinitionJSON(t *testing.T) {
	def := neta.Definition{
		ID:      "test-node-1",
		Type:    "edit-fields",
		Version: "1.0.0",
		Name:    "Test Node",
		Position: neta.Position{X: 100, Y: 200},
		Parameters: map[string]interface{}{
			"values": map[string]interface{}{
				"foo": "bar",
			},
		},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(def)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Unmarshal back
	var decoded neta.Definition
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Verify round-trip
	if decoded.ID != def.ID {
		t.Errorf("ID mismatch: got %s, want %s", decoded.ID, def.ID)
	}
}
```

**Step 3: Integration test for edit-fields neta**

Create `pkg/neta/library/editfields/editfields_test.go`:

```go
package editfields_test

import (
	"context"
	"testing"

	"github.com/yourusername/bento/pkg/neta/library/editfields"
)

func TestEditFields_SetStaticValues(t *testing.T) {
	ctx := context.Background()

	ef := editfields.New()

	params := map[string]interface{}{
		"values": map[string]interface{}{
			"name": "Product A",
			"sku":  "PROD-001",
		},
	}

	result, err := ef.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	output, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected map[string]interface{} result")
	}

	if output["name"] != "Product A" {
		t.Errorf("name = %v, want Product A", output["name"])
	}

	if output["sku"] != "PROD-001" {
		t.Errorf("sku = %v, want PROD-001", output["sku"])
	}
}

func TestEditFields_TemplateVariables(t *testing.T) {
	ctx := context.Background()

	ef := editfields.New()

	// Context from previous neta
	prevContext := map[string]interface{}{
		"product": map[string]interface{}{
			"name": "Widget",
			"id":   123,
		},
	}

	params := map[string]interface{}{
		"values": map[string]interface{}{
			"title":    "{{.product.name}}",
			"filename": "product-{{.product.id}}.png",
		},
	}

	// Pass context through params (itamae will do this)
	params["_context"] = prevContext

	result, err := ef.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	output := result.(map[string]interface{})

	if output["title"] != "Widget" {
		t.Errorf("title = %v, want Widget", output["title"])
	}

	if output["filename"] != "product-123.png" {
		t.Errorf("filename = %v, want product-123.png", output["filename"])
	}
}
```

### File Structure (Phase 1a)

```
pkg/neta/
├── definition.go          # Core Definition type (~150 lines)
├── executable.go          # Executable interface + base types (~100 lines)
├── port.go                # Port and Edge types (~80 lines)
├── position.go            # Position, Metadata types (~50 lines)
├── neta_test.go           # Core interface tests (~200 lines)
│
└── library/
    ├── editfields/
    │   ├── editfields.go       # edit-fields implementation (~150 lines)
    │   ├── template.go         # Template parsing logic (~100 lines)
    │   └── editfields_test.go  # Integration tests (~200 lines)
    │
    └── group/
        ├── group.go            # group implementation (~200 lines)
        ├── execution.go        # Sequential/parallel execution (~150 lines)
        └── group_test.go       # Integration tests (~250 lines)
```

### Implementation Guidance

**File: `pkg/neta/definition.go`**

```go
// Package neta provides the core node type definitions for the bento workflow system.
//
// In bento, workflows are composed of "neta" (ネタ - ingredients) - individual nodes
// that can be connected together to form complex automation workflows.
//
// Every neta must implement the Executable interface and can be serialized to/from JSON.
//
// Learn more about Go interfaces: https://go.dev/tour/methods/9
package neta

import "encoding/json"

// Definition represents a single neta (node) in a bento workflow.
//
// Each neta has:
//   - A unique ID within the workflow
//   - A type (e.g., "http-request", "edit-fields")
//   - Parameters specific to that type
//   - Input/output ports for connecting to other neta
//   - Position (for visual editor, optional for CLI)
//
// Example JSON:
//
//	{
//	  "id": "node-1",
//	  "type": "http-request",
//	  "version": "1.0.0",
//	  "name": "Fetch User Data",
//	  "parameters": {
//	    "url": "https://api.example.com/users",
//	    "method": "GET"
//	  }
//	}
type Definition struct {
	ID          string                 `json:"id"`                    // Unique identifier
	Type        string                 `json:"type"`                  // Neta type (http-request, loop, etc.)
	Version     string                 `json:"version"`               // Schema version
	ParentID    *string                `json:"parentId,omitempty"`    // Parent group ID (if nested)
	Name        string                 `json:"name"`                  // Human-readable name
	Position    Position               `json:"position"`              // Visual editor position
	Metadata    Metadata               `json:"metadata"`              // Additional metadata
	Parameters  map[string]interface{} `json:"parameters"`            // Neta-specific parameters
	Fields      *FieldsConfig          `json:"fields,omitempty"`      // Field configuration
	InputPorts  []Port                 `json:"inputPorts"`            // Input connection points
	OutputPorts []Port                 `json:"outputPorts"`           // Output connection points
	Nodes       []Definition           `json:"nodes,omitempty"`       // Child nodes (for group neta)
	Edges       []Edge                 `json:"edges,omitempty"`       // Connections between child nodes
}

// Position represents the visual location of a neta in the editor.
// For CLI-only usage, this can be zero values.
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Metadata contains additional neta information.
type Metadata struct {
	CreatedAt  string            `json:"createdAt,omitempty"`
	UpdatedAt  string            `json:"updatedAt,omitempty"`
	Tags       []string          `json:"tags,omitempty"`
	CustomData map[string]string `json:"customData,omitempty"`
}

// Port represents an input or output connection point.
type Port struct {
	ID     string `json:"id"`                // Unique port identifier
	Name   string `json:"name"`              // Human-readable name
	Handle string `json:"handle,omitempty"`  // Handle type (for visual editor)
}

// Edge represents a connection between two neta.
type Edge struct {
	ID           string `json:"id"`                      // Unique edge identifier
	Source       string `json:"source"`                  // Source neta ID
	Target       string `json:"target"`                  // Target neta ID
	SourceHandle string `json:"sourceHandle,omitempty"`  // Source port handle
	TargetHandle string `json:"targetHandle,omitempty"`  // Target port handle
}

// FieldsConfig represents field editor configuration.
type FieldsConfig struct {
	Values      map[string]interface{} `json:"values"`                // Field values
	KeepOnlySet bool                   `json:"keepOnlySet,omitempty"` // Only output set fields
}
```

**File: `pkg/neta/executable.go`**

```go
package neta

import "context"

// Executable is the core interface that all neta types must implement.
//
// The Execute method runs the neta's logic and returns a result.
// The result is passed as context to connected neta in the workflow.
//
// Go interfaces are implicit - if a type has an Execute method with this
// signature, it automatically implements Executable.
//
// Learn more: https://go.dev/tour/methods/10
type Executable interface {
	// Execute runs the neta's logic.
	//
	// ctx: Go context for cancellation and deadlines
	// params: Neta-specific parameters + execution context from previous neta
	//
	// Returns:
	//   - interface{}: Result data (usually map[string]interface{})
	//   - error: Any error that occurred during execution
	Execute(ctx context.Context, params map[string]interface{}) (interface{}, error)
}

// ExecutionContext contains data passed between neta during workflow execution.
// This is managed by the itamae (orchestration engine).
type ExecutionContext struct {
	Data   map[string]interface{} // Accumulated data from previous neta
	NodeID string                 // Current neta ID
	Depth  int                    // Execution depth (for nested groups)
}
```

### Common Go Pitfalls to Avoid

1. **nil maps**: Always initialize maps before use
   ```go
   // ❌ BAD - will panic
   var params map[string]interface{}
   params["foo"] = "bar"  // PANIC!

   // ✅ GOOD
   params := make(map[string]interface{})
   params["foo"] = "bar"
   ```

2. **interface{} type assertions**: Always check the "ok" value
   ```go
   // ❌ BAD - can panic
   str := result.(string)

   // ✅ GOOD
   str, ok := result.(string)
   if !ok {
       return nil, fmt.Errorf("expected string, got %T", result)
   }
   ```

3. **JSON field tags**: Required for proper serialization
   ```go
   // ❌ BAD - won't serialize correctly
   type Node struct {
       ID string
   }

   // ✅ GOOD
   type Node struct {
       ID string `json:"id"`
   }
   ```

---

## Phase 1b: I/O Neta

### Integration Tests

**File: `pkg/neta/library/http/http_test.go`**

```go
func TestHTTPRequest_GET(t *testing.T) {
	ctx := context.Background()

	// Start test HTTP server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"data":   []string{"item1", "item2"},
		})
	}))
	defer ts.Close()

	httpNeta := httpneta.New()

	params := map[string]interface{}{
		"url":    ts.URL,
		"method": "GET",
	}

	result, err := httpNeta.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Verify response
	response := result.(map[string]interface{})
	if response["status"] != "success" {
		t.Errorf("status = %v, want success", response["status"])
	}
}

func TestHTTPRequest_Timeout(t *testing.T) {
	ctx := context.Background()

	// Server that delays 5 seconds
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
	}))
	defer ts.Close()

	httpNeta := httpneta.New()

	params := map[string]interface{}{
		"url":     ts.URL,
		"method":  "GET",
		"timeout": 1, // 1 second timeout
	}

	_, err := httpNeta.Execute(ctx, params)
	if err == nil {
		t.Fatal("Expected timeout error, got nil")
	}

	if !strings.Contains(err.Error(), "timeout") {
		t.Errorf("Expected timeout error, got: %v", err)
	}
}
```

**File: `pkg/neta/library/shellcommand/shellcommand_test.go`**

```go
func TestShellCommand_BasicExecution(t *testing.T) {
	ctx := context.Background()

	sc := shellcommand.New()

	params := map[string]interface{}{
		"command": "echo",
		"args":    []string{"hello", "world"},
	}

	result, err := sc.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	output := result.(map[string]interface{})
	stdout := output["stdout"].(string)

	if !strings.Contains(stdout, "hello world") {
		t.Errorf("stdout = %q, want to contain 'hello world'", stdout)
	}

	if output["exitCode"] != 0 {
		t.Errorf("exitCode = %v, want 0", output["exitCode"])
	}
}

func TestShellCommand_LongRunning(t *testing.T) {
	// CRITICAL FOR PHASE 8: Blender renders take minutes
	ctx := context.Background()

	sc := shellcommand.New()

	params := map[string]interface{}{
		"command": "sleep",
		"args":    []string{"3"},  // 3 second sleep
		"timeout": 10,             // 10 second timeout (don't kill it)
	}

	start := time.Now()
	result, err := sc.Execute(ctx, params)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if duration < 3*time.Second {
		t.Errorf("Command finished too quickly: %v", duration)
	}

	output := result.(map[string]interface{})
	if output["exitCode"] != 0 {
		t.Errorf("exitCode = %v, want 0", output["exitCode"])
	}
}

func TestShellCommand_StreamingOutput(t *testing.T) {
	// CRITICAL FOR PHASE 8: Stream Blender render progress
	ctx := context.Background()

	sc := shellcommand.New()

	// Command that outputs multiple lines
	params := map[string]interface{}{
		"command": "bash",
		"args":    []string{"-c", "for i in 1 2 3; do echo $i; sleep 0.1; done"},
		"stream":  true,  // Enable streaming output
	}

	outputLines := []string{}

	// Mock callback for streaming (itamae will provide this)
	onOutput := func(line string) {
		outputLines = append(outputLines, line)
	}
	params["_onOutput"] = onOutput

	result, err := sc.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Verify we got streaming output
	if len(outputLines) < 3 {
		t.Errorf("Expected at least 3 output lines, got %d", len(outputLines))
	}
}
```

### File Structure (Phase 1b)

```
pkg/neta/library/
├── http/
│   ├── http.go           # HTTP request implementation (~200 lines)
│   ├── client.go         # HTTP client configuration (~100 lines)
│   └── http_test.go      # Integration tests (~300 lines)
│
├── filesystem/
│   ├── filesystem.go     # File operations (~200 lines)
│   ├── operations.go     # Read, write, copy, move (~150 lines)
│   └── filesystem_test.go # Integration tests (~250 lines)
│
└── shellcommand/
    ├── shellcommand.go   # Shell execution (~200 lines)
    ├── streaming.go      # Output streaming logic (~100 lines)
    └── shellcommand_test.go # Integration tests (~300 lines)
```

---

## Phase 1c: Advanced Neta

### Integration Tests

**File: `pkg/neta/library/loop/loop_test.go`**

```go
func TestLoop_ForEach(t *testing.T) {
	// CRITICAL FOR PHASE 8: Loop through CSV rows
	ctx := context.Background()

	loop := loopneta.New()

	// Simulate CSV rows
	csvData := []map[string]interface{}{
		{"sku": "PROD-001", "name": "Product A"},
		{"sku": "PROD-002", "name": "Product B"},
		{"sku": "PROD-003", "name": "Product C"},
	}

	params := map[string]interface{}{
		"mode":  "forEach",
		"items": csvData,
		"body": map[string]interface{}{
			// This would be a neta definition in real usage
			"type": "edit-fields",
			"parameters": map[string]interface{}{
				"values": map[string]interface{}{
					"folder": "products/{{.item.sku}}",
				},
			},
		},
	}

	result, err := loop.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	output := result.(map[string]interface{})
	iterations := output["iterations"].(int)

	if iterations != 3 {
		t.Errorf("iterations = %d, want 3", iterations)
	}
}
```

**File: `pkg/neta/library/spreadsheet/spreadsheet_test.go`**

```go
func TestSpreadsheet_ReadCSV(t *testing.T) {
	// CRITICAL FOR PHASE 8: Read product CSV
	ctx := context.Background()

	// Create test CSV
	csvContent := `sku,name,description
PROD-001,Product A,Description A
PROD-002,Product B,Description B
PROD-003,Product C,Description C`

	tmpfile, _ := os.CreateTemp("", "test-*.csv")
	tmpfile.WriteString(csvContent)
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	ss := spreadsheet.New()

	params := map[string]interface{}{
		"operation": "read",
		"format":    "csv",
		"path":      tmpfile.Name(),
	}

	result, err := ss.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	output := result.(map[string]interface{})
	rows := output["rows"].([]map[string]interface{})

	if len(rows) != 3 {
		t.Errorf("rows = %d, want 3", len(rows))
	}

	if rows[0]["sku"] != "PROD-001" {
		t.Errorf("rows[0].sku = %v, want PROD-001", rows[0]["sku"])
	}
}
```

**File: `pkg/neta/library/image/image_test.go`**

```go
func TestImage_OptimizeToWebP(t *testing.T) {
	// CRITICAL FOR PHASE 8: Optimize Blender output to webp
	ctx := context.Background()

	// Create test PNG image (simple 100x100 red square)
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	red := color.RGBA{255, 0, 0, 255}
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, red)
		}
	}

	tmpPNG, _ := os.CreateTemp("", "test-*.png")
	png.Encode(tmpPNG, img)
	tmpPNG.Close()
	defer os.Remove(tmpPNG.Name())

	imgNeta := imageneta.New()

	params := map[string]interface{}{
		"operation": "convert",
		"input":     tmpPNG.Name(),
		"output":    strings.Replace(tmpPNG.Name(), ".png", ".webp", 1),
		"format":    "webp",
		"quality":   80,
	}

	result, err := imgNeta.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	output := result.(map[string]interface{})

	// Verify output file exists
	if _, err := os.Stat(output["path"].(string)); os.IsNotExist(err) {
		t.Errorf("Output file not created: %v", output["path"])
	}

	// Verify size reduction (webp should be smaller)
	originalSize, _ := getFileSize(tmpPNG.Name())
	webpSize, _ := getFileSize(output["path"].(string))

	if webpSize >= originalSize {
		t.Logf("Warning: webp (%d bytes) not smaller than PNG (%d bytes)", webpSize, originalSize)
	}

	os.Remove(output["path"].(string))
}
```

### File Structure (Phase 1c)

```
pkg/neta/library/
├── loop/
│   ├── loop.go               # Loop implementation (~200 lines)
│   ├── modes.go              # forEach, times, while (~150 lines)
│   └── loop_test.go          # Integration tests (~250 lines)
│
├── parallel/
│   ├── parallel.go           # Parallel execution (~200 lines)
│   ├── workers.go            # Worker pool management (~150 lines)
│   └── parallel_test.go      # Integration tests (~200 lines)
│
├── spreadsheet/
│   ├── spreadsheet.go        # Spreadsheet operations (~200 lines)
│   ├── csv.go                # CSV parsing (~100 lines)
│   ├── excel.go              # Excel using excelize (~150 lines)
│   └── spreadsheet_test.go   # Integration tests (~300 lines)
│
├── image/
│   ├── image.go              # Image processing (~200 lines)
│   ├── govips.go             # govips wrapper (~150 lines)
│   ├── operations.go         # Resize, convert, optimize (~150 lines)
│   └── image_test.go         # Integration tests (~250 lines)
│
└── transform/
    ├── transform.go          # Transform using expr (~200 lines)
    ├── expr.go               # expr execution wrapper (~100 lines)
    └── transform_test.go     # Integration tests (~200 lines)
```

---

## Bento Box Principle Checklist

For each file created, verify:

- [ ] File < 250 lines (500 max)
- [ ] Functions < 20 lines (30 max)
- [ ] Single responsibility (file does ONE thing)
- [ ] No utility grab bags
- [ ] Clear interfaces
- [ ] File-level documentation explaining purpose
- [ ] Examples in doc comments

---

## Phase Completion

**Each sub-phase (1a, 1b, 1c) MUST end with:**

1. All tests passing (`go test ./pkg/neta/...`)
2. Run `/code-review` slash command
3. Address feedback from Karen and Colossus
4. Get explicit approval from both agents
5. Document any architectural decisions in `.claude/strategy/`

**Do not proceed to the next sub-phase until code review is approved.**

---

## Claude Prompt Template

### For Phase 1a

```
I need to implement Phase 1a of the neta package following TDD principles.

Please read:
- .claude/strategy/phase-1-neta.md (this file)
- .claude/BENTO_BOX_PRINCIPLE.md

Then:

1. Create `pkg/neta/neta_test.go` with integration tests for core types (Definition, Executable interface)
2. Watch the tests fail (they should - we haven't implemented anything yet)
3. Implement `pkg/neta/definition.go` and `pkg/neta/executable.go` to make tests pass
4. Create `pkg/neta/library/editfields/editfields_test.go` with integration tests for edit-fields neta
5. Implement `pkg/neta/library/editfields/editfields.go` to make tests pass
6. Repeat for group neta

Remember:
- Write tests FIRST
- Keep files < 250 lines
- Keep functions < 20 lines
- Add file-level documentation (I'm new to Go)
- Integration tests over unit tests

When complete, run `/code-review` and get Karen + Colossus approval.
```

### For Phase 1b

```
I need to implement Phase 1b of the neta package (I/O neta: http-request, file-system, shell-command).

Please read:
- .claude/strategy/phase-1-neta.md (Phase 1b section)
- .claude/BENTO_BOX_PRINCIPLE.md

Then, for EACH neta type (http, filesystem, shellcommand):

1. Create `pkg/neta/library/[type]/[type]_test.go` with integration tests FIRST
2. Watch tests fail
3. Implement `pkg/neta/library/[type]/[type].go` to make tests pass
4. Add file-level documentation

CRITICAL for shell-command:
- Support configurable timeouts (default 2 min, but Phase 8 needs longer for Blender)
- Support streaming output (callback for each line)
- Capture exit codes
- Test with long-running commands (sleep)

When complete, run `/code-review` and get Karen + Colossus approval.
```

### For Phase 1c

```
I need to implement Phase 1c of the neta package (advanced neta: loop, parallel, spreadsheet, image, transform).

Please read:
- .claude/strategy/phase-1-neta.md (Phase 1c section)
- .claude/BENTO_BOX_PRINCIPLE.md
- .claude/COMPLETE_NODE_INVENTORY.md (reference for TypeScript implementations)

Then, for EACH neta type:

1. Create integration tests FIRST
2. Implement to make tests pass
3. Add file-level documentation

CRITICAL requirements:
- spreadsheet: Must handle CSV (Phase 8 uses CSV for product data)
- loop: Must support forEach mode (Phase 8 loops through CSV rows)
- image: Must support webp optimization using govips (Phase 8 needs this)

Libraries to use:
- spreadsheet: github.com/xuri/excelize/v2 for Excel, encoding/csv for CSV
- image: github.com/davidbyttow/govips/v2
- transform: github.com/antonmedv/expr

When complete, run `/code-review` and get Karen + Colossus approval.
```

---

## Dependencies to Add (go.mod)

```bash
go get github.com/antonmedv/expr
go get github.com/davidbyttow/govips/v2
go get github.com/xuri/excelize/v2
```

---

## Notes

- This is the FOUNDATION of the entire system - take time to get it right
- Tests are design tools, not just validation
- File-level docs are critical (user is new to Go)
- Phase 8 depends heavily on: shell-command, loop, spreadsheet, image neta
- Shell-command must handle multi-minute timeouts for Blender renders
- Loop must handle iterating CSV rows
- Spreadsheet must parse CSV with headers
- Image must optimize to webp using govips

---

**Status:** Ready for implementation
**Next Phase:** Phase 2 (shoyu logger) - depends on completion of Phase 1a
