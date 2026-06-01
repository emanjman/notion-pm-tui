package notion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

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
