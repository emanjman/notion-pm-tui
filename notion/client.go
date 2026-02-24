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

type RelationIdsMsg struct {
	IDs      []string
	Err      error
	Duration time.Duration
}

func (c *Client) FetchAllRelationIds(pageID string, prop RelationProperty) tea.Cmd {
	return func() tea.Msg {
		start := time.Now()

		// populate w/ the initial set of relations
		ids := make([]string, len(prop.Relation))
		for i, r := range prop.Relation {
			ids[i] = r.ID
		}

		if !prop.HasMore {
			return RelationIdsMsg{IDs: ids, Duration: time.Since(start)}
		}

		// add the rest of relations
		cursor := ""
		for {
			url := baseUrl + "/pages/" + pageID + "/properties/" + prop.ID + "?page_size=100"
			if cursor != "" {
				url += "&start_cursor=" + cursor
			}

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				return RelationIdsMsg{Err: err, Duration: time.Since(start)}
			}

			var res RelationListResponse
			if err := c.do(req, &res); err != nil {
				return RelationIdsMsg{Err: err, Duration: time.Since(start)}
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

		return RelationIdsMsg{IDs: ids, Duration: time.Since(start)}
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
