package notion

import (
	"notion-project-tui/styles"
	"strings"

	lg "github.com/charmbracelet/lipgloss"
)

type TextContent struct {
	Content string `json:"content"`
}

type Annotations struct {
	Bold          bool  `json:"bold"`
	Italic        bool  `json:"italic"`
	Strikethrough bool  `json:"strikethrough"`
	Underline     bool  `json:"underline"`
	InlineCode    bool  `json:"code"`
	Color         Color `json:"color"`
}

type RichText struct {
	Text        TextContent  `json:"text"`
	PlainText   string       `json:"plain_text,omitempty"`
	Annotations *Annotations `json:"annotations,omitempty"`
}

// helper func to return arr of rich text into single plain text (string)
func ExtractPlainText(richTexts []RichText) string {
	var s strings.Builder
	for _, txt := range richTexts {
		style := lg.NewStyle()

		if txt.Annotations != nil {
			if txt.Annotations.Italic {
				style = style.Italic(true)
			}
			if txt.Annotations.Bold {
				style = style.Bold(true)
			}
			if txt.Annotations.Underline {
				style = style.Underline(true)
			}
			if txt.Annotations.Strikethrough {
				style = style.Strikethrough(true)
			}
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
			}
		}

		s.WriteString(style.Render(txt.PlainText))
	}
	return s.String()
}
