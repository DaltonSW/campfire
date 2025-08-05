package models

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/v2/viewport"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/compat"
)

var (
	titleColor    = compat.AdaptiveColor{Light: lipgloss.Color("#dd7878"), Dark: lipgloss.Color("#f2d5cf")}
	filenameColor = compat.AdaptiveColor{Light: lipgloss.Color("#fe640b"), Dark: lipgloss.Color("#ef9f76")}
	statsColor    = compat.AdaptiveColor{Light: lipgloss.Color("#7c7f93"), Dark: lipgloss.Color("#737994")}
	footerColor   = compat.AdaptiveColor{Light: lipgloss.Color("#7c7f93"), Dark: lipgloss.Color("#737994")}

	infoColor  = compat.AdaptiveColor{Light: lipgloss.Color("#40a02b"), Dark: lipgloss.Color("#a6d189")}
	warnColor  = compat.AdaptiveColor{Light: lipgloss.Color("#df8e1d"), Dark: lipgloss.Color("#e5c890")}
	errorColor = compat.AdaptiveColor{Light: lipgloss.Color("#d20f39"), Dark: lipgloss.Color("#e78284")}
	debugColor = compat.AdaptiveColor{Light: lipgloss.Color("#8839ef"), Dark: lipgloss.Color("#ca9ee6")}
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(titleColor).
			AlignHorizontal(lipgloss.Center).
			Bold(true).
			Underline(true)

	fileNameStyle = lipgloss.NewStyle().
			Foreground(filenameColor).
			AlignHorizontal(lipgloss.Center).
			Italic(true)

	statsStyle = lipgloss.NewStyle().
			Foreground(statsColor).
			AlignHorizontal(lipgloss.Right).
			Italic(true)

	viewportStyle = lipgloss.NewStyle().
			Align(lipgloss.Left, lipgloss.Top).
			Border(lipgloss.RoundedBorder())

	footerStyle = lipgloss.NewStyle().
			Foreground(footerColor).
			AlignHorizontal(lipgloss.Center).
			Italic(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(infoColor)

	warnStyle = lipgloss.NewStyle().
			Foreground(warnColor).
			Italic(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	debugStyle = lipgloss.NewStyle().
			Foreground(debugColor)
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

func StyleMessage(line string) string {
	switch {
	case strings.Contains(line, "INFO"):
		return infoStyle.Render(line)
	case strings.Contains(line, "WARN"):
		return warnStyle.Render(line)
	case strings.Contains(line, "ERRO"):
		return errorStyle.Render(line)
	case strings.Contains(line, "DEBU"):
		return debugStyle.Render(line)
	default:
		return line
	}
}
