package models

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/dustin/go-humanize"
)

// Header gets the above-viewport content. Title and file stats
func (m model) Header() string {
	cContent := titleStyle.Render("Campfire")

	rContent := ""
	if m.fileExists {
		filesize := humanize.Bytes(uint64(m.prevFileInfo.Size()))

		rContent = fmt.Sprintf(
			"%v %v",
			fileNameStyle.Render(m.filename),
			fmt.Sprintf("(Size: %v)", filesize),
		)
	} else {
		rContent = "File not found..."
	}
	rContent = statsStyle.Render(rContent)

	lContent := statsStyle.Italic(true).Render("Feedback: github.com/daltonsw/campfire")

	return align(m.width, lContent, cContent, rContent)
}

var visibleIcon = lipgloss.NewStyle().Foreground(lipgloss.Color("#a6da95")).Render("")
var invisibleIcon = lipgloss.NewStyle().Foreground(lipgloss.Color("#ed8796")).Render("")

// Footer prints the helptext and contact/repo info
func (m model) Footer() string {

	if !m.ready || m.width == 0 {
		return ""
	}

	var lContent, rContent []string

	infoIcon := ternary(m.filters.ShowInfo, visibleIcon, invisibleIcon)
	warnIcon := ternary(m.filters.ShowWarn, visibleIcon, invisibleIcon)
	errorIcon := ternary(m.filters.ShowError, visibleIcon, invisibleIcon)
	debugIcon := ternary(m.filters.ShowDebug, visibleIcon, invisibleIcon)
	fatalIcon := ternary(m.filters.ShowFatal, visibleIcon, invisibleIcon)
	otherIcon := ternary(m.filters.ShowOther, visibleIcon, invisibleIcon)

	lContent = append(lContent, "Text Filter: <not yet implemented>")
	lContent = append(lContent, fmt.Sprintf(
		"[1] INFO (%v) | [3] ERROR (%v) | [5] FATAL (%v)",
		infoIcon, errorIcon, fatalIcon,
	))
	lContent = append(lContent, fmt.Sprintf(
		"[2] WARN (%v) | [4] DEBUG (%v) | [6] OTHER (%v)",
		warnIcon, debugIcon, otherIcon,
	))

	rContent = append(rContent, "[^+u] 󱦒 pgup | [k/] up | [j/] down | [^+d] 󱦒 pgdn")
	rContent = append(rContent, "[1-6] log level | [ctrl+f] text filter | [q/ctrl+c] quit")
	rContent = append(rContent, "")

	var outContent string
	for i := range lContent {
		outContent += align(m.width, lContent[i], "", rContent[i]) + "\n"
	}

	return outContent
}
