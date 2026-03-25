package notion

type TaskPagesMsg struct {
	Pages        []TaskPage
	Err          error
	MilestoneIdx int
}

type TaskIDsMsg struct {
	IDs          []string
	Err          error
	MilestoneIdx int
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
