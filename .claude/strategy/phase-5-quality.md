# Phase 5: Code Quality - Best Practices Review

**Duration**: 2 days | **Complexity**: Medium | **Risk**: Low | **Status**: Not Started

## Overview
Comprehensive code quality review applying Bubbletea best practices, improving message handling, state management, and overall code quality.

## Goals
- Apply leg100's Bubbletea best practices
- Refactor message handling patterns
- Improve state management (immutable updates)
- Optimize rendering pipeline
- Enhance error handling
- Improve test coverage
- Optimize performance

## Related Original TODOs
- TODO #6: Best practices review

## Key Areas to Address

### 1. Message Handling
- Proper command batching with `tea.Batch()`
- No missed message handling
- Correct message routing
- Fix anti-patterns

### 2. State Management
- Ensure immutable updates
- No direct field mutations in `Update()`
- Return new model instances
- Proper Go struct update patterns

### 3. Rendering Optimization
- Minimize string allocations
- Cache computed layouts where appropriate
- No business logic in `View()`
- Use `strings.Builder` for concatenation

### 4. Error Handling
- Proper error propagation
- Consistent error patterns
- User-friendly error display
- No silent failures
- Error wrapping with `fmt.Errorf("%w")`

### 5. Code Quality
- Run `golangci-lint` and fix issues
- Verify all files < 250 lines (Bento Box)
- Verify all functions < 20 lines
- Check for circular dependencies
- Proper package organization

### 6. Testing
- Add unit tests for components
- Test message handling flows
- Test state transitions
- Test rendering edge cases
- Target > 80% coverage for business logic

### 7. Performance
- Profile key operations
- Identify bottlenecks
- Optimize hot paths
- Reduce unnecessary allocations

### 8. Documentation
- Add package documentation
- Document complex logic
- Add examples for key components
- Update architecture docs

## Implementation Details

### Message Handling Example
```go
// Before
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        m.screen = NewScreen
        return m, loadScreenCmd()
    }
    return m, nil
}

// After
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        newModel := m
        newModel.screen = NewScreen
        cmds = append(cmds, loadScreenCmd())
        return newModel, tea.Batch(cmds...)
    }
    return m, nil
}
```

### State Management Example
```go
// Before (mutating)
func (m Model) updateField(value string) Model {
    m.field = value  // Direct mutation
    return m
}

// After (immutable)
func (m Model) updateField(value string) Model {
    return Model{
        field: value,
        // Copy other fields
        width: m.width,
        height: m.height,
    }
}
```

## Files to Review

### High Priority
- `pkg/omise/omise.go` - Main update loop
- `pkg/omise/screens/*.go` - All screen Update/View methods
- `pkg/omise/components/*.go` - Component rendering
- `pkg/neta/*.go` - Core business logic

### Testing Priority
- Message handling in screens
- State transitions
- Form validation
- Node operations
- File I/O operations

## Dependencies
- All previous phases recommended (but not required)
- `golangci-lint` installed
- Go testing tools

## Bento Box Compliance
- [ ] All files < 250 lines
- [ ] All functions < 20 lines
- [ ] Zero circular dependencies
- [ ] Clear package boundaries

## Review Checklist
- [ ] All `Update()` methods use `tea.Batch()` properly
- [ ] All state updates are immutable
- [ ] No business logic in `View()` methods
- [ ] Error handling is consistent
- [ ] All files < 250 lines
- [ ] All functions < 20 lines
- [ ] `golangci-lint` passes
- [ ] Tests added for new functionality
- [ ] Documentation updated
- [ ] Performance profiled
- [ ] No obvious bottlenecks

## Success Criteria
- [ ] All code quality checks pass
- [ ] Test coverage > 80% for business logic
- [ ] `golangci-lint` passes with zero issues
- [ ] All Bento Box constraints met
- [ ] Performance improved or maintained
- [ ] Documentation complete
- [ ] Best practices applied consistently

## References
- https://leg100.github.io/en/posts/building-bubbletea-programs/
- https://github.com/charmbracelet/bubbletea/tree/main/tutorials
- Go standard library documentation

## Notes
- This is an optimization and quality phase
- Can be done incrementally
- Focus on high-impact improvements first
- Don't over-engineer simple solutions
- Apply YAGNI (You Aren't Gonna Need It)

---

## Claude Code Prompt for Guilliman

```
Implement Phase 5 of the Bento TUI enhancement plan: Code Quality - Best Practices Review.

IMPORTANT: Read the phase document at .claude/strategy/phase-5-quality.md for complete context, requirements, and success criteria before starting implementation.

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
