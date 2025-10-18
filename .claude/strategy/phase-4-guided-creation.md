# Phase 4: Guided Bento Creation with Huh

## Overview
Replace the JSON editor with a guided form-based creation flow using the `huh` library. This provides a clear, foolproof UX for creating bentos without requiring users to know JSON syntax or bento schema.

## Goals
- Simple, guided workflow for creating bentos
- Zero JSON knowledge required
- Impossible to create invalid bentos
- Fast implementation (afternoon of work)
- Better UX than text editor for 95% of use cases

## Non-Goals
- Advanced bulk editing
- Copy/paste between bentos
- Direct JSON manipulation (can add later as "Advanced" option)

## User Flow

### Creating a New Bento

```
Press 'n' in browser
↓
[Bento Metadata Form]
Name: _______________
Description: ________
Icon: 🍱

↓
[Node Creation Loop]
Add a node:
> HTTP GET Request
  HTTP POST Request
  JQ Transform
  For Loop
  Done (finish bento)

↓ (if HTTP GET selected)
[HTTP GET Configuration]
Node Name: _______________
URL: ____________________
Headers (optional): ______
Query Params (optional): __

↓
Add another node? (y/n)

↓ (after all nodes added)
Nodes will be connected sequentially.
Save bento? (y/n)

↓
Bento saved! ✓
```

### Editing an Existing Bento

```
Press 'e' on selected bento
↓
[Edit Bento Menu]
> Edit Metadata (name, desc, icon)
  Add Node
  Edit Node
  Remove Node
  Reorder Nodes
  Done (save changes)
```

## Implementation Architecture

### Files to Create

```
pkg/omise/screens/guided/
├── guided.go          # Main guided creation orchestrator
├── metadata.go        # Bento metadata form
├── nodes.go           # Node type selection
├── node_http.go       # HTTP node configuration forms
├── node_jq.go         # JQ node configuration forms
├── node_loop.go       # Loop node configuration forms
├── edges.go           # Edge creation/inference
└── guided_test.go     # Tests
```

### Core Types

```go
// GuidedCreator orchestrates the guided creation flow
type GuidedCreator struct {
    store       *jubako.Store
    workDir     string
    definition  *jubako.Definition
    editing     bool  // true if editing existing, false if new
}

// CreateBentoGuided runs the guided creation flow
func CreateBentoGuided(store *jubako.Store, workDir string) (*jubako.Definition, error)

// EditBentoGuided runs the guided editing flow
func EditBentoGuided(store *jubako.Store, def *jubako.Definition) (*jubako.Definition, error)
```

## Forms Design

### 1. Metadata Form

```go
func (g *GuidedCreator) promptMetadata() error {
    form := huh.NewForm(
        huh.NewGroup(
            huh.NewInput().
                Title("Bento Name").
                Value(&g.definition.Name).
                Validate(func(s string) error {
                    if s == "" {
                        return errors.New("name required")
                    }
                    return nil
                }),

            huh.NewText().
                Title("Description").
                Value(&g.definition.Description).
                CharLimit(200),

            huh.NewInput().
                Title("Icon (emoji)").
                Value(&g.definition.Icon).
                Placeholder("🍱"),
        ),
    )

    return form.Run()
}
```

### 2. Node Type Selection

```go
func (g *GuidedCreator) promptNodeType() (string, error) {
    var nodeType string

    form := huh.NewForm(
        huh.NewGroup(
            huh.NewSelect[string]().
                Title("Add a node").
                Options(
                    huh.NewOption("HTTP GET Request", "http.get"),
                    huh.NewOption("HTTP POST Request", "http.post"),
                    huh.NewOption("JQ Transform", "jq.transform"),
                    huh.NewOption("For Loop", "loop.for"),
                    huh.NewOption("Shell Command", "shell.exec"),
                    huh.NewOption("Done (finish)", "done"),
                ).
                Value(&nodeType),
        ),
    )

    if err := form.Run(); err != nil {
        return "", err
    }

    return nodeType, nil
}
```

### 3. HTTP GET Node Configuration

```go
func (g *GuidedCreator) promptHTTPGet() (*jubako.Node, error) {
    node := &jubako.Node{
        Type: "http.get",
        Config: make(map[string]interface{}),
    }

    var name, url, headers, query string

    form := huh.NewForm(
        huh.NewGroup(
            huh.NewInput().
                Title("Node Name").
                Value(&name).
                Validate(required),

            huh.NewInput().
                Title("URL").
                Value(&url).
                Validate(required).
                Placeholder("https://api.example.com/data"),

            huh.NewText().
                Title("Headers (JSON, optional)").
                Value(&headers).
                Placeholder(`{"Authorization": "Bearer token"}`),

            huh.NewText().
                Title("Query Params (JSON, optional)").
                Value(&query).
                Placeholder(`{"page": "1", "limit": "10"}`),
        ),
    )

    if err := form.Run(); err != nil {
        return nil, err
    }

    node.Name = name
    node.Config["url"] = url

    if headers != "" {
        var h map[string]string
        json.Unmarshal([]byte(headers), &h)
        node.Config["headers"] = h
    }

    if query != "" {
        var q map[string]string
        json.Unmarshal([]byte(query), &q)
        node.Config["query"] = q
    }

    return node, nil
}
```

### 4. JQ Transform Configuration

```go
func (g *GuidedCreator) promptJQTransform() (*jubako.Node, error) {
    node := &jubako.Node{
        Type: "jq.transform",
        Config: make(map[string]interface{}),
    }

    var name, filter string

    form := huh.NewForm(
        huh.NewGroup(
            huh.NewInput().
                Title("Node Name").
                Value(&name).
                Validate(required),

            huh.NewText().
                Title("JQ Filter").
                Value(&filter).
                Validate(required).
                Placeholder(".data | map(.id)").
                CharLimit(500),
        ),
    )

    if err := form.Run(); err != nil {
        return nil, err
    }

    node.Name = name
    node.Config["filter"] = filter

    return node, nil
}
```

## Edge Creation Strategy

### Option 1: Auto-Sequential (MVP)
Automatically connect nodes in the order they were created (sequential workflow).

```go
func (g *GuidedCreator) inferEdges() {
    g.definition.Edges = []jubako.Edge{}

    for i := 0; i < len(g.definition.Nodes)-1; i++ {
        edge := jubako.Edge{
            From: g.definition.Nodes[i].ID,
            To:   g.definition.Nodes[i+1].ID,
        }
        g.definition.Edges = append(g.definition.Edges, edge)
    }
}
```

### Option 2: Manual Connection (Future)
After all nodes added, show menu to manually connect:

```
Connect nodes:
  [✓] Fetch Data → Transform
  [ ] Fetch Data → Save
  [✓] Transform → Save
```

**Decision**: Use Option 1 (auto-sequential) for MVP. Add manual connection later if needed.

## Main Flow Implementation

```go
func CreateBentoGuided(store *jubako.Store, workDir string) (*jubako.Definition, error) {
    g := &GuidedCreator{
        store:   store,
        workDir: workDir,
        definition: &jubako.Definition{
            Version: "1.0",
            Type:    "group.sequence",
            Nodes:   []jubako.Node{},
            Edges:   []jubako.Edge{},
        },
    }

    // Step 1: Metadata
    if err := g.promptMetadata(); err != nil {
        return nil, err
    }

    // Step 2: Node creation loop
    for {
        nodeType, err := g.promptNodeType()
        if err != nil {
            return nil, err
        }

        if nodeType == "done" {
            break
        }

        node, err := g.createNode(nodeType)
        if err != nil {
            return nil, err
        }

        // Assign auto ID
        node.ID = fmt.Sprintf("node-%d", len(g.definition.Nodes)+1)
        g.definition.Nodes = append(g.definition.Nodes, *node)
    }

    // Step 3: Auto-create edges
    g.inferEdges()

    // Step 4: Confirm save
    var save bool
    confirmForm := huh.NewForm(
        huh.NewGroup(
            huh.NewConfirm().
                Title("Save bento?").
                Value(&save),
        ),
    )

    if err := confirmForm.Run(); err != nil {
        return nil, err
    }

    if !save {
        return nil, errors.New("cancelled")
    }

    // Save to store
    filename := strings.ReplaceAll(
        strings.ToLower(g.definition.Name),
        " ", "-",
    ) + ".bento.json"

    if err := store.Save(g.definition, filename); err != nil {
        return nil, err
    }

    return g.definition, nil
}

func (g *GuidedCreator) createNode(nodeType string) (*jubako.Node, error) {
    switch nodeType {
    case "http.get":
        return g.promptHTTPGet()
    case "http.post":
        return g.promptHTTPPost()
    case "jq.transform":
        return g.promptJQTransform()
    case "loop.for":
        return g.promptLoop()
    case "shell.exec":
        return g.promptShell()
    default:
        return nil, fmt.Errorf("unknown node type: %s", nodeType)
    }
}
```

## Integration with Browser

### Update browser_handlers.go

```go
// handleNew creates a new bento with guided prompts
func (b Browser) handleNew() (Browser, tea.Cmd) {
    return b, func() tea.Msg {
        // Run guided creation in background
        def, err := guided.CreateBentoGuided(b.store, b.store.WorkDir())
        if err != nil {
            return BentoOperationCompleteMsg{
                Success: false,
                Message: fmt.Sprintf("Failed to create bento: %v", err),
            }
        }

        return BentoOperationCompleteMsg{
            Success: true,
            Message: fmt.Sprintf("Created bento: %s", def.Name),
        }
    }
}
```

## Validation Helpers

```go
func required(s string) error {
    if strings.TrimSpace(s) == "" {
        return errors.New("this field is required")
    }
    return nil
}

func validURL(s string) error {
    if s == "" {
        return errors.New("URL is required")
    }
    _, err := url.Parse(s)
    if err != nil {
        return errors.New("invalid URL format")
    }
    return nil
}

func validJSON(s string) error {
    if s == "" {
        return nil  // optional
    }
    var js map[string]interface{}
    if err := json.Unmarshal([]byte(s), &js); err != nil {
        return errors.New("invalid JSON format")
    }
    return nil
}
```

## Testing Strategy

```go
func TestCreateBentoGuided(t *testing.T) {
    // Mock user input
    // Test metadata prompt
    // Test node creation
    // Test edge inference
    // Test save
}

func TestValidation(t *testing.T) {
    // Test required fields
    // Test URL validation
    // Test JSON validation
}
```

## UI/UX Considerations

### Loading State
While running huh forms, the TUI is blocked. This is fine - users are focused on the form.

### Error Handling
If form returns error (ESC pressed), return to browser without saving.

### Progress Indication
Show step numbers in form titles:
- "Step 1/3: Bento Metadata"
- "Step 2/3: Add Nodes"
- "Step 3/3: Review & Save"

### Form Styling
Use huh's theming to match app colors (pull from styles package).

## Implementation Plan

### Phase 1: Basic Structure (30 min)
1. Create `pkg/omise/screens/guided/` package
2. Implement `GuidedCreator` type
3. Implement metadata form
4. Wire up to browser 'n' key
5. Test: Can create bento with just metadata

### Phase 2: Node Creation (1 hour)
1. Implement node type selection
2. Implement HTTP GET form
3. Implement JQ transform form
4. Add node creation loop
5. Test: Can create bento with nodes

### Phase 3: Edges & Save (30 min)
1. Implement edge inference
2. Implement save confirmation
3. Test: Complete flow creates valid bento

### Phase 4: Additional Node Types (1 hour)
1. Implement HTTP POST form
2. Implement Loop form
3. Implement Shell form
4. Test: All node types work

### Phase 5: Edit Flow (1 hour - optional for MVP)
1. Implement edit menu
2. Implement node editing
3. Implement node removal
4. Wire up to browser 'e' key

**Total MVP Time**: ~3 hours (Phases 1-3)
**With all features**: ~5 hours (Phases 1-4)

## Success Criteria

1. User can create a complete bento without seeing JSON
2. All required fields are validated
3. Invalid bentos cannot be created
4. Created bentos execute correctly
5. UX is clearer than text editor

## Future Enhancements

- Manual edge creation
- Node reordering
- Copy node from another bento
- Import/export templates
- Preview generated JSON (for advanced users)
- Undo/redo in forms

## Risks & Mitigations

**Risk**: Huh forms block the TUI
- Mitigation: This is expected behavior, users are focused on forms

**Risk**: Complex configurations hard in forms (e.g., nested JSON)
- Mitigation: Start simple, add advanced fields later

**Risk**: Users want more control
- Mitigation: Add JSON editor later as "Advanced Edit" option

**Risk**: Form flow feels slow for experts
- Mitigation: Keep forms snappy, add keyboard shortcuts to skip through quickly
