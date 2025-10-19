# Phase 7: cmd/bento - CLI Commands

**Duration:** 1 week
**Package:** `cmd/bento/`
**Dependencies:** ALL packages (neta, itamae, pantry, hangiri, shoyu, omakase)

---

## TDD Philosophy

> **Write tests FIRST to define contracts**

CLI tests should verify:
1. `bento serve` executes a bento file
2. `bento inspect` validates without executing
3. `bento menu` lists available bentos
4. `bento new` creates a template bento
5. Flags are parsed correctly (--verbose, --timeout, etc.)
6. Exit codes are correct (0 for success, 1 for errors)
7. Output is user-friendly with sushi emojis

---

## Phase Overview

Phase 7 brings everything together into a beautiful CLI tool. The cmd/bento package implements all user-facing commands using Cobra, with playful sushi-themed naming.

Commands:
- **`bento serve`** - Execute a bento (like serving the finished dish)
- **`bento inspect`** - Validate a bento without executing
- **`bento menu`** - List available bentos (like a restaurant menu)
- **`bento new`** - Create a new bento template

---

## Commands

### `bento serve [file].bento.json`

Execute a bento workflow.

**Flags:**
- `--verbose`, `-v` - Verbose output
- `--timeout` - Execution timeout (default: 10m)
- `--var key=value` - Set template variables

**Examples:**
```bash
bento serve workflow.bento.json
bento serve workflow.bento.json --verbose
bento serve workflow.bento.json --timeout 30m
bento serve workflow.bento.json --var PRODUCT_ID=123
```

**Output:**
```
🍱 Serving bento: Product Automation
🍙 Executing neta 'read-csv' (spreadsheet)
✓ Completed in 124ms
🍙 Executing neta 'process-products' (loop)
  ⟳ Processing item 1/50: PROD-001
  ⟳ Processing item 2/50: PROD-002
  ...
✓ Completed in 3m 45s

🍱 Bento served successfully in 3m 46s
```

### `bento inspect [file].bento.json`

Validate a bento without executing it.

**Flags:**
- `--verbose`, `-v` - Show detailed validation results

**Examples:**
```bash
bento inspect workflow.bento.json
bento inspect workflow.bento.json --verbose
```

**Output:**
```
🍱 Inspecting bento: Product Automation
✓ Valid JSON structure
✓ All neta types registered
✓ All edges valid (sources and targets exist)
✓ Required parameters present
✓ Pre-flight checks passed

🍱 Bento is ready to serve!
```

### `bento menu [directory]`

List all available bentos in a directory.

**Flags:**
- `--recursive`, `-r` - Search subdirectories
- `--json` - Output as JSON

**Examples:**
```bash
bento menu
bento menu ~/workflows
bento menu ~/workflows --recursive
bento menu --json
```

**Output:**
```
🍱 Available Bentos:

  product-automation.bento.json
    Product Photo Automation
    10 neta, last modified 2 hours ago

  data-pipeline.bento.json
    Daily Data Pipeline
    5 neta, last modified 1 day ago

  image-batch.bento.json
    Batch Image Processing
    8 neta, last modified 3 days ago

3 bentos found
```

### `bento new [name]`

Create a new bento template.

**Flags:**
- `--type` - Template type (simple, loop, parallel)

**Examples:**
```bash
bento new my-workflow
bento new my-workflow --type loop
```

**Output:**
```
🍱 Creating new bento: my-workflow

Created: my-workflow.bento.json

Next steps:
  1. Edit my-workflow.bento.json
  2. Run: bento inspect my-workflow.bento.json
  3. Run: bento serve my-workflow.bento.json
```

---

## Success Criteria

**Phase 7 Complete When:**
- [ ] `bento serve` executes bentos with progress output
- [ ] `bento inspect` validates without executing
- [ ] `bento menu` lists available bentos
- [ ] `bento new` creates template bentos
- [ ] All flags work correctly
- [ ] Exit codes correct (0 success, 1 error)
- [ ] User-friendly output with emojis 🍱🍙
- [ ] Integration tests for each command
- [ ] Files < 250 lines each
- [ ] File-level documentation complete
- [ ] `/code-review` run with Karen + Colossus approval

---

## Test-First Approach

### Step 1: Test bento serve command

Create `cmd/bento/serve_test.go`:

```go
package main_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// Test: bento serve should execute a valid bento
func TestServeCommand_ValidBento(t *testing.T) {
	// Create a simple test bento
	bentoFile := createTestBento(t, "test.bento.json", `{
		"id": "test-bento",
		"type": "group",
		"version": "1.0.0",
		"name": "Test Bento",
		"nodes": [
			{
				"id": "node-1",
				"type": "edit-fields",
				"version": "1.0.0",
				"name": "Set Field",
				"parameters": {"values": {"foo": "bar"}}
			}
		],
		"edges": []
	}`)
	defer os.Remove(bentoFile)

	// Run: bento serve test.bento.json
	cmd := exec.Command("bento", "serve", bentoFile)
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("serve command failed: %v\nOutput: %s", err, string(output))
	}

	// Verify output contains success message
	outputStr := string(output)
	if !strings.Contains(outputStr, "served successfully") {
		t.Errorf("Output should contain 'served successfully': %s", outputStr)
	}

	// Verify emoji usage
	if !strings.Contains(outputStr, "🍱") {
		t.Error("Output should contain bento emoji 🍱")
	}
}

// Test: bento serve should fail on invalid bento
func TestServeCommand_InvalidBento(t *testing.T) {
	bentoFile := createTestBento(t, "invalid.bento.json", `{
		"id": "invalid",
		"type": "group",
		"nodes": []
	}`)
	defer os.Remove(bentoFile)

	cmd := exec.Command("bento", "serve", bentoFile)
	output, err := cmd.CombinedOutput()

	// Should exit with error code
	if err == nil {
		t.Fatal("serve command should fail on invalid bento")
	}

	// Check exit code
	if exitError, ok := err.(*exec.ExitError); ok {
		if exitError.ExitCode() != 1 {
			t.Errorf("Exit code = %d, want 1", exitError.ExitCode())
		}
	}

	// Output should mention the error
	outputStr := string(output)
	if !strings.Contains(outputStr, "error") || !strings.Contains(outputStr, "failed") {
		t.Errorf("Output should mention error: %s", outputStr)
	}
}

// Test: bento serve --verbose should show detailed output
func TestServeCommand_VerboseFlag(t *testing.T) {
	bentoFile := createTestBento(t, "test.bento.json", simpleValidBento())
	defer os.Remove(bentoFile)

	cmd := exec.Command("bento", "serve", bentoFile, "--verbose")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("serve command failed: %v", err)
	}

	outputStr := string(output)

	// Verbose mode should show node-level details
	if !strings.Contains(outputStr, "node-1") {
		t.Error("Verbose output should mention node IDs")
	}
}

func createTestBento(t *testing.T, name, content string) string {
	tmpfile, err := os.CreateTemp("", name)
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

### Step 2: Test bento inspect command

```go
// Test: bento inspect should validate without executing
func TestInspectCommand_ValidBento(t *testing.T) {
	bentoFile := createTestBento(t, "test.bento.json", simpleValidBento())
	defer os.Remove(bentoFile)

	cmd := exec.Command("bento", "inspect", bentoFile)
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("inspect command failed: %v\nOutput: %s", err, string(output))
	}

	outputStr := string(output)

	// Should show validation passed
	if !strings.Contains(outputStr, "ready to serve") {
		t.Errorf("Output should say 'ready to serve': %s", outputStr)
	}

	// Should NOT execute (check it doesn't say "served successfully")
	if strings.Contains(outputStr, "served successfully") {
		t.Error("inspect should NOT execute the bento")
	}
}

// Test: bento inspect should report validation errors clearly
func TestInspectCommand_InvalidBento(t *testing.T) {
	// Bento with missing required URL parameter
	bentoFile := createTestBento(t, "invalid.bento.json", `{
		"id": "invalid",
		"type": "group",
		"version": "1.0.0",
		"nodes": [{
			"id": "http-1",
			"type": "http-request",
			"version": "1.0.0",
			"parameters": {"method": "GET"}
		}],
		"edges": []
	}`)
	defer os.Remove(bentoFile)

	cmd := exec.Command("bento", "inspect", bentoFile)
	output, err := cmd.CombinedOutput()

	// Should exit with error
	if err == nil {
		t.Fatal("inspect should fail on invalid bento")
	}

	outputStr := string(output)

	// Should mention the missing URL
	if !strings.Contains(outputStr, "url") {
		t.Errorf("Error should mention missing 'url': %s", outputStr)
	}

	// Should mention node ID
	if !strings.Contains(outputStr, "http-1") {
		t.Errorf("Error should mention node ID 'http-1': %s", outputStr)
	}
}
```

### Step 3: Test bento menu command

```go
// Test: bento menu should list all bentos in directory
func TestMenuCommand_ListBentos(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test bentos
	createTestBentoInDir(t, tmpDir, "workflow1.bento.json", "Workflow 1")
	createTestBentoInDir(t, tmpDir, "workflow2.bento.json", "Workflow 2")

	cmd := exec.Command("bento", "menu", tmpDir)
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("menu command failed: %v\nOutput: %s", err, string(output))
	}

	outputStr := string(output)

	// Should list both bentos
	if !strings.Contains(outputStr, "workflow1.bento.json") {
		t.Error("Output should list workflow1.bento.json")
	}

	if !strings.Contains(outputStr, "workflow2.bento.json") {
		t.Error("Output should list workflow2.bento.json")
	}

	// Should show count
	if !strings.Contains(outputStr, "2 bentos found") {
		t.Error("Output should show '2 bentos found'")
	}
}
```

### Step 4: Test bento new command

```go
// Test: bento new should create a template bento
func TestNewCommand_CreateTemplate(t *testing.T) {
	tmpDir := t.TempDir()
	bentoPath := filepath.Join(tmpDir, "my-workflow.bento.json")

	// Change to temp dir for test
	oldDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldDir)

	cmd := exec.Command("bento", "new", "my-workflow")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("new command failed: %v\nOutput: %s", err, string(output))
	}

	// Verify file was created
	if _, err := os.Stat(bentoPath); os.IsNotExist(err) {
		t.Fatal("Bento file was not created")
	}

	// Read and verify it's valid JSON
	content, _ := os.ReadFile(bentoPath)
	var def map[string]interface{}
	if err := json.Unmarshal(content, &def); err != nil {
		t.Errorf("Created bento is not valid JSON: %v", err)
	}

	// Should have correct structure
	if def["id"] != "my-workflow" {
		t.Errorf("Bento ID = %v, want my-workflow", def["id"])
	}
}
```

---

## File Structure

```
cmd/bento/
├── main.go              # Entry point, Cobra root command (~150 lines)
├── serve.go             # serve command implementation (~200 lines)
├── inspect.go           # inspect command implementation (~150 lines)
├── menu.go              # menu command implementation (~150 lines)
├── new.go               # new command implementation (~150 lines)
├── config.go            # Configuration (Viper) (~100 lines)
├── output.go            # Pretty output formatting (~150 lines)
└── version.go           # Version command (~50 lines)

cmd/bento/
└── *_test.go            # Integration tests (~400 lines total)
```

---

## Implementation Guidance

**File: `cmd/bento/main.go`**

```go
// Package main implements the bento CLI.
//
// Bento is a high-performance workflow automation CLI written in Go.
// It uses playful sushi-themed commands to make automation fun.
//
// Commands:
//   - serve: Execute a bento workflow
//   - inspect: Validate a bento without executing
//   - menu: List available bentos
//   - new: Create a new bento template
//
// Learn more: https://github.com/Develonaut/bento
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev" // Set by build process

var rootCmd = &cobra.Command{
	Use:   "bento",
	Short: "🍱 High-performance workflow automation",
	Long: `Bento - Workflow automation with a taste of sushi 🍱

Bento lets you build powerful automation workflows using composable
"neta" (ingredients) that can be connected together like a carefully
crafted bento box.

Commands are playfully themed:
  • serve   - Execute a bento (like serving the finished dish)
  • inspect - Validate without executing (quality check)
  • menu    - List available bentos (restaurant menu)
  • new     - Create a new bento template`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(inspectCmd)
	rootCmd.AddCommand(menuCmd)
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(versionCmd)
}
```

**File: `cmd/bento/serve.go`**

```go
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Develonaut/bento/pkg/hangiri"
	"github.com/Develonaut/bento/pkg/itamae"
	"github.com/Develonaut/bento/pkg/pantry"
	"github.com/Develonaut/bento/pkg/shoyu"
	"github.com/spf13/cobra"
)

var (
	verboseFlag bool
	timeoutFlag time.Duration
)

var serveCmd = &cobra.Command{
	Use:   "serve [file].bento.json",
	Short: "🍱 Serve a bento (execute workflow)",
	Long: `Execute a bento workflow from start to finish.

Like a sushi chef serving the finished dish, this command executes
all neta in the bento and reports progress.

Examples:
  bento serve workflow.bento.json
  bento serve workflow.bento.json --verbose
  bento serve workflow.bento.json --timeout 30m`,
	Args: cobra.ExactArgs(1),
	RunE: runServe,
}

func init() {
	serveCmd.Flags().BoolVarP(&verboseFlag, "verbose", "v", false, "Verbose output")
	serveCmd.Flags().DurationVar(&timeoutFlag, "timeout", 10*time.Minute, "Execution timeout")
}

func runServe(cmd *cobra.Command, args []string) error {
	bentoPath := args[0]

	// Load bento
	h := hangiri.New()
	def, err := h.LoadBento(bentoPath)
	if err != nil {
		return fmt.Errorf("failed to load bento: %w", err)
	}

	fmt.Printf("🍱 Serving bento: %s\n", def.Name)

	// Create logger
	logLevel := shoyu.LevelInfo
	if verboseFlag {
		logLevel = shoyu.LevelDebug
	}

	logger := shoyu.New(shoyu.Config{
		Level:  logLevel,
		Format: shoyu.FormatConsole,
	})

	// Create itamae
	p := pantry.New()
	chef := itamae.New(p, logger)

	// Progress callback
	chef.OnProgress(func(nodeID, status string) {
		if status == "starting" {
			fmt.Printf("🍙 Executing neta '%s'...\n", nodeID)
		} else if status == "completed" {
			fmt.Printf("✓ Completed\n")
		}
	})

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeoutFlag)
	defer cancel()

	// Execute
	start := time.Now()
	result, err := chef.Serve(ctx, def)
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("\n❌ Bento execution failed: %v\n", err)
		return err
	}

	fmt.Printf("\n🍱 Bento served successfully in %v\n", duration)
	fmt.Printf("   %d neta executed\n", result.NodesExecuted)

	return nil
}
```

**File: `cmd/bento/output.go`**

```go
package main

import (
	"fmt"
	"time"
)

// formatDuration formats a duration for human readability.
func formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm %ds", int(d.Minutes()), int(d.Seconds())%60)
	}
	return fmt.Sprintf("%dh %dm", int(d.Hours()), int(d.Minutes())%60)
}

// successBox prints a success message in a box.
func successBox(message string) {
	fmt.Printf("\n")
	fmt.Printf("╭─────────────────────────────────────────────╮\n")
	fmt.Printf("│  ✓  %-40s │\n", message)
	fmt.Printf("╰─────────────────────────────────────────────╯\n")
}

// errorBox prints an error message in a box.
func errorBox(message string) {
	fmt.Printf("\n")
	fmt.Printf("╭─────────────────────────────────────────────╮\n")
	fmt.Printf("│  ✗  %-40s │\n", message)
	fmt.Printf("╰─────────────────────────────────────────────╯\n")
}
```

---

## Common Go Pitfalls to Avoid

1. **Cobra command setup**: AddCommand in init(), not main()
   ```go
   // ❌ BAD - commands won't be registered
   func main() {
       rootCmd.AddCommand(serveCmd)
       rootCmd.Execute()
   }

   // ✅ GOOD - use init()
   func init() {
       rootCmd.AddCommand(serveCmd)
   }
   ```

2. **Exit codes**: Return error from RunE, don't call os.Exit()
   ```go
   // ❌ BAD - Cobra can't handle the error
   func runServe(cmd *cobra.Command, args []string) error {
       if err != nil {
           os.Exit(1)
       }
   }

   // ✅ GOOD - return error
   func runServe(cmd *cobra.Command, args []string) error {
       if err != nil {
           return err
       }
   }
   ```

---

## Critical for Phase 8

**Progress Output:**
- `bento serve` must show real-time progress for long bentos
- "⟳ Rendering product 3/50... 45% complete"
- Streaming output from shell-command neta

**Timeout:**
- Default 10min timeout is too short for Phase 8 (Blender takes minutes per product)
- Use `--timeout 30m` or `--timeout 1h`

---

## Bento Box Principle Checklist

- [ ] Files < 250 lines (split into serve.go, inspect.go, menu.go, new.go)
- [ ] Functions < 20 lines
- [ ] Single responsibility per file
- [ ] User-friendly output (emojis, clear messages)
- [ ] File-level documentation

---

## Phase Completion

**Phase 7 MUST end with:**

1. All tests passing (`go test ./cmd/bento/...`)
2. Build binary (`go build -o bento ./cmd/bento`)
3. Manual testing of all commands
4. Run `/code-review` slash command
5. Address feedback from Karen and Colossus
6. Get explicit approval from both agents

**Do not proceed to Phase 8 until code review is approved.**

---

## Claude Prompt Template

```
I need to implement Phase 7: cmd/bento (CLI) following TDD principles.

This brings everything together into a user-facing tool!

Please read:
- .claude/strategy/phase-7-cmd-bento.md (this file)
- .claude/BENTO_BOX_PRINCIPLE.md
- .claude/EMOJIS.md (for approved emoji usage)

Then:

1. Create integration tests in cmd/bento/*_test.go for:
   - serve command (execute bento, show progress, handle errors)
   - inspect command (validate without executing)
   - menu command (list bentos in directory)
   - new command (create template bento)
   - Flags (--verbose, --timeout, --recursive, etc.)
   - Exit codes (0 for success, 1 for errors)

2. Watch the tests fail

3. Implement to make tests pass:
   - cmd/bento/main.go (~150 lines) - Cobra root
   - cmd/bento/serve.go (~200 lines) - serve command
   - cmd/bento/inspect.go (~150 lines) - inspect command
   - cmd/bento/menu.go (~150 lines) - menu command
   - cmd/bento/new.go (~150 lines) - new command
   - cmd/bento/output.go (~150 lines) - Pretty output formatting

4. Build and test manually: `go build -o bento ./cmd/bento && ./bento serve test.bento.json`

Remember:
- Write tests FIRST
- Files < 250 lines
- User-friendly output with emojis 🍱🍙
- Clear error messages
- Return errors from RunE (don't os.Exit)

When complete, run `/code-review` and get Karen + Colossus approval.
```

---

## Dependencies to Add

```bash
go get github.com/spf13/cobra
go get github.com/spf13/viper
```

---

## Notes

- This is what users interact with - make it delightful!
- Sushi emojis make the CLI feel playful and unique
- Clear error messages are critical (users won't have code access)
- Progress output is essential for long-running bentos (Phase 8)
- Exit codes matter for CI/CD integration

---

**Status:** Ready for implementation
**Next Phase:** Phase 8 (real-world integration test) - depends on completion of Phase 7
