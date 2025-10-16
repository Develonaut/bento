package components

import (
	"fmt"
	"hash/fnv"
	"strings"
	"time"

	"bento/pkg/omise/styles"

	"github.com/charmbracelet/lipgloss"
)

// Sushi-themed emojis for nodes
var sushiEmojis = []string{"🍣", "🍙", "🥢", "🍥"}

// Fun status word variations
var runningWords = []string{"Tasting", "Sampling", "Trying", "Enjoying", "Devouring", "Nibbling", "Savoring", "Testing"}
var completedWords = []string{"Savored", "Devoured", "Enjoyed", "Relished", "Finished", "Consumed", "Completed", "Perfected"}
var failedWords = []string{"Spoiled", "Burnt", "Dropped", "Ruined", "Failed", "Overcooked", "Undercooked"}

// StepStatus represents the execution status of a step
type StepStatus int

const (
	StepPending StepStatus = iota
	StepRunning
	StepCompleted
	StepFailed
)

// Step represents a single step in a sequence
type Step struct {
	Name     string
	Type     string
	Status   StepStatus
	Duration time.Duration
	Depth    int // Nesting level for indentation
}

// Sequence displays a list of execution steps with status indicators
type Sequence struct {
	spinner Spinner
	steps   []Step
}

// NewSequence creates a new sequence display
func NewSequence() Sequence {
	return Sequence{
		spinner: NewSpinner(),
		steps:   []Step{},
	}
}

// SetSteps updates the sequence with a list of steps
func (s Sequence) SetSteps(steps []Step) Sequence {
	s.steps = steps
	return s
}

// View renders the sequence of steps
func (s Sequence) View() string {
	if len(s.steps) == 0 {
		return ""
	}

	lines := []string{}
	for _, step := range s.steps {
		lines = append(lines, s.formatStep(step))
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// formatStep renders a single step with status and timing
func (s Sequence) formatStep(step Step) string {
	indent := strings.Repeat("  ", step.Depth)
	icon := s.getStepIcon(step.Status)

	// Get emoji (only for completed, not failed)
	var emoji string
	if step.Status == StepCompleted {
		emoji = getStepEmoji(step.Name)
	}

	// Get colored status word
	statusWord := getStatusLabel(step.Status, step.Name)
	coloredStatus := s.colorStatusWord(statusWord, step.Status)

	// Build parts
	var parts []string

	// Add indent and emoji (if present)
	if emoji != "" {
		parts = append(parts, indent+emoji)
		if icon != "" {
			parts = append(parts, icon)
		}
	} else {
		if icon != "" {
			parts = append(parts, indent+icon)
		} else {
			parts = append(parts, indent)
		}
	}

	// Add colored status and name
	parts = append(parts, coloredStatus, step.Name+"…")

	// Add duration only for completed or failed steps
	if (step.Status == StepCompleted || step.Status == StepFailed) && step.Duration > 0 {
		durationStr := step.Duration.Round(time.Millisecond).String()
		parts = append(parts, fmt.Sprintf("(%s)", durationStr))
	}

	return strings.Join(parts, " ")
}

// getStepEmoji returns a deterministic sushi emoji based on step name
// Uses FNV-1a hash to ensure the same step always gets the same emoji
func getStepEmoji(stepName string) string {
	h := fnv.New32a()
	h.Write([]byte(stepName))
	hash := h.Sum32()
	return sushiEmojis[hash%uint32(len(sushiEmojis))]
}

// colorStatusWord colors the status word based on status
func (s Sequence) colorStatusWord(word string, status StepStatus) string {
	switch status {
	case StepRunning:
		// Use Primary color for running status
		return lipgloss.NewStyle().Foreground(styles.Primary).Bold(true).Render(word)
	case StepCompleted:
		return styles.SuccessStyle.Render(word)
	case StepFailed:
		return styles.ErrorStyle.Render(word)
	default:
		return styles.Subtle.Render(word)
	}
}

// getStatusLabel returns a fun varied status word based on step name
// Uses deterministic hash to ensure same step gets same word
func getStatusLabel(status StepStatus, stepName string) string {
	h := fnv.New32a()
	h.Write([]byte(stepName))
	hash := h.Sum32()

	switch status {
	case StepPending:
		return "Preparing"
	case StepRunning:
		return runningWords[hash%uint32(len(runningWords))]
	case StepCompleted:
		return completedWords[hash%uint32(len(completedWords))]
	case StepFailed:
		return failedWords[hash%uint32(len(failedWords))]
	default:
		return "Preparing"
	}
}

// getStepIcon returns the icon for a step status
func (s Sequence) getStepIcon(status StepStatus) string {
	switch status {
	case StepRunning:
		return s.spinner.View() // Animated spinner
	case StepCompleted:
		return "" // No icon - rely on emoji and colors
	case StepFailed:
		return "❌" // Red X for failures
	default:
		return "•" // Pending dot
	}
}

// RebuildStyles updates the sequence component with current theme
func (s Sequence) RebuildStyles() Sequence {
	s.spinner = s.spinner.RebuildStyles()
	return s
}

// Update handles sequence messages (for spinner animation)
func (s Sequence) Update(spinner Spinner) Sequence {
	s.spinner = spinner
	return s
}
