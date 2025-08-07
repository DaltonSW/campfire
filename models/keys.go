package models

import (
	"github.com/charmbracelet/bubbles/v2/key"
)

func (k Keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Quit,
		k.LineUp, k.LineDn,
		k.HalfPgUp, k.HalfPgDn,
		k.GoToTop, k.GoToEnd,
		k.FocusFilter, k.NoFocusClearFilter,
		k.SaveFilter, k.FocusedClearFilter,
	}
}

func (k Keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

type Keymap struct {
	LineUp key.Binding
	LineDn key.Binding

	PageUp key.Binding
	PageDn key.Binding

	HalfPgUp key.Binding
	HalfPgDn key.Binding

	GoToTop key.Binding
	GoToEnd key.Binding

	FocusFilter key.Binding
	SaveFilter  key.Binding

	NoFocusClearFilter key.Binding
	FocusedClearFilter key.Binding

	ToggleInfo  key.Binding
	ToggleWarn  key.Binding
	ToggleError key.Binding
	ToggleDebug key.Binding
	ToggleFatal key.Binding
	ToggleOther key.Binding

	Quit key.Binding
}

func GetKeymap() Keymap {
	m := Keymap{}

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
		key.WithKeys("pgdown"),
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
	m.FocusedClearFilter.SetEnabled(false)

	m.SaveFilter = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "accept"),
	)
	m.SaveFilter.SetEnabled(false)

	m.ToggleInfo = key.NewBinding(key.WithKeys("1"))
	m.ToggleWarn = key.NewBinding(key.WithKeys("2"))
	m.ToggleError = key.NewBinding(key.WithKeys("3"))
	m.ToggleDebug = key.NewBinding(key.WithKeys("4"))
	m.ToggleFatal = key.NewBinding(key.WithKeys("5"))
	m.ToggleOther = key.NewBinding(key.WithKeys("6"))

	// Control
	m.Quit = key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	)

	return m
}
