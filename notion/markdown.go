package notion

type MDUpdateReq struct {
	Markdown string `json:"markdown"`
}

type MDSuccessRes struct {
	Object          string   `json:"object"`
	ID              string   `json:"id"`
	Markdown        string   `json:"markdown"`
	Truncated       bool     `json:"truncated"`
	UnknownBlockIDs []string `json:"unknown_block_ids"`
}

type MDFailRes struct {
	Object  string `json:"object"`
	Message string `json:"message"`
	Code    string `json:"code"`
	Status  int    `json:"status"`
}
