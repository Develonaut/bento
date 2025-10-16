package components

import (
	"fmt"
	"strings"
	"time"

	"bento/pkg/omise/styles"

	"github.com/charmbracelet/lipgloss"
)

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
	line := s.buildStepLine(indent, icon, step)
	return s.styleStepLine(line, step.Status)
}

// buildStepLine constructs the step line with icon, name, and duration
func (s Sequence) buildStepLine(indent, icon string, step Step) string {
	line := fmt.Sprintf("%s%s %s", indent, icon, step.Name)

	// Show type for running steps
	if step.Status == StepRunning && step.Type != "" {
		line = fmt.Sprintf("%s (%s)", line, step.Type)
	}

	// Show duration for completed/failed steps
	if step.Status == StepCompleted || step.Status == StepFailed {
		if step.Duration > 0 {
			durationStr := step.Duration.Round(time.Millisecond).String()
			line = fmt.Sprintf("%s (%s)", line, durationStr)
		}
	}

	return line
}

// styleStepLine applies styling based on step status
func (s Sequence) styleStepLine(line string, status StepStatus) string {
	switch status {
	case StepCompleted:
		return styles.SuccessStyle.Render(line)
	case StepFailed:
		return styles.ErrorStyle.Render(line)
	case StepRunning:
		return styles.Selected.Render(line)
	default:
		return styles.Subtle.Render(line)
	}
}

// getStepIcon returns the icon for a step status
func (s Sequence) getStepIcon(status StepStatus) string {
	switch status {
	case StepRunning:
		return s.spinner.View() // Animated spinner
	case StepCompleted:
		return "✓" // Success checkmark
	case StepFailed:
		return "✗" // Failure X
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
