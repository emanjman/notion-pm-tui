package notebook

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type VimMode int

const (
	NormalMode VimMode = iota
	InsertMode
)

func k(t tea.KeyType) tea.KeyMsg        { return tea.KeyMsg{Type: t} }
func altK(r rune) tea.KeyMsg            { return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}, Alt: true} }
func runeK(r rune) tea.KeyMsg           { return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}} }

func sendKeys(ta textarea.Model, keys ...tea.KeyMsg) (textarea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	for _, msg := range keys {
		var cmd tea.Cmd
		ta, cmd = ta.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return ta, tea.Batch(cmds...)
}

// handleNormalMode processes a key in vim Normal mode. Returns the updated
// model and whether a state transition to InsertMode occurred.
func handleNormalMode(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	ch := msg.String()

	// --- pending two-key sequences ---
	if m.pendingVimKey != "" {
		prefix := m.pendingVimKey
		m.pendingVimKey = ""
		switch prefix + ch {
		case "dd": // delete line: go to line start, delete to EOL, delete newline
			var cmd tea.Cmd
			m.editor, cmd = sendKeys(m.editor,
				k(tea.KeyCtrlA),  // line start
				k(tea.KeyCtrlK),  // delete to EOL
				k(tea.KeyCtrlD),  // delete newline char
			)
			return m, cmd
		case "dw": // delete word forward
			var cmd tea.Cmd
			m.editor, cmd = sendKeys(m.editor, altK('d'))
			return m, cmd
		case "gg": // go to input begin
			var cmd tea.Cmd
			m.editor, cmd = sendKeys(m.editor, k(tea.KeyCtrlHome))
			return m, cmd
		}
		// unrecognised sequence — fall through and treat second key normally
	}

	switch ch {
	// --- mode switches ---
	case "i":
		m.vimMode = InsertMode

	case "a": // insert after cursor
		m.editor, _ = sendKeys(m.editor, k(tea.KeyRight))
		m.vimMode = InsertMode

	case "A": // insert at end of line
		m.editor, _ = sendKeys(m.editor, k(tea.KeyCtrlE))
		m.vimMode = InsertMode

	case "I": // insert at start of line
		m.editor, _ = sendKeys(m.editor, k(tea.KeyCtrlA))
		m.vimMode = InsertMode

	case "o": // open line below
		m.editor, _ = sendKeys(m.editor, k(tea.KeyCtrlE), k(tea.KeyEnter))
		m.vimMode = InsertMode

	case "O": // open line above
		m.editor, _ = sendKeys(m.editor, k(tea.KeyCtrlA), k(tea.KeyEnter), k(tea.KeyUp))
		m.vimMode = InsertMode

	// --- motion ---
	case "h":
		var cmd tea.Cmd
		m.editor, cmd = sendKeys(m.editor, k(tea.KeyLeft))
		return m, cmd
	case "l":
		var cmd tea.Cmd
		m.editor, cmd = sendKeys(m.editor, k(tea.KeyRight))
		return m, cmd
	case "j":
		var cmd tea.Cmd
		m.editor, cmd = sendKeys(m.editor, k(tea.KeyDown))
		return m, cmd
	case "k":
		var cmd tea.Cmd
		m.editor, cmd = sendKeys(m.editor, k(tea.KeyUp))
		return m, cmd
	case "w":
		var cmd tea.Cmd
		m.editor, cmd = sendKeys(m.editor, altK('f'))
		return m, cmd
	case "b":
		var cmd tea.Cmd
		m.editor, cmd = sendKeys(m.editor, altK('b'))
		return m, cmd
	case "0":
		var cmd tea.Cmd
		m.editor, cmd = sendKeys(m.editor, k(tea.KeyCtrlA))
		return m, cmd
	case "$":
		var cmd tea.Cmd
		m.editor, cmd = sendKeys(m.editor, k(tea.KeyCtrlE))
		return m, cmd
	case "G": // go to input end
		var cmd tea.Cmd
		m.editor, cmd = sendKeys(m.editor, k(tea.KeyCtrlEnd))
		return m, cmd

	// --- edit in normal mode ---
	case "x": // delete char forward
		var cmd tea.Cmd
		m.editor, cmd = sendKeys(m.editor, k(tea.KeyCtrlD))
		return m, cmd
	case "X": // delete char backward
		var cmd tea.Cmd
		m.editor, cmd = sendKeys(m.editor, k(tea.KeyBackspace))
		return m, cmd
	case "D": // delete to EOL
		var cmd tea.Cmd
		m.editor, cmd = sendKeys(m.editor, k(tea.KeyCtrlK))
		return m, cmd
	case "p": // paste from clipboard (textarea handles ctrl+v)
		var cmd tea.Cmd
		m.editor, cmd = sendKeys(m.editor, runeK('v'))
		return m, cmd

	// --- pending prefix keys ---
	case "d", "g":
		m.pendingVimKey = ch
	}

	return m, nil
}
