package models

import "github.com/jomei/notionapi"

// Project is a thin wrapper around a Notion page.
type Project struct {
	Page notionapi.Page
}

// NewProject wraps a notionapi.Page in a Project.
func NewProject(page notionapi.Page) Project {
	return Project{Page: page}
}

// GetTitle returns the project title/name.
func (p *Project) GetTitle() string {
	prop, ok := p.Page.Properties["project"].(*notionapi.TitleProperty)
	if !ok || len(prop.Title) == 0 {
		return ""
	}
	return prop.Title[0].PlainText
}

// GetIcon returns the emoji icon (or empty string if none).
func (p *Project) GetIcon() string {
	if p.Page.Icon != nil && p.Page.Icon.Emoji != nil {
		return string(*p.Page.Icon.Emoji)
	}
	return ""
}

// GetStatus returns the status name.
func (p *Project) GetStatus() string {
	prop, ok := p.Page.Properties["status"].(*notionapi.StatusProperty)
	if !ok || prop.Status.Name == "" {
		return ""
	}
	return prop.Status.Name
}

// GetGroup returns the group name (e.g., "personal", "work").
func (p *Project) GetGroup() string {
	prop, ok := p.Page.Properties["group"].(*notionapi.SelectProperty)
	if !ok || prop.Select.Name == "" {
		return ""
	}
	return prop.Select.Name
}

// GetTypes returns all project types as a slice of strings.
func (p *Project) GetTypes() []string {
	prop, ok := p.Page.Properties["type"].(*notionapi.MultiSelectProperty)
	if !ok {
		return nil
	}
	types := make([]string, 0, len(prop.MultiSelect))
	for _, t := range prop.MultiSelect {
		types = append(types, t.Name)
	}
	return types
}

// GetBranchCount returns the number of task groups/branches.
func (p *Project) GetBranchCount() int {
	prop, ok := p.Page.Properties["[branches]"].(*notionapi.RelationProperty)
	if !ok {
		return 0
	}
	return len(prop.Relation)
}

// GetBranchIDs returns the IDs of all task groups/branches
func (p *Project) GetBranchIDs() []string {
	prop, ok := p.Page.Properties["[branches]"].(*notionapi.RelationProperty)
	if !ok {
		return nil
	}

	ids := make([]string, 0, len(prop.Relation))
	for _, rel := range prop.Relation {
		ids = append(ids, string(rel.ID))
	}
	return ids
}

// GetID returns the project's page ID
func (p *Project) GetID() string {
	return string(p.Page.ID)
}
