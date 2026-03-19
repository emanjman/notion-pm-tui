package notion

import "strings"

func BlocksToMarkdown(blocks []Block, depth int) string {
	var s strings.Builder
	for _, b := range blocks {
		s.WriteString(blockToMarkdown(b, depth))
		if b.HasChildren {
			s.WriteString(BlocksToMarkdown(b.Children, depth+1))
		}
	}
	return s.String()
}

func blockToMarkdown(b Block, depth int) string {
	indent := strings.Repeat("  ", depth)
	switch b.Type {
	case Paragraph:
		return indent + ExtractPlainText(b.Paragraph.RichText) + "\n\n"
	case Heading2:
		return "## " + ExtractPlainText(b.Heading2.RichText) + "\n\n"
	case Heading3:
		return "### " + ExtractPlainText(b.Heading3.RichText) + "\n\n"
	case BulletedListItem:
		return indent + "- " + ExtractPlainText(b.BulletedListItem.RichText) + "\n"
	case NumberedListItem:
		return indent + "1. " + ExtractPlainText(b.NumberedListItem.RichText) + "\n"
	case Code:
		lang := b.Code.Language
		content := ExtractPlainText(b.Code.RichText)
		return "```" + lang + "\n" + content + "\n```\n\n"
	case Divider:
		return "---\n\n"
	}
	return ""
}
