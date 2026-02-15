package tui

import (
	"sort"

	"github.com/emanuelpecson/notion-project-tui/internal/models"
)

// ProjectGroup represents a group of projects with the same status
type ProjectGroup struct {
	Status   string
	Projects []models.Project
}

// statusPriority defines the display order of statuses
var statusPriority = map[string]int{
	"developing":  1,
	"maintaining": 2,
	"planning":    3,
	"complete":    4,
	"idea":        5,
	// archived is intentionally omitted - it will be filtered out
}

// GroupProjectsByStatus organizes projects into groups based on their status
// Archived projects are excluded from the results
func GroupProjectsByStatus(projects []models.Project) []ProjectGroup {
	// Map to collect projects by status
	groupMap := make(map[string][]models.Project)

	for _, project := range projects {
		status := project.GetStatus()

		// Skip archived projects
		if status == "archived" {
			continue
		}

		if status == "" {
			status = "no status"
		}

		groupMap[status] = append(groupMap[status], project)
	}

	// Convert map to slice
	var groups []ProjectGroup
	for status, projs := range groupMap {
		groups = append(groups, ProjectGroup{
			Status:   status,
			Projects: projs,
		})
	}

	// Sort groups by priority order
	sort.Slice(groups, func(i, j int) bool {
		priorityI, okI := statusPriority[groups[i].Status]
		priorityJ, okJ := statusPriority[groups[j].Status]

		// If both have priorities, sort by priority (lower number = higher priority)
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

// GetTotalProjects returns the total number of projects across all groups
func GetTotalProjects(groups []ProjectGroup) int {
	total := 0
	for _, group := range groups {
		total += len(group.Projects)
	}
	return total
}
