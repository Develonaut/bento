package shared

import "github.com/charmbracelet/lipgloss"

// LayoutHelper provides common layout calculation utilities
type LayoutHelper struct {
	Width  int
	Height int
}

// NewLayoutHelper creates a new layout helper with terminal dimensions
func NewLayoutHelper(width, height int) LayoutHelper {
	return LayoutHelper{
		Width:  width,
		Height: height,
	}
}

// ContentHeight calculates available height after subtracting header and footer
func (l LayoutHelper) ContentHeight(headerHeight, footerHeight int) int {
	return l.Height - headerHeight - footerHeight
}

// ContentWidth calculates available width after subtracting margins
func (l LayoutHelper) ContentWidth(leftMargin, rightMargin int) int {
	return l.Width - leftMargin - rightMargin
}

// SplitHorizontal splits available width between left and right sections
// Returns (leftWidth, rightWidth)
func (l LayoutHelper) SplitHorizontal(leftRatio float64) (int, int) {
	if leftRatio < 0 {
		leftRatio = 0
	}
	if leftRatio > 1 {
		leftRatio = 1
	}

	leftWidth := int(float64(l.Width) * leftRatio)
	rightWidth := l.Width - leftWidth

	return leftWidth, rightWidth
}

// SplitVertical splits available height between top and bottom sections
// Returns (topHeight, bottomHeight)
func (l LayoutHelper) SplitVertical(topRatio float64) (int, int) {
	if topRatio < 0 {
		topRatio = 0
	}
	if topRatio > 1 {
		topRatio = 1
	}

	topHeight := int(float64(l.Height) * topRatio)
	bottomHeight := l.Height - topHeight

	return topHeight, bottomHeight
}

// FitContent calculates content size after subtracting rendered component sizes
func (l LayoutHelper) FitContent(renderedComponents ...string) (width, height int) {
	totalHeight := 0
	maxWidth := 0

	for _, component := range renderedComponents {
		h := lipgloss.Height(component)
		w := lipgloss.Width(component)

		totalHeight += h
		if w > maxWidth {
			maxWidth = w
		}
	}

	return l.Width - maxWidth, l.Height - totalHeight
}

// RemainingHeight calculates height remaining after subtracting rendered components
func (l LayoutHelper) RemainingHeight(renderedComponents ...string) int {
	totalHeight := 0
	for _, component := range renderedComponents {
		totalHeight += lipgloss.Height(component)
	}

	remaining := l.Height - totalHeight
	if remaining < 0 {
		return 0
	}
	return remaining
}

// RemainingWidth calculates width remaining after subtracting rendered components
func (l LayoutHelper) RemainingWidth(renderedComponents ...string) int {
	totalWidth := 0
	for _, component := range renderedComponents {
		totalWidth += lipgloss.Width(component)
	}

	remaining := l.Width - totalWidth
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Clamp ensures a value stays within min and max bounds
func Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// MaxInt returns the larger of two integers
func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// MinInt returns the smaller of two integers
func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
