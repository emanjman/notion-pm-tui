package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/emanuelpecson/notion-project-tui/internal/models"
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

	fmt.Printf("📊 Total projects found: %d\n\n", len(projects))

	fmt.Println("Projects:")
	fmt.Println("─────────────────────────────────────────────────")
	for i, project := range projects {
		icon := project.GetIcon()
		if icon == "" {
			icon = "•"
		}

		fmt.Printf("%d. %s %s\n", i+1, icon, project.GetTitle())
		fmt.Printf("   Status: %s | Group: %s | Branches: %d\n",
			project.GetStatus(),
			project.GetGroup(),
			project.GetBranchCount())

		types := project.GetTypes()
		if len(types) > 0 {
			fmt.Printf("   Types: %v\n", types)
		}
		fmt.Println()
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
