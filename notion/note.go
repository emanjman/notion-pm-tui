package notion

import "time"

type NoteMsg struct {
	Data     []NotePage
	Err      error
	Duration time.Duration
}

// * this can also act as the `selected, all we need to fetch is from the id`
type NotePage struct {
	ID         string         `json:"id"`
	Properties NoteProperties `json:"properties"`
}

type NoteProperties struct {
	Title        TitleProperty    `json:"name"`
	Project      RelationProperty `json:"@project"`
	CreatedLabel FormulaProperty  `json:"$created"` //type:string
}

type NoteSelectedMsg struct {
	Note NotePage
}
