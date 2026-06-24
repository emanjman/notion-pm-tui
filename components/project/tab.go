package project

type Tab int

// enum representation for better readability
const (
	ObjectiveTab Tab = iota
	NotebookTab
	BugsTab
	TechTab
)
const tabCount = 4 // todo: can inject labels/cnt into model itself?
