# Phase 8: Real-World Integration Test - Product Photo Automation

**Duration:** 1 week
**Goal:** Build and test a complete real-world bento that automates product photo generation
**Dependencies:** ALL phases (1-7) must be complete

---

## Overview

Phase 8 is the ultimate validation that the entire bento system works end-to-end. We'll build a real bento that automates your product photo generation workflow:

1. Read product data from CSV
2. For each product row:
   - Create folder structure (`products/[SKU]/`)
   - Call Figma API to generate overlay.png
   - Execute Blender Python script to composite overlay on render (2-5 min)
   - Output 8 product photos
   - Optimize all 8 photos to webp
   - Save to product folder

This workflow exercises **every critical feature** of the bento system:
- CSV parsing (spreadsheet neta)
- Looping through items (loop neta)
- HTTP API calls (http-request neta)
- File operations (file-system neta)
- Long-running shell commands (shell-command neta for Blender)
- Image optimization (image neta)
- Data transformation (transform neta)
- Orchestration (itamae)
- Error handling
- Progress tracking

---

## Success Criteria

**Phase 8 Complete When:**
- [ ] `product-automation.bento.json` created and tested
- [ ] Successfully processes at least 3 test products end-to-end
- [ ] Blender renders complete (2-5 minutes each)
- [ ] All 8 photos generated and optimized to webp
- [ ] Proper folder structure created
- [ ] Progress shown ("Rendering product 3/50... 45%")
- [ ] Errors handled gracefully (Figma API failure, Blender crash)
- [ ] Performance benchmarks documented
- [ ] Integration test written and passing
- [ ] `/code-review` run with Karen + Colossus approval

---

## The Product Automation Workflow

### Input: products.csv

```csv
sku,name,description,category
PROD-001,Widget A,High-quality widget,Widgets
PROD-002,Gadget B,Premium gadget,Gadgets
PROD-003,Doohickey C,Essential doohickey,Tools
```

### Output Structure

```
products/
├── PROD-001/
│   ├── overlay.png          # From Figma API
│   ├── render-1.webp        # Blender output, optimized
│   ├── render-2.webp
│   ├── render-3.webp
│   ├── render-4.webp
│   ├── render-5.webp
│   ├── render-6.webp
│   ├── render-7.webp
│   └── render-8.webp
├── PROD-002/
│   └── ...
└── PROD-003/
    └── ...
```

---

## Bento Definition

**File: `examples/product-automation.bento.json`**

```json
{
  "id": "product-photo-automation",
  "type": "group",
  "version": "1.0.0",
  "name": "Product Photo Automation",
  "metadata": {
    "description": "Automates product photo generation: CSV → Figma → Blender → WebP",
    "author": "Ryan",
    "tags": ["production", "photos", "blender", "figma"]
  },
  "nodes": [
    {
      "id": "read-products-csv",
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
      ]
    },
    {
      "id": "process-each-product",
      "type": "loop",
      "version": "1.0.0",
      "name": "Process Each Product",
      "parameters": {
        "mode": "forEach",
        "items": "{{.read-products-csv.rows}}",
        "showProgress": true
      },
      "inputPorts": [
        {"id": "in-1", "name": "items"}
      ],
      "nodes": [
        {
          "id": "create-product-folder",
          "type": "file-system",
          "version": "1.0.0",
          "name": "Create Product Folder",
          "parameters": {
            "operation": "mkdir",
            "path": "products/{{.item.sku}}",
            "recursive": true
          }
        },
        {
          "id": "generate-figma-overlay",
          "type": "http-request",
          "version": "1.0.0",
          "name": "Generate Figma Overlay",
          "parameters": {
            "url": "https://api.figma.com/v1/images/{{.FIGMA_FILE_ID}}",
            "method": "GET",
            "headers": {
              "X-Figma-Token": "{{.FIGMA_API_TOKEN}}"
            },
            "queryParams": {
              "ids": "{{.FIGMA_COMPONENT_ID}}",
              "format": "png",
              "scale": 2
            },
            "timeout": 30
          }
        },
        {
          "id": "download-overlay",
          "type": "http-request",
          "version": "1.0.0",
          "name": "Download Overlay Image",
          "parameters": {
            "url": "{{.generate-figma-overlay.images[FIGMA_COMPONENT_ID]}}",
            "method": "GET",
            "saveToFile": "products/{{.item.sku}}/overlay.png",
            "timeout": 60
          }
        },
        {
          "id": "render-with-blender",
          "type": "shell-command",
          "version": "1.0.0",
          "name": "Render Product Photos with Blender",
          "parameters": {
            "command": "blender",
            "args": [
              "--background",
              "{{.BLENDER_FILE}}",
              "--python", "{{.BLENDER_SCRIPT}}",
              "--",
              "--sku", "{{.item.sku}}",
              "--overlay", "products/{{.item.sku}}/overlay.png",
              "--output", "products/{{.item.sku}}/render"
            ],
            "timeout": 600,
            "stream": true,
            "workingDir": "."
          }
        },
        {
          "id": "optimize-to-webp",
          "type": "parallel",
          "version": "1.0.0",
          "name": "Optimize 8 Renders to WebP",
          "parameters": {
            "maxConcurrency": 4
          },
          "nodes": [
            {
              "id": "optimize-1",
              "type": "image",
              "version": "1.0.0",
              "name": "Optimize Render 1",
              "parameters": {
                "operation": "convert",
                "input": "products/{{.item.sku}}/render-1.png",
                "output": "products/{{.item.sku}}/render-1.webp",
                "format": "webp",
                "quality": 85
              }
            },
            {
              "id": "optimize-2",
              "type": "image",
              "version": "1.0.0",
              "parameters": {
                "operation": "convert",
                "input": "products/{{.item.sku}}/render-2.png",
                "output": "products/{{.item.sku}}/render-2.webp",
                "format": "webp",
                "quality": 85
              }
            },
            {
              "id": "optimize-3",
              "type": "image",
              "version": "1.0.0",
              "parameters": {
                "operation": "convert",
                "input": "products/{{.item.sku}}/render-3.png",
                "output": "products/{{.item.sku}}/render-3.webp",
                "format": "webp",
                "quality": 85
              }
            },
            {
              "id": "optimize-4",
              "type": "image",
              "version": "1.0.0",
              "parameters": {
                "operation": "convert",
                "input": "products/{{.item.sku}}/render-4.png",
                "output": "products/{{.item.sku}}/render-4.webp",
                "format": "webp",
                "quality": 85
              }
            },
            {
              "id": "optimize-5",
              "type": "image",
              "version": "1.0.0",
              "parameters": {
                "operation": "convert",
                "input": "products/{{.item.sku}}/render-5.png",
                "output": "products/{{.item.sku}}/render-5.webp",
                "format": "webp",
                "quality": 85
              }
            },
            {
              "id": "optimize-6",
              "type": "image",
              "version": "1.0.0",
              "parameters": {
                "operation": "convert",
                "input": "products/{{.item.sku}}/render-6.png",
                "output": "products/{{.item.sku}}/render-6.webp",
                "format": "webp",
                "quality": 85
              }
            },
            {
              "id": "optimize-7",
              "type": "image",
              "version": "1.0.0",
              "parameters": {
                "operation": "convert",
                "input": "products/{{.item.sku}}/render-7.png",
                "output": "products/{{.item.sku}}/render-7.webp",
                "format": "webp",
                "quality": 85
              }
            },
            {
              "id": "optimize-8",
              "type": "image",
              "version": "1.0.0",
              "parameters": {
                "operation": "convert",
                "input": "products/{{.item.sku}}/render-8.png",
                "output": "products/{{.item.sku}}/render-8.webp",
                "format": "webp",
                "quality": 85
              }
            }
          ]
        },
        {
          "id": "cleanup-pngs",
          "type": "file-system",
          "version": "1.0.0",
          "name": "Delete Original PNG Files",
          "parameters": {
            "operation": "delete",
            "path": "products/{{.item.sku}}/render-*.png"
          }
        }
      ],
      "edges": [
        {"id": "e1", "source": "create-product-folder", "target": "generate-figma-overlay"},
        {"id": "e2", "source": "generate-figma-overlay", "target": "download-overlay"},
        {"id": "e3", "source": "download-overlay", "target": "render-with-blender"},
        {"id": "e4", "source": "render-with-blender", "target": "optimize-to-webp"},
        {"id": "e5", "source": "optimize-to-webp", "target": "cleanup-pngs"}
      ]
    }
  ],
  "edges": [
    {"id": "main-edge", "source": "read-products-csv", "target": "process-each-product"}
  ],
  "inputPorts": [],
  "outputPorts": [],
  "position": {"x": 0, "y": 0},
  "parameters": {}
}
```

---

## Environment Setup

### Required Environment Variables

Create `.env` file:

```bash
# Figma API
FIGMA_API_TOKEN=your-figma-api-token
FIGMA_FILE_ID=your-figma-file-id
FIGMA_COMPONENT_ID=your-component-id

# Input/Output
INPUT_CSV=./products.csv
OUTPUT_DIR=./products

# Blender
BLENDER_FILE=./blender/product-template.blend
BLENDER_SCRIPT=./blender/render-product.py
```

### Blender Python Script

Create `blender/render-product.py`:

```python
import bpy
import sys
import os

# Parse command-line arguments
argv = sys.argv
argv = argv[argv.index("--") + 1:]  # Get args after "--"

sku = None
overlay_path = None
output_path = None

for i, arg in enumerate(argv):
    if arg == "--sku":
        sku = argv[i + 1]
    elif arg == "--overlay":
        overlay_path = argv[i + 1]
    elif arg == "--output":
        output_path = argv[i + 1]

print(f"Rendering product: {sku}")
print(f"Overlay: {overlay_path}")
print(f"Output: {output_path}")

# Load overlay into Blender scene
# (Your existing Blender setup code here)

# Render 8 angles
for i in range(1, 9):
    # Rotate camera to angle i
    # (Your rotation logic here)

    # Set output path
    bpy.context.scene.render.filepath = f"{output_path}-{i}.png"

    # Render
    print(f"Rendering angle {i}/8...")
    bpy.ops.render.render(write_still=True)

print(f"✓ Rendered 8 photos for {sku}")
```

---

## Test Strategy

### Integration Test

Create `tests/integration/product_automation_test.go`:

```go
package integration_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestProductAutomation_EndToEnd(t *testing.T) {
	// Skip if Blender not installed
	if _, err := exec.LookPath("blender"); err != nil {
		t.Skip("Blender not installed, skipping integration test")
	}

	// Create test CSV
	testCSV := createTestCSV(t)
	defer os.Remove(testCSV)

	// Set environment variables
	os.Setenv("INPUT_CSV", testCSV)
	os.Setenv("FIGMA_API_TOKEN", os.Getenv("TEST_FIGMA_TOKEN"))
	os.Setenv("FIGMA_FILE_ID", "test-file")
	os.Setenv("FIGMA_COMPONENT_ID", "test-component")

	// Run bento
	cmd := exec.Command("bento", "serve", "product-automation.bento.json", "--timeout", "30m")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Bento execution failed: %v\nOutput: %s", err, string(output))
	}

	// Verify outputs
	verifyProductOutput(t, "PROD-001")
	verifyProductOutput(t, "PROD-002")
	verifyProductOutput(t, "PROD-003")
}

func verifyProductOutput(t *testing.T, sku string) {
	productDir := filepath.Join("products", sku)

	// Check folder exists
	if _, err := os.Stat(productDir); os.IsNotExist(err) {
		t.Errorf("Product folder not created: %s", productDir)
		return
	}

	// Check overlay exists
	overlayPath := filepath.Join(productDir, "overlay.png")
	if _, err := os.Stat(overlayPath); os.IsNotExist(err) {
		t.Errorf("Overlay not created: %s", overlayPath)
	}

	// Check all 8 webp files exist
	for i := 1; i <= 8; i++ {
		webpPath := filepath.Join(productDir, fmt.Sprintf("render-%d.webp", i))
		if _, err := os.Stat(webpPath); os.IsNotExist(err) {
			t.Errorf("Render %d not created: %s", i, webpPath)
		}
	}

	// Check PNGs were deleted
	pngPath := filepath.Join(productDir, "render-1.png")
	if _, err := os.Stat(pngPath); !os.IsNotExist(err) {
		t.Error("PNG files should be deleted after webp conversion")
	}
}

func createTestCSV(t *testing.T) string {
	content := `sku,name,description,category
PROD-001,Test Widget A,Test product A,Widgets
PROD-002,Test Gadget B,Test product B,Gadgets
PROD-003,Test Tool C,Test product C,Tools`

	tmpfile, err := os.CreateTemp("", "test-products-*.csv")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}

	tmpfile.Close()
	return tmpfile.Name()
}
```

---

## Expected Output

```bash
$ bento serve product-automation.bento.json --timeout 30m

🍱 Serving bento: Product Photo Automation
🍙 Executing neta 'read-products-csv' (spreadsheet)
✓ Completed in 45ms - Read 50 products

🍙 Executing neta 'process-each-product' (loop)
  ⟳ Processing product 1/50: PROD-001
    🍙 Creating folder: products/PROD-001
    ✓ Completed in 12ms

    🍙 Calling Figma API for overlay
    ✓ Completed in 1.2s

    🍙 Downloading overlay image
    ✓ Completed in 847ms

    🍙 Rendering with Blender (this may take a few minutes)
      Fra:1 Mem:12.00M (Peak 12.00M) | Rendering 1/8
      Fra:2 Mem:12.00M (Peak 12.00M) | Rendering 2/8
      Fra:3 Mem:12.00M (Peak 12.00M) | Rendering 3/8
      Fra:4 Mem:12.00M (Peak 12.00M) | Rendering 4/8
      Fra:5 Mem:12.00M (Peak 12.00M) | Rendering 5/8
      Fra:6 Mem:12.00M (Peak 12.00M) | Rendering 6/8
      Fra:7 Mem:12.00M (Peak 12.00M) | Rendering 7/8
      Fra:8 Mem:12.00M (Peak 12.00M) | Rendering 8/8
    ✓ Completed in 3m 42s

    🍙 Optimizing 8 renders to WebP (parallel)
    ✓ Completed in 2.1s

    🍙 Cleaning up PNG files
    ✓ Completed in 134ms

  ✓ Product 1/50 complete in 3m 47s

  ⟳ Processing product 2/50: PROD-002
    ...

🍱 Bento served successfully in 3h 8m 42s
   302 neta executed
   50 products processed
   400 images generated
```

---

## Performance Benchmarks

Document performance in `BENCHMARKS.md`:

```markdown
# Performance Benchmarks

## Product Photo Automation

**Hardware:** M1 MacBook Pro, 16GB RAM
**Test Date:** 2025-10-19
**Products:** 50

| Stage | Time per Product | Total (50 products) |
|-------|------------------|---------------------|
| CSV Read | 1ms | 50ms |
| Folder Creation | 10ms | 500ms |
| Figma API | 1.2s | 1m |
| Overlay Download | 800ms | 40s |
| Blender Render | 3m 40s | 3h 3m |
| WebP Optimization | 2s | 1m 40s |
| PNG Cleanup | 100ms | 5s |
| **TOTAL** | **3m 47s** | **3h 8m** |

## Observations

- Blender rendering is the bottleneck (96% of total time)
- Parallel WebP optimization saves ~6s per product (8 images × 750ms sequential)
- Memory usage stays below 500MB throughout
- Could parallelize Blender renders if multiple licenses available
```

---

## Error Handling Scenarios

Test these error scenarios:

1. **Missing Figma API token**
   - Should fail fast at validation with clear error

2. **Figma API rate limit**
   - Should retry with exponential backoff

3. **Blender crash on specific product**
   - Should log error, skip product, continue to next

4. **Out of disk space**
   - Should fail gracefully with clear error

5. **Invalid CSV format**
   - Should fail at validation with row/column info

---

## Improvements for Future

Document potential improvements:

1. **Resume capability**: Save progress, resume from last successful product
2. **Parallel Blender**: Render multiple products simultaneously (if licenses permit)
3. **Cloud rendering**: Offload Blender to cloud instances
4. **Error recovery**: `--continue-on-error` flag to skip failed products
5. **Dry run**: `bento serve --dry-run` to validate without executing
6. **Notifications**: Slack/email notification when complete

---

## Phase Completion

**Phase 8 MUST end with:**

1. Bento file created and documented
2. Integration test passing
3. Successfully processed at least 3 real products end-to-end
4. Performance benchmarks documented
5. Error handling tested
6. Run `/code-review` slash command
7. Get approval from Karen and Colossus
8. Celebrate! 🍱🎉

---

## Claude Prompt Template

```
I need to implement Phase 8: Real-World Integration Test.

This is the FINAL phase - validating everything works with a real automation workflow!

Please read:
- .claude/strategy/phase-8-real-world-bento.md (this file)
- All previous phase documents (1-7)

Then:

1. Create `examples/product-automation.bento.json` from the template in phase-8 doc

2. Create integration test in `tests/integration/product_automation_test.go`

3. Set up test environment:
   - Create test CSV with 3 products
   - Set environment variables
   - Create mock Blender script (or use real one if available)

4. Run the bento end-to-end:
   - bento inspect product-automation.bento.json (should pass)
   - bento serve product-automation.bento.json --timeout 30m

5. Verify outputs:
   - All folders created
   - Overlay images downloaded
   - 8 webp files per product
   - PNGs cleaned up

6. Document performance benchmarks

7. Test error scenarios (missing Figma token, Blender crash, etc.)

This validates:
- All 10 neta types work correctly
- Context passing between neta
- Loop execution through CSV
- Long-running processes (Blender)
- Parallel execution (WebP optimization)
- Error handling
- Progress tracking

When complete, run `/code-review` and celebrate! 🍱
```

---

## Success Metrics

The bento system is production-ready when:

- ✅ All 50 products processed successfully
- ✅ 400 webp images generated (8 per product)
- ✅ Total time < 4 hours (automatic overnight processing acceptable)
- ✅ Memory usage < 1GB peak
- ✅ Clear progress shown throughout
- ✅ Errors handled gracefully
- ✅ Easy to resume if interrupted
- ✅ Integration test passes reliably

---

**Status:** Ready for final implementation and testing!
**This is it:** The culmination of all 7 phases. Let's automate your workflow! 🍱🚀
