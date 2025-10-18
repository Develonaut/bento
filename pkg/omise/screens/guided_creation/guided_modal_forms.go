package guided_creation

import (
	"fmt"

	"github.com/charmbracelet/huh"

	"bento/pkg/omise/emoji"
)

func (m *GuidedModal) createMetadataForm() *huh.Form {
	var name, description string

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("Name").
				Description("A short, descriptive name for this workflow").
				Value(&name).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("name is required")
					}
					return nil
				}),

			huh.NewInput().
				Key("description").
				Title("Description").
				Description("What does this bento do?").
				Value(&description).
				CharLimit(200),
		).Title("Bento Metadata:"),
	).
		WithWidth(m.width - 48).
		WithShowHelp(false).
		WithShowErrors(false)
}

func (m *GuidedModal) updateDefinitionFromForm() {
	// Update definition with current form values
	if name := m.form.GetString("name"); name != "" {
		m.definition.Name = name
	}
	if desc := m.form.GetString("description"); desc != "" {
		m.definition.Description = desc
	}
	// Assign random icon if not already set
	if m.definition.Icon == "" {
		m.definition.Icon = emoji.RandomSushi()
	}
}

func (m *GuidedModal) createNodeTypeSelectForm() *huh.Form {
	var nodeType string
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("node_type").
				Title("Node Type").
				Description("Choose the type of node to add").
				Options(m.nodeTypeOptions()...).
				Value(&nodeType),
		).Title("Add Node:"),
	).WithWidth(m.width - 48).WithShowHelp(false).WithShowErrors(false)
}

func (m *GuidedModal) nodeTypeOptions() []huh.Option[string] {
	return []huh.Option[string]{
		huh.NewOption("HTTP Request", "http"),
		huh.NewOption("Transform (jq)", "transform.jq"),
		huh.NewOption("File Write", "file.write"),
		huh.NewOption("Sequence Group", "group.sequence"),
		huh.NewOption("Parallel Group", "group.parallel"),
	}
}

func (m *GuidedModal) createNodeFormForType(nodeType string) *huh.Form {
	switch nodeType {
	case "http":
		return m.createHTTPNodeForm()
	case "transform.jq", "jq":
		return m.createJQNodeForm()
	case "file.write":
		return m.createFileWriteNodeForm()
	case "group.sequence", "sequence":
		return m.createSequenceNodeForm()
	case "group.parallel", "parallel":
		return m.createParallelNodeForm()
	default:
		// Fallback to generic form
		return m.createNodeTypeSelectForm()
	}
}

func (m *GuidedModal) createContinueForm() *huh.Form {
	var choice string
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("continue").
				Title("What would you like to do?").
				Options(
					huh.NewOption("Add another node", "add"),
					huh.NewOption("Done - Save bento", "done"),
				).
				Value(&choice),
		).Title("Node Added Successfully!"),
	).WithWidth(m.width - 48).WithShowHelp(false).WithShowErrors(false)
}
