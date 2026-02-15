package models

import "github.com/jomei/notionapi"

// TaskGroup represents a task group/branch within a project
type TaskGroup struct {
	Page notionapi.Page
}

// NewTaskGroup wraps a notionapi.Page in a TaskGroup
func NewTaskGroup(page notionapi.Page) TaskGroup {
	return TaskGroup{Page: page}
}

// GetTitle returns the task group name/title
func (tg *TaskGroup) GetTitle() string {
	// Task groups might use different property names
	// Try common variations: "Name", "name", "title", "Title"

	// Try "Name" first (most common for sub-pages)
	if prop, ok := tg.Page.Properties["Name"].(*notionapi.TitleProperty); ok && len(prop.Title) > 0 {
		return prop.Title[0].PlainText
	}

	// Try "name"
	if prop, ok := tg.Page.Properties["name"].(*notionapi.TitleProperty); ok && len(prop.Title) > 0 {
		return prop.Title[0].PlainText
	}

	// Try "title"
	if prop, ok := tg.Page.Properties["title"].(*notionapi.TitleProperty); ok && len(prop.Title) > 0 {
		return prop.Title[0].PlainText
	}

	// Fallback: use the page's own title if it exists
	// (Some pages store title directly in the page object, not properties)
	// For now, return empty string if no property found
	return ""
}

// GetIcon returns the emoji icon (or empty string if none)
func (tg *TaskGroup) GetIcon() string {
	if tg.Page.Icon != nil && tg.Page.Icon.Emoji != nil {
		return string(*tg.Page.Icon.Emoji)
	}
	return ""
}

// GetID returns the task group's page ID
func (tg *TaskGroup) GetID() string {
	return string(tg.Page.ID)
}

// GetStatus returns the status from the $status formula property
// Returns empty string if not found or not a string formula
func (tg *TaskGroup) GetStatus() string {
	prop, ok := tg.Page.Properties["$status"].(*notionapi.FormulaProperty)
	if !ok {
		return ""
	}

	// Check if formula result is a string
	if prop.Formula.Type == notionapi.FormulaTypeString {
		return prop.Formula.String
	}

	return ""
}

// GetTaskIDs returns the IDs of all tasks in this task group
func (tg *TaskGroup) GetTaskIDs() []string {
	prop, ok := tg.Page.Properties["@tasks"].(*notionapi.RelationProperty)
	if !ok {
		return nil
	}

	ids := make([]string, 0, len(prop.Relation))
	for _, rel := range prop.Relation {
		ids = append(ids, string(rel.ID))
	}
	return ids
}

// GetTaskCount returns the number of tasks in this task group
func (tg *TaskGroup) GetTaskCount() int {
	prop, ok := tg.Page.Properties["@tasks"].(*notionapi.RelationProperty)
	if !ok {
		return 0
	}
	return len(prop.Relation)
}
