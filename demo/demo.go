package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/log"
)

type model struct {
	filename    string
	logCount    int
	lastLevel   string
	lastMessage string
	file        *os.File
	width       int
	height      int
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	// Random interval between 1-3 seconds
	interval := time.Millisecond * time.Duration(250+rand.Intn(1000))
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type logGeneratedMsg struct {
	level   log.Level
	message string
}

var (
	logMessages = map[string][]string{
		"INFO": {
			"Application started successfully",
			"Processing user request for %s",
			"Database connection established",
			"Configuration loaded from %s",
			"Background task scheduled",
			"Authentication successful for user: %s",
			"Health check passed",
			"Backup completed successfully",
			"Transaction committed successfully",
			"Cache warmed up with %d entries",
		},
		"WARN": {
			"Cache miss for key: %s",
			"Rate limit approaching for user: %s",
			"Disk space low: %d%% remaining",
			"SSL certificate expires in %d days",
			"Connection pool nearly exhausted",
			"Deprecated API endpoint accessed: %s",
			"Memory usage above threshold: %d MB",
			"Queue backlog growing: %d items",
		},
		"ERROR": {
			"Connection timeout occurred",
			"Invalid input received from %s",
			"Service unavailable: %s",
			"Database query failed for table: %s",
			"File not found: %s",
			"Authentication failed for user: %s",
			"Permission denied accessing %s",
			"Network error: connection refused",
		},
		"DEBUG": {
			"Processing request with ID: %s",
			"Query executed in %d ms",
			"Loading configuration from %s",
			"Initializing module: %s",
			"Session created for user: %s",
			"Cache hit for key: %s",
			"Worker thread %d started",
			"Parsing configuration file: %s",
		},
	}

	randomValues = []string{
		"user123", "session_abc", "config.yaml", "database.db",
		"cache_key", "api_endpoint", "service_name", "worker_01",
		"GET /api/users", "POST /login", "file.txt", "backup.zip",
		"/etc/app/config", "192.168.1.100", "auth_token_xyz",
	}
)

func generateLogCmd() tea.Cmd {
	return func() tea.Msg {
		// Simple arrays for random selection
		levels := []log.Level{log.InfoLevel, log.WarnLevel, log.ErrorLevel, log.DebugLevel}
		levelNames := []string{"INFO", "WARN", "ERROR", "DEBUG"}

		// Pick same random index for both
		index := rand.Intn(len(levels))
		level := levels[index]
		levelName := levelNames[index]

		// Get random message template for this level
		messages := logMessages[levelName]
		template := messages[rand.Intn(len(messages))]

		// Fill in template with random values
		var message string
		switch rand.Intn(3) {
		case 0:
			message = fmt.Sprintf(template, randomValues[rand.Intn(len(randomValues))])
		case 1:
			message = fmt.Sprintf(template, rand.Intn(1000))
		default:
			message = template
		}

		// Log using appropriate level
		switch level {
		case log.ErrorLevel:
			log.Error(message)
		case log.WarnLevel:
			log.Warn(message)
		case log.InfoLevel:
			log.Info(message)
		case log.DebugLevel:
			log.Debug(message)
		}

		return logGeneratedMsg{
			level:   level,
			message: message,
		}
	}
}

// Styling
var (
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Bold(true)

	levelStyles = map[string]lipgloss.Style{
		"INFO":  lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")),
		"WARN":  lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")),
		"ERROR": lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")),
		"DEBUG": lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF")),
	}

	statStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Italic(true)
)

func initialModel(file *os.File) model {
	return model{
		filename: file.Name(),
		file:     file,
	}
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.file.Close()
			return m, tea.Quit
		case "r":
			// Reset log file
			m.file.Close()
			m.file, _ = os.OpenFile(m.filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			m.logCount = 0
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		return m, tea.Batch(generateLogCmd(), tickCmd())

	case logGeneratedMsg:
		m.logCount++
		m.lastLevel = msg.level.String()
		m.lastMessage = msg.message
	}

	return m, nil
}

func (m model) View() string {
	header := headerStyle.Render(" Log Generator ")

	content := fmt.Sprintf("📁 Output File: %s\n", m.filename)
	content += fmt.Sprintf("📊 Logs Generated: %d\n\n", m.logCount)

	if m.lastLevel != "" {
		levelStyle := levelStyles[m.lastLevel]
		content += "Last Entry:\n"
		content += fmt.Sprintf("  Level: %s\n", levelStyle.Render(m.lastLevel))
		content += fmt.Sprintf("  Message: %s\n\n", m.lastMessage)
	}

	controls := statStyle.Render("Controls: q/ctrl+c (quit) | r (reset log file)")

	// Center everything
	if m.width > 0 {
		header = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, header)
		controls = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, controls)
	}

	return fmt.Sprintf("%s\n\n%s\n%s", header, content, controls)
}

func main() {
	filename := "debug.log"
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	log.SetLevel(log.DebugLevel)
	log.SetOutput(file)
	defer file.Close()

	fmt.Printf("Starting log generator - writing to: %s\n", filename)

	p := tea.NewProgram(initialModel(file))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
