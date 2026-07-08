package notion

// -- types --

type ProjectPage struct {
	ID         string            `json:"id"`
	Properties ProjectProperties `json:"properties"`
	Icon       *Icon             `json:"icon"`
}

// -- msg --

type QueryProjectPagesMsg struct {
	Pages      []ProjectPage
	NextCursor *string // bookmark subseq notion-pages avail; unlikely
	Err        error   // failed fetch
}

// simple method of trickling down selected project to child models;
// in `notion` pkg to reduce dependency conflicts
type ProjectIDMsg struct {
	ID string
}
