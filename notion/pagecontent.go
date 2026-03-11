package notion

type PageContent struct {
	Results    []Block `json:"results"`
	HasMore    bool    `json:"has_more"`
	NextCursor *string `json:"next_cursor,omitempty"`
}
