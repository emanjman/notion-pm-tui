package notion

type MarkdownSuccessResponse struct {
	Object          string   `json:"object"`
	ID              string   `json:"id"`
	Markdown        string   `json:"markdown"`
	Truncated       bool     `json:"truncated"`
	UnknownBlockIDs []string `json:"unknown_block_ids"`
}

type MarkdownFailResponse struct {
	Object  struct{} `json:"object"`
	Message string   `json:"message"`
	Code    string   `json:"code"`
	Status  struct{} `json:"status"`
}
