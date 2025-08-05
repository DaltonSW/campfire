package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/v2/viewport"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/compat"
)

var (
	titleColor    = compat.AdaptiveColor{Light: lipgloss.Color("#dd7878"), Dark: lipgloss.Color("#f2d5cf")}
	filenameColor = compat.AdaptiveColor{Light: lipgloss.Color("#fe640b"), Dark: lipgloss.Color("#ef9f76")}
	statsColor    = compat.AdaptiveColor{}
	footerColor   = compat.AdaptiveColor{Light: lipgloss.Color("#7c7f93"), Dark: lipgloss.Color("#737994")}
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(titleColor).
			Bold(true).
			Underline(true)

	fileNameStyle = lipgloss.NewStyle().
			Foreground(filenameColor).
			Italic(true)

	viewportStyle = lipgloss.NewStyle().
			Align(lipgloss.Left, lipgloss.Top).
			Border(lipgloss.RoundedBorder())

	footerStyle = lipgloss.NewStyle().
			Foreground(footerColor).
			AlignHorizontal(lipgloss.Center).
			Italic(true)
)

func GutterFunc(info viewport.GutterContext) string {
	if info.Soft {
		return "     │ "
	}
	if info.Index >= info.TotalLines {
		return "   ~ │ "
	}
	return fmt.Sprintf("%4d │ ", info.Index+1)
}
