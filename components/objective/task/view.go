package task

func (m Model) View() string {
	// ! temp, styling
	// if m.loading {
	// 	return "Loading tasks..."
	// }

	return m.list.View()
}
