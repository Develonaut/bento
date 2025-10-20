# Phase 8.7: Master Integration Bento

**Duration:** 3 hours
**Goal:** Combine all 5 focused bentos into end-to-end product automation workflow
**Dependencies:** Phases 8.1-8.6 all complete

---

## Overview

**THE FINALE!** This phase combines all focused bentos (8.2-8.6) into one master workflow that processes multiple products end-to-end:

1. Read products from CSV (8.2)
2. Loop through each product:
   - Create folder structure (8.3)
   - Fetch overlay from mock Figma API (8.4)
   - Run mock Blender render (8.5)
   - Optimize to WebP in parallel (8.6)

**This validates EVERYTHING:**
- All 10 neta types
- Context passing between neta
- Looping with real data
- HTTP requests
- Long-running shell commands with streaming
- Parallel execution
- Error handling

---

## Prerequisites

- ✅ Phase 8.1 complete (mocks & fixtures)
- ✅ Phase 8.2 complete (CSV reader)
- ✅ Phase 8.3 complete (folder setup)
- ✅ Phase 8.4 complete (API fetch)
- ✅ Phase 8.5 complete (mock render)
- ✅ Phase 8.6 complete (image optimize)

---

## Master Bento Specification

**File:** `examples/phase8/product-automation.bento.json`

```json
{
  "id": "product-automation",
  "type": "group",
  "version": "1.0.0",
  "name": "Product Photo Automation",
  "metadata": {
    "description": "Complete workflow: CSV → Folders → Figma → Blender → WebP",
    "tags": ["phase8", "integration", "production", "automation"]
  },
  "nodes": [
    {
      "id": "read-products",
      "type": "spreadsheet",
      "version": "1.0.0",
      "name": "Read Products CSV",
      "parameters": {
        "operation": "read",
        "format": "csv",
        "path": "{{.INPUT_CSV}}",
        "hasHeaders": true
      },
      "outputPorts": [{"id": "out-1", "name": "rows"}],
      "position": {"x": 100, "y": 100}
    },
    {
      "id": "process-each-product",
      "type": "loop",
      "version": "1.0.0",
      "name": "Process Each Product",
      "parameters": {
        "mode": "forEach",
        "items": "{{.read-products.rows}}",
        "showProgress": true
      },
      "inputPorts": [{"id": "in-1", "name": "items"}],
      "position": {"x": 300, "y": 100},
      "nodes": [
        {
          "id": "create-folder",
          "type": "file-system",
          "version": "1.0.0",
          "name": "Create Product Folder",
          "parameters": {
            "operation": "mkdir",
            "path": "products/{{.item.sku}}",
            "recursive": true
          },
          "position": {"x": 50, "y": 50}
        },
        {
          "id": "call-figma",
          "type": "http-request",
          "version": "1.0.0",
          "name": "Get Figma Overlay URL",
          "parameters": {
            "url": "{{.FIGMA_API_URL}}",
            "method": "GET",
            "headers": {"X-Figma-Token": "{{.FIGMA_API_TOKEN}}"}
          },
          "position": {"x": 50, "y": 150}
        },
        {
          "id": "download-overlay",
          "type": "http-request",
          "version": "1.0.0",
          "name": "Download Overlay",
          "parameters": {
            "url": "{{.call-figma.body.images.test-component}}",
            "method": "GET",
            "saveToFile": "products/{{.item.sku}}/overlay.png"
          },
          "position": {"x": 50, "y": 250}
        },
        {
          "id": "render",
          "type": "shell-command",
          "version": "1.0.0",
          "name": "Run Blender Render",
          "parameters": {
            "command": "{{.BLENDER_MOCK_SCRIPT}}",
            "args": [
              "--", "--sku", "{{.item.sku}}",
              "--overlay", "products/{{.item.sku}}/overlay.png",
              "--output", "products/{{.item.sku}}/render"
            ],
            "timeout": 60,
            "stream": true
          },
          "position": {"x": 50, "y": 350}
        },
        {
          "id": "optimize",
          "type": "parallel",
          "version": "1.0.0",
          "name": "Optimize to WebP",
          "parameters": {"maxConcurrency": 4},
          "position": {"x": 50, "y": 450},
          "nodes": [
            // All 8 image optimization nodes (optimize-1 through optimize-8)
            // Each converts render-N.png to render-N.webp
          ]
        },
        {
          "id": "cleanup",
          "type": "file-system",
          "version": "1.0.0",
          "name": "Delete PNGs",
          "parameters": {
            "operation": "delete",
            "path": "products/{{.item.sku}}/render-*.png"
          },
          "position": {"x": 50, "y": 550}
        }
      ],
      "edges": [
        {"id": "e1", "source": "create-folder", "target": "call-figma"},
        {"id": "e2", "source": "call-figma", "target": "download-overlay"},
        {"id": "e3", "source": "download-overlay", "target": "render"},
        {"id": "e4", "source": "render", "target": "optimize"},
        {"id": "e5", "source": "optimize", "target": "cleanup"}
      ]
    }
  ],
  "edges": [
    {"id": "main", "source": "read-products", "target": "process-each-product"}
  ],
  "inputPorts": [],
  "outputPorts": [],
  "position": {"x": 0, "y": 0}
}
```

---

## TDD Workflow

### Step 1: Write Integration Test (RED)

**File:** `tests/integration/product_automation_test.go`

```go
func TestProductAutomation_EndToEnd(t *testing.T) {
	// This is THE test - validates everything!

	// Setup
	outputDir := filepath.Join(os.TempDir(), "bento-test-automation")
	defer CleanupTestDir(t, outputDir)

	originalDir, _ := os.Getwd()
	os.MkdirAll(outputDir, 0755)
	os.Chdir(outputDir)
	defer os.Chdir(originalDir)

	// Start mock Figma server
	figmaServer := mocks.NewFigmaServer()
	defer figmaServer.Close()

	testCSV := filepath.Join(originalDir, "../fixtures/products-test.csv")
	blenderScript := filepath.Join(originalDir, "../mocks/blender-mock.sh")

	envVars := map[string]string{
		"INPUT_CSV":           testCSV,
		"FIGMA_API_URL":       figmaServer.URL,
		"FIGMA_API_TOKEN":     "test-token",
		"BLENDER_MOCK_SCRIPT": blenderScript,
	}

	// Execute master bento
	bentoPath := filepath.Join(originalDir, "../../examples/phase8/product-automation.bento.json")
	output, err := RunBento(t, bentoPath, envVars)

	// Verify success
	require.NoError(t, err, "Product automation should complete successfully")

	// Verify all 3 products processed
	assert.Contains(t, output, "MOCK-001")
	assert.Contains(t, output, "MOCK-002")
	assert.Contains(t, output, "MOCK-003")

	// Verify folder structure
	for _, sku := range []string{"MOCK-001", "MOCK-002", "MOCK-003"} {
		productDir := filepath.Join("products", sku)

		// Folder exists
		VerifyFileExists(t, productDir)

		// Overlay downloaded
		VerifyFileExists(t, filepath.Join(productDir, "overlay.png"))

		// 8 WebP files created
		VerifyFileCount(t, productDir, "*.webp", 8)

		// PNGs cleaned up
		VerifyFileCount(t, productDir, "*.png", 1) // Only overlay.png remains
	}

	// Verify streaming output appeared
	assert.Contains(t, output, "Rendering 1/8", "Should show Blender progress")

	// Verify performance (should complete in < 30s)
	// (Can add timing if needed)
}

func TestProductAutomation_ErrorHandling(t *testing.T) {
	// Test missing Figma token
	// Test Blender script failure
	// Test missing CSV file
}
```

### Step 2: Create Master Bento (GREEN)

Combine all the pieces from phases 8.2-8.6 into one bento.

**Key:** This is mostly copy-paste from the focused bentos, but:
- All nodes go inside the loop
- Edges connect them sequentially
- Context variables use `{{.item.sku}}` inside loop

### Step 3: Run End-to-End Test

```bash
go test ./tests/integration -v -run TestProductAutomation
```

**Expected:**
- Test takes ~8-10 seconds (3 products × ~3s each)
- All assertions pass
- 24 WebP files created (8 per product × 3 products)

### Step 4: Manual Verification

```bash
# Clean up
rm -rf products/

# Start mock Figma server (or use one from test)
# Set environment variables
export INPUT_CSV=tests/fixtures/products-test.csv
export FIGMA_API_URL=http://localhost:9999
export FIGMA_API_TOKEN=test-token
export BLENDER_MOCK_SCRIPT=tests/mocks/blender-mock.sh

# Run the master bento!
bento run examples/phase8/product-automation.bento.json

# Watch the magic happen:
# - Reads CSV
# - Processing product 1/3: MOCK-001
# - Creating folder...
# - Calling Figma API...
# - Downloading overlay...
# - Rendering with Blender...
#   Fra:1 Rendering 1/8
#   Fra:2 Rendering 2/8
#   ...
# - Optimizing 8 images...
# - Cleaning up PNGs...
# - Processing product 2/3: MOCK-002
# ...

# Verify results
tree products/
# products/
# ├── MOCK-001/
# │   ├── overlay.png
# │   ├── render-1.webp
# │   ├── render-2.webp
# │   ... (8 webp files)
# ├── MOCK-002/
# │   ... (same)
# └── MOCK-003/
#     ... (same)
```

---

## Success Criteria

- [x] Integration test passes
- [x] Processes 3 products end-to-end
- [x] Creates 24 WebP images (8 per product × 3)
- [x] Shows streaming progress from Blender
- [x] Completes in < 30 seconds
- [x] All folder structure correct
- [x] PNGs cleaned up (only overlay.png remains)
- [x] Manual test works
- [x] Error handling tested
- [x] Performance documented
- [x] **FINAL CODE REVIEW** approved by Karen + Colossus

---

## Performance Benchmarks

Document in bento metadata or separate file:

```markdown
# Performance - Product Automation

**Test Date:** 2025-10-19
**System:** MacBook Pro M1
**Products:** 3

| Stage | Time per Product | Total (3 products) |
|-------|------------------|---------------------|
| CSV Read | 10ms | 30ms |
| Folder Creation | 5ms | 15ms |
| Figma API | 50ms | 150ms |
| Overlay Download | 100ms | 300ms |
| Blender Render | 1.6s | 4.8s |
| WebP Optimization | 400ms | 1.2s |
| PNG Cleanup | 10ms | 30ms |
| **TOTAL** | **~2.2s** | **~6.5s** |

**Observations:**
- Mock Blender much faster than real (1.6s vs 3-5min)
- Parallel optimization saves time (4 concurrent)
- Context passing works perfectly
- Streaming output works!
```

---

## Celebration Checklist

Phase 8 is COMPLETE when:

- [x] All 7 sub-phases (8.1-8.7) done
- [x] All tests pass (RED → GREEN → REFACTOR)
- [x] Integration test passes
- [x] Manual end-to-end test works
- [x] Performance benchmarked
- [x] Final code review approved
- [x] Phase 8 strategy document removed (mark complete!)

**THEN: CELEBRATE! 🍱🎉**

Phase 8 proves the entire bento system works end-to-end!

---

## Code Review Command

```bash
# Review the ENTIRE Phase 8
# /code-review examples/phase8/*.bento.json tests/integration/*_test.go tests/mocks/*
```

**What reviewers check:**
- All bentos follow best practices
- Tests are comprehensive
- Error handling works
- Performance is documented
- Context passing works correctly
- Streaming output works
- Code follows Bento Box Principle
- Ready for production use!

---

## Next Steps

**After Phase 8 Complete:**

1. Archive Phase 8 strategy doc:
   - Move to `.claude/completed/phase-8.md`
   - Update main README with Phase 8 completion

2. Production Readiness:
   - Document how to use real Figma API (not mock)
   - Document how to use real Blender (not mock)
   - Create user guide for product automation

3. Future Enhancements:
   - Resume capability (save progress)
   - Parallel Blender (multiple products at once)
   - Cloud rendering
   - Error recovery (`--continue-on-error`)
   - Notifications (Slack, email)

---

## Claude Code Prompt Template

```
Implement Phase 8.7: Master Integration Bento - THE FINALE!

Please read:
- .claude/strategy/phase-8-real-world-bento.md (overall Phase 8 context)
- .claude/BENTO_BOX_PRINCIPLE.md (coding standards)
- .claude/strategy/phase-8.1-test-infrastructure.md (test infrastructure)
- .claude/strategy/phase-8.2-csv-reader.md (CSV reader)
- .claude/strategy/phase-8.3-folder-setup.md (folder setup)
- .claude/strategy/phase-8.4-api-fetch.md (API fetch)
- .claude/strategy/phase-8.5-mock-render.md (mock render)
- .claude/strategy/phase-8.6-image-optimize.md (image optimization)
- .claude/strategy/phase-8.7-master-integration.md (this phase)

TDD Workflow:

1. Write integration test (RED):
   - tests/integration/product_automation_test.go
   - Test: End-to-end processes 3 products
   - Test: Creates 24 WebP files (8 × 3)
   - Test: Streaming output works
   - Test: Error handling
   - Run: go test ./tests/integration -v -run TestProductAutomation (FAIL)

2. Create master bento (GREEN):
   - examples/phase8/product-automation.bento.json
   - Combine nodes from phases 8.2-8.6
   - Put all processing inside loop neta
   - Connect with edges
   - Run: go test ./tests/integration -v -run TestProductAutomation (PASS!)

3. Verify manually:
   - Set all environment variables
   - Run: bento run examples/phase8/product-automation.bento.json
   - Watch streaming output
   - Verify: tree products/ shows correct structure
   - Verify: 24 WebP files exist

4. Document performance:
   - Time the execution
   - Note any bottlenecks
   - Add to bento metadata

5. Run FINAL /code-review:
   - Review ALL Phase 8 code
   - Get Karen + Colossus approval

6. CELEBRATE! Phase 8 complete! 🍱🎉

This is it - the complete product automation workflow!
All 10 neta types working together. THE BENTO SYSTEM WORKS!
```

---

**Status:** Ready for implementation
**Estimated Time:** 3 hours
**Previous Phase:** 8.6 - Image Optimization
**Next Step:** PRODUCTION READY! 🚀
