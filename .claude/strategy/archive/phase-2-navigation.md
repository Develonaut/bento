# Phase 2: Navigation - Tab System Implementation

**Duration**: 2-3 days | **Complexity**: Medium | **Risk**: Low | **Status**: Not Started

## Overview
Implement tab-based navigation with the Kitchen Workflow theme to organize the main views of the Bento TUI.

## Goals
- Implement tab-based navigation system
- Create visual tab indicators in header
- Maintain existing screen functionality
- Add keyboard shortcuts for tab switching
- Apply Kitchen Workflow naming theme

## Tab Names (Kitchen Workflow Theme)
- **🍱 Bentos** → Your prepared bento boxes (Home/Browser)
- **📖 Recipes** → Browse example bentos (Discover/Examples)
- **🔪 Mise** → Your workspace setup (Settings)
- **👨‍🍳 Sensei** → Guidance from the master (Help)

## Related Original TODOs
- TODO #1: Bubbletea Tabs implementation
- Complete TODO #4: Full header styling with tabs

## Visual Design

### Tab Structure
```
🍱 Bento v0.0.1
┌─────────┬─────────┬────────┬─────────┐
│🍱 Bentos│📖 Recipes│🔪 Mise │👨‍🍳 Sensei│
└─────────┴─────────┴────────┴─────────┘
```

### Active Tab Indicator
```
🍱 Bento v0.0.1
┌═════════┬─────────┬────────┬─────────┐
│🍱 Bentos│📖 Recipes│🔪 Mise │👨‍🍳 Sensei│ ← Active tab highlighted
└═════════┴─────────┴────────┴─────────┘
```

## Implementation Details

### Files to Create
- `pkg/omise/components/tabs.go` - Tab navigation logic

### Files to Modify
- `pkg/omise/components/header.go` - Display tabs
- `pkg/omise/omise.go` - Handle tab switching
- `pkg/omise/components/keys.go` - Add tab navigation keys
- Screen mappings for each tab

### Tab Component Structure
```go
// pkg/omise/components/tabs.go
type TabView struct {
    activeTab int
    tabs      []Tab
    width     int
}

type Tab struct {
    Name   string
    Icon   string
    Screen Screen
}

func (t TabView) View() string {
    // Render tabs with active state
}
```

### Key Bindings
- `1` - Switch to Bentos tab
- `2` - Switch to Recipes tab
- `3` - Switch to Mise tab
- `4` - Switch to Sensei tab
- `Tab` - Cycle through tabs

## Screen Mappings
- `Screen.Browser` → Bentos tab
- `Screen.Examples` → Recipes tab (if exists, otherwise plan for future)
- `Screen.Settings` → Mise tab
- `Screen.Help` → Sensei tab

## Dependencies
- Phase 1 completion (centralized key system)
- No new external dependencies

## Bento Box Compliance
- [ ] All files < 250 lines
- [ ] All functions < 20 lines
- [ ] Clear separation of concerns
- [ ] No circular dependencies

## Testing Requirements
- [ ] Tab switching works smoothly
- [ ] Screen state is preserved when switching tabs
- [ ] Keyboard shortcuts work correctly
- [ ] Visual indicators display properly
- [ ] No visual glitches during transitions
- [ ] All existing screen functionality intact

## Success Criteria
- [ ] Tab component implemented
- [ ] Header displays tabs with icons
- [ ] Active tab visually distinct
- [ ] All four tabs functional
- [ ] Keyboard navigation works
- [ ] Screen state preserved between switches
- [ ] Responsive to terminal width
- [ ] All tests pass

## References
- https://github.com/charmbracelet/bubbletea/blob/main/examples/tabs/main.go

## Notes
- Depends on Phase 1 completion
- If Recipes screen doesn't exist, create placeholder or stub
- Ensure smooth transitions between tabs
- Maintain existing screen functionality

---

## Claude Code Prompt for Guilliman

```
Implement Phase 2 of the Bento TUI enhancement plan: Navigation - Tab System Implementation.

IMPORTANT: Read the phase document at .claude/strategy/phase-2-navigation.md for complete context, requirements, and success criteria before starting implementation.

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
