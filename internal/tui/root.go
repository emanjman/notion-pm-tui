package tui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/emanuelpecson/notion-project-tui/internal/config"
	"github.com/emanuelpecson/notion-project-tui/internal/models"
	"github.com/emanuelpecson/notion-project-tui/internal/notion"
)

// ViewType represents which view we're currently on
type ViewType int

const (
	ViewProjects ViewType = iota
	ViewTaskGroups
	ViewTasks
	ViewLoading
)

// RootModel is the root Bubble Tea model that manages different views
type RootModel struct {
	notionClient    *notion.Client
	config          *config.Config
	currentView     ViewType
	projectsList    ProjectsListModel
	taskGroupsList  TaskGroupsListModel
	tasksList       TasksListModel
	loadingMessage  string
	cachedTasks     map[string][]models.Task // Cache tasks by task group ID
}

// NewRootModel creates the root model
func NewRootModel(projects []models.Project, notionClient *notion.Client, cfg *config.Config) RootModel {
	return RootModel{
		notionClient: notionClient,
		config:       cfg,
		currentView:  ViewProjects,
		projectsList: NewProjectsListModel(projects),
		cachedTasks:  make(map[string][]models.Task),
	}
}

// taskGroupsFetchedMsg is sent when task groups have been fetched
type taskGroupsFetchedMsg struct {
	project         models.Project
	taskGroups      []models.TaskGroup
	duration        time.Duration
	err             error
	prefetchedTasks []models.Task              // Tasks for the first task group (backward compat)
	prefetchedForID string                     // ID of the task group we prefetched for (backward compat)
	cachedTasks     map[string][]models.Task   // All pre-fetched tasks by task group ID
}

// tasksFetchedMsg is sent when tasks have been fetched
type tasksFetchedMsg struct {
	project   models.Project
	taskGroup models.TaskGroup
	tasks     []models.Task
	duration  time.Duration
	err       error
}

// fetchTaskGroupsCmd is a Bubble Tea command that fetches task groups
// Pre-fetching strategy:
// - Priority projects: pre-fetch ALL task groups' tasks
// - Other projects: pre-fetch only "under development" task groups' tasks
func (m RootModel) fetchTaskGroupsCmd(project models.Project) tea.Cmd {
	return func() tea.Msg {
		start := time.Now()

		ids := project.GetBranchIDs()
		if len(ids) == 0 {
			return taskGroupsFetchedMsg{
				project:    project,
				taskGroups: []models.TaskGroup{},
				duration:   time.Since(start),
				err:        nil,
			}
		}

		taskGroups, err := m.notionClient.FetchTaskGroups(ids)
		duration := time.Since(start)

		if err != nil {
			return taskGroupsFetchedMsg{
				project:    project,
				taskGroups: taskGroups,
				duration:   duration,
				err:        err,
			}
		}

		// Determine pre-fetch strategy based on project priority
		projectID := project.GetID()
		isPriority := m.config.IsPriorityProject(projectID)
		cachedTasks := make(map[string][]models.Task)

		// Debug output
		fmt.Printf("[DEBUG] Project ID: %s\n", projectID)
		fmt.Printf("[DEBUG] Is Priority: %v\n", isPriority)
		fmt.Printf("[DEBUG] Config Priority IDs: %v\n", m.config.GetPriorityProjectIDs())

		if isPriority {
			// Priority project: pre-fetch ALL task groups
			fmt.Printf("[DEBUG] Pre-fetching ALL %d task groups\n", len(taskGroups))
			for _, tg := range taskGroups {
				taskIDs := tg.GetTaskIDs()
				if len(taskIDs) > 0 {
					tasks, taskErr := m.notionClient.FetchTasks(taskIDs)
					if taskErr == nil {
						cachedTasks[tg.GetID()] = tasks
					}
				}
			}
		} else {
			// Non-priority project: pre-fetch only "under development" task groups
			devCount := 0
			for _, tg := range taskGroups {
				status := tg.GetStatus()
				if status == "🚧 under development" {
					devCount++
					taskIDs := tg.GetTaskIDs()
					if len(taskIDs) > 0 {
						tasks, taskErr := m.notionClient.FetchTasks(taskIDs)
						if taskErr == nil {
							cachedTasks[tg.GetID()] = tasks
						}
					}
				}
			}
			fmt.Printf("[DEBUG] Pre-fetching %d 'under development' task groups\n", devCount)
		}

		// Keep backward compatibility - set first cached task group as prefetched
		var prefetchedTasks []models.Task
		var prefetchedForID string
		if len(cachedTasks) > 0 {
			// Use the first task group we cached (most recent by activity)
			if len(taskGroups) > 0 {
				firstID := taskGroups[0].GetID()
				if tasks, ok := cachedTasks[firstID]; ok {
					prefetchedTasks = tasks
					prefetchedForID = firstID
				}
			}
		}

		return taskGroupsFetchedMsg{
			project:         project,
			taskGroups:      taskGroups,
			duration:        duration,
			err:             err,
			prefetchedTasks: prefetchedTasks,
			prefetchedForID: prefetchedForID,
			cachedTasks:     cachedTasks,
		}
	}
}

// fetchTasksCmd is a Bubble Tea command that fetches tasks
func (m RootModel) fetchTasksCmd(project models.Project, taskGroup models.TaskGroup) tea.Cmd {
	return func() tea.Msg {
		start := time.Now()

		ids := taskGroup.GetTaskIDs()
		if len(ids) == 0 {
			return tasksFetchedMsg{
				project:   project,
				taskGroup: taskGroup,
				tasks:     []models.Task{},
				duration:  time.Since(start),
				err:       nil,
			}
		}

		tasks, err := m.notionClient.FetchTasks(ids)
		duration := time.Since(start)

		return tasksFetchedMsg{
			project:   project,
			taskGroup: taskGroup,
			tasks:     tasks,
			duration:  duration,
			err:       err,
		}
	}
}

// Init initializes the root model
func (m RootModel) Init() tea.Cmd {
	return m.projectsList.Init()
}

// Update handles messages
func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.currentView {
	case ViewProjects:
		return m.updateProjectsView(msg)
	case ViewTaskGroups:
		return m.updateTaskGroupsView(msg)
	case ViewTasks:
		return m.updateTasksView(msg)
	case ViewLoading:
		return m.updateLoadingView(msg)
	}
	return m, nil
}

// updateProjectsView handles updates for the projects list view
func (m RootModel) updateProjectsView(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			// Get selected project
			if m.projectsList.cursor < len(m.projectsList.items) {
				item := m.projectsList.items[m.projectsList.cursor]
				if item.Type == ItemTypeProject && item.Project != nil {
					// Switch to loading view and fetch task groups
					m.currentView = ViewLoading
					m.loadingMessage = "Loading task groups..."
					return m, m.fetchTaskGroupsCmd(*item.Project)
				}
			}
		}
	}

	// Delegate to projects list
	var cmd tea.Cmd
	updatedModel, cmd := m.projectsList.Update(msg)
	m.projectsList = updatedModel.(ProjectsListModel)
	return m, cmd
}

// updateTaskGroupsView handles updates for the task groups view
func (m RootModel) updateTaskGroupsView(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" {
			// Go back to projects list
			m.currentView = ViewProjects
			return m, nil
		}

		// Check if Enter is pressed on a task group (not a header)
		if msg.String() == "enter" {
			if m.taskGroupsList.cursor < len(m.taskGroupsList.items) {
				item := m.taskGroupsList.items[m.taskGroupsList.cursor]
				if item.Type == TaskGroupItemTypeTaskGroup && item.TaskGroup != nil {
					taskGroupID := item.TaskGroup.GetID()

					// Check if we have cached tasks for this task group
					if cachedTasks, ok := m.cachedTasks[taskGroupID]; ok {
						// Use cached tasks - no loading needed!
						m.tasksList = NewTasksListModel(m.taskGroupsList.project, *item.TaskGroup, cachedTasks, "(cached)")
						m.currentView = ViewTasks
						return m, nil
					}

					// No cache - fetch tasks
					m.currentView = ViewLoading
					m.loadingMessage = "Loading tasks..."
					return m, m.fetchTasksCmd(m.taskGroupsList.project, *item.TaskGroup)
				}
			}
		}
	}

	// Delegate to task groups list
	var cmd tea.Cmd
	updatedModel, cmd := m.taskGroupsList.Update(msg)
	m.taskGroupsList = updatedModel.(TaskGroupsListModel)
	return m, cmd
}

// updateTasksView handles updates for the tasks view
func (m RootModel) updateTasksView(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" {
			// Go back to task groups list
			m.currentView = ViewTaskGroups
			return m, nil
		}
	}

	// Delegate to tasks list
	var cmd tea.Cmd
	updatedModel, cmd := m.tasksList.Update(msg)
	m.tasksList = updatedModel.(TasksListModel)
	return m, cmd
}

// updateLoadingView handles updates for the loading view
func (m RootModel) updateLoadingView(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case taskGroupsFetchedMsg:
		if msg.err != nil {
			// For now, go back to projects list on error
			// TODO: Show error message
			m.currentView = ViewProjects
			return m, nil
		}

		// Cache all pre-fetched tasks
		for taskGroupID, tasks := range msg.cachedTasks {
			m.cachedTasks[taskGroupID] = tasks
		}

		// Format the duration
		loadTime := formatDuration(msg.duration)

		// Switch to task groups view
		m.taskGroupsList = NewTaskGroupsListModel(msg.project, msg.taskGroups, loadTime)
		m.currentView = ViewTaskGroups
		return m, nil

	case tasksFetchedMsg:
		if msg.err != nil {
			// For now, go back to task groups list on error
			// TODO: Show error message
			m.currentView = ViewTaskGroups
			return m, nil
		}

		// Format the duration
		loadTime := formatDuration(msg.duration)

		// Switch to tasks view
		m.tasksList = NewTasksListModel(msg.project, msg.taskGroup, msg.tasks, loadTime)
		m.currentView = ViewTasks
		return m, nil

	case tea.KeyMsg:
		// Allow quitting even while loading
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	return m, nil
}

// formatDuration formats a duration into a human-readable string
func formatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%dμs", d.Microseconds())
	}
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

// View renders the current view
func (m RootModel) View() string {
	switch m.currentView {
	case ViewProjects:
		return m.projectsList.View()
	case ViewTaskGroups:
		return m.taskGroupsList.View()
	case ViewTasks:
		return m.tasksList.View()
	case ViewLoading:
		return m.renderLoading()
	}
	return ""
}

// renderLoading renders a loading screen
func (m RootModel) renderLoading() string {
	style := lipgloss.NewStyle().
		Padding(2, 4).
		Foreground(lipgloss.Color("12"))
	return style.Render(m.loadingMessage)
}
