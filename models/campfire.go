package models

import (
	"fmt"
	"io"
	"io/fs"
	"strings"

	"os"
	"time"

	"github.com/charmbracelet/bubbles/v2/help"
	"github.com/charmbracelet/bubbles/v2/textinput"
	"github.com/charmbracelet/bubbles/v2/viewport"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

const tickRate = time.Millisecond * 1000

// Messages

type tickMsg time.Time
type fileGoneMsg struct{}
type fileExistsMsg struct {
	info    fs.FileInfo
	content []byte
}
type fileErrorMsg error
type viewportUpdateMsg []string

// NewModel actually creates the main campfire model
func NewModel(filename string) *model {
	// Viewport is initialized in after window size message

	text := textinput.New()
	text.Placeholder = "<text filter>"
	text.Prompt = "Filter: "

	m := model{
		filename:   filename,
		navKeys:    GetNavKeymap(),
		filterKeys: GetFilterKeymap(),
		textInput:  text,
		filters: Filters{
			ShowInfo:  true,
			ShowWarn:  true,
			ShowError: true,
			ShowDebug: true,
			ShowFatal: true,
			ShowOther: false,
		},
	}

	return &m
}

// model is the BubbleTea model for campfire
type model struct {
	filename      string
	content       []LogMessage
	viewport      viewport.Model
	width, height int
	ready         bool

	filterKeys FilterKeymap
	navKeys    NavKeymap
	help       help.Model

	textInput  textinput.Model
	textActive bool

	filters Filters

	fileExists   bool
	prevFileInfo fs.FileInfo
}

// Init kicks off the ticking
func (m model) Init() tea.Cmd {
	return tea.Batch(tickCmd(), checkFile(m.filename))
}

// Update processes new messages for the model
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// TODO: Turn these into key.Matches() calls instead of string checks
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch m.textActive {
		case true:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "esc":
				m.textActive = false
				m.filters.FilterText = ""
				m.textInput.SetValue("")
				m.textInput.Blur()
			case "enter":
				m.textActive = false
				m.filters.FilterText = m.textInput.Value()
				m.textInput.Blur()
			default:
				m.textInput, cmd = m.textInput.Update(msg)
				cmds = append(cmds, cmd)
				m.filters.FilterText = m.textInput.Value()
			}

		case false:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "q":
				return m, tea.Quit
			case "esc":
				return m, tea.Quit

			// Level filter toggles
			case "1":
				m.filters.ShowInfo = !m.filters.ShowInfo
			case "2":
				m.filters.ShowWarn = !m.filters.ShowWarn
			case "3":
				m.filters.ShowError = !m.filters.ShowError
			case "4":
				m.filters.ShowDebug = !m.filters.ShowDebug
			case "5":
				m.filters.ShowFatal = !m.filters.ShowFatal
			case "6":
				m.filters.ShowOther = !m.filters.ShowOther

			// Keyword filtering
			case "f":
				m.textActive = true
				m.textInput.Focus()

			case "x":
				m.textInput.SetValue("")
				m.filters.FilterText = ""
			}

		}
		cmds = append(cmds, updateViewport(m.content, m.filters))

	// BUG: Something is funky here... Resizing does seemingly inconsistent stuff
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.Header())
		footerHeight := lipgloss.Height(m.Footer())
		verticalMarginHeight := headerHeight + footerHeight

		m.width = msg.Width - appStyle.GetHorizontalFrameSize()
		m.height = msg.Height - appStyle.GetVerticalFrameSize()

		m.help.Width = m.width
		m.textInput.SetWidth(int(m.width / 2))

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
			m.viewport.SetWidth(msg.Width - viewportStyle.GetHorizontalBorderSize())
			m.viewport.SetHeight(m.height - verticalMarginHeight - viewportStyle.GetVerticalBorderSize())
		}

	case fileExistsMsg:
		m.prevFileInfo = msg.info
		m.fileExists = true
		m.content = make([]LogMessage, 0)
		for i, message := range strings.Split(string(msg.content), "\n") {
			logMsg := NewLogMessage(i, message)
			m.content = append(m.content, logMsg)
		}

		cmds = append(cmds, updateViewport(m.content, m.filters))

	case fileGoneMsg:
		m.fileExists = false
		m.viewport.SetContent("")

	case fileErrorMsg:
		content := "❌ Error reading file: " + msg.Error()
		m.viewport.SetContent(content)

	case viewportUpdateMsg:
		m.viewport.SetContentLines(msg)

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

func updateViewport(content []LogMessage, filters Filters) tea.Cmd {
	return func() tea.Msg {
		var outContent []string

		for _, msg := range content {
			if filters.IncludeMessage(msg) {
				outContent = append(outContent, msg.String())
			}
		}

		return viewportUpdateMsg(outContent)
	}
}
