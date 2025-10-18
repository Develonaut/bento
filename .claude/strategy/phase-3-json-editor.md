# Phase 3: JSON Bento Editor - Implementation Plan

**Based on**: [phase-3-json-editor-research.md](./phase-3-json-editor-research.md)
**Duration**: 19 hours total (10h MVP + 5h enhancements + 4h migration)
**Status**: Ready for Implementation

---

## Overview

Implement a JSON-based text editor for creating and editing bento files, replacing the complex interactive editor with a simpler approach focused on direct JSON editing with helpful features.

### Goals

1. ✅ Text-based JSON editor (not visual/interactive builder)
2. ✅ Syntax highlighting for JSON
3. ✅ Live validation as you type
4. ✅ Template/snippet insertion via hotkeys
5. ✅ Auto-formatting on save or hotkey
6. ✅ Replace YAML with JSON format (clean migration)

### Non-Goals

- ❌ Visual drag-and-drop node builder
- ❌ Real-time graph visualization (future phase)
- ❌ Multi-file editing
- ❌ Git integration
- ❌ Dual YAML/JSON format support (JSON only!)

---

## Phase 3A: MVP Editor (10 hours)

### Deliverables

1. Basic textarea editor with line numbers
2. JSON validation (debounced 500ms)
3. Auto-formatting (Ctrl+F hotkey)
4. Template hotkeys for 7 common node types
5. Save/load integration with Jubako
6. Error display in footer

### Implementation Steps

#### Step 1: Create Package Structure (30 min)

**Files to create:**

```
pkg/omise/screens/json_editor/
├── editor.go         - Main editor model
├── validation.go     - Debounced validation
├── formatting.go     - JSON formatting
├── templates.go      - Node templates
├── keymap.go         - Key bindings
└── messages.go       - Custom messages
```

**Bento Box Compliance:**
- Each file < 250 lines
- Functions < 20 lines
- Clear separation of concerns

#### Step 2: Implement Basic Editor (2 hours)

**File**: `pkg/omise/screens/json_editor/editor.go`

**Requirements:**
- Use `bubbles/textarea` component
- Enable line numbers (`ShowLineNumbers = true`)
- Handle window resize
- Track file path and modified state
- Support loading existing JSON files

**Model Structure:**

```go
package json_editor

import (
    "github.com/charmbracelet/bubbles/textarea"
    "bento/pkg/neta"
)

type Model struct {
    // UI
    textarea textarea.Model
    width    int
    height   int

    // File state
    filePath string
    modified bool

    // Validation
    validator     *neta.Validator
    lastEdit      time.Time
    validationErr error

    // Templates
    nextNodeID int
}

func New(width, height int, filePath string) Model {
    ta := textarea.New()
    ta.ShowLineNumbers = true
    ta.CharLimit = 0
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
```

**Key Bindings:**
- `Ctrl+S` - Save
- `Ctrl+F` - Format
- `Ctrl+Q` or `Esc` - Cancel/Exit
- Arrow keys - Navigate
- `Ctrl+H` - Insert HTTP GET template
- `Ctrl+P` - Insert HTTP POST template
- `Ctrl+J` - Insert jq transform template
- `Ctrl+L` - Insert for loop template

**Testing:**
- Create new editor
- Load JSON file
- Basic text editing
- Window resize

#### Step 3: Implement Validation (2 hours)

**File**: `pkg/omise/screens/json_editor/validation.go`

**Requirements:**
- Debounce validation (500ms after last keystroke)
- Parse JSON with `encoding/json`
- Validate bento structure with `neta.Validator`
- Display errors in footer

**Implementation:**

```go
const ValidationDelay = 500 * time.Millisecond

type validationTickMsg struct{}

func (m Model) startValidationTicker() tea.Cmd {
    return tea.Tick(ValidationDelay, func(t time.Time) tea.Msg {
        return validationTickMsg{}
    })
}

func (m *Model) validate() error {
    content := m.textarea.Value()
    if content == "" {
        return nil
    }

    // Step 1: Parse JSON
    var def neta.Definition
    if err := json.Unmarshal([]byte(content), &def); err != nil {
        return fmt.Errorf("JSON syntax error: %w", err)
    }

    // Step 2: Validate version and structure
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

**Testing:**
- Valid JSON → no error
- Invalid JSON syntax → syntax error shown
- Missing required fields → validation error shown
- Validation happens 500ms after typing stops

#### Step 4: Implement Auto-Formatting (1 hour)

**File**: `pkg/omise/screens/json_editor/formatting.go`

**Requirements:**
- Format JSON with 2-space indentation
- Preserve cursor position (approximate)
- Handle malformed JSON gracefully
- Trigger on Ctrl+F or before save

**Implementation:**

```go
func (m *Model) formatJSON() error {
    content := m.textarea.Value()

    // Parse
    var data interface{}
    if err := json.Unmarshal([]byte(content), &data); err != nil {
        return fmt.Errorf("cannot format invalid JSON: %w", err)
    }

    // Format
    formatted, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        return err
    }

    // Update
    m.textarea.SetValue(string(formatted))
    m.modified = true

    return nil
}
```

**Testing:**
- Unformatted JSON → formatted with 2 spaces
- Invalid JSON → error message, no change
- Cursor position preserved (approximately)

#### Step 5: Implement Templates (3 hours)

**File**: `pkg/omise/screens/json_editor/templates.go`

**Requirements:**
- Define 7 node templates (HTTP GET/POST, jq, sequence, parallel, for loop, if, file write)
- Insert at cursor position
- Auto-generate unique IDs
- Replace placeholders with sensible defaults

**Template Definitions:**

```go
type Template struct {
    Name   string
    Hotkey tea.KeyType
    JSON   string
}

var NodeTemplates = []Template{
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
    {
        Name:   "HTTP POST",
        Hotkey: tea.KeyCtrlP,
        JSON: `{
  "id": "node-$ID",
  "type": "http",
  "name": "New POST Request",
  "parameters": {
    "method": "POST",
    "url": "https://example.com",
    "headers": {
      "Content-Type": "application/json"
    },
    "body": ""
  }
}`,
    },
    {
        Name:   "Transform (jq)",
        Hotkey: tea.KeyCtrlJ,
        JSON: `{
  "id": "node-$ID",
  "type": "transform.jq",
  "name": "Transform Data",
  "parameters": {
    "query": ".data"
  }
}`,
    },
    {
        Name:   "For Loop",
        Hotkey: tea.KeyCtrlL,
        JSON: `{
  "id": "node-$ID",
  "type": "loop.for",
  "name": "Process Items",
  "parameters": {
    "items": [],
    "body": {
      "type": "http",
      "name": "Process Item",
      "parameters": {
        "method": "GET",
        "url": "https://example.com"
      }
    }
  }
}`,
    },
    // ... more templates
}
```

**Insertion Logic:**

```go
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

**Testing:**
- Each hotkey inserts correct template
- IDs are unique and incremental
- Templates insert at cursor position
- JSON remains valid after insertion

#### Step 6: Implement Save/Load (1 hour)

**File**: `pkg/omise/screens/json_editor/editor.go` (add methods)

**Requirements:**
- Save JSON to file
- Format before saving
- Show success/error message
- Integrate with Jubako for file management

**Implementation:**

```go
func (m *Model) save() tea.Cmd {
    return func() tea.Msg {
        // Format first
        if err := m.formatJSON(); err != nil {
            return saveErrorMsg{err}
        }

        content := m.textarea.Value()

        // Validate before saving
        var def neta.Definition
        if err := json.Unmarshal([]byte(content), &def); err != nil {
            return saveErrorMsg{fmt.Errorf("invalid JSON: %w", err)}
        }

        // Save to file
        if err := os.WriteFile(m.filePath, []byte(content), 0644); err != nil {
            return saveErrorMsg{err}
        }

        return saveSuccessMsg{path: m.filePath}
    }
}

func (m *Model) load(filePath string) tea.Cmd {
    return func() tea.Msg {
        data, err := os.ReadFile(filePath)
        if err != nil {
            return loadErrorMsg{err}
        }

        return loadSuccessMsg{content: string(data)}
    }
}
```

**Testing:**
- Save valid JSON → file written
- Save invalid JSON → error shown, file not written
- Load existing file → content displayed
- Load non-existent file → error shown

#### Step 7: Implement UI Layout (1 hour)

**File**: `pkg/omise/screens/json_editor/editor.go` (View method)

**Requirements:**
- Full-screen modal layout
- Header with file path
- Main editor area with line numbers
- Footer with status and help

**Layout:**

```
┌─────────────────────────────────────────────────────────┐
│ 🍱 Bento Editor - example.bento.json                   │
├─────────────────────────────────────────────────────────┤
│   1 {                                                   │
│   2   "version": "1.0",                                 │
│   3   "name": "My Workflow",                            │
│   ...                                                   │
├─────────────────────────────────────────────────────────┤
│ ✓ Valid JSON | Ctrl+S: Save | Ctrl+F: Format | Esc: Exit │
└─────────────────────────────────────────────────────────┘
```

**Implementation:**

```go
func (m Model) View() string {
    header := renderHeader(m.filePath, m.width)
    editor := m.textarea.View()
    footer := renderFooter(m.validationErr, m.modified, m.width)

    return lipgloss.JoinVertical(
        lipgloss.Left,
        header,
        editor,
        footer,
    )
}

func renderHeader(filePath string, width int) string {
    title := lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("6")).
        Render("🍱 Bento Editor - " + filepath.Base(filePath))

    headerStyle := lipgloss.NewStyle().
        Width(width).
        BorderStyle(lipgloss.NormalBorder()).
        BorderBottom(true)

    return headerStyle.Render(title)
}

func renderFooter(err error, modified bool, width int) string {
    var status string
    if err != nil {
        status = lipgloss.NewStyle().
            Foreground(lipgloss.Color("1")).
            Render("✗ " + err.Error())
    } else {
        status = lipgloss.NewStyle().
            Foreground(lipgloss.Color("2")).
            Render("✓ Valid JSON")
    }

    if modified {
        status += " [Modified]"
    }

    help := lipgloss.NewStyle().
        Foreground(lipgloss.Color("8")).
        Render("Ctrl+S: Save | Ctrl+F: Format | Esc: Exit")

    footerStyle := lipgloss.NewStyle().
        Width(width).
        BorderStyle(lipgloss.NormalBorder()).
        BorderTop(true)

    return footerStyle.Render(status + "\n" + help)
}
```

**Testing:**
- Header shows file path
- Footer shows validation status
- Help text is visible
- Layout adapts to terminal size

#### Step 8: Integration with Browser (30 min)

**File**: `pkg/omise/update.go` and `pkg/omise/screens/browser_handlers.go`

**Requirements:**
- Launch editor when pressing 'e' in browser
- Load selected bento file
- Return to browser on exit
- Refresh browser list after save

**Implementation:**

```go
// In browser_handlers.go
func (m Model) handleEdit() (Model, tea.Cmd) {
    selected := m.browserModel.list.SelectedItem()
    if selected == nil {
        return m, nil
    }

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

    m.currentScreen = ScreenJSONEditor
    m.jsonEditorModel = editor

    return m, editor.Init()
}
```

**Testing:**
- Press 'e' in browser → editor opens
- Editor shows bento content
- Press Esc in editor → returns to browser
- Save in editor → browser list refreshes

---

## Phase 3B: Enhancements (5 hours)

### Deliverables

1. Syntax highlighting (split-view with Chroma)
2. Template wizard (Huh forms for complex nodes)
3. Improved error display (line numbers, context)

### Implementation Steps

#### Step 1: Add Chroma Dependency (15 min)

```bash
go get github.com/alecthomas/chroma/v2
```

**Update go.mod:**

```go
require (
    // ... existing deps
    github.com/alecthomas/chroma/v2 v2.14.0
)
```

#### Step 2: Implement Syntax Highlighting (3 hours)

**File**: `pkg/omise/screens/json_editor/highlighting.go`

**Requirements:**
- Use Chroma to highlight JSON
- Split-view: plain editor (left) + highlighted preview (right)
- Toggle split-view with Tab key
- Use terminal color profile (true color if supported)

**Implementation:**

```go
package json_editor

import (
    "bytes"
    "github.com/alecthomas/chroma/v2"
    "github.com/alecthomas/chroma/v2/formatters"
    "github.com/alecthomas/chroma/v2/lexers"
    "github.com/alecthomas/chroma/v2/styles"
)

func highlightJSON(content string) (string, error) {
    lexer := lexers.Get("json")
    if lexer == nil {
        lexer = lexers.Fallback
    }

    style := styles.Get("monokai")
    if style == nil {
        style = styles.Fallback
    }

    formatter := formatters.Get("terminal16m")
    if formatter == nil {
        formatter = formatters.Fallback
    }

    iterator, err := lexer.Tokenise(nil, content)
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

**Split-View Layout:**

```go
func (m Model) View() string {
    if !m.showPreview {
        // Single pane (editor only)
        return m.renderFullEditor()
    }

    // Split view
    editorPane := m.renderEditor(m.width / 2)
    previewPane := m.renderPreview(m.width / 2)

    return lipgloss.JoinHorizontal(
        lipgloss.Top,
        editorPane,
        previewPane,
    )
}

func (m Model) renderPreview(width int) string {
    content := m.textarea.Value()
    highlighted, err := highlightJSON(content)
    if err != nil {
        return "Error highlighting JSON"
    }

    previewStyle := lipgloss.NewStyle().
        Width(width).
        Height(m.height - 5).
        BorderStyle(lipgloss.NormalBorder()).
        BorderLeft(true)

    return previewStyle.Render(highlighted)
}
```

**Testing:**
- JSON is syntax highlighted
- Colors are visible in terminal
- Split-view toggle works
- Preview updates on edit

#### Step 3: Template Wizard (1.5 hours)

**File**: `pkg/omise/screens/json_editor/wizard.go`

**Requirements:**
- Show Huh form for complex templates
- Pre-fill with sensible defaults
- Insert generated JSON after form submission
- Support HTTP, Loop, and Conditional templates

**Implementation:**

```go
package json_editor

import "github.com/charmbracelet/huh"

func (m *Model) showHTTPWizard() tea.Cmd {
    var (
        name   string
        method string
        url    string
    )

    form := huh.NewForm(
        huh.NewGroup(
            huh.NewInput().
                Title("Node Name").
                Value(&name).
                Placeholder("My HTTP Request"),

            huh.NewSelect[string]().
                Title("HTTP Method").
                Options(
                    huh.NewOption("GET", "GET"),
                    huh.NewOption("POST", "POST"),
                    huh.NewOption("PUT", "PUT"),
                    huh.NewOption("DELETE", "DELETE"),
                ).
                Value(&method),

            huh.NewInput().
                Title("URL").
                Value(&url).
                Placeholder("https://api.example.com"),
        ),
    )

    return func() tea.Msg {
        if err := form.Run(); err != nil {
            return wizardCancelMsg{}
        }

        // Generate JSON
        nodeID := fmt.Sprintf("node-%d", m.nextNodeID)
        m.nextNodeID++

        json := fmt.Sprintf(`{
  "id": "%s",
  "type": "http",
  "name": "%s",
  "parameters": {
    "method": "%s",
    "url": "%s"
  }
}`, nodeID, name, method, url)

        return wizardCompleteMsg{json: json}
    }
}
```

**Testing:**
- Wizard form displays
- Values are captured
- JSON is generated correctly
- JSON is inserted at cursor

#### Step 4: Improved Error Display (30 min)

**File**: `pkg/omise/screens/json_editor/validation.go`

**Requirements:**
- Show line number of error (if available)
- Display error context (surrounding lines)
- Support multiple errors (show all)

**Implementation:**

```go
type ValidationError struct {
    Line    int
    Column  int
    Message string
}

func (m *Model) validate() []ValidationError {
    content := m.textarea.Value()
    var errors []ValidationError

    // Parse JSON
    var def neta.Definition
    if err := json.Unmarshal([]byte(content), &def); err != nil {
        // Extract line number from JSON syntax error
        if jsonErr, ok := err.(*json.SyntaxError); ok {
            line := calculateLineNumber(content, int(jsonErr.Offset))
            errors = append(errors, ValidationError{
                Line:    line,
                Message: "JSON syntax error: " + err.Error(),
            })
        } else {
            errors = append(errors, ValidationError{
                Message: "JSON error: " + err.Error(),
            })
        }
        return errors
    }

    // Validate structure
    if err := m.validator.ValidateRecursive(def); err != nil {
        errors = append(errors, ValidationError{
            Message: err.Error(),
        })
    }

    return errors
}

func calculateLineNumber(content string, offset int) int {
    return strings.Count(content[:offset], "\n") + 1
}
```

**Testing:**
- Syntax errors show line number
- Multiple errors are displayed
- Error context is helpful

---

## Phase 3C: YAML → JSON Migration (4 hours)

### Deliverables

1. Replace YAML parser with JSON parser
2. Remove YAML dependencies
3. Convert all example files to JSON
4. Update tests for JSON format only
5. Documentation updates
6. Clean up YAML remnants

### Implementation Steps

#### Step 1: Replace YAML Parser with JSON Parser (1.5 hours)

**File**: `pkg/jubako/parser.go`

**Requirements:**
- Replace YAML imports with JSON (stdlib)
- Simplify ParseBytes to only handle JSON
- Remove `normalizeDefinition()` (YAML-specific nodes extraction)
- Keep ID and edge generation (works for both formats)

**Implementation:**

```go
package jubako

import (
    "encoding/json"
    "fmt"
    "os"
    "bento/pkg/neta"
)

// Parser handles .bento.json file parsing.
type Parser struct{}

// NewParser creates a new parser.
func NewParser() *Parser {
    return &Parser{}
}

// Parse reads and parses a .bento.json file.
func (p *Parser) Parse(path string) (neta.Definition, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return neta.Definition{}, fmt.Errorf("read failed: %w", err)
    }

    return p.ParseBytes(data)
}

// ParseBytes parses .bento.json from bytes.
func (p *Parser) ParseBytes(data []byte) (neta.Definition, error) {
    var def neta.Definition
    if err := json.Unmarshal(data, &def); err != nil {
        return neta.Definition{}, fmt.Errorf("invalid JSON: %w", err)
    }

    // Assign IDs to nodes that don't have them
    def = assignNodeIDs(def)

    // Auto-generate edges if missing (backward compatibility)
    def = autoGenerateEdges(def)

    if err := validateDefinition(def); err != nil {
        return neta.Definition{}, fmt.Errorf("validation failed: %w", err)
    }

    return def, nil
}

// Format converts a definition to JSON.
func (p *Parser) Format(def neta.Definition) ([]byte, error) {
    data, err := json.MarshalIndent(def, "", "  ")
    if err != nil {
        return nil, fmt.Errorf("marshal failed: %w", err)
    }
    return data, nil
}

// validateDefinition, assignNodeIDs, autoGenerateEdges remain unchanged
```

**Files to also modify:**
- Remove `import "gopkg.in/yaml.v3"`
- Remove `normalizeDefinition()` function (was YAML-specific)
- Remove `isGroupType()` and `extractNodesFromParams()` (YAML-specific)

**Testing:**
- Parse valid JSON file → success
- Invalid JSON → clear error message
- IDs auto-assigned
- Edges auto-generated for sequences

#### Step 2: Update Discovery (30 min)

**File**: `pkg/jubako/discovery.go`

**Requirements:**
- Discover only `*.bento.json` files
- Simplify glob pattern

**Implementation:**

```go
func DiscoverBentos(dir string) ([]BentoMetadata, error) {
    // Find JSON files only
    pattern := filepath.Join(dir, "*.bento.json")
    jsonFiles, err := filepath.Glob(pattern)
    if err != nil {
        return nil, fmt.Errorf("glob failed: %w", err)
    }

    // Parse all files
    var bentos []BentoMetadata
    for _, path := range jsonFiles {
        meta, err := parseBentoMetadata(path)
        if err != nil {
            // Skip invalid files but log
            continue
        }
        bentos = append(bentos, meta)
    }

    return bentos, nil
}
```

**Testing:**
- Directory with JSON files → all discovered
- Empty directory → empty list
- Invalid JSON files → skipped gracefully

#### Step 3: Update Store (30 min)

**File**: `pkg/jubako/store.go`

**Requirements:**
- Save only in JSON format
- Use `.bento.json` extension

**Implementation:**

```go
func (s *Store) Save(def neta.Definition) error {
    // Marshal to JSON
    data, err := json.MarshalIndent(def, "", "  ")
    if err != nil {
        return fmt.Errorf("marshal failed: %w", err)
    }

    // Build path with .bento.json extension
    path := filepath.Join(s.dir, def.Name+".bento.json")

    // Write file
    if err := os.WriteFile(path, data, 0644); err != nil {
        return fmt.Errorf("write failed: %w", err)
    }

    return nil
}
```

**Testing:**
- Save bento → creates `.bento.json` file
- File contains formatted JSON (2 spaces)
- Existing method calls still work

#### Step 4: Convert Examples (1 hour)

**Files to convert:**

```
examples/http-get.bento.yaml       → examples/http-get.bento.json
examples/group-sequence.bento.yaml → examples/group-sequence.bento.json
examples/loop-for.bento.yaml       → examples/loop-for.bento.json
examples/conditional-if.bento.yaml → examples/conditional-if.bento.json
(8 files total)
```

**Conversion script:**

```bash
#!/bin/bash
# convert-examples.sh

for file in examples/*.bento.yaml; do
    base=$(basename "$file" .bento.yaml)
    echo "Converting $file..."

    # Use Go to convert (preserves types)
    go run ./tools/convert-yaml-to-json.go "$file" "examples/${base}.bento.json"
done

echo "Conversion complete!"
```

**Conversion tool** (`tools/convert-yaml-to-json.go`):

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"
    "bento/pkg/jubako"
)

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: convert-yaml-to-json <input.yaml> <output.json>")
        os.Exit(1)
    }

    inputPath := os.Args[1]
    outputPath := os.Args[2]

    // Parse YAML
    parser := jubako.NewParser()
    def, err := parser.Parse(inputPath)
    if err != nil {
        fmt.Printf("Error parsing YAML: %v\n", err)
        os.Exit(1)
    }

    // Write JSON
    jsonData, err := json.MarshalIndent(def, "", "  ")
    if err != nil {
        fmt.Printf("Error marshaling JSON: %v\n", err)
        os.Exit(1)
    }

    if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
        fmt.Printf("Error writing file: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Converted %s → %s\n", inputPath, outputPath)
}
```

**Testing:**
- Run conversion script
- Verify JSON files are valid
- Compare YAML and JSON (same structure)
- Keep YAML files for backward compatibility

#### Step 5: Update Tests (1 hour)

**Files to update:**

1. `pkg/jubako/parser_test.go`
2. `pkg/jubako/discovery_test.go`
3. `pkg/jubako/store_test.go`

**Update Existing Tests:**

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

    def, err := p.ParseBytes([]byte(jsonData))
    assert.NoError(t, err)
    assert.Equal(t, "1.0", def.Version)
    assert.Equal(t, "http", def.Type)
    assert.Equal(t, "Test", def.Name)
}

func TestParser_InvalidJSON(t *testing.T) {
    p := NewParser()

    invalidJSON := `{
        "version": "1.0",
        "type": "http",
        "name": "Test",
        // Missing closing brace
    `

    _, err := p.ParseBytes([]byte(invalidJSON))
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "invalid JSON")
}

func TestDiscovery_JSONFiles(t *testing.T) {
    tmpDir := t.TempDir()

    // Write JSON file
    jsonPath := filepath.Join(tmpDir, "test.bento.json")
    os.WriteFile(jsonPath, []byte(`{"version":"1.0","type":"http","name":"Test","parameters":{"method":"GET","url":"https://example.com"}}`), 0644)

    // Discover
    bentos, err := DiscoverBentos(tmpDir)
    assert.NoError(t, err)
    assert.Len(t, bentos, 1)
    assert.Contains(t, bentos[0].Path, ".bento.json")
}

func TestStore_SaveJSON(t *testing.T) {
    tmpDir := t.TempDir()
    store := NewStore(tmpDir)

    def := neta.Definition{
        Version: "1.0",
        Type:    "http",
        Name:    "test-bento",
        Parameters: map[string]interface{}{
            "method": "GET",
            "url":    "https://example.com",
        },
    }

    err := store.Save(def)
    assert.NoError(t, err)

    // Verify file exists
    path := filepath.Join(tmpDir, "test-bento.bento.json")
    assert.FileExists(t, path)

    // Verify content is valid JSON
    data, _ := os.ReadFile(path)
    var parsed neta.Definition
    err = json.Unmarshal(data, &parsed)
    assert.NoError(t, err)
}
```

**Remove:**
- All YAML-specific tests
- Format detection tests (no longer needed)
- YAML normalization tests

**Testing:**
- All tests pass with JSON-only parser
- No YAML dependencies remain

#### Step 6: Documentation (30 min)

**Files to update:**

1. `README.md` - Add JSON format examples
2. `.claude/strategy/phase-3-json-editor.md` - This file
3. `.claude/strategy/README.md` - Update phase 3 status

**README.md Updates:**

```markdown
## Bento File Format

Bento uses JSON format for workflow definitions:

### JSON Format

```json
{
  "version": "1.0",
  "name": "example-workflow",
  "description": "API data pipeline",
  "nodes": [
    {
      "id": "node-1",
      "type": "http",
      "name": "Fetch Users",
      "parameters": {
        "method": "GET",
        "url": "https://api.example.com/users"
      }
    }
  ]
}
```

Save as `example.bento.json`

### Features

- Clean, readable syntax
- Standard library support (encoding/json)
- Easy to edit manually or with the built-in JSON editor
- Syntax highlighting and validation
```

**Testing:**
- Documentation is clear
- Examples are valid
- Links work

---

## Success Criteria

### MVP (Phase 3A)

- [ ] Editor opens from browser (press 'e')
- [ ] Line numbers displayed
- [ ] JSON syntax validation works
- [ ] Errors shown in footer within 500ms of typing
- [ ] Ctrl+F formats JSON
- [ ] Ctrl+S saves file
- [ ] All 7 template hotkeys work
- [ ] Templates insert at cursor
- [ ] All files < 250 lines
- [ ] All functions < 20 lines

### Enhancements (Phase 3B)

- [ ] Syntax highlighting works (Chroma)
- [ ] Split-view toggle works (Tab key)
- [ ] Template wizard displays
- [ ] Error display shows line numbers
- [ ] Chroma dependency added to go.mod

### Migration (Phase 3C)

- [ ] JSON parser parses valid JSON
- [ ] YAML parser still works
- [ ] Auto-detection works by extension
- [ ] Discovery finds both formats
- [ ] Store saves in both formats
- [ ] All examples converted to JSON
- [ ] All tests pass
- [ ] Documentation updated

---

## Testing Strategy

### Unit Tests

1. **Editor Tests** (`editor_test.go`)
   - Create new editor
   - Load JSON content
   - Handle window resize
   - Track modified state

2. **Validation Tests** (`validation_test.go`)
   - Valid JSON → no error
   - Invalid syntax → error
   - Missing version → error
   - Invalid node params → error

3. **Formatting Tests** (`formatting_test.go`)
   - Unformatted JSON → formatted
   - Invalid JSON → error returned
   - 2-space indentation used

4. **Template Tests** (`templates_test.go`)
   - Each template is valid JSON
   - Placeholders replaced correctly
   - IDs are unique

5. **Parser Tests** (`parser_test.go`)
   - JSON parsing works
   - YAML parsing still works
   - Auto-detection works

### Integration Tests

1. **Browser → Editor Flow**
   - Open bento from browser
   - Edit and save
   - Return to browser
   - Browser refreshes

2. **Template Insertion**
   - Insert template via hotkey
   - Validate resulting JSON
   - Save successfully

3. **Error Handling**
   - Invalid JSON → error shown
   - File not found → error shown
   - Disk full → error shown

### Manual Testing

1. **Small Terminal** (80x24)
   - Layout still usable
   - No text overflow
   - Help is visible

2. **Large Terminal** (120x40)
   - Layout fills space
   - No empty gaps
   - Text is readable

3. **Long Files** (1000+ lines)
   - Scrolling works
   - Validation still fast
   - No lag in typing

---

## Rollout Plan

### Week 1: MVP Implementation

- Day 1-2: Package structure + basic editor
- Day 3: Validation + formatting
- Day 4: Templates
- Day 5: Testing + bug fixes

### Week 2: Enhancements

- Day 1-2: Syntax highlighting
- Day 3: Template wizard
- Day 4: Error improvements
- Day 5: Testing

### Week 3: Migration

- Day 1-2: JSON parser
- Day 3: Convert examples + update tests
- Day 4: Documentation
- Day 5: Final testing + release

---

## Risks & Mitigation

### Risk 1: Textarea Performance

**Risk**: Large files (1000+ nodes) may be slow to edit

**Mitigation**:
- Test with large files early
- Add file size warning
- Consider viewport/pagination for very large files

**Likelihood**: Low (most bentos <100 nodes)

### Risk 2: Chroma Integration

**Risk**: Chroma doesn't work well with textarea

**Mitigation**:
- Use split-view (plain + highlighted)
- Make highlighting optional
- Ship MVP without highlighting first

**Likelihood**: Medium (known limitation of textarea)

### Risk 3: User Adoption

**Risk**: Users don't want to switch to JSON

**Mitigation**:
- Support both formats indefinitely
- Make JSON default for new bentos only
- Provide clear migration docs

**Likelihood**: Low (JSON is simpler than YAML)

---

## Open Questions

1. ~~Should we support both JSON and YAML long-term?~~
   - **Decision**: Yes, dual support indefinitely (low cost)

2. ~~Should syntax highlighting be in MVP or enhancement?~~
   - **Decision**: Enhancement (Phase 3B)

3. ~~Should we auto-format on save?~~
   - **Decision**: Yes, but allow canceling if errors

4. ~~What should happen if user tries to save invalid JSON?~~
   - **Decision**: Allow draft saves, show warning

5. Should we add a "New Bento" wizard?
   - **Pending**: Consider for Phase 4

---

## Next Steps

1. ✅ Approve this implementation plan
2. Create git branch: `phase-3/json-editor`
3. Implement Phase 3A (MVP)
4. Run `/code-review` after MVP
5. Get Karen's approval
6. Implement Phase 3B (Enhancements)
7. Implement Phase 3C (Migration)
8. Final code review
9. Merge to main

---

**Implementation Plan Created**: 2025-10-17
**Author**: Claude (Sonnet 4.5)
**Status**: Ready for Development
