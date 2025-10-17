// Package screens provides individual TUI screens.
package screens

import (
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
	helpView      components.HelpView
	keys          components.BrowserKeyMap
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
		items = []list.Item{}
	}

	b := Browser{
		list:      components.NewStyledList(items, ""),
		store:     store,
		discovery: discovery,
		helpView:  components.NewHelpView(),
		keys:      components.NewBrowserKeyMap(),
	}
	b.list.SetSize(80, 20) // Set default size, will be updated on window resize
	return b, nil
}

// Init initializes the browser
func (b Browser) Init() tea.Cmd {
	return nil
}

// Update handles browser messages
func (b Browser) Update(msg tea.Msg) (Browser, tea.Cmd) {
	if b.confirmDialog != nil {
		return b.updateDialog(msg)
	}

	if newBrowser, cmd, handled := b.handleSpecialMsg(msg); handled {
		return newBrowser, cmd
	}

	// Update the list
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
	if newBrowser, cmd, handled := b.handleGlobalKey(msg); handled {
		return newBrowser, cmd
	}

	return b.handleItemKey(msg)
}

// handleGlobalKey handles keys that work without selection
func (b Browser) handleGlobalKey(msg tea.KeyMsg) (Browser, tea.Cmd, bool) {
	switch msg.String() {
	case "n":
		b, cmd := b.handleNew()
		return b, cmd, true
	case "?":
		b.helpView = b.helpView.Toggle()
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
	case "r":
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

// updateDialog handles dialog updates
func (b Browser) updateDialog(msg tea.Msg) (Browser, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "y", "enter":
			path := b.confirmDialog.context
			b.confirmDialog = nil
			return b, b.deleteBento(path)
		case "n", "esc":
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

// refreshList reloads bentos from disk
func (b Browser) refreshList() (Browser, tea.Cmd) {
	items, err := loadBentos(b.store)
	if err != nil {
		items = []list.Item{}
	}

	// Preserve current list size when refreshing
	width, height := b.list.Width(), b.list.Height()
	b.list = components.NewStyledList(items, "")
	b.list.SetSize(width, height)
	return b, nil
}

// ContextualKeys returns the most important contextual keys for the footer
func (b Browser) ContextualKeys() []components.KeyHelp {
	// Check if we have an item selected
	selected := b.getSelected()

	// If no item selected, show only new
	if selected == nil {
		return []components.KeyHelp{
			{Key: "n", Desc: "new bento"},
		}
	}

	// For selected items, show all action keys
	return []components.KeyHelp{
		{Key: "n", Desc: "new"},
		{Key: "r", Desc: "run"},
		{Key: "e", Desc: "edit"},
		{Key: "c", Desc: "copy"},
		{Key: "d", Desc: "delete"},
	}
}
