package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/emanuelpecson/notion-project-tui/internal/models"
)

// TaskGroupsListModel is the Bubble Tea model for the task groups list view
type TaskGroupsListModel struct {
	project         models.Project
	groups          []TaskGroupGroup
	items           []TaskGroupListItem // Flattened list of headers + task groups
	cursor          int
	loadTime        string // Human-readable load time
	collapsedGroups map[string]bool // Track which status groups are collapsed
}

// TaskGroupListItemType represents the type of list item (header or task group)
type TaskGroupListItemType int

const (
	TaskGroupItemTypeHeader TaskGroupListItemType = iota
	TaskGroupItemTypeTaskGroup
)

// TaskGroupListItem represents an item in the task groups list (either a group header or a task group)
type TaskGroupListItem struct {
	Type      TaskGroupListItemType
	Status    string               // For headers: the status name
	TaskGroup *models.TaskGroup    // For task groups: the actual task group
}

// NewTaskGroupsListModel creates a new task groups list model
func NewTaskGroupsListModel(project models.Project, taskGroups []models.TaskGroup, loadTime string) TaskGroupsListModel {
	groups := GroupTaskGroupsByStatus(taskGroups)
	collapsedGroups := make(map[string]bool)
	items := flattenTaskGroupsToItems(groups, collapsedGroups)

	return TaskGroupsListModel{
		project:         project,
		groups:          groups,
		items:           items,
		cursor:          0,
		loadTime:        loadTime,
		collapsedGroups: collapsedGroups,
	}
}

// flattenTaskGroupsToItems converts grouped task groups into a flat list of items
// Respects collapsed state - collapsed groups don't include their task groups
func flattenTaskGroupsToItems(groups []TaskGroupGroup, collapsedGroups map[string]bool) []TaskGroupListItem {
	var items []TaskGroupListItem

	for _, group := range groups {
		// Add group header
		items = append(items, TaskGroupListItem{
			Type:   TaskGroupItemTypeHeader,
			Status: group.Status,
		})

		// Only add task groups if this group is not collapsed
		if !collapsedGroups[group.Status] {
			for i := range group.TaskGroups {
				items = append(items, TaskGroupListItem{
					Type:      TaskGroupItemTypeTaskGroup,
					TaskGroup: &group.TaskGroups[i],
				})
			}
		}
	}

	return items
}

// Init initializes the model
func (m TaskGroupsListModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m TaskGroupsListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.Down):
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}

		case key.Matches(msg, keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}

		case key.Matches(msg, keys.Enter):
			// If cursor is on a header, toggle its collapsed state
			if m.cursor < len(m.items) {
				item := m.items[m.cursor]
				if item.Type == TaskGroupItemTypeHeader {
					// Toggle collapsed state
					m.collapsedGroups[item.Status] = !m.collapsedGroups[item.Status]
					// Rebuild items list
					m.items = flattenTaskGroupsToItems(m.groups, m.collapsedGroups)
					// Keep cursor in valid range
					if m.cursor >= len(m.items) {
						m.cursor = len(m.items) - 1
					}
				}
			}

		case msg.String() == "esc":
			// TODO: Return to projects list
			// For now, just quit
			return m, tea.Quit
		}
	}

	return m, nil
}

// View renders the model
func (m TaskGroupsListModel) View() string {
	var b strings.Builder

	// Header with project info
	projectIcon := m.project.GetIcon()
	if projectIcon == "" {
		projectIcon = "•"
	}

	header := lipgloss.NewStyle().
		Bold(true).
		Padding(1, 2).
		Render(fmt.Sprintf("%s %s", projectIcon, m.project.GetTitle()))
	b.WriteString(header)
	b.WriteString("\n\n")

	// Task groups title with load time
	totalCount := m.countTotalTaskGroups()
	titleText := fmt.Sprintf("Task Groups (%d)", totalCount)
	if m.loadTime != "" {
		titleText += fmt.Sprintf(" • Loaded in %s", m.loadTime)
	}

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		Padding(0, 2).
		Render(titleText)
	b.WriteString(title)
	b.WriteString("\n\n")

	// Render task groups by status groups
	if len(m.items) == 0 {
		noGroups := lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Padding(0, 4).
			Render("No task groups found")
		b.WriteString(noGroups)
		b.WriteString("\n")
	} else {
		for i, item := range m.items {
			if item.Type == TaskGroupItemTypeHeader {
				// Render status header with collapse indicator
				count := m.countTaskGroupsInStatus(item.Status)

				// Choose indicator based on collapsed state
				indicator := "▼" // Expanded
				if m.collapsedGroups[item.Status] {
					indicator = "▶" // Collapsed
				}

				// Highlight header if cursor is on it
				cursor := " "
				if i == m.cursor {
					cursor = ">"
				}

				headerText := fmt.Sprintf("%s %s %s (%d)", cursor, indicator, item.Status, count)

				style := lipgloss.NewStyle().
					Bold(true).
					Foreground(lipgloss.Color("12")).
					Padding(0, 2)

				if i == m.cursor {
					style = style.Foreground(lipgloss.Color("170"))
				}

				b.WriteString(style.Render(headerText))
				b.WriteString("\n")

			} else if item.Type == TaskGroupItemTypeTaskGroup && item.TaskGroup != nil {
				// Render task group
				tg := item.TaskGroup
				icon := tg.GetIcon()
				if icon == "" {
					icon = "•"
				}

				cursor := "  "
				if i == m.cursor {
					cursor = "> "
				}

				taskGroupText := fmt.Sprintf("%s %s %s", cursor, icon, tg.GetTitle())

				style := lipgloss.NewStyle().Padding(0, 2)
				if i == m.cursor {
					style = style.Foreground(lipgloss.Color("170"))
				}

				b.WriteString(style.Render(taskGroupText))
				b.WriteString("\n")
			}
		}
	}

	// Footer with help
	b.WriteString("\n")
	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("↑/↓: Navigate • Enter: Toggle Group • Esc: Back • q: Quit")
	b.WriteString(help)

	return b.String()
}

// countTotalTaskGroups returns the total number of task groups across all status groups
func (m TaskGroupsListModel) countTotalTaskGroups() int {
	total := 0
	for _, group := range m.groups {
		total += len(group.TaskGroups)
	}
	return total
}

// countTaskGroupsInStatus counts how many task groups are in a specific status group
func (m TaskGroupsListModel) countTaskGroupsInStatus(status string) int {
	for _, group := range m.groups {
		if group.Status == status {
			return len(group.TaskGroups)
		}
	}
	return 0
}
