package notion

// -- types --

type ProjectPage struct {
	ID         string            `json:"id"`
	Properties ProjectProperties `json:"properties"`
}

// -- msg --

type QueryProjectPagesMsg struct {
	Pages      []ProjectPage
	NextCursor *string // bookmark subseq notion-pages avail; unlikely
	Err        error   // failed fetch
}
