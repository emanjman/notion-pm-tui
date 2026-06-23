package notion

type ProjectProperties struct {
	Title TitleProperty `json:"project"`
}

var (
	projectPropTitle = "project"
)
