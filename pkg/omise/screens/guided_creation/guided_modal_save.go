package guided_creation

import (
	"encoding/json"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"bento/pkg/neta"
)

func (m *GuidedModal) saveBento() tea.Cmd {
	return func() tea.Msg {
		// Recover from panics during save to ensure terminal is restored
		var panicMsg *GuidedCompleteMsg
		defer func() {
			if r := recover(); r != nil {
				panicMsg = &GuidedCompleteMsg{
					Success:   false,
					Err:       fmt.Errorf("panic during save: %v", r),
					Cancelled: false,
				}
			}
		}()

		// Validate the bento before saving
		if err := m.validator.Validate(*m.definition); err != nil {
			return GuidedCompleteMsg{
				Success:   false,
				Err:       err,
				Cancelled: false,
			}
		}

		// Use original filename if editing, otherwise generate from name
		var filename string
		if m.editing && m.originalFilename != "" {
			filename = m.originalFilename
		} else {
			filename = strings.ReplaceAll(strings.ToLower(m.definition.Name), " ", "-")
		}

		// If editing, check if content actually changed before incrementing version
		if m.editing && !m.saveInProgress {
			// Load the original bento to compare
			originalDef, err := m.store.Load(filename)
			if err == nil {
				// Compare JSON representations (excluding version field)
				if !bentoContentChanged(originalDef, *m.definition) {
					// No changes - don't increment version
					m.definition.Version = originalDef.Version
				} else {
					// Content changed - increment version
					m.definition.Version = incrementVersion(originalDef.Version)
				}
			} else {
				// Couldn't load original (shouldn't happen), increment anyway
				m.definition.Version = incrementVersion(m.definition.Version)
			}
			m.saveInProgress = true
		}

		// Save to store
		if err := m.store.Save(filename, *m.definition); err != nil {
			return GuidedCompleteMsg{
				Success:   false,
				Err:       err,
				Cancelled: false,
			}
		}

		// Check if panic occurred and return panic message instead
		if panicMsg != nil {
			return *panicMsg
		}

		return GuidedCompleteMsg{
			Success:    true,
			Definition: m.definition,
			Cancelled:  false,
		}
	}
}

// bentoContentChanged compares two bentos ignoring version field
func bentoContentChanged(original, updated neta.Definition) bool {
	// Create copies with version normalized to compare content only
	origCopy := original
	updatedCopy := updated
	origCopy.Version = "0.0"
	updatedCopy.Version = "0.0"

	// Marshal to JSON for comparison
	origJSON, err1 := json.Marshal(origCopy)
	updatedJSON, err2 := json.Marshal(updatedCopy)

	if err1 != nil || err2 != nil {
		// If marshaling fails, assume changed
		return true
	}

	// Compare JSON strings
	return string(origJSON) != string(updatedJSON)
}

// incrementVersion increments a semantic version string (e.g., "1.0" -> "1.1")
func incrementVersion(version string) string {
	// Simple version increment: split on "." and increment the last number
	parts := strings.Split(version, ".")
	if len(parts) == 0 {
		return "1.0"
	}

	// Try to parse and increment the last part
	lastIdx := len(parts) - 1
	var lastNum int
	if _, err := fmt.Sscanf(parts[lastIdx], "%d", &lastNum); err == nil {
		lastNum++
		parts[lastIdx] = fmt.Sprintf("%d", lastNum)
		return strings.Join(parts, ".")
	}

	// If parsing failed, just append .1
	return version + ".1"
}
