package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"

	"bento/pkg/omise/styles"
)

// FormSelect wraps a Huh select form
type FormSelect struct {
	form *huh.Form
}

// NewFormSelect creates a new form select component
func NewFormSelect(title, description string, options []SelectOption, value *string) FormSelect {
	huhOptions := make([]huh.Option[string], len(options))
	for i, opt := range options {
		huhOptions[i] = huh.NewOption(opt.Label, opt.Value)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(title).
				Description(description).
				Options(huhOptions...).
				Value(value),
		),
	).WithTheme(styles.FormTheme())

	return FormSelect{form: form}
}

// SelectOption represents a select option
type SelectOption struct {
	Label string
	Value string
}

// Init initializes the form
func (f FormSelect) Init() tea.Cmd {
	return f.form.Init()
}

// Update updates the form
func (f FormSelect) Update(msg tea.Msg) (FormSelect, tea.Cmd) {
	form, cmd := f.form.Update(msg)
	f.form = form.(*huh.Form)
	return f, cmd
}

// View renders the form
func (f FormSelect) View() string {
	return f.form.View()
}

// IsCompleted returns true if the form is completed
func (f FormSelect) IsCompleted() bool {
	return f.form.State == huh.StateCompleted
}

// GetForm returns the underlying form (for direct access if needed)
func (f FormSelect) GetForm() *huh.Form {
	return f.form
}
