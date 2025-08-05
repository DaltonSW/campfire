package models

import (
	"fmt"
	"io"
	"io/fs"
	"math"
	"strings"

	"os"
	"time"

	"github.com/charmbracelet/bubbles/v2/help"
	"github.com/charmbracelet/bubbles/v2/viewport"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"

	"github.com/dustin/go-humanize"
)

const tickRate = time.Millisecond * 200

// Messages

type tickMsg time.Time
type fileGoneMsg struct{}
type fileExistsMsg struct {
	info    fs.FileInfo
	content []byte
}
type fileErrorMsg error

type VisualFilters struct {
	ShowInfo  bool
	ShowWarn  bool
	ShowError bool
	ShowDebug bool
}

// NewModel actually creates the main campfire model
func NewModel(filename string) *model {
	// Viewport is initialized in after window size message
	m := model{
		filename: filename,
		help:     help.New(),
		filters: VisualFilters{
			ShowInfo:  true,
			ShowWarn:  true,
			ShowError: true,
			ShowDebug: true,
		},
	}

	return &m
}

// model is the BubbleTea model for campfire
type model struct {
	filename      string
	content       string
	viewport      viewport.Model
	width, height int
	ready         bool
	help          help.Model

	filters VisualFilters

	fileExists   bool
	prevFileInfo fs.FileInfo
}

// Init kicks off the ticking
func (m model) Init() tea.Cmd {
	return tickCmd()
}

// Update processes new messages for the model
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			return m, tea.Quit
		case "esc":
			return m, tea.Quit

		case "1":
			m.filters.ShowInfo = !m.filters.ShowInfo
		case "2":
			m.filters.ShowWarn = !m.filters.ShowWarn
		case "3":
			m.filters.ShowError = !m.filters.ShowError
		case "4":
			m.filters.ShowDebug = !m.filters.ShowDebug
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.Header())
		footerHeight := lipgloss.Height(m.Footer())
		verticalMarginHeight := headerHeight + footerHeight

		m.width = msg.Width - appStyle.GetHorizontalFrameSize()
		m.height = msg.Height - appStyle.GetVerticalFrameSize()

		m.help.Width = m.width

		viewportStyle = viewportStyle.Width(m.width).Height(m.height - verticalMarginHeight)
		footerStyle.Width(m.width)

		if !m.ready {
			m.viewport = viewport.New(
				viewport.WithWidth(msg.Width-viewportStyle.GetHorizontalBorderSize()),
				viewport.WithHeight(m.height-verticalMarginHeight-viewportStyle.GetVerticalBorderSize()),
			)

			m.viewport.SoftWrap = true

			m.fileExists = false
			m.ready = true
		} else {
			m.viewport.SetWidth(msg.Width)
			m.viewport.SetHeight(msg.Height - verticalMarginHeight)
		}

	case fileExistsMsg:
		m.prevFileInfo = msg.info
		m.fileExists = true
		var outContent []string
		for i, message := range strings.Split(string(msg.content), "\n") {
			styled := StyleMessage(message, i, m.filters)
			if styled != "" {
				outContent = append(outContent, styled)
			}
		}
		m.viewport.SetContentLines(outContent)

	case fileGoneMsg:
		m.fileExists = false
		m.viewport.SetContent("")

	case fileErrorMsg:
		content := "❌ Error reading file: " + msg.Error()
		m.viewport.SetContent(content)

	case tickMsg:
		cmds = append(cmds, checkFile(m.filename))
		cmds = append(cmds, tickCmd())
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View displays the state of the model
func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	var centerContent string

	centerContent = m.viewport.View()

	return appStyle.Render(fmt.Sprintf("%s\n%s\n%s", m.Header(), viewportStyle.Render(centerContent), m.Footer()))
}

// Header gets the above-viewport content. Title and file stats
func (m model) Header() string {
	cWidth := len(m.filename)
	newLWidth := int(math.Floor(float64((m.width - cWidth) / 2)))

	cContent := titleStyle.AlignHorizontal(lipgloss.Right).Width(newLWidth + cWidth).Render("Campfire")

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
	rContent = statsStyle.Width(m.width - cWidth - newLWidth).Render(rContent)

	return cContent + rContent
}

var visibleIcon = lipgloss.NewStyle().Foreground(lipgloss.Color("#a6da95")).Render("")
var invisibleIcon = lipgloss.NewStyle().Foreground(lipgloss.Color("#ed8796")).Render("")

// Footer prints the helptext and contact/repo info
func (m model) Footer() string {

	if !m.ready || m.width == 0 {
		return ""
	}

	var lContent, cContent, rContent []string

	// TODO: Help

	infoIcon := ternary(m.filters.ShowInfo, visibleIcon, invisibleIcon)
	warnIcon := ternary(m.filters.ShowWarn, visibleIcon, invisibleIcon)
	errorIcon := ternary(m.filters.ShowError, visibleIcon, invisibleIcon)
	debugIcon := ternary(m.filters.ShowDebug, visibleIcon, invisibleIcon)

	lContent = append(lContent, fmt.Sprintf("[1] INFO (%v) | [3] ERROR (%v)", infoIcon, errorIcon))
	lContent = append(lContent, fmt.Sprintf("[2] WARN (%v) | [4] DEBUG (%v)", warnIcon, debugIcon))
	lContent = append(lContent, "Text Filter: <not yet implemented>")

	cContent = append(cContent, "    [^+u] 󱦒 pgup | [k/] up | [j/] down | [^+d] 󱦒 pgdn")
	cContent = append(cContent, "[1-4] filter log level | [ctrl+/] text filter")
	cContent = append(cContent, "[q/ctrl+c] quit")

	rContent = append(rContent, "Issues? Suggestions?")
	rContent = append(rContent, "Lemme know! feedback@dalton.dog")
	rContent = append(rContent, "https://github.com/daltonsw/campfire")

	var outContent string
	for i := range cContent {
		outContent += align(m.width, lContent[i], cContent[i], rContent[i]) + "\n"
	}

	return outContent
}

// ~~ Commands ~~

// tickCmd will send the same tick on a constant cadence
func tickCmd() tea.Cmd {
	return tea.Tick(tickRate, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// checkFile checks the current state of the file, returning a corresponding message
func checkFile(name string) tea.Cmd {
	return func() tea.Msg {
		info, err := os.Stat(name)

		// File doesn't exist
		if os.IsNotExist(err) {
			return fileGoneMsg{}
		}

		// File exists but error trying to access it
		if err != nil {
			return fileErrorMsg(err)
		}

		// Otherwise, open file
		file, err := os.Open(name)
		if err != nil {
			return fileErrorMsg(err)
		}

		// Can close the file at the end of this since we'll extract all the content prior
		defer file.Close()

		// Grab all the content, return it in a message
		content, err := io.ReadAll(file)
		if err != nil {
			return fileErrorMsg(err)
		}

		return fileExistsMsg{content: content, info: info}
	}
}
