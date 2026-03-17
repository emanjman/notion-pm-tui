package notion

type TaskPagesMsg struct {
	Pages []TaskPage
	Err   error
}

type TaskIDsMsg struct {
	IDs []string
	Err error
}

type TaskPage struct {
	ID         string         `json:"id"`
	Properties TaskProperties `json:"properties"`
}

type TaskProperties struct {
	Title     TitleProperty    `json:"task"`
	Status    StatusProperty   `json:"status"`
	Priority  SelectProperty   `json:"priority"`
	Type      SelectProperty   `json:"type"`
	Milestone RelationProperty `json:"@milestone"`
}
