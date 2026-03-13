package notion

import (
	"notion-project-tui/styles"
	"strings"

	lg "github.com/charmbracelet/lipgloss"
)

type RichText struct {
	PlainText   string `json:"plain_text"`
	Annotations struct {
		Bold          bool  `json:"bold"`
		Italic        bool  `json:"italic"`
		Strikethrough bool  `json:"strikethrough"`
		Underline     bool  `json:"underline"`
		InlineCode    bool  `json:"code"`
		Color         Color `json:"color"`
	} `json:"annotations"`
}

// helper func to return arr of rich text into single plain text (string)
func ExtractPlainText(richTexts []RichText) string {
	var s strings.Builder
	for _, txt := range richTexts {
		style := lg.NewStyle()

		if txt.Annotations.InlineCode {
			style = style.
				Background(styles.SelectedBackground).
				Foreground(styles.RedForeground)
		}

		switch txt.Annotations.Color {
		case Red:
			style = style.Foreground(styles.RedForeground)
		case Yellow:
			style = style.Foreground(styles.YellowForeground)
			// case Default:
			// 	if txt.Annotations.InlineCode {
			// 		style = style.Foreground(lg.Color("#ff6369"))
			// 	} else {
			// 		style = style.Foreground(styles.TechForeground)
			// 	}
		}

		s.WriteString(style.Render(txt.PlainText))
	}
	return s.String()
}
