// Package screens provides individual TUI screens.
package bentos

import (
	"bento/pkg/jubako"
	"bento/pkg/omise/components"
	"bento/pkg/omise/screens/guided_creation"
	"bento/pkg/omise/screens/shared"
	"bento/pkg/omise/styles"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Browser shows available bentos
type Browser struct {
	list          components.StyledList
	store         *jubako.Store
	discovery     *jubako.Discovery
	confirmDialog *shared.ConfirmDialog
	guidedModal   *guided_creation.GuidedModal // Guided creation/edit modal
	helpView      components.HelpView
	keys          components.BrowserKeyMap
	width         int // Current width
	height        int // Current height
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
	// PRIORITY 1: If guided modal is active, route ALL messages to it
	if b.guidedModal != nil {
		return b.updateGuidedModal(msg)
	}

	// PRIORITY 2: If confirmation dialog is active, handle it
	if b.confirmDialog != nil {
		return b.updateDialog(msg)
	}

	// PRIORITY 3: Normal browser updates
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
	case shared.BentoListRefreshMsg:
		newB, cmd := b.refreshList()
		return newB, cmd, true
	case shared.BentoOperationCompleteMsg:
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
	b.width = msg.Width - h
	b.height = msg.Height - v - 4
	b.list.SetSize(b.width, b.height)
	return b
}

// handleOperation handles operation completion
func (b Browser) handleOperation(msg shared.BentoOperationCompleteMsg) Browser {
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

// updateGuidedModal handles guided modal updates
func (b Browser) updateGuidedModal(msg tea.Msg) (Browser, tea.Cmd) {
	switch msg := msg.(type) {
	case guided_creation.GuidedCompleteMsg:
		// Modal finished - close it and refresh list
		b.guidedModal = nil
		if msg.Success {
			return b.refreshList()
		}
		return b, nil

	case tea.KeyMsg:
		// Allow Ctrl+C to quit even during guided flow
		if msg.String() == "ctrl+c" {
			return b, tea.Quit
		}
	}

	// All other messages go to the modal
	var cmd tea.Cmd
	newModal, cmd := b.guidedModal.Update(msg)
	if m, ok := newModal.(*guided_creation.GuidedModal); ok {
		b.guidedModal = m
	}
	return b, cmd
}

// updateDialog handles dialog updates
func (b Browser) updateDialog(msg tea.Msg) (Browser, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "y", "enter":
			path := b.confirmDialog.Context
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

// HasActiveModal returns true if a modal is currently active
func (b Browser) HasActiveModal() bool {
	return b.guidedModal != nil
}

// KeyBindings returns the contextual key bindings for the footer
func (b Browser) KeyBindings() []key.Binding {
	selected := b.getSelected()

	// If no item selected, show only new and search
	if selected == nil {
		return []key.Binding{
			b.keys.New,
			b.keys.Search,
		}
	}

	// For selected items, show action keys
	return []key.Binding{
		b.keys.Run,
		b.keys.Edit,
		b.keys.Copy,
		b.keys.Delete,
	}
}

