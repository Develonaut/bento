// Package miso provides terminal output "seasoning" - themed styling and progress display.
//
// Step sequence rendering with status words and sushi emojis.
package miso

import (
	"fmt"
	"hash/fnv"
	"strings"
	"time"
)

// Approved sushi emojis (from EMOJIS.md)
var sushiEmojis = []string{"🍣", "🍙", "🥢", "🍥"}

// Approved status words (from STATUS_WORDS.md)
var runningWords = []string{
	"Tasting", "Sampling", "Trying", "Enjoying",
	"Devouring", "Nibbling", "Savoring", "Testing",
}

var completedWords = []string{
	"Savored", "Devoured", "Enjoyed", "Relished",
	"Finished", "Consumed", "Completed", "Perfected",
}

var failedWords = []string{
	"Spoiled", "Burnt", "Dropped", "Ruined",
	"Failed", "Overcooked", "Undercooked",
}

// StepStatus represents the execution status of a step.
type StepStatus int

const (
	StepPending StepStatus = iota
	StepRunning
	StepCompleted
	StepFailed
)

// Step represents a single step in a sequence.
type Step struct {
	Name     string
	Type     string
	Status   StepStatus
	Duration time.Duration
	Depth    int // Nesting level for indentation
}

// Sequence displays a list of execution steps with status indicators.
type Sequence struct {
	steps []Step
}

// NewSequence creates a new sequence display.
func NewSequence() *Sequence {
	return &Sequence{
		steps: []Step{},
	}
}

// AddStep adds a step with depth 0.
func (s *Sequence) AddStep(name, nodeType string) {
	s.AddStepWithDepth(name, nodeType, 0)
}

// AddStepWithDepth adds a step with specified nesting depth.
func (s *Sequence) AddStepWithDepth(name, nodeType string, depth int) {
	s.steps = append(s.steps, Step{
		Name:     name,
		Type:     nodeType,
		Status:   StepPending,
		Duration: 0,
		Depth:    depth,
	})
}

// UpdateStep updates the status of a step by name.
func (s *Sequence) UpdateStep(name string, status StepStatus) {
	for i := range s.steps {
		if s.steps[i].Name == name {
			s.steps[i].Status = status
			return
		}
	}
}

// SetDuration sets the duration for a step by name.
func (s *Sequence) SetDuration(name string, duration time.Duration) {
	for i := range s.steps {
		if s.steps[i].Name == name {
			s.steps[i].Duration = duration
			return
		}
	}
}

// GetSteps returns all steps in the sequence.
func (s *Sequence) GetSteps() []Step {
	return s.steps
}

// View renders the sequence of steps.
func (s *Sequence) View() string {
	if len(s.steps) == 0 {
		return ""
	}

	lines := []string{}
	for _, step := range s.steps {
		lines = append(lines, s.formatStep(step))
	}

	return strings.Join(lines, "\n")
}

// formatStep renders a single step with status and timing.
func (s *Sequence) formatStep(step Step) string {
	prefix := buildStepPrefix(step)
	status := buildStepStatus(step)
	suffix := buildStepSuffix(step)

	parts := []string{prefix, status, step.Name + "…"}
	if suffix != "" {
		parts = append(parts, suffix)
	}

	return strings.Join(parts, " ")
}

// buildStepPrefix creates the indent/emoji/icon prefix for a step.
func buildStepPrefix(step Step) string {
	indent := strings.Repeat("  ", step.Depth)
	icon := getStepIcon(step.Status)
	emoji := ""

	if step.Status == StepCompleted {
		emoji = getStepEmoji(step.Name)
	}

	if emoji != "" {
		if icon != "" {
			return indent + emoji + " " + icon
		}
		return indent + emoji
	}

	if icon != "" {
		return indent + icon
	}

	return indent
}

// buildStepStatus creates the colored status word.
func buildStepStatus(step Step) string {
	statusWord := getStatusLabel(step.Status, step.Name)
	return colorStatusWord(statusWord, step.Status)
}

// buildStepSuffix creates the duration suffix if applicable.
func buildStepSuffix(step Step) string {
	if (step.Status == StepCompleted || step.Status == StepFailed) && step.Duration > 0 {
		return fmt.Sprintf("(%s)", step.Duration.Round(time.Millisecond).String())
	}
	return ""
}

// getStepEmoji returns a deterministic sushi emoji based on step name.
// Uses FNV-1a hash to ensure the same step always gets the same emoji.
func getStepEmoji(stepName string) string {
	h := fnv.New32a()
	h.Write([]byte(stepName))
	hash := h.Sum32()
	return sushiEmojis[hash%uint32(len(sushiEmojis))]
}

// colorStatusWord colors the status word based on status.
// Uses the default Maguro theme for now.
// TODO: Accept theme as parameter when integrating with daemon-combo.
func colorStatusWord(word string, status StepStatus) string {
	// Build default theme for standalone sequence usage
	theme := BuildTheme(GetPalette(VariantMaguro))

	switch status {
	case StepRunning:
		return theme.Title.Render(word) // Primary color
	case StepCompleted:
		return theme.Success.Render(word)
	case StepFailed:
		return theme.Error.Render(word)
	default:
		return theme.Subtle.Render(word)
	}
}

// getStatusLabel returns a fun varied status word based on step name.
// Uses deterministic hash to ensure same step gets same word.
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

// getStepIcon returns the icon for a step status.
func getStepIcon(status StepStatus) string {
	switch status {
	case StepRunning:
		return "⟳" // Loading spinner
	case StepCompleted:
		return "" // No icon - rely on emoji and colors
	case StepFailed:
		return "❌" // Red X for failures
	default:
		return "•" // Pending dot
	}
}
