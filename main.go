package main

import (
	"fmt"
	"notion-project-tui/app"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	p := tea.NewProgram(app.NewModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Program error: %v", err)
		os.Exit(1)
	}
}
