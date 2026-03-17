package pagecontent

import (
	"notion-project-tui/notion"
	"notion-project-tui/styles"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type PageContentModel struct {
	viewport viewport.Model
	notion   *notion.Client
	loading  bool
	error    error
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
	return notion.NewClient().FetchPageContent("30eb7273944b80ad80c4f91a4f5cfd8d")
}

func (m PageContentModel) View() string {
	if m.error != nil {
		style := lg.NewStyle().Foreground(styles.RedForeground)
		return style.Render(m.error.Error())
	}

	if m.loading {
		loadingStyle := lg.NewStyle().Height(m.viewport.Height)
		return loadingStyle.Render("Loading...")
	}

	style := lg.NewStyle().Padding(0, 1)
	return style.Render(m.viewport.View())
}

func (m PageContentModel) WithContent(blocks []notion.Block) PageContentModel {
	m.viewport.SetContent(renderBlocks(blocks, m.viewport.Width, 0))
	m.loading = false
	m.error = nil
	return m
}

func (m PageContentModel) Update(msg tea.Msg) (PageContentModel, tea.Cmd) {
	switch msg := msg.(type) {
	case notion.PageContentMsg:
		if msg.Err != nil {
			m.error = msg.Err
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
