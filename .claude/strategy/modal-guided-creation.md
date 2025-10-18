# Modal-Based Guided Creation Implementation

## Research Summary

After researching Charm libraries, there are **three viable approaches** for containing huh forms in a modal that prevents key bubbling:

## Option 1: Huh WithWidth/WithHeight + Lipgloss Place (Recommended)

### How It Works
1. **Huh provides sizing options**: `WithWidth()` and `WithHeight()` methods on forms
2. **Lipgloss provides positioning**: `Place()` function centers content in a defined space
3. **Bubble Tea handles routing**: Conditional Update() routing prevents key bubbling

### Implementation

```go
// In browser.go - Add modal state
type Browser struct {
    list          components.StyledList
    store         *jubako.Store
    discovery     *jubako.Discovery
    confirmDialog *ConfirmDialog
    helpView      components.HelpView
    keys          components.BrowserKeyMap
    guidedModal   *GuidedModal  // NEW: Modal wrapper
    width         int
    height        int
}

// New modal wrapper
type GuidedModal struct {
    active   bool
    huhModel tea.Model // huh form running as tea.Model
    result   *neta.Definition
    err      error
}

// Update routing - prevents key bubbling
func (b Browser) Update(msg tea.Msg) (Browser, tea.Cmd) {
    // If modal is active, ONLY update modal (blocks all keys from browser)
    if b.guidedModal != nil && b.guidedModal.active {
        return b.updateModal(msg)
    }

    // Normal browser updates
    // ...
}

func (b Browser) updateModal(msg tea.Msg) (Browser, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Modal captures ALL keys - nothing bubbles to browser
        switch msg.String() {
        case "ctrl+c":
            // Allow quit
            return b, tea.Quit
        default:
            // All other keys go to huh form
            var cmd tea.Cmd
            b.guidedModal.huhModel, cmd = b.guidedModal.huhModel.Update(msg)
            return b, cmd
        }
    case GuidedCompleteMsg:
        // Modal finished
        b.guidedModal = nil
        return b.refreshList()
    default:
        var cmd tea.Cmd
        b.guidedModal.huhModel, cmd = b.guidedModal.huhModel.Update(msg)
        return b, cmd
    }
}

// View with overlay
func (b Browser) View() string {
    baseView := b.list.View()

    // If modal active, overlay it on top
    if b.guidedModal != nil && b.guidedModal.active {
        modalView := b.guidedModal.View()

        // Use lipgloss.Place to center modal
        return lipgloss.Place(
            b.width,
            b.height,
            lipgloss.Center,
            lipgloss.Center,
            modalView,
            lipgloss.WithWhitespaceChars("░"), // Semi-transparent background
        )
    }

    return baseView
}
```

### Guided Modal Implementation

```go
// pkg/omise/screens/guided_modal.go
package screens

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/huh"
    "github.com/charmbracelet/lipgloss"

    "bento/pkg/jubako"
    "bento/pkg/neta"
    "bento/pkg/omise/screens/guided"
)

type GuidedModal struct {
    active     bool
    width      int
    height     int
    store      *jubako.Store
    workDir    string

    // Huh form wrapped as tea.Model
    form       *huh.Form
    formModel  tea.Model

    result     *neta.Definition
    err        error
}

func NewGuidedModal(store *jubako.Store, workDir string, width, height int) *GuidedModal {
    // Create huh form with constrained size
    // Modal should be 80% of screen size
    modalWidth := int(float64(width) * 0.8)
    modalHeight := int(float64(height) * 0.8)

    return &GuidedModal{
        active:  true,
        width:   modalWidth,
        height:  modalHeight,
        store:   store,
        workDir: workDir,
    }
}

func (m *GuidedModal) Init() tea.Cmd {
    // Launch guided creation in background
    return func() tea.Msg {
        def, err := guided.CreateBentoGuided(m.store, m.workDir)
        return GuidedCompleteMsg{def: def, err: err}
    }
}

func (m *GuidedModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case GuidedCompleteMsg:
        m.result = msg.def
        m.err = msg.err
        m.active = false
        return m, nil
    }
    return m, nil
}

func (m *GuidedModal) View() string {
    if !m.active {
        return ""
    }

    // Styled modal box
    modalStyle := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("205")).
        Padding(1, 2).
        Width(m.width).
        Height(m.height)

    content := "Creating bento..."

    return modalStyle.Render(content)
}

type GuidedCompleteMsg struct {
    def *neta.Definition
    err error
}
```

### Pros
- ✅ **No external dependencies** - Uses built-in Charm libs
- ✅ **Full control** - Can customize modal appearance
- ✅ **Lightweight** - Minimal code overhead
- ✅ **Key blocking built-in** - Conditional Update() routing
- ✅ **Sizes respected** - huh WithWidth/WithHeight work

### Cons
- ⚠️ **Manual implementation** - Need to write modal wrapper
- ⚠️ **Huh integration complexity** - Need to run huh forms as tea.Model

---

## Option 2: bubbletea-overlay Package

### How It Works
Uses `github.com/quickphosphat/bubbletea-overlay` package which wraps two tea.Models (background + foreground).

### Implementation

```go
import "github.com/quickphosphat/bubbletea-overlay"

type Browser struct {
    // ...
    overlay *overlay.Model
}

func (b Browser) Update(msg tea.Msg) (Browser, tea.Cmd) {
    if b.overlay != nil {
        // Overlay handles routing automatically
        var cmd tea.Cmd
        *b.overlay, cmd = b.overlay.Update(msg)
        return b, cmd
    }

    // Normal update
}

func (b Browser) View() string {
    if b.overlay != nil {
        return b.overlay.View() // Composites foreground onto background
    }

    return b.list.View()
}

func (b Browser) handleNew() (Browser, tea.Cmd) {
    // Create overlay with browser as background, guided as foreground
    background := &browserModel{browser: b}
    foreground := &guidedModel{store: b.store, workDir: b.workDir}

    b.overlay = overlay.New(background, foreground)

    return b, b.overlay.Init()
}
```

### Pros
- ✅ **Pre-built solution** - Less code to write
- ✅ **Tested** - Community-maintained package
- ✅ **Simple API** - Just wrap two models

### Cons
- ⚠️ **External dependency** - Not an official Charm package
- ⚠️ **Less control** - Limited customization options
- ⚠️ **May not handle huh specially** - Huh still uses full terminal

---

## Option 3: Dedicated Guided Screen (Full Screen Modal)

### How It Works
Switch to a dedicated screen for guided creation, not a modal overlay.

### Implementation

```go
// Just switch screens
func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch msg.String() {
    case "n":
        // Switch to guided screen
        m.screen = ScreenGuided
        return m, m.guided.Init()
    }
}

// Guided screen gets full terminal
func (m Model) renderContent() string {
    switch m.screen {
    case ScreenBrowser:
        return m.browser.View()
    case ScreenGuided:
        return m.guided.View() // Full screen
    // ...
    }
}
```

### Pros
- ✅ **Simplest** - No modal complexity
- ✅ **Huh gets full space** - No sizing issues
- ✅ **Clean separation** - Clear mental model

### Cons
- ⚠️ **Not a modal** - Loses context of browser
- ⚠️ **No overlay effect** - Different UX pattern

---

## Recommended Approach: Option 1 (Huh + Lipgloss + Conditional Routing)

### Why This is Best

1. **Native Charm ecosystem** - No external deps
2. **Full control over UX** - Can customize modal appearance
3. **Proper key blocking** - Conditional Update() routing is idiomatic Bubble Tea
4. **Size control** - huh.WithWidth/WithHeight ensures containment
5. **Professional appearance** - lipgloss.Place centers modal perfectly

### The Key Insight

The "modal" isn't a special component - it's just **conditional message routing**:

```go
func (b Browser) Update(msg tea.Msg) (Browser, tea.Cmd) {
    if b.modalActive {
        // ALL messages go to modal only
        return b.updateModal(msg)
    }

    // Otherwise, normal browser updates
    return b.handleBrowserUpdate(msg)
}
```

This pattern **completely blocks** keys from reaching the browser because the browser Update() method returns early when the modal is active.

### Implementation Complexity

**Phase 1: Basic Modal (30 min)**
- Add `guidedModal *GuidedModal` field to Browser
- Implement conditional routing in Browser.Update()
- Create GuidedModal with Init/Update/View

**Phase 2: Huh Integration (1 hour)**
- Modify guided.CreateBentoGuided to work as tea.Model
- Or: Create separate huh forms that can be embedded
- Handle completion message

**Phase 3: Styling (30 min)**
- Add lipgloss modal border/padding
- Implement lipgloss.Place for centering
- Add semi-transparent background overlay

**Total**: ~2 hours for full implementation

---

## Alternative: Simplified Huh Modal Pattern

If integrating huh as a tea.Model is complex, use this simpler pattern:

```go
func (b Browser) handleNew() (Browser, tea.Cmd) {
    return b, func() tea.Msg {
        // Create huh form with size constraints
        form := createGuidedForm(b.width, b.height)

        // Run it (blocks)
        err := form.Run()
        if err != nil {
            return GuidedCancelledMsg{}
        }

        // Extract data and save
        def := extractDefinitionFromForm(form)
        b.store.Save(def)

        return GuidedCompleteMsg{def: def}
    }
}

func createGuidedForm(screenWidth, screenHeight int) *huh.Form {
    // Size form to 80% of screen
    formWidth := int(float64(screenWidth) * 0.8)
    formHeight := int(float64(screenHeight) * 0.8)

    return huh.NewForm(
        // ... groups
    ).WithWidth(formWidth).WithHeight(formHeight)
}
```

**This blocks the browser Update() during huh execution**, effectively preventing key bubbling.

When the command returns `GuidedCompleteMsg`, the browser refreshes the list.

---

## Testing Strategy

### Unit Tests
```go
func TestGuidedModal_KeyBlocking(t *testing.T) {
    modal := &GuidedModal{active: true}
    browser := Browser{guidedModal: modal}

    // Send tab key
    msg := tea.KeyMsg{Type: tea.KeyTab}
    browser, _ = browser.Update(msg)

    // Browser should NOT have switched tabs
    // Modal should have received the key
}
```

### Manual Tests
1. Start guided creation
2. Press Tab - should navigate form fields, NOT switch app tabs
3. Press 's' - should type 's', NOT open settings
4. Press '?' - should type '?', NOT open help
5. Press Ctrl+C - should quit app (only exception)
6. Complete form - should return to browser with new bento

---

## Implementation Plan

### Step 1: Create GuidedModal wrapper
- `pkg/omise/screens/guided_modal.go`
- Basic Init/Update/View

### Step 2: Add modal field to Browser
- `guidedModal *GuidedModal`
- Update Browser.Update() with conditional routing

### Step 3: Implement lipgloss overlay
- Center modal with lipgloss.Place
- Add styled border
- Add semi-transparent background

### Step 4: Integrate huh forms
- Modify CreateBentoGuided or create new form constructor
- Add size constraints with WithWidth/WithHeight
- Return completion message

### Step 5: Test thoroughly
- Manual testing of all key blocking scenarios
- Verify tab containment
- Test on different terminal sizes

---

## Code Example: Complete Pattern

```go
// browser.go
func (b Browser) Update(msg tea.Msg) (Browser, tea.Cmd) {
    // MODAL ACTIVE: Block all keys from browser
    if b.guidedModal != nil && b.guidedModal.active {
        switch msg := msg.(type) {
        case tea.KeyMsg:
            if msg.String() == "ctrl+c" {
                return b, tea.Quit
            }
            // All other keys → modal only
            var cmd tea.Cmd
            b.guidedModal, cmd = b.guidedModal.Update(msg)
            return b, cmd
        case GuidedCompleteMsg:
            b.guidedModal = nil
            return b.refreshList()
        default:
            var cmd tea.Cmd
            b.guidedModal, cmd = b.guidedModal.Update(msg)
            return b, cmd
        }
    }

    // Normal browser updates
    return b.handleNormalUpdate(msg)
}

func (b Browser) View() string {
    baseView := b.list.View()

    if b.guidedModal != nil && b.guidedModal.active {
        // Overlay modal on top
        modalView := b.guidedModal.View()
        return lipgloss.Place(
            b.width, b.height,
            lipgloss.Center, lipgloss.Center,
            modalView,
            lipgloss.WithWhitespaceChars("░"),
        )
    }

    return baseView
}
```

This pattern **guarantees** no key bubbling because the browser Update() never processes keys when modal is active.

---

## Conclusion

**Recommended**: Option 1 with conditional Update() routing

**Why**: Native Charm ecosystem, full control, proven Bubble Tea pattern

**Effort**: ~2 hours implementation

**Result**: Professional modal that blocks ALL keys except Ctrl+C and respects tab boundaries
