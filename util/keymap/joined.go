package keymap

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type JoinedKeyMap struct {
	Primary   help.KeyMap
	Secondary help.KeyMap
}

func (m JoinedKeyMap) ShortHelp() []key.Binding {
	return append(m.Primary.ShortHelp(), m.Secondary.ShortHelp()...)
}

func (m JoinedKeyMap) FullHelp() [][]key.Binding {
	return append(m.Primary.FullHelp(), m.Secondary.FullHelp()...)
}
