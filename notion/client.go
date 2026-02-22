package notion

import (
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const baseUrl = "https://api.notion.com/v1"
const version = "2025-09-03"

type Client struct {
	token  string
	projId string
}

// constructor
func NewClient() *Client {
	// address of newly created client
	return &Client{
		token:  os.Getenv("NOTION_API_TOKEN"),
		projId: os.Getenv("NOTION_HOOP_ARCHIVES_ID"),
	}
}

// cmd func returns a tea.Msg
func (c *Client) FetchProjectById() tea.Cmd {
	return func() tea.Msg {
		url := baseUrl + "/pages/" + c.projId

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return ProjectMsg{Err: err}
		}

		req.Header.Add("Notion-Version", version)
		req.Header.Add("Authorization", "Bearer "+c.token)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return ProjectMsg{Err: err}
		}

		// todo: to complete later
		return ProjectMsg{Data: nil}
	}
}
