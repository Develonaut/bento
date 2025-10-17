# Phase 1: Foundation - Help & Key System Enhancement

**Duration**: 1-2 days | **Complexity**: Low | **Risk**: Low | **Status**: Not Started

## Overview
Implement proper Bubbles help and key components to establish a foundation for better user guidance and consistent keybindings across all screens.

## Goals
- Create centralized key management system
- Implement proper Bubbles help component
- Standardize keybindings across all screens
- Create consistent footer help menu
- Establish foundation for better user guidance

## Related Original TODOs
- TODO #2: Bubbles Help and Key components
- Partial TODO #4: Header styling preparation

## Implementation Details

### Files to Create
- `pkg/omise/components/keys.go` - Centralized key definitions

### Files to Modify
- `pkg/omise/components/footer.go` - Enhanced with help.Model integration
- `pkg/omise/screens/browser.go` - Use centralized keys
- `pkg/omise/screens/editor.go` - Use centralized keys
- `pkg/omise/screens/settings.go` - Use centralized keys
- `pkg/omise/screens/helpview.go` - Use centralized keys

### Key Structure
```go
// pkg/omise/components/keys.go
type KeyMap struct {
    Navigation NavigationKeys
    Browser    BrowserKeys
    Editor     EditorKeys
    Settings   SettingsKeys
    Global     GlobalKeys
}

type NavigationKeys struct {
    Up    key.Binding
    Down  key.Binding
    Left  key.Binding
    Right key.Binding
}

// ... other key groups
```

### Footer Enhancement
```go
// pkg/omise/components/footer.go
type Footer struct {
    help   help.Model
    keys   KeyMap
    width  int
}

func (f Footer) View() string {
    return f.help.View(f.keys)
}
```

## Dependencies
- Already have `charmbracelet/bubbles` dependency
- Existing `helpview.go` can be enhanced
- No new external dependencies required

## Bento Box Compliance
- [ ] All files < 250 lines
- [ ] All functions < 20 lines
- [ ] No circular dependencies
- [ ] Clear separation of concerns

## Testing Requirements
- [ ] All screens still function correctly
- [ ] Key bindings work as expected
- [ ] Help text displays properly
- [ ] Footer responds to width changes
- [ ] No regressions in existing functionality

## Success Criteria
- [ ] Centralized KeyMap implemented
- [ ] All screens use centralized keys
- [ ] Footer displays help using Bubbles help.Model
- [ ] Help text is consistent across screens
- [ ] Full/short help modes work
- [ ] All tests pass

## References
- https://github.com/charmbracelet/bubbletea/blob/main/examples/help/main.go
- https://github.com/charmbracelet/bubbles?tab=readme-ov-file#key

## Notes
- This phase establishes the foundation for Phase 2 (tab navigation)
- Keep changes minimal and focused on key management
- Ensure backward compatibility with existing screens

---

## Claude Code Prompt for Guilliman

```
Implement Phase 1 of the Bento TUI enhancement plan: Foundation - Help & Key System Enhancement.

IMPORTANT: Read the phase document at .claude/strategy/phase-1-foundation.md for complete context, requirements, and success criteria before starting implementation.

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
