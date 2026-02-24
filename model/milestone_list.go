package model

import (
	"fmt"
	"io"
	"notion-project-tui/notion"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type MilestoneListModel struct {
	list    list.Model
	loading bool
}

func NewMilestoneListModel() MilestoneListModel {
	// l := list.New([]list.Item{}, NewMilestoneListDelegate(), 0, 0)
	l := list.New(mockMilestoneItems(), NewMilestoneListDelegate(), 0, 0)

	// custom configs
	l.Title = "Milestones"
	l.SetShowHelp(false)

	return MilestoneListModel{list: l, loading: true}
}

// just forward the list.Update(msg)
// and forward its returned response
func (m MilestoneListModel) Update(msg tea.Msg) (MilestoneListModel, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)

	case notion.MilestoneMsg:
		if msg.Err != nil {
			return m, nil
		}

		// create the list items
		items := make([]list.Item, len(msg.Data))
		for i, page := range msg.Data {
			items[i] = NewMilestoneListItem(page)
		}

		m.list.SetItems(items)
		m.loading = false

	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// just forward the list.View()
func (m MilestoneListModel) View() string {
	// ! temp, styling
	// if m.loading {
	// 	return "Loading milestones..."
	// }
	return m.list.View()
}

// --------------------------------------------

// implementation for the `list.Item` interface
type MilestoneListItem struct {
	ID           string
	Name         string
	Status       string
	LastActivity time.Time
	Progress     float64
	Tags         []string
}

// func (m MilestoneListItem) Title() string       { return m.Name }
// func (m MilestoneListItem) Description() string { return m.Status }
func (m MilestoneListItem) FilterValue() string { return m.Name }

func NewMilestoneListItem(page notion.MilestonePage) MilestoneListItem {
	title := notion.ExtractPlainText(page.Properties.Title.Title)

	status := ""
	if page.Properties.Status.Formula.String != nil {
		status = *page.Properties.Status.Formula.String
	}

	progress := 0.0
	if page.Properties.Progress.Formula.Number != nil {
		progress = *page.Properties.Progress.Formula.Number
	}

	tags := make([]string, len(page.Properties.Tags.MultiSelect))
	for i, tag := range page.Properties.Tags.MultiSelect {
		tags[i] = tag.Name
	}

	return MilestoneListItem{
		ID:       page.ID,
		Name:     title,
		Status:   status,
		Progress: progress,
		Tags:     tags,
	}
}

// -------------------------------------------

// implementation for delegate
type MilestoneListDelegate struct {
	defaultStyle  lg.Style
	selectedStyle lg.Style
}

func (d MilestoneListDelegate) Height() int  { return 3 }
func (d MilestoneListDelegate) Spacing() int { return 0 }
func (d MilestoneListDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d MilestoneListDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	milestone := item.(MilestoneListItem) // conform as specific list item type
	selected := index == m.Index()

	row1 := padBetween(milestone.Name, milestone.Status, m.Width())
	row2 := fmt.Sprintf("%.0f%%", milestone.Progress*100)
	block := row1 + "\n" + row2

	style := d.defaultStyle
	if selected {
		style = d.selectedStyle
	}

	// write to `w`
	fmt.Fprint(w, style.Width(m.Width()).Render(block))
}

func padBetween(left, right string, windowWidth int) string {
	// use lg.Width to only consider visible cells
	padding := windowWidth - lg.Width(left) - lg.Width(right)
	if padding < 0 {
		padding = 0
	}

	return left + strings.Repeat(" ", padding) + right
}

func NewMilestoneListDelegate() MilestoneListDelegate {
	base := lg.NewStyle().
		Border(lg.NormalBorder(), false, false, true, false).
		BorderForeground(lg.Color("238"))

	return MilestoneListDelegate{
		defaultStyle: base,

		selectedStyle: base.
			Bold(true).
			Foreground(lg.Color("205")),
	}
}

// ---

func mockMilestoneItems() []list.Item {
	return []list.Item{
		MilestoneListItem{
			ID:       "1",
			Name:     "Setup Project Structure",
			Status:   "Completed",
			Progress: 1.0,
			Tags:     []string{"backend", "setup"},
		},
		MilestoneListItem{
			ID:       "2",
			Name:     "Implement Notion API",
			Status:   "In Progress",
			Progress: 0.75,
			Tags:     []string{"backend", "api"},
		},
		MilestoneListItem{
			ID:       "3",
			Name:     "Build TUI Dashboard",
			Status:   "In Progress",
			Progress: 0.4,
			Tags:     []string{"frontend", "tui"},
		},
		MilestoneListItem{
			ID:       "4",
			Name:     "Authentication System",
			Status:   "Not Started",
			Progress: 0.0,
			Tags:     []string{"backend", "auth"},
		},
		MilestoneListItem{
			ID:       "5",
			Name:     "Data Persistence Layer",
			Status:   "Not Started",
			Progress: 0.0,
			Tags:     []string{"backend", "database"},
		},
		MilestoneListItem{
			ID:       "6",
			Name:     "Testing & QA",
			Status:   "Not Started",
			Progress: 0.0,
			Tags:     []string{"testing"},
		},
		MilestoneListItem{
			ID:       "7",
			Name:     "Documentation",
			Status:   "In Progress",
			Progress: 0.2,
			Tags:     []string{"docs"},
		},
	}
}
