package pagecontent

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"notion-project-tui/notion"
)

type PageContentModel struct {
	viewport viewport.Model
	notion   *notion.Client
	loading  bool
}

func NewPageContentModel(n *notion.Client) PageContentModel {
	vp := viewport.New(0, 0)

	return PageContentModel{
		viewport: vp,
		notion:   n,
		loading:  true,
	}
}

func (m PageContentModel) Init() tea.Cmd {
	return notion.NewClient().FetchPageContent("1e3b7273944b8059a15cd994116f24a9")
}

func (m PageContentModel) View() string {
	if m.loading {
		return "Loading..."
	}
	return m.viewport.View()
}

func (m PageContentModel) Update(msg tea.Msg) (PageContentModel, tea.Cmd) {
	switch msg := msg.(type) {
	case notion.PageContentMsg:
		if msg.Err != nil {
			return m, nil
		}

		m.viewport.SetContent(renderBlocks(msg.Data, m.viewport.Width, 0))
		m.loading = false

	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height

	// forward other commands to list
	default:
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}

	return m, nil
}
