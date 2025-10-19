# Phase 8.5: Mock Render Bento

**Duration:** 2 hours
**Goal:** TDD a bento that runs mock Blender script with streaming progress output
**Dependencies:** Phase 8.1 (mock Blender script) complete

---

## Overview

**CRITICAL PHASE**: Validates streaming output from long-running shell commands. This is the most important feature for real Blender renders.

Uses mock Blender script from Phase 8.1 that outputs progress lines like real Blender.

---

## Prerequisites

- ✅ Phase 8.1 complete (`tests/mocks/blender-mock.sh` exists)
- ✅ Shell-command neta supports `stream: true`
- ✅ Shoyu logger supports streaming callback

---

## Bento Specification

**File:** `examples/phase8/mock-render.bento.json`

```json
{
  "id": "mock-render",
  "type": "group",
  "version": "1.0.0",
  "name": "Mock Blender Render",
  "metadata": {
    "description": "Runs mock Blender script that creates 8 product renders with streaming progress",
    "tags": ["phase8", "blender", "render", "streaming"]
  },
  "nodes": [
    {
      "id": "render-product",
      "type": "shell-command",
      "version": "1.0.0",
      "name": "Run Mock Blender",
      "parameters": {
        "command": "{{.BLENDER_MOCK_SCRIPT}}",
        "args": [
          "--",
          "--sku", "{{.PRODUCT_SKU}}",
          "--overlay", "products/{{.PRODUCT_SKU}}/overlay.png",
          "--output", "products/{{.PRODUCT_SKU}}/render"
        ],
        "timeout": 60,
        "stream": true,
        "workingDir": "."
      },
      "position": {"x": 100, "y": 100},
      "metadata": {
        "description": "Streams output: 'Rendering 1/8... 2/8...' Shows real-time progress."
      }
    }
  ],
  "edges": [],
  "inputPorts": [],
  "outputPorts": [],
  "position": {"x": 0, "y": 0}
}
```

---

## TDD Workflow

### Step 1: Write Test (RED)

**File:** `tests/integration/mock_render_test.go`

```go
func TestMockRender_CreatesPNGs(t *testing.T) {
	// Setup
	os.MkdirAll("products/MOCK-001", 0755)
	defer os.RemoveAll("products")

	// Create dummy overlay.png
	os.WriteFile("products/MOCK-001/overlay.png", []byte("fake"), 0644)

	scriptPath, _ := filepath.Abs("../mocks/blender-mock.sh")

	envVars := map[string]string{
		"BLENDER_MOCK_SCRIPT": scriptPath,
		"PRODUCT_SKU":         "MOCK-001",
	}

	output, err := RunBento(t, "../../examples/phase8/mock-render.bento.json", envVars)
	require.NoError(t, err)

	// Verify streaming output appeared
	assert.Contains(t, output, "Rendering 1/8", "Should show progress line 1")
	assert.Contains(t, output, "Rendering 8/8", "Should show progress line 8")

	// Verify 8 PNGs created
	VerifyFileCount(t, "products/MOCK-001", "render-*.png", 8)
}

func TestMockRender_StreamingProgress(t *testing.T) {
	// Similar test but focus on streaming behavior
	// Verify we see incremental output, not all at once
}
```

### Step 2: Create Bento (GREEN)

Create bento file with `stream: true` parameter.

**Key:** The `stream: true` flag tells shell-command neta to output lines in real-time instead of buffering.

### Step 3: Verify Manually

```bash
# Make sure mock script exists and is executable
chmod +x tests/mocks/blender-mock.sh

# Create test structure
mkdir -p products/MOCK-001
touch products/MOCK-001/overlay.png

# Run bento
export BLENDER_MOCK_SCRIPT=tests/mocks/blender-mock.sh
export PRODUCT_SKU=MOCK-001

bento savor examples/phase8/mock-render.bento.json

# You should see:
# Fra:1 Mem:12.00M (Peak 12.00M) | Rendering 1/8
# Fra:2 Mem:12.00M (Peak 12.00M) | Rendering 2/8
# ...
# Fra:8 Mem:12.00M (Peak 12.00M) | Rendering 8/8

# Verify PNGs created
ls products/MOCK-001/render-*.png | wc -l  # Should be 8
```

---

## Success Criteria

- [x] Bento executes mock Blender script
- [x] Streaming output shows progress lines in real-time
- [x] Creates 8 PNG files (render-1.png through render-8.png)
- [x] Tests verify streaming behavior
- [x] Manual test shows incremental output
- [x] Code review approved (Karen + Colossus check streaming!)

---

## Testing Streaming Output

**Important:** Verify output appears incrementally, not all at once.

The shell-command neta should call the logger's `Stream()` method for each line, which triggers the OnStream callback.

---

## Common Issues

### Issue: No streaming output visible

**Solution:** Check that:
1. `stream: true` is set in bento parameters
2. shell-command neta implements streaming correctly
3. shoyu logger OnStream callback is configured

### Issue: Script not executable

**Solution:**
```bash
chmod +x tests/mocks/blender-mock.sh
```

### Issue: PNGs not created

**Solution:** Check script permissions and working directory

---

## Claude Code Prompt

```
Implement Phase 8.5: Mock Render Bento using TDD.

Read: .claude/strategy/phase-8.5-mock-render.md

TDD Workflow:
1. Verify tests/mocks/blender-mock.sh is executable
2. Write tests/integration/mock_render_test.go (RED)
   - Test: Creates 8 PNGs
   - Test: Shows streaming progress
3. Create examples/phase8/mock-render.bento.json (GREEN)
   - Use shell-command neta with stream: true
4. Verify manually - should see incremental output
5. Run /code-review

CRITICAL: Streaming must work! This is key for real Blender.
```

---

**Next Phase:** 8.6 - Image Optimization (parallel WebP conversion)
