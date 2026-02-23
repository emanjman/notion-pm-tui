package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"

	"notion-project-tui/model"
)

func init() {
	// load .env
	godotenv.Load()
}

func main() {
	p := tea.NewProgram(model.InitProjectModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Program error: %v", err)
		os.Exit(1)
	}
}
