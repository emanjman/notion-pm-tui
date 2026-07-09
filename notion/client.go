package notion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Client struct {
	http    *http.Client
	token   string
	version string
	baseURL string

	projectsDatasourceID   string
	tasksDatasourceID      string
	milestonesDatasourceID string
	versionsDatasourceID   string
	notesDatasourceID      string
}

// constructor
func NewClient() *Client {
	// address of newly created client
	return &Client{
		http:    &http.Client{Timeout: 10 * time.Second},
		token:   os.Getenv("NOTION_API_TOKEN"),
		version: os.Getenv("NOTION_VERSION"),
		baseURL: os.Getenv("NOTION_API_URL"),

		projectsDatasourceID:   os.Getenv("NOTION_PROJECTS_DS_ID"),
		tasksDatasourceID:      os.Getenv("NOTION_TASKS_DS_ID"),
		milestonesDatasourceID: os.Getenv("NOTION_MILESTONES_DS_ID"),
		versionsDatasourceID:   os.Getenv("NOTION_VERSIONS_DS_ID"),
		notesDatasourceID:      os.Getenv("NOTION_PROJECT_NOTES_DS_ID"),
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

type PaginationResponse[T any] struct {
	Results    []T     `json:"results"`
	NextCursor *string `json:"next_cursor"`
	HasMore    bool    `json:"has_more"`
}

func queryDatasource[T any](c *Client, datasourceID string, body map[string]any, cursor string, filterProps []string) (*PaginationResponse[T], error) {
	// setup url + queries
	url, err := url.Parse(c.baseURL + "/data_sources/" + datasourceID + "/query")
	if err != nil {
		return nil, err
	}
	q := url.Query()
	for _, fp := range filterProps {
		q.Add("filter_properties[]", fp)
	}
	url.RawQuery = q.Encode()
	// build body + inject cursor (if loading more)
	if cursor != "" {
		body["start_cursor"] = cursor
	}
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	// send req + prep response
	req, err := http.NewRequest("POST", url.String(), bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	var res PaginationResponse[T]
	if err := c.do(req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// returns error on failed req
func (c *Client) TrashPage(ID string) error {
	url := c.baseURL + "/pages/" + ID
	body := map[string]bool{
		"in_trash": true,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	var res any

	if err := c.do(req, &res); err != nil {
		return err
	}

	return nil
}
