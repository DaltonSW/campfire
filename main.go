package main

import (
	"os"

	"go.dalton.dog/campfire/models"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/log"
)

// TODO: Consider implementing cobra, but not sure this is complex enough to warrant that

// Entry point of the program, starts up the BubbleTea program
func main() {
	if len(os.Args) < 2 {
		log.Info("Usage: campfire <path/to/logfile>")
		return
	}

	logfilePath := os.Args[1]

	model := models.NewModel(logfilePath)

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),       // Use the full size of the terminal
		tea.WithMouseCellMotion(), // Enable tracking the mouse wheel
	)

	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running program:\n%v", err)
	}
}
