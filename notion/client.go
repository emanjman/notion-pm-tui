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

const baseURL = "https://api.notion.com/v1"
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

		url := baseURL + "/pages/" + c.projId

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

func (c *Client) FetchPageContent(pageID string) tea.Cmd {
	return func() tea.Msg {
		start := time.Now()

		blocks, err := c.fetchBlocksRecursive(pageID)

		if err != nil {
			return PageContentMsg{Err: err, Duration: time.Since(start)}
		}

		return PageContentMsg{Data: blocks, Duration: time.Since(start)}
	}
}

// return every block of the page by dfs
func (c *Client) fetchBlocksRecursive(pageID string) ([]Block, error) {
	blocks, err := c.fetchAllChildrenBlocks(pageID)
	if err != nil {
		return nil, err
	}

	for i, block := range blocks {
		if block.HasChildren {
			children, err := c.fetchBlocksRecursive(block.ID)

			if err != nil {
				return nil, err
			}

			blocks[i].Children = children
		}
	}

	return blocks, nil
}

// fetches all children blocks of a given page/block
func (c *Client) fetchAllChildrenBlocks(blockID string) ([]Block, error) {
	var blocks []Block
	cursor := ""

	for {
		url := baseURL + "/blocks/" + blockID + "/children?page_size=100"
		if cursor != "" {
			url += "&start_cursor=" + cursor
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		var result struct {
			Results    []Block `json:"results"`
			HasMore    bool    `json:"has_more"`
			NextCursor string  `json:"next_cursor"`
		}

		if err := c.do(req, &result); err != nil {
			return nil, err
		}

		// merge additional blocks
		blocks = append(blocks, result.Results...)

		// terminate if no more sibling blocks to fetch
		if !result.HasMore {
			break
		}

		cursor = result.NextCursor
	}

	// return final set of blocks
	return blocks, nil
}

// ---

func (c *Client) FetchRelationIDs(pageID string, propID string) ([]string, error) {
	ids := []string{}
	cursor := ""

	for {
		url := baseURL + "/pages/" + pageID + "/properties/" + propID + "?page_size=100"
		if cursor != "" {
			url += "&start_cursor=" + cursor
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		var res RelationListResponse
		if err := c.do(req, &res); err != nil {
			return nil, err
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

	return ids, nil
}

func FetchPages[T any](c *Client, ids []string) ([]T, error) {
	relations := make([]T, 0, len(ids))
	for _, id := range ids {
		url := baseURL + "/pages/" + id
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		var relation T
		if err := c.do(req, &relation); err != nil {
			return nil, err
		}
		relations = append(relations, relation)
	}
	return relations, nil
}
