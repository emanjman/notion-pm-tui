package pagecontent

import (
	"fmt"
	"notion-project-tui/notion"
	"notion-project-tui/styles"
	"strings"

	lg "github.com/charmbracelet/lipgloss"
)

func renderBlocks(bs []notion.Block, windowWidth int) string {
	var s strings.Builder
	counter := 0
	var counterType *notion.ListFormatType = nil

	for _, b := range bs {
		if b.Type == notion.NumberedListItem {
			counter++
			if b.NumberedListItem.ListFormat != nil {
				counterType = b.NumberedListItem.ListFormat
			}
		} else {
			counter = 0
			counterType = nil
		}

		s.WriteString(renderBlock(b, windowWidth, 0, counter, counterType))
		s.WriteString("\n")
	}
	return s.String()
}

func renderBlock(b notion.Block, windowWidth int, depth int, counter int, counterType *notion.ListFormatType) string {
	switch b.Type {
	case notion.Divider:
		return lg.NewStyle().
			Foreground(styles.BorderForeground).
			Render(strings.Repeat("—", windowWidth))

	case notion.Callout:
		return lg.NewStyle().
			Border(lg.NormalBorder(), true, true, true, true).
			BorderForeground(styles.BorderForeground).
			Render("expect child here")

	case notion.Heading2:
		return lg.NewStyle().
			Bold(true).
			Render("\n" + notion.ExtractPlainText(b.Heading2.RichText) + "\n")

	case notion.Heading3:
		return lg.NewStyle().
			Bold(true).
			Underline(true).
			Render("\n" + notion.ExtractPlainText(b.Heading3.RichText) + "\n")

	case notion.BulletedListItem:
		return "• " + notion.ExtractPlainText(b.BulletedListItem.RichText)

	case notion.NumberedListItem:
		var pt string

		if counterType != nil {
			if *counterType == notion.Numbers {
				pt = fmt.Sprintf("%d.", counter)
			} else if *counterType == notion.Letters {
				pt = letterPoints[counter]
			} else if *counterType == notion.Roman {
				pt = romanPoints[counter]
			}
		}

		return fmt.Sprintf("%s %s", pt, notion.ExtractPlainText(b.NumberedListItem.RichText))

	case notion.Paragraph:
		return notion.ExtractPlainText(b.Paragraph.RichText)
	}

	return "--"
}

var letterPoints = map[int]string{
	1: "a.",
	2: "b.",
	3: "c.",
	4: "d.",
	5: "e.",
}

var romanPoints = map[int]string{
	1: "i.",
	2: "ii.",
	3: "iii.",
	4: "iv.",
	5: "v.",
}
