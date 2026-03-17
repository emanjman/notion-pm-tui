package notion

import "time"

type ProjectNoteMsg struct {
	Data     []ProjectNotePage
	Err      error
	Duration time.Duration
}

type ProjectNoteRelationIdsMsg struct {
	IDs      []string
	Err      error
	Duration time.Duration
}

type ProjectNoteSelectedMsg struct {
	ID string
}

type ProjectNotePreviewMsg struct {
	PageID string
	Blocks []Block
	Err    error
}

// ------

type ProjectNotePage struct {
	ID          string                `json:"id"`
	CreatedTime string                `json:"created_time"`
	Properties  ProjectNoteProperties `json:"properties"`
}

type ProjectNoteProperties struct {
	Title TitleProperty `json:"title"`
}
