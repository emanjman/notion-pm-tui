package notion

type NoteProperties struct {
	Title        TitleProperty    `json:"name"`
	Project      RelationProperty `json:"@project"`
	CreatedLabel FormulaProperty  `json:"$created-at-label"` // type:string
}

var (
	notePropTitle           = "name"
	notePropProjectRelation = "@project"
	notePropCreatedLabel    = "$created-at-label"
)
