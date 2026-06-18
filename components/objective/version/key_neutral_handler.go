package version

import tea "github.com/charmbracelet/bubbletea"

func (m Model) onNeutralPrev() (Model, tea.Cmd) {
	n := len(m.pages)
	if m.pageIdx == 0 {
		m.pageIdx = n - 1
	} else {
		m.pageIdx -= 1
	}

	m.CurrPage = &m.pages[m.pageIdx]
	versionID := m.CurrPage.ID

	return m, fetchInitVersionMilestones(versionID, m.notion)
}

func (m Model) onNeutralNext() (Model, tea.Cmd) {
	n := len(m.pages)
	m.pageIdx = (m.pageIdx + 1) % n

	m.CurrPage = &m.pages[m.pageIdx]
	versionID := m.CurrPage.ID

	return m, fetchInitVersionMilestones(versionID, m.notion)
}
