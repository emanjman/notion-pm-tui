package notion

type TitleProperty struct {
	ID    string     `json:"id"`
	Title []RichText `json:"title"`
}

// ---------------------------------

type RelationProperty struct {
	ID       string         `json:"id"`
	HasMore  bool           `json:"has_more"`
	Relation []RelationItem `json:"relation"`
}

// each relation comes in as { id: xxx-xxx-xxxx }
type RelationItem struct {
	ID string `json:"id"`
}

// for handling paginated list properties
type RelationListResponse struct {
	Results []struct {
		Relation RelationItem `json:"relation"`
	} `json:"results"`
	NextCursor *string `json:"next_cursor"`
	HasMore    bool    `json:"has_more"`
}

// ---------------------------------

type FormulaProperty struct {
	ID      string      `json:"id"`
	Formula FormulaItem `json:"formula"`
}

type FormulaItem struct {
	Type string `json:"type"`

	// possible (optional) types
	String  *string       `json:"string,omitempty"`
	Number  *float64      `json:"number,omitempty"`
	Boolean *bool         `json:"boolean,omitempty"`
	Date    *DateProperty `json:"date,omitempty"`
}

// ------------------------------

type DateProperty struct {
	ID   string   `json:"id"`
	Date DateItem `json:"date"`
}

type DateItem struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// ------------------------------

type MultiSelectProperty struct {
	ID          string       `json:"id"`
	MultiSelect []SelectItem `json:"multi_select"`
}

type SelectProperty struct {
	ID     string     `json:"id"`
	Select SelectItem `json:"select"`
}

type StatusProperty struct {
	ID     string     `json:"id"`
	Status SelectItem `json:"status"`
}

type SelectItem struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}
