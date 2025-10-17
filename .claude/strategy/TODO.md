# Bento TUI Enhancement - Phased Implementation Plan

This document outlines the phased approach to improving the Bento TUI codebase. Each phase includes a prompt that can be used with Claude Code to have Guilliman perform the work.

## Tab Names (Kitchen Workflow Theme)
- **🍱 Bentos** → Your prepared bento boxes (Home/Browser)
- **📖 Recipes** → Browse example bentos (Discover/Examples)
- **🔪 Mise** → Your workspace setup (Settings)
- **👨‍🍳 Sensei** → Guidance from the master (Help)

---

## Phase 1: Foundation - Help & Key System Enhancement
**Duration**: 1-2 days | **Complexity**: Low | **Risk**: Low

### What Will Be Accomplished
- Implement proper Bubbles help and key components
- Standardize keybindings across all screens
- Create consistent footer help menu
- Establish foundation for better user guidance

### Related Original TODOs
- TODO #2: Bubbles Help and Key components
- Partial TODO #4: Header styling preparation

### Implementation Details
- Create centralized `pkg/omise/components/keys.go` for key definitions
- Enhance `pkg/omise/components/footer.go` with help.Model integration
- Standardize key bindings across all screens
- Reference: https://github.com/charmbracelet/bubbletea/blob/main/examples/help/main.go
- Reference: https://github.com/charmbracelet/bubbles?tab=readme-ov-file#key

### Claude Code Prompt for Guilliman
```
Implement Phase 1 of the Bento TUI enhancement plan: Foundation - Help & Key System Enhancement.

Requirements:
1. Create pkg/omise/components/keys.go with centralized KeyMap structure containing:
   - NavigationKeys (common navigation across all screens)
   - BrowserKeys (browser-specific keys)
   - EditorKeys (editor-specific keys)
   - SettingsKeys (settings-specific keys)
   - GlobalKeys (app-wide keys like quit)

2. Enhance pkg/omise/components/footer.go to use Bubbles help.Model:
   - Integrate help component for dynamic key display
   - Support full/short help modes
   - Ensure width-responsive rendering

3. Update all existing screens to use the new centralized key system:
   - Browser, Editor, Settings, Help screens
   - Ensure consistent key bindings
   - Update their View() methods to use the new footer

4. Follow Bento Box principles:
   - Files must be < 250 lines
   - Functions must be < 20 lines
   - No circular dependencies
   - Use Go standard library where possible

5. References:
   - https://github.com/charmbracelet/bubbletea/blob/main/examples/help/main.go
   - https://github.com/charmbracelet/bubbles?tab=readme-ov-file#key

6. Ensure all changes are tested and the app still runs correctly.

Return a summary of changes made, files modified/created, and any issues encountered.
```

---

## Phase 2: Navigation - Tab System Implementation
**Duration**: 2-3 days | **Complexity**: Medium | **Risk**: Low

### What Will Be Accomplished
- Implement tab-based navigation with Kitchen Workflow theme
- Create visual tab indicators in header
- Maintain existing screen functionality
- Add keyboard shortcuts for tab switching

### Related Original TODOs
- TODO #1: Bubbletea Tabs implementation
- Complete TODO #4: Full header styling with tabs

### Tab Structure
```
🍱 Bento v0.0.1
┌─────────┬─────────┬────────┬─────────┐
│🍱 Bentos│📖 Recipes│🔪 Mise │👨‍🍳 Sensei│
└─────────┴─────────┴────────┴─────────┘
```

### Implementation Details
- Create `pkg/omise/components/tabs.go` for tab navigation logic
- Update header component to display tabs with active state
- Map existing screens to tabs: Browser→Bentos, Examples→Recipes, Settings→Mise, Help→Sensei
- Implement tab switching with number keys (1-4) or Tab key
- Reference: https://github.com/charmbracelet/bubbletea/blob/main/examples/tabs/main.go

### Claude Code Prompt for Guilliman
```
Implement Phase 2 of the Bento TUI enhancement plan: Navigation - Tab System Implementation.

Requirements:
1. Create pkg/omise/components/tabs.go (<250 lines) with:
   - TabView struct managing active tab state
   - Tab definitions for: Bentos (🍱), Recipes (📖), Mise (🔪), Sensei (👨‍🍳)
   - Tab rendering with active/inactive states
   - Tab switching logic

2. Update header component to display tabs:
   - Two-line header with tabs
   - Visual indicators for active tab
   - Consistent styling with existing theme
   - Emoji icons for each tab

3. Map existing screens to tabs:
   - Browser screen → Bentos tab
   - Examples/Discovery screen → Recipes tab (if exists, otherwise plan for future)
   - Settings screen → Mise tab
   - Help screen → Sensei tab

4. Implement tab navigation:
   - Number keys 1-4 for direct tab access
   - Tab key for cycling through tabs
   - Update centralized KeyMap from Phase 1

5. Ensure the main model in pkg/omise/omise.go properly handles:
   - Tab switching messages
   - Screen transitions
   - State preservation when switching tabs

6. Follow Bento Box principles:
   - Files < 250 lines
   - Functions < 20 lines
   - Clear separation of concerns
   - No circular dependencies

7. References:
   - https://github.com/charmbracelet/bubbletea/blob/main/examples/tabs/main.go

8. Test thoroughly to ensure:
   - Tab switching works smoothly
   - Screen state is preserved
   - No visual glitches

Return a summary of implementation, files changed, and any design decisions made.
```

---

## Phase 3: Editor Simplification - Two-Section Design
**Duration**: 3-4 days | **Complexity**: High | **Risk**: Medium

### What Will Be Accomplished
- Remove complex editor code
- Implement table view for nodes (Bubbles table)
- Create Huh form integration for node editing
- Support nested display for loops/groups
- Feature-flag new editor for safe rollout

### Related Original TODOs
- TODO #3: Complete editor redesign

### New Editor Structure
```
┌─────────────────────────────────┐
│ Node Table (Bubbles Table)      │
│ ┌─────┬──────┬─────────────┐   │
│ │Name │Type  │Parameters    │   │
│ ├─────┼──────┼─────────────┤   │
│ │http1│HTTP  │GET /api/...  │   │
│ │loop1│Loop  │[nested table]│   │
│ └─────┴──────┴─────────────┘   │
├─────────────────────────────────┤
│ Node Editor (Huh Forms)         │
│ [Dynamic form based on node]    │
└─────────────────────────────────┘
```

### Implementation Details
- Create new package `pkg/omise/screens/editor_v2/`
- Split into focused files: `editor.go`, `table.go`, `forms.go`, `messages.go`
- Each file < 250 lines
- Keep old editor during transition with feature flag
- References:
  - https://github.com/charmbracelet/bubbletea/blob/main/examples/table/main.go
  - https://github.com/charmbracelet/huh

### Claude Code Prompt for Guilliman
```
Implement Phase 3 of the Bento TUI enhancement plan: Editor Simplification - Two-Section Design.

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

---

## Phase 4: Documentation & Demo - VHS Integration
**Duration**: 1 day | **Complexity**: Low | **Risk**: Minimal

### What Will Be Accomplished
- Add VHS scripts for capturing TUI demos
- Create animated GIFs for README
- Document key workflows visually
- Create demo scripts for common use cases

### Related Original TODOs
- TODO #5: VHS Charm functionality

### Implementation Details
- Create `.vhs/` directory with tape scripts
- Create demos for: browsing bentos, creating new bento, editing nodes, running bento
- Generate GIFs for README
- Reference: https://github.com/charmbracelet/vhs

### Claude Code Prompt for Guilliman
```
Implement Phase 4 of the Bento TUI enhancement plan: Documentation & Demo - VHS Integration.

Requirements:
1. Create .vhs/ directory with VHS tape scripts:

   a. demo-overview.tape:
      - Show main tab navigation
      - Browse existing bentos
      - Quick tour of all tabs

   b. demo-create-bento.tape:
      - Create new bento from scratch
      - Add multiple nodes
      - Save and run

   c. demo-editor.tape:
      - Show new editor interface (if Phase 3 complete)
      - Table navigation
      - Form editing

   d. demo-recipes.tape:
      - Browse example bentos
      - Copy example to own collection
      - Customize and run

2. VHS script features:
   - Set appropriate terminal size (1200x800)
   - Use readable font size (14-16)
   - Include pauses for readability
   - Add typed commands with realistic timing
   - Capture key interactions

3. Documentation:
   - Create .vhs/README.md with instructions
   - Document how to generate GIFs
   - List prerequisites (VHS installation)
   - Provide regeneration commands

4. Update main README.md:
   - Add animated GIF demos
   - Show key features visually
   - Link to detailed documentation

5. Reference:
   - https://github.com/charmbracelet/vhs

Example tape structure:
```tape
Output demo-overview.gif
Set FontSize 14
Set Width 1200
Set Height 800
Set Theme "Catppuccin Mocha"

Type "bento taste"
Sleep 2s
Enter
Sleep 1s
Type "1"  # Switch to Bentos tab
Sleep 1s
Screenshot
```

6. Generate all GIFs and verify:
   - Clear and readable
   - Show key features
   - Appropriate length (5-15s each)
   - High quality rendering

Return a summary of:
- Tape scripts created
- GIFs generated
- README updates made
- Instructions for regenerating demos
```

---

## Phase 5: Code Quality - Best Practices Review
**Duration**: 2 days | **Complexity**: Medium | **Risk**: Low

### What Will Be Accomplished
- Apply leg100's Bubbletea best practices
- Refactor message handling patterns
- Improve state management (immutable updates)
- Optimize rendering pipeline
- Comprehensive code quality review

### Related Original TODOs
- TODO #6: Best practices review

### Key Improvements
- Command batching patterns
- Immutable state updates
- Efficient view composition
- Proper error handling
- Testing coverage

### Implementation Details
- Review and refactor all Update() methods
- Ensure proper command batching
- Verify immutable state updates
- Optimize View() rendering
- Add comprehensive tests
- Reference: https://leg100.github.io/en/posts/building-bubbletea-programs/

### Claude Code Prompt for Guilliman
```
Implement Phase 5 of the Bento TUI enhancement plan: Code Quality - Best Practices Review.

Requirements:
1. Read and apply best practices from:
   - https://leg100.github.io/en/posts/building-bubbletea-programs/

2. Message Handling Review:
   - Audit all Update() methods across the codebase
   - Ensure proper command batching with tea.Batch()
   - Verify no missed message handling
   - Check for proper message routing
   - Fix any message handling anti-patterns

3. State Management Review:
   - Ensure all model updates are immutable
   - No direct field mutations in Update()
   - Return new model instances, not mutated originals
   - Review all struct updates for proper Go patterns

4. Rendering Optimization:
   - Audit all View() methods
   - Minimize string allocations
   - Cache computed layouts where appropriate
   - Ensure no business logic in View()
   - Optimize string building with strings.Builder

5. Error Handling:
   - Ensure all errors properly propagated
   - Consistent error message patterns
   - User-friendly error display
   - No silent failures
   - Proper error wrapping with fmt.Errorf("%w")

6. Code Quality Checks:
   - Run golangci-lint and fix all issues
   - Verify all files < 250 lines (Bento Box)
   - Verify all functions < 20 lines
   - Check for circular dependencies
   - Ensure proper package organization

7. Testing:
   - Add unit tests for components
   - Test message handling flows
   - Test state transitions
   - Test rendering edge cases
   - Target > 80% coverage for business logic

8. Performance:
   - Profile key operations
   - Identify bottlenecks
   - Optimize hot paths
   - Reduce unnecessary allocations

9. Documentation:
   - Add package documentation
   - Document complex logic
   - Add examples for key components
   - Update architecture docs

10. Review checklist:
    - [ ] All Update() methods use tea.Batch() properly
    - [ ] All state updates are immutable
    - [ ] No business logic in View() methods
    - [ ] Error handling is consistent
    - [ ] All files < 250 lines
    - [ ] All functions < 20 lines
    - [ ] golangci-lint passes
    - [ ] Tests added for new functionality
    - [ ] Documentation updated

Return a comprehensive report:
- Issues found and fixed
- Performance improvements made
- Test coverage achieved
- Remaining technical debt
- Recommendations for future improvements
```

---

## Timeline & Execution Order

### Recommended Sequence
- **Week 1**: Phase 1 & 2 (Foundation + Navigation)
- **Week 2**: Phase 3 (Editor Redesign)
- **Week 3**: Phase 4 & 5 (Documentation + Quality)

### Critical Path
Phase 1 → Phase 2 → Phase 3 (Must be sequential)
Phase 4 & 5 can run in parallel after Phase 3

---

## Success Metrics

### Code Quality
- [ ] All files < 250 lines
- [ ] All functions < 20 lines
- [ ] Zero circular dependencies
- [ ] golangci-lint passes
- [ ] Test coverage > 80%

### User Experience
- [ ] Tab navigation intuitive
- [ ] Help system comprehensive
- [ ] Editor simpler to use
- [ ] Visual feedback clear

### Bento Box Compliance
- [ ] Single responsibility per package
- [ ] No utils/helpers packages
- [ ] Clear module boundaries
- [ ] Composable components

---

---

## Phase 6: Remove Shift+C Key from Executor
**Duration**: < 1 hour | **Complexity**: Very Low | **Risk**: Minimal

### What Will Be Accomplished
- Remove the Shift+C keybinding from the executor screen
- Clean up any related key handling code
- Update help documentation if needed

### Claude Code Prompt for Guilliman
```
Remove the Shift+C key binding from the executor screen in the Bento TUI.

Requirements:
1. Locate the executor screen implementation
2. Remove Shift+C key handling from the Update() method
3. Remove any references to Shift+C in help text or key maps
4. Test to ensure the key no longer triggers any action
5. Verify no other functionality is affected

Return a summary of:
- Files modified
- What the Shift+C key was doing (for documentation)
- Confirmation that it's been removed
```

---

## Notes
- Each phase builds upon the previous
- Prompts are designed for Guilliman to execute independently
- Follow Go idioms and Bento Box principles throughout
- Test thoroughly between phases
- Keep old code during transitions with feature flags
- Individual phase documents are located in `.claude/strategy/`
- Completed phases are moved to `.claude/strategy/archive/` after code review and explicit user approval
- **IMPORTANT**: No work should be auto-committed for any of these phases
- **IMPORTANT**: Must wait for explicit user consent before moving forward with implementation
