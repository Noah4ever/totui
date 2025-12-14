package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/noah4ever/totui/internal/app"
	"github.com/noah4ever/totui/internal/storage"
)

func main() {
	dataPath := "totui_data.json"
	lists, err := storage.Load(dataPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load data: %v\n", err)
	}

	model := app.NewModelWithData(lists)
	model.StoragePath = dataPath
	model.SaveFunc = storage.Save

	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}
}
