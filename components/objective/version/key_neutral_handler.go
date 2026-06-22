package version

import tea "github.com/charmbracelet/bubbletea"

func (m Model) onNeutralPrev() (Model, tea.Cmd) {
	n := len(m.pages)
	if m.PageIdx == 0 {
		m.PageIdx = n - 1
	} else {
		m.PageIdx -= 1
	}

	return m, nil
}

func (m Model) onNeutralNext() (Model, tea.Cmd) {
	n := len(m.pages)
	m.PageIdx = (m.PageIdx + 1) % n

	return m, nil
}
