package styles

import (
	"strings"

	lg "github.com/charmbracelet/lipgloss"
)

// foreground
var (
	PrimaryForeground = lg.Color("#c0caf5")
	MutedForeground   = lg.Color("#565f89")
)

// background
var (
	SelectedBackground = lg.Color("#2e2e32")
)

// borders
var (
	BorderForeground = lg.Color("236")
)

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
