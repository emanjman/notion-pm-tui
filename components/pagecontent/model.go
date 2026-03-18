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
	PageID   string
}

func New(pageID string, n *notion.Client) Model {
	vp := viewport.New(0, 0)

	m := Model{
		viewport: vp,
		notion:   n,
		loading:  true,
	}

	if pageID != "" {
		m.PageID = pageID
	}

	return m
}

func (m Model) Init() tea.Cmd {
	if m.PageID != "" {
		return func() tea.Msg {
			blocks, err := m.notion.FetchPageContent(m.PageID)
			return notion.PageContentMsg{Data: blocks, Err: err}
		}
	}
	return nil
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
	content := ""
	if m.PageID != "" {
		content = m.viewport.View()
	} else {
		content = "No content loaded..."
	}
	return style.Render(content)
}
