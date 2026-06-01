package notion

import "net/http"

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
