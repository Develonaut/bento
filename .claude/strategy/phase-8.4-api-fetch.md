# Phase 8.4: API Fetch Bento

**Duration:** 2 hours
**Goal:** TDD a bento that calls mock Figma API and downloads overlay images
**Dependencies:** Phase 8.1 (mocks) and 8.3 (folders) complete

---

## Overview

Demonstrates HTTP requests to mock Figma API, downloading images, and error handling. Uses mock server from Phase 8.1.

---

## Prerequisites

- ✅ Phase 8.1 complete (mock Figma server)
- ✅ Phase 8.3 complete (product folders exist)
- ✅ HTTP-request neta working

---

## Bento Specification

**File:** `examples/phase8/api-fetch.bento.json`

```json
{
  "id": "api-fetch",
  "type": "group",
  "version": "1.0.0",
  "name": "Figma API Fetch",
  "metadata": {
    "description": "Calls Figma API to get overlay image URL, then downloads it",
    "tags": ["phase8", "http", "figma", "api"]
  },
  "nodes": [
    {
      "id": "call-figma-api",
      "type": "http-request",
      "version": "1.0.0",
      "name": "Get Figma Image URL",
      "parameters": {
        "url": "{{.FIGMA_API_URL}}/v1/images/{{.FIGMA_FILE_ID}}",
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
      },
      "outputPorts": [{"id": "out-1", "name": "response"}],
      "position": {"x": 100, "y": 100}
    },
    {
      "id": "download-image",
      "type": "http-request",
      "version": "1.0.0",
      "name": "Download Overlay Image",
      "parameters": {
        "url": "{{.call-figma-api.body.images.test-component}}",
        "method": "GET",
        "saveToFile": "products/{{.PRODUCT_SKU}}/overlay.png",
        "timeout": 60
      },
      "inputPorts": [{"id": "in-1", "name": "imageUrl"}],
      "position": {"x": 300, "y": 100}
    }
  ],
  "edges": [
    {"id": "e1", "source": "call-figma-api", "target": "download-image"}
  ],
  "inputPorts": [],
  "outputPorts": [],
  "position": {"x": 0, "y": 0}
}
```

---

## TDD Workflow

### Step 1: Write Test (RED)

**File:** `tests/integration/api_fetch_test.go`

```go
func TestAPIFetch_DownloadOverlay(t *testing.T) {
	// Start mock Figma server
	server := mocks.NewFigmaServer()
	defer server.Close()

	// Create test product folder
	os.MkdirAll("products/MOCK-001", 0755)
	defer os.RemoveAll("products")

	envVars := map[string]string{
		"FIGMA_API_URL":      server.URL,
		"FIGMA_API_TOKEN":    "test-token",
		"FIGMA_FILE_ID":      "test-file",
		"FIGMA_COMPONENT_ID": "test-component",
		"PRODUCT_SKU":        "MOCK-001",
	}

	output, err := RunBento(t, "../../examples/phase8/api-fetch.bento.json", envVars)
	require.NoError(t, err)

	// Verify overlay downloaded
	VerifyFileExists(t, "products/MOCK-001/overlay.png")
}

func TestAPIFetch_MissingToken(t *testing.T) {
	server := mocks.NewFigmaServer()
	defer server.Close()

	envVars := map[string]string{
		"FIGMA_API_URL": server.URL,
		// Missing FIGMA_API_TOKEN
	}

	_, err := RunBento(t, "../../examples/phase8/api-fetch.bento.json", envVars)
	assert.Error(t, err, "Should fail without API token")
}
```

### Step 2: Create Bento (GREEN)

Create bento file as specified above.

### Step 3: Verify Manually

```bash
# Start mock server (in Go test or separate terminal)
# Then:
export FIGMA_API_URL=http://localhost:8080
export FIGMA_API_TOKEN=test-token
export FIGMA_FILE_ID=test-file
export FIGMA_COMPONENT_ID=test-component
export PRODUCT_SKU=MOCK-001

mkdir -p products/MOCK-001
bento savor examples/phase8/api-fetch.bento.json

ls products/MOCK-001/overlay.png  # Should exist
```

---

## Success Criteria

- [x] Bento calls mock Figma API
- [x] Downloads overlay.png to product folder
- [x] Handles missing auth token (fails gracefully)
- [x] Tests pass (2/2)
- [x] Code review approved

---

## Claude Code Prompt

```
Implement Phase 8.4: API Fetch Bento using TDD.

Please read:
- .claude/strategy/phase-8-real-world-bento.md (overall Phase 8 context)
- .claude/BENTO_BOX_PRINCIPLE.md (coding standards)
- .claude/strategy/phase-8.1-test-infrastructure.md (mock Figma server)
- .claude/strategy/phase-8.4-api-fetch.md (this phase)

TDD Workflow:
1. Write tests/integration/api_fetch_test.go (RED)
2. Create examples/phase8/api-fetch.bento.json (GREEN)
3. Verify with mock Figma server from Phase 8.1
4. Test error handling (missing token)
5. Run /code-review

This demonstrates HTTP requests and file downloads!
```

---

**Next Phase:** 8.5 - Mock Render (Blender simulation with streaming output)
