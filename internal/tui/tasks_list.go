package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/emanuelpecson/notion-project-tui/internal/models"
)

// TasksListModel is the Bubble Tea model for the tasks list view
type TasksListModel struct {
	project         models.Project
	taskGroup       models.TaskGroup
	groups          []TaskStatusGroup
	items           []TaskListItem
	cursor          int
	loadTime        string
	collapsedGroups map[string]bool
}

// TaskListItemType represents the type of list item (header or task)
type TaskListItemType int

const (
	TaskItemTypeHeader TaskListItemType = iota
	TaskItemTypeTask
)

// TaskListItem represents an item in the tasks list (either a group header or a task)
type TaskListItem struct {
	Type   TaskListItemType
	Status string        // For headers: the status name
	Task   *models.Task  // For tasks: the actual task
}

// NewTasksListModel creates a new tasks list model
func NewTasksListModel(project models.Project, taskGroup models.TaskGroup, tasks []models.Task, loadTime string) TasksListModel {
	groups := GroupTasksByStatus(tasks)
	collapsedGroups := make(map[string]bool)
	items := flattenTasksToItems(groups, collapsedGroups)

	return TasksListModel{
		project:         project,
		taskGroup:       taskGroup,
		groups:          groups,
		items:           items,
		cursor:          0,
		loadTime:        loadTime,
		collapsedGroups: collapsedGroups,
	}
}

// flattenTasksToItems converts grouped tasks into a flat list of items
func flattenTasksToItems(groups []TaskStatusGroup, collapsedGroups map[string]bool) []TaskListItem {
	var items []TaskListItem

	for _, group := range groups {
		// Add group header
		items = append(items, TaskListItem{
			Type:   TaskItemTypeHeader,
			Status: group.Status,
		})

		// Only add tasks if this group is not collapsed
		if !collapsedGroups[group.Status] {
			for i := range group.Tasks {
				items = append(items, TaskListItem{
					Type: TaskItemTypeTask,
					Task: &group.Tasks[i],
				})
			}
		}
	}

	return items
}

// Init initializes the model
func (m TasksListModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m TasksListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				if item.Type == TaskItemTypeHeader {
					// Toggle collapsed state
					m.collapsedGroups[item.Status] = !m.collapsedGroups[item.Status]
					// Rebuild items list
					m.items = flattenTasksToItems(m.groups, m.collapsedGroups)
					// Keep cursor in valid range
					if m.cursor >= len(m.items) {
						m.cursor = len(m.items) - 1
					}
				}
			}

		case msg.String() == "esc":
			// TODO: Return to task groups list
			// For now, just quit
			return m, tea.Quit
		}
	}

	return m, nil
}

// View renders the model
func (m TasksListModel) View() string {
	var b strings.Builder

	// Header with breadcrumb
	projectIcon := m.project.GetIcon()
	if projectIcon == "" {
		projectIcon = "•"
	}

	taskGroupIcon := m.taskGroup.GetIcon()
	if taskGroupIcon == "" {
		taskGroupIcon = "•"
	}

	header := lipgloss.NewStyle().
		Bold(true).
		Padding(1, 2).
		Render(fmt.Sprintf("%s %s > %s %s",
			projectIcon, m.project.GetTitle(),
			taskGroupIcon, m.taskGroup.GetTitle()))
	b.WriteString(header)
	b.WriteString("\n\n")

	// Tasks title with load time
	totalCount := m.countTotalTasks()
	titleText := fmt.Sprintf("Tasks (%d)", totalCount)
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

	// Render tasks by status groups
	if len(m.items) == 0 {
		noTasks := lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Padding(0, 4).
			Render("No tasks found")
		b.WriteString(noTasks)
		b.WriteString("\n")
	} else {
		for i, item := range m.items {
			if item.Type == TaskItemTypeHeader {
				// Render status header
				count := m.countTasksInStatus(item.Status)

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

			} else if item.Type == TaskItemTypeTask && item.Task != nil {
				// Render task
				task := item.Task

				cursor := "  "
				if i == m.cursor {
					cursor = "> "
				}

				// Build task text: cursor + type → title
				taskType := task.GetType()
				if taskType == "" {
					taskType = "•"
				}

				leftPart := fmt.Sprintf("%s[%s] → %s", cursor, taskType, task.GetText())

				// Priority on the right
				priority := task.GetPriority()
				rightPart := fmt.Sprintf("P%.0f", priority)

				// Combine with spacing
				// Use a fixed width to align priorities
				const lineWidth = 80
				leftWidth := lipgloss.Width(leftPart)
				rightWidth := lipgloss.Width(rightPart)
				spacingWidth := lineWidth - leftWidth - rightWidth - 4 // 4 for padding

				var taskLine string
				if spacingWidth > 0 {
					taskLine = fmt.Sprintf("%s%s%s", leftPart, strings.Repeat(" ", spacingWidth), rightPart)
				} else {
					// If line is too long, just append priority
					taskLine = fmt.Sprintf("%s  %s", leftPart, rightPart)
				}

				style := lipgloss.NewStyle().Padding(0, 2)
				if i == m.cursor {
					style = style.Foreground(lipgloss.Color("170"))
				}

				b.WriteString(style.Render(taskLine))
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

// countTotalTasks returns the total number of tasks across all status groups
func (m TasksListModel) countTotalTasks() int {
	total := 0
	for _, group := range m.groups {
		total += len(group.Tasks)
	}
	return total
}

// countTasksInStatus counts how many tasks are in a specific status group
func (m TasksListModel) countTasksInStatus(status string) int {
	for _, group := range m.groups {
		if group.Status == status {
			return len(group.Tasks)
		}
	}
	return 0
}
