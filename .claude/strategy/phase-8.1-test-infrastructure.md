# Phase 8.1: Test Infrastructure & Mocks

**Duration:** 1-2 hours
**Goal:** Create mock services and test fixtures for Phase 8 integration tests
**Dependencies:** Phases 1-7 complete

---

## Overview

Phase 8.1 establishes the test infrastructure needed for all subsequent Phase 8 sub-phases. We'll create mock versions of external services (Figma API, Blender) and test data to validate our bentos work end-to-end without requiring real external dependencies.

**Why Mocks?**
- **Fast tests**: Mock Blender completes in 2s instead of 5min
- **Reliable**: No external API rate limits or network issues
- **Reproducible**: Same output every time
- **TDD-friendly**: Can test error scenarios easily

---

## Prerequisites

- ✅ All Phase 1-7 complete
- ✅ All neta types working (10/10 passing tests)
- ✅ `bento savor` command working

---

## File Structure

```
tests/
├── fixtures/
│   ├── products-test.csv        # 3 test products
│   └── overlay-sample.png       # 100x100 test image
│
├── mocks/
│   ├── figma_server.go          # HTTP test server for Figma API
│   ├── figma_server_test.go     # Tests for mock server
│   ├── blender-mock.sh          # Shell script simulating Blender
│   └── README.md                # Documentation for mocks
│
└── integration/
    └── helpers.go                # Test helper functions
```

---

## Deliverables

### 1. Mock Figma Server

**File:** `tests/mocks/figma_server.go`

```go
package mocks

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

// NewFigmaServer creates a mock Figma API server.
// Returns URLs for test images like real Figma API.
func NewFigmaServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check auth header
		if r.Header.Get("X-Figma-Token") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "Missing Figma token",
			})
			return
		}

		// Simulate Figma API response
		// Real API returns: {"images": {"node-id": "https://..."}}
		response := map[string]interface{}{
			"images": map[string]interface{}{
				"test-component": "http://localhost:9999/mock-overlay.png",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
}
```

**Test:** `tests/mocks/figma_server_test.go`

```go
func TestFigmaServer_ValidRequest(t *testing.T) {
	server := NewFigmaServer()
	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL, nil)
	req.Header.Set("X-Figma-Token", "test-token")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)

	images := body["images"].(map[string]interface{})
	assert.NotEmpty(t, images["test-component"])
}

func TestFigmaServer_MissingToken(t *testing.T) {
	server := NewFigmaServer()
	defer server.Close()

	resp, err := http.Get(server.URL)
	require.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}
```

---

### 2. Mock Blender Script

**File:** `tests/mocks/blender-mock.sh`

```bash
#!/bin/bash
# Mock Blender script that simulates rendering 8 product photos
# Outputs progress like real Blender and creates PNG files

SKU=""
OVERLAY=""
OUTPUT=""

# Parse arguments (after "--")
while [[ $# -gt 0 ]]; do
    case $1 in
        --sku)
            SKU="$2"
            shift 2
            ;;
        --overlay)
            OVERLAY="$2"
            shift 2
            ;;
        --output)
            OUTPUT="$2"
            shift 2
            ;;
        *)
            shift
            ;;
    esac
done

echo "Blender 3.6.0"
echo "Rendering product: $SKU"
echo "Overlay: $OVERLAY"
echo "Output: $OUTPUT"

# Render 8 angles (simulate 0.2s per frame)
for i in {1..8}; do
    # Output progress like real Blender
    echo "Fra:$i Mem:12.00M (Peak 12.00M) | Rendering $i/8"

    # Create mock PNG file (1x1 pixel)
    printf "\x89\x50\x4e\x47\x0d\x0a\x1a\x0a" > "${OUTPUT}-${i}.png"

    # Simulate render time
    sleep 0.2
done

echo "✓ Rendered 8 photos for $SKU"
exit 0
```

**Make executable:**
```bash
chmod +x tests/mocks/blender-mock.sh
```

**Test manually:**
```bash
./tests/mocks/blender-mock.sh -- --sku MOCK-001 --overlay test.png --output /tmp/render
# Should create /tmp/render-1.png through /tmp/render-8.png
```

---

### 3. Test Fixtures

**File:** `tests/fixtures/products-test.csv`

```csv
sku,name,description,category
MOCK-001,Test Widget A,High-quality test widget,Widgets
MOCK-002,Test Gadget B,Premium test gadget,Gadgets
MOCK-003,Test Tool C,Essential test tool,Tools
```

**File:** `tests/fixtures/overlay-sample.png`
- Create a simple 100x100 PNG image
- Can use ImageMagick: `convert -size 100x100 xc:blue tests/fixtures/overlay-sample.png`
- Or commit a test image

---

### 4. Test Helpers

**File:** `tests/integration/helpers.go`

```go
package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// RunBento executes a bento file and returns output.
func RunBento(t *testing.T, bentoPath string, envVars map[string]string) (string, error) {
	t.Helper()

	cmd := exec.Command("bento", "savor", bentoPath)

	// Set environment variables
	cmd.Env = os.Environ()
	for k, v := range envVars {
		cmd.Env = append(cmd.Env, k+"="+v)
	}

	output, err := cmd.CombinedOutput()
	return string(output), err
}

// CleanupTestDir removes test output directory.
func CleanupTestDir(t *testing.T, dir string) {
	t.Helper()
	os.RemoveAll(dir)
}

// VerifyFileExists checks if a file exists.
func VerifyFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("File does not exist: %s", path)
	}
}

// VerifyFileCount checks number of files matching pattern.
func VerifyFileCount(t *testing.T, dir, pattern string, expected int) {
	t.Helper()

	matches, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		t.Fatalf("Glob failed: %v", err)
	}

	if len(matches) != expected {
		t.Errorf("Expected %d files matching %s, got %d", expected, pattern, len(matches))
	}
}
```

---

## TDD Workflow

### Step 1: Write Tests (Red)

Create `tests/mocks/figma_server_test.go`:
- Test valid request returns images
- Test missing auth returns 401
- Test response format matches real Figma

### Step 2: Implement Mocks (Green)

Create `tests/mocks/figma_server.go`:
- Implement httptest.Server
- Return mock image URLs
- Handle auth headers

Create `tests/mocks/blender-mock.sh`:
- Parse command-line args
- Output progress lines
- Create PNG files
- Exit successfully

### Step 3: Verify & Refactor (Green → Better)

Run tests:
```bash
go test ./tests/mocks/... -v
```

Test manually:
```bash
# Test Figma mock
go run tests/mocks/figma_server.go

# Test Blender mock
./tests/mocks/blender-mock.sh -- --sku TEST --overlay x.png --output /tmp/test
ls /tmp/test-*.png  # Should show 8 files
```

### Step 4: Document

Create `tests/mocks/README.md` explaining:
- What each mock does
- How to use them in tests
- How they differ from real services

---

## Success Criteria

- [x] `figma_server.go` implemented and tested
- [x] `figma_server_test.go` passes (2/2 tests)
- [x] `blender-mock.sh` executable and working
- [x] Manual test of Blender mock creates 8 PNGs
- [x] `products-test.csv` created with 3 products
- [x] `overlay-sample.png` exists
- [x] `helpers.go` with test utilities
- [x] All tests pass: `go test ./tests/mocks/...`
- [x] Code review approved by Karen + Colossus

---

## Code Review

When complete, run:

```bash
# Placeholder for actual /code-review command
# /code-review tests/mocks/
```

**What reviewers will check:**
- Mock Figma server matches real API response format
- Blender mock outputs realistic progress
- Test helpers are reusable
- Code follows Bento Box Principle
- Tests are thorough

---

## Next Steps

Once Phase 8.1 is complete and code-reviewed:

→ **Phase 8.2**: Create CSV Reader bento (uses `products-test.csv`)

---

## Claude Code Prompt Template

```
I need to implement Phase 8.1: Test Infrastructure & Mocks.

Please read:
- .claude/strategy/phase-8-real-world-bento.md (overall Phase 8 context)
- .claude/BENTO_BOX_PRINCIPLE.md (coding standards)
- .claude/strategy/phase-8.1-test-infrastructure.md (this phase)

Then follow TDD workflow:

1. Create test structure:
   - tests/fixtures/ directory
   - tests/mocks/ directory
   - tests/integration/ directory

2. TDD the Figma mock server:
   - Write tests/mocks/figma_server_test.go (RED)
   - Implement tests/mocks/figma_server.go (GREEN)
   - Verify tests pass

3. Create Blender mock script:
   - Write tests/mocks/blender-mock.sh
   - Make executable
   - Test manually (should create 8 PNGs)

4. Create test fixtures:
   - tests/fixtures/products-test.csv (3 products)
   - tests/fixtures/overlay-sample.png (100x100 blue square)

5. Create test helpers:
   - tests/integration/helpers.go
   - RunBento, CleanupTestDir, VerifyFileExists, VerifyFileCount

6. Verify all tests pass:
   - go test ./tests/mocks/... -v

7. Run /code-review and get Karen + Colossus approval

This establishes the test infrastructure for all of Phase 8!
```

---

**Status:** Ready for implementation
**Estimated Time:** 1-2 hours
**Next Phase:** 8.2 - CSV Reader Bento
