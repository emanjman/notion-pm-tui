package notion

import "strings"

// AddBlockSpacing inserts blank lines between top-level blocks for readability.
// Grouping rules:
//   - List items (- / 1. etc.) are kept tight with adjacent list items
//   - Table rows (|) are kept tight with adjacent table rows
//   - Toggle blocks (<details>) are kept tight with adjacent toggles
//   - Callout blocks (<callout>) are isolated — blank lines around them
//   - Code blocks (```) and the interiors of <details>/<callout> are emitted as-is
func AddBlockSpacing(md string) string {
	lines := strings.Split(md, "\n")
	out := make([]string, 0, len(lines)*2)

	inCode := false
	inCallout := false
	inDetails := false

	prevTopLine := "" // last emitted line that was not inside a fence

	isListLine := func(l string) bool {
		t := strings.TrimSpace(l)
		if strings.HasPrefix(t, "- ") || strings.HasPrefix(t, "* ") {
			return true
		}
		// numbered: "1. " "a. " "i. " etc.
		for i, c := range t {
			if c == '.' && i > 0 {
				return strings.HasPrefix(t[i:], ". ")
			}
			if c != ' ' && !(c >= 'a' && c <= 'z') && !(c >= 'A' && c <= 'Z') && !(c >= '0' && c <= '9') {
				break
			}
		}
		return false
	}

	isTableRow := func(l string) bool {
		return strings.HasPrefix(strings.TrimSpace(l), "|")
	}

	isDetailsOpen := func(l string) bool { return strings.TrimSpace(l) == "<details>" }
	isDetailsClose := func(l string) bool { return strings.TrimSpace(l) == "</details>" }
	isCalloutOpen := func(l string) bool { return strings.HasPrefix(strings.TrimSpace(l), "<callout") }
	isCalloutClose := func(l string) bool { return strings.TrimSpace(l) == "</callout>" }
	isCodeFence := func(l string) bool { return strings.HasPrefix(strings.TrimSpace(l), "```") }

	needsSpaceBefore := func(curr, prev string) bool {
		if prev == "" {
			return false
		}
		// callout open/close always gets a blank line before it
		if isCalloutOpen(curr) || isCalloutClose(curr) {
			return true
		}
		// blank line after callout close
		if isCalloutClose(prev) {
			return true
		}
		// group: list with list
		if isListLine(curr) && isListLine(prev) {
			return false
		}
		// group: table row with table row
		if isTableRow(curr) && isTableRow(prev) {
			return false
		}
		// group: details open tight with prior details close
		if isDetailsOpen(curr) && isDetailsClose(prev) {
			return false
		}
		return true
	}

	for _, line := range lines {
		switch {
		case !inCode && !inCallout && !inDetails && isCodeFence(line):
			inCode = true
			if needsSpaceBefore(line, prevTopLine) {
				out = append(out, "")
			}
			out = append(out, line)
			prevTopLine = line

		case inCode && isCodeFence(line):
			inCode = false
			out = append(out, line)
			prevTopLine = line

		case !inCode && !inCallout && !inDetails && isCalloutOpen(line):
			inCallout = true
			if needsSpaceBefore(line, prevTopLine) {
				out = append(out, "")
			}
			out = append(out, line)

		case inCallout && isCalloutClose(line):
			inCallout = false
			out = append(out, line)
			prevTopLine = line

		case !inCode && !inCallout && !inDetails && isDetailsOpen(line):
			inDetails = true
			if needsSpaceBefore(line, prevTopLine) {
				out = append(out, "")
			}
			out = append(out, line)

		case inDetails && isDetailsClose(line):
			inDetails = false
			out = append(out, line)
			prevTopLine = line

		case inCode || inCallout || inDetails:
			// inside a fence — emit as-is, no spacing logic
			out = append(out, line)

		default:
			if needsSpaceBefore(line, prevTopLine) {
				out = append(out, "")
			}
			out = append(out, line)
			if line != "" {
				prevTopLine = line
			}
		}
	}

	return strings.Join(out, "\n")
}

type MDSuccessRes struct {
	Object          string   `json:"object"`
	ID              string   `json:"id"`
	Markdown        string   `json:"markdown"`
	Truncated       bool     `json:"truncated"`
	UnknownBlockIDs []string `json:"unknown_block_ids"`
}

type MDFailRes struct {
	Object  string `json:"object"`
	Message string `json:"message"`
	Code    string `json:"code"`
	Status  int    `json:"status"`
}

type ReplaceContent struct {
	NewStr string `json:"new_str"`
}
type MDReplaceReq struct {
	Type           string         `json:"type"` // replace_content
	ReplaceContent ReplaceContent `json:"replace_content"`
}
