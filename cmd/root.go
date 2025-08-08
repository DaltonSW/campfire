package cmd

import (
	"context"
	"fmt"
	"os"

	"go.dalton.dog/campfire/internal/models"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

const Version = "0.9.2"

var rootCmd = &cobra.Command{
	Use:   "campfire <./path/to/file>",
	Short: "A quick and stylish log viewer",
	Long:  "Get cozy with your logs with campfire, a fast and beautiful log viewer!",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		model := models.NewModel(args[0])

		p := tea.NewProgram(
			model,
			tea.WithAltScreen(),       // Use the full size of the terminal
			tea.WithMouseCellMotion(), // Enable tracking the mouse wheel
		)

		if _, err := p.Run(); err != nil {
			log.Fatalf("Error running program:\n%v", err)
		}
	},
}

func Execute() {
	if err := fang.Execute(context.Background(), rootCmd, fang.WithoutManpage(), fang.WithoutCompletions(), fang.WithVersion(Version)); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
