package milestone

import "notion-project-tui/notion"

// carries updated task groups for the selected milestone to the task panel
type MilestoneTasksMsg struct {
	Groups notion.TaskGroups // task pages keyed by status, each w/ pagination + visibility state
}

// result of notion milestone title rename
type UpdateNotionTitleMsg struct{ Err error }
