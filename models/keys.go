package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/lipgloss/v2"
)

var helpSep = lipgloss.NewStyle().Foreground(helpKeyColor).Render("•")

func StyleKey(key key.Binding) string {
	return fmt.Sprintf(
		"%s %s",
		helpKeyStyle.Render("["+key.Help().Key+"]"),
		helpDescStyle.Render(key.Help().Desc))
}

type NavKeymap struct {
	LineUp key.Binding
	LineDn key.Binding

	PageUp key.Binding
	PageDn key.Binding

	HalfPgUp key.Binding
	HalfPgDn key.Binding

	GoToTop key.Binding
	GoToEnd key.Binding

	Quit key.Binding
}

func (f NavKeymap) String() []string {
	var out []string
	out = append(out, fmt.Sprintf("%v %v %v %v %v",
		StyleKey(f.GoToTop), helpSep,
		StyleKey(f.HalfPgUp), helpSep,
		StyleKey(f.LineUp)))
	out = append(out, fmt.Sprintf("%v %v %v %v %v",
		StyleKey(f.GoToEnd), helpSep,
		StyleKey(f.HalfPgDn), helpSep,
		StyleKey(f.LineDn)))

	return out
}

func GetNavKeymap() NavKeymap {
	m := NavKeymap{}

	// Navigation/Scrolling
	m.LineUp = key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("k/↑", "up"),
	)

	m.LineDn = key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("j/↓", "dn"),
	)

	m.PageUp = key.NewBinding(
		key.WithKeys("pgup"),
		key.WithHelp("pgup", "pgup"),
	)

	m.PageDn = key.NewBinding(
		key.WithKeys("pgdn"),
		key.WithHelp("pgdn", "pgdn"),
	)

	m.HalfPgUp = key.NewBinding(
		key.WithKeys("u", "ctrl+u"),
		key.WithHelp("u", "½ pgup"),
	)

	m.HalfPgDn = key.NewBinding(
		key.WithKeys("d", "ctrl+d"),
		key.WithHelp("d", "½ pgdn"),
	)

	m.GoToTop = key.NewBinding(
		key.WithKeys("g", "home"),
		key.WithHelp("g", "top"),
	)

	m.GoToEnd = key.NewBinding(
		key.WithKeys("G", "end"),
		key.WithHelp("G", "end"),
	)

	// Control
	m.Quit = key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/^+c", "quit"),
	)

	return m
}

type FilterKeymap struct {
	FocusFilter key.Binding
	SaveFilter  key.Binding

	NoFocusClearFilter key.Binding
	FocusedClearFilter key.Binding

	Quit key.Binding
}

func (f FilterKeymap) String(textFocused bool) string {
	mainKeys := ""
	if textFocused {
		mainKeys = fmt.Sprintf("%v %v %v", StyleKey(f.SaveFilter), helpSep, StyleKey(f.FocusedClearFilter))
	} else {
		mainKeys = fmt.Sprintf("%v %v %v", StyleKey(f.FocusFilter), helpSep, StyleKey(f.NoFocusClearFilter))
	}

	return fmt.Sprintf("%v %v %v", StyleKey(f.Quit), helpSep, mainKeys)
}

func GetFilterKeymap() FilterKeymap {
	m := FilterKeymap{}

	// Filtering
	m.FocusFilter = key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "focus filter"),
	)

	m.NoFocusClearFilter = key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "clear"),
	)

	m.FocusedClearFilter = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "clear"),
	)

	m.SaveFilter = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "accept"),
	)

	// Control
	m.Quit = key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/^+c", "quit"),
	)

	return m
}
