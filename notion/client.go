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
	http          *http.Client
	token         string
	version       string
	baseURL       string
	projId        string
	tasksDsId     string
	milestoneDsId string
}

// constructor
func NewClient() *Client {
	// address of newly created client
	return &Client{
		http:          &http.Client{Timeout: 10 * time.Second},
		token:         os.Getenv("NOTION_API_TOKEN"),
		version:       os.Getenv("NOTION_VERSION"),
		baseURL:       os.Getenv("NOTION_API_URL"),
		projId:        os.Getenv("NOTION_HOOP_ARCHIVES_ID"),
		tasksDsId:     os.Getenv("NOTION_TASKS_DS_ID"),
		milestoneDsId: os.Getenv("NOTION_MILESTONES_DS_ID"),
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
	url := c.baseURL + "/pages" + ID
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
