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
	f, err := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	p := tea.NewProgram(app.New(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Program error: %v", err)
		os.Exit(1)
	}
}
