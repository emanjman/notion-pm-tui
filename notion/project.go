package notion

type ProjectMsg struct {
	Data ProjectPage
	Err  error
}

// -------

type ProjectPage struct {
	ID         string            `json:"id"`
	Properties ProjectProperties `json:"properties"`
}

type ProjectProperties struct {
	Milestones   RelationProperty `json:"@milestones"`
	OverviewPage RelationProperty `json:"$overview"`
	ProjectNotes RelationProperty `json:"@notes"`
	DebugNotes   RelationProperty `json:"@debug-notes"`
}
