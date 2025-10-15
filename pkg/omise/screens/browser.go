// Package screens provides individual TUI screens.
package screens

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"bento/pkg/jubako"
	"bento/pkg/omise/components"
	"bento/pkg/omise/styles"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Browser shows available bentos
type Browser struct {
	list          components.StyledList
	store         *jubako.Store
	discovery     *jubako.Discovery
	confirmDialog *ConfirmDialog
	showingHelp   bool
}

// NewBrowser creates a browser screen with Jubako integration
func NewBrowser(workDir string) (Browser, error) {
	store, err := jubako.NewStore(workDir)
	if err != nil {
		return Browser{}, err
	}

	discovery := jubako.NewDiscovery(workDir)

	items, err := loadBentos(store)
	if err != nil {
		items = []list.Item{} // Empty list on error
	}

	l := components.NewStyledList(items, "Available Bentos")

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

	// Delegate to specific message handlers
	if newBrowser, cmd, handled := b.handleSpecialMsg(msg); handled {
		return newBrowser, cmd
	}

	// Default: update list
	var cmd tea.Cmd
	b.list.Model, cmd = b.list.Model.Update(msg)
	return b, cmd
}

// handleSpecialMsg handles special message types
func (b Browser) handleSpecialMsg(msg tea.Msg) (Browser, tea.Cmd, bool) {
	switch msg := msg.(type) {
	case styles.ThemeChangedMsg:
		b.list = b.list.RebuildStyles()
		return b, nil, true
	case tea.WindowSizeMsg:
		return b.handleResize(msg), nil, true
	case BentoListRefreshMsg:
		newB, cmd := b.refreshList()
		return newB, cmd, true
	case BentoOperationCompleteMsg:
		return b.handleOperation(msg), nil, true
	case tea.KeyMsg:
		newB, cmd := b.handleKey(msg)
		return newB, cmd, true
	}
	return b, nil, false
}

// handleResize updates browser dimensions
func (b Browser) handleResize(msg tea.WindowSizeMsg) Browser {
	h, v := lipgloss.NewStyle().Margin(2, 2).GetFrameSize()
	b.list.SetSize(msg.Width-h, msg.Height-v-4)
	return b
}

// handleOperation handles operation completion
func (b Browser) handleOperation(msg BentoOperationCompleteMsg) Browser {
	if msg.Success {
		b, _ = b.refreshList()
	}
	return b
}

// handleKey processes keyboard input
func (b Browser) handleKey(msg tea.KeyMsg) (Browser, tea.Cmd) {
	// Handle keys that don't require selection
	if newBrowser, cmd, handled := b.handleGlobalKey(msg); handled {
		return newBrowser, cmd
	}

	// Handle keys that require selection
	return b.handleItemKey(msg)
}

// handleGlobalKey handles keys that work without selection
func (b Browser) handleGlobalKey(msg tea.KeyMsg) (Browser, tea.Cmd, bool) {
	switch msg.String() {
	case "n":
		b, cmd := b.handleNew()
		return b, cmd, true
	case "?":
		b.showingHelp = !b.showingHelp
		return b, nil, true
	}
	return b, nil, false
}

// handleItemKey handles keys that require an item selected
func (b Browser) handleItemKey(msg tea.KeyMsg) (Browser, tea.Cmd) {
	selected := b.getSelected()
	if selected == nil {
		var cmd tea.Cmd
		b.list.Model, cmd = b.list.Model.Update(msg)
		return b, cmd
	}

	switch msg.String() {
	case "enter", " ", "r":
		return b.handleRun(selected)
	case "e":
		return b.handleEdit(selected)
	case "c":
		return b.handleCopy(selected)
	case "d":
		return b.handleDelete(selected)
	default:
		var cmd tea.Cmd
		b.list.Model, cmd = b.list.Model.Update(msg)
		return b, cmd
	}
}

// handleRun runs the selected bento or creates new if special item
func (b Browser) handleRun(item *bentoItem) (Browser, tea.Cmd) {
	if item.isNewItem {
		return b, func() tea.Msg {
			return CreateBentoMsg{}
		}
	}

	return b, func() tea.Msg {
		return WorkflowSelectedMsg{
			Name: item.name,
			Path: item.path,
		}
	}
}

// handleEdit edits the selected bento
func (b Browser) handleEdit(item *bentoItem) (Browser, tea.Cmd) {
	if item.isNewItem {
		// Create new instead of edit for the special item
		return b, func() tea.Msg {
			return CreateBentoMsg{}
		}
	}

	return b, func() tea.Msg {
		return EditBentoMsg{
			Name: item.name,
			Path: item.path,
		}
	}
}

// handleCopy initiates bento copy
func (b Browser) handleCopy(item *bentoItem) (Browser, tea.Cmd) {
	if item.isNewItem {
		return b, nil // Can't copy the "Create New" item
	}
	return b, b.copyBento(item)
}

// handleDelete shows delete confirmation
func (b Browser) handleDelete(item *bentoItem) (Browser, tea.Cmd) {
	if item.isNewItem {
		return b, nil // Can't delete the "Create New" item
	}

	b.confirmDialog = NewConfirmDialog(
		"Delete Bento",
		fmt.Sprintf("Are you sure you want to delete '%s'?", item.name),
		item.path,
	)
	return b, nil
}

// handleNew creates a new bento
func (b Browser) handleNew() (Browser, tea.Cmd) {
	return b, func() tea.Msg {
		return CreateBentoMsg{}
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

// getSelected returns the selected bento item
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
		newName := generateCopyName(item.name)
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
		name := extractBentoName(path)
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
	items, err := loadBentos(b.store)
	if err != nil {
		items = []list.Item{}
	}

	b.list = components.NewStyledList(items, "Available Bentos")
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

	return lipgloss.JoinVertical(
		lipgloss.Left,
		b.list.View(),
		"",
		b.renderFooter(),
	)
}

// renderFooter shows keyboard shortcuts
func (b Browser) renderFooter() string {
	shortcuts := "n: New • enter: Run • e: Edit • c: Copy • d: Delete • ?: Help"
	return styles.Subtle.Render(shortcuts)
}

// helpView renders keyboard shortcuts
func (b Browser) helpView() string {
	help := `
Keyboard Shortcuts:

  enter/space/r  Run bento
  e              Edit bento (Phase 7)
  c              Copy bento
  d              Delete bento
  n              Create new bento (Phase 7)
  ?              Toggle this help
  tab            Next screen
  q              Quit

Press ? again to return to list.
`
	return styles.Subtle.Render(help)
}

// bentoItem represents a bento in the list
type bentoItem struct {
	name      string
	path      string
	version   string
	nodeType  string
	modified  time.Time
	isNewItem bool // Special item for creating new bentos
}

// Title returns the item title
func (i bentoItem) Title() string {
	if i.isNewItem {
		return "+ Create New Bento"
	}
	return fmt.Sprintf("%s (v%s)", i.name, i.version)
}

// Description returns the item description
func (i bentoItem) Description() string {
	if i.isNewItem {
		return "Start building a new bento from scratch"
	}
	return fmt.Sprintf("%s • Modified: %s", i.nodeType, i.modified.Format("2006-01-02 15:04"))
}

// FilterValue returns the value to filter by
func (i bentoItem) FilterValue() string {
	if i.isNewItem {
		return "new create"
	}
	return i.name
}

// loadBentos loads bentos from store
func loadBentos(store *jubako.Store) ([]list.Item, error) {
	infos, err := store.List()
	if err != nil {
		return nil, err
	}

	// Start with "Create New" item
	items := make([]list.Item, 0, len(infos)+1)
	items = append(items, bentoItem{
		isNewItem: true,
	})

	// Add existing bentos
	for _, info := range infos {
		def, err := store.Load(extractBentoName(info.Name))
		if err != nil {
			continue // Skip invalid files
		}

		items = append(items, bentoItem{
			name:     extractBentoName(info.Name),
			path:     info.Path,
			version:  def.Version,
			nodeType: def.Type,
			modified: info.Modified,
		})
	}

	return items, nil
}

// generateCopyName creates a unique name for a copied bento
func generateCopyName(name string) string {
	base := strings.TrimSuffix(name, ".bento.yaml")
	return fmt.Sprintf("%s-copy", base)
}

// extractBentoName extracts the bento name from a path or filename
func extractBentoName(pathOrName string) string {
	base := filepath.Base(pathOrName)
	return strings.TrimSuffix(base, ".bento.yaml")
}
