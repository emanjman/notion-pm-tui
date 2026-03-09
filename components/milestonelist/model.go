package milestonelist

import (
	"notion-project-tui/notion"
	listutil "notion-project-tui/util/list"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type MilestoneListModel struct {
	list    list.Model
	loading bool
	groups  map[string][]MilestoneListItem // header to items
	hidden  map[string]bool                // hidden group

	ActiveKeyMap    help.KeyMap // for help focus view
	neutralKeyMap   NeutralKeyMap
	selectingKeyMap SelectingKeyMap
	writingKeyMap   WritingKeyMap

	Focus *FocusState
}

var statusOrder = []string{"🚧 under development", "😴 idle", "🎉 complete"}

func NewMilestoneListModel() MilestoneListModel {
	f := FocusState{}

	l := list.New([]list.Item{}, NewMilestoneListDelegate(true, &f), 0, 0)
	l.Title = "Milestones"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	l.DisableQuitKeybindings()

	m := MilestoneListModel{
		list:    l,
		loading: true,
		groups:  listutil.GroupByKey(mockMilestoneItems()),
		hidden:  map[string]bool{},

		ActiveKeyMap:    NeutralKeyMapper, // default map view
		neutralKeyMap:   NeutralKeyMapper,
		selectingKeyMap: SelectingKeyMapper,
		writingKeyMap:   WritingKeyMapper,

		Focus: &f,
	}
	m.list.SetItems(listutil.BuildGroupList(m.groups, m.hidden, statusOrder))

	return m
}

func (m MilestoneListModel) Update(msg tea.Msg) (MilestoneListModel, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		// handle title editing (writing mode)
		if m.Focus.Mode == WritingMode {
			switch {
			case key.Matches(msg, m.writingKeyMap.Save):
				// update item in list
				if milestone, ok := m.list.SelectedItem().(MilestoneListItem); ok {
					milestone.Name = m.Focus.tempTitle.Value()
					m.list.SetItem(m.list.Index(), milestone)
					m.updateMilestoneInGroups(milestone)
				}

				m.ActiveKeyMap = NeutralKeyMapper
				m.Focus.Mode = NeutralMode

				// todo: send command to update milestone title in notion
				return m, nil

			// forward all keys into the textinput model
			default:
				var cmd tea.Cmd
				m.Focus.tempTitle, cmd = m.Focus.tempTitle.Update(msg)
				return m, cmd
			}
		}

		if m.Focus.Mode == SelectingMode {
			switch {

			// on exit, save updates via notion api
			case key.Matches(msg, m.selectingKeyMap.Exit):
				m.Focus.Mode = NeutralMode
				m.ActiveKeyMap = NeutralKeyMapper

				// todo: send command to update milestone changes in notion
				return m, nil

			// switch between fields (vertical navigation)
			case key.Matches(msg, m.selectingKeyMap.Up):
				if m.Focus.field == MilestoneTitle {
					m.Focus.field = fieldCnt - 1
				} else {
					m.Focus.field = (m.Focus.field - 1) % fieldCnt
				}
				return m, nil
			case key.Matches(msg, m.selectingKeyMap.Down):
				m.Focus.field = (m.Focus.field + 1) % fieldCnt
				return m, nil

			// enter field: cycle select or enter writing mode
			case key.Matches(msg, m.selectingKeyMap.Select):
				selected := m.list.SelectedItem()
				if milestone, ok := selected.(MilestoneListItem); ok {
					switch m.Focus.field {
					case MilestoneTag:
						// cycle tag, stay in selecting mode
						milestone.Tag = cycleTagField(milestone.Tag, 1)
						m.list.SetItem(m.Focus.milestoneIdx, milestone)
						m.updateMilestoneInGroups(milestone)
					case MilestoneTitle:
						// enter writing mode for title
						m.Focus.Mode = WritingMode
						m.ActiveKeyMap = WritingKeyMapper

						if item, ok := m.list.SelectedItem().(MilestoneListItem); ok {
							m.Focus.tempTitle = initTempTitle(item)
						}
					}

					return m, nil
				}
			}

			// consume all keys, don't forward to list navigations
			return m, nil
		}

		if m.Focus.Mode == NeutralMode {
			switch {
			case key.Matches(msg, m.neutralKeyMap.Select):
				selected := m.list.SelectedItem()

				// if selected item is header, toggle + rebuild list
				if header, ok := selected.(listutil.ListItemGroupHeader); ok {
					m.hidden[header.Label] = !m.hidden[header.Label]
					m.list.SetItems(listutil.BuildGroupList(m.groups, m.hidden, statusOrder))
				} else if milestone, ok := selected.(MilestoneListItem); ok {
					// initialize the focus state
					m.Focus.milestoneID = milestone.ID
					m.Focus.milestoneIdx = m.list.Index()
					m.Focus.field = MilestoneTitle // default field

					m.ActiveKeyMap = SelectingKeyMapper
					m.Focus.Mode = SelectingMode
				}

				return m, nil
			}
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

		m.groups = listutil.GroupByKey(tempItems)
		items := listutil.BuildGroupList(m.groups, m.hidden, statusOrder)

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

	containerStyle := lg.NewStyle().PaddingRight(1)
	return containerStyle.Render(m.list.View())
}

func (m MilestoneListModel) SelectedMilestone() notion.SelectedMilestone {
	item := m.list.SelectedItem()

	switch item := item.(type) {
	// get first milestone id of this group
	case listutil.ListItemGroupHeader:
		milestone := m.groups[item.Label][0]
		return notion.SelectedMilestone{
			ID:          milestone.ID,
			TasksPropID: milestone.TasksPropID,
		}
	// otherwise, on milestone, return its id
	case MilestoneListItem:
		return notion.SelectedMilestone{
			ID:          item.ID,
			TasksPropID: item.TasksPropID,
		}
	}

	return notion.SelectedMilestone{}
}

func (m *MilestoneListModel) SetItemDelegate(d list.ItemDelegate) {
	m.list.SetDelegate(d)
}

func (m MilestoneListModel) updateMilestoneInGroups(updated MilestoneListItem) MilestoneListModel {
	group := m.groups[updated.Status]

	// overwrite task in m.groups
	for i, t := range group {
		if t.ID == updated.ID {
			m.groups[updated.Status][i] = updated
			break
		}
	}

	// then rebuild item list
	m.list.SetItems(listutil.BuildGroupList(m.groups, m.hidden, statusOrder))
	return m
}
