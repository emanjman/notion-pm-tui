package notion

type NotePagesMsg struct {
	Pages []NotePage
	Err   error
}

type NoteIDsMsg struct {
	IDs []string
	Err error
}

// * this can also act as the `selected, all we need to fetch is from the id`
type NotePage struct {
	ID         string         `json:"id"`
	Properties NoteProperties `json:"properties"`
}

type NoteProperties struct {
	Title        TitleProperty       `json:"name"`
	Project      RelationProperty    `json:"@project"`
	CreatedDate  CreatedTimeProperty `json:"created"`
	CreatedLabel FormulaProperty     `json:"$created"` //type:string
}

type NoteSelectedMsg struct {
	Note NotePage
}
