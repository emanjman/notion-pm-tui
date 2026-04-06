package milestone

import "github.com/charmbracelet/bubbles/list"

type LoadMoreItem struct {
	Status  string
	Loading bool
}

var _ list.Item = (*LoadMoreItem)(nil) // conform

func (l LoadMoreItem) FilterValue() string { return "" }
