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
	"github.com/charmbracelet/log"

	"github.com/dustin/go-humanize"
)

const tickRate = time.Millisecond * 500

// Messages

type tickMsg time.Time
type fileGoneMsg struct{}
type fileExistsMsg struct {
	info    fs.FileInfo
	content []byte
}
type fileErrorMsg error

// NewModel actually creates the main campfire model
func NewModel(filename string) *model {
	// Viewport is initialized in after window size message
	m := model{
		filename: filename,
		help:     help.New(),
	}

	log.Info("NewModel function")

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
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.Header())
		footerHeight := lipgloss.Height(m.Footer())
		verticalMarginHeight := headerHeight + footerHeight

		m.width = msg.Width
		m.height = msg.Height

		m.help.Width = msg.Width

		viewportStyle = viewportStyle.Width(msg.Width).Height(m.height - verticalMarginHeight)
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
		for message := range strings.SplitSeq(string(msg.content), "\n") {
			outContent = append(outContent, StyleMessage(message))
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

	return fmt.Sprintf("%s\n%s\n%s", m.Header(), viewportStyle.Render(centerContent), m.Footer())
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

	// cContent := fileNameStyle.Render(m.filename)
	// cContent = lContent + cContent[lWidth:]
	// cContent = cContent[0:lipgloss.Width(cContent)-rWidth-1] + rContent

	return cContent + rContent
}

// Footer prints the helptext and contact/repo info
func (m model) Footer() string {
	// TODO: Help

	// outStr := footerStyle.Render(m.help.View(viewport.DefaultKeyMap()) + "\n")
	outStr := footerStyle.Render("Issues? Suggestions? Discussions? Lemme know -- https://github.com/daltonsw/campfire")
	return lipgloss.PlaceHorizontal(m.width, lipgloss.Center, outStr)
}

// Commands
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
