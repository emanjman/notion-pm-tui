package notion

import "strings"

// helper func to return arr of rich text into single plain text (string)
func ExtractPlainText(richTexts []RichText) string {
	var s strings.Builder
	for _, rt := range richTexts {
		s.WriteString(rt.PlainText)
	}
	return s.String()
}
