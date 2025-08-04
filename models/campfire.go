package models

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"

	"os"
	"time"

	"github.com/charmbracelet/bubbles/v2/viewport"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/log"
)

var (
	// titleStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#DDDDDD")).Bold(true).Underline(true)
	fileNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#DD8800")).Italic(true)
	viewportStyle = lipgloss.NewStyle().Align(lipgloss.Center, lipgloss.Top).Border(lipgloss.RoundedBorder())
	footerStyle   = lipgloss.NewStyle().AlignHorizontal(lipgloss.Center).Italic(true)
)

func NewModel(filename string) *model {
	// Viewport is initialized in after window size message
	m := model{
		filename: filename,
	}

	log.Info("NewModel function")

	return &m
}

func (m model) CloseModel() {
	if m.file != nil {
		m.file.Close()
	}
}

type model struct {
	filename      string
	content       string
	viewport      viewport.Model
	width, height int
	ready         bool

	file   *os.File
	reader *bufio.Reader

	prevFileInfo fs.FileInfo
}

func (m model) Init() tea.Cmd {
	// TODO: Try and load file here so there isn't a "trying to find file" visible every time
	return tickCmd()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.header())
		footerHeight := lipgloss.Height(m.footer())
		verticalMarginHeight := headerHeight + footerHeight

		m.width = msg.Width
		m.height = msg.Height

		viewportStyle = viewportStyle.Width(msg.Width).Height(m.height - verticalMarginHeight)
		footerStyle.Width(m.width)

		if !m.ready {
			m.viewport = viewport.New(viewport.WithWidth(msg.Width), viewport.WithHeight(m.height-verticalMarginHeight-viewportStyle.GetVerticalBorderSize()))
			m.viewport.SetContent("Waiting for file to be found...")
			m.ready = true
		} else {
			m.viewport.SetWidth(msg.Width)
			m.viewport.SetHeight(msg.Height - verticalMarginHeight)
		}

	case fileExistsMsg:
		m.viewport.SetContent(string(msg.content))
		m.viewport.GotoTop()

	case fileGoneMsg:
		content := "❌ File not found: " + m.filename
		content += "\n\n" + "Waiting for file to be created..."
		m.viewport.SetContent(content)

	case fileErrorMsg:
		content := "❌ Error reading file: " + msg.Error()
		m.viewport.SetContent(content)

	case tickMsg:
		// log.Info("Tick msg")
		cmds = append(cmds, checkFile(m.filename))
		cmds = append(cmds, tickCmd())
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	var centerContent string

	centerContent = m.viewport.View()

	return fmt.Sprintf("%s\n%s\n%s", m.header(), viewportStyle.Render(centerContent), m.footer())
}

func (m model) header() string {
	outStr := "Current File: " + fileNameStyle.Render(m.filename)
	return lipgloss.PlaceHorizontal(m.width, lipgloss.Center, outStr)
}
func (m model) footer() string {
	outStr := footerStyle.Render("Eventual helptext location and whatever")
	return lipgloss.PlaceHorizontal(m.width, lipgloss.Center, outStr)
}

const tickRate = time.Millisecond * 500

type tickMsg time.Time
type fileReadyMsg *os.File
type errorMsg error
type newLineMsg string

type fileGoneMsg struct{}
type fileExistsMsg struct {
	info    fs.FileInfo
	content []byte
}
type fileErrorMsg error

// region: Commands

func tickCmd() tea.Cmd {
	return tea.Tick(tickRate, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

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

		// Can close the file at the end of this since we've already extracted content
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			return fileErrorMsg(err)
		}

		return fileExistsMsg{content: content, info: info}
	}
}
