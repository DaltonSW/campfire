package models

import "github.com/charmbracelet/lipgloss/v2"

// Ternary emulates a ternary operator.
// If condition is true, returns ifTrue, otherwise returns ifFalse
func ternary(condition bool, ifTrue, ifFalse string) string {
	if condition {
		return ifTrue
	}

	return ifFalse
}

// Aligns the 3 pieces of text inline with each other. Notably, ensures that 'center'
// is aligned within the total width and not just between the other two elements
func align(width int, left, center, right string) string {
	cWidth := lipgloss.Width(center)

	// Create a styled string for the center, taking up its true width
	styledCenter := lipgloss.NewStyle().Width(cWidth).Render(center)

	// Determine proper widths for other two
	leftPaddingForCenter := (width - cWidth) / 2
	rightPaddingForCenter := width - cWidth - leftPaddingForCenter

	// Render out other sides
	renderedLeft := lipgloss.NewStyle().Width(leftPaddingForCenter).Align(lipgloss.Left).Render(left)
	renderedRight := lipgloss.NewStyle().Width(rightPaddingForCenter).Align(lipgloss.Right).Render(right)

	// Combine together and return
	return renderedLeft + styledCenter + renderedRight
}
