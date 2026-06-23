package notion

// todo: eventually might wanna filter out archived projects
func queryProjectBody() map[string]any {
	return map[string]any{
		"page_size": 25, // hardcoded sane size
	}
}
