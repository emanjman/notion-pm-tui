package milestone

import (
	"fmt"
	"notion-project-tui/notion"
	"notion-project-tui/styles"
	"strings"

	lg "github.com/charmbracelet/lipgloss"
)

func createProgressBar(progress float64, width int, baseStyle lg.Style) string {
	wfilled := int(progress * float64(width))
	wempty := width - wfilled

	filled := strings.Repeat("▬", wfilled)
	empty := strings.Repeat("▬", wempty)

	return baseStyle.Foreground(styles.TechForeground).Render(filled) +
		baseStyle.Foreground(lg.Color("#2a2a2a")).Render(empty)
}

func renderItem(d ItemDelegate, item DefaultItem, selected bool, noBorder bool, windowWidth int) string {
	contStyle := d.style.itemContainer.base
	segStyle := d.style.itemSegment.base
	nameStyle, stateStyle, countStyle := lg.Style{}, lg.Style{}, lg.Style{}

	// handle field highlighting by mode
	if selected {
		if d.focus.Mode == NeutralMode {
			segStyle = d.style.itemSegment.selected
			countStyle = d.style.itemContainer.selected

			// apply select highlight row-wide
			nameStyle, stateStyle, countStyle = segStyle, segStyle, segStyle
		} else {
			nameStyle = nameStyle.Inherit(d.style.itemSegment.selected)
			countStyle = countStyle.Inherit(d.style.itemSegment.selected)
		}
	}

	// apply final field-specific styles
	nameStyle = nameStyle.Foreground(styles.PrimaryForeground)
	countStyle = countStyle.Foreground(styles.MutedForeground)

	if item.Icon == "" {
		item.Icon = "  "
	}
	name := nameStyle.Render(item.Icon + " " + item.Name)
	count := countStyle.Render(fmt.Sprintf("%d", item.TaskCount))

	var state string
	switch item.FetchStatus {
	case FetchIdle:
		state = stateStyle.Foreground(styles.MutedForeground).Render("◌")
	case FetchPending:
		state = stateStyle.Foreground(styles.MutedForeground).Render("↻")
	case FetchFailed:
		state = stateStyle.Foreground(lg.Color("#e0af68")).Render("⚠")
	}
	space := segStyle.Render(" ")

	// hide progress bar for completed milestones
	var progress string
	if item.MilestoneStatus != notion.MilestoneComplete {
		// completion := segStyle.
		// 	Foreground(styles.MutedForeground).
		// 	Render(fmt.Sprintf("%.0f%%", item.Progress*100))
		pbar := createProgressBar(item.Progress, windowWidth/4, segStyle)
		progress = pbar + segStyle.Render(" ") + count
	}

	// calculate max title width
	leftOffset, rightOffset := 3, 2
	offset := leftOffset + rightOffset
	nameMaxWidth := windowWidth - lg.Width(progress+space+state) - offset

	if selected && d.focus.Mode == EditMode {
		// use textinput component in edit mode
		d.focus.tempTitle.Width = nameMaxWidth
		name = d.focus.tempTitle.View()
	} else if lg.Width(name) > nameMaxWidth {
		// if past the max width, truncate until valid
		n := item.Name
		for lg.Width(n+"...") > nameMaxWidth && len(n) > 0 {
			n = n[:len(n)-1]
		}
		n = n + "..."
		name = nameStyle.Render(n)
	}

	if noBorder {
		contStyle = contStyle.Border(lg.NormalBorder(), false)
	}

	left := name
	right := progress + space + state
	px := styles.GetPaddingBetween(left, right, windowWidth, contStyle)
	content := left + styles.RenderPadding(segStyle, px) + right

	rendered := contStyle.Width(windowWidth).Render(content)
	if noBorder {
		return rendered + "\n" + lg.NewStyle().Render("")
	}
	return rendered
}
