package notion

// -- types --

type VersionPage struct {
	ID         string            `json:"id"`
	Properties VersionProperties `json:"properties"`
}

// -- msg --

type QueryVersionPagesMsg struct {
	Pages      []VersionPage
	NextCursor *string // bookmark subseq notion-pages avail; unlikely
	Err        error   // failed fetch
}
