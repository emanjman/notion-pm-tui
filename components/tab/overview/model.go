package overview

import (
	"notion-project-tui/components/pagecontent"
	"notion-project-tui/notion"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	content pagecontent.PageContentModel
}

func New(n *notion.Client) Model {
	return Model{
		content: pagecontent.NewPageContentModel(n),
	}
}

func (m Model) Init() tea.Cmd {
	return m.content.Init()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.content, cmd = m.content.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.content.View()
}
