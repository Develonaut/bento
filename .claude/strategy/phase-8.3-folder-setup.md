# Phase 8.3: Folder Setup Bento

**Duration:** 1 hour
**Goal:** TDD a bento that creates product folder structure using loop neta
**Dependencies:** Phase 8.2 complete

---

## Overview

Phase 8.3 demonstrates looping through CSV data and performing file system operations for each item. This bento reads the CSV (using the bento from 8.2) and creates a folder for each product.

**Why This Phase?**
- Demonstrates loop neta with forEach mode
- Combines multiple neta types (spreadsheet + file-system + loop)
- Tests context passing between neta (CSV data → loop items)
- Real-world pattern: "for each X, create Y"

---

## Prerequisites

- ✅ Phase 8.1 complete (test fixtures)
- ✅ Phase 8.2 complete (CSV reader bento)
- ✅ `tests/fixtures/products-test.csv` has 3 products
- ✅ Loop neta and file-system neta working

---

## File Structure

```
examples/phase8/
├── csv-reader.bento.json      # From 8.2
└── folder-setup.bento.json    # NEW - what we're building

tests/integration/
├── csv_reader_test.go          # From 8.2
└── folder_setup_test.go        # NEW - TDD test
```

---

## Bento Specification

**File:** `examples/phase8/folder-setup.bento.json`

```json
{
  "id": "folder-setup",
  "type": "group",
  "version": "1.0.0",
  "name": "Product Folder Setup",
  "metadata": {
    "description": "Creates folder structure for each product from CSV",
    "usage": "INPUT_CSV=products.csv bento savor folder-setup.bento.json",
    "tags": ["phase8", "folders", "loop", "filesystem"]
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
      "outputPorts": [
        {"id": "out-1", "name": "rows"}
      ],
      "position": {"x": 100, "y": 100},
      "metadata": {}
    },
    {
      "id": "create-folders",
      "type": "loop",
      "version": "1.0.0",
      "name": "Create Folder for Each Product",
      "parameters": {
        "mode": "forEach",
        "items": "{{.read-csv.rows}}",
        "showProgress": true
      },
      "inputPorts": [
        {"id": "in-1", "name": "items"}
      ],
      "position": {"x": 300, "y": 100},
      "metadata": {
        "description": "Loops through each product row and creates its folder"
      },
      "nodes": [
        {
          "id": "mkdir",
          "type": "file-system",
          "version": "1.0.0",
          "name": "Create Product Directory",
          "parameters": {
            "operation": "mkdir",
            "path": "products/{{.item.sku}}",
            "recursive": true
          },
          "position": {"x": 100, "y": 100},
          "metadata": {
            "description": "Creates products/[SKU]/ directory. {{.item}} is the current row."
          }
        }
      ],
      "edges": []
    }
  ],
  "edges": [
    {
      "id": "e1",
      "source": "read-csv",
      "target": "create-folders",
      "sourcePort": "out-1",
      "targetPort": "in-1"
    }
  ],
  "inputPorts": [],
  "outputPorts": [],
  "position": {"x": 0, "y": 0},
  "parameters": {}
}
```

**Key Concepts:**
- `{{.read-csv.rows}}` - Reference output from previous neta
- `{{.item.sku}}` - Inside loop, `item` is current CSV row
- `recursive: true` - Creates parent directories if needed
- Nested nodes inside loop neta

---

## TDD Workflow

### Step 1: Write Test First (RED)

**File:** `tests/integration/folder_setup_test.go`

```go
package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFolderSetup_CreateProductFolders tests the folder setup bento.
func TestFolderSetup_CreateProductFolders(t *testing.T) {
	// Setup - create temp output directory
	outputDir := filepath.Join(os.TempDir(), "bento-test-folders")
	defer CleanupTestDir(t, outputDir)

	// Change to temp directory so products/ is created there
	originalDir, _ := os.Getwd()
	os.MkdirAll(outputDir, 0755)
	os.Chdir(outputDir)
	defer os.Chdir(originalDir)

	testCSV := filepath.Join(originalDir, "../fixtures/products-test.csv")

	envVars := map[string]string{
		"INPUT_CSV": testCSV,
	}

	// Execute bento
	bentoPath := filepath.Join(originalDir, "../../examples/phase8/folder-setup.bento.json")
	output, err := RunBento(t, bentoPath, envVars)

	// Verify execution succeeded
	require.NoError(t, err, "Folder setup bento should execute successfully")

	// Verify output shows progress
	assert.Contains(t, output, "Savoring", "Should show execution progress")

	// Verify folders were created
	VerifyFileExists(t, "products/MOCK-001")
	VerifyFileExists(t, "products/MOCK-002")
	VerifyFileExists(t, "products/MOCK-003")

	// Verify they are directories
	info, _ := os.Stat("products/MOCK-001")
	assert.True(t, info.IsDir(), "products/MOCK-001 should be a directory")
}

// TestFolderSetup_AlreadyExists tests idempotency.
func TestFolderSetup_AlreadyExists(t *testing.T) {
	outputDir := filepath.Join(os.TempDir(), "bento-test-folders-2")
	defer CleanupTestDir(t, outputDir)

	originalDir, _ := os.Getwd()
	os.MkdirAll(outputDir, 0755)
	os.Chdir(outputDir)
	defer os.Chdir(originalDir)

	// Pre-create one folder
	os.MkdirAll("products/MOCK-001", 0755)

	testCSV := filepath.Join(originalDir, "../fixtures/products-test.csv")
	envVars := map[string]string{
		"INPUT_CSV": testCSV,
	}

	bentoPath := filepath.Join(originalDir, "../../examples/phase8/folder-setup.bento.json")
	output, err := RunBento(t, bentoPath, envVars)

	// Should still succeed (idempotent)
	require.NoError(t, err, "Should succeed even if folder exists")

	// All folders should exist
	VerifyFileExists(t, "products/MOCK-001")
	VerifyFileExists(t, "products/MOCK-002")
	VerifyFileExists(t, "products/MOCK-003")
}
```

**Run test (should FAIL - RED):**
```bash
go test ./tests/integration -v -run TestFolderSetup
# FAIL: bento file doesn't exist yet
```

---

### Step 2: Create Bento (RED → GREEN)

Create `examples/phase8/folder-setup.bento.json` (see specification above)

**Run test again:**
```bash
go test ./tests/integration -v -run TestFolderSetup
# Should now PASS (GREEN)
```

---

### Step 3: Verify Manually

```bash
# Clean up any previous test
rm -rf products/

# Set environment variable
export INPUT_CSV=tests/fixtures/products-test.csv

# Run bento
bento savor examples/phase8/folder-setup.bento.json

# Verify folders created
ls -la products/
# Should show:
# products/MOCK-001/
# products/MOCK-002/
# products/MOCK-003/
```

---

### Step 4: Refactor (GREEN → Better)

**Improvements to consider:**
- Add better progress messages
- Add metadata explaining the loop structure
- Consider if we need error handling for mkdir failures

**Example refactor:**
```json
{
  "id": "mkdir",
  "metadata": {
    "description": "Creates folder for product. Uses {{.item.sku}} from loop context.",
    "example": "MOCK-001 → products/MOCK-001/"
  }
}
```

---

## Success Criteria

- [x] `examples/phase8/folder-setup.bento.json` created
- [x] `tests/integration/folder_setup_test.go` passes (2/2 tests)
- [x] Creates `products/MOCK-001/`, `MOCK-002/`, `MOCK-003/`
- [x] Idempotent (running twice doesn't fail)
- [x] Manual test: folders created correctly
- [x] Loop context (`{{.item.sku}}`) works correctly
- [x] Edge connection from CSV reader to loop works
- [x] Code review approved by Karen + Colossus

---

## Code Review

When complete, run:

```bash
# Placeholder for actual /code-review command
# /code-review examples/phase8/folder-setup.bento.json tests/integration/folder_setup_test.go
```

**What reviewers will check:**
- Loop neta configured correctly
- Context passing works (CSV rows → loop items)
- Template variables used properly (`{{.item.sku}}`)
- Tests cover success and edge cases
- Folder creation is idempotent

---

## Common Issues & Solutions

### Issue: "Template variable .item not found"

**Symptom:** Error about `.item` not found inside loop

**Solution:** Check loop configuration:
- `mode` must be "forEach"
- `items` must reference CSV output: `{{.read-csv.rows}}`
- Inside loop, use `{{.item.fieldname}}`

### Issue: Folders not created

**Symptom:** Test fails, no folders exist

**Solution:** Check:
1. Edge connects CSV reader to loop
2. `recursive: true` on mkdir operation
3. Path uses correct template variable: `products/{{.item.sku}}`

### Issue: "products" directory not found

**Symptom:** Error creating products/MOCK-001

**Solution:** Ensure `recursive: true` so parent directory is created

---

## Next Steps

Once Phase 8.3 is complete and code-reviewed:

→ **Phase 8.4**: Create API Fetch bento (calls mock Figma API, downloads overlay images)

---

## Claude Code Prompt Template

```
I need to implement Phase 8.3: Folder Setup Bento using TDD.

Please read:
- .claude/strategy/phase-8.3-folder-setup.md (this file)
- .claude/strategy/phase-8.2-csv-reader.md (CSV reader reference)
- .claude/strategy/phase-8.1-test-infrastructure.md (test helpers)

Then follow TDD workflow:

1. Write failing test (RED):
   - tests/integration/folder_setup_test.go
   - Test: Creates folders for 3 products
   - Test: Idempotent (handles existing folders)
   - Run: go test ./tests/integration -v -run TestFolderSetup (should FAIL)

2. Create bento (RED → GREEN):
   - examples/phase8/folder-setup.bento.json
   - Use spreadsheet neta to read CSV
   - Use loop neta with mode: "forEach"
   - Inside loop, use file-system neta with operation: "mkdir"
   - Connect CSV output to loop input with edge
   - Run: go test ./tests/integration -v -run TestFolderSetup (should PASS)

3. Verify manually:
   - rm -rf products/
   - export INPUT_CSV=tests/fixtures/products-test.csv
   - bento savor examples/phase8/folder-setup.bento.json
   - ls -la products/  # Should show 3 folders

4. Refactor (GREEN → Better):
   - Add metadata explaining loop structure
   - Document template variable usage
   - Ensure clear naming

5. Run /code-review and get Karen + Colossus approval

This demonstrates loops and context passing!
```

---

**Status:** Ready for implementation
**Estimated Time:** 1 hour
**Previous Phase:** 8.2 - CSV Reader
**Next Phase:** 8.4 - API Fetch Bento
