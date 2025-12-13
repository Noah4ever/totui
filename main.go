package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/noah4ever/totui/internal/app"
)

func main() {
	model := app.NewModel()
	p := tea.NewProgram(model)

	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}
}
