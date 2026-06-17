package notion

type VersionProperties struct {
	Title TitleProperty `json:"name"`
	// Description RichTextProperty    `json:"description"`
	Project    RelationProperty    `json:"@project"`
	CreatedAt  CreatedTimeProperty `json:"created-at"`
	Milestones RelationProperty    `json:"@milestones"`
	// Completion  RollupProperty      `json:"%completion"` // type:number, avg of milestone $progress
}

var (
	versionPropTitle = "name"
	// versionPropDescription        = "description"
	versionPropProjectRelation    = "@project"
	versionPropCreatedAt          = "created-at"
	versionPropMilestonesRelation = "@milestones"
	// versionPropCompletionRollup   = "%completion"
)
