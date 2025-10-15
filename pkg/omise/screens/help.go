package screens

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"bento/pkg/omise/styles"
)

// Help shows keyboard shortcuts and usage information
type Help struct{}

// NewHelp creates a help screen
func NewHelp() Help {
	return Help{}
}

// Init initializes the help screen
func (h Help) Init() tea.Cmd {
	return nil
}

// Update handles help messages
func (h Help) Update(msg tea.Msg) (Help, tea.Cmd) {
	// Handle theme changes (styles are global, no rebuild needed for Help)
	if _, ok := msg.(styles.ThemeChangedMsg); ok {
		return h, nil
	}
	return h, nil
}

// View renders the help
func (h Help) View() string {
	title := styles.Title.Render("Help - Keyboard Shortcuts")
	content := h.renderSections()
	about := h.renderAbout()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		content,
		about,
	)
}

// renderSections renders all help sections
func (h Help) renderSections() string {
	sections := helpSections()
	var content string
	for _, section := range sections {
		content += h.renderSection(section)
	}
	return content
}

// renderSection renders a single help section
func (h Help) renderSection(section helpSection) string {
	var content string
	content += "\n" + styles.Selected.Render(section.title) + "\n"
	for _, item := range section.items {
		key := styles.HelpKey.Render(item[0])
		desc := styles.HelpDesc.Render(item[1])
		content += "  " + key + "  " + desc + "\n"
	}
	return content
}

// renderAbout renders the about section
func (h Help) renderAbout() string {
	return "\n" + styles.Subtle.Render(
		"Bento - Organized bento orchestration\n"+
			"Version 0.1.0 (Phase 4)\n"+
			"Omise (お店) - The shop where bentos are served",
	)
}

type helpSection struct {
	title string
	items [][]string
}

// helpSections returns the help section data
func helpSections() []helpSection {
	return []helpSection{
		{
			title: "Navigation",
			items: [][]string{
				{"tab", "Next screen"},
				{"shift+tab", "Previous screen"},
				{"↑/k", "Move up"},
				{"↓/j", "Move down"},
				{"?/h", "Show help"},
			},
		},
		{
			title: "Browser Screen",
			items: [][]string{
				{"enter/space", "Execute selected bento"},
				{"/", "Search bentos"},
				{"esc", "Clear search"},
			},
		},
		{
			title: "Pantry Screen",
			items: [][]string{
				{"↑/↓", "Navigate neta types"},
				{"enter/space", "View details (coming soon)"},
			},
		},
		{
			title: "General",
			items: [][]string{
				{"q", "Quit application"},
				{"ctrl+c", "Force quit"},
			},
		},
	}
}
