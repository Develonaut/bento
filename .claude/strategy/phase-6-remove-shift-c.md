# Phase 6: Remove Shift+C Key from Executor

**Duration**: < 1 hour | **Complexity**: Very Low | **Risk**: Minimal | **Status**: Not Started

## Overview
Remove the Shift+C keybinding from the executor screen, as it's no longer needed or causing issues.

## Goals
- Remove Shift+C key handling from executor
- Clean up any related code
- Update help documentation if needed
- Ensure no regressions

## Related Original TODOs
- User request to remove Shift+C from executor

## Implementation Details

### Files to Investigate
- `pkg/omise/screens/executor.go` - Main executor implementation
- Key handling in `Update()` method
- Help text/key map references

### What to Remove
1. Shift+C case in `Update()` switch statement
2. Any help text mentioning Shift+C
3. Key binding definition (if in centralized keys)

### Example Change
```go
// Before
case tea.KeyMsg:
    switch {
    case key.Matches(msg, keys.ShiftC):
        // Do something
        return m, someCmd
    }

// After
// Remove the entire case block
```

## Dependencies
- None

## Testing Requirements
- [ ] Shift+C no longer triggers any action in executor
- [ ] No other keys affected
- [ ] Help text updated (if applicable)
- [ ] Executor still functions normally

## Success Criteria
- [ ] Shift+C key handling removed from executor
- [ ] Related code cleaned up
- [ ] Help text updated (if needed)
- [ ] No regressions in executor functionality
- [ ] Testing confirms removal

## Notes
- Very simple change
- Can be done independently
- Low risk of issues
- Quick win

---

## Claude Code Prompt for Guilliman

```
Remove the Shift+C key binding from the executor screen in the Bento TUI.

IMPORTANT: Read the phase document at .claude/strategy/phase-6-remove-shift-c.md for complete context, requirements, and success criteria before starting implementation.

Requirements:
1. Locate the executor screen implementation (likely pkg/omise/screens/executor.go)
2. Remove Shift+C key handling from the Update() method
3. Remove any references to Shift+C in help text or key maps
4. Test to ensure the key no longer triggers any action
5. Verify no other functionality is affected

Return a summary of:
- Files modified
- What the Shift+C key was doing (for documentation)
- Confirmation that it's been removed
- Any issues encountered
```
