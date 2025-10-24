package miso

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// showForm creates and displays the variable form
func (m Model) showForm() (tea.Model, tea.Cmd) {
	// Sort variables to show path variables first
	sortedVars := sortVariablesByPriority(m.bentoVars)

	// Create value holders for each variable
	valueHolders := make(map[string]*string)
	for _, v := range sortedVars {
		holder := v.DefaultValue
		valueHolders[v.Name] = &holder
	}

	// Build form fields - use full terminal height for sizing
	fields := make([]huh.Field, 0, len(sortedVars))
	for _, v := range sortedVars {
		fields = append(fields, buildFieldWithHeight(v, valueHolders[v.Name], m.height))
	}

	// Create form with dimensions
	m.form = huh.NewForm(
		huh.NewGroup(fields...),
	).WithTheme(huh.ThemeCharm()).
		WithWidth(m.width).
		WithHeight(m.height)

	// Store value holders so we can extract values after form completion
	m.varHolders = valueHolders

	m.currentView = formView
	return m, m.form.Init()
}

// sortVariablesByPriority sorts variables to show path variables first
func sortVariablesByPriority(vars []Variable) []Variable {
	// Create a copy to avoid modifying original
	sorted := make([]Variable, len(vars))
	copy(sorted, vars)

	// Sort with custom comparator:
	// 1. Path variables (containing PATH, DIR, FOLDER) first
	// 2. Within each group, alphabetically
	var pathVars, otherVars []Variable

	for _, v := range sorted {
		upperName := strings.ToUpper(v.Name)
		if strings.Contains(upperName, "PATH") ||
			strings.Contains(upperName, "DIR") ||
			strings.Contains(upperName, "DIRECTORY") ||
			strings.Contains(upperName, "FOLDER") {
			pathVars = append(pathVars, v)
		} else {
			otherVars = append(otherVars, v)
		}
	}

	// Combine: path vars first, then others
	result := make([]Variable, 0, len(sorted))
	result = append(result, pathVars...)
	result = append(result, otherVars...)

	return result
}

// startExecution runs the bento with collected variables
func (m Model) startExecution() (tea.Model, tea.Cmd) {
	// Clear logs from previous execution
	m.logs = fmt.Sprintf("🍱 Executing: %s\n\n", filepath.Base(m.selectedBento))

	// Set environment variables from form value holders
	if m.varHolders != nil {
		for name, valuePtr := range m.varHolders {
			if valuePtr != nil && *valuePtr != "" {
				os.Setenv(name, *valuePtr)
				// Log the variable being set for debugging
				m.logs += fmt.Sprintf("Set %s = %s\n", name, *valuePtr)
			}
		}
		m.logs += "\n"
	}

	m.currentView = executionView
	m.executing = true

	// Initialize viewport for log display
	// Account for: title (3 lines) + help (1 line) + border (2 lines) + border padding (2 lines)
	viewportHeight := m.height - 8
	if viewportHeight < 5 {
		viewportHeight = 5
	}
	// Account for: border (1 char each side) + border padding (2 chars each side)
	// Add small buffer for safety margin
	viewportWidth := m.width - 6 - 2
	if viewportWidth < 40 {
		viewportWidth = 40
	}
	m.logViewport = viewport.New(viewportWidth, viewportHeight)
	// Wrap initial log content to viewport width
	wrappedLogs := wrapLogContent(m.logs, viewportWidth)
	m.logViewport.SetContent(wrappedLogs)

	// Create log channel for streaming logs
	m.logChan = make(chan string, 100)

	// Start async execution and log listener
	execCmd, startCmd := m.executeBentoAsync(m.logChan)
	return m, tea.Batch(
		startCmd, // Send start message with cancel function
		execCmd,  // Start execution
		listenForLogs(m.logChan),
	)
}
