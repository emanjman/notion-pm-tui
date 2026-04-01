package styles

import (
	"strings"

	lg "github.com/charmbracelet/lipgloss"
)

// foreground
var (
	PrimaryForeground = lg.Color("#c0caf5")
	MutedForeground   = lg.Color("#565f89")
	TechForeground    = lg.Color("#76FF83")
	RedForeground     = lg.Color("#ff6369")
	YellowForeground  = lg.Color("#ffed4f")
)

// background
var (
	SelectedBackground = lg.Color("#1d1d1f")
	ErrorBackground    = lg.Color("#2d1a1a")
)

// borders
var (
	BorderForeground = lg.Color("236")
)

// RenderWithTitle renders content with a rounded border and embeds title into the top edge.
func RenderWithTitle(title string, content string, style lg.Style) string {
	rendered := style.Render(content)
	lines := strings.Split(rendered, "\n")
	if len(lines) == 0 {
		return rendered
	}

	topLine := []rune(lines[0])
	insert := []rune("─ " + title + " ")

	if len(insert) < len(topLine)-1 {
		copy(topLine[1:], insert)
		lines[0] = string(topLine)
	}

	return strings.Join(lines, "\n")
}

// util
func GetPaddingBetween(left string, right string, windowWidth int, externalStyle lg.Style) int {
	availWidth := windowWidth - externalStyle.GetHorizontalPadding()

	padWidth := availWidth - lg.Width(left) - lg.Width(right)
	if padWidth < 0 {
		padWidth = 0
	}

	return padWidth
}

func RenderPadding(segmentStyle lg.Style, width int) string {
	return segmentStyle.Render(strings.Repeat(" ", width))
}
