# Guided Form Enhancement Strategy

**Status**: Design Document
**Created**: 2025-10-18
**Related Documents**: BENTO_BOX_PRINCIPLE.md, CHARM_STACK_GUIDE.md

---

## Executive Summary

This document outlines the architecture and implementation strategy for three critical UX enhancements to the guided bento creation flow:

1. **Nested Group Node Support** - Add child nodes to sequence/parallel groups
2. **Navigation & Editing** - Navigate up/down through form stages to edit previous entries
3. **Edit Mode** - Edit existing bentos with full CRUD operations on nodes and metadata

These enhancements will transform the guided flow from a linear creation-only tool into a full-featured bento editor while maintaining the Bento Box Principle of clean, compartmentalized code.

---

## Current State Analysis

### Current Architecture

The guided creation modal (`pkg/omise/screens/guided_creation/`) consists of:

```
guided_modal.go              - Core model and state machine
guided_modal_forms.go        - Metadata and navigation forms
guided_modal_node_forms.go   - Type-specific node forms
guided_modal_transitions.go  - Stage transition logic
guided_modal_view.go         - Rendering logic
guided_modal_save.go         - Bento persistence
guided_modal_styles.go       - Lipgloss styling
```

### Current Flow (Linear State Machine)

```
guidedStageMetadata
    ↓
guidedStageNodeTypeSelect
    ↓
guidedStageNodeParameters
    ↓
guidedStageContinue
    ↓ (if "add")
guidedStageNodeTypeSelect (loop)
    ↓ (if "done")
Save & Exit
```

### Current Limitations

1. **No Nested Groups**: When you add a group node (sequence/parallel), there's no way to add children to it
2. **No Back Navigation**: Cannot edit metadata or previous nodes once you've moved forward
3. **No Edit Mode**: Cannot load and edit existing bentos
4. **No Node Deletion**: Cannot remove nodes once added
5. **Flat Structure Only**: All nodes are added to the root definition's `Nodes` array

---

## Design Principles

### 1. Maintain Bento Box Compartmentalization

Each enhancement should be isolated in its own file/package:
- `guided_modal_navigation.go` - Navigation state and history
- `guided_modal_groups.go` - Group hierarchy management
- `guided_modal_editing.go` - Edit mode operations

### 2. Preserve Charm Stack Patterns

- Use Huh forms for all data collection
- Use Bubble Tea Update/View patterns
- Use Lip Gloss for consistent styling
- Maintain single-responsibility functions (<20 lines)

### 3. Backward Compatibility

- Existing bentos must continue to work
- Current creation flow remains the default happy path
- New features are opt-in (keyboard shortcuts, menu choices)

---

## Enhancement 1: Nested Group Node Support

### Problem Statement

When a user adds a `group.sequence` or `group.parallel` node, the system should:
1. Recognize it's a container node
2. Prompt: "Add child to [group name]?" vs "Add another top-level node?"
3. Maintain a context stack to track which group is being edited

### Proposed Architecture

#### Data Structure Changes

```go
// guided_modal.go
type GuidedModal struct {
    // ... existing fields ...

    // Group hierarchy tracking
    nodeStack     []*neta.Definition  // Stack of parent nodes
    currentParent *neta.Definition     // Current parent being edited (nil = root)
}
```

#### New Stage

```go
// guided_modal.go
const (
    // ... existing stages ...
    guidedStageGroupContext guidedStage = iota + 10  // After continue stage
)
```

#### New Forms

```go
// guided_modal_groups.go
package guided_creation

// createGroupContextForm prompts user to add to group or return to parent
func (m *GuidedModal) createGroupContextForm(groupName string) *huh.Form {
    var choice string

    formWidth := m.width - 40 - 8

    return huh.NewForm(
        huh.NewGroup(
            huh.NewSelect[string]().
                Key("group_context").
                Title(fmt.Sprintf("Group '%s' created. What next?", groupName)).
                Options(
                    huh.NewOption(fmt.Sprintf("Add child to '%s'", groupName), "add_child"),
                    huh.NewOption("Add another node at current level", "add_sibling"),
                    huh.NewOption("Done with current level", "done_level"),
                    huh.NewOption("Save bento", "save"),
                ).
                Value(&choice),
        ).Title("Group Context:"),
    ).WithWidth(formWidth).WithShowHelp(false).WithShowErrors(false)
}
```

#### Transition Logic Changes

```go
// guided_modal_transitions.go (updated)
case guidedStageNodeParameters:
    if err := m.validateCurrentNode(); err != nil {
        m.validationErr = err
        return m, nil
    }

    m.validationErr = nil

    // Add node to current parent (or root if no parent)
    if m.currentParent != nil {
        m.currentParent.Nodes = append(m.currentParent.Nodes, *m.currentNode)
    } else {
        m.definition.Nodes = append(m.definition.Nodes, *m.currentNode)
    }

    // Check if the node just added is a group
    if m.currentNode.Type == "group.sequence" || m.currentNode.Type == "group.parallel" {
        // Initialize the group's Nodes array if needed
        if m.currentNode.Nodes == nil {
            m.currentNode.Nodes = []neta.Definition{}
        }

        // Move to group context prompt
        m.stage = guidedStageGroupContext
        m.form = m.createGroupContextForm(m.currentNode.Name)
        return m, m.form.Init()
    }

    // Not a group, proceed to normal continue stage
    m.stage = guidedStageContinue
    m.form = m.createContinueForm()
    return m, m.form.Init()

case guidedStageGroupContext:
    choice := m.form.GetString("group_context")

    switch choice {
    case "add_child":
        // Push current node onto stack and make it the parent
        m.nodeStack = append(m.nodeStack, m.currentParent)
        m.currentParent = m.currentNode

        // Reset current node and go to type selection
        m.currentNode = nil
        m.stage = guidedStageNodeTypeSelect
        m.form = m.createNodeTypeSelectForm()
        return m, m.form.Init()

    case "add_sibling":
        // Add another node at the same level
        m.currentNode = nil
        m.stage = guidedStageNodeTypeSelect
        m.form = m.createNodeTypeSelectForm()
        return m, m.form.Init()

    case "done_level":
        // Pop back to parent level
        if len(m.nodeStack) > 0 {
            m.currentParent = m.nodeStack[len(m.nodeStack)-1]
            m.nodeStack = m.nodeStack[:len(m.nodeStack)-1]

            // Go to continue form at parent level
            m.stage = guidedStageContinue
            m.form = m.createContinueForm()
            return m, m.form.Init()
        } else {
            // Already at root, save
            m.state = guidedStateCompleted
            return m, m.saveBento()
        }

    case "save":
        // Save and exit
        m.state = guidedStateCompleted
        return m, m.saveBento()
    }
```

#### Visual Indicator

Update the view to show current hierarchy:

```go
// guided_modal_view.go (enhancement)
func (m *GuidedModal) renderBreadcrumb() string {
    if m.currentParent == nil {
        return m.styles.Breadcrumb.Render("Root")
    }

    // Build breadcrumb from stack
    parts := []string{"Root"}
    for _, parent := range m.nodeStack {
        if parent != nil {
            parts = append(parts, parent.Name)
        }
    }
    parts = append(parts, m.currentParent.Name)

    return m.styles.Breadcrumb.Render(strings.Join(parts, " > "))
}
```

### Implementation Checklist

- [ ] Add `nodeStack` and `currentParent` fields to GuidedModal
- [ ] Create `guided_modal_groups.go` with group context form
- [ ] Add `guidedStageGroupContext` stage constant
- [ ] Update transition logic to detect group nodes
- [ ] Implement stack push/pop operations
- [ ] Add breadcrumb rendering to view
- [ ] Update preview to show hierarchy visually
- [ ] Write tests for nested group creation
- [ ] Add validation for maximum nesting depth (e.g., 5 levels)

---

## Enhancement 2: Navigation & Editing

### Problem Statement

Users need to:
1. Navigate backward through stages (Up arrow = previous stage)
2. Edit previously entered data (metadata, node parameters)
3. Navigate forward after editing (Down arrow = next stage)
4. Delete nodes from the definition

### Proposed Architecture

#### Data Structure Changes

```go
// guided_modal_navigation.go
package guided_creation

type navigationHistory struct {
    stages []guidedStage       // History of visited stages
    forms  []*huh.Form          // Snapshot of forms at each stage
    nodes  []*neta.Definition   // Snapshot of nodes at each stage
}

// GuidedModal (updated)
type GuidedModal struct {
    // ... existing fields ...

    history         navigationHistory
    historyIndex    int              // Current position in history (-1 = no history mode)
    editMode        bool             // Are we editing vs creating?
}
```

#### New Keyboard Shortcuts

```go
// guided_modal.go (Update method)
case tea.KeyMsg:
    switch msg.String() {
    case "ctrl+c":
        return m, tea.Interrupt
    case "esc":
        // Cancel or exit edit mode
        if m.editMode && m.historyIndex >= 0 {
            // Exit history navigation mode
            m.historyIndex = -1
            m.editMode = false
            return m, nil
        }
        // Cancel the entire flow
        m.state = guidedStateCancelled
        return m, func() tea.Msg {
            return GuidedCompleteMsg{Cancelled: true}
        }

    case "ctrl+up", "ctrl+k":
        // Navigate to previous stage in history
        return m, m.navigateHistory(-1)

    case "ctrl+down", "ctrl+j":
        // Navigate to next stage in history
        return m, m.navigateHistory(1)

    case "ctrl+d":
        // Delete current node (if on node edit stage)
        if m.stage == guidedStageNodeParameters && m.editMode {
            return m, m.deleteCurrentNode()
        }
    }
```

#### Navigation Functions

```go
// guided_modal_navigation.go
package guided_creation

import tea "github.com/charmbracelet/bubbletea"

// navigateHistory moves backward or forward in history
func (m *GuidedModal) navigateHistory(direction int) tea.Cmd {
    newIndex := m.historyIndex + direction

    // Validate bounds
    if newIndex < 0 || newIndex >= len(m.history.stages) {
        return nil
    }

    // Update index
    m.historyIndex = newIndex
    m.editMode = true

    // Restore stage and form from history
    m.stage = m.history.stages[newIndex]
    m.form = m.history.forms[newIndex]

    // Restore node state if applicable
    if newIndex < len(m.history.nodes) && m.history.nodes[newIndex] != nil {
        m.currentNode = m.history.nodes[newIndex]
    }

    return m.form.Init()
}

// captureHistorySnapshot saves current state to history
func (m *GuidedModal) captureHistorySnapshot() {
    // Only capture if not in history navigation mode
    if m.historyIndex >= 0 {
        return
    }

    // Clone form (create new instance with same values)
    formClone := m.cloneCurrentForm()

    // Clone current node if it exists
    var nodeClone *neta.Definition
    if m.currentNode != nil {
        clone := *m.currentNode
        nodeClone = &clone
    }

    // Append to history
    m.history.stages = append(m.history.stages, m.stage)
    m.history.forms = append(m.history.forms, formClone)
    m.history.nodes = append(m.history.nodes, nodeClone)
}

// cloneCurrentForm creates a deep copy of the current form
func (m *GuidedModal) cloneCurrentForm() *huh.Form {
    // This depends on which stage we're on
    switch m.stage {
    case guidedStageMetadata:
        return m.createMetadataForm()
    case guidedStageNodeTypeSelect:
        return m.createNodeTypeSelectForm()
    case guidedStageNodeParameters:
        if m.currentNode != nil {
            return m.createNodeFormForType(m.currentNode.Type)
        }
        return m.createNodeTypeSelectForm()
    case guidedStageContinue:
        return m.createContinueForm()
    case guidedStageGroupContext:
        if m.currentNode != nil {
            return m.createGroupContextForm(m.currentNode.Name)
        }
        return m.createContinueForm()
    default:
        return m.form
    }
}

// deleteCurrentNode removes the node being edited from the definition
func (m *GuidedModal) deleteCurrentNode() tea.Cmd {
    if m.currentNode == nil {
        return nil
    }

    // Find and remove node from parent or root
    if m.currentParent != nil {
        m.currentParent.Nodes = removeNodeByName(m.currentParent.Nodes, m.currentNode.Name)
    } else {
        m.definition.Nodes = removeNodeByName(m.definition.Nodes, m.currentNode.Name)
    }

    // Go back to continue stage
    m.currentNode = nil
    m.stage = guidedStageContinue
    m.form = m.createContinueForm()

    return m.form.Init()
}

// removeNodeByName removes a node from a slice by name
func removeNodeByName(nodes []neta.Definition, name string) []neta.Definition {
    result := make([]neta.Definition, 0, len(nodes))
    for _, node := range nodes {
        if node.Name != name {
            result = append(result, node)
        }
    }
    return result
}
```

#### View Updates

```go
// guided_modal_view.go (enhancement)
func (m *GuidedModal) renderHelpBar() string {
    help := []string{
        "esc: cancel",
        "ctrl+↑/↓: navigate history",
    }

    if m.editMode && m.stage == guidedStageNodeParameters {
        help = append(help, "ctrl+d: delete node")
    }

    return m.styles.Help.Render(strings.Join(help, " • "))
}
```

#### Transition Logic Updates

```go
// guided_modal_transitions.go (enhancement)
func (m *GuidedModal) handleStageTransition() (*GuidedModal, tea.Cmd) {
    // Capture history before transitioning
    if !m.editMode {
        m.captureHistorySnapshot()
    }

    // ... existing transition logic ...
}
```

### Implementation Checklist

- [ ] Create `guided_modal_navigation.go` with history management
- [ ] Add `history` and `historyIndex` fields to GuidedModal
- [ ] Implement `navigateHistory()` function
- [ ] Implement `captureHistorySnapshot()` function
- [ ] Implement `cloneCurrentForm()` function
- [ ] Implement `deleteCurrentNode()` function
- [ ] Add keyboard shortcuts (Ctrl+Up/Down, Ctrl+D)
- [ ] Update view to show help bar with shortcuts
- [ ] Update transition logic to capture history
- [ ] Write tests for navigation and deletion
- [ ] Add visual indicator when in edit mode

---

## Enhancement 3: Edit Mode for Existing Bentos

### Problem Statement

Users need to:
1. Open an existing bento for editing
2. Edit metadata (name, description, icon)
3. Edit individual nodes (parameters, type)
4. Add/delete nodes
5. Restructure the graph (move nodes, change parent)

### Proposed Architecture

#### Entry Point

```go
// guided_modal.go
// NewGuidedModalForEdit creates a modal for editing an existing bento
func NewGuidedModalForEdit(
    store *jubako.Store,
    workDir string,
    width, height int,
    existingDef *neta.Definition,
) *GuidedModal {
    m := NewGuidedModal(store, workDir, width, height)

    // Mark as editing mode
    m.editing = true

    // Load existing definition
    m.definition = existingDef

    // Start at edit menu instead of metadata
    m.stage = guidedStageEditMenu
    m.form = m.createEditMenuForm()

    return m
}
```

#### New Stages

```go
// guided_modal.go
const (
    // ... existing stages ...
    guidedStageEditMenu     guidedStage = iota + 20
    guidedStageEditMetadata guidedStage = iota + 21
    guidedStageNodeList     guidedStage = iota + 22
    guidedStageNodeEdit     guidedStage = iota + 23
)
```

#### New Forms

```go
// guided_modal_editing.go
package guided_creation

import "github.com/charmbracelet/huh"

// createEditMenuForm presents top-level editing options
func (m *GuidedModal) createEditMenuForm() *huh.Form {
    var choice string

    formWidth := m.width - 40 - 8

    return huh.NewForm(
        huh.NewGroup(
            huh.NewSelect[string]().
                Key("edit_choice").
                Title("What would you like to edit?").
                Options(
                    huh.NewOption("Edit metadata (name, description, icon)", "metadata"),
                    huh.NewOption("Add a new node", "add_node"),
                    huh.NewOption("Edit existing nodes", "edit_nodes"),
                    huh.NewOption("Delete nodes", "delete_nodes"),
                    huh.NewOption("Save and exit", "save"),
                    huh.NewOption("Cancel without saving", "cancel"),
                ).
                Value(&choice),
        ).Title(fmt.Sprintf("Editing: %s", m.definition.Name)),
    ).WithWidth(formWidth).WithShowHelp(false).WithShowErrors(false)
}

// createNodeListForm shows all nodes for selection
func (m *GuidedModal) createNodeListForm() *huh.Form {
    var selectedNode string

    formWidth := m.width - 40 - 8

    // Build options from existing nodes
    options := make([]huh.Option[string], 0, len(m.definition.Nodes))
    for _, node := range m.definition.Nodes {
        label := fmt.Sprintf("%s (%s)", node.Name, node.Type)
        options = append(options, huh.NewOption(label, node.Name))
    }

    if len(options) == 0 {
        options = append(options, huh.NewOption("No nodes to edit", ""))
    }

    return huh.NewForm(
        huh.NewGroup(
            huh.NewSelect[string]().
                Key("selected_node").
                Title("Select a node to edit").
                Options(options...).
                Value(&selectedNode),
        ).Title("Nodes:"),
    ).WithWidth(formWidth).WithShowHelp(false).WithShowErrors(false)
}
```

#### Edit Mode Transitions

```go
// guided_modal_transitions.go (additions)
case guidedStageEditMenu:
    choice := m.form.GetString("edit_choice")

    switch choice {
    case "metadata":
        // Populate form with existing metadata
        m.stage = guidedStageEditMetadata
        m.form = m.createMetadataFormWithValues(
            m.definition.Name,
            m.definition.Description,
            m.definition.Icon,
        )
        return m, m.form.Init()

    case "add_node":
        // Go to standard node creation flow
        m.stage = guidedStageNodeTypeSelect
        m.form = m.createNodeTypeSelectForm()
        return m, m.form.Init()

    case "edit_nodes":
        // Show list of nodes to select
        m.stage = guidedStageNodeList
        m.form = m.createNodeListForm()
        return m, m.form.Init()

    case "delete_nodes":
        // Show list of nodes with delete option
        m.stage = guidedStageNodeList
        m.form = m.createNodeListFormWithDelete()
        return m, m.form.Init()

    case "save":
        // Save and exit
        m.state = guidedStateCompleted
        return m, m.saveBento()

    case "cancel":
        // Cancel without saving
        m.state = guidedStateCancelled
        return m, func() tea.Msg {
            return GuidedCompleteMsg{Cancelled: true}
        }
    }

case guidedStageEditMetadata:
    // Update definition and return to menu
    m.updateDefinitionFromForm()
    m.stage = guidedStageEditMenu
    m.form = m.createEditMenuForm()
    return m, m.form.Init()

case guidedStageNodeList:
    // Load selected node for editing
    nodeName := m.form.GetString("selected_node")
    if nodeName == "" {
        // Return to menu
        m.stage = guidedStageEditMenu
        m.form = m.createEditMenuForm()
        return m, m.form.Init()
    }

    // Find node in definition
    node := m.findNodeByName(nodeName)
    if node == nil {
        // Node not found, return to menu
        m.stage = guidedStageEditMenu
        m.form = m.createEditMenuForm()
        return m, m.form.Init()
    }

    // Set as current node and load appropriate form
    m.currentNode = node
    m.stage = guidedStageNodeEdit
    m.form = m.createNodeFormForTypeWithValues(node.Type, node)
    return m, m.form.Init()

case guidedStageNodeEdit:
    // Update node and return to menu
    m.updateCurrentNodeFromNodeForm(m.currentNode.Type)
    m.stage = guidedStageEditMenu
    m.form = m.createEditMenuForm()
    return m, m.form.Init()
```

#### Helper Functions

```go
// guided_modal_editing.go
// findNodeByName recursively searches for a node by name
func (m *GuidedModal) findNodeByName(name string) *neta.Definition {
    return findNodeInTree(m.definition.Nodes, name)
}

func findNodeInTree(nodes []neta.Definition, name string) *neta.Definition {
    for i := range nodes {
        if nodes[i].Name == name {
            return &nodes[i]
        }

        // Search children if this is a group
        if len(nodes[i].Nodes) > 0 {
            if found := findNodeInTree(nodes[i].Nodes, name); found != nil {
                return found
            }
        }
    }
    return nil
}

// createMetadataFormWithValues creates a metadata form pre-populated with values
func (m *GuidedModal) createMetadataFormWithValues(name, description, icon string) *huh.Form {
    // Similar to createMetadataForm but with initial values set
    // ... implementation ...
}

// createNodeFormForTypeWithValues creates a node form pre-populated with values
func (m *GuidedModal) createNodeFormForTypeWithValues(nodeType string, node *neta.Definition) *huh.Form {
    // Similar to createNodeFormForType but with values from node
    // ... implementation ...
}
```

### Implementation Checklist

- [ ] Create `guided_modal_editing.go` with edit-specific forms
- [ ] Add edit mode stages (EditMenu, EditMetadata, NodeList, NodeEdit)
- [ ] Implement `NewGuidedModalForEdit()` constructor
- [ ] Create `createEditMenuForm()` function
- [ ] Create `createNodeListForm()` function
- [ ] Implement `findNodeByName()` helper
- [ ] Create form builders with pre-populated values
- [ ] Add edit mode transitions
- [ ] Update main screen to support "Edit Bento" action
- [ ] Write tests for edit operations
- [ ] Add confirmation dialog for destructive actions (delete, cancel)

---

## Integration & Rollout Plan

### Phase 1: Nested Group Support (Week 1-2)

**Goal**: Users can add children to group nodes

**Tasks**:
1. Implement group detection and context form
2. Add node stack management
3. Update transitions for group hierarchy
4. Add breadcrumb visualization
5. Test with deeply nested structures

**Success Criteria**:
- Can create sequence with 3 child nodes
- Can create parallel group inside sequence group
- Breadcrumb shows current context
- Preview renders hierarchy correctly

### Phase 2: Navigation & Editing (Week 3-4)

**Goal**: Users can navigate backward/forward and edit previous entries

**Tasks**:
1. Implement history capture system
2. Add keyboard shortcuts for navigation
3. Implement form cloning for snapshots
4. Add delete node functionality
5. Update UI to show edit mode state

**Success Criteria**:
- Ctrl+Up/Down navigates through history
- Can edit metadata after adding nodes
- Can delete a node with Ctrl+D
- Help bar shows available shortcuts
- History doesn't break form validation

### Phase 3: Edit Mode (Week 5-6)

**Goal**: Users can open and edit existing bentos

**Tasks**:
1. Implement edit menu and forms
2. Add node selection list
3. Create form pre-population logic
4. Integrate with main screen
5. Add save/cancel confirmation

**Success Criteria**:
- Can load existing bento for editing
- Can change metadata without affecting nodes
- Can edit individual node parameters
- Can add nodes to existing bento
- Changes save correctly

### Phase 4: Polish & Testing (Week 7)

**Goal**: Production-ready features with comprehensive tests

**Tasks**:
1. Write unit tests for all new functions
2. Write integration tests for flows
3. Add validation for edge cases
4. Improve error messages
5. Performance optimization
6. Documentation updates

**Success Criteria**:
- >80% test coverage for new code
- No regressions in existing functionality
- Clean separation of concerns (Bento Box compliance)
- All files <250 lines
- Documentation updated

---

## File Structure (Post-Implementation)

```
pkg/omise/screens/guided_creation/
├── guided_modal.go                    # Core model (updated with new fields)
├── guided_modal_forms.go              # Base forms (metadata, continue)
├── guided_modal_node_forms.go         # Node type forms
├── guided_modal_groups.go             # NEW: Group hierarchy management
├── guided_modal_navigation.go         # NEW: History and navigation
├── guided_modal_editing.go            # NEW: Edit mode operations
├── guided_modal_transitions.go        # Transition logic (updated)
├── guided_modal_view.go               # Rendering (updated with breadcrumb, help)
├── guided_modal_save.go               # Persistence
├── guided_modal_styles.go             # Styling
├── guided_modal_messages.go           # Message types
├── guided_modal_test.go               # Tests (updated)
├── guided_modal_test_helpers.go       # Test utilities
└── guided_modal_preview.go            # Preview rendering
```

**Bento Box Compliance**:
- Each file has single responsibility
- New features in separate files
- No file exceeds 250 lines
- Clean interfaces between components

---

## Technical Considerations

### 1. Form State Management

**Challenge**: Huh forms are stateful, so cloning for history is non-trivial

**Solution**: Instead of cloning form state, clone the *data* and recreate forms on navigation

```go
// Don't clone form
m.history.forms = append(m.history.forms, m.form.Clone()) // ❌

// Clone data and recreate form
snapshot := m.extractFormData()
m.history.data = append(m.history.data, snapshot)
// Later: m.form = m.createFormFromSnapshot(snapshot) // ✅
```

### 2. Validation During Navigation

**Challenge**: When navigating backward, should we re-validate?

**Solution**: Skip validation in history navigation mode

```go
if m.editMode && m.historyIndex >= 0 {
    // Don't validate, just restore state
    return m, nil
}
```

### 3. Group Node Modification

**Challenge**: If user edits a group node that already has children, what happens?

**Solution**: Warn user and offer options

```go
if len(node.Nodes) > 0 {
    // Show warning: "This group has N children. Editing may affect them."
    // Options: "Continue", "Cancel"
}
```

### 4. Undo/Redo

**Future Enhancement**: History system sets up foundation for undo/redo

```go
// Could add later:
case "ctrl+z":
    return m, m.undo()
case "ctrl+y":
    return m, m.redo()
```

### 5. Large Bentos

**Challenge**: Editing bentos with 50+ nodes becomes unwieldy

**Solution**: Add search/filter to node list

```go
// Future enhancement
func (m *GuidedModal) createNodeListFormWithSearch() *huh.Form {
    return huh.NewForm(
        huh.NewGroup(
            huh.NewInput().
                Key("search").
                Title("Search nodes").
                Placeholder("Type to filter..."),
            // Filtered list based on search
        ),
    )
}
```

---

## Risks & Mitigations

### Risk 1: Complexity Explosion

**Risk**: Adding navigation + groups + editing could make code unmaintainable

**Mitigation**: Strict adherence to Bento Box Principle
- One responsibility per file
- Functions <20 lines
- Clear separation of concerns
- Regular refactoring reviews

### Risk 2: State Synchronization Bugs

**Risk**: Multiple state variables (stage, history, nodeStack) could get out of sync

**Mitigation**:
- Single source of truth for each piece of state
- Invariant checking in tests
- State transition logging for debugging

### Risk 3: UX Confusion

**Risk**: Too many options/shortcuts could overwhelm users

**Mitigation**:
- Progressive disclosure (show advanced options only when relevant)
- Clear visual indicators of mode (create vs edit vs history)
- Help text and keyboard shortcut guide
- User testing before release

### Risk 4: Performance Degradation

**Risk**: History snapshots and deep cloning could slow down large bentos

**Mitigation**:
- Lazy cloning (clone only when needed)
- History size limit (e.g., max 50 snapshots)
- Benchmark tests for large bentos

---

## Success Metrics

### User Experience

- **Task Completion Time**: Creating a 5-node bento should take <2 minutes
- **Error Rate**: <5% of save operations should fail validation
- **User Satisfaction**: Positive feedback from beta testers

### Code Quality

- **Test Coverage**: >80% for new code
- **File Size**: All files <250 lines (Bento Box compliance)
- **Function Size**: 90% of functions <20 lines
- **Cyclomatic Complexity**: <10 for all functions

### Technical

- **Performance**: Form transitions <50ms
- **Memory**: History uses <10MB for typical bentos
- **Stability**: Zero panics in production

---

## Future Enhancements (Out of Scope)

These are intentionally deferred to avoid scope creep:

1. **Drag & Drop Node Reordering** - Requires terminal mouse support
2. **Visual Graph Editor** - Requires full TUI canvas library
3. **Multi-Bento Editing** - Edit multiple bentos simultaneously
4. **Git Integration** - Track changes, commit, diff
5. **Templates & Snippets** - Reusable node patterns
6. **Collaborative Editing** - Multi-user real-time editing
7. **AI Assistance** - Suggest nodes based on context

---

## Alternatives Considered

### Alternative 1: Separate Edit Screen

**Approach**: Create entirely separate screen for editing vs creating

**Pros**:
- Cleaner separation of concerns
- Easier to reason about each mode

**Cons**:
- Code duplication (forms, validation, etc.)
- Inconsistent UX between create and edit
- Harder to transition between modes

**Decision**: Rejected - Unified modal provides better UX and code reuse

### Alternative 2: Form-Based Navigation Only (No Keyboard)

**Approach**: Only allow navigation through form selections (no Ctrl+Up/Down)

**Pros**:
- Simpler implementation
- More discoverable for non-technical users

**Cons**:
- Slower for power users
- More clicks to accomplish tasks
- Less intuitive for editing workflows

**Decision**: Rejected - Keyboard shortcuts are essential for efficiency

### Alternative 3: Flat Structure with References

**Approach**: Keep all nodes flat in root, use references for hierarchy

**Pros**:
- Simpler data structure
- Easier to serialize

**Cons**:
- Doesn't match neta.Definition structure
- Harder to visualize hierarchy
- Reference management adds complexity

**Decision**: Rejected - Hierarchical structure matches domain model

---

## Questions for Review

Before implementation, please review and answer:

1. **Scope**: Is this the right level of detail, or should we drill deeper into specific areas?

2. **Priorities**: Should we implement phases in a different order? (e.g., edit mode before navigation?)

3. **Architecture**: Any concerns about the proposed data structures or file organization?

4. **UX**: Are the keyboard shortcuts intuitive? Should we use different keys?

5. **Testing**: What scenarios are most critical to test?

6. **Compatibility**: Any concerns about backward compatibility with existing bentos?

7. **Performance**: Should we add explicit performance benchmarks for large bentos (50+ nodes)?

8. **Documentation**: What additional documentation is needed beyond this strategy doc?

---

## Appendix: State Machine Diagram

### Current State Machine (Linear)

```
┌─────────────────┐
│ Metadata        │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Node Type       │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Node Parameters │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Continue?       │
└────────┬────────┘
         │
    ┌────┴────┐
    │         │
    ▼         ▼
  Add       Done
    │         │
    │         ▼
    │    ┌────────┐
    │    │  Save  │
    │    └────────┘
    │
    └──────┐
           │
           ▼
    ┌─────────────────┐
    │ Node Type       │ (loop)
    └─────────────────┘
```

### Enhanced State Machine (With Nested Groups, Navigation, Edit)

```
                    ┌─────────────────────┐
                    │   Create or Edit?   │
                    └──────────┬──────────┘
                               │
                    ┌──────────┴──────────┐
                    │                     │
                    ▼                     ▼
            ┌─────────────┐      ┌──────────────┐
            │  Metadata   │      │  Edit Menu   │
            └──────┬──────┘      └──────┬───────┘
                   │                    │
                   │         ┌──────────┼──────────┬──────────┐
                   │         │          │          │          │
                   │         ▼          ▼          ▼          ▼
                   │    ┌─────────┐ ┌────┐ ┌─────────┐ ┌────────┐
                   │    │Edit Meta│ │Add │ │Edit Node│ │Delete  │
                   │    └─────────┘ └────┘ └─────────┘ └────────┘
                   │         │         │         │          │
                   │         └─────────┴─────────┴──────────┘
                   │                   │
                   ▼                   ▼
            ┌─────────────┐     ┌─────────────┐
            │ Node Type   │ ◄───┤ Node Type   │
            └──────┬──────┘     └──────┬──────┘
                   │                   │
                   ▼                   ▼
            ┌─────────────┐     ┌─────────────┐
            │ Node Params │     │ Node Params │
            └──────┬──────┘     └──────┬──────┘
                   │                   │
                   ▼                   ▼
            ┌─────────────┐     ┌─────────────┐
            │ Is Group?   │     │   Update    │
            └──────┬──────┘     └──────┬──────┘
                   │                   │
            ┌──────┴──────┐            │
            │             │            │
            ▼             ▼            │
         ┌────┐    ┌────────────┐     │
         │ No │    │ Group Menu │     │
         └─┬──┘    └─────┬──────┘     │
           │             │            │
           │    ┌────────┼───────┐   │
           │    │        │       │   │
           │    ▼        ▼       ▼   │
           │  Add     Add     Done   │
           │  Child   Sib     Level  │
           │    │       │       │    │
           │    └───────┴───────┘    │
           │            │             │
           ▼            ▼             │
      ┌─────────┐  ┌─────────┐       │
      │Continue │  │Push/Pop │       │
      └────┬────┘  │  Stack  │       │
           │       └────┬────┘       │
           │            │            │
           ▼            ▼            │
      ┌─────────┐  ┌─────────┐      │
      │  Add?   │  │Continue │      │
      └────┬────┘  └────┬────┘      │
           │            │            │
      ┌────┴────┐       │            │
      │         │       │            │
      ▼         ▼       ▼            ▼
    Add       Done    Save ◄─────── Save
      │         │
      └─────────┘
           │
           ▼
       ┌──────┐
       │ Save │
       └──────┘

    [Ctrl+↑/↓ = Navigate History - Available from any stage]
    [Ctrl+D = Delete Node - Available on Node Params stage]
    [ESC = Cancel/Exit - Available from any stage]
```

---

**Document End**

This strategy document provides a comprehensive roadmap for implementing nested group support, navigation/editing, and edit mode for the guided bento creation flow. Implementation should follow the phased approach outlined above, with continuous testing and adherence to the Bento Box Principle.
