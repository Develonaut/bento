# Phase 8: Enhanced Executor - Foundation & Layout

**Status**: Pending
**Duration**: 4-5 hours
**Prerequisites**: Phase 7 complete, Karen approved

## Overview

Transform the executor screen from a simple progress indicator into a beautiful, informative three-section display with syntax-highlighted output and sushi-themed lifecycle indicators. This phase establishes the foundation for Phase 9's real-time per-node progress tracking.

**Current Executor:**
```
┌─────────────────────────────────────┐
│ Bento Executor                      │
├─────────────────────────────────────┤
│ Bento: hello-world                  │
│ Path: ~/.bento/bentos/hello-world   │
│                                     │
│ ⏳ Executing bento...                │
│ [Progress Bar ████████░░] 80%      │
│                                     │
│ Output: {[{{...raw Go struct...}}]} │
└─────────────────────────────────────┘
```

**New Enhanced Executor:**
```
┌─────────────────────────────────────┐
│ Bento Executor                      │
├─────────────────────────────────────┤
│ Bento: hello-world                  │
│ Path: ~/.bento/bentos/hello-world   │
│                                     │
│ 🍱 Unpacking bento...                │
│ ⏳ Executing bento...                │
│ 🍱 Bento complete!                   │
│                                     │
│ {                                   │
│   "slideshow": {                    │  ← Syntax-highlighted
│     "author": "Yours Truly",        │  ← with colors!
│     "title": "Sample Slide Show"    │
│   }                                 │
│ }                                   │
│                                     │
│ ✓ Success                           │
│ Execution time: 1.234s              │
│ [Progress Bar ██████████] 100%     │
└─────────────────────────────────────┘
```

This phase focuses on **visual design and user experience** without the complexity of per-node tracking. Users get immediate feedback through sushi emojis, beautiful JSON output, and precise timing.

## Pre-Work Checklist

Before starting, you MUST:

1. ✅ Read [BENTO_BOX_PRINCIPLE.md](../BENTO_BOX_PRINCIPLE.md)
2. ✅ Read [CHARM_STACK_GUIDE.md](../CHARM_STACK_GUIDE.md)
3. ✅ Confirm: "I understand the Bento Box Principle and will follow it"
4. ✅ Use TodoWrite to track all tasks
5. ✅ Phase 7 approved by Karen

## Goals

1. Implement three-section layout (header, center content, footer)
2. Add execution timer with millisecond precision
3. Create sushi emoji lifecycle system
4. Integrate Glamour for syntax-highlighted JSON output
5. Improve viewport integration for scrollable output
6. Maintain simple progress tracking (no per-node yet)
7. Ensure all files remain < 300 lines
8. Validate Bento Box compliance

## Three-Section Layout Design

### Section 1: Header (Bento Information)
```
Bento: hello-world
Path: ~/.bento/bentos/hello-world.bento.yaml
```

**Purpose**: Show what's being executed
**Always visible**: Top of viewport

### Section 2: Center (Execution Content)
```
🍱 Unpacking bento...
⏳ Executing bento...
🍱 Bento complete!

[Syntax-highlighted JSON output]
```

**Purpose**: Real-time status and results
**Scrollable**: Can grow beyond viewport height
**Dynamic**: Changes based on execution state

### Section 3: Footer (Status & Metrics)
```
✓ Success
Execution time: 1.234s
[Progress Bar ██████████] 100%
```

**Purpose**: Final status and performance metrics
**Always visible**: Bottom of viewport

## Sushi Emoji Lifecycle System

### Lifecycle States

```
State Machine:
idle → unpacking → executing → complete/failed
```

### Emoji Meanings

| Emoji | State | Usage |
|-------|-------|-------|
| 🍱 | Bento lifecycle | Start ("Unpacking bento...") and end ("Bento complete!") |
| ⏳ | Executing | Currently running ("Executing bento...") |
| ✓ | Success | Completed successfully |
| ✗ | Failure | Completed with error |
| 🍣 | Reserved | Future per-node indicator (Phase 9) |
| 🍙 | Reserved | Future transform indicator (Phase 9) |
| 🥢 | Reserved | Future preparation indicator (Phase 9) |

### Implementation Pattern

**File**: `pkg/omise/screens/executor.go`

```go
// Emoji constants for lifecycle states
const (
	emojiBento     = "🍱"
	emojiExecuting = "⏳"
	emojiSuccess   = "✓"
	emojiFailure   = "✗"
)

// getLifecycleEmoji returns emoji for current state
func (e Executor) getLifecycleEmoji() string {
	if !e.running && !e.complete {
		return emojiBento
	}
	if e.running {
		return emojiExecuting
	}
	if e.success {
		return emojiSuccess
	}
	return emojiFailure
}

// getLifecycleMessage returns status message with emoji
func (e Executor) getLifecycleMessage() string {
	if !e.running && !e.complete {
		return emojiBento + " Ready to execute"
	}
	if e.running {
		return emojiExecuting + " Executing bento..."
	}
	if e.success {
		return emojiSuccess + " Bento complete!"
	}
	return emojiFailure + " Bento failed"
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: Emoji selection
- ✅ Functions < 20 lines
- ✅ Clear naming
- ✅ No side effects

## Glamour Integration for JSON

### Why Glamour?

- ✅ **Charm ecosystem** - Same family as Bubble Tea and Lipgloss
- ✅ **Markdown-based** - Simple to use with code fences
- ✅ **Syntax highlighting** - Beautiful JSON colors
- ✅ **Fallback graceful** - Returns plain text on error
- ✅ **Auto-styling** - Adapts to terminal capabilities

### Implementation Pattern

**File**: `pkg/omise/screens/executor.go`

```go
import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/glamour"
	"bento/pkg/neta"
)

// formatResult formats execution result with syntax highlighting
func formatResult(result interface{}) string {
	if result == nil {
		return "No output"
	}

	// Type assert to neta.Result
	netaResult, ok := result.(neta.Result)
	if !ok {
		return fmt.Sprintf("%v", result)
	}

	// Handle nil output
	if netaResult.Output == nil {
		return "No output"
	}

	// Marshal to pretty JSON
	jsonBytes, err := json.MarshalIndent(netaResult.Output, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", netaResult.Output)
	}

	// Wrap in markdown code fence for Glamour
	markdown := fmt.Sprintf("```json\n%s\n```", string(jsonBytes))

	// Render with Glamour
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		return string(jsonBytes) // Fallback to plain JSON
	}

	highlighted, err := renderer.Render(markdown)
	if err != nil {
		return string(jsonBytes) // Fallback to plain JSON
	}

	return highlighted
}
```

**Key Features**:
- ✅ Triple fallback: Glamour → plain JSON → raw struct
- ✅ Auto-styling adapts to terminal
- ✅ Word wrap at 80 characters
- ✅ Clean error handling

**Bento Box Compliance**:
- ✅ Single responsibility: Format output
- ✅ Function < 20 lines (with helpers)
- ✅ Clear error path
- ✅ No side effects

## Execution Timer

### Implementation

**Add to Executor struct**:
```go
type Executor struct {
	// ... existing fields
	startTime time.Time
	endTime   time.Time
}
```

**Track timing in StartBento**:
```go
func (e Executor) StartBento(name, path, workDir string) Executor {
	// ... existing code
	e.startTime = time.Now()
	e.endTime = time.Time{} // Zero value
	return e
}
```

**Update on completion** (ExecutionCompleteMsg handler):
```go
case ExecutionCompleteMsg:
	e.endTime = time.Now()
	e.running = false
	e.complete = true
	// ... rest of handler
```

**Display in footer**:
```go
func (e Executor) renderFooter() string {
	var elapsed time.Duration
	if !e.endTime.IsZero() {
		elapsed = e.endTime.Sub(e.startTime)
	} else if !e.startTime.IsZero() {
		elapsed = time.Since(e.startTime)
	}

	lines := []string{}

	// Status line
	if e.complete {
		status := emojiSuccess + " Success"
		if !e.success {
			status = emojiFailure + " Failed"
		}
		lines = append(lines, status)
	}

	// Timer line
	if elapsed > 0 {
		timerText := fmt.Sprintf("Execution time: %s",
			elapsed.Round(time.Millisecond))
		lines = append(lines, styles.Subtle.Render(timerText))
	}

	// Progress bar
	lines = append(lines, e.progress.ViewAs(e.progressPercent))

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}
```

**Bento Box Compliance**:
- ✅ Time tracking isolated to timer logic
- ✅ Display logic separate from calculation
- ✅ Clear variable names
- ✅ Millisecond precision

## Updated View Structure

### Idle View (Not Running)

```go
func (e Executor) idleView(title string) string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		styles.Subtle.Render(emojiBento + " Ready to execute bentos"),
		"",
		styles.Subtle.Render("Select a bento from the Browser and press 'r' to run."),
	)
}
```

### Running View (Executing)

```go
func (e Executor) runningView(title string) string {
	// HEADER SECTION
	header := []string{
		title,
		"",
		styles.Subtle.Render("Bento: " + e.bentoName),
		styles.Subtle.Render("Path: " + e.bentoPath),
		"",
	}

	// CENTER SECTION (lifecycle + status)
	center := []string{
		emojiBento + " Unpacking bento...",
		e.spinner.View() + " " + emojiExecuting + " " + e.status,
		"",
	}

	// FOOTER SECTION
	elapsed := time.Since(e.startTime)
	footer := []string{
		"",
		styles.Subtle.Render(fmt.Sprintf("Elapsed: %s",
			elapsed.Round(time.Millisecond))),
		"",
		e.progress.ViewAs(e.progressPercent),
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		append(append(header, center...), footer...)...,
	)
}
```

### Complete View (Finished)

```go
func (e Executor) completeView(title string) string {
	// HEADER SECTION
	header := []string{
		title,
		"",
		styles.Subtle.Render("Bento: " + e.bentoName),
		styles.Subtle.Render("Path: " + e.bentoPath),
		"",
	}

	// CENTER SECTION (lifecycle + output)
	center := []string{
		emojiBento + " Bento complete!",
		"",
	}

	// Add error message if failed
	if !e.success && e.errorMsg != "" {
		center = append(center,
			styles.ErrorStyle.Render("Error:"),
			styles.ErrorStyle.Render(e.errorMsg),
			"",
		)
	}

	// Add output if successful
	if e.success && e.result != "" {
		center = append(center,
			styles.Subtle.Render("Output:"),
			"",
			e.result, // Glamour-highlighted JSON
			"",
		)
	}

	// FOOTER SECTION
	footer := []string{e.renderFooter()}

	// Show copy feedback if present
	if e.copyFeedback != "" {
		footer = append(footer, "",
			styles.SuccessStyle.Render(e.copyFeedback))
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		append(append(header, center...), footer...)...,
	)
}
```

**Bento Box Compliance**:
- ✅ Each view function < 20 lines
- ✅ Clear section separation
- ✅ Consistent structure
- ✅ Easy to maintain

## File Structure

```
pkg/omise/screens/
├── executor.go           # MODIFY - Add timer, layout, Glamour
├── executor_messages.go  # Existing - No changes needed
└── executor_cmd.go       # MODIFY - Remove complex progress tracking

go.mod                    # ADD - github.com/charmbracelet/glamour

Dependencies added:
- github.com/charmbracelet/glamour (latest)
```

**Size Targets**:
- `executor.go`: < 300 lines (currently ~210)
- `executor_cmd.go`: < 150 lines (currently ~145)
- All functions: < 20 lines

## Deliverables

### 1. Enhanced Executor Model

**File**: `pkg/omise/screens/executor.go` (MODIFY)
**Target Size**: < 300 lines

**Changes needed**:
1. Add timing fields:
   ```go
   type Executor struct {
       // ... existing fields
       startTime time.Time
       endTime   time.Time
   }
   ```

2. Add emoji constants (at package level)
3. Add `formatResult()` with Glamour
4. Add `renderFooter()` helper
5. Update `idleView()`, `runningView()`, `completeView()`
6. Update `StartBento()` to set startTime
7. Update `ExecutionCompleteMsg` handler to set endTime

**Bento Box Compliance**:
- ✅ Single responsibility: Bento execution display
- ✅ Clear state machine (idle → running → complete)
- ✅ All functions < 20 lines
- ✅ File < 300 lines

### 2. Simplified Progress Tracking

**File**: `pkg/omise/screens/executor_cmd.go` (MODIFY)
**Target Size**: < 150 lines

**Simplifications**:
1. Remove complex polling loops
2. Keep basic progress: 10% start → 90% during → 100% complete
3. Timer-based increments (simpler than per-node)

**Rationale**: Phase 9 will add proper per-node tracking. For Phase 8, keep it simple and focus on layout/UX.

**Bento Box Compliance**:
- ✅ Simplified logic
- ✅ Clear progression
- ✅ Easy to replace in Phase 9

### 3. Add Glamour Dependency

**File**: `go.mod` (MODIFY)

```bash
go get github.com/charmbracelet/glamour
```

**Version**: Use latest stable (Glamour v0.8.0+)

## Testing Strategy

### Manual Testing

```bash
# Build
go build -o bento ./cmd/bento

# Run TUI
./bento

# Test flow:
# 1. Navigate to Browser
# 2. Select hello-world bento
# 3. Press 'r' to run
# 4. Observe:
#    - Lifecycle emojis (🍱 → ⏳ → 🍱)
#    - Timer counting up during execution
#    - Final timer shows total time
#    - JSON output is syntax-highlighted
#    - Three-section layout is clear
# 5. Try copying output with 'c'
# 6. Verify copy includes formatted JSON
```

### Visual Checks

**Idle State**:
```
✅ Shows "Ready to execute bentos"
✅ Displays 🍱 emoji
✅ Instructions are clear
```

**Running State**:
```
✅ Shows "Unpacking bento..."
✅ Shows spinner + ⏳ + status
✅ Timer updates in real-time
✅ Progress bar increments
```

**Complete State (Success)**:
```
✅ Shows "Bento complete!" with 🍱
✅ JSON is syntax-highlighted (colors!)
✅ Footer shows ✓ Success
✅ Timer shows final time (e.g., "1.234s")
✅ Progress bar at 100%
```

**Complete State (Failure)**:
```
✅ Shows error message in red
✅ Footer shows ✗ Failed
✅ Timer still shows elapsed time
✅ Progress bar state preserved
```

### Unit Tests

**File**: `pkg/omise/screens/executor_test.go`

```go
func TestFormatResult_WithGlamour(t *testing.T) {
	testData := map[string]interface{}{
		"key": "value",
	}

	result := neta.Result{Output: testData}
	formatted := formatResult(result)

	// Should contain JSON
	if !strings.Contains(formatted, "key") {
		t.Error("missing JSON content")
	}
	// Glamour adds ANSI codes (or falls back to plain)
	// Just verify it doesn't crash and returns content
}

func TestExecutorTimer(t *testing.T) {
	e := NewExecutor()
	e = e.StartBento("test", "/path", ".")

	// Start time should be set
	if e.startTime.IsZero() {
		t.Error("start time not set")
	}

	// Simulate completion
	time.Sleep(10 * time.Millisecond)
	e.endTime = time.Now()

	elapsed := e.endTime.Sub(e.startTime)
	if elapsed < 10*time.Millisecond {
		t.Error("elapsed time too short")
	}
}
```

## Success Criteria

Phase 8 is complete when:

1. ✅ Three-section layout implemented (header, center, footer)
2. ✅ Sushi emoji lifecycle working (🍱 → ⏳ → ✓/✗)
3. ✅ Glamour rendering JSON with syntax highlighting
4. ✅ Execution timer displays millisecond-precision timing
5. ✅ Timer counts up during execution
6. ✅ Timer shows final time on completion
7. ✅ Viewport integration working for scrollable output
8. ✅ Long JSON output scrolls properly
9. ✅ Copy function includes formatted JSON
10. ✅ All files < 300 lines
11. ✅ All functions < 20 lines
12. ✅ Tests passing
13. ✅ Manual testing successful
14. ✅ **Karen's approval granted**

## Common Pitfalls to Avoid

1. ❌ **Overcomplicating layout logic**
   - Keep three sections simple
   - Use `lipgloss.JoinVertical` consistently
   - Don't nest too deeply

2. ❌ **Not handling Glamour errors**
   - Always check errors from `NewTermRenderer()`
   - Always check errors from `Render()`
   - Provide fallback to plain JSON

3. ❌ **Hardcoding emojis everywhere**
   - Use constants at package level
   - Create helper functions for emoji selection
   - Keep emoji logic centralized

4. ❌ **Timer drift or inaccuracy**
   - Use `time.Time` not `time.Duration` for tracking
   - Calculate elapsed on-demand
   - Round to milliseconds for display

5. ❌ **Breaking viewport integration**
   - Remember viewport needs content via `SetContent()`
   - Viewport handles scrolling automatically
   - Don't try to manage scroll position manually

6. ❌ **Mixing Phase 8 and Phase 9 concerns**
   - Don't add per-node tracking yet
   - Keep progress simple (0% → 90% → 100%)
   - Save complex logic for Phase 9

## Validation Commands

```bash
# Build and test
go build -o bento ./cmd/bento
go test -v ./pkg/omise/screens/

# Run and verify
./bento

# In TUI:
# 1. Go to Browser
# 2. Run hello-world bento
# 3. Verify layout looks correct
# 4. Verify JSON has colors
# 5. Verify timer shows accurate time
# 6. Copy output and check clipboard

# Check file sizes
wc -l pkg/omise/screens/executor.go        # Should be < 300
wc -l pkg/omise/screens/executor_cmd.go    # Should be < 150

# Check function sizes (none should be > 20 lines)
# Manual review of functions
```

## Integration with Existing Code

### No Breaking Changes

Phase 8 is **purely additive**:
- ✅ Same message types (ExecutionProgressMsg, ExecutionCompleteMsg)
- ✅ Same public API (StartBento, ExecuteCmd, etc.)
- ✅ Same Update/View pattern
- ✅ Backward compatible

### Works with Current Flow

```
Browser → Press 'r' → Executor.StartBento() → ExecuteCmd()
   ↓
ExecuteBentoCmd (goroutine)
   ↓
ExecutionProgressMsg (periodic updates)
   ↓
ExecutionCompleteMsg (when done)
   ↓
Executor shows result with new layout
```

### Prepares for Phase 9

Phase 8 sets foundation for Phase 9:
- ✅ Three-section layout ready for per-node progress
- ✅ Emoji system extensible (🍣 🍙 reserved)
- ✅ Timer architecture supports per-node timing
- ✅ Clean separation of concerns

## Dependencies

### New Dependencies

```toml
# go.mod additions
require (
    github.com/charmbracelet/glamour v0.8.0
)
```

### Existing Dependencies (unchanged)

- github.com/charmbracelet/bubbletea
- github.com/charmbracelet/lipgloss
- github.com/charmbracelet/bubbles/progress
- github.com/charmbracelet/bubbles/spinner
- github.com/charmbracelet/bubbles/viewport

## Next Phase

After Karen approval, proceed to **[Phase 9: Enhanced Executor - Real-time Progress](./phase-9-enhanced-executor-progress.md)** to:

- Add ProgressMessenger interface to itamae
- Emit per-node execution events
- Display real-time node status with spinners/checkmarks
- Show individual node timing
- Handle recursive bento trees
- Calculate accurate progress based on node completion

**Phase 9 builds on Phase 8's foundation**:
- Same three-section layout
- Enhanced center section with per-node details
- More accurate progress bar
- Better user insight into execution

## Execution Prompt

```
I'm ready to begin Phase 8: Enhanced Executor - Foundation & Layout.

I have read the Bento Box Principle and will follow it.

Please enhance the executor with:
- Three-section layout (header, center, footer)
- Sushi emoji lifecycle (🍱 ⏳ ✓ ✗)
- Glamour syntax-highlighted JSON output
- Execution timer with millisecond precision
- Improved viewport integration

Keep progress tracking simple (Phase 9 will add per-node tracking).

Each file < 300 lines, functions < 20 lines. I will use TodoWrite to track progress and get Karen's approval before completing.
```

---

**Phase 8 Enhanced Executor**: Beautiful layout with sushi emojis and syntax-highlighted output 🍱✨
