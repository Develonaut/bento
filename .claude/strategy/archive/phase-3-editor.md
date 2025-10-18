# Phase 3: Editor Simplification - Two-Section Design

**Duration**: 3-4 days | **Complexity**: High | **Risk**: Medium | **Status**: Not Started

## Overview
Complete redesign of the editor with a simplified two-section approach: table view for browsing nodes and Huh forms for editing.

## Goals
- Simplify editor interface
- Implement table view for nodes using Bubbles table
- Create Huh form integration for node editing
- Support nested display for loops/groups
- Feature-flag new editor for safe rollout

## Related Original TODOs
- TODO #3: Complete editor redesign

## Visual Design

### Two-Section Layout
```
┌─────────────────────────────────┐
│ Node Table (Bubbles Table)      │ ← Top 60%
│ ┌─────┬──────┬─────────────┐   │
│ │Name │Type  │Parameters    │   │
│ ├─────┼──────┼─────────────┤   │
│ │http1│HTTP  │GET /api/...  │   │
│ │loop1│Loop  │[nested table]│   │
│ │grp1 │Group │[nested table]│   │
│ └─────┴──────┴─────────────┘   │
├─────────────────────────────────┤
│ Node Editor (Huh Forms)         │ ← Bottom 40%
│ [Dynamic form based on node]    │
│                                 │
└─────────────────────────────────┘
```

## Implementation Details

### Package Structure
Create new package: `pkg/omise/screens/editor_v2/`

### Files to Create
```
pkg/omise/screens/editor_v2/
├── editor.go      (<250 lines) - Main editor logic
├── table.go       (<250 lines) - Table view & navigation
├── forms.go       (<250 lines) - Huh form generation
└── messages.go    (<250 lines) - Custom messages
```

### Editor Structure
```go
// pkg/omise/screens/editor_v2/editor.go
type EditorV2 struct {
    nodeTable table.Model
    nodeForm  *huh.Form
    nodes     []neta.Definition
    mode      EditMode
    selected  int
}

type EditMode int

const (
    ModeBrowsing EditMode = iota
    ModeEditing
    ModeCreating
)
```

### Table Features
- Display all nodes in current bento
- Columns: Name, Type, Parameters/Summary
- Show node hierarchy (nested for loops/groups)
- Keyboard navigation (arrow keys, vim keys)
- Selection highlighting
- Summary view of node parameters

### Form Features
- Dynamic forms based on selected node type
- Support all existing neta.Definition types:
  - HTTP requests
  - Shell commands
  - Loops
  - Groups
  - Assertions
  - Variables
- Validation before submission
- Clear success/error feedback
- Cancel option (Esc key)

## Feature Flag

### Implementation Options
1. **Environment Variable**: `BENTO_EDITOR_V2=true`
2. **Config File**: Add to settings
3. **Command Flag**: `bento taste --editor-v2`

### Default Behavior
- Default to old editor initially
- Easy toggle in settings
- Keep old editor code intact during transition
- Clear migration path for users

## Dependencies
- Existing Bubbles table component
- Existing Huh integration
- Phase 1 & 2 completion recommended but not required

## Bento Box Compliance
- [ ] Each file < 250 lines
- [ ] Functions < 20 lines
- [ ] Clear separation of concerns
- [ ] No circular dependencies
- [ ] Use standard library where possible

## Testing Requirements
- [ ] Basic node operations work (add, edit, delete)
- [ ] Table navigation smooth
- [ ] Form validation works
- [ ] Complex bentos handled (nested loops/groups)
- [ ] Feature flag toggles correctly
- [ ] Old editor still works
- [ ] No data loss during editing

## Risk Mitigation
- Don't delete old editor code
- Comprehensive error handling
- Graceful degradation if issues occur
- Easy rollback via feature flag
- Extensive testing before default switch

## Success Criteria
- [ ] EditorV2 package created with all files
- [ ] Table view displays nodes correctly
- [ ] Nested tables work for loops/groups
- [ ] Forms generate dynamically per node type
- [ ] All node types supported
- [ ] Validation prevents invalid data
- [ ] Feature flag implemented
- [ ] Old editor preserved and functional
- [ ] All tests pass
- [ ] No regressions

## References
- https://github.com/charmbracelet/bubbletea/blob/main/examples/table/main.go
- https://github.com/charmbracelet/huh

## Notes
- This is the most complex phase
- Take time to get it right
- Test thoroughly with real bentos
- Consider user feedback before making default
- Keep old editor until confident in new one

---

## Claude Code Prompt for Guilliman

```
Implement Phase 3 of the Bento TUI enhancement plan: Editor Simplification - Two-Section Design.

IMPORTANT: Read the phase document at .claude/strategy/phase-3-editor.md for complete context, requirements, and success criteria before starting implementation.


IMPORTANT: This is a major redesign. Create the new editor alongside the existing one with a feature flag for safe rollout.

Requirements:
1. Create new package pkg/omise/screens/editor_v2/ with these files:

   a. editor.go (<250 lines):
      - EditorV2 struct with table.Model and huh.Form
      - Main Update() and View() methods
      - Mode management (browsing table vs editing node)
      - Integration with existing bento file operations

   b. table.go (<250 lines):
      - Table view configuration for nodes
      - Column definitions: Name, Type, Parameters/Summary
      - Row rendering logic
      - Support for nested tables (loops/groups)
      - Table navigation handlers

   c. forms.go (<250 lines):
      - Huh form generation based on node type
      - Dynamic form fields for different neta.Definition types
      - Form validation logic
      - Form submission handlers

   d. messages.go (<250 lines):
      - Custom messages for editor operations
      - Node add/edit/delete messages
      - Table selection messages
      - Form submission messages

2. Two-section layout:
   - Top 60%: Bubbles table showing nodes
   - Bottom 40%: Huh form for editing selected node
   - Clear visual separator between sections
   - Responsive to terminal size

3. Table features:
   - Display all nodes in current bento
   - Show node hierarchy (nested for loops/groups)
   - Keyboard navigation (arrow keys)
   - Selection highlighting
   - Summary view of node parameters

4. Form features:
   - Dynamic forms based on selected node type
   - Support all existing neta.Definition types
   - Validation before submission
   - Clear success/error feedback
   - Cancel option

5. Feature flag:
   - Add config option to enable EditorV2
   - Default to old editor
   - Easy toggle in settings or env var
   - Keep old editor code intact

6. Follow Bento Box principles:
   - Each file < 250 lines
   - Functions < 20 lines
   - Clear separation of concerns
   - No circular dependencies
   - Use standard library where possible

7. References:
   - https://github.com/charmbracelet/bubbletea/blob/main/examples/table/main.go
   - https://github.com/charmbracelet/huh

8. Testing:
   - Ensure basic node operations work
   - Test with complex bentos (nested loops/groups)
   - Verify form validation
   - Test table navigation

9. Risk mitigation:
   - Don't delete old editor
   - Comprehensive error handling
   - Graceful degradation if issues occur

Return a detailed summary of:
- Files created and their responsibilities
- Design decisions made
- How feature flag works
- Testing results
- Any issues or limitations encountered
```
