package notion

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/emanuelpecson/notion-project-tui/internal/models"
	"github.com/jomei/notionapi"
)

// Client wraps the Notion API client with helper methods
type Client struct {
	api *notionapi.Client
}

// NewClient creates a new Notion client wrapper
func NewClient(apiToken string) *Client {
	return &Client{
		api: notionapi.NewClient(notionapi.Token(apiToken)),
	}
}

// FetchAllProjects queries the Notion database and returns all projects, handling pagination
func (c *Client) FetchAllProjects(dbID string) ([]models.Project, error) {
	var projects []models.Project
	var cursor notionapi.Cursor

	for {
		req := &notionapi.DatabaseQueryRequest{
			PageSize: 100,
		}
		if cursor != "" {
			req.StartCursor = cursor
		}

		resp, err := c.api.Database.Query(context.Background(), notionapi.DatabaseID(dbID), req)
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

// FetchTaskGroups fetches task group pages given their IDs
// Returns them sorted by latest activity (most recent first)
func (c *Client) FetchTaskGroups(ids []string) ([]models.TaskGroup, error) {
	taskGroups := make([]models.TaskGroup, 0, len(ids))

	for _, id := range ids {
		page, err := c.api.Page.Get(context.Background(), notionapi.PageID(id))
		if err != nil {
			// Log error but continue with other pages
			fmt.Printf("Warning: failed to fetch task group %s: %v\n", id, err)
			continue
		}

		taskGroups = append(taskGroups, models.NewTaskGroup(*page))
	}

	// Sort by $latest-activity-at formula property (most recent first)
	sort.Slice(taskGroups, func(i, j int) bool {
		timeI := getLatestActivityTime(taskGroups[i])
		timeJ := getLatestActivityTime(taskGroups[j])

		// If both have valid times, compare them
		if timeI != nil && timeJ != nil {
			return (*time.Time)(timeI).After(*(*time.Time)(timeJ))
		}

		// Pages with activity time come before those without
		if timeI != nil {
			return true
		}
		if timeJ != nil {
			return false
		}

		// If neither has activity time, maintain order
		return false
	})

	return taskGroups, nil
}

// FetchTasks fetches task pages given their IDs
// Returns them sorted by priority (desc) then created time (asc)
func (c *Client) FetchTasks(ids []string) ([]models.Task, error) {
	tasks := make([]models.Task, 0, len(ids))

	for _, id := range ids {
		page, err := c.api.Page.Get(context.Background(), notionapi.PageID(id))
		if err != nil {
			// Log error but continue with other pages
			fmt.Printf("Warning: failed to fetch task %s: %v\n", id, err)
			continue
		}

		tasks = append(tasks, models.NewTask(*page))
	}

	// Sort by priority (descending - higher first), then by created time (ascending - oldest first)
	sort.Slice(tasks, func(i, j int) bool {
		priorityI := tasks[i].GetPriority()
		priorityJ := tasks[j].GetPriority()

		// If priorities are different, sort by priority (higher first)
		if priorityI != priorityJ {
			return priorityI > priorityJ
		}

		// If priorities are equal, sort by created time (older first)
		return tasks[i].GetCreatedTime().Before(tasks[j].GetCreatedTime())
	})

	return tasks, nil
}

// getLatestActivityTime extracts the date from the $latest-activity-at formula property
// Returns nil if the property doesn't exist or isn't a date
func getLatestActivityTime(tg models.TaskGroup) *notionapi.Date {
	prop, ok := tg.Page.Properties["$latest-activity-at"].(*notionapi.FormulaProperty)
	if !ok {
		return nil
	}

	// Formula can return different types - check if it's a date
	if prop.Formula.Type == notionapi.FormulaTypeDate && prop.Formula.Date != nil && prop.Formula.Date.Start != nil {
		return prop.Formula.Date.Start
	}

	return nil
}
