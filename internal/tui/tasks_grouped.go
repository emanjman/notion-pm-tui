package tui

import (
	"sort"

	"github.com/emanuelpecson/notion-project-tui/internal/models"
)

// TaskStatusGroup represents a group of tasks with the same status
type TaskStatusGroup struct {
	Status string
	Tasks  []models.Task
}

// statusPriorityTasks defines the display order of task statuses
// Order: dev, idle, complete, archive
var statusPriorityTasks = map[string]int{
	"dev":      1,
	"idle":     2,
	"complete": 3,
	"archive":  4,
}

// GroupTasksByStatus organizes tasks into groups based on their status
// Tasks are already sorted by priority (desc) and created time (asc) from the fetch
func GroupTasksByStatus(tasks []models.Task) []TaskStatusGroup {
	// Map to collect tasks by status
	groupMap := make(map[string][]models.Task)

	for _, task := range tasks {
		status := task.GetStatus()
		if status == "" {
			status = "no status"
		}
		groupMap[status] = append(groupMap[status], task)
	}

	// Convert map to slice
	var groups []TaskStatusGroup
	for status, taskList := range groupMap {
		groups = append(groups, TaskStatusGroup{
			Status: status,
			Tasks:  taskList,
		})
	}

	// Sort groups by priority order
	sort.Slice(groups, func(i, j int) bool {
		priorityI, okI := statusPriorityTasks[groups[i].Status]
		priorityJ, okJ := statusPriorityTasks[groups[j].Status]

		// If both have priorities, sort by priority
		if okI && okJ {
			return priorityI < priorityJ
		}

		// Statuses with priority come first
		if okI {
			return true
		}
		if okJ {
			return false
		}

		// Otherwise sort alphabetically
		return groups[i].Status < groups[j].Status
	})

	return groups
}
