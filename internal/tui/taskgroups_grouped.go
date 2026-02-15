package tui

import (
	"sort"

	"github.com/emanuelpecson/notion-project-tui/internal/models"
)

// TaskGroupGroup represents a group of task groups with the same status
type TaskGroupGroup struct {
	Status     string
	TaskGroups []models.TaskGroup
}

// statusPriorityTaskGroups defines the display order of task group statuses
var statusPriorityTaskGroups = map[string]int{
	"🚧 under development": 1,
	"🎉 complete":          2,
	"😴 idle":              3,
}

// GroupTaskGroupsByStatus organizes task groups into groups based on their status
func GroupTaskGroupsByStatus(taskGroups []models.TaskGroup) []TaskGroupGroup {
	// Map to collect task groups by status
	groupMap := make(map[string][]models.TaskGroup)

	for _, tg := range taskGroups {
		status := tg.GetStatus()
		if status == "" {
			status = "no status"
		}
		groupMap[status] = append(groupMap[status], tg)
	}

	// Convert map to slice
	var groups []TaskGroupGroup
	for status, tgs := range groupMap {
		groups = append(groups, TaskGroupGroup{
			Status:     status,
			TaskGroups: tgs,
		})
	}

	// Sort groups by priority order
	sort.Slice(groups, func(i, j int) bool {
		priorityI, okI := statusPriorityTaskGroups[groups[i].Status]
		priorityJ, okJ := statusPriorityTaskGroups[groups[j].Status]

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
