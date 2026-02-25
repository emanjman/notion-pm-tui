package main

import (
	"fmt"
	"notion-project-tui/app"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

func init() {
	// load .env
	godotenv.Load()
}

func main() {
	p := tea.NewProgram(app.InitProjectModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Program error: %v", err)
		os.Exit(1)
	}
}
