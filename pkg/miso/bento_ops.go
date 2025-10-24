package miso

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Develonaut/bento/pkg/itamae"
	"github.com/Develonaut/bento/pkg/neta"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// loadBentos scans configured bento home for bento files
func loadBentos() ([]list.Item, error) {
	bentoHome := LoadBentoHome()
	bentosDir := filepath.Join(bentoHome, "bentos")

	entries, err := os.ReadDir(bentosDir)
	if err != nil {
		return nil, err
	}

	var items []list.Item
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".bento.json") {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".bento.json")
		items = append(items, BentoItem{
			Name:     name,
			FilePath: filepath.Join(bentosDir, entry.Name()),
		})
	}

	return items, nil
}

// loadBentoDefinition loads a bento from file
func loadBentoDefinition(path string) (*neta.Definition, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read bento: %w", err)
	}

	var def neta.Definition
	if err := json.Unmarshal(data, &def); err != nil {
		return nil, fmt.Errorf("failed to parse bento: %w", err)
	}

	return &def, nil
}

// runBento executes the selected bento
func (m Model) runBento() (tea.Model, tea.Cmd) {
	// Read bento file
	bentoJSON, err := os.ReadFile(m.selectedBento)
	if err != nil {
		m.logs = fmt.Sprintf("Failed to read bento: %v", err)
		m.currentView = executionView
		return m, nil
	}

	// Parse metadata to check for variables
	meta, err := ParseBentoMetadata(bentoJSON)
	if err != nil {
		m.logs = fmt.Sprintf("Failed to parse bento: %v", err)
		m.currentView = executionView
		return m, nil
	}

	// If variables exist, show form
	if len(meta.Variables) > 0 {
		m.bentoVars = meta.Variables
		return m.showForm()
	}

	// No variables, go straight to execution
	return m.startExecution()
}

// executeBentoAsync runs the bento in a goroutine and returns a tea.Cmd
func (m Model) executeBentoAsync(logChan chan string) tea.Cmd {
	return func() tea.Msg {
		defer close(logChan)

		// Load bento definition
		def, err := loadBentoDefinition(m.selectedBento)
		if err != nil {
			return executionCompleteMsg{err: err}
		}

		// Create pantry
		p := createTUIPantry()

		// Create logger that writes to both file and TUI
		logFile, logger, err := createTUILogger(logChan)
		if err != nil {
			return executionCompleteMsg{err: err}
		}
		defer logFile.Close()

		// Create chef with logger (no messenger needed)
		chef := itamae.NewWithMessenger(p, logger, nil)

		// Execute with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 45*time.Minute)
		defer cancel()

		start := time.Now()
		_, err = chef.Serve(ctx, def)
		duration := time.Since(start)

		return executionCompleteMsg{
			err:      err,
			duration: duration,
		}
	}
}
