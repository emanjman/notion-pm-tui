package models

import (
	"time"

	"github.com/jomei/notionapi"
)

// Task represents a task within a task group
type Task struct {
	Page notionapi.Page
}

// NewTask wraps a notionapi.Page in a Task
func NewTask(page notionapi.Page) Task {
	return Task{Page: page}
}

// GetText returns the task text/description
func (t *Task) GetText() string {
	// Try common property names for task text
	// Try "task" first
	if prop, ok := t.Page.Properties["task"].(*notionapi.TitleProperty); ok && len(prop.Title) > 0 {
		return prop.Title[0].PlainText
	}

	// Try "text"
	if prop, ok := t.Page.Properties["text"].(*notionapi.TitleProperty); ok && len(prop.Title) > 0 {
		return prop.Title[0].PlainText
	}

	// Try "Text"
	if prop, ok := t.Page.Properties["Text"].(*notionapi.TitleProperty); ok && len(prop.Title) > 0 {
		return prop.Title[0].PlainText
	}

	// Try "Name"
	if prop, ok := t.Page.Properties["Name"].(*notionapi.TitleProperty); ok && len(prop.Title) > 0 {
		return prop.Title[0].PlainText
	}

	return ""
}

// GetStatus returns the task status
func (t *Task) GetStatus() string {
	prop, ok := t.Page.Properties["status"].(*notionapi.StatusProperty)
	if !ok || prop.Status.Name == "" {
		return ""
	}
	return prop.Status.Name
}

// GetType returns the task type
func (t *Task) GetType() string {
	prop, ok := t.Page.Properties["type"].(*notionapi.SelectProperty)
	if !ok || prop.Select.Name == "" {
		return ""
	}
	return prop.Select.Name
}

// GetPriority returns the task priority
// Assumes priority is a number property
func (t *Task) GetPriority() float64 {
	prop, ok := t.Page.Properties["priority"].(*notionapi.NumberProperty)
	if !ok {
		return 0
	}
	return prop.Number
}

// GetCreatedTime returns when the task was created
func (t *Task) GetCreatedTime() time.Time {
	return time.Time(t.Page.CreatedTime)
}

// GetID returns the task's page ID
func (t *Task) GetID() string {
	return string(t.Page.ID)
}
