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
	SelectedBackground = lg.Color("#1d1d1f")
)

// borders
var (
	BorderForeground = lg.Color("236")
)

// util
func PadBetween(left, right string, windowWidth int, extStyle lg.Style) string {
	// use lg.Width to only consider visible cells
	padding := windowWidth - lg.Width(left) - lg.Width(right) - extStyle.GetHorizontalPadding()
	if padding < 0 {
		padding = 0
	}

	return left + strings.Repeat(" ", padding) + right
}
