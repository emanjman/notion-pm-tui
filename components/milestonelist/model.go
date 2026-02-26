package milestonelist

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"notion-project-tui/notion"
)

type MilestoneListModel struct {
	list    list.Model
	loading bool
	grouped map[string][]MilestoneListItem // header to items
	hidden  map[string]bool                // hidden group
	Keys    KeyMap
}

func NewMilestoneListModel() MilestoneListModel {
	l := list.New([]list.Item{}, NewMilestoneListDelegate(), 0, 0)
	l.Title = "Milestones"
	l.SetShowHelp(false)

	m := MilestoneListModel{
		list:    l,
		loading: true,
		grouped: groupByStatus(mockMilestoneItems()),
		hidden:  map[string]bool{},
		Keys:    DefaultKeyMap,
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
		case key.Matches(msg, m.Keys.Select):
			selected := m.list.SelectedItem()

			if header, ok := selected.(MilestoneListItemHeader); ok {
				m.hidden[header.Label] = !m.hidden[header.Label] // toggle
				m.list.SetItems(m.buildGroupedList())            // rebuild list
			}

			return m, nil
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
	m.list, cmd = m.list.Update(msg) // handles up/down nav
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
