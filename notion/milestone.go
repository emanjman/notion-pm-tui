package notion

// -- enum --

type QueryMilestonePagesSource int

const (
	VersionChange QueryMilestonePagesSource = iota
	MoreMilestones
)

// -- types --

// notion page from @project-milestones
type MilestonePage struct {
	ID         string              `json:"id"`
	Properties MilestoneProperties `json:"properties"`
	Icon       *Icon               `json:"icon"`
}

// packaging grouped notion pages w/ state info
type MilestoneGroup struct {
	Milestones []MilestonePage
	NextCursor *string // bookmarks subseq notion-pages available
	Hide       bool    // is hidden
	Loading    bool    // is fetching notion pages
}

// map milestone-status to milestone-group
type MilestoneGroups map[MilestoneStatus]MilestoneGroup

// -- msg --

type AddMilestonePageMsg struct {
	TempID string         // optimistic temp-id to be reconciled
	Page   *MilestonePage // created notion-page (w/ real notion-id)
	Err    error          // failed notion-page creation
}

type QueryMilestonePagesMsg struct {
	Pages      []MilestonePage
	Status     MilestoneStatus // grouping-key
	NextCursor *string         // bookmarks subseq notion-pages available
	Err        error           // failed notion-page fetch
	Source     QueryMilestonePagesSource
	VersionID  string // version these milestones were fetched for
}
