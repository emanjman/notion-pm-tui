package models

import "time"

// Project represents a project from Notion
type Project struct {
	ID           string    `json:"id"`
	Icon         Icon      `json:"icon"`
	Properties   Properties `json:"properties"`
	CreatedTime  time.Time `json:"created_time"`
	LastEditedTime time.Time `json:"last_edited_time"`
	URL          string    `json:"url"`
}

// Icon represents the icon of a page (emoji or external image)
type Icon struct {
	Type  string `json:"type"`  // "emoji" or "external"
	Emoji string `json:"emoji"` // only present if type is "emoji"
}

// Properties contains all the database properties
type Properties struct {
	Project  TitleProperty  `json:"project"`  // The title/name of the project
	Status   StatusProperty `json:"status"`   // Project status
	Group    SelectProperty `json:"group"`    // Group (e.g., "personal", "work")
	Type     MultiSelectProperty `json:"type"` // Project types (e.g., "web", "mobile")
	Branches RelationProperty `json:"[branches]"` // Relation to task groups
}

// TitleProperty represents a title property (the project name)
type TitleProperty struct {
	Type  string      `json:"type"`
	Title []RichText  `json:"title"`
}

// StatusProperty represents a status property
type StatusProperty struct {
	Type   string       `json:"type"`
	Status StatusValue  `json:"status"`
}

// StatusValue contains the actual status data
type StatusValue struct {
	Name  string `json:"name"`  // e.g., "active", "archived", "on hold"
	Color string `json:"color"`
}

// SelectProperty represents a select property (single choice)
type SelectProperty struct {
	Type   string      `json:"type"`
	Select SelectValue `json:"select"`
}

// SelectValue contains the actual select data
type SelectValue struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// MultiSelectProperty represents a multi-select property (multiple choices)
type MultiSelectProperty struct {
	Type        string        `json:"type"`
	MultiSelect []SelectValue `json:"multi_select"`
}

// RelationProperty represents a relation to other pages
type RelationProperty struct {
	Type     string     `json:"type"`
	Relation []Relation `json:"relation"`
	HasMore  bool       `json:"has_more"`
}

// Relation represents a single relation entry
type Relation struct {
	ID string `json:"id"`
}

// RichText represents rich text content
type RichText struct {
	Type      string   `json:"type"`
	PlainText string   `json:"plain_text"`
	Text      TextContent `json:"text,omitempty"`
}

// TextContent contains the actual text content
type TextContent struct {
	Content string `json:"content"`
}

// Helper methods to extract commonly used data

// GetTitle returns the project title/name
func (p *Project) GetTitle() string {
	if len(p.Properties.Project.Title) > 0 {
		return p.Properties.Project.Title[0].PlainText
	}
	return ""
}

// GetIcon returns the emoji icon (or empty string if none)
func (p *Project) GetIcon() string {
	if p.Icon.Type == "emoji" {
		return p.Icon.Emoji
	}
	return ""
}

// GetStatus returns the status name
func (p *Project) GetStatus() string {
	return p.Properties.Status.Status.Name
}

// GetGroup returns the group name (e.g., "personal", "work")
func (p *Project) GetGroup() string {
	return p.Properties.Group.Select.Name
}

// GetTypes returns all project types as a slice of strings
func (p *Project) GetTypes() []string {
	types := make([]string, 0, len(p.Properties.Type.MultiSelect))
	for _, t := range p.Properties.Type.MultiSelect {
		types = append(types, t.Name)
	}
	return types
}

// GetBranchCount returns the number of task groups/branches
func (p *Project) GetBranchCount() int {
	return len(p.Properties.Branches.Relation)
}
