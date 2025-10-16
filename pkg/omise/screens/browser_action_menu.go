package screens

import (
	"bento/pkg/omise/styles"

	"github.com/charmbracelet/huh"
)

// BentoAction represents an action that can be performed on a bento
type BentoAction string

const (
	ActionRun    BentoAction = "run"
	ActionEdit   BentoAction = "edit"
	ActionCopy   BentoAction = "copy"
	ActionDelete BentoAction = "delete"
	ActionCancel BentoAction = "cancel"
)

// BentoActionMenu shows available actions for a selected bento
type BentoActionMenu struct {
	form     *huh.Form
	item     *bentoItem
	selected string
}

// NewBentoActionMenu creates an action menu for a bento item
func NewBentoActionMenu(item *bentoItem) *BentoActionMenu {
	menu := &BentoActionMenu{
		item: item,
	}

	menu.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What would you like to do with '"+item.name+"'?").
				Options(
					huh.NewOption("Run Bento", string(ActionRun)),
					huh.NewOption("Edit Bento", string(ActionEdit)),
					huh.NewOption("Copy Bento", string(ActionCopy)),
					huh.NewOption("Delete Bento", string(ActionDelete)),
					huh.NewOption("Cancel", string(ActionCancel)),
				).
				Value(&menu.selected),
		),
	).WithTheme(styles.FormTheme())

	return menu
}

// GetSelectedAction returns the selected action after form completion
func (m *BentoActionMenu) GetSelectedAction() BentoAction {
	return BentoAction(m.selected)
}
