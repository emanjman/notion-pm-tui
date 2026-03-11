package notion

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const baseUrl = "https://api.notion.com/v1"
const version = "2025-09-03"

type Client struct {
	http   *http.Client
	token  string
	projId string
}

// constructor
func NewClient() *Client {
	// address of newly created client
	return &Client{
		http:   &http.Client{Timeout: 10 * time.Second},
		token:  os.Getenv("NOTION_API_TOKEN"),
		projId: os.Getenv("NOTION_HOOP_ARCHIVES_ID"),
	}
}

// helper func to decode result as json
func (c *Client) do(req *http.Request, target interface{}) error {
	req.Header.Add("Notion-Version", version)
	req.Header.Add("Authorization", "Bearer "+c.token)

	res, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close() // close by end-of-life

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("Notion API error: %d: %s", res.StatusCode, body)
	}

	// parse as json
	return json.NewDecoder(res.Body).Decode(target)
}

// cmd func returns a tea.Msg
func (c *Client) FetchProject() tea.Cmd {
	return func() tea.Msg {
		start := time.Now()

		url := baseUrl + "/pages/" + c.projId

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return ProjectMsg{Err: err, Duration: time.Since(start)}
		}

		// parse as json
		var proj ProjectPage
		if err := c.do(req, &proj); err != nil {
			return ProjectMsg{Err: err, Duration: time.Since(start)}
		}
		return ProjectMsg{Data: proj, Duration: time.Since(start)}
	}
}

type TaskRelationIdsMsg struct {
	IDs      []string
	Err      error
	Duration time.Duration
}

type MilestoneRelationIdsMsg struct {
	IDs      []string
	Err      error
	Duration time.Duration
}

func (c *Client) FetchMilestoneRelationIds(pageID string, propID string) tea.Cmd {
	return func() tea.Msg {
		start := time.Now()
		ids := []string{}
		cursor := ""

		for {
			url := baseUrl + "/pages/" + pageID + "/properties/" + propID + "?page_size=100"
			if cursor != "" {
				url += "&start_cursor=" + cursor
			}

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				return MilestoneRelationIdsMsg{Err: err, Duration: time.Since(start)}
			}

			var res RelationListResponse
			if err := c.do(req, &res); err != nil {
				return MilestoneRelationIdsMsg{Err: err, Duration: time.Since(start)}
			}

			for _, result := range res.Results {
				ids = append(ids, result.Relation.ID)
			}

			if !res.HasMore || res.NextCursor == nil {
				break
			}

			cursor = *res.NextCursor
		}

		return MilestoneRelationIdsMsg{IDs: ids, Duration: time.Since(start)}
	}
}

func (c *Client) FetchTaskRelationIds(pageID string, propID string) tea.Cmd {
	return func() tea.Msg {
		start := time.Now()
		ids := []string{}
		cursor := ""

		for {
			url := baseUrl + "/pages/" + pageID + "/properties/" + propID + "?page_size=100"
			if cursor != "" {
				url += "&start_cursor=" + cursor
			}

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				return TaskRelationIdsMsg{Err: err, Duration: time.Since(start)}
			}

			var res RelationListResponse
			if err := c.do(req, &res); err != nil {
				return TaskRelationIdsMsg{Err: err, Duration: time.Since(start)}
			}

			for _, result := range res.Results {
				ids = append(ids, result.Relation.ID)
			}

			// exit if we've exhausted all relations
			if !res.HasMore || res.NextCursor == nil {
				break
			}

			cursor = *res.NextCursor
		}

		return TaskRelationIdsMsg{IDs: ids, Duration: time.Since(start)}
	}
}

func (c *Client) FetchMilestones(ids []string) tea.Cmd {
	return func() tea.Msg {
		start := time.Now()
		milestones := make([]MilestonePage, 0, len(ids))

		for _, id := range ids {
			url := baseUrl + "/pages/" + id
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				return MilestoneMsg{Err: err, Duration: time.Since(start)}
			}

			var milestone MilestonePage
			if err := c.do(req, &milestone); err != nil {
				return MilestoneMsg{Err: err, Duration: time.Since(start)}
			}

			milestones = append(milestones, milestone)
		}

		return MilestoneMsg{Data: milestones, Duration: time.Since(start)}
	}
}

func (c *Client) FetchTasks(ids []string) tea.Cmd {
	return func() tea.Msg {
		start := time.Now()
		tasks := make([]TaskPage, 0, len(ids))

		for _, id := range ids {
			url := baseUrl + "/pages/" + id
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				return TaskMsg{Err: err, Duration: time.Since(start)}
			}

			var task TaskPage
			if err := c.do(req, &task); err != nil {
				return TaskMsg{Err: err, Duration: time.Since(start)}
			}

			tasks = append(tasks, task)
		}

		return TaskMsg{Data: tasks, Duration: time.Since(start)}
	}
}

func (c *Client) FetchBlockChildren(pageID string) ([]Block, error) {
	url := baseUrl + "/blocks/" + pageID + "/children"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Blocks     []Block `json:"results"`
		HasMore    bool    `json:"has_more"`
		NextCursor string  `json:"next_cursor"`
	}

	if err := c.do(req, &result); err != nil {
		return nil, err
	}

	return result.Blocks, nil
}
