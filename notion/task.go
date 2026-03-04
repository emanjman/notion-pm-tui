package notion

import "time"

type TaskMsg struct {
	Data     []TaskPage
	Err      error
	Duration time.Duration
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

var PrioritySelectValues = []int{1, 2, 3, 4, 5}
var TypeSelectValues = []string{"chore", "feat", "fix", "refactor", "style"}
