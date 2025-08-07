package models

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/v2"
)

const LevelFilterString = "[1] INFO %v | [2] WARN %v | [3] ERROR %v | [4] DEBUG %v| [5] FATAL %v"

var visibleIcon = lipgloss.NewStyle().Foreground(lipgloss.Color("#a6da95")).Render("✔")
var invisibleIcon = lipgloss.NewStyle().Foreground(lipgloss.Color("#ed8796")).Render("✘")

var borderStyle = lipgloss.NewStyle().AlignHorizontal(lipgloss.Center).Border(lipgloss.RoundedBorder()).PaddingLeft(1)

// Footer prints the helptext and contact/repo info
func (m model) Footer() string {

	if !m.ready || m.width == 0 {
		return ""
	}

	infoIcon := ternary(m.filters.ShowInfo, visibleIcon, invisibleIcon)
	warnIcon := ternary(m.filters.ShowWarn, visibleIcon, invisibleIcon)
	errorIcon := ternary(m.filters.ShowError, visibleIcon, invisibleIcon)
	debugIcon := ternary(m.filters.ShowDebug, visibleIcon, invisibleIcon)
	fatalIcon := ternary(m.filters.ShowFatal, visibleIcon, invisibleIcon)

	levelFilter := borderStyle.Render(
		fmt.Sprintf(LevelFilterString,
			infoIcon, warnIcon, errorIcon, debugIcon, fatalIcon,
		) + " | " + m.textInput.View())

	// outContent = lipgloss.JoinHorizontal(lipgloss.Center, levelFilter, borderStyle.Render(m.textInput.View()))

	return levelFilter + "\n" + m.help.ShortHelpView(m.keys.ShortHelp())
}
