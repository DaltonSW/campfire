package models

import (
	"fmt"
	"io"
	"io/fs"

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
			m.viewport.LeftGutterFunc = GutterFunc

			m.fileExists = false
			m.ready = true
		} else {
			m.viewport.SetWidth(msg.Width)
			m.viewport.SetHeight(msg.Height - verticalMarginHeight)
		}

	case fileExistsMsg:
		m.viewport.LeftGutterFunc = GutterFunc
		m.prevFileInfo = msg.info
		m.fileExists = true
		m.viewport.SetContent(string(msg.content))
		m.viewport.GotoTop()

	case fileGoneMsg:
		m.viewport.LeftGutterFunc = nil
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
	outStr := fmt.Sprintf("Current File: %v\n", fileNameStyle.Render(m.filename))
	if m.fileExists {
		outStr += fmt.Sprintf("File Size: %v ~~ Last Modified: %v", humanize.Bytes(uint64(m.prevFileInfo.Size())), m.prevFileInfo.ModTime().Format("03:04:05.0000 PM"))
	} else {
		outStr += "File being monitored doesn't seem to exist..."
	}
	return lipgloss.PlaceHorizontal(m.width, lipgloss.Center, outStr)

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
