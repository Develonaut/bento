# Phase 6: Enhanced Browser & CRUD Operations

**Status**: Pending
**Duration**: 3-4 hours
**Prerequisites**: Phase 5.5 complete, Karen approved

## Overview

Transform the browser from a simple list into a full-featured bento management interface with keyboard shortcuts for run, edit, copy, and delete operations. Integrate with Jubako for real-time discovery and CRUD operations.

## Pre-Work Checklist

Before starting, you MUST:

1. ✅ Read [BENTO_BOX_PRINCIPLE.md](../BENTO_BOX_PRINCIPLE.md)
2. ✅ Read [CHARM_STACK_GUIDE.md](../CHARM_STACK_GUIDE.md)
3. ✅ Confirm: "I understand the Bento Box Principle and will follow it"
4. ✅ Use TodoWrite to track all tasks
5. ✅ Phase 5.5 approved by Karen

## Goals

1. Add keyboard shortcuts: `r` (run), `e` (edit), `c` (copy), `d` (delete)
2. Integrate Jubako for dynamic bento discovery
3. Implement confirmation dialog for delete
4. Add "Create New Bento" action
5. Display bento metadata (version, type, last modified)
6. Real-time bento list updates
7. Validate Bento Box compliance

## Current vs Enhanced Flow

### Current Flow
```
Browser Screen:
- Show hardcoded list
- Press Enter/Space → Run bento
- Tab → Next screen
```

### Enhanced Flow
```
Browser Screen:
- Show dynamic list from Jubako
- Press Enter/Space/r → Run bento
- Press e → Edit bento (→ Editor screen)
- Press c → Copy bento (duplicate file)
- Press d → Delete bento (with confirmation)
- Press n → Create new bento (→ Editor screen)
- Tab → Next screen

Bentos show: Name, Version, Type, Last Modified
```

## Screen Structure

```
pkg/omise/screens/
├── browser.go           # Enhanced browser (modify)
├── browser_test.go      # Browser tests (modify)
├── confirm.go           # Confirmation dialog (NEW)
└── confirm_test.go      # Dialog tests (NEW)

pkg/omise/components/
└── list.go              # Already exists, may need updates
```

## Deliverables

### 1. Message Types

**File**: `pkg/omise/screens/messages.go` (NEW)
**Target Size**: < 100 lines

```go
package screens

// WorkflowSelectedMsg signals workflow selected for execution
type WorkflowSelectedMsg struct {
	Name string
	Path string
}

// EditBentoMsg signals user wants to edit a bento
type EditBentoMsg struct {
	Name string
	Path string
}

// CreateBentoMsg signals user wants to create new bento
type CreateBentoMsg struct{}

// CopyBentoMsg signals user wants to copy a bento
type CopyBentoMsg struct {
	Name string
	Path string
}

// DeleteBentoMsg signals user confirmed deletion
type DeleteBentoMsg struct {
	Name string
	Path string
}

// BentoListRefreshMsg signals bento list should reload
type BentoListRefreshMsg struct{}

// BentoOperationCompleteMsg signals operation completed
type BentoOperationCompleteMsg struct {
	Operation string // "copy", "delete", "create"
	Success   bool
	Error     error
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: Message definitions
- ✅ Clear event-driven architecture
- ✅ File < 100 lines

### 2. Enhanced Browser

**File**: `pkg/omise/screens/browser.go` (modify existing)
**Target Size**: < 250 lines

```go
package screens

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"bento/pkg/jubako"
	"bento/pkg/omise/components"
	"bento/pkg/omise/styles"
)

// Browser shows available workflows
type Browser struct {
	list         components.StyledList
	store        *jubako.Store
	discovery    *jubako.Discovery
	confirmDialog *ConfirmDialog
	showingHelp   bool
}

// NewBrowser creates a browser screen
func NewBrowser(workDir string) (Browser, error) {
	store, err := jubako.NewStore(workDir)
	if err != nil {
		return Browser{}, err
	}

	discovery := jubako.NewDiscovery(workDir)

	items, err := loadBentos(store, discovery)
	if err != nil {
		items = []list.Item{} // Empty list on error
	}

	l := components.NewStyledList(items, "🍱 Available Bentos")

	return Browser{
		list:      l,
		store:     store,
		discovery: discovery,
	}, nil
}

// Init initializes the browser
func (b Browser) Init() tea.Cmd {
	return nil
}

// Update handles browser messages
func (b Browser) Update(msg tea.Msg) (Browser, tea.Cmd) {
	// Handle confirmation dialog if active
	if b.confirmDialog != nil {
		return b.updateDialog(msg)
	}

	// Handle theme changes
	if _, ok := msg.(styles.ThemeChangedMsg); ok {
		b.list = b.list.RebuildStyles()
	}

	// Handle window resize
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		h, v := lipgloss.NewStyle().Margin(2, 2).GetFrameSize()
		b.list.SetSize(msg.Width-h, msg.Height-v-4)
	}

	// Handle bento list refresh
	if _, ok := msg.(BentoListRefreshMsg); ok {
		return b.refreshList()
	}

	// Handle keyboard shortcuts
	if msg, ok := msg.(tea.KeyMsg); ok {
		return b.handleKey(msg)
	}

	var cmd tea.Cmd
	b.list.Model, cmd = b.list.Model.Update(msg)
	return b, cmd
}

// handleKey processes keyboard input
func (b Browser) handleKey(msg tea.KeyMsg) (Browser, tea.Cmd) {
	selected := b.getSelected()
	if selected == nil {
		// No item selected, only handle list navigation
		var cmd tea.Cmd
		b.list.Model, cmd = b.list.Model.Update(msg)
		return b, cmd
	}

	switch msg.String() {
	case "enter", " ", "r":
		// Run bento
		return b, func() tea.Msg {
			return WorkflowSelectedMsg{
				Name: selected.name,
				Path: selected.path,
			}
		}

	case "e":
		// Edit bento
		return b, func() tea.Msg {
			return EditBentoMsg{
				Name: selected.name,
				Path: selected.path,
			}
		}

	case "c":
		// Copy bento
		return b, b.copyBento(selected)

	case "d":
		// Show delete confirmation
		b.confirmDialog = NewConfirmDialog(
			"Delete Bento",
			fmt.Sprintf("Are you sure you want to delete '%s'?", selected.name),
			selected.path,
		)
		return b, nil

	case "n":
		// Create new bento
		return b, func() tea.Msg {
			return CreateBentoMsg{}
		}

	case "?":
		// Toggle help
		b.showingHelp = !b.showingHelp
		return b, nil

	default:
		var cmd tea.Cmd
		b.list.Model, cmd = b.list.Model.Update(msg)
		return b, cmd
	}
}

// updateDialog handles dialog updates
func (b Browser) updateDialog(msg tea.Msg) (Browser, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "y", "enter":
			// Confirmed deletion
			path := b.confirmDialog.context
			b.confirmDialog = nil
			return b, b.deleteBento(path)

		case "n", "esc":
			// Cancelled
			b.confirmDialog = nil
			return b, nil
		}
	}

	return b, nil
}

// getSelected returns the selected workflow item
func (b Browser) getSelected() *bentoItem {
	if item, ok := b.list.SelectedItem().(bentoItem); ok {
		return &item
	}
	return nil
}

// copyBento duplicates a bento file
func (b Browser) copyBento(item *bentoItem) tea.Cmd {
	return func() tea.Msg {
		def, err := b.store.Load(item.name)
		if err != nil {
			return BentoOperationCompleteMsg{
				Operation: "copy",
				Success:   false,
				Error:     err,
			}
		}

		// Create new name
		newName := fmt.Sprintf("%s-copy", item.name)
		def.Name = newName

		if err := b.store.Save(newName, def); err != nil {
			return BentoOperationCompleteMsg{
				Operation: "copy",
				Success:   false,
				Error:     err,
			}
		}

		return BentoOperationCompleteMsg{
			Operation: "copy",
			Success:   true,
		}
	}
}

// deleteBento removes a bento file
func (b Browser) deleteBento(path string) tea.Cmd {
	return func() tea.Msg {
		name := filepath.Base(path)
		if err := b.store.Delete(name); err != nil {
			return BentoOperationCompleteMsg{
				Operation: "delete",
				Success:   false,
				Error:     err,
			}
		}

		return BentoOperationCompleteMsg{
			Operation: "delete",
			Success:   true,
		}
	}
}

// refreshList reloads bentos from disk
func (b Browser) refreshList() (Browser, tea.Cmd) {
	items, err := loadBentos(b.store, b.discovery)
	if err != nil {
		items = []list.Item{}
	}

	b.list = components.NewStyledList(items, "🍱 Available Bentos")
	return b, nil
}

// View renders the browser
func (b Browser) View() string {
	if b.confirmDialog != nil {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			b.list.View(),
			"",
			b.confirmDialog.View(),
		)
	}

	if b.showingHelp {
		return b.helpView()
	}

	return b.list.View()
}

// helpView renders keyboard shortcuts
func (b Browser) helpView() string {
	help := `
Keyboard Shortcuts:

  enter/space/r  Run bento
  e              Edit bento
  c              Copy bento
  d              Delete bento
  n              Create new bento
  ?              Toggle this help
  tab            Next screen
  q              Quit

Press ? again to return to list.
`
	return styles.Subtle.Render(help)
}

// bentoItem represents a bento in the list
type bentoItem struct {
	name     string
	path     string
	version  string
	nodeType string
	modified time.Time
}

// Title returns the item title
func (i bentoItem) Title() string {
	return fmt.Sprintf("%s (v%s)", i.name, i.version)
}

// Description returns the item description
func (i bentoItem) Description() string {
	return fmt.Sprintf("%s • Modified: %s", i.nodeType, i.modified.Format("2006-01-02 15:04"))
}

// FilterValue returns the value to filter by
func (i bentoItem) FilterValue() string {
	return i.name
}

// loadBentos loads bentos from store
func loadBentos(store *jubako.Store, discovery *jubako.Discovery) ([]list.Item, error) {
	infos, err := store.List()
	if err != nil {
		return nil, err
	}

	items := make([]list.Item, len(infos))
	for i, info := range infos {
		def, err := store.Load(info.Name)
		if err != nil {
			continue
		}

		items[i] = bentoItem{
			name:     info.Name,
			path:     info.Path,
			version:  def.Version,
			nodeType: def.Type,
			modified: info.Modified,
		}
	}

	return items, nil
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: Bento browsing and CRUD
- ✅ Functions < 20 lines
- ✅ Clear message-driven architecture
- ✅ File < 250 lines

### 3. Confirmation Dialog

**File**: `pkg/omise/screens/confirm.go` (NEW)
**Target Size**: < 100 lines

```go
package screens

import (
	"github.com/charmbracelet/lipgloss"
	"bento/pkg/omise/styles"
)

// ConfirmDialog is a simple yes/no confirmation dialog
type ConfirmDialog struct {
	title   string
	message string
	context string // Context data (e.g., path to delete)
}

// NewConfirmDialog creates a confirmation dialog
func NewConfirmDialog(title, message, context string) *ConfirmDialog {
	return &ConfirmDialog{
		title:   title,
		message: message,
		context: context,
	}
}

// View renders the confirmation dialog
func (c *ConfirmDialog) View() string {
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Warning).
		Padding(1, 2).
		Width(50)

	title := styles.Bold.Foreground(styles.Warning).Render(c.title)
	message := styles.Subtle.Render(c.message)
	prompt := styles.Subtle.Render("\nPress Y to confirm, N to cancel")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		message,
		prompt,
	)

	return box.Render(content)
}
```

**Bento Box Compliance**:
- ✅ Single responsibility: Confirmation dialog
- ✅ Simple, reusable component
- ✅ File < 100 lines

### 4. Update Root Model

**File**: `pkg/omise/update.go` (modify handleWorkflowSelected, add handlers)

```go
// Add new message handlers

// handleEditBento switches to editor for existing bento
func (m Model) handleEditBento(msg screens.EditBentoMsg) (tea.Model, tea.Cmd) {
	// Phase 7: Will switch to editor screen
	// For now, just acknowledge
	return m, nil
}

// handleCreateBento switches to editor for new bento
func (m Model) handleCreateBento(msg screens.CreateBentoMsg) (tea.Model, tea.Cmd) {
	// Phase 7: Will switch to editor screen
	// For now, just acknowledge
	return m, nil
}

// handleBentoOperation handles completion of copy/delete operations
func (m Model) handleBentoOperation(msg screens.BentoOperationCompleteMsg) (tea.Model, tea.Cmd) {
	if msg.Success {
		// Refresh browser list
		return m.updateScreen(screens.BentoListRefreshMsg{})
	}
	// Show error (could add status bar in future)
	return m, nil
}

// Update Update function to handle new messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleResize(msg)
	case tea.KeyMsg:
		return m.handleKey(msg)
	case screens.WorkflowSelectedMsg:
		return m.handleWorkflowSelected(msg)
	case screens.EditBentoMsg:
		return m.handleEditBento(msg)
	case screens.CreateBentoMsg:
		return m.handleCreateBento(msg)
	case screens.BentoOperationCompleteMsg:
		return m.handleBentoOperation(msg)
	case styles.ThemeChangedMsg:
		return m.handleThemeChanged(msg)
	default:
		return m.updateScreen(msg)
	}
}
```

## Integration with Jubako

### Store Setup

The browser needs a configured Jubako store:

**File**: `pkg/omise/app.go` (modify)

```go
import (
	"os"
	"path/filepath"
)

// Launch starts the TUI application
func Launch() error {
	// Get or create work directory
	workDir, err := getWorkDir()
	if err != nil {
		return err
	}

	// Initialize model with work directory
	m, err := NewModelWithWorkDir(workDir)
	if err != nil {
		return err
	}

	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	_, err = p.Run()
	return err
}

// getWorkDir returns the bento work directory
func getWorkDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	workDir := filepath.Join(home, ".bento", "workflows")
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return "", err
	}

	return workDir, nil
}
```

**File**: `pkg/omise/model.go` (modify)

```go
// NewModelWithWorkDir creates model with configured work directory
func NewModelWithWorkDir(workDir string) (Model, error) {
	browser, err := screens.NewBrowser(workDir)
	if err != nil {
		return Model{}, err
	}

	return Model{
		screen:   ScreenBrowser,
		browser:  browser,
		executor: screens.NewExecutor(),
		pantry:   screens.NewPantry(),
		settings: screens.NewSettings(),
		help:     screens.NewHelp(),
	}, nil
}
```

## Testing Strategy

### Browser Tests

**File**: `pkg/omise/screens/browser_test.go`

```go
func TestBrowser_KeyboardShortcuts(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		wantMsg  interface{}
	}{
		{
			name:    "r runs bento",
			key:     "r",
			wantMsg: WorkflowSelectedMsg{},
		},
		{
			name:    "e edits bento",
			key:     "e",
			wantMsg: EditBentoMsg{},
		},
		{
			name:    "c copies bento",
			key:     "c",
			wantMsg: BentoOperationCompleteMsg{},
		},
		{
			name:    "d shows confirmation",
			key:     "d",
			wantMsg: nil, // Shows dialog
		},
		{
			name:    "n creates bento",
			key:     "n",
			wantMsg: CreateBentoMsg{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test keyboard shortcuts
		})
	}
}

func TestBrowser_LoadBentos(t *testing.T) {
	// Test loading bentos from Jubako
	workDir := t.TempDir()
	store, err := jubako.NewStore(workDir)
	if err != nil {
		t.Fatal(err)
	}

	// Create test bento
	def := neta.Definition{
		Version: "1.0",
		Type:    "http",
		Name:    "test",
	}

	if err := store.Save("test", def); err != nil {
		t.Fatal(err)
	}

	// Load bentos
	browser, err := screens.NewBrowser(workDir)
	if err != nil {
		t.Fatal(err)
	}

	// Verify bento appears in list
	// ...
}
```

## Validation Commands

```bash
# Format
go fmt ./pkg/omise/...

# Lint
golangci-lint run ./pkg/omise/...

# Test
go test -v ./pkg/omise/screens/

# Integration test
# 1. Create test bentos
mkdir -p ~/.bento/workflows

cat > ~/.bento/workflows/test1.bento.yaml <<EOF
version: "1.0"
type: http
name: Test 1
parameters:
  url: https://httpbin.org/get
EOF

cat > ~/.bento/workflows/test2.bento.yaml <<EOF
version: "1.0"
type: http
name: Test 2
parameters:
  url: https://httpbin.org/post
EOF

# 2. Launch TUI
./bento

# 3. Test keyboard shortcuts:
#    - Navigate list with arrows
#    - Press 'r' to run
#    - Press 'c' to copy
#    - Press 'd' to delete (confirm with 'y')
#    - Press 'n' to create (Phase 7)
#    - Press '?' for help

# File size check
find pkg/omise/screens -name "*.go" -exec wc -l {} + | sort -rn
```

## Success Criteria

Phase 6 is complete when:

1. ✅ Browser shows bentos from Jubako
2. ✅ Keyboard shortcuts work (r, e, c, d, n)
3. ✅ Delete confirmation dialog working
4. ✅ Copy creates duplicate file
5. ✅ Bento metadata displayed (version, type, modified)
6. ✅ Help screen shows shortcuts
7. ✅ All files < 250 lines
8. ✅ All functions < 20 lines
9. ✅ Tests passing
10. ✅ **Karen's approval granted**

## Common Pitfalls to Avoid

1. ❌ **No error handling** - File operations can fail
2. ❌ **Not refreshing list** - Must reload after copy/delete
3. ❌ **Hardcoded paths** - Use configurable work directory
4. ❌ **No confirmation** - Always confirm destructive operations
5. ❌ **Poor UX** - Show clear feedback for operations

## Next Phase

After Karen approval, proceed to **[Phase 7: Bento Editor - Node Builder](./phase-7-bento-editor-builder.md)** to:
- Create editor screen
- Add pantry-based node selection
- Implement Huh configuration wizards
- Build Definition structures

## Execution Prompt

```
I'm ready to begin Phase 6: Enhanced Browser & CRUD Operations.

I have read the Bento Box Principle and will follow it.

Please enhance the browser with:
- Keyboard shortcuts (r, e, c, d, n)
- Jubako integration for dynamic lists
- Confirmation dialogs
- Copy and delete operations
- Metadata display

Each file < 250 lines, functions < 20 lines. I will use TodoWrite to track progress and get Karen's approval before completing.
```

---

**Phase 6 Enhanced Browser**: Full-featured bento management 🗂️
