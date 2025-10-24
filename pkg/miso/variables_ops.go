package miso

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// loadVariablesView loads variables and switches to variables view
func (m Model) loadVariablesView() (tea.Model, tea.Cmd) {
	mgr, err := NewVariablesManager()
	if err != nil {
		m.logs = fmt.Sprintf("Failed to load variables: %v", err)
		m.currentView = executionView
		return m, nil
	}

	vars := mgr.GetAll()

	// Build variables list
	items := make([]list.Item, 0, len(vars))
	for key, value := range vars {
		items = append(items, VariableItem{
			Key:   key,
			Value: value,
		})
	}

	m.variablesList.SetItems(items)
	m.currentView = variablesView
	return m, nil
}

// addVariable prompts for a new variable
func (m Model) addVariable() (tea.Model, tea.Cmd) {
	// Use Huh form to get variable key and value
	var key, value string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Variable Key").
				Description("Uppercase letters, numbers, and underscores").
				Placeholder("PRODUCTS_URL").
				Value(&key),
			huh.NewInput().
				Title("Variable Value").
				Description("The value to store").
				Placeholder("/Users/you/Products").
				Value(&value),
		),
	).WithTheme(huh.ThemeCharm())

	if err := form.Run(); err != nil {
		// User cancelled
		return m, nil
	}

	// Store variable
	mgr, err := NewVariablesManager()
	if err != nil {
		m.logs = fmt.Sprintf("Failed to initialize variables: %v", err)
		return m, nil
	}

	if err := mgr.Set(key, value); err != nil {
		m.logs = fmt.Sprintf("Failed to store variable: %v", err)
		return m, nil
	}

	// Reload variables view
	return m.loadVariablesView()
}

// deleteVariable removes a variable
func (m Model) deleteVariable(key string) (tea.Model, tea.Cmd) {
	mgr, err := NewVariablesManager()
	if err != nil {
		m.logs = fmt.Sprintf("Failed to initialize variables: %v", err)
		return m, nil
	}

	if err := mgr.Delete(key); err != nil {
		m.logs = fmt.Sprintf("Failed to delete variable: %v", err)
		return m, nil
	}

	// Reload variables view
	return m.loadVariablesView()
}
