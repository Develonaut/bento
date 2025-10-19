# Phase 3: Omakase Package (おまかせ - "Chef's Choice")

**Duration:** 2-3 days
**Package:** `pkg/omakase/`
**Dependencies:** `pkg/neta`

---

## TDD Philosophy

> **Write tests FIRST to define contracts**

Validation tests should verify:
1. Struct tag validation works (required fields, format constraints)
2. Custom validators for bento-specific rules
3. Clear, actionable error messages
4. Pre-flight checks (Blender installed, Figma API credentials exist)
5. Nested structure validation (group neta with child nodes)

---

## Phase Overview

The omakase ("chef's choice") package provides validation for bento files and neta definitions. Like a sushi chef who validates ingredients before serving, omakase ensures all neta are properly configured before execution.

It wraps `go-playground/validator` with bento-specific validation rules:

- **Struct tag validation:** Required fields, format constraints
- **Custom validators:** Bento-specific business rules
- **Clear error messages:** "Missing required field 'url' in http-request neta node-1"
- **Pre-flight checks:** Verify environment (Blender, Figma API) before execution
- **Nested validation:** Validate group neta with child nodes and edges

### Why "Omakase"?

"Omakase" (おまかせ) means "I'll leave it up to you" or "chef's choice." Just as you trust the sushi chef to select the best ingredients, you trust omakase to validate your bento has all the right ingredients configured correctly.

---

## Success Criteria

**Phase 3 Complete When:**
- [ ] Struct tag validation for all neta types
- [ ] Custom validators for bento-specific rules
- [ ] Pre-flight environment checks
- [ ] Clear, actionable error messages
- [ ] Integration tests for each neta type's validation
- [ ] Files < 250 lines each
- [ ] File-level documentation complete
- [ ] `/code-review` run with Karen + Colossus approval

---

## Test-First Approach

### Step 1: Define validation interface via tests

Create `pkg/omakase/omakase_test.go`:

```go
package omakase_test

import (
	"testing"

	"github.com/yourusername/bento/pkg/neta"
	"github.com/yourusername/bento/pkg/omakase"
)

// Test: Valid definition should pass validation
func TestValidator_ValidDefinition(t *testing.T) {
	validator := omakase.New()

	def := &neta.Definition{
		ID:      "node-1",
		Type:    "http-request",
		Version: "1.0.0",
		Name:    "Fetch User Data",
		Parameters: map[string]interface{}{
			"url":    "https://api.example.com/users",
			"method": "GET",
		},
		InputPorts: []neta.Port{
			{ID: "in-1", Name: "input"},
		},
		OutputPorts: []neta.Port{
			{ID: "out-1", Name: "output"},
		},
	}

	err := validator.Validate(def)
	if err != nil {
		t.Fatalf("Valid definition should pass validation: %v", err)
	}
}

// Test: Missing required field should fail with clear error
func TestValidator_MissingRequiredField(t *testing.T) {
	validator := omakase.New()

	def := &neta.Definition{
		ID:      "node-1",
		Type:    "http-request",
		Version: "1.0.0",
		Name:    "Fetch Data",
		Parameters: map[string]interface{}{
			// Missing required "url" field
			"method": "GET",
		},
	}

	err := validator.Validate(def)
	if err == nil {
		t.Fatal("Expected validation error for missing URL")
	}

	// Error message should be clear and actionable
	errMsg := err.Error()
	if !contains(errMsg, "url") {
		t.Errorf("Error should mention missing 'url': %s", errMsg)
	}

	if !contains(errMsg, "node-1") {
		t.Errorf("Error should mention node ID 'node-1': %s", errMsg)
	}
}

// Test: Invalid URL format should fail
func TestValidator_InvalidURLFormat(t *testing.T) {
	validator := omakase.New()

	def := &neta.Definition{
		ID:      "node-1",
		Type:    "http-request",
		Version: "1.0.0",
		Name:    "Fetch Data",
		Parameters: map[string]interface{}{
			"url":    "not-a-valid-url",  // Invalid URL
			"method": "GET",
		},
	}

	err := validator.Validate(def)
	if err == nil {
		t.Fatal("Expected validation error for invalid URL")
	}

	errMsg := err.Error()
	if !contains(errMsg, "invalid URL") {
		t.Errorf("Error should mention invalid URL: %s", errMsg)
	}
}

// Test: Validate group neta with nested nodes
func TestValidator_GroupWithNestedNodes(t *testing.T) {
	validator := omakase.New()

	def := &neta.Definition{
		ID:      "group-1",
		Type:    "group",
		Version: "1.0.0",
		Name:    "Main Group",
		Nodes: []neta.Definition{
			{
				ID:      "node-1",
				Type:    "edit-fields",
				Version: "1.0.0",
				Name:    "Set Fields",
				Parameters: map[string]interface{}{
					"values": map[string]interface{}{
						"foo": "bar",
					},
				},
			},
			{
				ID:      "node-2",
				Type:    "http-request",
				Version: "1.0.0",
				Name:    "Fetch Data",
				Parameters: map[string]interface{}{
					"url":    "https://api.example.com",
					"method": "GET",
				},
			},
		},
		Edges: []neta.Edge{
			{
				ID:     "edge-1",
				Source: "node-1",
				Target: "node-2",
			},
		},
	}

	err := validator.Validate(def)
	if err != nil {
		t.Fatalf("Valid group should pass validation: %v", err)
	}
}

// Test: Invalid edge (source doesn't exist) should fail
func TestValidator_InvalidEdge(t *testing.T) {
	validator := omakase.New()

	def := &neta.Definition{
		ID:      "group-1",
		Type:    "group",
		Version: "1.0.0",
		Name:    "Main Group",
		Nodes: []neta.Definition{
			{
				ID:      "node-1",
				Type:    "edit-fields",
				Version: "1.0.0",
				Name:    "Set Fields",
			},
		},
		Edges: []neta.Edge{
			{
				ID:     "edge-1",
				Source: "nonexistent-node",  // Doesn't exist!
				Target: "node-1",
			},
		},
	}

	err := validator.Validate(def)
	if err == nil {
		t.Fatal("Expected validation error for invalid edge")
	}

	errMsg := err.Error()
	if !contains(errMsg, "nonexistent-node") {
		t.Errorf("Error should mention missing source node: %s", errMsg)
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
```

### Step 2: Test neta-specific validators

```go
// Test: HTTP request should validate method
func TestValidator_HTTPMethod(t *testing.T) {
	validator := omakase.New()

	def := &neta.Definition{
		ID:      "node-1",
		Type:    "http-request",
		Version: "1.0.0",
		Name:    "Fetch Data",
		Parameters: map[string]interface{}{
			"url":    "https://api.example.com",
			"method": "INVALID_METHOD",  // Should be GET, POST, etc.
		},
	}

	err := validator.Validate(def)
	if err == nil {
		t.Fatal("Expected validation error for invalid HTTP method")
	}
}

// Test: Loop neta should validate mode
func TestValidator_LoopMode(t *testing.T) {
	validator := omakase.New()

	def := &neta.Definition{
		ID:      "node-1",
		Type:    "loop",
		Version: "1.0.0",
		Name:    "Loop Through Items",
		Parameters: map[string]interface{}{
			"mode": "invalidMode",  // Should be forEach, times, or while
		},
	}

	err := validator.Validate(def)
	if err == nil {
		t.Fatal("Expected validation error for invalid loop mode")
	}
}

// Test: File-system neta should validate operation
func TestValidator_FileSystemOperation(t *testing.T) {
	validator := omakase.New()

	def := &neta.Definition{
		ID:      "node-1",
		Type:    "file-system",
		Version: "1.0.0",
		Name:    "File Operation",
		Parameters: map[string]interface{}{
			"operation": "invalidOp",  // Should be read, write, copy, move, delete
		},
	}

	err := validator.Validate(def)
	if err == nil {
		t.Fatal("Expected validation error for invalid file operation")
	}
}
```

### Step 3: Test pre-flight checks

```go
// Test: Pre-flight should check for required commands
func TestValidator_PreflightCheck(t *testing.T) {
	validator := omakase.New()

	def := &neta.Definition{
		ID:      "node-1",
		Type:    "shell-command",
		Version: "1.0.0",
		Name:    "Run Blender",
		Parameters: map[string]interface{}{
			"command": "blender",
			"args":    []string{"--version"},
		},
	}

	// Preflight check should verify Blender is installed
	err := validator.PreflightCheck(def)

	// This might fail if Blender isn't installed - that's OK for the test
	// We're testing that the check RUNS, not that it passes
	if err != nil {
		t.Logf("Preflight check failed (expected if Blender not installed): %v", err)
	}
}

// Test: Pre-flight should check for environment variables
func TestValidator_PreflightEnvVars(t *testing.T) {
	validator := omakase.New()

	def := &neta.Definition{
		ID:      "node-1",
		Type:    "http-request",
		Version: "1.0.0",
		Name:    "Call Figma API",
		Parameters: map[string]interface{}{
			"url": "https://api.figma.com/v1/files/{{.FIGMA_FILE_ID}}",
			"headers": map[string]string{
				"X-Figma-Token": "{{.FIGMA_API_TOKEN}}",
			},
		},
	}

	// Should check that FIGMA_API_TOKEN env var exists
	err := validator.PreflightCheck(def)

	if err != nil {
		// Expected if env var not set
		t.Logf("Preflight check failed (expected if FIGMA_API_TOKEN not set): %v", err)
	}
}
```

---

## File Structure

```
pkg/omakase/
├── omakase.go           # Main validator implementation (~200 lines)
├── validators.go        # Neta-specific validators (~200 lines)
├── preflight.go         # Pre-flight environment checks (~150 lines)
├── errors.go            # Clear error message formatting (~100 lines)
└── omakase_test.go      # Integration tests (~350 lines)
```

---

## Implementation Guidance

**File: `pkg/omakase/omakase.go`**

```go
// Package omakase provides validation for bento files and neta definitions.
//
// "Omakase" (おまかせ - "chef's choice") ensures all neta are properly configured
// before execution, just as a sushi chef validates ingredients before serving.
//
// It uses go-playground/validator for struct tag validation with custom
// validators for bento-specific business rules.
//
// Usage:
//
//	validator := omakase.New()
//
//	// Validate a neta definition
//	if err := validator.Validate(netaDef); err != nil {
//	    log.Fatalf("Invalid neta: %v", err)
//	}
//
//	// Pre-flight checks (environment, commands, API keys)
//	if err := validator.PreflightCheck(netaDef); err != nil {
//	    log.Fatalf("Pre-flight check failed: %v", err)
//	}
//
// Learn more about go-playground/validator:
// https://github.com/go-playground/validator
package omakase

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/yourusername/bento/pkg/neta"
)

// Validator validates neta definitions and performs pre-flight checks.
type Validator struct {
	v *validator.Validate
}

// New creates a new Validator with all custom validators registered.
func New() *Validator {
	v := validator.New()

	// Register custom validators
	v.RegisterValidation("netaType", validateNetaType)
	v.RegisterValidation("httpMethod", validateHTTPMethod)
	v.RegisterValidation("loopMode", validateLoopMode)
	v.RegisterValidation("fileOperation", validateFileOperation)

	return &Validator{v: v}
}

// Validate validates a neta definition.
//
// Returns a ValidationError with clear, actionable error messages if validation fails.
func (val *Validator) Validate(def *neta.Definition) error {
	// Validate core structure
	if err := val.validateCore(def); err != nil {
		return err
	}

	// Validate type-specific parameters
	if err := val.validateTypeSpecific(def); err != nil {
		return err
	}

	// Validate nested nodes (for group neta)
	if def.Type == "group" {
		if err := val.validateGroup(def); err != nil {
			return err
		}
	}

	return nil
}

// validateCore validates core fields (ID, Type, Version, etc.)
func (val *Validator) validateCore(def *neta.Definition) error {
	if def.ID == "" {
		return fmt.Errorf("neta is missing required field 'id'")
	}

	if def.Type == "" {
		return fmt.Errorf("neta '%s' is missing required field 'type'", def.ID)
	}

	if def.Version == "" {
		return fmt.Errorf("neta '%s' is missing required field 'version'", def.ID)
	}

	return nil
}

// validateTypeSpecific validates neta-specific parameters
func (val *Validator) validateTypeSpecific(def *neta.Definition) error {
	switch def.Type {
	case "http-request":
		return val.validateHTTPRequest(def)
	case "file-system":
		return val.validateFileSystem(def)
	case "shell-command":
		return val.validateShellCommand(def)
	case "loop":
		return val.validateLoop(def)
	case "spreadsheet":
		return val.validateSpreadsheet(def)
	case "image":
		return val.validateImage(def)
	case "transform":
		return val.validateTransform(def)
	case "edit-fields":
		return val.validateEditFields(def)
	case "group":
		return nil // Validated in validateGroup
	case "parallel":
		return val.validateParallel(def)
	default:
		return fmt.Errorf("neta '%s' has unknown type '%s'", def.ID, def.Type)
	}
}

// validateGroup validates group neta (nested nodes and edges)
func (val *Validator) validateGroup(def *neta.Definition) error {
	// Validate all child nodes
	for _, child := range def.Nodes {
		if err := val.Validate(&child); err != nil {
			return fmt.Errorf("invalid child node in group '%s': %w", def.ID, err)
		}
	}

	// Validate edges (source and target must exist)
	nodeIDs := make(map[string]bool)
	for _, node := range def.Nodes {
		nodeIDs[node.ID] = true
	}

	for _, edge := range def.Edges {
		if !nodeIDs[edge.Source] {
			return fmt.Errorf("edge '%s' in group '%s' has invalid source '%s' (node doesn't exist)",
				edge.ID, def.ID, edge.Source)
		}

		if !nodeIDs[edge.Target] {
			return fmt.Errorf("edge '%s' in group '%s' has invalid target '%s' (node doesn't exist)",
				edge.ID, def.ID, edge.Target)
		}
	}

	return nil
}

// PreflightCheck performs environment checks before execution.
//
// This includes:
//   - Verifying required commands are installed (blender, etc.)
//   - Checking environment variables exist (FIGMA_API_TOKEN, etc.)
//   - Validating file paths exist
func (val *Validator) PreflightCheck(def *neta.Definition) error {
	switch def.Type {
	case "shell-command":
		return val.preflightShellCommand(def)
	case "http-request":
		return val.preflightHTTPRequest(def)
	case "file-system":
		return val.preflightFileSystem(def)
	}

	return nil
}
```

**File: `pkg/omakase/validators.go`**

```go
package omakase

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/yourusername/bento/pkg/neta"
)

// validateHTTPRequest validates http-request neta parameters
func (val *Validator) validateHTTPRequest(def *neta.Definition) error {
	url, ok := def.Parameters["url"].(string)
	if !ok || url == "" {
		return fmt.Errorf("http-request neta '%s' missing required parameter 'url'", def.ID)
	}

	method, ok := def.Parameters["method"].(string)
	if !ok || method == "" {
		return fmt.Errorf("http-request neta '%s' missing required parameter 'method'", def.ID)
	}

	// Validate HTTP method
	validMethods := map[string]bool{
		"GET": true, "POST": true, "PUT": true, "PATCH": true, "DELETE": true,
		"HEAD": true, "OPTIONS": true,
	}

	if !validMethods[method] {
		return fmt.Errorf("http-request neta '%s' has invalid method '%s' (must be GET, POST, PUT, PATCH, DELETE, HEAD, or OPTIONS)",
			def.ID, method)
	}

	return nil
}

// validateLoop validates loop neta parameters
func (val *Validator) validateLoop(def *neta.Definition) error {
	mode, ok := def.Parameters["mode"].(string)
	if !ok || mode == "" {
		return fmt.Errorf("loop neta '%s' missing required parameter 'mode'", def.ID)
	}

	validModes := map[string]bool{
		"forEach": true, "times": true, "while": true,
	}

	if !validModes[mode] {
		return fmt.Errorf("loop neta '%s' has invalid mode '%s' (must be forEach, times, or while)",
			def.ID, mode)
	}

	// Mode-specific validation
	switch mode {
	case "forEach":
		if def.Parameters["items"] == nil {
			return fmt.Errorf("loop neta '%s' with mode 'forEach' missing required parameter 'items'", def.ID)
		}
	case "times":
		if def.Parameters["count"] == nil {
			return fmt.Errorf("loop neta '%s' with mode 'times' missing required parameter 'count'", def.ID)
		}
	case "while":
		if def.Parameters["condition"] == nil {
			return fmt.Errorf("loop neta '%s' with mode 'while' missing required parameter 'condition'", def.ID)
		}
	}

	return nil
}

// validateFileSystem validates file-system neta parameters
func (val *Validator) validateFileSystem(def *neta.Definition) error {
	operation, ok := def.Parameters["operation"].(string)
	if !ok || operation == "" {
		return fmt.Errorf("file-system neta '%s' missing required parameter 'operation'", def.ID)
	}

	validOps := map[string]bool{
		"read": true, "write": true, "copy": true, "move": true, "delete": true,
		"mkdir": true, "exists": true,
	}

	if !validOps[operation] {
		return fmt.Errorf("file-system neta '%s' has invalid operation '%s'", def.ID, operation)
	}

	return nil
}

// validateShellCommand validates shell-command neta parameters
func (val *Validator) validateShellCommand(def *neta.Definition) error {
	command, ok := def.Parameters["command"].(string)
	if !ok || command == "" {
		return fmt.Errorf("shell-command neta '%s' missing required parameter 'command'", def.ID)
	}

	return nil
}

// Additional validators for other neta types...
func (val *Validator) validateEditFields(def *neta.Definition) error { return nil }
func (val *Validator) validateSpreadsheet(def *neta.Definition) error { return nil }
func (val *Validator) validateImage(def *neta.Definition) error { return nil }
func (val *Validator) validateTransform(def *neta.Definition) error { return nil }
func (val *Validator) validateParallel(def *neta.Definition) error { return nil }

// Custom validator functions for go-playground/validator
func validateNetaType(fl validator.FieldLevel) bool {
	// Validate neta type is one of the 10 valid types
	validTypes := map[string]bool{
		"edit-fields": true, "http-request": true, "file-system": true,
		"shell-command": true, "group": true, "loop": true, "parallel": true,
		"spreadsheet": true, "image": true, "transform": true,
	}

	return validTypes[fl.Field().String()]
}

func validateHTTPMethod(fl validator.FieldLevel) bool {
	validMethods := map[string]bool{
		"GET": true, "POST": true, "PUT": true, "PATCH": true, "DELETE": true,
	}

	return validMethods[fl.Field().String()]
}

func validateLoopMode(fl validator.FieldLevel) bool {
	validModes := map[string]bool{"forEach": true, "times": true, "while": true}
	return validModes[fl.Field().String()]
}

func validateFileOperation(fl validator.FieldLevel) bool {
	validOps := map[string]bool{
		"read": true, "write": true, "copy": true, "move": true, "delete": true,
	}

	return validOps[fl.Field().String()]
}
```

**File: `pkg/omakase/preflight.go`**

```go
package omakase

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/yourusername/bento/pkg/neta"
)

// preflightShellCommand checks if the command exists
func (val *Validator) preflightShellCommand(def *neta.Definition) error {
	command, ok := def.Parameters["command"].(string)
	if !ok {
		return nil // Already validated by validateShellCommand
	}

	// Check if command exists in PATH
	_, err := exec.LookPath(command)
	if err != nil {
		return fmt.Errorf("shell-command neta '%s': command '%s' not found in PATH. Please install it first.",
			def.ID, command)
	}

	return nil
}

// preflightHTTPRequest checks for required environment variables in URL/headers
func (val *Validator) preflightHTTPRequest(def *neta.Definition) error {
	// Check URL for template variables
	url, _ := def.Parameters["url"].(string)
	if envVars := extractEnvVars(url); len(envVars) > 0 {
		for _, envVar := range envVars {
			if os.Getenv(envVar) == "" {
				return fmt.Errorf("http-request neta '%s': environment variable '%s' not set (required in URL)",
					def.ID, envVar)
			}
		}
	}

	// Check headers for template variables
	if headers, ok := def.Parameters["headers"].(map[string]string); ok {
		for key, value := range headers {
			if envVars := extractEnvVars(value); len(envVars) > 0 {
				for _, envVar := range envVars {
					if os.Getenv(envVar) == "" {
						return fmt.Errorf("http-request neta '%s': environment variable '%s' not set (required in header '%s')",
							def.ID, envVar, key)
					}
				}
			}
		}
	}

	return nil
}

// preflightFileSystem checks if file paths exist
func (val *Validator) preflightFileSystem(def *neta.Definition) error {
	operation, _ := def.Parameters["operation"].(string)

	// For read operations, check file exists
	if operation == "read" {
		if path, ok := def.Parameters["path"].(string); ok {
			if _, err := os.Stat(path); os.IsNotExist(err) {
				return fmt.Errorf("file-system neta '%s': file not found: %s", def.ID, path)
			}
		}
	}

	return nil
}

// extractEnvVars finds {{.VAR_NAME}} patterns in a string
func extractEnvVars(s string) []string {
	var vars []string

	// Simple regex-free approach: look for {{.WORD}}
	for {
		start := strings.Index(s, "{{.")
		if start == -1 {
			break
		}

		end := strings.Index(s[start:], "}}")
		if end == -1 {
			break
		}

		varName := s[start+3 : start+end]
		vars = append(vars, varName)

		s = s[start+end+2:]
	}

	return vars
}
```

---

## Critical for Phase 8

**Pre-flight Checks:**
- Must check Blender is installed before trying to render
- Must check Figma API token exists before making requests
- Must check CSV file exists before trying to read it
- Fail fast with clear error messages (don't wait until row 23 to discover Blender isn't installed)

**Clear Error Messages:**
- "shell-command neta 'render-product': command 'blender' not found in PATH. Please install Blender first."
- "http-request neta 'figma-overlay': environment variable 'FIGMA_API_TOKEN' not set (required in header 'X-Figma-Token')"
- "spreadsheet neta 'read-products': file not found: products.csv"

---

## Bento Box Principle Checklist

- [ ] Files < 250 lines (omakase.go ~200, validators.go ~200, preflight.go ~150)
- [ ] Functions < 20 lines
- [ ] Single responsibility (validation only)
- [ ] Clear error messages (mention neta ID, field name, expected format)
- [ ] File-level documentation

---

## Phase Completion

**Phase 3 MUST end with:**

1. All tests passing (`go test ./pkg/omakase/...`)
2. Run `/code-review` slash command
3. Address feedback from Karen and Colossus
4. Get explicit approval from both agents
5. Document any decisions in `.claude/strategy/`

**Do not proceed to Phase 4 until code review is approved.**

---

## Claude Prompt Template

```
I need to implement Phase 3: omakase (validation package) following TDD principles.

Please read:
- .claude/strategy/phase-3-omakase.md (this file)
- .claude/BENTO_BOX_PRINCIPLE.md
- .claude/COMPLETE_NODE_INVENTORY.md (for neta parameter requirements)

Then:

1. Create `pkg/omakase/omakase_test.go` with integration tests for:
   - Core field validation (ID, Type, Version)
   - Neta-specific parameter validation (each of the 10 types)
   - Group neta with nested nodes and edges
   - Invalid edges (source/target doesn't exist)
   - Pre-flight checks (commands, env vars, file existence)
   - Clear error messages (mention neta ID)

2. Watch the tests fail

3. Implement to make tests pass:
   - pkg/omakase/omakase.go (~200 lines)
   - pkg/omakase/validators.go (~200 lines)
   - pkg/omakase/preflight.go (~150 lines)
   - pkg/omakase/errors.go (~100 lines) - optional

4. Add file-level documentation

Remember:
- Write tests FIRST
- Clear, actionable error messages
- Pre-flight checks are CRITICAL for Phase 8 (fail fast!)
- Files < 250 lines
- Functions < 20 lines

When complete, run `/code-review` and get Karen + Colossus approval.
```

---

## Dependencies to Add

```bash
go get github.com/go-playground/validator/v10
```

---

## Notes

- Validation prevents cryptic runtime errors
- Pre-flight checks save time (don't start rendering 50 products only to fail on #23)
- Clear error messages are critical for debugging
- Group validation must be recursive (nested groups)
- Template variable extraction ({{.VAR}}) is simple string parsing, no regex needed

---

**Status:** Ready for implementation
**Next Phase:** Phase 4 (hangiri storage) - depends on completion of Phases 1-3
