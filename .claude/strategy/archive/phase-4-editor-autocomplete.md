# Phase 4: JSON Editor Autocomplete & Managed Fields

## Overview
Enhance the JSON editor with intelligent autocomplete for bento structures and automatic management of system-level fields that users shouldn't manually edit.

## Problem Statement
Currently, users must manually type all JSON structures including:
- Node type definitions (http.get, jq.transform, etc.)
- Field names and structure
- System fields like `id`, `version`, and `type` that should be managed automatically

This is error-prone and provides poor UX compared to modern code editors.

## Goals
1. **Managed Fields**: Automatically handle `id`, `version`, and top-level `type` fields
2. **Node Autocomplete**: Suggest node templates when creating nodes in the `nodes` array
3. **Field Autocomplete**: Suggest valid field names based on current context
4. **Type Value Autocomplete**: Suggest valid values for `type` fields within nodes

## Non-Goals (for this phase)
- Full LSP-style autocomplete for all JSON
- Autocomplete for arbitrary user data
- Inline documentation/help text (can add later)

## Managed Fields Strategy

### Fields to Auto-Manage
1. **`version`** - Always set to `"1.0"`
2. **`type`** - Always set to `"group.sequence"` at root level
3. **`id`** - Generate UUIDs for nodes automatically

### Implementation Approach

#### Option 1: Hidden Fields (Recommended)
- Strip managed fields from editor content when loading
- Add them back when saving
- Simpler UX - user never sees or worries about these fields

```go
// When loading into editor
func prepareForEdit(def *jubako.Definition) string {
    // Remove version, type, and node IDs
    cleaned := removeAutoFields(def)
    return marshal(cleaned)
}

// When saving from editor
func prepareForSave(content string) (*jubako.Definition, error) {
    def := parse(content)
    // Add version and type
    def.Version = "1.0"
    def.Type = "group.sequence"
    // Generate IDs for nodes that don't have them
    assignNodeIDs(def.Nodes)
    return def
}
```

#### Option 2: Read-only Fields
- Show fields but make them non-editable
- More complex to implement in textarea
- User sees fields but might be confused why they can't edit

**Decision**: Use Option 1 (Hidden Fields) for better UX.

### Template Structure
Users will edit a simplified structure:
```json
{
  "name": "My Workflow",
  "icon": "🍱",
  "description": "Does cool things",
  "nodes": [
    {
      "type": "http.get",
      "name": "Fetch Data",
      "url": "https://api.example.com"
    }
  ],
  "edges": []
}
```

System adds on save:
```json
{
  "version": "1.0",
  "type": "group.sequence",
  "name": "My Workflow",
  "icon": "🍱",
  "description": "Does cool things",
  "nodes": [
    {
      "id": "node-550e8400-e29b",
      "type": "http.get",
      "name": "Fetch Data",
      "url": "https://api.example.com"
    }
  ],
  "edges": []
}
```

## Autocomplete System Design

### Architecture Components

#### 1. Context Analyzer
Determines where cursor is in JSON structure:
```go
type JSONContext struct {
    Path       []string      // e.g., ["nodes", "0", "type"]
    InArray    bool          // Inside an array
    ArrayType  string        // "nodes", "edges", etc.
    CurrentKey string        // Current field being edited
    Expecting  ExpectedType  // What comes next
}

type ExpectedType int
const (
    ExpectKey ExpectedType = iota
    ExpectValue
    ExpectArrayElement
)
```

#### 2. Suggestion Provider
Provides autocomplete suggestions based on context:
```go
type Suggestion struct {
    Label       string   // Display text
    InsertText  string   // Text to insert
    Detail      string   // Description
    Kind        SuggestKind
}

type SuggestKind int
const (
    SuggestNodeType SuggestKind = iota
    SuggestField
    SuggestValue
    SuggestTemplate
)

type SuggestionProvider interface {
    GetSuggestions(ctx JSONContext) []Suggestion
}
```

#### 3. Autocomplete UI Component
Popup menu overlaying the textarea:
```go
type AutocompleteMenu struct {
    suggestions []Suggestion
    selected    int
    visible     bool
    x, y        int  // Position in terminal
}
```

### Autocomplete Triggers

#### 1. Node Template Autocomplete
**Trigger**: Type `{` inside `nodes` array
**Context**: `["nodes", "<index>"]`

**Suggestions**:
- HTTP GET Request
- HTTP POST Request
- JQ Transform
- For Loop
- Conditional
- Shell Command
- File Read/Write

**Insert Template Example** (HTTP GET):
```json
{
  "type": "http.get",
  "name": "HTTP Request",
  "url": "https://api.example.com",
  "headers": {},
  "query": {}
}
```

#### 2. Field Name Autocomplete
**Trigger**: Type `"` for a new key inside an object
**Context**: Detect we're starting a new field

**Suggestions** (context-dependent):
- Inside node: `type`, `name`, `url`, `headers`, `body`, `method`, etc.
- At root: `name`, `icon`, `description`, `nodes`, `edges`

#### 3. Type Value Autocomplete
**Trigger**: Editing value of `"type"` field
**Context**: `currentKey == "type"`

**Suggestions**:
- `http.get`
- `http.post`
- `jq.transform`
- `loop.for`
- `conditional.if`
- `shell.exec`
- `file.read`
- `file.write`

### Trigger Detection Strategy

#### Option 1: Character-based Triggers
Watch for specific characters:
- `{` → Node template
- `"` → Field name or value
- `.` → Type continuation (e.g., `http.`)

#### Option 2: Hotkey Trigger
User presses `Ctrl+Space` to invoke autocomplete

**Decision**: Hybrid approach
- Automatic for `{` in nodes array (most common use case)
- `Ctrl+Space` for manual invocation anywhere else
- ESC to close autocomplete menu

### User Workflow Examples

#### Creating a Node
1. User types `{` in nodes array
2. Autocomplete menu appears with node types
3. User arrows down to "HTTP GET Request", presses Enter
4. Full template inserted with cursor on `name` field
5. User fills in name, tabs to next field

#### Editing Type Field
1. User positions cursor on type value
2. Presses `Ctrl+Space`
3. Autocomplete shows valid type values
4. User selects from list

## Implementation Plan

### Phase 4.1: Managed Fields
**Files to modify:**
- `pkg/omise/screens/json_editor/save.go` - Strip/add fields on save
- `pkg/omise/screens/browser_handlers.go` - Strip fields when loading into editor
- `pkg/jubako/types.go` - Helper functions for field management

**Tasks:**
1. Create `prepareForEdit()` function to remove managed fields
2. Create `prepareForSave()` function to add managed fields back
3. Implement UUID generation for node IDs
4. Update editor load/save flow to use these functions
5. Test that bentos round-trip correctly

### Phase 4.2: Context Analysis
**New files:**
- `pkg/omise/screens/json_editor/context.go` - JSON context analyzer

**Tasks:**
1. Implement cursor position → JSON path mapping
2. Detect array context (nodes vs edges vs other)
3. Detect key vs value context
4. Create unit tests for context detection

### Phase 4.3: Suggestion System
**New files:**
- `pkg/omise/screens/json_editor/suggestions.go` - Suggestion provider
- `pkg/omise/screens/json_editor/templates.go` - Node templates (can reuse existing)

**Tasks:**
1. Define node type templates (reuse existing template system)
2. Implement suggestion provider for each trigger type
3. Create suggestion ranking/filtering logic

### Phase 4.4: Autocomplete UI
**New files:**
- `pkg/omise/screens/json_editor/autocomplete.go` - Autocomplete menu component

**Tasks:**
1. Build popup menu component
2. Integrate with editor Update() cycle
3. Handle arrow keys and Enter for selection
4. Handle ESC to cancel
5. Position menu relative to cursor

### Phase 4.5: Integration & Testing
**Tasks:**
1. Wire up autocomplete triggers to context analyzer
2. Add `Ctrl+Space` hotkey handling
3. Test all autocomplete scenarios
4. Update integration tests

## Technical Considerations

### Cursor Position Mapping
Textarea provides cursor position as character offset. Need to:
1. Parse JSON up to cursor position
2. Build JSON path from parsed structure
3. Handle partially-typed JSON gracefully

### Partial JSON Parsing
Standard JSON parsers fail on incomplete JSON. Options:
- Use lenient parser that handles trailing commas, missing brackets
- Parse line-by-line and infer context from indentation
- Use JSON5 or similar relaxed format

**Decision**: Implement simple line-based context detection:
- Look at current line and previous lines
- Use indentation level to infer array/object depth
- Check for `"nodes": [` or `"edges": [` in parent lines

### Performance
- Context analysis on every keystroke could be slow
- Debounce analysis by 100-200ms
- Cache parsed context between keystrokes

### Alternative: Structured Editor (Future)
Instead of text editing with autocomplete, use a structured form-based editor:
- Node list with add/remove buttons
- Each node expands to show fields
- Field types enforce valid values

**Decision**: Stick with text editor + autocomplete for now. More flexible for advanced users.

## Open Questions

1. **Should edges get autocomplete too?**
   - Edges reference node IDs, which we're hiding
   - May need different UX for edges (dropdown of existing nodes?)

2. **Icon picker?**
   - Autocomplete with emoji list when editing `icon` field?
   - Or simple text field is fine?

3. **Field validation?**
   - Should autocomplete only suggest valid fields?
   - Or allow custom fields for extensibility?

4. **Multi-cursor support?**
   - Out of scope for now

## Success Criteria

1. User can create a new bento without manually typing `version`, `type`, or node `id`
2. User can type `{` in nodes array and get node type suggestions
3. User can press `Ctrl+Space` on any field to see valid completions
4. Saved bentos have all managed fields properly populated
5. No breaking changes to existing bento file format

## Risks & Mitigations

**Risk**: Autocomplete feels slow or janky
- Mitigation: Debounce, optimize context detection, keep suggestion lists small

**Risk**: Hiding fields confuses users who need to see/debug them
- Mitigation: Add "Show Advanced Fields" toggle in normal mode (future enhancement)

**Risk**: Context detection fails on complex/nested JSON
- Mitigation: Graceful degradation - if context unclear, don't show autocomplete

**Risk**: Users want to edit managed fields (e.g., manual node IDs)
- Mitigation: Document that IDs are auto-managed, provide export/import for advanced users

## Timeline Estimate

- Phase 4.1 (Managed Fields): 2-3 hours
- Phase 4.2 (Context Analysis): 3-4 hours
- Phase 4.3 (Suggestions): 2-3 hours
- Phase 4.4 (UI Component): 4-5 hours
- Phase 4.5 (Integration): 2-3 hours
- **Total**: ~15-20 hours

## Next Steps

1. Review this design doc
2. Decide on any open questions
3. Start with Phase 4.1 (Managed Fields) as it's independent and immediately useful
4. Validate managed field approach before building autocomplete on top
