package notion

import (
	"notion-project-tui/notion"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type PageContentModel struct {
	viewport viewport.Model
	blocks   []notion.Block
	notion   *notion.Client
	loading  bool
}

func NewContentModel(n *notion.Client) PageContentModel {
	v := viewport.New(0, 0)
	return PageContentModel{
		viewport: v,
		blocks:   []notion.Block{},
		notion:   n,
		loading:  true,
	}
}

func (m PageContentModel) Init() tea.Cmd {
	return nil
}

func (m PageContentModel) View() string {
	return m.viewport.View()
}

func (m PageContentModel) Update() (PageContentModel, tea.Cmd) {
	// ! mock data
	m.viewport.SetContent(m.renderBlocks())
	return m, nil
}

func (m PageContentModel) renderBlock(block notion.Block) string {
	switch block.Type {
	case notion.Paragraph:
		return notion.ExtractPlainText(block.Paragraph.RichText)
	case notion.Divider:
		return strings.Repeat("─", m.viewport.Width)
	}

	return "Unhandled block"
}

func (m PageContentModel) renderBlocks() string {
	var s strings.Builder
	for _, block := range m.blocks {
		s.WriteString(m.renderBlock(block))
		s.WriteString("\n")
	}
	return s.String()
}
