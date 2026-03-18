package pagecontent

import (
	"notion-project-tui/notion"
	"notion-project-tui/styles"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type Model struct {
	viewport viewport.Model
	notion   *notion.Client
	loading  bool
	error    error
}

func New(n *notion.Client) Model {
	vp := viewport.New(0, 0)

	return Model{
		viewport: vp,
		notion:   n,
		loading:  true,
	}
}

func (m Model) Init() tea.Cmd {
	return m.notion.FetchPageContent("30eb7273944b80ad80c4f91a4f5cfd8d")
}

func (m Model) View() string {
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

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
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
