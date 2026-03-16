package overview

import (
	"notion-project-tui/components/pagecontent"
	"notion-project-tui/notion"

	tea "github.com/charmbracelet/bubbletea"
)

type OverviewModel struct {
	content pagecontent.PageContentModel
}

func NewOverviewModel(n *notion.Client) OverviewModel {
	return OverviewModel{
		content: pagecontent.NewPageContentModel(n),
	}
}

func (m OverviewModel) Init() tea.Cmd {
	return m.content.Init()
}

func (m OverviewModel) Update(msg tea.Msg) (OverviewModel, tea.Cmd) {
	var cmd tea.Cmd
	m.content, cmd = m.content.Update(msg)
	return m, cmd
}

func (m OverviewModel) View() string {
	return m.content.View()
}
