package pagecontent

import (
	"encoding/json"
	"notion-project-tui/notion"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
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

type BlockMsg struct {
	Data     []notion.Block
	Err      error
	Duration time.Duration
}

func (m PageContentModel) Init() tea.Cmd {
	start := time.Now()

	return func() tea.Msg {
		data, err := os.ReadFile("mock/keycloak-migration.json")
		if err != nil {
			return BlockMsg{Err: err, Duration: time.Since(start)}
		}

		var content notion.PageContent
		if err := json.Unmarshal(data, &content); err != nil {
			return BlockMsg{Err: err, Duration: time.Since(start)}
		}

		// todo: if HasMore is true, we gotta keep fetching
		return BlockMsg{Data: content.Results, Duration: time.Since(start)}
	}
}

func (m PageContentModel) View() string {
	if m.loading {
		return "Loading..."
	}
	return m.viewport.View()
}

func (m PageContentModel) Update(msg tea.Msg) (PageContentModel, tea.Cmd) {
	switch msg := msg.(type) {
	case BlockMsg:
		if msg.Err != nil {
			return m, nil
		}

		m.viewport.SetContent(renderBlocks(msg.Data, m.viewport.Width))
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
