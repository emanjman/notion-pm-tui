package notion

import (
	"fmt"
	"notion-project-tui/styles"
	"strings"

	lg "github.com/charmbracelet/lipgloss"
)

func RenderBlocks(bs []Block, windowWidth int, depth int) string {
	var s strings.Builder
	counter := 0
	var counterType *ListFormatType = nil

	for i, b := range bs {
		// skip possible block types to hide
		if (i == 1 && b.Type == Divider) || b.Type == Breadcrumb {
			continue
		}

		if b.Type == NumberedListItem {
			counter++
			if b.NumberedListItem.ListFormat != nil {
				counterType = b.NumberedListItem.ListFormat
			}
		} else {
			counter = 0
			counterType = nil
		}

		s.WriteString(renderBlock(b, windowWidth, depth, counter, counterType))
		s.WriteString("\n")

		if b.HasChildren && b.Type != Callout {
			s.WriteString(RenderBlocks(b.Children, windowWidth, depth+1))
		}
	}

	return s.String()
}

func renderBlock(b Block, windowWidth int, depth int, counter int, counterType *ListFormatType) string {
	base := lg.NewStyle().PaddingLeft(depth * 3)

	switch b.Type {
	case Divider:
		return base.
			Foreground(styles.BorderForeground).
			PaddingTop(1).
			PaddingBottom(1).
			Render(strings.Repeat("—", windowWidth))

	case Callout:
		content := ExtractPlainText(b.Callout.RichText)

		if b.HasChildren {
			for _, child := range b.Children {
				content += renderBlock(child, windowWidth, depth, counter, counterType)
			}
		}

		return base.
			Background(styles.SelectedBackground).
			Render(content)

	case Heading2:
		txt := ExtractPlainText(b.Heading2.RichText)
		parts := strings.Fields(txt)
		icon, header := parts[0], parts[1]

		return base.
			Bold(true).
			Foreground(styles.PrimaryForeground).
			Render(fmt.Sprintf(" %s %s ", icon, header))

	case Heading3:
		return base.
			Bold(true).
			Foreground(styles.PrimaryForeground).
			Render("\n" + ExtractPlainText(b.Heading3.RichText) + "\n")

	case BulletedListItem:
		pt := lg.NewStyle().
			Foreground(styles.MutedForeground).
			Render("-  ")
		txt := lg.NewStyle().
			Foreground(styles.TechForeground).
			Render(ExtractPlainText(b.BulletedListItem.RichText))
		return base.Render(pt + txt)

	case NumberedListItem:
		var pt string
		format := getListFormat(counterType, depth)

		switch format {
		case Numbers:
			pt = fmt.Sprintf("%d ", counter)
		case Letters:
			pt = fmt.Sprintf("%s ", string(letterPoints[counter-1]))
		case Roman:
			pt = fmt.Sprintf("%s.", toRomanNumeral(counter))
		}

		pt = lg.NewStyle().
			Foreground(styles.MutedForeground).
			Render(pt)
		txt := lg.NewStyle().
			Foreground(styles.TechForeground).
			Render(ExtractPlainText(b.NumberedListItem.RichText))
		return base.Render(fmt.Sprintf("%s %s", pt, txt))

	case Toggle:
		chevron := lg.NewStyle().Foreground(styles.MutedForeground).Render("▼  ")
		return base.Render(chevron + ExtractPlainText(b.Toggle.RichText))

	case Paragraph:
		return base.Render(ExtractPlainText(b.Paragraph.RichText))
	}

	return lg.NewStyle().
		Foreground(styles.MutedForeground).
		Render("<Unhandled block>")
}

func getListFormat(explicit *ListFormatType, depth int) ListFormatType {
	if explicit != nil {
		return *explicit
	}

	switch depth % 3 {
	case 0:
		return Numbers
	case 1:
		return Letters
	case 2:
		return Roman
	}
	return Numbers
}

var letterPoints = "abcdefghijklmnopqrstuvwxyz"

func toRomanNumeral(num int) string {
	if num <= 0 || num > 3999 {
		return ""
	}

	values := []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
	symbols := []string{"m", "cm", "d", "cd", "c", "xc", "l", "xl", "x", "ix", "v", "iv", "i"}

	var result strings.Builder
	for i := 0; i < len(values); i++ {
		for num >= values[i] {
			result.WriteString(symbols[i])
			num -= values[i]
		}
	}
	return result.String()
}
