package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/emanuelpecson/notion-project-tui/internal/config"
	"github.com/emanuelpecson/notion-project-tui/internal/notion"
	"github.com/emanuelpecson/notion-project-tui/internal/tui"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiToken := os.Getenv("NOTION_API_TOKEN")
	databaseID := os.Getenv("NOTION_PROJECTS_DB_ID")

	if apiToken == "" || databaseID == "" {
		log.Fatal("Missing NOTION_API_TOKEN or NOTION_PROJECTS_DB_ID in .env file")
	}

	fmt.Printf("📊 Database ID: %s\n\n", databaseID)

	// Create Notion client
	notionClient := notion.NewClient(apiToken)

	projects, err := notionClient.FetchAllProjects(databaseID)
	if err != nil {
		log.Fatalf("Error querying database: %v", err)
	}

	fmt.Printf("📊 Total projects found: %d\n", len(projects))

	// Load config (optional - provides empty config if file doesn't exist)
	cfg := config.LoadConfig(".notion-tui-config")
	if len(cfg.PriorityProjectIDs) > 0 {
		fmt.Printf("⚙️  Priority projects configured: %d\n", len(cfg.PriorityProjectIDs))
	}

	fmt.Println("Starting TUI...\n")

	// Create and run the TUI with the root model
	tuiModel := tui.NewRootModel(projects, notionClient, cfg)
	p := tea.NewProgram(tuiModel)

	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running TUI: %v", err)
	}
}
