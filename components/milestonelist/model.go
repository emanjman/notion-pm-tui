package milestonelist

import (
	"notion-project-tui/notion"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type MilestoneListModel struct {
	list    list.Model
	loading bool

	grouped map[string][]MilestoneListItem // header to items
	hidden  map[string]bool                // hidden group

	// keys KeyMap
}

func NewMilestoneListModel() MilestoneListModel {
	l := list.New([]list.Item{}, NewMilestoneListDelegate(), 0, 0)

	// custom configs
	l.Title = "Milestones"
	l.SetShowHelp(false)

	m := MilestoneListModel{
		list:    l,
		loading: true,

		grouped: groupByStatus(mockMilestoneItems()),
		hidden:  map[string]bool{},
	}

	m.list.SetItems(m.buildGroupedList())
	return m
}

// just forward the list.Update(msg)
// and forward its returned response
func (m MilestoneListModel) Update(msg tea.Msg) (MilestoneListModel, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch {
		// case key.Matches(msg, )
		}

	case notion.MilestoneMsg:
		if msg.Err != nil {
			return m, nil
		}

		// create the list items
		tempItems := make([]MilestoneListItem, len(msg.Data))
		for i, page := range msg.Data {
			tempItems[i] = NewMilestoneListItem(page)
		}

		m.grouped = groupByStatus(tempItems)
		items := m.buildGroupedList()

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

type MilestoneListItemHeader struct {
	Label  string
	Hidden bool
	Count  int
}

// exclude header in filter-search (required item function)
func (g MilestoneListItemHeader) FilterValue() string { return "" }

// helper func to group items into a map (keyed by their status)
func groupByStatus(items []MilestoneListItem) map[string][]MilestoneListItem {

	groups := map[string][]MilestoneListItem{}
	for _, item := range items {
		groups[item.Status] = append(groups[item.Status], item)
	}
	return groups
}

var statusOrder = []string{"🚧 under development", "😴 idle", "🎉 complete"}

// conforms the group-map into a []list.Item + add headers
func (m MilestoneListModel) buildGroupedList() []list.Item {
	var items []list.Item

	// build list in this group order
	for _, status := range statusOrder {
		group, ok := m.grouped[status]
		if !ok {
			continue
		}

		// add group header
		items = append(items, MilestoneListItemHeader{
			Label:  status,
			Hidden: m.hidden[status],
			Count:  len(group),
		})

		// add the group's items (if not hidden)
		if !m.hidden[status] {
			for _, milestone := range group {
				items = append(items, milestone)
			}
		}
	}

	return items
}

// -------------------------------------------

func mockMilestoneItems() []MilestoneListItem {
	return []MilestoneListItem{
		MilestoneListItem{
			ID:       "1",
			Name:     "Setup Project Structure",
			Status:   "🎉 complete",
			Progress: 1.0,
			Tags:     []string{"backend", "setup"},
		},
		MilestoneListItem{
			ID:       "2",
			Name:     "Implement Notion API",
			Status:   "🚧 under development",
			Progress: 0.75,
			Tags:     []string{"backend", "api"},
		},
		MilestoneListItem{
			ID:       "3",
			Name:     "Build TUI Dashboard",
			Status:   "🚧 under development",
			Progress: 0.4,
			Tags:     []string{"frontend", "tui"},
		},
		MilestoneListItem{
			ID:       "4",
			Name:     "Authentication System",
			Status:   "😴 idle",
			Progress: 0.0,
			Tags:     []string{"backend", "auth"},
		},
		MilestoneListItem{
			ID:       "5",
			Name:     "Data Persistence Layer",
			Status:   "😴 idle",
			Progress: 0.0,
			Tags:     []string{"backend", "database"},
		},
		MilestoneListItem{
			ID:       "6",
			Name:     "Testing & QA",
			Status:   "😴 idle",
			Progress: 0.0,
			Tags:     []string{"testing"},
		},
		MilestoneListItem{
			ID:       "7",
			Name:     "Documentation",
			Status:   "🚧 under development",
			Progress: 0.2,
			Tags:     []string{"docs"},
		},
	}
}
