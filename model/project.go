package model

type Tab int

// enum representation for better readability
const (
	MilestonesTab Tab = iota
	OverviewTab
	ProjectNotesTab
	DebugNotesTab
)

type ProjectModel struct {
	activeTab    Tab
	milestones   views.MilestonesListModel
	overview     views.PageContentModel
	projectNotes views.NotesListModel
	debugNotes   views.NotesListModel
}
