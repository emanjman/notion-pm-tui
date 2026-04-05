package notion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Client struct {
	http      *http.Client
	token     string
	version   string
	baseURL   string
	projId    string
	tasksDsId string
}

// constructor
func NewClient() *Client {
	// address of newly created client
	return &Client{
		http:      &http.Client{Timeout: 10 * time.Second},
		token:     os.Getenv("NOTION_API_TOKEN"),
		version:   os.Getenv("NOTION_VERSION"),
		baseURL:   os.Getenv("NOTION_API_URL"),
		projId:    os.Getenv("NOTION_HOOP_ARCHIVES_ID"),
		tasksDsId: os.Getenv("NOTION_TASKS_DS_ID"),
	}
}

// helper func to decode result as json
func (c *Client) do(req *http.Request, target interface{}) error {
	req.Header.Add("Notion-Version", c.version)
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

		url := c.baseURL + "/pages/" + c.projId

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

func (c *Client) FetchPageBlocks(pageID string) ([]Block, error) {
	blocks, err := c.fetchBlocksRecursive(pageID)
	if err != nil {
		return nil, err
	}
	return blocks, nil
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
		url := c.baseURL + "/blocks/" + blockID + "/children?page_size=100"
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
		endpt := c.baseURL + "/pages/" + pageID + "/properties/" + propID + "?page_size=100"
		if cursor != "" {
			endpt += "&start_cursor=" + cursor
		}
		req, err := http.NewRequest("GET", endpt, nil)
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
		url := c.baseURL + "/pages/" + id
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

func (c *Client) FetchPageMarkdown(pageID string) (string, error) {
	url := c.baseURL + "/pages/" + pageID + "/markdown"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	var res MDSuccessRes
	if err = c.do(req, &res); err != nil {
		return "", err
	}
	return res.Markdown, nil
}

func (c *Client) ReplaceContentByMarkdown(pageID string, md string) (string, error) {
	url := c.baseURL + "/pages/" + pageID + "/markdown"

	log.Printf("entered replace content func")

	reqBody := MDReplaceReq{
		Type:           "replace_content",
		ReplaceContent: ReplaceContent{NewStr: md}}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return md, fmt.Errorf("failed to marshal markdown: %w", err)
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewReader(body))
	if err != nil {
		return md, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")

	var res MDSuccessRes
	if err := c.do(req, &res); err != nil {
		return md, err
	}

	log.Printf("replacing content")

	return res.Markdown, nil
}

func (c *Client) UpdatePageProperties(pageID string, props any) error {
	url := c.baseURL + "/pages/" + pageID

	body, err := json.Marshal(struct {
		Properties any `json:"properties"`
	}{Properties: props})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	var res struct{}
	return c.do(req, &res)
}

func (c *Client) QueryTasks(milestoneID, status, cursor string, milestoneIdx int) tea.Cmd {
	return func() tea.Msg {
		url := c.baseURL + "/data_sources/" + c.tasksDsId + "/query"

		body := map[string]any{
			"filter": map[string]any{
				"and": []map[string]any{
					{
						"property": "status",
						"status":   map[string]any{"equals": status},
					},
					{
						"property": "@milestone",
						"relation": map[string]any{"contains": milestoneID},
					},
				},
			},
			"sorts": []map[string]any{
				{
					"property":  "priority",
					"direction": "descending",
				},
				{
					"property":  "created-at",
					"direction": "ascending",
				},
			},
			"page_size": 5,
		}
		if cursor != "" {
			body["start_cursor"] = cursor
		}

		b, err := json.Marshal(body)
		if err != nil {
			return TaskQueryMsg{Err: err, Status: status, MilestoneIdx: milestoneIdx}
		}

		req, err := http.NewRequest("POST", url, bytes.NewReader(b))
		if err != nil {
			return TaskQueryMsg{Err: err, Status: status, MilestoneIdx: milestoneIdx}
		}
		req.Header.Add("Content-Type", "application/json")

		var res struct {
			Results    []TaskPage `json:"results"`
			NextCursor *string    `json:"next_cursor"`
			HasMore    bool       `json:"has_more"`
		}
		if err := c.do(req, &res); err != nil {
			return TaskQueryMsg{Err: err, Status: status, MilestoneIdx: milestoneIdx}
		}

		var nextCursor *string
		if res.HasMore {
			nextCursor = res.NextCursor
		}

		return TaskQueryMsg{
			Pages:        res.Results,
			NextCursor:   nextCursor,
			Status:       status,
			MilestoneIdx: milestoneIdx,
		}
	}
}

// options for `type` property, e.g. "style" "feat" "refactor"
func (c *Client) FetchTaskTypeOptions() tea.Cmd {
	return func() tea.Msg {
		url := c.baseURL + "/data_sources/" + c.tasksDsId
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return TaskTypeOptionsMsg{Err: err}
		}
		var res TaskDatasourceResponse
		if err := c.do(req, &res); err != nil {
			return TaskTypeOptionsMsg{Err: err}
		}
		return TaskTypeOptionsMsg{Options: res.Properties.Type.Select.Options}
	}
}
