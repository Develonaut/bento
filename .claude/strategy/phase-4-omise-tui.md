# Phase 4: Omise (TUI)

**Status**: Pending
**Duration**: 4-5 hours
**Prerequisites**: Phase 3 complete, Karen approved

## Overview

Build the Bubble Tea TUI - "Omise" (お店 - shop). This is the interactive customer experience where users browse workflows, execute them with progress visualization, explore the pantry, and configure settings.

## Pre-Work Checklist

Before starting, you MUST:

1. ✅ Read [BENTO_BOX_PRINCIPLE.md](../BENTO_BOX_PRINCIPLE.md)
2. ✅ Read [CHARM_STACK_GUIDE.md](../ CHARM_STACK_GUIDE.md)
3. ✅ Confirm: "I understand the Bento Box Principle and will follow it"
4. ✅ Use TodoWrite to track all tasks
5. ✅ Phase 3 approved by Karen

## Goals

1. Create Bubble Tea application structure
2. Implement 5 main screens (browser, executor, pantry, settings, help)
3. Style with Lip Gloss
4. Integrate Bubbles components (list, spinner, progress, viewport)
5. Use Huh for forms and wizards
6. Keyboard navigation and shortcuts
7. Validate Bento Box compliance

## TUI Structure

```
pkg/omise/
├── go.mod
├── app.go              # Main Bubble Tea app
├── app_test.go
├── model.go            # App state model
├── update.go           # Update logic
├── view.go             # View rendering
├── screens/
│   ├── browser.go      # Workflow browser
│   ├── browser_test.go
│   ├── executor.go     # Execution viewer
│   ├── executor_test.go
│   ├── pantry.go       # Neta type explorer
│   ├── pantry_test.go
│   ├── settings.go     # Settings screen
│   └── help.go         # Help screen
├── components/
│   ├── header.go       # App header
│   ├── footer.go       # Keyboard shortcuts
│   ├── progress.go     # Custom progress bar
│   └── table.go        # Data tables
└── styles/
    └── theme.go        # Lip Gloss styles
```

## Deliverables

### 1. Main App Structure

**File**: `pkg/omise/app.go`
**File Size Target**: < 100 lines

```go
// Package omise provides the Bubble Tea TUI for Bento.
// Omise (お店) means "shop" - the customer interaction point.
package omise

import (
    "github.com/charmbracelet/bubbles/key"
    tea "github.com/charmbracelet/bubbletea"

    "bento/pkg/omise/screens"
)

// Launch starts the TUI application.
func Launch() error {
    p := tea.NewProgram(
        newModel(),
        tea.WithAltScreen(),
        tea.WithMouseCellMotion(),
    )

    _, err := p.Run()
    return err
}

// newModel creates the initial application model.
func newModel() Model {
    return Model{
        screen:   ScreenBrowser,
        browser:  screens.NewBrowser(),
        executor: screens.NewExecutor(),
        pantry:   screens.NewPantry(),
        settings: screens.NewSettings(),
        help:     screens.NewHelp(),
        keys:     DefaultKeyMap(),
    }
}

// DefaultKeyMap returns the default keyboard shortcuts.
func DefaultKeyMap() KeyMap {
    return KeyMap{
        Quit: key.NewBinding(
            key.WithKeys("q", "ctrl+c"),
            key.WithHelp("q", "quit"),
        ),
        Tab: key.NewBinding(
            key.WithKeys("tab"),
            key.WithHelp("tab", "next screen"),
        ),
        // ... more key bindings
    }
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: App launch and setup
- ✅ Composable screen structure
- ✅ File < 100 lines

### 2. Model

**File**: `pkg/omise/model.go`
**File Size Target**: < 150 lines

```go
package omise

import (
    tea "github.com/charmbracelet/bubbletea"

    "bento/pkg/omise/screens"
)

// Screen identifies which screen is active.
type Screen int

const (
    ScreenBrowser Screen = iota
    ScreenExecutor
    ScreenPantry
    ScreenSettings
    ScreenHelp
)

// Model is the root Bubble Tea model.
type Model struct {
    screen Screen
    width  int
    height int

    // Screens
    browser  screens.Browser
    executor screens.Executor
    pantry   screens.Pantry
    settings screens.Settings
    help     screens.Help

    // State
    keys       KeyMap
    quitting   bool
    lastError  error
}

// KeyMap defines keyboard shortcuts.
type KeyMap struct {
    Quit key.Binding
    Tab  key.Binding
    Back key.Binding
    Help key.Binding
}

// Init initializes the model.
func (m Model) Init() tea.Cmd {
    return m.currentScreen().Init()
}

// currentScreen returns the active screen model.
func (m Model) currentScreen() tea.Model {
    switch m.screen {
    case ScreenBrowser:
        return m.browser
    case ScreenExecutor:
        return m.executor
    case ScreenPantry:
        return m.pantry
    case ScreenSettings:
        return m.settings
    case ScreenHelp:
        return m.help
    default:
        return m.browser
    }
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: State management
- ✅ Clear screen separation
- ✅ File < 150 lines

### 3. Update Logic

**File**: `pkg/omise/update.go`
**File Size Target**: < 150 lines

```go
package omise

import (
    tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages and updates the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        return m.handleResize(msg)
    case tea.KeyMsg:
        return m.handleKey(msg)
    case ExecutionCompleteMsg:
        return m.handleExecutionComplete(msg)
    default:
        return m.updateScreen(msg)
    }
}

// handleResize updates dimensions.
func (m Model) handleResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
    m.width = msg.Width
    m.height = msg.Height
    return m, nil
}

// handleKey processes keyboard input.
func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch {
    case key.Matches(msg, m.keys.Quit):
        m.quitting = true
        return m, tea.Quit

    case key.Matches(msg, m.keys.Tab):
        m.screen = nextScreen(m.screen)
        return m, m.currentScreen().Init()

    default:
        return m.updateScreen(msg)
    }
}

// updateScreen delegates to the current screen.
func (m Model) updateScreen(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd

    switch m.screen {
    case ScreenBrowser:
        m.browser, cmd = m.browser.Update(msg)
    case ScreenExecutor:
        m.executor, cmd = m.executor.Update(msg)
    case ScreenPantry:
        m.pantry, cmd = m.pantry.Update(msg)
    case ScreenSettings:
        m.settings, cmd = m.settings.Update(msg)
    case ScreenHelp:
        m.help, cmd = m.help.Update(msg)
    }

    return m, cmd
}

// nextScreen cycles to the next screen.
func nextScreen(current Screen) Screen {
    return (current + 1) % 5
}

// ExecutionCompleteMsg signals workflow completion.
type ExecutionCompleteMsg struct {
    Result interface{}
    Error  error
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: Message handling
- ✅ Functions < 20 lines
- ✅ Clear delegation pattern
- ✅ File < 150 lines

### 4. View Rendering

**File**: `pkg/omise/view.go`
**File Size Target**: < 100 lines

```go
package omise

import (
    "github.com/charmbracelet/lipgloss"

    "bento/pkg/omise/components"
    "bento/pkg/omise/styles"
)

// View renders the TUI.
func (m Model) View() string {
    if m.quitting {
        return styles.Goodbye.Render("Thanks for using Bento! 🍱\n")
    }

    header := components.Header(m.screen, m.width)
    content := m.renderContent()
    footer := components.Footer(m.keys, m.width)

    return lipgloss.JoinVertical(
        lipgloss.Left,
        header,
        content,
        footer,
    )
}

// renderContent renders the active screen.
func (m Model) renderContent() string {
    contentHeight := m.height - 6 // Header + Footer

    switch m.screen {
    case ScreenBrowser:
        return m.browser.View()
    case ScreenExecutor:
        return m.executor.View()
    case ScreenPantry:
        return m.pantry.View()
    case ScreenSettings:
        return m.settings.View()
    case ScreenHelp:
        return m.help.View()
    default:
        return ""
    }
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: View composition
- ✅ Lip Gloss for layout
- ✅ File < 100 lines

### 5. Browser Screen

**File**: `pkg/omise/screens/browser.go`
**File Size Target**: < 200 lines

```go
// Package screens provides individual TUI screens.
package screens

import (
    "github.com/charmbracelet/bubbles/list"
    tea "github.com/charmbracelet/bubbletea"

    "bento/pkg/jubako" // Phase 5
)

// Browser shows available workflows.
type Browser struct {
    list     list.Model
    selected string
}

// NewBrowser creates a browser screen.
func NewBrowser() Browser {
    items := []list.Item{
        // Load from jubako (Phase 5)
    }

    l := list.New(items, list.NewDefaultDelegate(), 0, 0)
    l.Title = "🍱 Workflows"
    l.SetShowStatusBar(false)

    return Browser{list: l}
}

// Init initializes the browser.
func (b Browser) Init() tea.Cmd {
    return nil
}

// Update handles browser messages.
func (b Browser) Update(msg tea.Msg) (Browser, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "enter" {
            // Execute selected workflow
            if i, ok := b.list.SelectedItem().(workflowItem); ok {
                return b, executeWorkflow(i.path)
            }
        }
    }

    var cmd tea.Cmd
    b.list, cmd = b.list.Update(msg)
    return b, cmd
}

// View renders the browser.
func (b Browser) View() string {
    return b.list.View()
}

// workflowItem represents a .bento.yaml file.
type workflowItem struct {
    name string
    path string
}

func (i workflowItem) Title() string       { return i.name }
func (i workflowItem) Description() string { return i.path }
func (i workflowItem) FilterValue() string { return i.name }

// executeWorkflow starts workflow execution.
func executeWorkflow(path string) tea.Cmd {
    return func() tea.Msg {
        // Execute via itamae
        return ExecutionStartedMsg{Path: path}
    }
}

// ExecutionStartedMsg signals workflow start.
type ExecutionStartedMsg struct {
    Path string
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: Workflow browsing
- ✅ Bubbles list component
- ✅ Functions < 20 lines
- ✅ File < 200 lines

### 6. Executor Screen

**File**: `pkg/omise/screens/executor.go`
**File Size Target**: < 200 lines

```go
package screens

import (
    "github.com/charmbracelet/bubbles/progress"
    "github.com/charmbracelet/bubbles/spinner"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

// Executor shows workflow execution progress.
type Executor struct {
    spinner  spinner.Model
    progress progress.Model
    status   string
    running  bool
}

// NewExecutor creates an executor screen.
func NewExecutor() Executor {
    s := spinner.New()
    s.Spinner = spinner.Dot
    s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

    p := progress.New(progress.WithDefaultGradient())

    return Executor{
        spinner:  s,
        progress: p,
        status:   "Ready",
    }
}

// Init initializes the executor.
func (e Executor) Init() tea.Cmd {
    return nil
}

// Update handles executor messages.
func (e Executor) Update(msg tea.Msg) (Executor, tea.Cmd) {
    switch msg := msg.(type) {
    case ExecutionStartedMsg:
        e.running = true
        e.status = "Executing: " + msg.Path
        return e, tea.Batch(
            e.spinner.Tick,
            executeAsync(msg.Path),
        )

    case NodeProgressMsg:
        e.status = msg.NodeName
        e.progress.SetPercent(msg.Progress)
        return e, nil

    case ExecutionCompleteMsg:
        e.running = false
        if msg.Error != nil {
            e.status = "Failed: " + msg.Error.Error()
        } else {
            e.status = "Complete! ✅"
        }
        return e, nil
    }

    var cmd tea.Cmd
    if e.running {
        e.spinner, cmd = e.spinner.Update(msg)
    }
    return e, cmd
}

// View renders the executor.
func (e Executor) View() string {
    if !e.running {
        return e.status
    }

    return lipgloss.JoinVertical(
        lipgloss.Left,
        e.spinner.View()+" "+e.status,
        "",
        e.progress.View(),
    )
}

// NodeProgressMsg reports node execution progress.
type NodeProgressMsg struct {
    NodeName string
    Progress float64
}

// ExecutionCompleteMsg signals completion.
type ExecutionCompleteMsg struct {
    Result interface{}
    Error  error
}

// executeAsync runs workflow in background.
func executeAsync(path string) tea.Cmd {
    return func() tea.Msg {
        // Execute via itamae
        // Send NodeProgressMsg updates
        // Return ExecutionCompleteMsg when done
        return ExecutionCompleteMsg{}
    }
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: Execution visualization
- ✅ Bubbles spinner and progress
- ✅ Functions < 20 lines
- ✅ File < 200 lines

### 7. Pantry Screen

**File**: `pkg/omise/screens/pantry.go`
**File Size Target**: < 150 lines

```go
package screens

import (
    "github.com/charmbracelet/bubbles/table"
    tea "github.com/charmbracelet/bubbletea"
)

// Pantry shows available neta types.
type Pantry struct {
    table table.Model
}

// NewPantry creates a pantry screen.
func NewPantry() Pantry {
    columns := []table.Column{
        {Title: "Type", Width: 30},
        {Title: "Category", Width: 20},
        {Title: "Description", Width: 50},
    }

    rows := []table.Row{
        {"http", "Network", "HTTP request execution"},
        {"transform.jq", "Data", "JQ transformation"},
        {"conditional.if", "Control", "If/else logic"},
        {"loop.for", "Control", "For loop iteration"},
        {"group.sequence", "Group", "Sequential execution"},
    }

    t := table.New(
        table.WithColumns(columns),
        table.WithRows(rows),
        table.WithFocused(true),
        table.WithHeight(10),
    )

    return Pantry{table: t}
}

// Init initializes the pantry.
func (p Pantry) Init() tea.Cmd {
    return nil
}

// Update handles pantry messages.
func (p Pantry) Update(msg tea.Msg) (Pantry, tea.Cmd) {
    var cmd tea.Cmd
    p.table, cmd = p.table.Update(msg)
    return p, cmd
}

// View renders the pantry.
func (p Pantry) View() string {
    return p.table.View()
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: Type catalog
- ✅ Bubbles table component
- ✅ File < 150 lines

### 8. Styles

**File**: `pkg/omise/styles/theme.go`
**File Size Target**: < 150 lines

```go
// Package styles provides Lip Gloss styles for the TUI.
package styles

import "github.com/charmbracelet/lipgloss"

var (
    // Primary colors
    Primary   = lipgloss.Color("205") // Pink for bento theme
    Secondary = lipgloss.Color("99")  // Purple
    Success   = lipgloss.Color("42")  // Green
    Error     = lipgloss.Color("196") // Red
    Muted     = lipgloss.Color("241") // Gray

    // Title style
    Title = lipgloss.NewStyle().
        Bold(true).
        Foreground(Primary).
        MarginBottom(1)

    // Header style
    Header = lipgloss.NewStyle().
        Bold(true).
        Padding(0, 1).
        Background(Primary).
        Foreground(lipgloss.Color("230"))

    // Footer style
    Footer = lipgloss.NewStyle().
        Foreground(Muted).
        MarginTop(1)

    // Error style
    ErrorStyle = lipgloss.NewStyle().
        Foreground(Error).
        Bold(true)

    // Success style
    SuccessStyle = lipgloss.NewStyle().
        Foreground(Success).
        Bold(true)

    // Goodbye message
    Goodbye = lipgloss.NewStyle().
        Bold(true).
        Foreground(Primary).
        Padding(1)
)
```

**Bento Box Compliance**:
- ✅ Single responsibility: Style definitions
- ✅ Centralized theme
- ✅ File < 150 lines

### 9. Components

**File**: `pkg/omise/components/header.go`

```go
// Package components provides reusable TUI components.
package components

import (
    "github.com/charmbracelet/lipgloss"
    "bento/pkg/omise/styles"
)

// Header renders the app header.
func Header(screen int, width int) string {
    title := "🍱 Bento"
    screenName := getScreenName(screen)

    header := lipgloss.JoinHorizontal(
        lipgloss.Left,
        title,
        " | ",
        screenName,
    )

    return styles.Header.Width(width).Render(header)
}

// getScreenName returns the screen name.
func getScreenName(screen int) string {
    names := []string{"Browser", "Executor", "Pantry", "Settings", "Help"}
    if screen < 0 || screen >= len(names) {
        return "Unknown"
    }
    return names[screen]
}
```

## Integration with Root Command

Update `cmd/bento/root.go`:

```go
Run: func(cmd *cobra.Command, args []string) {
    if err := omise.Launch(); err != nil {
        fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
        os.Exit(1)
    }
},
```

## Testing Strategy

TUI testing with bubbletea/tea testing utilities:

```go
func TestBrowser_Navigation(t *testing.T) {
    m := NewBrowser()

    // Simulate key presses
    m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
    m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})

    // Assert state changes
}
```

## Validation Commands

Before marking phase complete:

```bash
# Format
go fmt ./...

# Lint
golangci-lint run

# Test
go test -v ./pkg/omise/...

# Build and launch TUI
go build ./cmd/bento
./bento  # Should launch TUI!

# File size check
find pkg/omise -name "*.go" -exec wc -l {} + | sort -rn | head -10
```

## Success Criteria

Phase 4 is complete when:

1. ✅ Bubble Tea app structure complete
2. ✅ 5 screens implemented and working
3. ✅ Lip Gloss styling applied
4. ✅ Bubbles components integrated
5. ✅ Keyboard navigation functional
6. ✅ TUI launches from `bento` command
7. ✅ All files < 250 lines
8. ✅ All functions < 20 lines
9. ✅ Tests passing
10. ✅ **Karen's approval granted**

## Common Pitfalls to Avoid

1. ❌ **God model** - Keep screen models separate
2. ❌ **Mixing view logic** - Components should be pure rendering
3. ❌ **Large update functions** - Delegate to screen-specific handlers
4. ❌ **Blocking operations** - Use tea.Cmd for async work
5. ❌ **No keyboard shortcuts** - Every action needs a shortcut

## Next Phase

After Karen approval, proceed to **[Phase 5: Jubako Storage](./phase-5-jubako.md)** to:
- Implement .bento.yaml file management
- Add workflow history
- Create import/export functionality

## Execution Prompt

```
I'm ready to begin Phase 4: Omise TUI.

I have read the Bento Box Principle and Charm Stack Guide and will follow them.

Please build the Bubble Tea TUI:
- App structure (model, update, view)
- 5 screens (browser, executor, pantry, settings, help)
- Lip Gloss styling
- Bubbles components

Each screen in focused file < 200 lines. I will use TodoWrite to track progress and get Karen's approval before completing.
```

---

**Phase 4 Omise TUI**: Interactive shop experience 🍱
