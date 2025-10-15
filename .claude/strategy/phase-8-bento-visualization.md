# Phase 8: Bento Editor - Visualization & Navigation

**Status**: Pending
**Duration**: 4-5 hours
**Prerequisites**: Phase 7 complete, Karen approved

## Overview

Add visual bento box representation showing compartments (nodes) and enable arrow key navigation to edit, move, or delete individual nodes. Users can also run the bento directly from the editor with progress visualization.

## Pre-Work Checklist

Before starting, you MUST:

1. ✅ Read [BENTO_BOX_PRINCIPLE.md](../BENTO_BOX_PRINCIPLE.md)
2. ✅ Read [CHARM_STACK_GUIDE.md](../CHARM_STACK_GUIDE.md)
3. ✅ Confirm: "I understand the Bento Box Principle and will follow it"
4. ✅ Use TodoWrite to track all tasks
5. ✅ Phase 7 approved by Karen

## Goals

1. Visual bento box showing nodes as compartments
2. Arrow key navigation between nodes
3. Edit node in-place (re-configure parameters)
4. Move node (reorder in sequence)
5. Delete node from bento
6. Run bento from editor with progress
7. Clear visual feedback for selected node
8. Validate Bento Box compliance

## Visualization Design

### Simple Text-Based Bento Box

```
┌─────────────────────────────────────────────┐
│ Bento: My API Workflow (v1.0)               │
├─────────────────────────────────────────────┤
│ ╔═════════════════════════════════════════╗ │
│ ║ 1. Fetch User Data          [http]      ║ │ ← Selected (highlighted)
│ ║    GET https://api.example.com/users    ║ │
│ ╚═════════════════════════════════════════╝ │
│                                             │
│ ┌───────────────────────────────────────┐   │
│ │ 2. Transform Data        [transform.jq]│   │
│ │    .data | map(.id)                    │   │
│ └───────────────────────────────────────┘   │
│                                             │
│ ┌───────────────────────────────────────┐   │
│ │ 3. Filter Active       [conditional.if]│   │
│ │    .status == "active"                 │   │
│ └───────────────────────────────────────┘   │
├─────────────────────────────────────────────┤
│ ↑/↓: Navigate • e: Edit • m: Move •         │
│ d: Delete • r: Run • esc: Back              │
└─────────────────────────────────────────────┘
```

### Future: Visual Compartments (Optional Enhancement)

```
┌────────────────────────────────────┐
│     🍱 My API Workflow (v1.0)      │
├────────────────────────────────────┤
│  ╔══════════╗  ╔══════════╗        │
│  ║  Fetch   ║  ║Transform ║        │
│  ║  Users   ║→ ║  Data    ║→  ...  │
│  ║  [HTTP]  ║  ║  [JQ]    ║        │
│  ╚══════════╝  ╚══════════╝        │
└────────────────────────────────────┘
```

## Screen Structure

```
pkg/omise/screens/
├── editor.go              # Add visualization state (modify)
├── editor_visual.go       # Visualization rendering (NEW)
├── editor_visual_test.go  # Visual tests (NEW)
└── editor_nav.go          # Navigation logic (NEW)

pkg/omise/components/
└── bentobox.go            # Bento box component (NEW)
```

## Deliverables

### 1. Visual State

**File**: `pkg/omise/screens/editor.go` (modify existing)

Add to Editor struct:

```go
// Editor screen additions
type Editor struct {
	// ... existing fields

	// Visual navigation
	selectedNodeIndex int
	viewMode          ViewMode
}

// ViewMode defines the view mode
type ViewMode int

const (
	ViewModeList ViewMode = iota  // List view (current Phase 7)
	ViewModeVisual                // Visual bento box (Phase 8)
)
```

Add state to StateReview:

```go
// Update handleReviewKey to support navigation
func (e Editor) handleReviewKey(msg tea.KeyMsg) (Editor, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		// Navigate up
		if e.selectedNodeIndex > 0 {
			e.selectedNodeIndex--
		}
		return e, nil

	case "down", "j":
		// Navigate down
		nodeCount := len(e.def.Nodes)
		if e.def.Type != "group.sequence" && e.def.Type != "group.parallel" {
			nodeCount = 1 // Single node bento
		}
		if e.selectedNodeIndex < nodeCount-1 {
			e.selectedNodeIndex++
		}
		return e, nil

	case "e":
		// Edit selected node
		return e, e.editNode(e.selectedNodeIndex)

	case "m":
		// Move selected node
		return e.moveNode(e.selectedNodeIndex)

	case "d":
		// Delete selected node
		return e.deleteNode(e.selectedNodeIndex)

	case "r":
		// Run bento
		return e, e.runBento()

	case "a":
		// Add another node
		e.state = StateSelectingType
		return e, nil

	case "s", "enter":
		// Save and exit
		return e, e.saveBento()

	case "v":
		// Toggle view mode
		if e.viewMode == ViewModeList {
			e.viewMode = ViewModeVisual
		} else {
			e.viewMode = ViewModeList
		}
		return e, nil
	}

	return e, nil
}
```

### 2. Node Operations

**File**: `pkg/omise/screens/editor_nav.go` (NEW)
**Target Size**: < 200 lines

```go
package screens

import (
	"bento/pkg/neta"
	tea "github.com/charmbracelet/bubbletea"
)

// editNode launches wizard to edit node
func (e Editor) editNode(index int) tea.Cmd {
	node := e.getNode(index)
	if node == nil {
		return nil
	}

	// Launch wizard with existing values
	return func() tea.Msg {
		return EditNodeMsg{
			Index: index,
			Node:  *node,
		}
	}
}

// moveNode prompts for new position
func (e Editor) moveNode(index int) (Editor, tea.Cmd) {
	// For now, simple swap with next
	nodes := e.getNodes()
	if index >= len(nodes)-1 {
		return e, nil // Can't move last node down
	}

	// Swap with next
	nodes[index], nodes[index+1] = nodes[index+1], nodes[index]
	e.setNodes(nodes)

	e.message = "Node moved"
	return e, nil
}

// deleteNode removes node
func (e Editor) deleteNode(index int) (Editor, tea.Cmd) {
	nodes := e.getNodes()
	if index >= len(nodes) {
		return e, nil
	}

	// Remove node
	nodes = append(nodes[:index], nodes[index+1:]...)
	e.setNodes(nodes)

	// Adjust selection
	if e.selectedNodeIndex >= len(nodes) && e.selectedNodeIndex > 0 {
		e.selectedNodeIndex--
	}

	e.message = "Node deleted"
	return e, nil
}

// runBento executes the bento
func (e Editor) runBento() tea.Cmd {
	return func() tea.Msg {
		return RunBentoFromEditorMsg{
			Def: e.def,
		}
	}
}

// getNode returns node at index
func (e Editor) getNode(index int) *neta.Definition {
	nodes := e.getNodes()
	if index < 0 || index >= len(nodes) {
		return nil
	}
	return &nodes[index]
}

// getNodes returns all nodes
func (e Editor) getNodes() []neta.Definition {
	// If single-node bento
	if e.def.Type != "" && e.def.Type != "group.sequence" && e.def.Type != "group.parallel" {
		return []neta.Definition{e.def}
	}
	// Multi-node bento
	return e.def.Nodes
}

// setNodes updates nodes
func (e Editor) setNodes(nodes []neta.Definition) {
	// If converting to/from single node
	if len(nodes) == 1 && (e.def.Type == "group.sequence" || e.def.Type == "group.parallel") {
		// Convert to single node
		e.def = nodes[0]
	} else if len(nodes) > 1 && e.def.Type != "group.sequence" && e.def.Type != "group.parallel" {
		// Convert to group
		e.def = neta.Definition{
			Version: neta.CurrentVersion,
			Type:    "group.sequence",
			Name:    e.def.Name,
			Nodes:   nodes,
		}
	} else {
		e.def.Nodes = nodes
	}
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: Node operations
- ✅ Functions < 20 lines
- ✅ Clear separation from rendering
- ✅ File < 200 lines

### 3. Visual Rendering

**File**: `pkg/omise/screens/editor_visual.go` (NEW)
**Target Size**: < 250 lines

```go
package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"bento/pkg/neta"
	"bento/pkg/omise/styles"
)

// renderReview renders bento review (update existing method)
func (e Editor) renderReview() string {
	if e.viewMode == ViewModeVisual {
		return e.renderVisualBentoBox()
	}
	return e.renderListView()
}

// renderListView renders simple list (existing Phase 7 view)
func (e Editor) renderListView() string {
	content := fmt.Sprintf("Bento: %s (v%s)\n", e.def.Name, e.def.Version)
	content += fmt.Sprintf("Type: %s\n\n", e.def.Type)

	nodes := e.getNodes()
	if len(nodes) > 0 {
		content += "Nodes:\n"
		for i, node := range nodes {
			selected := ""
			if i == e.selectedNodeIndex {
				selected = "→ "
			}
			content += fmt.Sprintf("  %s%d. %s (%s)\n", selected, i+1, node.Name, node.Type)
		}
	}

	return styles.Subtle.Render(content)
}

// renderVisualBentoBox renders visual bento box
func (e Editor) renderVisualBentoBox() string {
	nodes := e.getNodes()

	// Header
	header := fmt.Sprintf("🍱 %s (v%s)", e.def.Name, e.def.Version)
	headerBox := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, true, false, true).
		BorderForeground(styles.Primary).
		Padding(0, 1).
		Width(60).
		Render(header)

	// Nodes as compartments
	compartments := []string{}
	for i, node := range nodes {
		compartment := e.renderCompartment(node, i == e.selectedNodeIndex)
		compartments = append(compartments, compartment)
	}

	// Join vertically
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		headerBox,
		"",
		strings.Join(compartments, "\n\n"),
	)

	return content
}

// renderCompartment renders a node as a compartment
func (e Editor) renderCompartment(node neta.Definition, selected bool) string {
	// Border style
	border := lipgloss.RoundedBorder()
	borderColor := styles.Muted

	if selected {
		border = lipgloss.ThickBorder()
		borderColor = styles.Primary
	}

	// Content
	title := fmt.Sprintf("%s [%s]", node.Name, node.Type)
	params := e.formatParameters(node)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		styles.Bold.Render(title),
		styles.Subtle.Render(params),
	)

	// Box
	box := lipgloss.NewStyle().
		Border(border).
		BorderForeground(borderColor).
		Padding(1, 2).
		Width(56).
		Render(content)

	return box
}

// formatParameters formats node parameters for display
func (e Editor) formatParameters(node neta.Definition) string {
	switch node.Type {
	case "http":
		return e.formatHTTPParams(node.Parameters)
	case "transform.jq":
		return e.formatJQParams(node.Parameters)
	case "conditional.if":
		return e.formatConditionalParams(node.Parameters)
	default:
		return e.formatGenericParams(node.Parameters)
	}
}

// formatHTTPParams formats HTTP parameters
func (e Editor) formatHTTPParams(params map[string]interface{}) string {
	method := "GET"
	url := ""

	if m, ok := params["method"].(string); ok {
		method = m
	}
	if u, ok := params["url"].(string); ok {
		url = u
	}

	return fmt.Sprintf("%s %s", method, url)
}

// formatJQParams formats JQ parameters
func (e Editor) formatJQParams(params map[string]interface{}) string {
	if filter, ok := params["filter"].(string); ok {
		return filter
	}
	return "No filter"
}

// formatConditionalParams formats conditional parameters
func (e Editor) formatConditionalParams(params map[string]interface{}) string {
	if cond, ok := params["condition"].(string); ok {
		return cond
	}
	return "No condition"
}

// formatGenericParams formats generic parameters
func (e Editor) formatGenericParams(params map[string]interface{}) string {
	if len(params) == 0 {
		return "No parameters"
	}

	lines := []string{}
	for key, val := range params {
		lines = append(lines, fmt.Sprintf("%s: %v", key, val))
	}

	return strings.Join(lines, "\n")
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: Visual rendering
- ✅ Functions < 20 lines
- ✅ Clear formatting per node type
- ✅ File < 250 lines

### 4. Bento Box Component

**File**: `pkg/omise/components/bentobox.go` (NEW)
**Target Size**: < 150 lines

```go
// Package components provides reusable TUI components
package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"bento/pkg/neta"
	"bento/pkg/omise/styles"
)

// BentoBox renders a visual bento box
type BentoBox struct {
	def            neta.Definition
	selectedIndex  int
	width          int
}

// NewBentoBox creates a bento box component
func NewBentoBox(def neta.Definition, selectedIndex int) BentoBox {
	return BentoBox{
		def:           def,
		selectedIndex: selectedIndex,
		width:         60,
	}
}

// SetWidth sets the box width
func (b BentoBox) SetWidth(width int) BentoBox {
	b.width = width
	return b
}

// View renders the bento box
func (b BentoBox) View() string {
	// Header
	header := b.renderHeader()

	// Compartments
	compartments := b.renderCompartments()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		compartments,
	)
}

// renderHeader renders the bento box header
func (b BentoBox) renderHeader() string {
	title := fmt.Sprintf("🍱 %s (v%s)", b.def.Name, b.def.Version)

	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, true, false, true).
		BorderForeground(styles.Primary).
		Padding(0, 1).
		Width(b.width).
		Render(title)
}

// renderCompartments renders node compartments
func (b BentoBox) renderCompartments() string {
	nodes := b.getNodes()
	compartments := []string{}

	for i, node := range nodes {
		comp := b.renderCompartment(node, i == b.selectedIndex)
		compartments = append(compartments, comp)
	}

	return lipgloss.JoinVertical(lipgloss.Left, compartments...)
}

// renderCompartment renders a single compartment
func (b BentoBox) renderCompartment(node neta.Definition, selected bool) string {
	border := lipgloss.RoundedBorder()
	borderColor := styles.Muted

	if selected {
		border = lipgloss.ThickBorder()
		borderColor = styles.Primary
	}

	title := fmt.Sprintf("%s [%s]", node.Name, node.Type)

	return lipgloss.NewStyle().
		Border(border).
		BorderForeground(borderColor).
		Padding(1, 2).
		Width(b.width - 4).
		Render(title)
}

// getNodes returns all nodes
func (b BentoBox) getNodes() []neta.Definition {
	if b.def.Type != "" && !b.def.IsGroup() {
		return []neta.Definition{b.def}
	}
	return b.def.Nodes
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: Bento box rendering
- ✅ Reusable component
- ✅ Functions < 20 lines
- ✅ File < 150 lines

### 5. Run from Editor

**File**: `pkg/omise/update.go` (add handler)

```go
// handleRunBentoFromEditor executes bento and shows progress
func (m Model) handleRunBentoFromEditor(msg screens.RunBentoFromEditorMsg) (tea.Model, tea.Cmd) {
	// Switch to executor
	m.screen = ScreenExecutor
	m.executor = m.executor.StartWorkflow(msg.Def.Name, "")

	// Stay in editor context (return after execution)
	return m, m.executor.ExecuteCmd()
}

// Update add case
case screens.RunBentoFromEditorMsg:
	return m.handleRunBentoFromEditor(msg)
```

### 6. Messages

**File**: `pkg/omise/screens/editor_messages.go` (add)

```go
// EditNodeMsg signals node should be edited
type EditNodeMsg struct {
	Index int
	Node  neta.Definition
}

// RunBentoFromEditorMsg signals run bento from editor
type RunBentoFromEditorMsg struct {
	Def neta.Definition
}
```

## Testing Strategy

**File**: `pkg/omise/screens/editor_visual_test.go`

```go
func TestEditor_Navigation(t *testing.T) {
	editor := createTestEditor()

	// Add multiple nodes
	editor.def.Nodes = []neta.Definition{
		{Type: "http", Name: "Node 1"},
		{Type: "transform.jq", Name: "Node 2"},
		{Type: "conditional.if", Name: "Node 3"},
	}

	// Test down navigation
	editor, _ = editor.Update(tea.KeyMsg{Type: tea.KeyDown})
	if editor.selectedNodeIndex != 1 {
		t.Error("should select node 1")
	}

	// Test up navigation
	editor, _ = editor.Update(tea.KeyMsg{Type: tea.KeyUp})
	if editor.selectedNodeIndex != 0 {
		t.Error("should select node 0")
	}
}

func TestEditor_MoveNode(t *testing.T) {
	editor := createTestEditor()

	editor.def.Nodes = []neta.Definition{
		{Type: "http", Name: "Node 1"},
		{Type: "transform.jq", Name: "Node 2"},
	}

	// Select first node
	editor.selectedNodeIndex = 0

	// Move down
	editor, _ = editor.moveNode(0)

	// Should swap
	if editor.def.Nodes[0].Name != "Node 2" {
		t.Error("nodes not swapped")
	}
}

func TestEditor_DeleteNode(t *testing.T) {
	editor := createTestEditor()

	editor.def.Nodes = []neta.Definition{
		{Type: "http", Name: "Node 1"},
		{Type: "transform.jq", Name: "Node 2"},
	}

	// Delete first node
	editor, _ = editor.deleteNode(0)

	if len(editor.def.Nodes) != 1 {
		t.Error("node not deleted")
	}

	if editor.def.Nodes[0].Name != "Node 2" {
		t.Error("wrong node deleted")
	}
}
```

## Validation Commands

```bash
# Test
go test -v ./pkg/omise/screens/

# Integration test
./bento

# In editor (after creating/editing bento with multiple nodes):
# 1. Press ↑/↓ to navigate between nodes
# 2. Selected node should be highlighted
# 3. Press 'e' to edit selected node
# 4. Press 'm' to move node
# 5. Press 'd' to delete node
# 6. Press 'r' to run bento
# 7. Press 'v' to toggle visual/list view
# 8. Press 's' to save
```

## Success Criteria

Phase 8 is complete when:

1. ✅ Visual bento box rendering works
2. ✅ Arrow key navigation between nodes
3. ✅ Selected node highlighted
4. ✅ Edit node in-place working
5. ✅ Move node reorders correctly
6. ✅ Delete node removes correctly
7. ✅ Run from editor executes bento
8. ✅ Toggle between list and visual views
9. ✅ All files < 250 lines
10. ✅ All functions < 20 lines
11. ✅ Tests passing
12. ✅ **Karen's approval granted**

## Common Pitfalls to Avoid

1. ❌ **Complex rendering** - Keep visual simple, text-based
2. ❌ **No bounds checking** - Validate node index always
3. ❌ **Losing selection** - Maintain selection after operations
4. ❌ **Breaking single-node bentos** - Handle both single and multi-node
5. ❌ **Poor visual feedback** - Clear indication of selected node

## Next Phase

After Karen approval, proceed to **[Phase 9: Examples & Templates](./phase-9-examples.md)** to:
- Embed example bentos in binary
- Create examples section in browser
- Add copy-from-template functionality
- Provide starting points for users

## Execution Prompt

```
I'm ready to begin Phase 8: Bento Editor - Visualization & Navigation.

I have read the Bento Box Principle and will follow it.

Please add visual bento box with:
- Text-based compartment rendering
- Arrow key navigation
- Edit/move/delete operations
- Run from editor
- View mode toggle

Each file < 250 lines, functions < 20 lines. I will use TodoWrite to track progress and get Karen's approval before completing.
```

---

**Phase 8 Bento Visualization**: Visual bento box with navigation 📦
