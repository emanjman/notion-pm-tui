package main

import (
	"fmt"
	"log"
	"notion-project-tui/app"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	// set up log file
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	p := tea.NewProgram(app.New(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Program error: %v", err)
		os.Exit(1)
	}
}
