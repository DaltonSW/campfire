package cmd

import (
	"os"

	"go.dalton.dog/campfire/models"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/log"
)

func Run() {
	if len(os.Args) < 2 {
		log.Info("Usage: campfire <path/to/logfile>")
		return
	}

	logfilePath := os.Args[1]

	model := models.NewModel(logfilePath)

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	defer model.CloseModel()

	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running program:\n%v", err)
	}
}
