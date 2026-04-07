package milestone

func (m Model) View() string {
	if m.pendingFetches > 0 {
		return "Loading milestones..."
	}
	return m.list.View()
}
