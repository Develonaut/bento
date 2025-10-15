# Phase 10: Real-World Proof-of-Concept - Etsy Product Image Pipeline

**Status**: Pending
**Duration**: 3-4 hours (includes new node types)
**Prerequisites**: Phase 9 complete, Karen approved

## Overview

Build a real-world bento for automating Etsy product image generation using **composable bento architecture**. This validates the core principle: **"A node is a node is a node"** - bentos can contain other bentos, creating reusable, testable components.

This workflow demonstrates:
- **Bento composition** - One bento calling another as a node
- **Real file operations** - CSV reading, folder creation, image conversion
- **External API integration** - Figma API with secure token storage
- **Shell command execution** - Running Blender scripts
- **Complete production workflow** - Not a toy example!

## Pre-Work Checklist

Before starting, you MUST:

1. ✅ Read [BENTO_BOX_PRINCIPLE.md](../BENTO_BOX_PRINCIPLE.md)
2. ✅ Confirm: "I understand the Bento Box Principle and will follow it"
3. ✅ Use TodoWrite to track all tasks
4. ✅ Phases 5.5-9 complete and Karen approved

## User's Workflow

### Business Context
- Etsy shop selling 3D printed products
- Uses Blender + Figma to generate product images
- Processes multiple products from CSV manifest
- Automates repetitive image generation

### Workflow Architecture (Composable Bentos)

```
┌─────────────────────────────────────────────────────────┐
│  etsy-product-pipeline.bento.yaml (MAIN)                │
├─────────────────────────────────────────────────────────┤
│  1. Read manifest.csv (file.csv.read)                   │
│  2. Loop through products (loop.for)                    │
│     ├─ Create folder (file.mkdir)                       │
│     ├─ Execute: generate-figma-image.bento ──────────┐  │
│     ├─ Run Blender script (shell.exec)               │  │
│     └─ Convert to WebP (image.convert)               │  │
└─────────────────────────────────────────────────────────┘
                                                          │
                                                          ▼
┌─────────────────────────────────────────────────────────┐
│  generate-figma-image.bento.yaml (REUSABLE)            │
├─────────────────────────────────────────────────────────┤
│  1. Call Figma API (http + auth token)                 │
│  2. Download exported PNG (http)                       │
│  3. Return file path                                   │
└─────────────────────────────────────────────────────────┘
```

**Key Insight:** The main bento executes the Figma bento as a node using `bento.execute`. This proves **"a node is a node is a node"** - compartments can contain other bento boxes!

### Example CSV Structure

```csv
sku,name,model_file,figma_component_id,color,material
SKU001,Dragon Miniature,models/dragon.stl,123:456,red,PLA
SKU002,Castle Tower,models/tower.stl,123:457,gray,PETG
SKU003,Sword Prop,models/sword.stl,123:458,silver,ABS
```

## Goals

1. **Implement Missing Node Types**
   - CSV reader
   - File system operations
   - Shell command execution
   - Image conversion

2. **Add Secure Configuration**
   - API key storage
   - Settings UI for managing tokens
   - Environment variable support

3. **Build the Real Bento**
   - Use the editor to create the workflow
   - Validate all phases work together
   - Test with actual CSV data

4. **Validate System**
   - Editor supports complex workflows
   - All node types integrate properly
   - Secure config works
   - System is production-ready

## Core Architecture: Composable Bentos

### The `bento.execute` Node Type

**This is the key innovation** - bentos can call other bentos as nodes!

**File**: `pkg/neta/bento/execute.go` (NEW)
**Type**: `bento.execute`

```go
// Package bento provides bento composition
package bento

import (
	"context"
	"fmt"

	"bento/pkg/itamae"
	"bento/pkg/jubako"
	"bento/pkg/neta"
)

// Execute loads and executes another bento
type Execute struct {
	store *jubako.Store
	chef  *itamae.Itamae
}

// NewExecute creates a bento executor
func NewExecute(store *jubako.Store, chef *itamae.Itamae) *Execute {
	return &Execute{
		store: store,
		chef:  chef,
	}
}

// Execute runs another bento as a node
func (e *Execute) Execute(ctx context.Context, params neta.Params) (interface{}, error) {
	bentoName, err := params.GetString("bento")
	if err != nil {
		return nil, fmt.Errorf("bento name required: %w", err)
	}

	// Load the bento definition
	def, err := e.store.Load(bentoName)
	if err != nil {
		return nil, fmt.Errorf("load bento %s: %w", bentoName, err)
	}

	// Extract inputs for the sub-bento
	inputs := params.GetMap("inputs", map[string]interface{}{})

	// Execute the bento (recursive!)
	result, err := e.chef.Execute(ctx, def, inputs)
	if err != nil {
		return nil, fmt.Errorf("execute bento %s: %w", bentoName, err)
	}

	return result.Output, nil
}
```

**YAML Example:**
```yaml
version: "1.0"
type: bento.execute
name: Generate Product Overlay
parameters:
  bento: generate-figma-image
  inputs:
    figma_component_id: "{{item.figma_component_id}}"
    output_path: "output/{{item.name}}/overlay.png"
```

**Why This Matters:**
- ✅ Reusable components - Write once, use everywhere
- ✅ Testable in isolation - Test each bento independently
- ✅ Shareable - Build a library of common operations
- ✅ Natural composition - Complex workflows = small pieces
- ✅ Proves the architecture - "A node is a node is a node"

---

## New Node Types Required

### 1. CSV Reader Node

**File**: `pkg/neta/file/csv.go` (NEW)
**Type**: `file.csv`

```go
// Package file provides file system operations
package file

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"

	"bento/pkg/neta"
)

// CSVReader reads CSV files
type CSVReader struct{}

// Execute reads CSV and returns rows
func (c CSVReader) Execute(ctx context.Context, params neta.Params) (interface{}, error) {
	path, err := params.GetString("path")
	if err != nil {
		return nil, err
	}

	hasHeader := params.GetBool("has_header", true)

	return readCSV(path, hasHeader)
}

// readCSV reads and parses CSV file
func readCSV(path string, hasHeader bool) ([]map[string]interface{}, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read csv: %w", err)
	}

	if len(records) == 0 {
		return []map[string]interface{}{}, nil
	}

	// Extract headers
	var headers []string
	startRow := 0

	if hasHeader {
		headers = records[0]
		startRow = 1
	} else {
		// Generate column names: col0, col1, etc.
		headers = make([]string, len(records[0]))
		for i := range headers {
			headers[i] = fmt.Sprintf("col%d", i)
		}
	}

	// Convert to maps
	rows := make([]map[string]interface{}, 0, len(records)-startRow)
	for _, record := range records[startRow:] {
		row := make(map[string]interface{})
		for i, value := range record {
			if i < len(headers) {
				row[headers[i]] = value
			}
		}
		rows = append(rows, row)
	}

	return rows, nil
}
```

**YAML Example:**
```yaml
version: "1.0"
type: file.csv
name: Read Product Manifest
parameters:
  path: manifest.csv
  has_header: true
```

### 2. Directory Creation Node

**File**: `pkg/neta/file/mkdir.go` (NEW)
**Type**: `file.mkdir`

```go
package file

import (
	"context"
	"fmt"
	"os"

	"bento/pkg/neta"
)

// MkDir creates directories
type MkDir struct{}

// Execute creates directory
func (m MkDir) Execute(ctx context.Context, params neta.Params) (interface{}, error) {
	path, err := params.GetString("path")
	if err != nil {
		return nil, err
	}

	recursive := params.GetBool("recursive", true)
	mode := params.GetInt("mode", 0755)

	return createDir(path, recursive, mode)
}

// createDir creates directory with permissions
func createDir(path string, recursive bool, mode int) (interface{}, error) {
	if recursive {
		if err := os.MkdirAll(path, os.FileMode(mode)); err != nil {
			return nil, fmt.Errorf("create dir: %w", err)
		}
	} else {
		if err := os.Mkdir(path, os.FileMode(mode)); err != nil {
			return nil, fmt.Errorf("create dir: %w", err)
		}
	}

	return map[string]interface{}{
		"path":    path,
		"created": true,
	}, nil
}
```

**YAML Example:**
```yaml
version: "1.0"
type: file.mkdir
name: Create Product Folder
parameters:
  path: "output/{{product_name}}"
  recursive: true
```

### 3. Shell Command Node

**File**: `pkg/neta/shell/command.go` (NEW)
**Type**: `shell.command`

```go
// Package shell provides shell command execution
package shell

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"bento/pkg/neta"
)

// Command executes shell commands
type Command struct{}

// Execute runs shell command
func (c Command) Execute(ctx context.Context, params neta.Params) (interface{}, error) {
	command, err := params.GetString("command")
	if err != nil {
		return nil, err
	}

	args := params.GetStringSlice("args", []string{})
	workDir := params.GetString("working_dir", "")

	return runCommand(ctx, command, args, workDir)
}

// runCommand executes command and returns output
func runCommand(ctx context.Context, command string, args []string, workDir string) (interface{}, error) {
	cmd := exec.CommandContext(ctx, command, args...)

	if workDir != "" {
		cmd.Dir = workDir
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("command failed: %w\nOutput: %s", err, string(output))
	}

	return map[string]interface{}{
		"stdout":    string(output),
		"exit_code": 0,
	}, nil
}
```

**YAML Example:**
```yaml
version: "1.0"
type: shell.command
name: Run Blender Script
parameters:
  command: blender
  args:
    - --background
    - --python
    - render_product.py
    - --
    - "--overlay={{overlay_path}}"
    - "--model={{stl_path}}"
    - "--output={{output_path}}"
```

### 4. WebP Conversion Node

**File**: `pkg/neta/image/webp.go` (NEW)
**Type**: `image.webp`

```go
// Package image provides image processing operations
package image

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"bento/pkg/neta"
)

// WebPConverter converts images to WebP
type WebPConverter struct{}

// Execute converts image to WebP
func (w WebPConverter) Execute(ctx context.Context, params neta.Params) (interface{}, error) {
	input, err := params.GetString("input")
	if err != nil {
		return nil, err
	}

	output := params.GetString("output", "")
	quality := params.GetInt("quality", 80)

	if output == "" {
		// Auto-generate output path
		ext := filepath.Ext(input)
		output = strings.TrimSuffix(input, ext) + ".webp"
	}

	return convertToWebP(ctx, input, output, quality)
}

// convertToWebP uses cwebp command
func convertToWebP(ctx context.Context, input, output string, quality int) (interface{}, error) {
	// Check if cwebp is available
	if _, err := exec.LookPath("cwebp"); err != nil {
		return nil, fmt.Errorf("cwebp not found: install webp tools")
	}

	cmd := exec.CommandContext(ctx, "cwebp",
		"-q", fmt.Sprintf("%d", quality),
		input,
		"-o", output,
	)

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("webp conversion failed: %w", err)
	}

	return map[string]interface{}{
		"input":  input,
		"output": output,
	}, nil
}
```

**YAML Example:**
```yaml
version: "1.0"
type: image.webp
name: Convert to WebP
parameters:
  input: "output/product/render.png"
  quality: 80
```

## Secure Configuration System

### Settings Storage

**File**: `pkg/omise/config/secrets.go` (NEW)

```go
// Package config provides configuration management
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Secrets manages sensitive configuration
type Secrets struct {
	configPath string
	data       map[string]string
}

// NewSecrets creates secrets manager
func NewSecrets() (*Secrets, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(home, ".bento", "secrets.json")

	s := &Secrets{
		configPath: configPath,
		data:       make(map[string]string),
	}

	// Load existing secrets
	if err := s.Load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return s, nil
}

// Get retrieves a secret
func (s *Secrets) Get(key string) (string, bool) {
	// Check environment variable first
	if val := os.Getenv(key); val != "" {
		return val, true
	}

	// Check stored secrets
	val, ok := s.data[key]
	return val, ok
}

// Set stores a secret
func (s *Secrets) Set(key, value string) error {
	s.data[key] = value
	return s.Save()
}

// Delete removes a secret
func (s *Secrets) Delete(key string) error {
	delete(s.data, key)
	return s.Save()
}

// List returns all secret keys (not values)
func (s *Secrets) List() []string {
	keys := make([]string, 0, len(s.data))
	for key := range s.data {
		keys = append(keys, key)
	}
	return keys
}

// Load reads secrets from disk
func (s *Secrets) Load() error {
	data, err := os.ReadFile(s.configPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &s.data)
}

// Save writes secrets to disk
func (s *Secrets) Save() error {
	// Ensure directory exists
	dir := filepath.Dir(s.configPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}

	// Write with restricted permissions
	return os.WriteFile(s.configPath, data, 0600)
}
```

### Settings Screen Enhancement

**File**: `pkg/omise/screens/settings_secrets.go` (NEW)

Add secrets management to settings screen:

```go
package screens

import (
	"github.com/charmbracelet/huh"
	"bento/pkg/omise/config"
)

// SecretsForm manages API keys and tokens
func SecretsForm(secrets *config.Secrets) *huh.Form {
	var (
		action     string
		key        string
		value      string
		deleteKey  string
	)

	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Action").
				Options(
					huh.NewOption("Add/Update Secret", "add"),
					huh.NewOption("Delete Secret", "delete"),
					huh.NewOption("List Secrets", "list"),
				).
				Value(&action),
		),

		// Add/Update
		huh.NewGroup(
			huh.NewInput().
				Title("Secret Name").
				Placeholder("FIGMA_API_TOKEN").
				Value(&key),

			huh.NewInput().
				Title("Secret Value").
				Placeholder("figd_...").
				EchoMode(huh.EchoModePassword).
				Value(&value),
		).WithHideFunc(func() bool {
			return action != "add"
		}),

		// Delete
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select Secret to Delete").
				OptionsFunc(func() []huh.Option[string] {
					keys := secrets.List()
					opts := make([]huh.Option[string], len(keys))
					for i, k := range keys {
						opts[i] = huh.NewOption(k, k)
					}
					return opts
				}, &deleteKey).
				Value(&deleteKey),
		).WithHideFunc(func() bool {
			return action != "delete"
		}),
	)
}
```

## The Complete Bentos

### Bento 1: generate-figma-image.bento.yaml (Reusable Component)

```yaml
version: "1.0"
type: group.sequence
name: Generate Figma Image

# This bento accepts parameters from the caller
parameters:
  figma_component_id: ""
  output_path: ""

nodes:
  # Step 1: Request Figma export
  - type: http
    name: Request Figma Export
    parameters:
      method: GET
      url: "https://api.figma.com/v1/images/{{config.figma_file_id}}"
      headers:
        X-Figma-Token: "{{config.figma_api_key}}"
      query:
        ids: "{{params.figma_component_id}}"
        format: png
        scale: 2

  # Step 2: Download the exported PNG
  - type: http
    name: Download Exported PNG
    parameters:
      method: GET
      url: "{{previous.images[params.figma_component_id]}}"
      output_file: "{{params.output_path}}"

  # Return the file path
  # Output: { png_path: "..." }
```

**This bento can be:**
- ✅ Tested independently
- ✅ Used in multiple workflows
- ✅ Shared across projects

---

### Bento 2: etsy-product-pipeline.bento.yaml (Main Workflow)

```yaml
version: "1.0"
type: group.sequence
name: Etsy Product Image Pipeline

nodes:
  # Step 1: Read CSV manifest
  - type: file.csv.read
    name: Read Product Manifest
    parameters:
      path: manifest.csv
      has_header: true

  # Step 2: Loop through each product
  - type: loop.for
    name: Process Each Product
    parameters:
      items: "{{previous}}"
      body:
        type: group.sequence
        name: Product Processing Steps
        nodes:
          # 2a: Create output folder
          - type: file.mkdir
            name: Create Product Folder
            parameters:
              path: "output/{{item.name}}"
              recursive: true

          # 2b: Execute the Figma bento (composition!)
          - type: bento.execute
            name: Generate Figma Overlay
            parameters:
              bento: generate-figma-image
              inputs:
                figma_component_id: "{{item.figma_component_id}}"
                output_path: "output/{{item.name}}/overlay.png"

          # 2c: Run Blender rendering
          - type: shell.exec
            name: Render with Blender
            parameters:
              command: blender
              args:
                - --background
                - --python
                - render_product.py
                - --
                - "--overlay=output/{{item.name}}/overlay.png"
                - "--model=models/{{item.model_file}}"
                - "--output=output/{{item.name}}/render.png"
                - "--sku={{item.sku}}"

          # 2d: Convert to WebP
          - type: image.convert
            name: Convert to WebP
            parameters:
              input: "output/{{item.name}}/render.png"
              output: "output/{{item.name}}/{{item.sku}}.webp"
              format: webp
              quality: 85
```

**This demonstrates:**
- ✅ Real-world complexity
- ✅ Bento composition (calling generate-figma-image)
- ✅ All new node types working together
- ✅ Production-ready workflow

## Implementation Plan

### Phase 10a: Bento Composition (1 hour)

**Priority 1: This enables everything else!**

1. **`bento.execute` Node Type**
   - Load bento from Jubako
   - Pass inputs to sub-bento
   - Execute with Itamae (recursive)
   - Return output
   - Tests

2. **Parameter Resolution**
   - Support `{{params.xxx}}` in sub-bentos
   - Support `{{config.xxx}}` for secrets
   - Variable interpolation

3. **Integration with Itamae**
   - Ensure Itamae can handle recursive execution
   - Context propagation
   - Error handling in nested bentos

**Deliverables:**
- ✅ `pkg/neta/bento/execute.go`
- ✅ Tests proving recursive execution works
- ✅ Example: bento-a calls bento-b

---

### Phase 10b: File & Shell Operations (1 hour)

1. **CSV Reader** (`file.csv.read`)
   - Parse CSV files
   - Handle headers
   - Return array of row objects
   - Tests

2. **Directory Operations** (`file.mkdir`)
   - Create directories
   - Recursive creation
   - Cross-platform paths
   - Tests

3. **Shell Execution** (`shell.exec`)
   - Run external commands
   - Argument handling
   - Working directory support
   - Output capture
   - Tests

**Deliverables:**
- ✅ `pkg/neta/file/csv.go`
- ✅ `pkg/neta/file/mkdir.go`
- ✅ `pkg/neta/shell/exec.go`
- ✅ All tests passing

---

### Phase 10c: Image Operations (30 min)

1. **Image Conversion** (`image.convert`)
   - WebP conversion (via cwebp or Go imaging)
   - Quality settings
   - Format detection
   - Tests

**Deliverables:**
- ✅ `pkg/neta/image/convert.go`
- ✅ Tests with sample images

### Phase 10d: Secure Configuration (30 min)

1. **Secrets Manager**
   - File-based storage (~/.bento/secrets.json)
   - Permissions: 0600 (owner read/write only)
   - Environment variable fallback

2. **Settings UI**
   - Add "API Keys & Secrets" section
   - Huh form for managing secrets
   - List/add/update/delete

3. **Template Variables**
   - Support `{{SECRET_NAME}}` in parameters
   - Resolve from secrets at runtime

### Phase 10e: Build & Test the Bentos (30-45 min)

1. **Create generate-figma-image.bento.yaml**
   - Use editor to build
   - Configure Figma API calls
   - Test independently with mock data
   - Validate output

2. **Create etsy-product-pipeline.bento.yaml**
   - Use editor to build main workflow
   - Add bento.execute node to call Figma bento
   - Configure all steps
   - Save

3. **End-to-End Test**
   - Prepare test CSV with 2-3 products
   - Configure Figma API key in Settings
   - Run: `bento pack etsy-product-pipeline.bento.yaml`
   - Verify:
     - ✅ Folders created
     - ✅ Figma overlays downloaded
     - ✅ Blender renders generated (if Blender available)
     - ✅ WebP files created
     - ✅ Sub-bento executed correctly

4. **Success Criteria**
   - You can run your actual workflow
   - All outputs generated correctly
   - Bento composition works
   - System is production-ready!

## Testing Strategy

### Unit Tests

```bash
# Test new node types
go test -v ./pkg/neta/file/
go test -v ./pkg/neta/shell/
go test -v ./pkg/neta/image/

# Test secrets manager
go test -v ./pkg/omise/config/
```

### Integration Test

```bash
# Create test CSV
cat > test-manifest.csv <<EOF
sku,name,model_file,figma_component_id,color
TEST001,Test Product,test.stl,123:456,red
EOF

# Create test bento
./bento  # Use editor to build workflow

# Run bento
./bento pack test-product-pipeline.bento.yaml

# Verify outputs
ls -la output/Test\ Product/
# Should see: overlay.png, render.png, TEST001.webp
```

## Success Criteria

Phase 10 is complete when:

1. ✅ **`bento.execute` node implemented** - Bentos can call other bentos
2. ✅ **Composable architecture proven** - "A node is a node is a node"
3. ✅ CSV reader node implemented and tested
4. ✅ File operation nodes working (mkdir, etc.)
5. ✅ Shell command node executing properly
6. ✅ Image conversion working
7. ✅ Secrets management implemented
8. ✅ Settings UI for API keys working
9. ✅ **generate-figma-image.bento.yaml** created and tested
10. ✅ **etsy-product-pipeline.bento.yaml** created using editor
11. ✅ **You can run your actual workflow** with real CSV data
12. ✅ All outputs generated correctly
13. ✅ Sub-bento execution validated
14. ✅ System is production-ready
15. ✅ **Karen's approval granted**

**The Ultimate Test:**
```bash
# You run this command with your actual manifest.csv:
bento pack etsy-product-pipeline.bento.yaml

# And it produces all your product images!
```

## Deliverables

### Code
- ✅ **`bento.execute` node** - Composable bento architecture
- ✅ 4 new node type packages (file, shell, image)
- ✅ Secrets manager
- ✅ Enhanced settings screen
- ✅ All tests passing

### Documentation
- ✅ Composable bento guide
- ✅ Node type documentation
- ✅ Secrets setup guide
- ✅ Example bentos for new nodes
- ✅ Etsy pipeline tutorial

### Working Bentos
- ✅ **`generate-figma-image.bento.yaml`** - Reusable component
- ✅ **`etsy-product-pipeline.bento.yaml`** - Main workflow
- ✅ Tested with real CSV data
- ✅ All steps executing correctly
- ✅ **You can run your production workflow!**

## Future Enhancements (Post-Phase 10)

These are identified needs but not required for Phase 10:

1. **Etsy API Integration**
   - Upload product images
   - Update product listings
   - Manage inventory

2. **Advanced Image Processing**
   - Batch operations
   - Filters and effects
   - Format conversions

3. **Error Recovery**
   - Retry failed steps
   - Partial completion handling
   - Resume from checkpoint

4. **Parallel Processing**
   - Process multiple products simultaneously
   - Thread pool management

## Documentation Example

**README Section:**

```markdown
## 🎨 Real-World Example: Etsy Product Pipeline

This example demonstrates using Bento to automate Etsy product image generation.

### What it does
1. Reads product data from CSV
2. Generates Figma overlays for each product
3. Renders images with Blender
4. Converts to WebP for web optimization

### Prerequisites
- Figma API token
- Blender installed
- `cwebp` installed (WebP conversion)

### Setup
1. Store Figma token:
   ```bash
   bento  # Launch TUI
   # Navigate to Settings → API Keys
   # Add: FIGMA_API_TOKEN=your_token_here
   ```

2. Prepare manifest.csv:
   ```csv
   sku,name,model_file,figma_component_id
   SKU001,Product 1,model1.stl,123:456
   ```

3. Create bento:
   ```bash
   bento  # Use editor to build pipeline
   # Or copy example:
   cp examples/etsy-pipeline.bento.yaml my-pipeline.bento.yaml
   ```

4. Run:
   ```bash
   bento pack my-pipeline.bento.yaml
   ```

### Output
```
output/
├── Product 1/
│   ├── overlay.png
│   ├── render.png
│   └── SKU001.webp
```
```

## Execution Prompt

```
I'm ready to begin Phase 10: Real-World Proof-of-Concept.

I have read the Bento Box Principle and will follow it.

This phase validates the core principle: "A node is a node is a node"

Please implement:
1. bento.execute node (PRIORITY - enables composition)
2. New node types (CSV, file ops, shell, image)
3. Secrets management system
4. Build TWO bentos:
   - generate-figma-image.bento.yaml (reusable)
   - etsy-product-pipeline.bento.yaml (main)
5. Test with real workflow

Success = I can run: bento pack etsy-product-pipeline.bento.yaml
And it generates my product images!

Each file < 250 lines, functions < 20 lines. I will use TodoWrite to track progress and get Karen's approval before completing.
```

---

**Phase 10 Proof-of-Concept**: Composable bentos + Real-world validation 🎨🍱

**After this phase**:
- ✅ Bento composition proven
- ✅ Production-ready system
- ✅ Your Etsy pipeline running!
- ✅ Battle-tested architecture 🚀
