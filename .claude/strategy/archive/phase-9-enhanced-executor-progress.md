# Phase 9: Enhanced Executor - Real-time Progress Tracking

**Status**: Pending
**Duration**: 5-6 hours
**Prerequisites**: Phase 8 complete, Karen approved

## Overview

Add real-time per-node progress tracking to the executor using Bubble Tea's messaging system. Users see exactly which nodes are running, completed, or failed, with individual timing for each step. This phase implements the ProgressMessenger interface in itamae and handles recursive bento trees for complete execution visibility.

**Phase 8 Executor (Simple Progress):**
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
│ [JSON output]                       │
│                                     │
│ ✓ Success                           │
│ Execution time: 1.234s              │
│ [Progress Bar ██████████] 100%     │
└─────────────────────────────────────┘
```

**Phase 9 Executor (Per-Node Progress):**
```
┌─────────────────────────────────────┐
│ Bento Executor                      │
├─────────────────────────────────────┤
│ Bento: hello-world                  │
│ Path: ~/.bento/bentos/hello-world   │
│                                     │
│ 🍱 Unpacking bento...                │
│ ✓ Get Hello Message (234ms)         │
│ ⏳ Extract Slideshow Title           │  ← Real-time!
│ • Echo Result                       │
│                                     │
│ [JSON output when complete]         │
│                                     │
│ ✓ Success                           │
│ Execution time: 1.234s              │
│ [Progress Bar ████████░░] 66%      │  ← Accurate!
└─────────────────────────────────────┘
```

**Key Innovation**: Bubble Tea's `Program.Send()` allows itamae (running in a goroutine) to send messages directly to the TUI update loop, enabling real-time progress updates without polling or race conditions.

## Pre-Work Checklist

Before starting, you MUST:

1. ✅ Read [BENTO_BOX_PRINCIPLE.md](../BENTO_BOX_PRINCIPLE.md)
2. ✅ Read [CHARM_STACK_GUIDE.md](../CHARM_STACK_GUIDE.md)
3. ✅ Confirm: "I understand the Bento Box Principle and will follow it"
4. ✅ Use TodoWrite to track all tasks
5. ✅ Phase 8 approved by Karen
6. ✅ Understand Bubble Tea messaging (read send-msg example)
7. ✅ Understand recursion handling for nested bentos

## Goals

1. Add ProgressMessenger interface to itamae
2. Emit NodeStarted/NodeCompleted events during execution
3. Implement ExecutorMessenger adapter using Program.Send()
4. Track NodeState for each node in bento
5. Display real-time node status with spinners/checkmarks/X
6. Show individual node timing
7. Handle recursive bento trees (nested groups, bento.execute)
8. Calculate accurate progress based on node completion
9. Support viewport scrolling for long node lists
10. Maintain Bento Box compliance (< 300 lines per file)

## Architecture: Bubble Tea Messaging Pattern

### Message Flow

```
┌──────────────┐                    ┌──────────────────┐
│   Executor   │                    │  ExecuteBentoCmd │
│   (TUI)      │                    │   (goroutine)    │
└──────┬───────┘                    └────────┬─────────┘
       │                                     │
       │ 1. ExecuteCmd()                    │
       │───────────────────────────────────>│
       │                                     │
       │                                     │ 2. Load bento
       │                                     │ 3. Create chef with messenger
       │                                     │ 4. chef.Execute(def)
       │                                     │
       │                                     ▼
       │                            ┌──────────────┐
       │                            │   Itamae     │
       │                            └──────┬───────┘
       │                                   │
       │ 5. NodeStartedMsg                 │ For each node:
       │<──────────────────────────────────┤ - Send NodeStarted
       │                                   │ - Execute node
       │ 6. Update UI                      │ - Send NodeCompleted
       │                                   │
       │ 7. NodeCompletedMsg               │
       │<──────────────────────────────────┤
       │                                   │
       │ 8. Update UI                      │
       │                                   │
       │ 9. ExecutionCompleteMsg           │
       │<──────────────────────────────────┘
       │
       │ 10. Show final result
       ▼
```

### Key Pattern: Program.Send()

**From Bubble Tea send-msg example:**

```go
// In background goroutine
program.Send(customMsg{data: "something"})

// TUI Update() receives it
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case customMsg:
        // Handle the message
        m.data = msg.data
        return m, nil
    }
}
```

**Thread-safe**: `Program.Send()` is safe to call from goroutines
**Queue-based**: Messages are queued and processed in order
**Non-blocking**: Sending doesn't wait for processing

## ProgressMessenger Interface

### Interface Definition

**File**: `pkg/itamae/itamae.go` (MODIFY)

```go
// ProgressMessenger receives execution progress events
// Used by TUI to display real-time progress
// Optional - nil check before use for non-TUI execution
type ProgressMessenger interface {
	// SendNodeStarted notifies that a node has started execution
	// path: node path in tree (e.g., "0", "1.2", "0.1.3")
	// name: human-readable node name
	// nodeType: node type (e.g., "http", "transform.jq")
	SendNodeStarted(path, name, nodeType string)

	// SendNodeCompleted notifies that a node has finished execution
	// path: node path in tree
	// duration: how long the node took to execute
	// err: error if node failed, nil if successful
	SendNodeCompleted(path string, duration time.Duration, err error)
}

// Itamae orchestrates the execution of neta definitions
type Itamae struct {
	pantry    Registry
	messenger ProgressMessenger // Optional - can be nil
}

// NewWithMessenger creates Itamae with progress messaging
func NewWithMessenger(registry Registry, messenger ProgressMessenger) *Itamae {
	return &Itamae{
		pantry:    registry,
		messenger: messenger,
	}
}
```

**Design Decisions**:
- ✅ **Optional messenger** - Nil check allows non-TUI usage
- ✅ **Path-based identification** - Supports nested structures
- ✅ **Simple interface** - Only 2 methods needed
- ✅ **No return values** - Fire-and-forget messaging
- ✅ **Error in completion** - Single place for success/failure

### Modified Execution Logic

**File**: `pkg/itamae/itamae.go` (MODIFY)

```go
// executeGroup runs a group of nodes in sequence
func (i *Itamae) executeGroup(ctx context.Context, def neta.Definition, basePath string) (neta.Result, error) {
	results := make([]neta.Result, 0, len(def.Nodes))

	for idx, child := range def.Nodes {
		// Build path for this node
		var nodePath string
		if basePath == "" {
			nodePath = fmt.Sprintf("%d", idx)
		} else {
			nodePath = fmt.Sprintf("%s.%d", basePath, idx)
		}

		// Emit start message
		if i.messenger != nil {
			i.messenger.SendNodeStarted(nodePath, child.Name, child.Type)
		}

		// Execute node with timing
		start := time.Now()
		result, err := i.Execute(ctx, child)
		duration := time.Since(start)

		// Emit completion message
		if i.messenger != nil {
			i.messenger.SendNodeCompleted(nodePath, duration, err)
		}

		// Handle error
		if err != nil {
			return neta.Result{}, err
		}

		results = append(results, result)
	}

	return neta.Result{Output: results}, nil
}

// executeSingle runs a single node
func (i *Itamae) executeSingle(ctx context.Context, def neta.Definition, path string) (neta.Result, error) {
	// Emit start message
	if i.messenger != nil {
		i.messenger.SendNodeStarted(path, def.Name, def.Type)
	}

	// Get executable
	exec, err := i.pantry.Get(def.Type)
	if err != nil {
		// Emit failure
		if i.messenger != nil {
			i.messenger.SendNodeCompleted(path, 0, err)
		}
		return neta.Result{}, fmt.Errorf("node type not found: %s: %w", def.Type, err)
	}

	// Execute with timing
	start := time.Now()
	result, err := exec.Execute(ctx, def.Parameters)
	duration := time.Since(start)

	// Emit completion
	if i.messenger != nil {
		i.messenger.SendNodeCompleted(path, duration, err)
	}

	return result, err
}
```

**Bento Box Compliance**:
- ✅ Messenger optional (nil checks)
- ✅ Path generation isolated
- ✅ Timing logic clear
- ✅ Error handling preserved

## Node Path System

### Why Paths?

Bentos can be deeply nested:
```yaml
nodes:
  - name: Step 1        # path: "0"
    nodes:
      - name: Sub A     # path: "0.0"
      - name: Sub B     # path: "0.1"
        nodes:
          - name: C     # path: "0.1.0"
  - name: Step 2        # path: "1"
```

Paths uniquely identify each node for UI updates.

### Path Format

- **Root nodes**: `"0"`, `"1"`, `"2"`
- **First level children**: `"0.0"`, `"0.1"`, `"0.2"`
- **Nested**: `"0.1.2"` = 3rd child of 2nd child of 1st node
- **Deep nesting**: `"0.1.2.3.4"` = 5 levels deep

### Implementation

```go
// buildPath constructs node path
func buildPath(basePath string, index int) string {
	if basePath == "" {
		return fmt.Sprintf("%d", index)
	}
	return fmt.Sprintf("%s.%d", basePath, index)
}

// parseDepth calculates nesting level from path
func parseDepth(path string) int {
	if path == "" {
		return 0
	}
	return strings.Count(path, ".") + 1
}
```

## ExecutorMessenger Adapter

### Implementation

**File**: `pkg/omise/screens/executor_cmd.go` (MODIFY)

```go
// executorMessenger sends progress messages to TUI
type executorMessenger struct {
	program *tea.Program
}

// SendNodeStarted sends node start message
func (m *executorMessenger) SendNodeStarted(path, name, nodeType string) {
	m.program.Send(NodeStartedMsg{
		Path:     path,
		Name:     name,
		NodeType: nodeType,
	})
}

// SendNodeCompleted sends node completion message
func (m *executorMessenger) SendNodeCompleted(path string, duration time.Duration, err error) {
	m.program.Send(NodeCompletedMsg{
		Path:     path,
		Duration: duration,
		Error:    err,
	})
}

// ExecuteBentoCmd creates command with progress messaging
func ExecuteBentoCmd(bentoName, workDir string, program *tea.Program) tea.Cmd {
	return func() tea.Msg {
		// Load bento
		store, err := jubako.NewStore(workDir)
		if err != nil {
			return ExecutionErrorMsg{Error: err}
		}

		def, err := store.Load(bentoName)
		if err != nil {
			return ExecutionErrorMsg{Error: err}
		}

		// Create registry and chef with messenger
		registry := pantry.New()
		messenger := &executorMessenger{program: program}
		chef := itamae.NewWithMessenger(registry, messenger)

		// Register node types
		registerNodeTypes(registry, chef)

		// Execute with context
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		result, err := chef.Execute(ctx, def)
		if err != nil {
			return ExecutionCompleteMsg{Success: false, Error: err}
		}

		return ExecutionCompleteMsg{Success: true, Result: result}
	}
}

// registerNodeTypes registers all standard node types
func registerNodeTypes(registry *pantry.Pantry, chef *itamae.Itamae) {
	_ = registry.Register("http", http.New())
	_ = registry.Register("transform.jq", transform.NewJQ())
	_ = registry.Register("group.sequence", group.NewSequence(chef))
	_ = registry.Register("group.parallel", group.NewParallel(chef))
	_ = registry.Register("conditional.if", conditional.NewIf(chef))
	_ = registry.Register("loop.for", loop.NewFor(chef))
}
```

**Bento Box Compliance**:
- ✅ Adapter pattern clean
- ✅ Program reference encapsulated
- ✅ Registration logic extracted
- ✅ Functions < 20 lines

## NodeState Tracking

### Data Structure

**File**: `pkg/omise/screens/executor.go` (MODIFY)

```go
// NodeStatus represents node execution state
type NodeStatus int

const (
	NodePending NodeStatus = iota
	NodeRunning
	NodeCompleted
	NodeFailed
)

// NodeState tracks individual node execution
type NodeState struct {
	path      string
	name      string
	nodeType  string
	status    NodeStatus
	startTime time.Time
	duration  time.Duration
	depth     int  // Nesting level for indentation
}

// Executor model additions
type Executor struct {
	// ... existing fields
	nodeStates []NodeState
}

// initializeNodeStates prepares tracking for all nodes
func (e *Executor) initializeNodeStates(def neta.Definition) {
	e.nodeStates = flattenDefinition(def, "")
}

// flattenDefinition converts tree to flat list with paths
func flattenDefinition(def neta.Definition, basePath string) []NodeState {
	states := []NodeState{}

	if def.IsGroup() {
		for idx, child := range def.Nodes {
			path := buildPath(basePath, idx)
			depth := parseDepth(path)

			states = append(states, NodeState{
				path:     path,
				name:     child.Name,
				nodeType: child.Type,
				status:   NodePending,
				depth:    depth,
			})

			// Recursively flatten children
			if child.IsGroup() {
				childStates := flattenDefinition(child, path)
				states = append(states, childStates...)
			}
		}
	} else {
		// Single node bento
		states = append(states, NodeState{
			path:     basePath,
			name:     def.Name,
			nodeType: def.Type,
			status:   NodePending,
			depth:    0,
		})
	}

	return states
}
```

**Design Decisions**:
- ✅ **Pre-flatten tree** - All nodes known before execution
- ✅ **Path-based lookup** - Fast updates by path
- ✅ **Depth tracking** - Enables visual indentation
- ✅ **Status enum** - Clear state machine

### Message Handlers

**File**: `pkg/omise/screens/executor.go` (MODIFY)

```go
// Update handler additions
func (e Executor) Update(msg tea.Msg) (Executor, tea.Cmd) {
	// ... existing handlers

	switch msg := msg.(type) {
	case NodeStartedMsg:
		return e.handleNodeStarted(msg)
	case NodeCompletedMsg:
		return e.handleNodeCompleted(msg)
	// ... other cases
	}

	return e, nil
}

// handleNodeStarted updates node to running state
func (e Executor) handleNodeStarted(msg NodeStartedMsg) (Executor, tea.Cmd) {
	for i := range e.nodeStates {
		if e.nodeStates[i].path == msg.Path {
			e.nodeStates[i].status = NodeRunning
			e.nodeStates[i].startTime = time.Now()
			break
		}
	}
	return e, nil
}

// handleNodeCompleted updates node to completed/failed state
func (e Executor) handleNodeCompleted(msg NodeCompletedMsg) (Executor, tea.Cmd) {
	for i := range e.nodeStates {
		if e.nodeStates[i].path == msg.Path {
			e.nodeStates[i].duration = msg.Duration
			if msg.Error != nil {
				e.nodeStates[i].status = NodeFailed
			} else {
				e.nodeStates[i].status = NodeCompleted
			}
			break
		}
	}

	// Update progress based on completion
	e.progressPercent = e.calculateProgress()

	return e, nil
}

// calculateProgress returns completion percentage
func (e Executor) calculateProgress() float64 {
	if len(e.nodeStates) == 0 {
		return 0.0
	}

	completed := 0
	for _, node := range e.nodeStates {
		if node.status == NodeCompleted || node.status == NodeFailed {
			completed++
		}
	}

	return float64(completed) / float64(len(e.nodeStates))
}
```

**Bento Box Compliance**:
- ✅ Linear search acceptable (small lists)
- ✅ Clear update logic
- ✅ Progress calculation isolated
- ✅ Functions < 20 lines

## Enhanced View Rendering

### Running View with Per-Node Status

**File**: `pkg/omise/screens/executor.go` (MODIFY)

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

	// CENTER SECTION (node progress)
	center := []string{
		emojiBento + " Unpacking bento...",
		"",
	}

	// Render each node with status
	for _, node := range e.nodeStates {
		line := e.formatNodeLine(node)
		center = append(center, line)
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

// formatNodeLine renders single node with status icon and timing
func (e Executor) formatNodeLine(node NodeState) string {
	// Indentation based on depth
	indent := strings.Repeat("  ", node.depth)

	// Status icon
	icon := e.getNodeIcon(node.status)

	// Base line
	line := fmt.Sprintf("%s%s %s", indent, icon, node.name)

	// Add duration if completed
	if node.status == NodeCompleted || node.status == NodeFailed {
		durationStr := node.duration.Round(time.Millisecond).String()
		line = fmt.Sprintf("%s (%s)", line, durationStr)
	}

	// Add type hint if running
	if node.status == NodeRunning {
		line = fmt.Sprintf("%s [%s]", line, node.nodeType)
	}

	// Style based on status
	switch node.status {
	case NodeCompleted:
		return styles.SuccessStyle.Render(line)
	case NodeFailed:
		return styles.ErrorStyle.Render(line)
	case NodeRunning:
		return line // Default style with spinner
	default:
		return styles.Subtle.Render(line)
	}
}

// getNodeIcon returns emoji/character for node status
func (e Executor) getNodeIcon(status NodeStatus) string {
	switch status {
	case NodeRunning:
		return e.spinner.View() // Animated spinner
	case NodeCompleted:
		return emojiSuccess // ✓
	case NodeFailed:
		return emojiFailure // ✗
	default:
		return "•" // Pending
	}
}
```

**Visual Example**:
```
🍱 Unpacking bento...

✓ Get Hello Message (234ms)
⏳ Extract Slideshow Title [transform.jq]
• Echo Result
  • Nested Step
```

**Bento Box Compliance**:
- ✅ Formatting logic isolated
- ✅ Clear status → icon mapping
- ✅ Indentation based on depth
- ✅ Functions < 20 lines

## Message Types

### New Messages

**File**: `pkg/omise/screens/executor_messages.go` (MODIFY)

```go
// NodeStartedMsg signals a node has started execution
type NodeStartedMsg struct {
	Path     string
	Name     string
	NodeType string
}

// NodeCompletedMsg signals a node has finished execution
type NodeCompletedMsg struct {
	Path     string
	Duration time.Duration
	Error    error
}
```

**Simple and clear** - Just what we need for progress tracking

## Threading Program Reference

### Problem

`ExecuteBentoCmd` needs `*tea.Program` to create messenger, but it's called from `Executor.ExecuteCmd()` which doesn't have it.

### Solution

Pass program reference through the call chain.

**File**: `pkg/omise/update.go` (MODIFY)

```go
// handleWorkflowSelected passes program reference
func (m Model) handleWorkflowSelected(msg screens.WorkflowSelectedMsg) (tea.Model, tea.Cmd) {
	m.screen = ScreenExecutor
	m.executor = m.executor.StartBento(msg.Name, msg.Path, m.workDir)

	// Get program reference from context or store in Model
	// For now, we'll add a method to Model
	return m, m.executor.ExecuteCmdWithProgram(m.program)
}
```

**File**: `pkg/omise/model.go` (MODIFY)

```go
type Model struct {
	// ... existing fields
	program *tea.Program // Set in main
}

// SetProgram stores program reference for messaging
func (m *Model) SetProgram(p *tea.Program) {
	m.program = p
}
```

**File**: `cmd/bento/main.go` (MODIFY)

```go
func main() {
	// ... existing setup
	model := omise.NewModel()

	program := tea.NewProgram(model, tea.WithAltScreen())

	// Set program reference for messaging
	model.SetProgram(program)

	if err := program.Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
```

**Alternative**: Use tea.WindowSizeMsg to capture program in Update

```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Bubble Tea passes program in context
		// We can capture it here if needed
	}
}
```

**Simplest**: Store in Model at creation

## File Structure

```
pkg/itamae/
├── itamae.go             # MODIFY - Add ProgressMessenger interface
└── itamae_test.go        # MODIFY - Test messaging

pkg/omise/screens/
├── executor.go           # MODIFY - Add NodeState, handlers, rendering
├── executor_messages.go  # MODIFY - Add NodeStarted/CompletedMsg
└── executor_cmd.go       # MODIFY - Add ExecutorMessenger adapter

pkg/omise/
├── model.go              # MODIFY - Add program field
└── update.go             # MODIFY - Thread program reference

cmd/bento/
└── main.go               # MODIFY - Set program reference

No new files!
```

**Size Targets**:
- `itamae/itamae.go`: < 150 lines (currently ~60)
- `executor.go`: < 350 lines (currently ~300)
- `executor_cmd.go`: < 200 lines (currently ~145)
- `executor_messages.go`: < 100 lines (currently ~30)

## Testing Strategy

### Unit Tests

**File**: `pkg/itamae/itamae_test.go`

```go
// testMessenger implements ProgressMessenger for testing
type testMessenger struct {
	started   []string
	completed []string
}

func (t *testMessenger) SendNodeStarted(path, name, nodeType string) {
	t.started = append(t.started, path)
}

func (t *testMessenger) SendNodeCompleted(path string, duration time.Duration, err error) {
	t.completed = append(t.completed, path)
}

func TestItamae_WithMessenger(t *testing.T) {
	registry := pantry.New()
	messenger := &testMessenger{}
	chef := itamae.NewWithMessenger(registry, messenger)

	// Create test definition with 3 nodes
	def := neta.Definition{
		Type: "group.sequence",
		Nodes: []neta.Definition{
			{Name: "Node 1", Type: "test"},
			{Name: "Node 2", Type: "test"},
			{Name: "Node 3", Type: "test"},
		},
	}

	// Execute
	_, _ = chef.Execute(context.Background(), def)

	// Verify messages sent
	if len(messenger.started) != 3 {
		t.Errorf("expected 3 started messages, got %d", len(messenger.started))
	}
	if len(messenger.completed) != 3 {
		t.Errorf("expected 3 completed messages, got %d", len(messenger.completed))
	}

	// Verify paths
	expectedPaths := []string{"0", "1", "2"}
	for i, path := range messenger.started {
		if path != expectedPaths[i] {
			t.Errorf("path %d: expected %s, got %s", i, expectedPaths[i], path)
		}
	}
}

func TestItamae_WithoutMessenger(t *testing.T) {
	// Ensure nil messenger doesn't crash
	registry := pantry.New()
	chef := itamae.New(registry) // No messenger

	def := neta.Definition{Type: "test"}

	// Should not crash
	_, _ = chef.Execute(context.Background(), def)
}
```

### Integration Tests

**Test with hello-world bento**:

```bash
# Run TUI
./bento

# In browser:
# 1. Select hello-world
# 2. Press 'r' to run
# 3. Verify:
#    - Each node shows "•" pending initially
#    - Nodes update to spinner when running
#    - Nodes show checkmark + timing when complete
#    - Progress bar advances with each completion
#    - Final progress is 100%
```

**Test with nested bento**:

Create test bento with nested groups:
```yaml
version: "1.0"
type: group.sequence
name: Nested Test
nodes:
  - type: group.sequence
    name: Group 1
    nodes:
      - type: http
        name: Request A
      - type: http
        name: Request B
  - type: http
    name: Request C
```

Verify:
- ✅ Nested nodes show with indentation
- ✅ Paths are correct ("0.0", "0.1", "1")
- ✅ All nodes tracked individually
- ✅ Progress accurate across levels

## Success Criteria

Phase 9 is complete when:

1. ✅ ProgressMessenger interface added to itamae
2. ✅ NewWithMessenger constructor created
3. ✅ Itamae emits NodeStarted messages
4. ✅ Itamae emits NodeCompleted messages
5. ✅ Node paths generated correctly
6. ✅ ExecutorMessenger adapter implemented
7. ✅ Program reference threaded through
8. ✅ NodeState tracking implemented
9. ✅ Tree flattening works for nested bentos
10. ✅ Message handlers update node states
11. ✅ Per-node view rendering working
12. ✅ Spinners show for running nodes
13. ✅ Checkmarks show for completed nodes
14. ✅ X shows for failed nodes
15. ✅ Individual node timing displayed
16. ✅ Indentation based on nesting depth
17. ✅ Progress bar accurate to node completion
18. ✅ Long node lists scroll in viewport
19. ✅ All files < 350 lines
20. ✅ All functions < 20 lines
21. ✅ Tests passing (unit + integration)
22. ✅ **Karen's approval granted**

## Common Pitfalls to Avoid

1. ❌ **Not checking for nil messenger**
   - Always: `if i.messenger != nil {`
   - Itamae must work without TUI

2. ❌ **Race conditions with Program.Send()**
   - Don't worry! `Program.Send()` is thread-safe
   - It's designed for this exact use case

3. ❌ **Not pre-flattening tree**
   - Flatten definition before execution starts
   - Don't try to build states on-the-fly

4. ❌ **Progress calculation errors**
   - Count only leaf nodes OR
   - Count all nodes including groups
   - Be consistent!

5. ❌ **Path generation bugs**
   - Test with nested bentos early
   - Verify paths match between itamae and executor

6. ❌ **Not testing with real nested bentos**
   - Test hello-world (flat)
   - Test with group.sequence (one level)
   - Test with bento.execute (recursion)

7. ❌ **Forgetting to update spinner**
   - Spinner needs tick messages
   - Keep spinner.Update() in main Update loop

8. ❌ **Thread-safety paranoia**
   - Bubble Tea handles concurrency
   - Just use Program.Send()
   - Don't add mutexes or channels

## Validation Commands

```bash
# Build
go build -o bento ./cmd/bento

# Test
go test -v ./pkg/itamae/
go test -v ./pkg/omise/screens/

# Run and verify
./bento

# Test cases:
# 1. hello-world (3 sequential nodes)
# 2. Create nested bento
# 3. Verify per-node progress
# 4. Verify timing accuracy
# 5. Test failure case (bad URL)
# 6. Verify failed node shows X
# 7. Check viewport scrolling with many nodes

# Check file sizes
wc -l pkg/itamae/itamae.go              # < 150
wc -l pkg/omise/screens/executor.go     # < 350
wc -l pkg/omise/screens/executor_cmd.go # < 200
```

## Integration with Phase 8

Phase 9 **extends** Phase 8, not replaces it:

### Unchanged from Phase 8
- ✅ Three-section layout
- ✅ Sushi emoji lifecycle (🍱 ⏳ ✓ ✗)
- ✅ Glamour JSON rendering
- ✅ Execution timer
- ✅ Viewport integration

### Enhanced in Phase 9
- ✅ Center section now shows per-node progress
- ✅ Progress bar now accurate (was estimated)
- ✅ Timer still shows total time
- ✅ Layout structure unchanged

**Phase 8 + Phase 9 = Complete Enhanced Executor**

## Dependencies

No new dependencies! Uses existing:
- github.com/charmbracelet/bubbletea (for Program.Send)
- github.com/charmbracelet/glamour (from Phase 8)
- github.com/charmbracelet/lipgloss
- github.com/charmbracelet/bubbles/progress
- github.com/charmbracelet/bubbles/spinner

## Next Phase

After Karen approval, proceed to **[Phase 10: Real-World Proof-of-Concept](./phase-10-proof-of-concept.md)** to:

- Build real Etsy product pipeline bento
- Validate composable bento architecture
- Test with actual CSV data and Figma API
- Prove system is production-ready

**Phase 9's per-node progress** makes Phase 10's complex workflow visible and debuggable. Users see exactly where the pipeline is in processing their products!

## Execution Prompt

```
I'm ready to begin Phase 9: Enhanced Executor - Real-time Progress Tracking.

I have read the Bento Box Principle and will follow it.
I have read and understand the Bubble Tea send-msg pattern.

Please implement:
- ProgressMessenger interface in itamae
- Node path system for tree tracking
- ExecutorMessenger adapter with Program.Send()
- NodeState tracking and flattening
- Per-node view rendering with spinners/checkmarks
- Individual node timing display
- Accurate progress calculation

This builds on Phase 8's layout to show real-time execution progress.

Each file < 350 lines, functions < 20 lines. I will use TodoWrite to track progress and get Karen's approval before completing.
```

---

**Phase 9 Enhanced Executor**: Real-time per-node progress with Bubble Tea messaging 🍣⏳✨
