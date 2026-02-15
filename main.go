package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/emanuelpecson/notion-project-tui/internal/models"
	"github.com/emanuelpecson/notion-project-tui/internal/tui"
	"github.com/jomei/notionapi"
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

	client := notionapi.NewClient(notionapi.Token(apiToken))

	projects, err := fetchAllProjects(client, notionapi.DatabaseID(databaseID))
	if err != nil {
		log.Fatalf("Error querying database: %v", err)
	}

	fmt.Printf("📊 Total projects found: %d\n", len(projects))
	fmt.Println("Starting TUI...\n")

	// Create and run the TUI
	tuiModel := tui.NewProjectsListModel(projects)
	p := tea.NewProgram(tuiModel)

	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running TUI: %v", err)
	}
}

// fetchAllProjects queries the Notion database and returns all projects, handling pagination.
func fetchAllProjects(client *notionapi.Client, dbID notionapi.DatabaseID) ([]models.Project, error) {
	var projects []models.Project
	var cursor notionapi.Cursor

	for {
		req := &notionapi.DatabaseQueryRequest{
			PageSize: 100,
		}
		if cursor != "" {
			req.StartCursor = cursor
		}

		resp, err := client.Database.Query(context.Background(), dbID, req)
		if err != nil {
			return nil, fmt.Errorf("database query failed: %w", err)
		}

		for _, page := range resp.Results {
			projects = append(projects, models.NewProject(page))
		}

		if !resp.HasMore || resp.NextCursor == "" {
			break
		}
		cursor = notionapi.Cursor(resp.NextCursor)
	}

	return projects, nil
}
