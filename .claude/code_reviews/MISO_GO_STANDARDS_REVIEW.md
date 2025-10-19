# Miso Package - Go Standards Review

**Date**: 2025-10-19
**Reviewer**: Colossus (Go Standards Guardian)
**Package**: `pkg/miso/`
**Phase**: 7 (Charm CLI Integration)

---

## Executive Summary

**Overall Grade: A- (Excellent)**

The miso package demonstrates strong adherence to Go idioms and the Bento Box Principle. The code is clean, well-tested, and properly structured. Minor issues identified are primarily around **global state management** and **potential race conditions** that should be addressed.

**Key Strengths**:
- Excellent file size discipline (all files < 250 lines)
- 100% test coverage with proper mocking
- Effective use of stdlib (minimal dependencies)
- Clean, flat package structure
- Idiomatic error handling

**Critical Issues**: 1 (Race condition in theme.go)
**Medium Issues**: 2 (Global state patterns)
**Minor Issues**: 3 (Naming, comments)

---

## Critical Issues

### 1. RACE CONDITION: Unsynchronized Global State in `theme.go`

**Severity**: CRITICAL
**Location**: `/Users/Ryan/Code/bento/pkg/miso/theme.go:18-19, 36-38`

```go
// PROBLEM: Unsynchronized global variable
var currentTheme *Theme

// Called from multiple goroutines potentially
func ApplyPalette(p Palette) {
    currentTheme = buildTheme(p)  // RACE!
}

func GetTheme() *Theme {
    if currentTheme == nil {        // RACE!
        palette := GetPalette(VariantMaguro)
        currentTheme = buildTheme(palette)
    }
    return currentTheme              // RACE!
}
```

**Why This Is Critical**:
- `currentTheme` is a package-level global variable
- Both `ApplyPalette()` and `GetTheme()` read/write without synchronization
- The race detector passed because tests don't have concurrent access
- In production CLI usage with TUI (bubbletea), multiple goroutines will access this

**Go Proverb Violated**: "Don't communicate by sharing memory; share memory by communicating"

**Solutions** (Pick ONE):

#### Option A: Use sync.RWMutex (Simplest)
```go
var (
    currentTheme *Theme
    themeMu      sync.RWMutex
)

func GetTheme() *Theme {
    themeMu.RLock()
    if currentTheme != nil {
        defer themeMu.RUnlock()
        return currentTheme
    }
    themeMu.RUnlock()

    themeMu.Lock()
    defer themeMu.Unlock()

    // Double-check after acquiring write lock
    if currentTheme == nil {
        palette := GetPalette(VariantMaguro)
        currentTheme = buildTheme(palette)
    }
    return currentTheme
}

func ApplyPalette(p Palette) {
    themeMu.Lock()
    defer themeMu.Unlock()
    currentTheme = buildTheme(p)
}
```

#### Option B: Use sync.Once + atomic.Value (More complex, lock-free reads)
```go
var (
    currentTheme atomic.Value // stores *Theme
    themeOnce    sync.Once
)

func GetTheme() *Theme {
    themeOnce.Do(func() {
        palette := GetPalette(VariantMaguro)
        currentTheme.Store(buildTheme(palette))
    })
    return currentTheme.Load().(*Theme)
}

func ApplyPalette(p Palette) {
    currentTheme.Store(buildTheme(p))
}
```

#### Option C: Eliminate Global State (Most Go-idiomatic)
```go
// Remove global variable entirely
// Make Theme part of Manager instead

type Manager struct {
    variant Variant
    palette Palette
    theme   *Theme  // Add theme here
}

func (m *Manager) GetTheme() *Theme {
    return m.theme
}

func (m *Manager) SetVariant(v Variant) {
    m.variant = v
    m.palette = GetPalette(v)
    m.theme = buildTheme(m.palette)  // Update theme
    currentVariant = v
    _ = SaveTheme(v)
}
```

**Recommendation**: Use **Option C** - eliminate global state. This aligns with "Accept interfaces, return structs" and makes the Manager the single source of truth.

---

## Medium Issues

### 2. Global State in `manager.go`

**Severity**: MEDIUM
**Location**: `/Users/Ryan/Code/bento/pkg/miso/manager.go:12-13, 50-52`

```go
// Global mutable state
var currentVariant Variant

func (m *Manager) SetVariant(v Variant) {
    // ...
    currentVariant = v  // Why is this needed?
}

func CurrentVariant() Variant {
    return currentVariant
}
```

**Issues**:
1. **Unclear ownership**: Who owns `currentVariant`? The Manager or the package?
2. **Redundant state**: Manager already tracks variant in `m.variant`
3. **No synchronization**: Same race condition risk as theme.go
4. **Testing complexity**: Global state makes tests harder to isolate

**Why This Exists**:
The `CurrentVariant()` function provides package-level access to the variant. This is useful for `sequence.go` which calls `GetTheme()` without having a Manager reference.

**Solution**:
```go
// REMOVE the global variable entirely
// REMOVE CurrentVariant() function

// In sequence.go, pass Manager or Theme explicitly:
func (s *Sequence) ViewWithTheme(theme *Theme) string {
    // Use theme instead of GetTheme()
}

// OR: Accept theme in NewSequence
func NewSequence(theme *Theme) *Sequence {
    return &Sequence{
        steps: []Step{},
        theme: theme,
    }
}
```

**Recommendation**: Eliminate `currentVariant` global. Pass theme/manager explicitly through function parameters. This follows "Clear is better than clever".

---

### 3. `init()` Usage in `theme.go`

**Severity**: MEDIUM
**Location**: `/Users/Ryan/Code/bento/pkg/miso/theme.go:21-25`

```go
// init initializes theme with default Maguro palette.
func init() {
    palette := GetPalette(VariantMaguro)
    currentTheme = buildTheme(palette)
}
```

**Issues**:
1. **Hidden initialization**: `init()` runs at import time, not obvious to callers
2. **Testing complexity**: Can't control when this runs in tests
3. **Global state initialization**: Sets up global state before any configuration

**Go Wisdom**: "Make the zero value useful" and avoid `init()` when possible.

**Better Approach**:
```go
// Remove init() entirely

func GetTheme() *Theme {
    if currentTheme == nil {
        palette := GetPalette(VariantMaguro)
        currentTheme = buildTheme(palette)
    }
    return currentTheme
}

// Even better: Use sync.Once
var (
    currentTheme *Theme
    themeOnce    sync.Once
)

func GetTheme() *Theme {
    themeOnce.Do(func() {
        palette := GetPalette(VariantMaguro)
        currentTheme = buildTheme(palette)
    })
    return currentTheme
}
```

The `GetTheme()` function already has lazy initialization logic (lines 28-33), making the `init()` redundant.

**Recommendation**: Remove `init()`. The lazy initialization in `GetTheme()` is sufficient and more testable.

---

## Minor Issues

### 4. Mocking Pattern in `config.go`

**Severity**: MINOR
**Location**: `/Users/Ryan/Code/bento/pkg/miso/config.go:12-20`

```go
// configDir returns the bento config directory path.
// Mutable var allows mocking in tests.
var configDir = func() (string, error) {
    home, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }
    return filepath.Join(home, ".bento"), nil
}
```

**Assessment**: ACCEPTABLE

This is a pragmatic approach to mocking for tests. While purists might prefer dependency injection, this pattern:
- Works well for simple cases
- Is explicitly documented ("Mutable var allows mocking in tests")
- Tests properly restore original value with `t.Cleanup()`
- Avoids over-engineering for a simple config path

**Alternative** (if you wanted to be more "pure"):
```go
type ConfigProvider interface {
    ConfigDir() (string, error)
}

type defaultConfigProvider struct{}

func (d defaultConfigProvider) ConfigDir() (string, error) {
    // ...
}
```

But this adds complexity for minimal benefit in a CLI tool.

**Verdict**: Keep as-is. This is fine for the Bento project's needs.

---

### 5. Function Naming: `GetPalette` vs `Palette`

**Severity**: MINOR
**Location**: `/Users/Ryan/Code/bento/pkg/miso/variants.go:45`

```go
func GetPalette(v Variant) Palette {
    // ...
}
```

**Go Convention**: "Get" prefix is often unnecessary in Go.

**More Idiomatic**:
```go
func (v Variant) Palette() Palette {
    switch v {
    case VariantNasu:
        return nasuPalette()
    // ...
    }
}

// Usage:
palette := VariantMaguro.Palette()
```

This follows the Go pattern of making Variant a first-class type with methods.

**Impact**: Low - current naming is still acceptable and widely understood.

**Recommendation**: Consider refactoring to method on Variant type, but not urgent.

---

### 6. Exported Helper Functions

**Severity**: MINOR
**Location**: `/Users/Ryan/Code/bento/pkg/miso/variants.go:22-32`

```go
// AllVariants returns all available theme variants in order.
func AllVariants() []Variant {
    return []Variant{
        VariantNasu,
        // ...
    }
}
```

**Question**: Does this need to be exported?

Looking at usage:
- Used in tests extensively
- Used in Manager.NextVariant() (line 57 of manager.go)
- May be used by CLI commands (theme selection UI)

**Verdict**: Export is justified. This is part of the public API.

**Suggestion**: Add a comment explaining the ordering matters for `NextVariant()` cycling:
```go
// AllVariants returns all available theme variants in display order.
// The order determines the sequence for NextVariant() cycling.
func AllVariants() []Variant {
```

---

### 7. Magic Numbers in `sequence.go`

**Severity**: MINOR
**Location**: `/Users/Ryan/Code/bento/pkg/miso/sequence.go:120, 155`

```go
indent := strings.Repeat("  ", step.Depth)  // "  " is magic

if step.Duration > 0 {  // 0 is arbitrary threshold
```

**Better**:
```go
const (
    indentSpaces = "  "
    minDurationDisplay = 0
)

indent := strings.Repeat(indentSpaces, step.Depth)
if step.Duration > minDurationDisplay {
```

**Impact**: Very low - these are reasonable inline literals.

**Recommendation**: Optional improvement, not required.

---

## Strengths (Things Done Right)

### 1. Excellent Bento Box Compliance

**File Sizes**:
```
  66 pkg/miso/theme.go          ✓
  71 pkg/miso/manager.go         ✓
  75 pkg/miso/config.go          ✓
 157 pkg/miso/variants.go        ✓
 220 pkg/miso/sequence.go        ✓
```

All files under 250 lines. Largest file (sequence.go) is well-organized with clear sections.

---

### 2. Proper Error Handling

All error handling follows Go conventions:
- Errors returned as last return value
- Errors checked at call site
- Fallback to defaults where appropriate (LoadSavedTheme)
- No panics in library code

Example from config.go:
```go
func LoadSavedTheme() Variant {
    path, err := themeConfigPath()
    if err != nil {
        return VariantMaguro  // Sensible default
    }
    // ...
}
```

---

### 3. Effective Use of Stdlib

Dependencies:
- Standard library: `fmt`, `os`, `path/filepath`, `strings`, `time`, `hash/fnv`
- External: Only `lipgloss` (required for the domain)

No unnecessary dependencies. Good adherence to "A little copying is better than a little dependency".

Example of good stdlib usage (hash/fnv for deterministic emoji selection):
```go
func getStepEmoji(stepName string) string {
    h := fnv.New32a()
    h.Write([]byte(stepName))
    hash := h.Sum32()
    return sushiEmojis[hash%uint32(len(sushiEmojis))]
}
```

This is clever but clear - uses stdlib hash for deterministic randomness.

---

### 4. Clean Package Structure

```
pkg/miso/
├── config.go       # Persistence
├── manager.go      # State management
├── sequence.go     # Display logic
├── theme.go        # Style definitions
└── variants.go     # Color palettes
```

Each file has a single, clear responsibility. No "utils.go" grab bag.

---

### 5. Excellent Test Coverage

Test highlights:
- Table-driven tests where appropriate
- Proper use of `t.TempDir()` for filesystem tests
- Cleanup with `t.Cleanup()`
- Tests for edge cases (invalid content, missing files)
- Mocking done properly with variable substitution

Example of good test structure:
```go
func TestSaveAndLoadTheme(t *testing.T) {
    tmpDir := t.TempDir()
    originalConfigDir := configDir

    configDir = func() (string, error) {
        return tmpDir, nil
    }
    t.Cleanup(func() {
        configDir = originalConfigDir
    })

    for _, variant := range AllVariants() {
        t.Run(string(variant), func(t *testing.T) {
            // Test logic
        })
    }
}
```

---

### 6. Clear Type Definitions

Types are well-defined with clear semantics:
```go
type Variant string  // String-based enum

const (
    VariantNasu   Variant = "Nasu"   // Purple (eggplant sushi)
    // Clear naming and documentation
)

type Palette struct {
    Primary   lipgloss.Color // Main theme color (brand)
    Secondary lipgloss.Color // Accent color
    // Each field has clear purpose
}
```

---

### 7. Good Use of Receiver Names

```go
func (s *Sequence) AddStep(name, nodeType string)
func (m *Manager) SetVariant(v Variant)
func (p Palette) ...  // Would be good if methods existed
```

Receiver names are 1-2 letters, following Go convention.

---

## Go Proverbs Adherence

| Proverb | Status | Notes |
|---------|--------|-------|
| Simple is better than complex | ✓ | Code is straightforward, no over-abstraction |
| Clear is better than clever | ✓ | Clear intent throughout |
| Errors are values | ✓ | Proper error handling everywhere |
| Don't panic | ✓ | No panics in library code |
| Accept interfaces, return structs | ✓ | Returns concrete types (Manager, Sequence) |
| A little copying is better than a little dependency | ✓ | Minimal external dependencies |
| **Don't communicate by sharing memory** | ✗ | **Global state in theme.go and manager.go** |
| Make the zero value useful | ~ | Could be better (Manager requires NewManager) |

---

## Specific Recommendations

### Immediate (Before Merge)

1. **FIX: Race condition in theme.go**
   - Add sync.RWMutex OR eliminate global state
   - This is a production bug waiting to happen

2. **CONSIDER: Remove currentVariant global in manager.go**
   - Pass theme/manager explicitly
   - Makes ownership clear

3. **REMOVE: init() in theme.go**
   - Lazy initialization is already implemented
   - Less magic, more testable

### Soon (Next Refactor)

4. **Refactor: Make GetPalette a method on Variant**
   ```go
   palette := variant.Palette()
   ```

5. **Add: Concurrency test**
   ```go
   func TestThemeConcurrency(t *testing.T) {
       var wg sync.WaitGroup
       for i := 0; i < 100; i++ {
           wg.Add(1)
           go func() {
               defer wg.Done()
               theme := GetTheme()
               _ = theme.Title.Render("test")
           }()
       }
       wg.Wait()
   }
   ```

6. **Consider: Make Manager methods thread-safe**
   - If Manager will be shared across goroutines in TUI

---

## Final Verdict

**Grade: A- (Excellent with Critical Race Condition)**

The miso package demonstrates strong Go fundamentals and excellent code organization. The critical issue is the **unsynchronized global state** in theme.go, which will cause race conditions in concurrent usage (TUI with bubbletea).

**Before merging Phase 7**, fix the race condition. Everything else is optional polish.

**What Colossus Says**:
> "This is clean, well-tested code that respects the Bento Box Principle. Fix the race condition in theme.go, and you'll have A+ code. The Go community would be proud of this package structure - no utils packages, clean boundaries, proper testing. Just eliminate the global state, and you're golden."

---

## Questions for Design Review

1. **Will Manager be shared across goroutines?**
   - If yes, add mutex to Manager
   - If no, document "not thread-safe" clearly

2. **Should Theme be immutable?**
   - Would make concurrent access trivial
   - `ApplyPalette()` could return new Theme instead of mutating

3. **Is CurrentVariant() function needed?**
   - Who calls it?
   - Can we pass Manager/Theme explicitly instead?

---

## References

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Proverbs](https://go-proverbs.github.io/)
- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)
- [Bento Box Principle](./../BENTO_BOX_PRINCIPLE.md)

---

**Next Steps**:
1. Address critical race condition
2. Run `go test -race` again
3. Review global state usage
4. Document thread-safety guarantees
