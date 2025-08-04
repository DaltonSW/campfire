package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/v2/viewport"
	"github.com/charmbracelet/lipgloss/v2"
)

var (
	// titleStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#DDDDDD")).Bold(true).Underline(true)
	fileNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#DD8800")).Italic(true)
	viewportStyle = lipgloss.NewStyle().Align(lipgloss.Left, lipgloss.Top).Border(lipgloss.RoundedBorder())
	footerStyle   = lipgloss.NewStyle().AlignHorizontal(lipgloss.Center).Italic(true)
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
