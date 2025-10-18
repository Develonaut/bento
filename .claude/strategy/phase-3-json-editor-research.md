# Phase 3: JSON Bento Editor - Research Report

**Date**: 2025-10-17
**Status**: Research Complete
**Next Phase**: Implementation Planning

---

## Executive Summary

### Key Findings

1. **Current Format**: Bento currently uses YAML exclusively (`*.bento.yaml`) with nested node structures in the `parameters.nodes` field for groups
2. **Recommended Editor**: Use `bubbles/textarea` with `chroma` for syntax highlighting - proven, simple, and already in ecosystem
3. **Migration Complexity**: LOW - JSON is drop-in compatible with existing struct tags, minimal code changes needed
4. **Implementation Time**: 8-12 hours for MVP (editor + migration)
5. **Risk Level**: LOW - Clean switch to JSON, structs already have json tags

### Critical Recommendations

- ✅ **JSON over YAML**: Simpler syntax, better for manual editing, stdlib support
- ✅ **Flat array structure**: Use `parentId` field instead of nested objects (already partially supported in structs)
- ✅ **textarea + chroma**: Don't reinvent the wheel, use proven components
- ✅ **Debounced validation**: Validate on pause (500ms), not every keystroke
- ✅ **Hotkey templates**: Ctrl+N for new node templates, not auto-insertion
- ✅ **Line numbers**: textarea supports `ShowLineNumbers` - use it
- ✅ **Clean migration**: Replace YAML parser with JSON parser, convert all examples in one go

---

## 1. Current State Analysis

### File Format: YAML

**Files Found:**
- `examples/http-get.bento.yaml`
- `examples/group-sequence.bento.yaml`
- `examples/loop-for.bento.yaml`
- `examples/conditional-if.bento.yaml`
- 8+ total example files

**Current Structure (YAML):**
```yaml
version: "1.0"
type: group.sequence
name: API Pipeline
parameters:
    nodes:
        - name: Fetch User
          type: http
          parameters:
            method: GET
            url: https://api.github.com/users/octocat
        - name: Extract Login
          type: transform.jq
          parameters:
            query: .login
```

### Parser: `pkg/jubako/parser.go`

**Key Functions:**
- `Parse(path string)` - reads .bento.yaml files
- `ParseBytes(data []byte)` - unmarshals YAML
- `Format(def neta.Definition)` - marshals to YAML
- `normalizeDefinition()` - extracts child nodes from `parameters.nodes`
- `assignNodeIDs()` - auto-generates IDs for nodes without them
- `autoGenerateEdges()` - generates sequential edges for group types

**YAML-Specific Behavior:**
```go
import "gopkg.in/yaml.v3"

func (p *Parser) ParseBytes(data []byte) (neta.Definition, error) {
    var def neta.Definition
    if err := yaml.Unmarshal(data, &def); err != nil {
        return neta.Definition{}, fmt.Errorf("invalid YAML: %w", err)
    }
    // ... normalization logic
}
```

### Data Model: `pkg/neta/definition.go`

**Struct Definition:**
```go
type Definition struct {
    ID          string                 `yaml:"id,omitempty" json:"id,omitempty"`
    ParentID    string                 `yaml:"parentId,omitempty" json:"parentId,omitempty"` // ✅ Already supports flat structure!
    Version     string                 `yaml:"version" json:"version"`
    Type        string                 `yaml:"type" json:"type"`
    Name        string                 `yaml:"name" json:"name"`
    Icon        string                 `yaml:"icon,omitempty" json:"icon,omitempty"`
    Description string                 `yaml:"description,omitempty" json:"description,omitempty"`
    Parameters  map[string]interface{} `yaml:"parameters,omitempty" json:"parameters,omitempty"`
    Nodes       []Definition           `yaml:"nodes,omitempty" json:"nodes,omitempty"`
    Edges       []NodeEdge             `yaml:"edges,omitempty" json:"edges,omitempty"`
}
```

**Observation**: Struct already has dual `yaml` and `json` tags! Migration will be trivial.

### Validation: `pkg/neta/validator.go` + `pkg/neta/schemas/`

**Node Types with Schemas:**
1. `http` - GET/POST/PUT/DELETE with URL, headers, body
2. `transform.jq` - jq query transformation
3. `group.sequence` - sequential execution
4. `group.parallel` - parallel execution (with optional `max_concurrent`)
5. `loop.for` - iterate over items array
6. `conditional.if` - boolean condition with then/else branches
7. `file.write` - write to file path

**Validation Framework:**
- Each schema implements `Schema` interface with `Validate()` and `Fields()`
- `Fields()` returns metadata for UI form generation (Huh integration)
- Validation errors use `ValidationErrors` type with field-level detail

### Keymap System: `pkg/omise/components/`

**Existing Keymaps:**
- `keymap_global.go` - App-wide keys (quit, help)
- `keymap_navigation.go` - Up/down/tab navigation
- `keymap_browser.go` - Browser-specific (r=run, e=edit, c=copy, d=delete, n=new)
- `keymap_settings.go` - Settings screen keys
- `keymap_picker.go` - Directory picker keys

**Pattern:**
```go
type BrowserKeyMap struct {
    Navigation  NavigationKeyMap
    Execute     key.Binding
    Run         key.Binding
    New         key.Binding
    // ...
}

func NewBrowserKeyMap() BrowserKeyMap {
    return BrowserKeyMap{
        New: key.NewBinding(
            key.WithKeys("n"),
            key.WithHelp("n", "new bento"),
        ),
    }
}
```

**Recommendation**: Create `keymap_editor.go` following this pattern

---

## 2. JSON Editor Components

### Primary Component: `bubbles/textarea`

**Already in go.mod:**
```
github.com/charmbracelet/bubbles v0.21.1-0.20250623103423-23b8fd6302d7
```

**Capabilities:**
- ✅ Multi-line text editing
- ✅ Line numbers (`ShowLineNumbers` option)
- ✅ Cursor movement (up/down/left/right, home/end, page up/down)
- ✅ Text selection
- ✅ Clipboard paste support
- ✅ Vertical scrolling
- ✅ Width/height constraints
- ✅ Focused/blurred styling
- ✅ Unicode support
- ✅ Character limit

**Basic Usage:**
```go
import "github.com/charmbracelet/bubbles/textarea"

type EditorModel struct {
    textarea textarea.Model
}

func NewEditor() EditorModel {
    ta := textarea.New()
    ta.Placeholder = "Enter JSON here..."
    ta.ShowLineNumbers = true
    ta.Focus()

    return EditorModel{textarea: ta}
}

func (m EditorModel) Update(msg tea.Msg) (EditorModel, tea.Cmd) {
    var cmd tea.Cmd
    m.textarea, cmd = m.textarea.Update(msg)
    return m, cmd
}

func (m EditorModel) View() string {
    return m.textarea.View()
}
```

**Pros:**
- Part of Charm ecosystem (consistent with existing code)
- Battle-tested in production apps
- Line numbers built-in
- Clipboard integration
- Minimal code required

**Cons:**
- No syntax highlighting built-in
- No code-specific features (auto-indent, bracket matching)
- Large files may have performance issues

**Recommendation**: ✅ Use `bubbles/textarea` as the foundation

---

## 3. Syntax Highlighting

### Option 1: Chroma (Recommended)

**Library**: `github.com/alecthomas/chroma`

**Features:**
- ✅ Pure Go syntax highlighter (Pygments port)
- ✅ JSON lexer built-in
- ✅ Terminal output formatters (8-color, 256-color, true-color)
- ✅ Multiple color schemes
- ✅ Fast and lightweight

**Terminal Formatters:**
- `TTY` - 8-color (basic ANSI)
- `TTY256` - 256-color palette
- `TTY16m` - True color (16 million colors)

**Example Code:**
```go
import (
    "github.com/alecthomas/chroma"
    "github.com/alecthomas/chroma/formatters"
    "github.com/alecthomas/chroma/lexers"
    "github.com/alecthomas/chroma/styles"
)

func HighlightJSON(jsonText string) (string, error) {
    // Get JSON lexer
    lexer := lexers.Get("json")
    if lexer == nil {
        lexer = lexers.Fallback
    }

    // Choose style (monokai, dracula, etc.)
    style := styles.Get("monokai")
    if style == nil {
        style = styles.Fallback
    }

    // Format for terminal (true color)
    formatter := formatters.Get("terminal16m")
    if formatter == nil {
        formatter = formatters.Fallback
    }

    // Tokenize and format
    iterator, err := lexer.Tokenise(nil, jsonText)
    if err != nil {
        return "", err
    }

    var buf bytes.Buffer
    err = formatter.Format(&buf, style, iterator)
    if err != nil {
        return "", err
    }

    return buf.String(), nil
}
```

**Integration with textarea:**
```go
func (m EditorModel) View() string {
    content := m.textarea.Value()

    // Apply syntax highlighting
    highlighted, err := HighlightJSON(content)
    if err != nil {
        // Fallback to plain text
        return m.textarea.View()
    }

    // Replace textarea content with highlighted version
    // Note: This is conceptual - actual implementation needs care
    // to preserve cursor position and editing state
    return highlighted
}
```

**Challenge**: textarea doesn't natively support styled content rendering - it's plain text. We have two approaches:

**Approach A: Post-Render Highlighting (Recommended)**
- Render textarea normally
- Apply syntax highlighting to the *displayed* content only
- Keep editing in plain text
- Refresh highlights on change

**Approach B: Background Preview**
- Split view: plain textarea (left) + highlighted preview (right)
- User edits plain text, sees highlighted version alongside
- Simpler implementation, no cursor issues

**Recommendation**: ✅ Use **Approach B** (split view) for MVP, investigate Approach A for future enhancement

### Option 2: Glamour

**Library**: `github.com/charmbracelet/glamour`

**Features:**
- ✅ Markdown renderer for terminals
- ✅ Part of Charm ecosystem
- ❌ NOT suitable for JSON (markdown only)

**Recommendation**: ❌ Do not use - designed for markdown, not JSON

### Option 3: Custom Highlighting

**Approach**: Manually parse JSON and apply lipgloss styles

**Pros:**
- Full control over colors
- No external dependency

**Cons:**
- Reinventing the wheel
- Complex to implement correctly
- Fragile (must handle all JSON edge cases)
- More code to maintain

**Recommendation**: ❌ Do not implement - use Chroma instead

### Final Recommendation: Syntax Highlighting

✅ **Use Chroma with split-view approach**:
- Left pane: `bubbles/textarea` (plain text editing)
- Right pane: Chroma-highlighted preview (read-only)
- Toggle view with hotkey (Tab key or F2)
- For MVP, optional - can ship without highlighting first

---

## 4. Template/Hotkey Insertion System

### Design: Template Library + Hotkey Bindings

**Templates Needed:**

1. **HTTP GET**
2. **HTTP POST**
3. **Shell Command**
4. **Transform (jq)**
5. **Group (Sequence)**
6. **Group (Parallel)**
7. **Loop (For)**
8. **Conditional (If)**
9. **File Write**

**Template Structure:**

```go
package templates

// Template represents a JSON node template
type Template struct {
    Name        string
    Description string
    Hotkey      string
    JSON        string
}

var NodeTemplates = []Template{
    {
        Name:        "HTTP GET",
        Description: "HTTP GET request",
        Hotkey:      "ctrl+h",
        JSON: `{
  "id": "node-$ID",
  "type": "http",
  "name": "$NAME",
  "parameters": {
    "method": "GET",
    "url": "$URL"
  }
}`,
    },
    {
        Name:        "HTTP POST",
        Description: "HTTP POST request with JSON body",
        Hotkey:      "ctrl+p",
        JSON: `{
  "id": "node-$ID",
  "type": "http",
  "name": "$NAME",
  "parameters": {
    "method": "POST",
    "url": "$URL",
    "headers": {
      "Content-Type": "application/json"
    },
    "body": "$BODY"
  }
}`,
    },
    {
        Name:        "Transform (jq)",
        Description: "jq data transformation",
        Hotkey:      "ctrl+j",
        JSON: `{
  "id": "node-$ID",
  "type": "transform.jq",
  "name": "$NAME",
  "parameters": {
    "query": "$QUERY"
  }
}`,
    },
    {
        Name:        "Sequence Group",
        Description: "Sequential execution group",
        Hotkey:      "ctrl+s",
        JSON: `{
  "id": "node-$ID",
  "type": "group.sequence",
  "name": "$NAME",
  "nodes": []
}`,
    },
    {
        Name:        "For Loop",
        Description: "Iterate over array items",
        Hotkey:      "ctrl+l",
        JSON: `{
  "id": "node-$ID",
  "type": "loop.for",
  "name": "$NAME",
  "parameters": {
    "items": [],
    "body": {
      "type": "$BODY_TYPE",
      "name": "$BODY_NAME",
      "parameters": {}
    }
  }
}`,
    },
}
```

### Insertion Logic

**On Hotkey Press:**

1. Get cursor position in textarea
2. Get current line content
3. Calculate indentation level
4. Insert template at cursor with proper indentation
5. Replace placeholders (`$ID`, `$NAME`, etc.) with defaults or prompts
6. Move cursor to first placeholder for editing

**Implementation:**

```go
type EditorModel struct {
    textarea textarea.Model
    nextID   int
}

func (m *EditorModel) InsertTemplate(tmpl Template) {
    // Generate unique ID
    nodeID := fmt.Sprintf("node-%d", m.nextID)
    m.nextID++

    // Replace placeholders
    content := strings.ReplaceAll(tmpl.JSON, "$ID", nodeID)
    content = strings.ReplaceAll(content, "$NAME", "New Node")
    content = strings.ReplaceAll(content, "$URL", "https://example.com")
    // ... other replacements

    // Get cursor position
    cursorPos := m.textarea.CursorPosition()

    // Insert at cursor
    current := m.textarea.Value()
    before := current[:cursorPos]
    after := current[cursorPos:]

    newContent := before + content + after
    m.textarea.SetValue(newContent)

    // Move cursor to end of inserted content
    m.textarea.SetCursor(cursorPos + len(content))
}

func (m EditorModel) Update(msg tea.Msg) (EditorModel, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Check for template hotkeys
        if msg.Type == tea.KeyCtrlH {
            m.InsertTemplate(templates.NodeTemplates[0]) // HTTP GET
            return m, nil
        }
        if msg.Type == tea.KeyCtrlP {
            m.InsertTemplate(templates.NodeTemplates[1]) // HTTP POST
            return m, nil
        }
        // ... other hotkeys
    }

    // Default textarea handling
    var cmd tea.Cmd
    m.textarea, cmd = m.textarea.Update(msg)
    return m, cmd
}
```

### Smart Indentation

**Challenge**: Insert template with correct indentation for current context

**Solution**: Detect indentation level from current line

```go
func detectIndentation(line string) int {
    spaces := 0
    for _, ch := range line {
        if ch == ' ' {
            spaces++
        } else if ch == '\t' {
            spaces += 2 // Assume 2 spaces per tab
        } else {
            break
        }
    }
    return spaces
}

func indentTemplate(template string, level int) string {
    indent := strings.Repeat(" ", level)
    lines := strings.Split(template, "\n")

    var indented []string
    for _, line := range lines {
        if line != "" {
            indented = append(indented, indent+line)
        } else {
            indented = append(indented, line)
        }
    }

    return strings.Join(indented, "\n")
}
```

### Cursor Positioning

**Goal**: After insertion, place cursor at first placeholder for immediate editing

**Approach**: Use special marker like `$|$` to mark cursor position in template

```go
JSON: `{
  "id": "node-$ID",
  "name": "$|$",$  // Cursor goes here
  "type": "http",
  "parameters": {
    "url": ""
  }
}`

func (m *EditorModel) InsertTemplate(tmpl Template) {
    content := processPlaceholders(tmpl.JSON)

    // Find cursor marker
    cursorMarker := "$|$"
    markerPos := strings.Index(content, cursorMarker)

    if markerPos >= 0 {
        // Remove marker
        content = strings.ReplaceAll(content, cursorMarker, "")

        // Insert content
        // ... insertion logic

        // Set cursor to marker position
        m.textarea.SetCursor(insertPos + markerPos)
    }
}
```

### Alternative: Template Wizard

**Instead of direct insertion, show a quick form:**

```go
func (m EditorModel) ShowTemplateWizard(tmpl Template) tea.Cmd {
    // Use Huh form to prompt for values
    form := huh.NewForm(
        huh.NewGroup(
            huh.NewInput().
                Title("Node Name").
                Value(&m.templateName),
            huh.NewInput().
                Title("URL").
                Value(&m.templateURL),
        ),
    )

    return form.Run()
}
```

**Recommendation**: ✅ **Start with direct insertion** (simpler), add wizard later if needed

---

## 5. Live Validation

### Strategy: Debounced JSON Validation

**Requirements:**
1. Parse JSON on user pause (not every keystroke)
2. Show syntax errors (invalid JSON)
3. Show schema errors (invalid bento structure)
4. Display errors clearly without blocking editing

### Implementation

**Debouncing:**

```go
type EditorModel struct {
    textarea       textarea.Model
    validator      *neta.Validator
    lastEdit       time.Time
    validationErr  error
    validationTick int
}

const ValidationDelay = 500 * time.Millisecond

func (m EditorModel) Init() tea.Cmd {
    return tea.Tick(ValidationDelay, func(t time.Time) tea.Msg {
        return validationTickMsg{}
    })
}

type validationTickMsg struct{}

func (m EditorModel) Update(msg tea.Msg) (EditorModel, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Update last edit time
        m.lastEdit = time.Now()
        m.validationErr = nil // Clear old error while typing

        // Handle key
        var cmd tea.Cmd
        m.textarea, cmd = m.textarea.Update(msg)
        return m, cmd

    case validationTickMsg:
        // Check if enough time has passed since last edit
        if time.Since(m.lastEdit) >= ValidationDelay {
            m.validationErr = m.validate()
        }

        // Schedule next tick
        return m, tea.Tick(ValidationDelay, func(t time.Time) tea.Msg {
            return validationTickMsg{}
        })
    }

    return m, nil
}

func (m *EditorModel) validate() error {
    content := m.textarea.Value()

    // Step 1: Parse JSON
    var def neta.Definition
    if err := json.Unmarshal([]byte(content), &def); err != nil {
        return fmt.Errorf("JSON syntax error: %w", err)
    }

    // Step 2: Validate structure (version, type, etc.)
    if err := neta.ValidateVersion(def.Version); err != nil {
        return err
    }

    // Step 3: Validate node parameters
    if err := m.validator.ValidateRecursive(def); err != nil {
        return err
    }

    return nil
}
```

### Error Display

**Options:**

1. **Footer status line** (recommended for MVP)
2. **Sidebar error panel** (more complex)
3. **Inline annotations** (very complex, requires custom rendering)

**Footer Approach:**

```go
func (m EditorModel) View() string {
    var b strings.Builder

    // Main editor
    b.WriteString(m.textarea.View())
    b.WriteString("\n")

    // Status line
    if m.validationErr != nil {
        errorStyle := lipgloss.NewStyle().
            Foreground(lipgloss.Color("1")). // Red
            Bold(true)

        b.WriteString(errorStyle.Render("✗ " + m.validationErr.Error()))
    } else if m.textarea.Value() != "" {
        successStyle := lipgloss.NewStyle().
            Foreground(lipgloss.Color("2")). // Green
            Bold(true)

        b.WriteString(successStyle.Render("✓ Valid JSON"))
    }

    return b.String()
}
```

### Performance Considerations

**Large Files:**
- 100 nodes = ~5KB JSON (fast)
- 1000 nodes = ~50KB JSON (still fast with encoding/json)
- 10000+ nodes = may need optimization

**Optimization Strategies:**
1. Only validate visible portion for very large files
2. Cache validation results, only re-validate changed sections
3. Limit debounce delay to 500ms (good balance)

**Recommendation**: ✅ Start simple (full validation), optimize only if needed

---

## 6. Auto-Formatting

### Strategy: `json.MarshalIndent` from stdlib

**When to format:**
- Option A: On save (automatic)
- Option B: On hotkey (Ctrl+F or Ctrl+Shift+F)
- Option C: Both (recommended)

### Implementation

```go
func (m *EditorModel) FormatJSON() error {
    content := m.textarea.Value()

    // Parse JSON
    var data interface{}
    if err := json.Unmarshal([]byte(content), &data); err != nil {
        return fmt.Errorf("cannot format invalid JSON: %w", err)
    }

    // Format with 2-space indentation
    formatted, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        return err
    }

    // Update textarea
    m.textarea.SetValue(string(formatted))

    return nil
}

func (m EditorModel) Update(msg tea.Msg) (EditorModel, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Ctrl+F to format
        if msg.Type == tea.KeyCtrlF {
            if err := m.FormatJSON(); err != nil {
                // Show error in status
                m.validationErr = err
            }
            return m, nil
        }

        // Ctrl+S to save (with auto-format)
        if msg.Type == tea.KeyCtrlS {
            // Format before saving
            if err := m.FormatJSON(); err == nil {
                return m, m.save()
            }
            return m, nil
        }
    }

    // Default handling
    var cmd tea.Cmd
    m.textarea, cmd = m.textarea.Update(msg)
    return m, cmd
}
```

### Preserving Cursor Position

**Challenge**: Formatting changes text, cursor position becomes invalid

**Solution**: Calculate relative position before format, restore after

```go
func (m *EditorModel) FormatJSON() error {
    // Save cursor position (line and column)
    cursorBefore := m.textarea.CursorPosition()
    contentBefore := m.textarea.Value()

    // Count line number of cursor
    linesBefore := strings.Count(contentBefore[:cursorBefore], "\n")

    // Format
    var data interface{}
    if err := json.Unmarshal([]byte(contentBefore), &data); err != nil {
        return err
    }

    formatted, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        return err
    }

    // Update content
    m.textarea.SetValue(string(formatted))

    // Restore cursor to same line (approximate)
    linesAfter := strings.Split(string(formatted), "\n")
    if linesBefore < len(linesAfter) {
        // Move cursor to end of that line
        targetLine := linesAfter[linesBefore]
        newPos := 0
        for i := 0; i < linesBefore; i++ {
            newPos += len(linesAfter[i]) + 1 // +1 for newline
        }
        newPos += len(targetLine)
        m.textarea.SetCursor(newPos)
    }

    return nil
}
```

### Indentation: 2 Spaces vs 4 Spaces

**Recommendation**: ✅ **2 spaces** (more compact, standard for JSON)

---

## 7. Editor Layout

### Recommended: Full-Screen Modal

**Layout:**

```
┌─────────────────────────────────────────────────────────┐
│ 🍱 Bento Editor - bento.json                           │
├─────────────────────────────────────────────────────────┤
│   1 {                                                   │
│   2   "version": "1.0",                                 │
│   3   "name": "My Workflow",                            │
│   4   "description": "API data pipeline",               │
│   5   "nodes": [                                        │
│   6     {                                               │
│   7       "id": "node-1",                               │
│   8       "type": "http",                               │
│   9       "name": "Fetch Users",                        │
│  10       "parameters": {                               │
│  11         "method": "GET",                            │
│  12         "url": "https://api.example.com/users"      │
│  13       }                                             │
│  14     }                                               │
│  15   ]                                                 │
│  16 }                                                   │
│                                                         │
│                                                         │
├─────────────────────────────────────────────────────────┤
│ ✓ Valid JSON | Ctrl+S: Save | Ctrl+F: Format | Esc: Cancel │
└─────────────────────────────────────────────────────────┘
```

**Why Full-Screen:**
- JSON files can be large (need vertical space)
- Simpler implementation (no split-view complexity)
- Focus on editing (no distractions)
- Easy to add help panel later (toggle with ?)

**Alternative: Split-View**

```
┌─────────────────────┬───────────────────────┐
│ Editor (Plain)      │ Preview (Highlighted) │
│                     │                       │
│ {                   │ {                     │
│   "version": "1.0", │   "version": "1.0",   │
│   "name": "Test"    │   "name": "Test"      │
│ }                   │ }                     │
│                     │                       │
└─────────────────────┴───────────────────────┘
```

**Pros of split-view:**
- See highlighted syntax while editing
- Can compare formatted output

**Cons of split-view:**
- Reduces horizontal space (bad for deeply nested JSON)
- More complex to implement
- Not needed if we auto-format frequently

**Recommendation**: ✅ **Full-screen modal** for MVP, add split-view in Phase 4 if desired

### Responsive Layout

**Minimum Terminal Size**: 80x24 (standard)

**Layout Calculations:**

```go
func (m EditorModel) View() string {
    // Get terminal dimensions
    width := m.width
    height := m.height

    // Reserve space for header (3 lines) and footer (2 lines)
    editorHeight := height - 5

    // Set textarea dimensions
    m.textarea.SetWidth(width - 4) // -4 for borders
    m.textarea.SetHeight(editorHeight)

    // Render components
    header := renderHeader(width)
    editor := renderEditor(m.textarea, width, editorHeight)
    footer := renderFooter(m.validationErr, width)

    return lipgloss.JoinVertical(
        lipgloss.Left,
        header,
        editor,
        footer,
    )
}
```

---

## 8. Validation UX

### Error Display Strategy

**Live vs On-Save:**
- ✅ **Show errors live** (debounced 500ms after typing stops)
- ✅ **Allow saving invalid JSON** (draft saves)
- Reasoning: Don't block user workflow, allow work-in-progress

**Error Positioning:**

1. **Footer status line** (MVP - recommended)
   - Simple to implement
   - Always visible
   - Non-intrusive

2. **Inline highlights** (Future enhancement)
   - Show exact error location
   - Requires custom rendering
   - More complex

**Multiple Errors:**
- Show first error only (keep it simple)
- Add "Show all errors" command later (Ctrl+E)

### Visual Indicators

**Colors:**
- ✅ Valid: Green checkmark
- ❌ Invalid: Red X with error message
- ⏳ Validating: Yellow spinner (during debounce)

**Example:**

```go
func renderFooter(err error, validating bool) string {
    if validating {
        return lipgloss.NewStyle().
            Foreground(lipgloss.Color("3")). // Yellow
            Render("⏳ Validating...")
    }

    if err != nil {
        return lipgloss.NewStyle().
            Foreground(lipgloss.Color("1")). // Red
            Render("✗ " + err.Error())
    }

    return lipgloss.NewStyle().
        Foreground(lipgloss.Color("2")). // Green
        Render("✓ Valid JSON")
}
```

### Draft Saves

**Allow saving malformed JSON:**

```go
func (m *EditorModel) Save() error {
    content := m.textarea.Value()

    // Warn if invalid, but allow saving
    if err := m.validate(); err != nil {
        // Log warning or show confirmation
        // "Save invalid JSON? It may not run correctly."
    }

    // Save anyway
    return os.WriteFile(m.filePath, []byte(content), 0644)
}
```

---

## 9. YAML → JSON Migration Plan

### Migration Checklist

#### Phase 1: Replace YAML Parser with JSON Parser

**Files to Modify:**

1. **`pkg/jubako/parser.go`**
   - Replace YAML parser with JSON parser
   - Remove `gopkg.in/yaml.v3` import
   - Use `encoding/json` instead

```go
func (p *Parser) Parse(path string) (neta.Definition, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return neta.Definition{}, err
    }

    return p.ParseBytes(data)
}

func (p *Parser) ParseBytes(data []byte) (neta.Definition, error) {
    var def neta.Definition
    if err := json.Unmarshal(data, &def); err != nil {
        return neta.Definition{}, fmt.Errorf("invalid JSON: %w", err)
    }

    // Assign IDs and edges
    def = assignNodeIDs(def)
    def = autoGenerateEdges(def)

    if err := validateDefinition(def); err != nil {
        return neta.Definition{}, err
    }

    return def, nil
}

func (p *Parser) Format(def neta.Definition) ([]byte, error) {
    return json.MarshalIndent(def, "", "  ")
}
```

2. **`pkg/jubako/discovery.go`**
   - Update file glob to use `*.bento.json` only

```go
func DiscoverBentos(dir string) ([]BentoMetadata, error) {
    // Find .json files only
    jsonFiles, err := filepath.Glob(filepath.Join(dir, "*.bento.json"))
    if err != nil {
        return nil, err
    }

    // ... rest of discovery logic
}
```

3. **`pkg/jubako/store.go`**
   - Update `Save()` to use JSON format only

```go
func (s *Store) Save(def neta.Definition) error {
    data, err := json.MarshalIndent(def, "", "  ")
    if err != nil {
        return fmt.Errorf("marshal failed: %w", err)
    }

    path := filepath.Join(s.dir, def.Name+".bento.json")
    return os.WriteFile(path, data, 0644)
}
```

#### Phase 2: Update Examples

**Convert all example files to JSON:**

```bash
# Script to convert examples
for file in examples/*.yaml; do
    base=$(basename "$file" .bento.yaml)
    yq -o=json '.' "$file" > "examples/${base}.bento.json"
done
```

**Files to convert:**
- `examples/http-get.bento.yaml` → `examples/http-get.bento.json`
- `examples/group-sequence.bento.yaml` → `examples/group-sequence.bento.json`
- `examples/loop-for.bento.yaml` → `examples/loop-for.bento.json`
- (8 files total)

#### Phase 3: Update Templates

**Files to modify:**

1. **`pkg/omise/examples/templates/`**
   - Convert all `.bento.yaml` templates to `.bento.json`
   - Update `examples.go` to embed JSON files

```go
//go:embed templates/*.bento.json
var templateFS embed.FS
```

#### Phase 4: Update Tests

**Files to modify:**

1. **`pkg/jubako/parser_test.go`**
   - Add JSON parsing tests
   - Test dual format support

```go
func TestParser_ParseJSON(t *testing.T) {
    p := NewParser()

    jsonData := `{
        "version": "1.0",
        "type": "http",
        "name": "Test",
        "parameters": {
            "method": "GET",
            "url": "https://example.com"
        }
    }`

    def, err := p.ParseJSON([]byte(jsonData))
    assert.NoError(t, err)
    assert.Equal(t, "1.0", def.Version)
    assert.Equal(t, "http", def.Type)
}
```

2. **Other test files**
   - Update fixtures to use JSON where appropriate
   - Keep YAML tests for backward compatibility

#### Phase 5: Documentation

**Files to create/update:**

1. **`README.md`**
   - Document JSON format as primary
   - Note YAML still supported
   - Show JSON examples

2. **`.claude/strategy/phase-3-json-editor.md`**
   - Implementation guide (this will be created after research)

#### Phase 6: Remove YAML Support

**Clean up YAML remnants:**

1. Remove `gopkg.in/yaml.v3` from go.mod
2. Delete old `.bento.yaml` examples after verifying `.bento.json` versions
3. Update all documentation to show JSON only
4. Remove any YAML-specific normalization code (e.g., `normalizeDefinition()` that extracts nodes from parameters)

### Migration Complexity Estimate

- **Phase 1**: 1-2 hours (replace YAML parser with JSON parser)
- **Phase 2**: 30 minutes (convert examples)
- **Phase 3**: 1 hour (update templates)
- **Phase 4**: 1-2 hours (update tests)
- **Phase 5**: 30 minutes (documentation)
- **Phase 6**: 30 minutes (cleanup YAML remnants)

**Total**: 4-6 hours (simpler than dual format!)

---

## 10. JSON Bento Format Design

### Flat Array Structure with `parentId`

**Root Bento:**

```json
{
  "version": "1.0",
  "name": "example-workflow",
  "description": "Complete API workflow example",
  "icon": "🔄",
  "nodes": [
    {
      "id": "node-1",
      "type": "http",
      "name": "Fetch Users",
      "parameters": {
        "method": "GET",
        "url": "https://api.example.com/users",
        "headers": {
          "Accept": "application/json"
        }
      }
    },
    {
      "id": "node-2",
      "type": "transform.jq",
      "name": "Extract IDs",
      "parameters": {
        "query": ".data | map(.id)"
      }
    },
    {
      "id": "node-3",
      "type": "file.write",
      "name": "Save Results",
      "parameters": {
        "path": "/tmp/user-ids.json"
      }
    }
  ],
  "edges": [
    {
      "id": "edge-1",
      "source": "node-1",
      "target": "node-2"
    },
    {
      "id": "edge-2",
      "source": "node-2",
      "target": "node-3"
    }
  ]
}
```

### Node Type Schemas

#### 1. HTTP Node

```json
{
  "id": "http-1",
  "type": "http",
  "name": "API Request",
  "description": "Fetch data from REST API",
  "parameters": {
    "method": "POST",
    "url": "https://api.example.com/data",
    "headers": {
      "Content-Type": "application/json",
      "Authorization": "Bearer TOKEN"
    },
    "body": "{\"key\": \"value\"}"
  }
}
```

**Parameters:**
- `method` (string, optional, default: "GET"): HTTP verb (GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS)
- `url` (string, required): Full HTTP(S) URL
- `headers` (object, optional): Key-value pairs of headers
- `body` (string, optional): Request body for POST/PUT/PATCH

#### 2. Transform (jq) Node

```json
{
  "id": "jq-1",
  "type": "transform.jq",
  "name": "Extract Fields",
  "parameters": {
    "query": ".users | map({id: .id, name: .name})",
    "input": "{\"users\": [{\"id\": 1, \"name\": \"Alice\"}]}"
  }
}
```

**Parameters:**
- `query` (string, required): jq filter expression
- `input` (string, optional): Static input data (if not using previous node output)

#### 3. File Write Node

```json
{
  "id": "file-1",
  "type": "file.write",
  "name": "Save Output",
  "parameters": {
    "path": "/tmp/output.json",
    "content": "{\"result\": \"success\"}"
  }
}
```

**Parameters:**
- `path` (string, required): File path to write to
- `content` (string, optional): Content to write (uses previous node output if omitted)

#### 4. Sequence Group Node

```json
{
  "id": "seq-1",
  "type": "group.sequence",
  "name": "Sequential Pipeline",
  "nodes": [
    {
      "id": "seq-1-node-1",
      "parentId": "seq-1",
      "type": "http",
      "name": "Step 1",
      "parameters": {
        "method": "GET",
        "url": "https://example.com/step1"
      }
    },
    {
      "id": "seq-1-node-2",
      "parentId": "seq-1",
      "type": "http",
      "name": "Step 2",
      "parameters": {
        "method": "GET",
        "url": "https://example.com/step2"
      }
    }
  ],
  "edges": [
    {
      "id": "seq-1-edge-1",
      "source": "seq-1-node-1",
      "target": "seq-1-node-2"
    }
  ]
}
```

**Parameters:**
- `nodes` (array, required): Child nodes to execute sequentially
- `edges` (array, optional): Explicit edges (auto-generated if omitted)

#### 5. Parallel Group Node

```json
{
  "id": "par-1",
  "type": "group.parallel",
  "name": "Parallel Fetch",
  "parameters": {
    "max_concurrent": 3
  },
  "nodes": [
    {
      "id": "par-1-node-1",
      "parentId": "par-1",
      "type": "http",
      "name": "Fetch A",
      "parameters": {
        "method": "GET",
        "url": "https://example.com/a"
      }
    },
    {
      "id": "par-1-node-2",
      "parentId": "par-1",
      "type": "http",
      "name": "Fetch B",
      "parameters": {
        "method": "GET",
        "url": "https://example.com/b"
      }
    }
  ]
}
```

**Parameters:**
- `max_concurrent` (integer, optional, default: 0): Max parallel executions (0 = unlimited)
- `nodes` (array, required): Child nodes to execute in parallel

#### 6. For Loop Node

```json
{
  "id": "loop-1",
  "type": "loop.for",
  "name": "Process Items",
  "parameters": {
    "items": ["item1", "item2", "item3"],
    "body": {
      "type": "http",
      "name": "Process Item",
      "parameters": {
        "method": "POST",
        "url": "https://example.com/process"
      }
    }
  }
}
```

**Parameters:**
- `items` (array, required): Array to iterate over
- `body` (object, required): Node definition to execute for each item

#### 7. Conditional (If) Node

```json
{
  "id": "if-1",
  "type": "conditional.if",
  "name": "Check Status",
  "parameters": {
    "condition": true,
    "then": {
      "type": "http",
      "name": "Success Path",
      "parameters": {
        "method": "GET",
        "url": "https://example.com/success"
      }
    },
    "else": {
      "type": "http",
      "name": "Failure Path",
      "parameters": {
        "method": "GET",
        "url": "https://example.com/failure"
      }
    }
  }
}
```

**Parameters:**
- `condition` (boolean, required): Condition to evaluate
- `then` (object, optional): Node to execute if condition is true
- `else` (object, optional): Node to execute if condition is false

### Complete Example: Nested Workflow

```json
{
  "version": "1.0",
  "name": "complex-workflow",
  "description": "API pipeline with loops and conditions",
  "icon": "⚙️",
  "nodes": [
    {
      "id": "fetch-users",
      "type": "http",
      "name": "Fetch All Users",
      "parameters": {
        "method": "GET",
        "url": "https://api.example.com/users"
      }
    },
    {
      "id": "extract-active",
      "type": "transform.jq",
      "name": "Filter Active Users",
      "parameters": {
        "query": ".users | map(select(.active == true))"
      }
    },
    {
      "id": "process-users",
      "type": "loop.for",
      "name": "Process Each User",
      "parameters": {
        "items": [],
        "body": {
          "type": "group.sequence",
          "name": "User Processing",
          "nodes": [
            {
              "type": "http",
              "name": "Fetch User Details",
              "parameters": {
                "method": "GET",
                "url": "https://api.example.com/users/{{item.id}}"
              }
            },
            {
              "type": "conditional.if",
              "name": "Check Premium Status",
              "parameters": {
                "condition": false,
                "then": {
                  "type": "http",
                  "name": "Send Premium Email",
                  "parameters": {
                    "method": "POST",
                    "url": "https://api.example.com/email/premium"
                  }
                },
                "else": {
                  "type": "http",
                  "name": "Send Standard Email",
                  "parameters": {
                    "method": "POST",
                    "url": "https://api.example.com/email/standard"
                  }
                }
              }
            }
          ]
        }
      }
    },
    {
      "id": "save-results",
      "type": "file.write",
      "name": "Save Processing Results",
      "parameters": {
        "path": "/tmp/results.json"
      }
    }
  ],
  "edges": [
    {
      "id": "edge-1",
      "source": "fetch-users",
      "target": "extract-active"
    },
    {
      "id": "edge-2",
      "source": "extract-active",
      "target": "process-users"
    },
    {
      "id": "edge-3",
      "source": "process-users",
      "target": "save-results"
    }
  ]
}
```

---

## 11. Architecture Recommendation

### Component Structure

```
pkg/omise/screens/json_editor/
├── editor.go         (<250 lines) - Main editor model
├── validation.go     (<200 lines) - Debounced validation logic
├── formatting.go     (<150 lines) - JSON formatting utilities
├── templates.go      (<200 lines) - Node template definitions
├── keymap.go         (<150 lines) - Editor-specific key bindings
└── messages.go       (<100 lines) - Custom Bubble Tea messages
```

### Editor Model

```go
package json_editor

import (
    "time"
    "github.com/charmbracelet/bubbles/textarea"
    "github.com/charmbracelet/bubbletea"
    "bento/pkg/neta"
)

type Model struct {
    // UI components
    textarea   textarea.Model
    width      int
    height     int

    // Validation state
    validator      *neta.Validator
    lastEdit       time.Time
    validationErr  error
    isValidating   bool

    // File state
    filePath   string
    modified   bool

    // Template state
    nextNodeID int
}

func New(width, height int, filePath string) Model {
    ta := textarea.New()
    ta.ShowLineNumbers = true
    ta.Focus()

    return Model{
        textarea:   ta,
        width:      width,
        height:     height,
        filePath:   filePath,
        validator:  neta.NewValidator(),
        nextNodeID: 1,
    }
}

func (m Model) Init() tea.Cmd {
    return tea.Batch(
        textarea.Blink,
        m.startValidationTicker(),
    )
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Handle messages
    // ...
}

func (m Model) View() string {
    // Render editor
    // ...
}
```

### Integration with Browser

**Launch editor from browser:**

```go
// In browser Update()
case tea.KeyMsg:
    if msg.String() == "e" {
        // Get selected bento
        selected := m.list.SelectedItem()

        // Load bento
        def, err := m.jubako.Load(selected.Path)
        if err != nil {
            return m, showError(err)
        }

        // Convert to JSON
        jsonData, err := json.MarshalIndent(def, "", "  ")
        if err != nil {
            return m, showError(err)
        }

        // Create editor
        editor := json_editor.New(m.width, m.height, selected.Path)
        editor.SetContent(string(jsonData))

        return editor, editor.Init()
    }
```

### Save Flow

```go
func (m Model) save() tea.Cmd {
    return func() tea.Msg {
        content := m.textarea.Value()

        // Parse JSON
        var def neta.Definition
        if err := json.Unmarshal([]byte(content), &def); err != nil {
            return saveErrorMsg{err}
        }

        // Save to file
        if err := os.WriteFile(m.filePath, []byte(content), 0644); err != nil {
            return saveErrorMsg{err}
        }

        return saveSuccessMsg{}
    }
}
```

---

## 12. Implementation Complexity

### Feature Breakdown

| Feature | Complexity | Est. Time | Dependencies |
|---------|------------|-----------|--------------|
| Basic textarea editor | Low | 1 hour | bubbles/textarea |
| Line numbers | Trivial | 15 min | bubbles (built-in) |
| JSON validation | Low | 2 hours | encoding/json, neta.Validator |
| Debounced validation | Medium | 1 hour | time.Ticker |
| Auto-formatting | Low | 1 hour | encoding/json |
| Template system | Medium | 3 hours | Custom code |
| Hotkey bindings | Low | 1 hour | bubbles/key |
| Syntax highlighting (split-view) | Medium | 3 hours | chroma |
| Error display (footer) | Low | 1 hour | lipgloss |
| Save/load integration | Low | 2 hours | jubako |
| YAML→JSON migration | Medium | 5-7 hours | Multiple files |

**MVP (Without syntax highlighting):**
- Basic editor: 1 hour
- Validation: 3 hours
- Formatting: 1 hour
- Templates: 3 hours
- Integration: 2 hours
- **Total**: ~10 hours

**Full Implementation (With syntax highlighting):**
- MVP: 10 hours
- Syntax highlighting: 3 hours
- Polish/testing: 2 hours
- **Total**: ~15 hours

**Including Migration:**
- Editor: 15 hours
- Migration: 7 hours
- **Total**: ~22 hours

### Bento Box Compliance

**File Size Targets:**
- `editor.go`: ~200 lines (model definition + core logic)
- `validation.go`: ~150 lines (validation ticker + error handling)
- `formatting.go`: ~100 lines (JSON formatting utilities)
- `templates.go`: ~200 lines (template definitions)
- `keymap.go`: ~100 lines (key bindings)
- `messages.go`: ~50 lines (custom messages)

**Total**: ~800 lines across 6 files (avg 133 lines/file) ✅ Well under 250 line limit

---

## 13. Dependencies

### Required Libraries (Already in go.mod)

✅ **No new dependencies needed!**

1. `github.com/charmbracelet/bubbles` - textarea component
2. `github.com/charmbracelet/bubbletea` - TUI framework
3. `github.com/charmbracelet/lipgloss` - Styling
4. `encoding/json` - JSON parsing (stdlib)
5. `bento/pkg/neta` - Validation framework (internal)

### Optional Dependencies (Future Enhancements)

1. `github.com/alecthomas/chroma` - Syntax highlighting
   - Version: `v2.14.0` (latest stable)
   - License: MIT
   - Size: ~2MB

**Add chroma:**
```bash
go get github.com/alecthomas/chroma/v2
```

---

## 14. Risks & Concerns

### 1. Textarea Performance with Large Files

**Risk**: Slow rendering/editing with large bentos (1000+ nodes, 50KB+ JSON)

**Mitigation**:
- Test with large files early
- Consider viewport/pagination for very large files
- Add file size warning (e.g., "File >100KB, editor may be slow")

**Likelihood**: Low (most bentos will be <10KB)

### 2. Cursor Position After Formatting

**Risk**: Cursor jumps to wrong position after auto-format

**Mitigation**:
- Implement cursor position preservation (see section 6)
- Test thoroughly with various cursor positions
- Make formatting opt-in (hotkey) rather than automatic

**Likelihood**: Medium (known issue with text editors)

### 3. Syntax Highlighting Integration

**Risk**: Chroma highlighting doesn't work well with textarea (plain text component)

**Mitigation**:
- Use split-view approach (editor + preview)
- Make highlighting optional feature
- Start without highlighting for MVP

**Likelihood**: High (textarea doesn't natively support styled content)

### 4. Validation Performance

**Risk**: Validation slows down editing on large files

**Mitigation**:
- Use debouncing (500ms delay)
- Cache validation results
- Make validation optional (toggle with hotkey)

**Likelihood**: Low (JSON parsing is fast)

### 5. Template Insertion Complexity

**Risk**: Templates don't insert at correct indentation level

**Mitigation**:
- Detect indentation from current line
- Test with various nesting levels
- Provide "Format" hotkey to fix indentation issues

**Likelihood**: Medium (indentation detection can be tricky)

### 6. YAML→JSON Compatibility

**Risk**: Existing YAML files don't convert cleanly to JSON

**Mitigation**:
- Maintain YAML parser indefinitely
- Provide migration tool with validation
- Test migration with all existing examples

**Likelihood**: Low (struct tags already support both formats)

### 7. User Adoption

**Risk**: Users prefer YAML or don't want to switch

**Mitigation**:
- Support both formats long-term
- Make JSON the default for new bentos
- Provide clear migration documentation

**Likelihood**: Medium (user preference varies)

### 8. Breaking Changes

**Risk**: JSON format introduces breaking changes

**Mitigation**:
- Keep version field ("1.0") for compatibility
- Add migration guide
- Support YAML indefinitely
- Feature-flag JSON editor initially

**Likelihood**: Low (dual format support prevents breakage)

---

## 15. References

### Libraries & Documentation

1. **Bubbles Textarea**
   - Docs: https://pkg.go.dev/github.com/charmbracelet/bubbles/textarea
   - Examples: https://github.com/charmbracelet/bubbletea/tree/main/examples
   - Key bindings: https://github.com/charmbracelet/bubbles?tab=readme-ov-file#key

2. **Chroma Syntax Highlighting**
   - GitHub: https://github.com/alecthomas/chroma
   - Docs: https://pkg.go.dev/github.com/alecthomas/chroma
   - Terminal formatters: https://pkg.go.dev/github.com/alecthomas/chroma/formatters

3. **Bubble Tea Best Practices**
   - Building Bubble Tea Programs: https://leg100.github.io/en/posts/building-bubbletea-programs/
   - Official Tutorial: https://github.com/charmbracelet/bubbletea/tree/main/tutorials

4. **Go JSON Package**
   - Encoding/JSON: https://pkg.go.dev/encoding/json
   - JSON and Go: https://go.dev/blog/json

### Similar Projects

1. **JSON Path Evaluator** (Bubble Tea + JSON)
   - GitHub: https://github.com/tyler-mairose-sp/json-path-evaluator
   - Uses textarea for JSON editing

2. **jqp** (Interactive jq processor)
   - GitHub: https://github.com/noahgorstein/jqp
   - Bubble Tea TUI for jq queries

### Bento Codebase

- Strategy docs: `.claude/strategy/`
- Bento Box Principle: `.claude/BENTO_BOX_PRINCIPLE.md`
- Phase 2 (Tab Navigation): `.claude/strategy/phase-2-navigation.md`
- Validation framework: `pkg/neta/schemas/`

---

## 16. Prototype Code Snippets

### 1. Basic Editor Setup

```go
package json_editor

import (
    "github.com/charmbracelet/bubbles/textarea"
    "github.com/charmbracelet/bubbletea"
)

type Model struct {
    textarea textarea.Model
}

func New() Model {
    ta := textarea.New()
    ta.Placeholder = "Enter JSON here..."
    ta.ShowLineNumbers = true
    ta.CharLimit = 0 // No limit
    ta.SetWidth(80)
    ta.SetHeight(24)
    ta.Focus()

    return Model{textarea: ta}
}

func (m Model) Init() tea.Cmd {
    return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyCtrlC, tea.KeyEsc:
            return m, tea.Quit
        }
    }

    m.textarea, cmd = m.textarea.Update(msg)
    return m, cmd
}

func (m Model) View() string {
    return m.textarea.View()
}
```

### 2. JSON Validation

```go
package json_editor

import (
    "encoding/json"
    "fmt"
    "time"

    "bento/pkg/neta"
    tea "github.com/charmbracelet/bubbletea"
)

const ValidationDelay = 500 * time.Millisecond

type validationTickMsg struct{}

func (m Model) startValidationTicker() tea.Cmd {
    return tea.Tick(ValidationDelay, func(t time.Time) tea.Msg {
        return validationTickMsg{}
    })
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Mark as edited
        m.lastEdit = time.Now()
        m.validationErr = nil

    case validationTickMsg:
        // Validate if enough time passed
        if time.Since(m.lastEdit) >= ValidationDelay {
            m.validationErr = m.validate()
        }

        // Schedule next tick
        return m, m.startValidationTicker()
    }

    var cmd tea.Cmd
    m.textarea, cmd = m.textarea.Update(msg)
    return m, cmd
}

func (m Model) validate() error {
    content := m.textarea.Value()
    if content == "" {
        return nil
    }

    // Parse JSON
    var def neta.Definition
    if err := json.Unmarshal([]byte(content), &def); err != nil {
        return fmt.Errorf("JSON error: %w", err)
    }

    // Validate structure
    if err := neta.ValidateVersion(def.Version); err != nil {
        return err
    }

    // Validate node parameters
    validator := neta.NewValidator()
    if err := validator.ValidateRecursive(def); err != nil {
        return err
    }

    return nil
}
```

### 3. Auto-Formatting

```go
package json_editor

import (
    "encoding/json"
    tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Ctrl+F to format
        if msg.Type == tea.KeyCtrlF {
            if err := m.formatJSON(); err != nil {
                m.validationErr = err
            }
            return m, nil
        }

        // Ctrl+S to save (with auto-format)
        if msg.Type == tea.KeyCtrlS {
            m.formatJSON() // Ignore errors, save anyway
            return m, m.save()
        }
    }

    var cmd tea.Cmd
    m.textarea, cmd = m.textarea.Update(msg)
    return m, cmd
}

func (m *Model) formatJSON() error {
    content := m.textarea.Value()

    // Parse
    var data interface{}
    if err := json.Unmarshal([]byte(content), &data); err != nil {
        return err
    }

    // Format with 2-space indent
    formatted, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        return err
    }

    // Update textarea
    m.textarea.SetValue(string(formatted))
    m.modified = true

    return nil
}
```

### 4. Template Insertion

```go
package json_editor

import (
    "fmt"
    "strings"
    tea "github.com/charmbracelet/bubbletea"
)

type Template struct {
    Name   string
    Hotkey tea.KeyType
    JSON   string
}

var templates = []Template{
    {
        Name:   "HTTP GET",
        Hotkey: tea.KeyCtrlH,
        JSON: `{
  "id": "node-$ID",
  "type": "http",
  "name": "New Request",
  "parameters": {
    "method": "GET",
    "url": "https://example.com"
  }
}`,
    },
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Check template hotkeys
        for _, tmpl := range templates {
            if msg.Type == tmpl.Hotkey {
                m.insertTemplate(tmpl)
                return m, nil
            }
        }
    }

    var cmd tea.Cmd
    m.textarea, cmd = m.textarea.Update(msg)
    return m, cmd
}

func (m *Model) insertTemplate(tmpl Template) {
    // Generate ID
    nodeID := fmt.Sprintf("node-%d", m.nextNodeID)
    m.nextNodeID++

    // Replace placeholders
    content := strings.ReplaceAll(tmpl.JSON, "$ID", nodeID)

    // Insert at cursor
    current := m.textarea.Value()
    cursorPos := m.textarea.CursorPosition()

    newContent := current[:cursorPos] + content + current[cursorPos:]
    m.textarea.SetValue(newContent)
    m.modified = true
}
```

### 5. Error Display (Footer)

```go
package json_editor

import (
    "strings"
    "github.com/charmbracelet/lipgloss"
)

var (
    errorStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("1")).
        Bold(true)

    successStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("2")).
        Bold(true)
)

func (m Model) View() string {
    var b strings.Builder

    // Main editor
    b.WriteString(m.textarea.View())
    b.WriteString("\n\n")

    // Status line
    if m.validationErr != nil {
        b.WriteString(errorStyle.Render("✗ " + m.validationErr.Error()))
    } else if m.textarea.Value() != "" {
        b.WriteString(successStyle.Render("✓ Valid JSON"))
    }

    // Help line
    b.WriteString("\n")
    b.WriteString(lipgloss.NewStyle().
        Foreground(lipgloss.Color("8")).
        Render("Ctrl+S: Save | Ctrl+F: Format | Esc: Cancel"))

    return b.String()
}
```

---

## Conclusion

### Summary

The JSON editor for Bento is **highly feasible** and can be implemented with **low risk** using existing Charm ecosystem components. The migration from YAML to JSON is **straightforward** thanks to existing dual struct tags.

### Recommended Approach

1. **Start with MVP** (10 hours):
   - Basic textarea editor with line numbers
   - JSON validation (debounced)
   - Auto-formatting (Ctrl+F)
   - Template hotkeys for common nodes
   - Save/load integration

2. **Add Enhancements** (5 hours):
   - Syntax highlighting (split-view with Chroma)
   - Template wizard (Huh forms for complex nodes)
   - Inline error annotations

3. **Perform Migration** (7 hours):
   - Add JSON parser alongside YAML
   - Convert all examples to JSON
   - Update documentation
   - Support both formats long-term

### Next Steps

1. ✅ **Approve research findings**
2. Create implementation plan: `.claude/strategy/phase-3-json-editor.md`
3. Begin MVP implementation
4. Test with existing bentos
5. Gather user feedback
6. Iterate on enhancements

**Total Estimated Time**: 22 hours (MVP + enhancements + migration)

---

**Research Completed**: 2025-10-17
**Researcher**: Claude (Sonnet 4.5)
**Status**: Ready for Implementation Planning
