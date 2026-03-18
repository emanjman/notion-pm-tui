package pagecontent

import (
	"log"
	"notion-project-tui/notion"
	"notion-project-tui/styles"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type SwitchContentMsg struct {
	Content string
}

type Model struct {
	Viewport viewport.Model
	notion   *notion.Client
	loading  bool
	error    error
	Content  string
}

func New(pageID string, n *notion.Client) Model {
	vp := viewport.New(0, 0)
	m := Model{
		Viewport: vp,
		notion:   n,
		loading:  true,
		error:    nil,
		Content:  "",
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	case SwitchContentMsg:
		log.Printf("page content received [SwitchContentMsg]")
		m.Viewport.SetContent(msg.Content)
		m.loading = false

	case tea.WindowSizeMsg:
		m.Viewport.Width = msg.Width
		m.Viewport.Height = msg.Height

	// forward other commands to list
	default:
		var cmd tea.Cmd
		m.Viewport, cmd = m.Viewport.Update(msg)
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
		loadingStyle := lg.NewStyle().Height(m.Viewport.Height)
		return loadingStyle.Render("loading...")
	}

	style := lg.NewStyle().Padding(0, 1)
	return style.Render(m.Viewport.View())
}
