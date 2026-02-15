package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/emanuelpecson/notion-project-tui/internal/models"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the API token and data source ID from environment variables
	apiToken := os.Getenv("NOTION_API_TOKEN")
	dataSourceID := os.Getenv("NOTION_PROJECTS_DS_ID")

	if apiToken == "" || dataSourceID == "" {
		log.Fatal("Missing NOTION_API_TOKEN or NOTION_PROJECTS_DS_ID in .env file")
	}

	fmt.Println("🔗 Connecting to Notion API (v2025-09-03)...")
	fmt.Printf("📊 Data Source ID: %s\n\n", dataSourceID)

	// Query the data source to get all projects
	projects, err := queryDataSource(apiToken, dataSourceID)
	if err != nil {
		log.Fatalf("Error querying database: %v", err)
	}

	fmt.Println("✅ Successfully fetched projects!")

	// Parse the results into Project structs
	projectList, err := parseProjects(projects)
	if err != nil {
		log.Fatalf("Error parsing projects: %v", err)
	}

	fmt.Printf("\n📊 Total projects found: %d\n\n", len(projectList))

	// Display all projects in a clean format
	fmt.Println("Projects:")
	fmt.Println("─────────────────────────────────────────────────")
	for i, project := range projectList {
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

// parseProjects converts the raw Notion API response into a slice of Project structs
func parseProjects(response map[string]interface{}) ([]models.Project, error) {
	// Extract the results array
	resultsRaw, ok := response["results"]
	if !ok {
		return nil, fmt.Errorf("no 'results' field in response")
	}

	// Convert back to JSON then unmarshal into Project structs
	// (This is easier than manually type asserting the nested structure)
	jsonData, err := json.Marshal(resultsRaw)
	if err != nil {
		return nil, err
	}

	var projects []models.Project
	err = json.Unmarshal(jsonData, &projects)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

// queryDataSource queries a Notion data source and returns the results
func queryDataSource(apiToken, dataSourceID string) (map[string]interface{}, error) {
	// Notion API endpoint for querying a data source (API version 2025-09-03)
	url := fmt.Sprintf("https://api.notion.com/v1/data_sources/%s/query", dataSourceID)

	// Create an empty request body (we're fetching all projects)
	requestBody := map[string]interface{}{}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	// Set required headers
	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("Notion-Version", "2025-09-03")
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("notion API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse the JSON response
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
