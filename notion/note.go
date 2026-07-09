package notion

type QueryNotePagesMsg struct {
	Pages      []NotePage
	NextCursor *string
	Err        error
}

// * this can also act as the `selected, all we need to fetch is from the id`
type NotePage struct {
	ID          string         `json:"id"`
	Properties  NoteProperties `json:"properties"`
	Icon        *Icon          `json:"icon"`
	CreatedTime string         `json:"created_time"`
}

type NoteSelectedMsg struct {
	Note NotePage
}
