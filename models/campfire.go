package models

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/v2/viewport"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
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

	// Things for if no file exists
	fileExists bool
	// fileSpinner spinner.Model
}

func (m model) Init() tea.Cmd { return waitForFileCmd(m.filename) }

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

	case fileReadyMsg:
		m.file = msg
		m.reader = bufio.NewReader(m.file)
		m.fileExists = true
		m.viewport.SetContent("")
		cmds = append(cmds, tailFileCmd(m.reader, m.file))

	case newLineMsg:
		m.viewport.SetContent(m.viewport.GetContent() + string(msg))
		m.viewport.GotoBottom()
		cmds = append(cmds, tailFileCmd(m.reader, m.file))

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

	if m.fileExists {
		centerContent = m.viewport.View()
	} else {
		centerContent = "Waiting for file to exist..."
	}

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

// region: Commands

// waitForFileCmd waits for a file to exist and then sends a fileReadyMsg.
func waitForFileCmd(path string) tea.Cmd {
	return func() tea.Msg {
		for {
			f, err := os.Open(path)
			if err == nil {
				return fileReadyMsg(f)
			}
			if !os.IsNotExist(err) {
				return errorMsg(err)
			}
			time.Sleep(1 * time.Second)
		}
	}
}

// tailFileCmd reads one line from the file and sends a newLineMsg.
// It's designed to be called repeatedly in a command loop.
func tailFileCmd(reader *bufio.Reader, file *os.File) tea.Cmd {
	return func() tea.Msg {
		line, err := reader.ReadString('\n')

		if err != nil {
			// If we're at the end, just wait a bit and try again.
			if err == io.EOF {
				time.Sleep(250 * time.Millisecond)
				reader.Reset(file)
				return tailFileCmd(reader, file)() // Recurse to continue tailing
			}
			return errorMsg(err)
		}

		return newLineMsg(line)
	}
}

func tickCmd() tea.Msg {
	return tea.Tick(tickRate, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
