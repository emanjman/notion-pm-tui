package notion

type TaskTypeOptionsMsg struct {
	Options []SelectItem
	Err     error
}

type TaskDatasourceResponse struct {
	Properties struct {
		Type struct {
			Select struct {
				Options []SelectItem `json:"options"`
			} `json:"select"`
		} `json:"type"`
	} `json:"properties"`
}

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
