package notion

type RelationProperty struct {
	ID       string         `json:"id"`
	HasMore  bool           `json:"has_more"`
	Relation []RelationItem `json:"relation"`
}

// each relation comes in as { id: xxx-xxx-xxxx }
type RelationItem struct {
	ID string `json:"id"`
}
