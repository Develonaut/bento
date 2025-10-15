# Charm Stack Component Guide

**Understanding the Charm Ecosystem for TUI Development**

---

## Quick Overview

The Charm stack has 4 main components that work together:

```
┌─────────────────────────────────────────────┐
│         Your TUI Application                │
├─────────────────────────────────────────────┤
│  Huh (Forms) │ Bubbles (Components)         │  ← Pre-built UI components
├─────────────────────────────────────────────┤
│  Lip Gloss (Styling/Layout)                 │  ← Make it pretty
├─────────────────────────────────────────────┤
│  Bubble Tea (Framework)                     │  ← Core TUI engine
├─────────────────────────────────────────────┤
│  Terminal                                   │
└─────────────────────────────────────────────┘
```

---

## 1. Bubble Tea (The Framework) 🫧

**What it does:** The core TUI framework - handles terminal events, screen rendering, and app state.

**Think of it as:** React for the terminal (but uses Elm Architecture)

**You always need this** - It's the foundation.

```go
package main

import (
    "fmt"
    tea "github.com/charmbracelet/bubbletea"
)

// Your app is a Model
type model struct {
    choices  []string
    cursor   int
    selected map[int]struct{}
}

// Initialize the app
func initialModel() model {
    return model{
        choices:  []string{"Run Flow", "Create Flow", "Settings", "Quit"},
        selected: make(map[int]struct{}),
    }
}

// Handle events (keyboard, mouse, etc.)
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "up", "k":
            if m.cursor > 0 {
                m.cursor--
            }
        case "down", "j":
            if m.cursor < len(m.choices)-1 {
                m.cursor++
            }
        case "enter":
            m.selected[m.cursor] = struct{}{}
        }
    }
    return m, nil
}

// Render the UI
func (m model) View() string {
    s := "What should we do?\n\n"
    for i, choice := range m.choices {
        cursor := " "
        if m.cursor == i {
            cursor = ">"
        }
        s += fmt.Sprintf("%s %s\n", cursor, choice)
    }
    return s
}

func main() {
    p := tea.NewProgram(initialModel())
    p.Run()
}
```

**Key Concepts:**
- **Model** - Your app's state (like React state)
- **Update** - Handles events and updates state (like React's setState)
- **View** - Renders UI from state (like React's render)
- **Cmd** - Side effects (like fetching data)

---

## 2. Lip Gloss (The Styler) 💅

**What it does:** Makes your terminal UI beautiful with colors, borders, padding, alignment, etc.

**Think of it as:** CSS for the terminal (but better!)

**Use when:** You want to style text, create layouts, add colors/borders.

```go
package main

import (
    "fmt"
    "github.com/charmbracelet/lipgloss"
)

var (
    // Define a style (like CSS class)
    titleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#FF79C6")).  // Hot pink
        Background(lipgloss.Color("#282A36")).  // Dark background
        Padding(1, 2).                          // Padding: top/bottom, left/right
        MarginTop(1)

    // Another style
    selectedStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#7CD5FA")).  // Cyan
        Bold(true)

    normalStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#F8F8F2"))   // White
)

func main() {
    // Apply styles to text
    title := titleStyle.Render("🚀 Atomiton")
    selected := selectedStyle.Render("> Run Flow")
    normal := normalStyle.Render("  Create Flow")

    fmt.Println(title)
    fmt.Println(selected)
    fmt.Println(normal)
}
```

**Advanced Layout:**

```go
// Create a box with border
boxStyle := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("#874BFD")).
    Padding(1, 2).
    Width(40)

content := boxStyle.Render("This is inside a pretty box!")

// Join things horizontally
left := lipgloss.NewStyle().Width(20).Render("Left column")
right := lipgloss.NewStyle().Width(20).Render("Right column")
row := lipgloss.JoinHorizontal(lipgloss.Top, left, right)

// Join things vertically
header := titleStyle.Render("Header")
body := normalStyle.Render("Body content")
footer := lipgloss.NewStyle().Faint(true).Render("Footer")
page := lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
```

**Lip Gloss does:**
- ✅ Colors (foreground, background)
- ✅ Typography (bold, italic, underline)
- ✅ Spacing (padding, margin)
- ✅ Borders (rounded, thick, double, custom)
- ✅ Alignment (left, center, right)
- ✅ Width/height constraints
- ✅ Layout (join horizontal/vertical)

---

## 3. Bubbles (Pre-built Components) 🎨

**What it does:** Ready-to-use UI components (list, spinner, text input, progress bar, etc.)

**Think of it as:** Component library (like Material-UI or Ant Design)

**Use when:** You don't want to build common UI components from scratch.

```go
package main

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/bubbles/list"
    "github.com/charmbracelet/bubbles/spinner"
    "github.com/charmbracelet/bubbles/textinput"
    "github.com/charmbracelet/bubbles/progress"
)

type model struct {
    list      list.Model      // List component
    spinner   spinner.Model   // Loading spinner
    input     textinput.Model // Text input field
    progress  progress.Model  // Progress bar
}

func initialModel() model {
    // Create a list
    items := []list.Item{
        item{title: "my-flow.yaml", desc: "HTTP request pipeline"},
        item{title: "image-processor.yaml", desc: "Batch image processing"},
    }
    l := list.New(items, list.NewDefaultDelegate(), 0, 0)
    l.Title = "Available Flows"

    // Create a spinner
    s := spinner.New()
    s.Spinner = spinner.Dot
    s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

    // Create a text input
    ti := textinput.New()
    ti.Placeholder = "Enter flow name..."
    ti.Focus()

    // Create a progress bar
    p := progress.New(progress.WithDefaultGradient())

    return model{
        list:     l,
        spinner:  s,
        input:    ti,
        progress: p,
    }
}
```

**Available Bubbles Components:**

| Component | Purpose | Example Use |
|-----------|---------|-------------|
| `list` | Scrollable list with items | Flow browser, menu |
| `textinput` | Single-line text input | Search, file name |
| `textarea` | Multi-line text input | Script editor |
| `spinner` | Loading indicator | During execution |
| `progress` | Progress bar | File upload, processing |
| `table` | Data table | Display results |
| `paginator` | Page navigation | Large lists |
| `viewport` | Scrollable content | Log viewer |
| `filepicker` | File browser | Select flow file |
| `stopwatch` | Timer/stopwatch | Execution time |

**Each component:**
- Has its own `Model` (state)
- Has `Update(msg)` method (handles events)
- Has `View()` method (renders UI)
- Integrates seamlessly with Bubble Tea

---

## 4. Huh (Forms) 📝

**What it does:** High-level form builder for collecting user input with validation.

**Think of it as:** Form library (like Formik or React Hook Form)

**Use when:** You need to collect multiple fields of data (create flow wizard, settings form, etc.)

```go
package main

import (
    "fmt"
    "github.com/charmbracelet/huh"
)

func main() {
    var (
        flowName    string
        description string
        nodeType    string
        concurrent  bool
    )

    // Define a form
    form := huh.NewForm(
        huh.NewGroup(
            // Text input
            huh.NewInput().
                Title("Flow Name").
                Placeholder("my-awesome-flow").
                Value(&flowName).
                Validate(func(s string) error {
                    if len(s) < 3 {
                        return fmt.Errorf("name must be at least 3 characters")
                    }
                    return nil
                }),

            // Text area
            huh.NewText().
                Title("Description").
                Placeholder("What does this flow do?").
                Value(&description),

            // Select dropdown
            huh.NewSelect[string]().
                Title("Node Type").
                Options(
                    huh.NewOption("HTTP Request", "httpRequest"),
                    huh.NewOption("Transform", "transform"),
                    huh.NewOption("File System", "fileSystem"),
                ).
                Value(&nodeType),

            // Checkbox
            huh.NewConfirm().
                Title("Run nodes concurrently?").
                Value(&concurrent),
        ),
    )

    // Run the form
    err := form.Run()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    // Use the collected data
    fmt.Printf("Creating flow: %s\n", flowName)
    fmt.Printf("Description: %s\n", description)
    fmt.Printf("Node type: %s\n", nodeType)
    fmt.Printf("Concurrent: %v\n", concurrent)
}
```

**Huh Field Types:**

| Field Type | Use Case | Example |
|------------|----------|---------|
| `Input` | Single-line text | Flow name, URL |
| `Text` | Multi-line text | Description, notes |
| `Select` | Choose one option | Node type, theme |
| `MultiSelect` | Choose multiple | Tags, categories |
| `Confirm` | Yes/No question | Enable feature? |
| `FilePicker` | Select file/directory | Choose flow file |

**Features:**
- Built-in validation
- Keyboard navigation (Tab, Enter, Esc)
- Accessible design
- Themeable with Lip Gloss
- **Automatic form flow** - You don't need to wire up state!

---

## How They Work Together

### Example: Atomiton Flow Creator

```go
package main

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/bubbles/list"
    "github.com/charmbracelet/lipgloss"
    "github.com/charmbracelet/huh"
)

// Bubble Tea: The framework (required)
type model struct {
    screen   string      // Which screen to show
    flowList list.Model  // Bubbles: Pre-built list component
    form     *huh.Form   // Huh: Form for creating flows
}

// Bubble Tea: Initialize
func initialModel() model {
    items := []list.Item{
        // ... flow items
    }
    l := list.New(items, list.NewDefaultDelegate(), 0, 0)

    return model{
        screen:   "menu",
        flowList: l,
    }
}

// Bubble Tea: Handle events
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "c" {
            // Switch to create screen with Huh form
            var flowName string
            m.form = huh.NewForm(
                huh.NewGroup(
                    huh.NewInput().
                        Title("Flow Name").
                        Value(&flowName),
                ),
            )
            m.screen = "create"
        }
    }

    // Update the active component
    if m.screen == "list" {
        var cmd tea.Cmd
        m.flowList, cmd = m.flowList.Update(msg)
        return m, cmd
    }

    return m, nil
}

// Bubble Tea: Render UI
func (m model) View() string {
    // Lip Gloss: Style the UI
    titleStyle := lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#FF79C6"))

    title := titleStyle.Render("🚀 Atomiton")

    // Render different screens
    var content string
    switch m.screen {
    case "menu":
        content = "1. Run Flow\n2. Create Flow\n3. Settings"
    case "list":
        // Bubbles: Use pre-built list
        content = m.flowList.View()
    case "create":
        // Huh: Use form
        content = m.form.View()
    }

    // Lip Gloss: Layout
    return lipgloss.JoinVertical(
        lipgloss.Left,
        title,
        content,
    )
}

func main() {
    p := tea.NewProgram(initialModel())
    p.Run()
}
```

---

## Recommended Stack for Atomiton

### Core (Required)
1. **Bubble Tea** - TUI framework (always needed)
2. **Lip Gloss** - Styling and layout (always needed)

### Components (Choose based on needs)
3. **Bubbles** - Use for common components:
   - `list` → Flow browser
   - `textinput` → Search bar
   - `progress` → Execution progress
   - `spinner` → Loading states
   - `viewport` → Log viewer

4. **Huh** - Use for forms:
   - Create flow wizard
   - Settings configuration
   - Multi-step forms

### Architecture

```
Your Atomiton TUI
├── Screens (Bubble Tea Models)
│   ├── MenuScreen
│   ├── FlowListScreen (uses Bubbles list)
│   ├── CreateFlowScreen (uses Huh form)
│   └── SettingsScreen (uses Huh form)
├── Components (if needed)
│   ├── Header (Lip Gloss)
│   ├── Footer (Lip Gloss)
│   └── StatusBar (Lip Gloss)
└── Theme (Lip Gloss styles)
```

---

## Comparison to Your Current Stack

### Current (TypeScript + Ink)

```typescript
// Ink (React-based)
<Box flexDirection="column">
  <Text bold color="magenta">Atomiton</Text>
  <SelectInput items={flows} onSelect={handleSelect} />
</Box>
```

### Future (Go + Charm)

```go
// Bubble Tea + Lip Gloss
titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("magenta"))
title := titleStyle.Render("Atomiton")

// Use Bubbles list component
flowList := list.New(flows, list.NewDefaultDelegate(), 0, 0)
```

**Advantages:**
- ✅ **Better performance** - Native Go, no JS runtime
- ✅ **More features** - Bubbles has more components than Ink
- ✅ **Better styling** - Lip Gloss > Chalk/Ink styles
- ✅ **Form handling** - Huh makes forms trivial
- ✅ **Single binary** - No npm/node_modules

---

## Learning Path

### Week 1: Bubble Tea Basics
- Build simple menu app
- Learn Model-Update-View pattern
- Handle keyboard events
- Understand commands (tea.Cmd)

### Week 2: Lip Gloss Styling
- Add colors and styles
- Create layouts (borders, padding)
- Build reusable style themes
- Master JoinHorizontal/JoinVertical

### Week 3: Bubbles Components
- Integrate list component (flow browser)
- Add text input (search)
- Use spinner (loading states)
- Add progress bar (execution)

### Week 4: Huh Forms
- Build create flow wizard
- Add validation
- Multi-step forms
- Settings screen

### Week 5: Integration
- Connect to Atomiton backend
- Real flow execution
- Error handling
- Polish UX

---

## Official Resources

- **Bubble Tea**: https://github.com/charmbracelet/bubbletea
- **Lip Gloss**: https://github.com/charmbracelet/lipgloss
- **Bubbles**: https://github.com/charmbracelet/bubbles
- **Huh**: https://github.com/charmbracelet/huh
- **Examples**: https://github.com/charmbracelet/bubbletea/tree/master/examples
- **Tutorials**: https://charm.sh/blog/

---

## Summary

**You need:**
- ✅ **Bubble Tea** (framework) - Always
- ✅ **Lip Gloss** (styling) - Always
- ✅ **Bubbles** (components) - For common UI elements
- ✅ **Huh** (forms) - For data collection

**Think of it like web development:**
- Bubble Tea = React (framework)
- Lip Gloss = CSS (styling)
- Bubbles = Component library (Material-UI)
- Huh = Form library (Formik)

**Start simple:** Bubble Tea + Lip Gloss only, add Bubbles/Huh as needed.
