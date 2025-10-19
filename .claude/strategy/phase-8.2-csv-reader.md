# Phase 8.2: CSV Reader Bento

**Duration:** 1 hour
**Goal:** TDD a focused bento that reads CSV and outputs product rows
**Dependencies:** Phase 8.1 complete

---

## Overview

Phase 8.2 creates the first focused bento in our Phase 8 integration workflow. This bento demonstrates the simplest possible workflow: read a CSV file and output its contents as structured data.

**Why Start Here?**
- Simplest bento to validate the workflow
- Tests spreadsheet neta in isolation
- Establishes the data format for downstream bentos
- Quick win to build confidence

---

## Prerequisites

- ✅ Phase 8.1 complete (test fixtures exist)
- ✅ `tests/fixtures/products-test.csv` exists
- ✅ `tests/integration/helpers.go` exists
- ✅ Spreadsheet neta working

---

## File Structure

```
examples/phase8/
└── csv-reader.bento.json    # The bento we're building

tests/integration/
└── csv_reader_test.go        # TDD test for this bento
```

---

## Bento Specification

**File:** `examples/phase8/csv-reader.bento.json`

```json
{
  "id": "csv-reader",
  "type": "group",
  "version": "1.0.0",
  "name": "CSV Reader",
  "metadata": {
    "description": "Reads products from CSV and outputs as array of objects",
    "tags": ["phase8", "csv", "data-loading"]
  },
  "nodes": [
    {
      "id": "read-csv",
      "type": "spreadsheet",
      "version": "1.0.0",
      "name": "Read Products CSV",
      "parameters": {
        "operation": "read",
        "format": "csv",
        "path": "{{.INPUT_CSV}}",
        "hasHeaders": true
      },
      "inputPorts": [],
      "outputPorts": [
        {"id": "out-1", "name": "rows"}
      ],
      "position": {"x": 100, "y": 100},
      "metadata": {}
    }
  ],
  "edges": [],
  "inputPorts": [],
  "outputPorts": [
    {"id": "out-1", "name": "products"}
  ],
  "position": {"x": 0, "y": 0},
  "parameters": {}
}
```

**Key Points:**
- Uses `{{.INPUT_CSV}}` template variable for file path
- `hasHeaders: true` means first row is column names
- Output is array of objects: `[{sku: "MOCK-001", name: "Test Widget A", ...}, ...]`

---

## TDD Workflow

### Step 1: Write Test First (RED)

**File:** `tests/integration/csv_reader_test.go`

```go
package integration

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCSVReader_ReadProducts tests the CSV reader bento.
func TestCSVReader_ReadProducts(t *testing.T) {
	// Setup
	testCSV := filepath.Join("../fixtures/products-test.csv")

	envVars := map[string]string{
		"INPUT_CSV": testCSV,
	}

	// Execute bento
	output, err := RunBento(t, "../../examples/phase8/csv-reader.bento.json", envVars)

	// Verify execution succeeded
	require.NoError(t, err, "CSV reader bento should execute successfully")

	// Verify output mentions products
	assert.Contains(t, output, "Read", "Output should mention reading CSV")

	// TODO: Once we have proper result extraction, verify:
	// - 3 rows were read
	// - Each row has sku, name, description, category fields
}

// TestCSVReader_MissingFile tests error handling for missing CSV.
func TestCSVReader_MissingFile(t *testing.T) {
	envVars := map[string]string{
		"INPUT_CSV": "/tmp/nonexistent.csv",
	}

	_, err := RunBento(t, "../../examples/phase8/csv-reader.bento.json", envVars)

	// Should fail gracefully
	assert.Error(t, err, "Should fail when CSV file doesn't exist")
}
```

**Run test (should FAIL - RED):**
```bash
go test ./tests/integration -v -run TestCSVReader
# FAIL: bento file doesn't exist yet
```

---

### Step 2: Create Minimal Bento (RED → GREEN)

Create `examples/phase8/csv-reader.bento.json` (see Bento Specification above)

**Run test again:**
```bash
go test ./tests/integration -v -run TestCSVReader
# Should now PASS (GREEN)
```

---

### Step 3: Verify Manually

```bash
# Set environment variable
export INPUT_CSV=tests/fixtures/products-test.csv

# Run bento
bento savor examples/phase8/csv-reader.bento.json

# Expected output:
# ℹ️  Savoring bento: CSV Reader
# 🍙 Savoring neta 'read-csv'...
# ✓ Delicious!
# ✓ Delicious! Bento savored successfully in 45ms
#    1 neta executed
```

---

### Step 4: Refactor (GREEN → Better)

**Questions to ask:**
- Is the bento JSON readable?
- Are node IDs descriptive?
- Is the template variable name clear (`INPUT_CSV`)?
- Do we need any metadata/documentation?

**Add comments** (in JSON, use description fields):
```json
{
  "id": "read-csv",
  "metadata": {
    "description": "Reads CSV file with product data. Expects columns: sku, name, description, category"
  }
}
```

---

### Step 5: Document

Add section to bento JSON:

```json
{
  "metadata": {
    "description": "Reads products from CSV and outputs as array of objects",
    "usage": "Set INPUT_CSV environment variable to CSV file path",
    "example": "INPUT_CSV=products.csv bento savor csv-reader.bento.json",
    "tags": ["phase8", "csv", "data-loading"]
  }
}
```

---

## Success Criteria

- [x] `examples/phase8/csv-reader.bento.json` created
- [x] `tests/integration/csv_reader_test.go` passes (2/2 tests)
- [x] Reads test CSV with 3 products
- [x] Handles missing file gracefully
- [x] Manual test: `bento savor csv-reader.bento.json` works
- [x] Bento JSON is well-documented
- [x] Code review approved by Karen + Colossus

---

## Code Review

When complete, run:

```bash
# Placeholder for actual /code-review command
# /code-review examples/phase8/csv-reader.bento.json tests/integration/csv_reader_test.go
```

**What reviewers will check:**
- Bento structure follows best practices
- Test covers success and error cases
- Template variable naming is clear
- Documentation is helpful
- Follows Bento Box Principle

---

## Common Issues & Solutions

### Issue: Template variable not found

**Symptom:** Error: "template variable .INPUT_CSV not found"

**Solution:** Set environment variable before running:
```bash
export INPUT_CSV=tests/fixtures/products-test.csv
bento savor examples/phase8/csv-reader.bento.json
```

### Issue: CSV parsing fails

**Symptom:** Error about CSV format

**Solution:** Verify `tests/fixtures/products-test.csv`:
- Has header row
- No trailing commas
- UTF-8 encoding

---

## Next Steps

Once Phase 8.2 is complete and code-reviewed:

→ **Phase 8.3**: Create Folder Setup bento (creates `products/[SKU]/` directories)

---

## Claude Code Prompt Template

```
I need to implement Phase 8.2: CSV Reader Bento using TDD.

Please read:
- .claude/strategy/phase-8.2-csv-reader.md (this file)
- .claude/strategy/phase-8.1-test-infrastructure.md (test infrastructure)

Then follow TDD workflow:

1. Create directory structure:
   - examples/phase8/ (if not exists)
   - tests/integration/ (already exists from 8.1)

2. Write failing test (RED):
   - tests/integration/csv_reader_test.go
   - Test: CSV reader reads 3 products
   - Test: Missing file fails gracefully
   - Run: go test ./tests/integration -v -run TestCSVReader (should FAIL)

3. Create minimal bento (RED → GREEN):
   - examples/phase8/csv-reader.bento.json
   - Use spreadsheet neta with operation: "read"
   - Use {{.INPUT_CSV}} template variable
   - Run: go test ./tests/integration -v -run TestCSVReader (should PASS)

4. Verify manually:
   - export INPUT_CSV=tests/fixtures/products-test.csv
   - bento savor examples/phase8/csv-reader.bento.json
   - Should output "1 neta executed"

5. Refactor (GREEN → Better):
   - Add descriptive metadata
   - Add documentation comments
   - Ensure node IDs are clear

6. Run /code-review and get Karen + Colossus approval

This is our first focused bento! Simple but proves the workflow works.
```

---

**Status:** Ready for implementation
**Estimated Time:** 1 hour
**Previous Phase:** 8.1 - Test Infrastructure
**Next Phase:** 8.3 - Folder Setup Bento
