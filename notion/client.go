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

// cmd func returns a tea.Msg
func (c *Client) FetchProjectById() tea.Cmd {
	return func() tea.Msg {
		start := time.Now()

		url := baseUrl + "/pages/" + c.projId

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return ProjectMsg{Err: err, Duration: time.Since(start)}
		}

		req.Header.Add("Notion-Version", version)
		req.Header.Add("Authorization", "Bearer "+c.token)

		res, err := c.http.Do(req)
		if err != nil {
			return ProjectMsg{Err: err, Duration: time.Since(start)}
		}
		defer res.Body.Close() // close by end-of-life

		if res.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(res.Body)
			return ProjectMsg{
				Err:      fmt.Errorf("Notion API error: %d: %s", res.StatusCode, body),
				Duration: time.Since(start),
			}
		}

		// parse as json
		var proj ProjectPage
		if err := json.NewDecoder(res.Body).Decode(&proj); err != nil {
			return ProjectMsg{Err: err, Duration: time.Since(start)}
		}

		return ProjectMsg{Data: proj, Duration: time.Since(start)}
	}
}
