# Phase 8.6: Image Optimization Bento

**Duration:** 2 hours
**Goal:** TDD a bento that converts 8 PNGs to WebP in parallel, then cleans up
**Dependencies:** Phase 8.5 (mock render creates PNGs) complete

---

## Overview

Demonstrates parallel execution with the parallel neta. Converts 8 PNG files to WebP concurrently (4 at a time), then deletes the original PNGs.

**Why Parallel?**
- Converting 8 images sequentially: ~6 seconds (750ms each)
- Converting 8 images with concurrency 4: ~1.5 seconds
- 4x speedup!

---

## Prerequisites

- ✅ Phase 8.5 complete (8 PNGs exist in product folder)
- ✅ Parallel neta working
- ✅ Image neta working (govips)
- ✅ File-system neta working (delete operation)

---

## Bento Specification

**File:** `examples/phase8/image-optimize.bento.json`

```json
{
  "id": "image-optimize",
  "type": "group",
  "version": "1.0.0",
  "name": "Image Optimization",
  "metadata": {
    "description": "Converts 8 PNGs to WebP in parallel, then deletes PNGs",
    "tags": ["phase8", "image", "webp", "parallel", "optimization"]
  },
  "nodes": [
    {
      "id": "optimize-parallel",
      "type": "parallel",
      "version": "1.0.0",
      "name": "Convert to WebP (Parallel)",
      "parameters": {
        "maxConcurrency": 4
      },
      "position": {"x": 100, "y": 100},
      "nodes": [
        {
          "id": "optimize-1",
          "type": "image",
          "version": "1.0.0",
          "name": "Optimize Render 1",
          "parameters": {
            "operation": "convert",
            "input": "products/{{.PRODUCT_SKU}}/render-1.png",
            "output": "products/{{.PRODUCT_SKU}}/render-1.webp",
            "format": "webp",
            "quality": 85
          },
          "position": {"x": 50, "y": 50}
        },
        {
          "id": "optimize-2",
          "type": "image",
          "version": "1.0.0",
          "parameters": {
            "operation": "convert",
            "input": "products/{{.PRODUCT_SKU}}/render-2.png",
            "output": "products/{{.PRODUCT_SKU}}/render-2.webp",
            "format": "webp",
            "quality": 85
          },
          "position": {"x": 50, "y": 150}
        }
        // ... repeat for render-3 through render-8 ...
      ],
      "edges": []
    },
    {
      "id": "cleanup-pngs",
      "type": "file-system",
      "version": "1.0.0",
      "name": "Delete Original PNGs",
      "parameters": {
        "operation": "delete",
        "path": "products/{{.PRODUCT_SKU}}/render-*.png"
      },
      "position": {"x": 300, "y": 100}
    }
  ],
  "edges": [
    {"id": "e1", "source": "optimize-parallel", "target": "cleanup-pngs"}
  ],
  "inputPorts": [],
  "outputPorts": [],
  "position": {"x": 0, "y": 0}
}
```

**Note:** Full bento has all 8 optimize nodes (optimize-1 through optimize-8). Abbreviated here for brevity.

---

## TDD Workflow

### Step 1: Write Test (RED)

**File:** `tests/integration/image_optimize_test.go`

```go
func TestImageOptimize_ConvertsToWebP(t *testing.T) {
	// Setup - create 8 test PNGs
	os.MkdirAll("products/MOCK-001", 0755)
	defer os.RemoveAll("products")

	// Create 8 tiny PNG files (1x1 pixel)
	for i := 1; i <= 8; i++ {
		pngPath := fmt.Sprintf("products/MOCK-001/render-%d.png", i)
		// Use govips or simple PNG bytes
		createTestPNG(t, pngPath)
	}

	envVars := map[string]string{
		"PRODUCT_SKU": "MOCK-001",
	}

	output, err := RunBento(t, "../../examples/phase8/image-optimize.bento.json", envVars)
	require.NoError(t, err)

	// Verify 8 WebP files created
	VerifyFileCount(t, "products/MOCK-001", "*.webp", 8)

	// Verify PNGs deleted
	VerifyFileCount(t, "products/MOCK-001", "*.png", 0)
}

func TestImageOptimize_ParallelExecution(t *testing.T) {
	// Test that parallel execution actually happens
	// (Can check timing or logs for concurrency)
}

func createTestPNG(t *testing.T, path string) {
	// Create minimal valid PNG file
	// Can use govips or write PNG bytes directly
}
```

### Step 2: Create Bento (GREEN)

Create bento with:
- Parallel neta containing 8 image neta nodes
- Each image neta converts PNG → WebP
- File-system neta to delete PNGs after conversion

### Step 3: Verify Manually

```bash
# Create test PNGs (using ImageMagick)
mkdir -p products/MOCK-001
for i in {1..8}; do
  convert -size 100x100 xc:blue products/MOCK-001/render-$i.png
done

# Run bento
export PRODUCT_SKU=MOCK-001
bento savor examples/phase8/image-optimize.bento.json

# Verify WebP files exist
ls products/MOCK-001/*.webp | wc -l  # Should be 8

# Verify PNGs deleted
ls products/MOCK-001/*.png 2>/dev/null || echo "No PNGs (correct!)"
```

---

## Success Criteria

- [x] Converts 8 PNGs to WebP
- [x] Parallel execution (4 concurrent)
- [x] Deletes original PNGs
- [x] Tests pass (2/2)
- [x] Manual test works
- [x] Code review approved

---

## Performance Notes

Document timing:
- Sequential (1 at a time): ~6s
- Parallel (4 concurrent): ~1.5s
- Speedup: 4x

This proves parallel neta works!

---

## Claude Code Prompt

```
Implement Phase 8.6: Image Optimization Bento using TDD.

Please read:
- .claude/strategy/phase-8-real-world-bento.md (overall Phase 8 context)
- .claude/BENTO_BOX_PRINCIPLE.md (coding standards)
- .claude/strategy/phase-8.5-mock-render.md (PNGs from render phase)
- .claude/strategy/phase-8.6-image-optimize.md (this phase)

TDD Workflow:
1. Write tests/integration/image_optimize_test.go (RED)
   - Test: Converts 8 PNGs to WebP
   - Test: Deletes original PNGs
2. Create examples/phase8/image-optimize.bento.json (GREEN)
   - Use parallel neta with maxConcurrency: 4
   - 8 image neta nodes inside parallel
   - file-system neta to cleanup PNGs
3. Verify manually with ImageMagick test images
4. Document performance (time the operation)
5. Run /code-review

This demonstrates parallel execution!
```

---

**Next Phase:** 8.7 - Master Integration (combine all 5 bentos!)
