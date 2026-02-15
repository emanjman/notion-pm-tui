package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/emanuelpecson/notion-project-tui/internal/models"
)

// ListItemType represents the type of list item (header or project)
type ListItemType int

const (
	ItemTypeHeader ListItemType = iota
	ItemTypeProject
)

// ListItem represents an item in the projects list (either a group header or a project)
type ListItem struct {
	Type    ListItemType
	Group   string // For headers: the status name
	Project *models.Project // For projects: the actual project
}

// ProjectsListModel is the Bubble Tea model for the projects list view
type ProjectsListModel struct {
	groups          []ProjectGroup
	items           []ListItem      // Flattened list of headers + projects
	cursor          int             // Currently selected item index
	selectedItem    *ListItem       // The currently selected item
	collapsedGroups map[string]bool // Track which status groups are collapsed
}

// keyMap defines the keyboard shortcuts
type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Enter  key.Binding
	Quit   key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// NewProjectsListModel creates a new projects list model
func NewProjectsListModel(projects []models.Project) ProjectsListModel {
	groups := GroupProjectsByStatus(projects)
	collapsedGroups := make(map[string]bool)
	items := flattenGroupsToItems(groups, collapsedGroups)

	return ProjectsListModel{
		groups:          groups,
		items:           items,
		cursor:          0,
		collapsedGroups: collapsedGroups,
	}
}

// flattenGroupsToItems converts grouped projects into a flat list of items
// with headers and projects interleaved
// Respects collapsed state - collapsed groups don't include their projects
func flattenGroupsToItems(groups []ProjectGroup, collapsedGroups map[string]bool) []ListItem {
	var items []ListItem

	for _, group := range groups {
		// Add group header
		items = append(items, ListItem{
			Type:  ItemTypeHeader,
			Group: group.Status,
		})

		// Only add projects if this group is not collapsed
		if !collapsedGroups[group.Status] {
			for i := range group.Projects {
				items = append(items, ListItem{
					Type:    ItemTypeProject,
					Project: &group.Projects[i],
				})
			}
		}
	}

	return items
}

// Init initializes the model
func (m ProjectsListModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m ProjectsListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.cursor < len(m.items) {
				item := m.items[m.cursor]

				// If cursor is on a header, toggle its collapsed state
				if item.Type == ItemTypeHeader {
					m.collapsedGroups[item.Group] = !m.collapsedGroups[item.Group]
					// Rebuild items list
					m.items = flattenGroupsToItems(m.groups, m.collapsedGroups)
					// Keep cursor in valid range
					if m.cursor >= len(m.items) {
						m.cursor = len(m.items) - 1
					}
				} else {
					// If cursor is on a project, store it for selection
					m.selectedItem = &m.items[m.cursor]
				}
			}
		}
	}

	return m, nil
}

// View renders the model
func (m ProjectsListModel) View() string {
	var b strings.Builder

	// Header
	header := lipgloss.NewStyle().
		Bold(true).
		Padding(1, 2).
		Render("Notion Projects")
	b.WriteString(header)
	b.WriteString("\n\n")

	// Render each item
	for i, item := range m.items {
		if item.Type == ItemTypeHeader {
			// Render group header with collapse indicator
			count := m.countProjectsInGroup(item.Group)

			// Choose indicator based on collapsed state
			indicator := "▼" // Expanded
			if m.collapsedGroups[item.Group] {
				indicator = "▶" // Collapsed
			}

			// Highlight header if cursor is on it
			cursor := " "
			if i == m.cursor {
				cursor = ">"
			}

			headerText := fmt.Sprintf("%s %s %s (%d)", cursor, indicator, capitalizeStatus(item.Group), count)

			style := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("12")).
				Padding(0, 2)

			if i == m.cursor {
				style = style.Foreground(lipgloss.Color("170"))
			}

			b.WriteString(style.Render(headerText))
			b.WriteString("\n")

		} else if item.Type == ItemTypeProject && item.Project != nil {
			// Render project
			project := item.Project
			icon := project.GetIcon()
			if icon == "" {
				icon = "•"
			}

			cursor := "  " // No cursor by default
			if i == m.cursor {
				cursor = "> " // Cursor indicator
			}

			projectText := fmt.Sprintf("%s %s %s",
				cursor,
				icon,
				project.GetTitle())

			// Highlight selected item
			style := lipgloss.NewStyle().Padding(0, 2)
			if i == m.cursor {
				style = style.Foreground(lipgloss.Color("170"))
			}

			b.WriteString(style.Render(projectText))
			b.WriteString("\n")
		}
	}

	// Footer with help
	b.WriteString("\n")
	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("↑/↓: Navigate • Enter: Select/Toggle • q: Quit")
	b.WriteString(help)

	return b.String()
}

// countProjectsInGroup counts how many projects are in a specific status group
func (m ProjectsListModel) countProjectsInGroup(status string) int {
	for _, group := range m.groups {
		if group.Status == status {
			return len(group.Projects)
		}
	}
	return 0
}

// capitalizeStatus capitalizes the status for display
func capitalizeStatus(status string) string {
	if status == "" {
		return ""
	}
	return strings.ToUpper(status[:1]) + status[1:]
}
