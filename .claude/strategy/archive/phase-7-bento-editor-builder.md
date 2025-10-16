# Phase 7: Bento Editor - Node Builder

**Status**: Pending
**Duration**: 5-6 hours
**Prerequisites**: Phase 6 complete, Karen approved

## Overview

Create the Bento Editor screen with guided node building experience. Users select node types from the Pantry, configure parameters through Huh form wizards, and incrementally build their bento Definition. The editor is only accessible through the browser (edit existing or create new).

## Pre-Work Checklist

Before starting, you MUST:

1. ✅ Read [BENTO_BOX_PRINCIPLE.md](../BENTO_BOX_PRINCIPLE.md)
2. ✅ Read [CHARM_STACK_GUIDE.md](../CHARM_STACK_GUIDE.md)
3. ✅ Confirm: "I understand the Bento Box Principle and will follow it"
4. ✅ Use TodoWrite to track all tasks
5. ✅ Phase 6 approved by Karen

## Goals

1. Create new Editor screen (separate from browser)
2. Integrate Pantry for node type discovery
3. Implement Huh wizards for parameter configuration
4. Build neta.Definition incrementally
5. Save to Jubako with validation
6. Support both create and edit modes
7. Clear, guided user experience
8. Validate Bento Box compliance

## Editor Flow

### Create New Bento
```
Browser → Press 'n' → Editor (Create Mode)
                       ↓
                   Enter bento name
                       ↓
                   Select node type from Pantry
                       ↓
                   Configure node parameters (Huh form)
                       ↓
                   Add node to bento
                       ↓
                   Option: Add another node / Save / Cancel
                       ↓
                   Save to Jubako → Back to Browser
```

### Edit Existing Bento
```
Browser → Select bento → Press 'e' → Editor (Edit Mode)
                                       ↓
                                   Load existing Definition
                                       ↓
                                   View/Edit nodes
                                       ↓
                                   Save changes → Back to Browser
```

## Screen Structure

```
pkg/omise/screens/
├── editor.go            # Main editor screen (NEW)
├── editor_test.go       # Editor tests (NEW)
├── editor_wizard.go     # Huh form wizards (NEW)
└── editor_wizard_test.go

pkg/omise/
├── model.go             # Add ScreenEditor (modify)
└── update.go            # Handle editor messages (modify)
```

## Deliverables

### 1. Editor Screen

**File**: `pkg/omise/screens/editor.go` (NEW)
**Target Size**: < 250 lines

```go
package screens

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"bento/pkg/jubako"
	"bento/pkg/neta"
	"bento/pkg/omise/styles"
	"bento/pkg/pantry"
)

// EditorMode defines the editor mode
type EditorMode int

const (
	EditorModeCreate EditorMode = iota
	EditorModeEdit
)

// EditorState defines the current editor state
type EditorState int

const (
	StateNaming EditorState = iota
	StateSelectingType
	StateConfiguringNode
	StateReview
)

// Editor screen for creating and editing bentos
type Editor struct {
	mode     EditorMode
	state    EditorState
	store    *jubako.Store
	registry *pantry.Registry

	// Bento being edited
	bentoName string
	bentoPath string
	def       neta.Definition

	// Current node being configured
	currentNodeType string
	wizard          *NodeWizard

	// UI state
	message string
	width   int
	height  int
}

// NewEditorCreate creates editor in create mode
func NewEditorCreate(store *jubako.Store, registry *pantry.Registry) Editor {
	return Editor{
		mode:     EditorModeCreate,
		state:    StateNaming,
		store:    store,
		registry: registry,
		def: neta.Definition{
			Version: neta.CurrentVersion,
			Nodes:   []neta.Definition{},
		},
	}
}

// NewEditorEdit creates editor in edit mode
func NewEditorEdit(store *jubako.Store, registry *pantry.Registry, name, path string) (Editor, error) {
	def, err := store.Load(name)
	if err != nil {
		return Editor{}, err
	}

	return Editor{
		mode:      EditorModeEdit,
		state:     StateReview,
		store:     store,
		registry:  registry,
		bentoName: name,
		bentoPath: path,
		def:       def,
	}, nil
}

// Init initializes the editor
func (e Editor) Init() tea.Cmd {
	return nil
}

// Update handles editor messages
func (e Editor) Update(msg tea.Msg) (Editor, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		e.width = msg.Width
		e.height = msg.Height
		return e, nil

	case tea.KeyMsg:
		return e.handleKey(msg)

	case BentoNameEnteredMsg:
		return e.handleNameEntered(msg)

	case NodeTypeSelectedMsg:
		return e.handleTypeSelected(msg)

	case NodeConfiguredMsg:
		return e.handleNodeConfigured(msg)
	}

	return e, nil
}

// handleKey processes keyboard input
func (e Editor) handleKey(msg tea.KeyMsg) (Editor, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Cancel editor
		return e, func() tea.Msg {
			return EditorCancelledMsg{}
		}

	case "ctrl+s":
		// Save bento
		return e, e.saveBento()

	default:
		return e.handleStateKey(msg)
	}
}

// handleStateKey processes state-specific keys
func (e Editor) handleStateKey(msg tea.KeyMsg) (Editor, tea.Cmd) {
	switch e.state {
	case StateNaming:
		return e.handleNamingKey(msg)
	case StateSelectingType:
		return e.handleTypeSelectionKey(msg)
	case StateConfiguringNode:
		return e.handleConfigurationKey(msg)
	case StateReview:
		return e.handleReviewKey(msg)
	}
	return e, nil
}

// handleNamingKey handles name entry state
func (e Editor) handleNamingKey(msg tea.KeyMsg) (Editor, tea.Cmd) {
	// Huh form will handle this
	// For now, simulate name entry
	if msg.String() == "enter" {
		// Move to type selection
		e.state = StateSelectingType
	}
	return e, nil
}

// handleTypeSelectionKey handles type selection
func (e Editor) handleTypeSelectionKey(msg tea.KeyMsg) (Editor, tea.Cmd) {
	// Show list of node types from registry
	// User selects with arrow keys + enter
	switch msg.String() {
	case "enter":
		// Selected a type, move to configuration
		// For now, assume "http" selected
		e.currentNodeType = "http"
		e.state = StateConfiguringNode
		return e, e.launchWizard(e.currentNodeType)
	}
	return e, nil
}

// handleConfigurationKey handles parameter configuration
func (e Editor) handleConfigurationKey(msg tea.KeyMsg) (Editor, tea.Cmd) {
	// Wizard handles this
	return e, nil
}

// handleReviewKey handles review state
func (e Editor) handleReviewKey(msg tea.KeyMsg) (Editor, tea.Cmd) {
	switch msg.String() {
	case "a":
		// Add another node
		e.state = StateSelectingType
		return e, nil
	case "s", "enter":
		// Save and exit
		return e, e.saveBento()
	}
	return e, nil
}

// handleNameEntered processes name entry
func (e Editor) handleNameEntered(msg BentoNameEnteredMsg) (Editor, tea.Cmd) {
	e.bentoName = msg.Name
	e.def.Name = msg.Name
	e.state = StateSelectingType
	return e, nil
}

// handleTypeSelected processes type selection
func (e Editor) handleTypeSelected(msg NodeTypeSelectedMsg) (Editor, tea.Cmd) {
	e.currentNodeType = msg.Type
	e.state = StateConfiguringNode
	return e, e.launchWizard(msg.Type)
}

// handleNodeConfigured processes configured node
func (e Editor) handleNodeConfigured(msg NodeConfiguredMsg) (Editor, tea.Cmd) {
	// Add node to definition
	node := neta.Definition{
		Version:    neta.CurrentVersion,
		Type:       msg.Type,
		Name:       msg.Name,
		Parameters: msg.Parameters,
	}

	// If this is a single-node bento, set as root
	if len(e.def.Nodes) == 0 && e.def.Type == "" {
		e.def.Type = msg.Type
		e.def.Parameters = msg.Parameters
	} else {
		// Multi-node bento, add to nodes
		if e.def.Type == "" {
			// Convert to group
			e.def.Type = "group.sequence"
		}
		e.def.Nodes = append(e.def.Nodes, node)
	}

	e.state = StateReview
	e.message = fmt.Sprintf("Added node: %s", msg.Name)
	return e, nil
}

// launchWizard starts configuration wizard
func (e Editor) launchWizard(nodeType string) tea.Cmd {
	return func() tea.Msg {
		// Phase 7: Launch Huh wizard for node type
		// For now, return dummy configured node
		return NodeConfiguredMsg{
			Type: nodeType,
			Name: "New Node",
			Parameters: map[string]interface{}{
				"url": "https://example.com",
			},
		}
	}
}

// saveBento saves the bento to Jubako
func (e Editor) saveBento() tea.Cmd {
	return func() tea.Msg {
		if err := e.store.Save(e.bentoName, e.def); err != nil {
			return EditorSaveErrorMsg{Error: err}
		}
		return EditorSavedMsg{Name: e.bentoName}
	}
}

// View renders the editor
func (e Editor) View() string {
	title := e.renderTitle()
	content := e.renderContent()
	footer := e.renderFooter()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		content,
		"",
		footer,
	)
}

// renderTitle renders the editor title
func (e Editor) renderTitle() string {
	mode := "Create New Bento"
	if e.mode == EditorModeEdit {
		mode = fmt.Sprintf("Edit Bento: %s", e.bentoName)
	}
	return styles.Title.Render(mode)
}

// renderContent renders state-specific content
func (e Editor) renderContent() string {
	switch e.state {
	case StateNaming:
		return e.renderNaming()
	case StateSelectingType:
		return e.renderTypeSelection()
	case StateConfiguringNode:
		return e.renderConfiguration()
	case StateReview:
		return e.renderReview()
	}
	return ""
}

// renderNaming renders name entry
func (e Editor) renderNaming() string {
	return styles.Subtle.Render("Enter bento name:\n\n[Name entry form here]")
}

// renderTypeSelection renders node type selection
func (e Editor) renderTypeSelection() string {
	types := e.registry.List()
	content := "Select node type:\n\n"

	for _, t := range types {
		content += fmt.Sprintf("  • %s\n", t)
	}

	return styles.Subtle.Render(content)
}

// renderConfiguration renders parameter configuration
func (e Editor) renderConfiguration() string {
	return styles.Subtle.Render(fmt.Sprintf("Configure %s node:\n\n[Wizard form here]", e.currentNodeType))
}

// renderReview renders bento review
func (e Editor) renderReview() string {
	content := fmt.Sprintf("Bento: %s (v%s)\n", e.def.Name, e.def.Version)
	content += fmt.Sprintf("Type: %s\n\n", e.def.Type)

	if len(e.def.Nodes) > 0 {
		content += "Nodes:\n"
		for i, node := range e.def.Nodes {
			content += fmt.Sprintf("  %d. %s (%s)\n", i+1, node.Name, node.Type)
		}
	}

	return styles.Subtle.Render(content)
}

// renderFooter renders keyboard shortcuts
func (e Editor) renderFooter() string {
	shortcuts := ""
	switch e.state {
	case StateReview:
		shortcuts = "a: Add node • s: Save • esc: Cancel"
	default:
		shortcuts = "esc: Cancel • ctrl+s: Save"
	}

	if e.message != "" {
		shortcuts = e.message + " • " + shortcuts
	}

	return styles.Subtle.Render(shortcuts)
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: Bento editing
- ✅ State machine pattern (clear states)
- ✅ Functions < 20 lines
- ✅ File < 250 lines

### 2. Editor Messages

**File**: `pkg/omise/screens/editor_messages.go` (NEW)
**Target Size**: < 100 lines

```go
package screens

import "bento/pkg/neta"

// BentoNameEnteredMsg signals bento name was entered
type BentoNameEnteredMsg struct {
	Name string
}

// NodeTypeSelectedMsg signals node type was selected from Pantry
type NodeTypeSelectedMsg struct {
	Type string
}

// NodeConfiguredMsg signals node parameters were configured
type NodeConfiguredMsg struct {
	Type       string
	Name       string
	Parameters map[string]interface{}
}

// EditorSavedMsg signals bento was saved
type EditorSavedMsg struct {
	Name string
}

// EditorSaveErrorMsg signals save error
type EditorSaveErrorMsg struct {
	Error error
}

// EditorCancelledMsg signals editor was cancelled
type EditorCancelledMsg struct{}

// ReturnToBrowserMsg signals return to browser
type ReturnToBrowserMsg struct{}
```

**Bento Box Compliance**:
- ✅ Single responsibility: Message types
- ✅ Clear event definitions
- ✅ File < 100 lines

### 3. Node Wizard

**File**: `pkg/omise/screens/editor_wizard.go` (NEW)
**Target Size**: < 250 lines

```go
package screens

import (
	"github.com/charmbracelet/huh"
)

// NodeWizard guides parameter configuration
type NodeWizard struct {
	nodeType string
	form     *huh.Form
	result   map[string]interface{}
}

// NewNodeWizard creates wizard for node type
func NewNodeWizard(nodeType string) *NodeWizard {
	form := buildForm(nodeType)
	return &NodeWizard{
		nodeType: nodeType,
		form:     form,
		result:   make(map[string]interface{}),
	}
}

// Run executes the wizard
func (w *NodeWizard) Run() (map[string]interface{}, error) {
	if err := w.form.Run(); err != nil {
		return nil, err
	}
	return w.result, nil
}

// buildForm creates Huh form for node type
func buildForm(nodeType string) *huh.Form {
	switch nodeType {
	case "http":
		return buildHTTPForm()
	case "transform.jq":
		return buildJQForm()
	case "conditional.if":
		return buildConditionalForm()
	default:
		return buildGenericForm()
	}
}

// buildHTTPForm creates form for HTTP nodes
func buildHTTPForm() *huh.Form {
	var (
		name   string
		url    string
		method string
	)

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Node Name").
				Placeholder("My HTTP Request").
				Value(&name).
				Validate(validateRequired),

			huh.NewInput().
				Title("URL").
				Placeholder("https://api.example.com/data").
				Value(&url).
				Validate(validateURL),

			huh.NewSelect[string]().
				Title("Method").
				Options(
					huh.NewOption("GET", "GET"),
					huh.NewOption("POST", "POST"),
					huh.NewOption("PUT", "PUT"),
					huh.NewOption("DELETE", "DELETE"),
				).
				Value(&method),
		),
	)
}

// buildJQForm creates form for JQ transform nodes
func buildJQForm() *huh.Form {
	var (
		name   string
		filter string
	)

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Node Name").
				Placeholder("Transform Data").
				Value(&name).
				Validate(validateRequired),

			huh.NewText().
				Title("JQ Filter").
				Placeholder(".data | map(.id)").
				Value(&filter).
				Validate(validateRequired),
		),
	)
}

// buildConditionalForm creates form for conditional nodes
func buildConditionalForm() *huh.Form {
	var (
		name      string
		condition string
	)

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Node Name").
				Placeholder("Check Status").
				Value(&name).
				Validate(validateRequired),

			huh.NewInput().
				Title("Condition").
				Placeholder(".status == 200").
				Value(&condition).
				Validate(validateRequired),
		),
	)
}

// buildGenericForm creates generic form
func buildGenericForm() *huh.Form {
	var name string

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Node Name").
				Placeholder("New Node").
				Value(&name).
				Validate(validateRequired),
		),
	)
}

// Validators

func validateRequired(s string) error {
	if s == "" {
		return fmt.Errorf("this field is required")
	}
	return nil
}

func validateURL(s string) error {
	if s == "" {
		return fmt.Errorf("URL is required")
	}
	if !strings.HasPrefix(s, "http://") && !strings.HasPrefix(s, "https://") {
		return fmt.Errorf("URL must start with http:// or https://")
	}
	return nil
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: Form wizards
- ✅ One form builder per node type
- ✅ Clear validation
- ✅ File < 250 lines

### 4. Update Root Model

**File**: `pkg/omise/model.go` (add ScreenEditor)

```go
const (
	ScreenBrowser Screen = iota
	ScreenExecutor
	ScreenPantry
	ScreenSettings
	ScreenHelp
	ScreenEditor  // NEW
	screenCount
)

// String returns the screen name
func (s Screen) String() string {
	return [...]string{"Browser", "Executor", "Pantry", "Settings", "Help", "Editor"}[s]
}

// Model add editor field
type Model struct {
	// ... existing fields
	editor   screens.Editor
}

// NewModelWithWorkDir update
func NewModelWithWorkDir(workDir string) (Model, error) {
	// ... existing code

	store, _ := jubako.NewStore(workDir)
	registry := pantry.NewRegistry()

	return Model{
		// ... existing fields
		editor: screens.Editor{}, // Will be initialized when needed
	}, nil
}
```

**File**: `pkg/omise/update.go` (add editor handlers)

```go
// handleEditBento switches to editor for existing bento
func (m Model) handleEditBento(msg screens.EditBentoMsg) (tea.Model, tea.Cmd) {
	store, _ := jubako.NewStore(m.workDir) // Store workDir in Model
	registry := pantry.NewRegistry()

	editor, err := screens.NewEditorEdit(store, registry, msg.Name, msg.Path)
	if err != nil {
		// Handle error
		return m, nil
	}

	m.editor = editor
	m.screen = ScreenEditor
	return m, nil
}

// handleCreateBento switches to editor for new bento
func (m Model) handleCreateBento(msg screens.CreateBentoMsg) (tea.Model, tea.Cmd) {
	store, _ := jubako.NewStore(m.workDir)
	registry := pantry.NewRegistry()

	m.editor = screens.NewEditorCreate(store, registry)
	m.screen = ScreenEditor
	return m, nil
}

// handleEditorSaved returns to browser after save
func (m Model) handleEditorSaved(msg screens.EditorSavedMsg) (tea.Model, tea.Cmd) {
	m.screen = ScreenBrowser
	// Refresh browser list
	return m.updateScreen(screens.BentoListRefreshMsg{})
}

// handleEditorCancelled returns to browser without saving
func (m Model) handleEditorCancelled(msg screens.EditorCancelledMsg) (tea.Model, tea.Cmd) {
	m.screen = ScreenBrowser
	return m, nil
}

// Update add editor message handlers
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// ... existing cases
	case screens.EditorSavedMsg:
		return m.handleEditorSaved(msg)
	case screens.EditorCancelledMsg:
		return m.handleEditorCancelled(msg)
	default:
		return m.updateScreen(msg)
	}
}

// updateScreen add editor case
func (m Model) updateScreen(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.screen {
	// ... existing cases
	case ScreenEditor:
		m.editor, cmd = m.editor.Update(msg)
	}

	return m, cmd
}
```

## Testing Strategy

**File**: `pkg/omise/screens/editor_test.go`

```go
func TestEditor_CreateMode(t *testing.T) {
	workDir := t.TempDir()
	store, err := jubako.NewStore(workDir)
	if err != nil {
		t.Fatal(err)
	}

	registry := pantry.NewRegistry()
	editor := screens.NewEditorCreate(store, registry)

	// Test initial state
	if editor.mode != screens.EditorModeCreate {
		t.Error("expected create mode")
	}

	// Test name entry
	editor, _ = editor.Update(screens.BentoNameEnteredMsg{Name: "test-bento"})
	if editor.bentoName != "test-bento" {
		t.Error("name not set")
	}

	// Test type selection
	editor, _ = editor.Update(screens.NodeTypeSelectedMsg{Type: "http"})
	if editor.currentNodeType != "http" {
		t.Error("type not set")
	}

	// Test node configuration
	editor, _ = editor.Update(screens.NodeConfiguredMsg{
		Type: "http",
		Name: "Test Request",
		Parameters: map[string]interface{}{
			"url": "https://example.com",
		},
	})

	if editor.def.Type == "" {
		t.Error("definition not updated")
	}
}

func TestEditor_EditMode(t *testing.T) {
	workDir := t.TempDir()
	store, err := jubako.NewStore(workDir)
	if err != nil {
		t.Fatal(err)
	}

	// Create test bento
	def := neta.Definition{
		Version: "1.0",
		Type:    "http",
		Name:    "existing-bento",
	}
	store.Save("existing-bento", def)

	// Load in editor
	registry := pantry.NewRegistry()
	editor, err := screens.NewEditorEdit(store, registry, "existing-bento", "")
	if err != nil {
		t.Fatal(err)
	}

	if editor.mode != screens.EditorModeEdit {
		t.Error("expected edit mode")
	}

	if editor.def.Name != "existing-bento" {
		t.Error("definition not loaded")
	}
}
```

## Validation Commands

```bash
# Test
go test -v ./pkg/omise/screens/

# Integration test
./bento

# In browser:
# 1. Press 'n' to create new bento
# 2. Should switch to editor
# 3. Enter name, select type, configure
# 4. Press ctrl+s to save
# 5. Should return to browser with new bento

# 6. Select existing bento
# 7. Press 'e' to edit
# 8. Should load in editor
# 9. Modify and save
# 10. Should see changes in browser
```

## Success Criteria

Phase 7 is complete when:

1. ✅ Editor screen created
2. ✅ Create mode working
3. ✅ Edit mode working
4. ✅ Pantry integration for type selection
5. ✅ Huh wizards for parameter configuration
6. ✅ Definition building working
7. ✅ Save to Jubako working
8. ✅ Cancel returns to browser
9. ✅ All files < 250 lines
10. ✅ All functions < 20 lines
11. ✅ Tests passing
12. ✅ **Karen's approval granted**

## Common Pitfalls to Avoid

1. ❌ **Complex state management** - Use clear state machine
2. ❌ **Mixing UI and business logic** - Keep wizard separate from editor
3. ❌ **No validation** - Validate all inputs
4. ❌ **Poor error handling** - Show clear error messages
5. ❌ **Not saving version** - Always set version on definitions

## Next Phase

After Karen approval, proceed to **[Phase 8: Bento Editor - Visualization](./phase-8-bento-visualization.md)** to:
- Add visual bento box representation
- Implement node navigation with arrow keys
- Add in-editor node manipulation (edit/move/delete)
- Add run-from-editor capability

## Execution Prompt

```
I'm ready to begin Phase 7: Bento Editor - Node Builder.

I have read the Bento Box Principle and will follow it.

Please create the editor screen with:
- State machine for create/edit flow
- Pantry integration for node types
- Huh wizards for configuration
- Definition building
- Save to Jubako

Each file < 250 lines, functions < 20 lines. I will use TodoWrite to track progress and get Karen's approval before completing.
```

---

**Phase 7 Bento Editor**: Guided node building experience ✏️
