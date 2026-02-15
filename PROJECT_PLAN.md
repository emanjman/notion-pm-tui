# Notion Project TUI

## Overview

TUI application built with Bubble Tea framework to interact with Notion programming projects, enabling keyboard-driven workflows without switching to the web UI.

## Architecture

### Core components

- **TUI Layer**: Bubble Tea models and views for rendering
- **Notion Client**: API wrapper for Notion operations
- **Data Models**: structured representations of projects, task groups, tasks, notes
- **State Management**: application state and navigation
- **Configuration**: API credentials, workspace settings

### Tech stack

- `bubbletea` - TUI framework
- `lipgloss` - styling and layout
- `bubbles` - reusable components (list, table, viewport)
- `notion-sdk-go` or custom HTTP client for Notion API

## Data Model

### Project hierarchy

```
Project (Database Item)
├── Title
├── Icon (emoji) - optional
├── Status (e.g., Active, On Hold, Complete)
├── Task Groups (Named Sections/Pages)
│   ├── Task Group: "Database Migration"
│   │   ├── Icon (emoji) - optional
│   │   ├── Task
│   │   │   ├── Status (e.g., TODO, In Progress, Done)
│   │   │   ├── Text/Description
│   │   │   ├── Type (e.g., Feature, Bug, Refactor)
│   │   │   └── Metadata (created, updated, priority, etc.)
│   │   └── Task...
│   ├── Task Group: "Authentication"
│   │   ├── Icon (emoji) - optional
│   │   └── Tasks...
│   └── Task Group: "UI Improvements"
│       ├── Icon (emoji) - optional
│       └── Tasks...
├── Project Overview (Page)
├── Debug Notes (Database)
│   └── Note entries
└── General Notes (Database)
    └── Note entries
```

## Features

### Phase 1: Core functionality

#### Project navigation

- list all projects from dashboard
- display project emoji icons (if set)
- filter by status (active, on hold, complete)
- search projects by name
- select project to view details
- set/update project emoji icon

#### Task group navigation

- view all task groups within a project
- display task group emoji icons (if set)
- select task group to view tasks
- create new task group with optional emoji icon
- rename task group
- update task group emoji icon
- delete task group

#### Task management

- view all tasks within a task group
- filter tasks by status/type
- create new task
- update task status
- edit task text/description
- delete task

#### Basic operations

- quick status toggle (TODO -> In Progress -> Done)
- keyboard shortcuts for common actions
- real-time sync with Notion

### Phase 2: Enhanced features

#### Notes management

- view debug notes database
- view general notes database
- create new notes
- edit existing notes
- delete notes

#### Project overview

- display project overview page content
- edit overview (if supported)

#### Advanced filtering

- multi-filter support (status + type)
- custom filter queries
- saved filter presets

### Phase 3: Power features

#### Bulk operations

- bulk status updates
- bulk task creation
- bulk delete with confirmation

#### Templates

- task templates for common types
- task group templates
- project templates

#### Offline mode

- cache data locally
- queue operations when offline
- sync when connection restored

## UI Layout

### Main views

#### Projects list view

```
┌─────────────────────────────────────────────┐
│ Notion Projects                    [? Help] │
├─────────────────────────────────────────────┤
│                                             │
│ > Active Projects (3)                       │
│   🚀 Project Alpha             [Active]     │
│   💻 Project Beta              [Active]     │
│   🍾 Project Gamma              [Active]     │
│                                             │
│   On Hold Projects (1)                      │
│   ⏸️  Project Delta             [On Hold]    │
│                                             │
│   Complete Projects (2)                     │
│   ✅ Project Epsilon           [Complete]   │
│   🎉 Project Zeta              [Complete]   │
│                                             │
├─────────────────────────────────────────────┤
│ ↑/↓: Navigate  Enter: Select  /: Search    │
│ q: Quit  f: Filter  i: Set Icon             │
└─────────────────────────────────────────────┘
```

#### Project detail view (task groups)

```
┌─────────────────────────────────────────────┐
│ 🚀 Project Alpha                   [Active] │
├─────────────────────────────────────────────┤
│ [1] Tasks  [2] Overview  [3] Debug  [4] Notes│
├─────────────────────────────────────────────┤
│                                             │
│ Task Groups (5)              + New Group    │
│                                             │
│ > 🗄️  Database Migration       [12 tasks]   │
│   🔐 Authentication            [8 tasks]    │
│   🎨 UI Improvements           [15 tasks]   │
│   🔧 API Refactor              [6 tasks]    │
│   📝 Testing                    [20 tasks]   │
│                                             │
│                                             │
│                                             │
│                                             │
│                                             │
│                                             │
│                                             │
├─────────────────────────────────────────────┤
│ ↑/↓: Navigate  Enter: Open  n: New Group   │
│ d: Delete  r: Rename  i: Set Icon  Esc: Back│
└─────────────────────────────────────────────┘
```

#### Task group view (tasks list)

```
┌─────────────────────────────────────────────┐
│ 🚀 Project Alpha > 🗄️  Database Migration   │
├─────────────────────────────────────────────┤
│                                             │
│ Tasks (12)                    + New Task    │
│                                             │
│ TODO (4)                                    │
│ > Setup migration framework    [Feature]    │
│   • Create initial schema      [Feature]    │
│   • Add rollback support       [Feature]    │
│   • Fix migration bug          [Bug]        │
│                                             │
│ In Progress (3)                             │
│   • Write migration scripts    [Feature]    │
│   • Test migration locally     [Feature]    │
│   • Update docs                [Docs]       │
│                                             │
│ Done (5)                                    │
│   • Research tools             [Research]   │
│   • Setup dev environment      [Feature]    │
│   ...                                       │
│                                             │
├─────────────────────────────────────────────┤
│ ↑/↓: Navigate  Enter: Edit  Space: Toggle  │
│ n: New  d: Delete  Esc: Back to Groups      │
└─────────────────────────────────────────────┘
```

#### Task edit modal

```
┌─────────────────────────────────────────────┐
│                 Edit Task                   │
├─────────────────────────────────────────────┤
│                                             │
│ Text:                                       │
│ ┌─────────────────────────────────────────┐ │
│ │ Setup migration framework               │ │
│ └─────────────────────────────────────────┘ │
│                                             │
│ Status:  [•] TODO  [ ] In Progress  [ ] Done│
│                                             │
│ Type:    [•] Feature  [ ] Bug  [ ] Refactor │
│          [ ] Docs     [ ] Research          │
│                                             │
│                                             │
│         [Save]           [Cancel]           │
│                                             │
└─────────────────────────────────────────────┘
```

#### Task group create/rename modal

```
┌─────────────────────────────────────────────┐
│              New Task Group                 │
├─────────────────────────────────────────────┤
│                                             │
│ Icon (optional):                            │
│ ┌────┐                                      │
│ │ 🍾 │  (type emoji using Cmd+Ctrl+Space)  │
│ └────┘                                      │
│                                             │
│ Name:                                       │
│ ┌─────────────────────────────────────────┐ │
│ │ Performance Optimization                │ │
│ └─────────────────────────────────────────┘ │
│                                             │
│         [Create]         [Cancel]           │
│                                             │
└─────────────────────────────────────────────┘
```

#### Project icon modal

```
┌─────────────────────────────────────────────┐
│            Set Project Icon                 │
├─────────────────────────────────────────────┤
│                                             │
│ Icon (leave empty to remove):               │
│ ┌────┐                                      │
│ │ 🚀 │  (type emoji using Cmd+Ctrl+Space)  │
│ └────┘                                      │
│                                             │
│                                             │
│         [Save]           [Cancel]           │
│                                             │
└─────────────────────────────────────────────┘
```

### Navigation hierarchy

```
Projects List
    ↓ (select project)
Project Detail (tabs: Tasks, Overview, Debug, Notes)
    ↓ (tab 1: Tasks - shows task groups)
Task Groups List
    ↓ (select task group)
Task Group Detail (shows tasks)
    ↓ (select task)
Task Edit Modal
```

## Key Mappings

### Global

- `q` - quit application
- `?` - help/keybindings
- `/` - search/filter
- `Esc` - go back/cancel
- `Ctrl+C` - force quit

### Projects list

- `↑/↓` or `j/k` - navigate
- `Enter` - select project
- `i` - set/update icon for selected project
- `f` - filter by status
- `r` - refresh from Notion

### Project detail (tabs)

- `1-4` - switch tabs (tasks, overview, debug, notes)
- `r` - refresh

### Task groups list (tab 1)

- `↑/↓` or `j/k` - navigate
- `Enter` - open task group
- `n` - new task group
- `r` - rename selected group
- `i` - set/update icon for selected group
- `d` - delete selected group
- `Esc` - back to projects

### Task list (within group)

- `↑/↓` or `j/k` - navigate
- `Enter` - edit selected task
- `Space` - toggle status (quick action)
- `n` - new task
- `d` - delete selected task
- `t` - filter by type
- `s` - filter by status
- `c` - clear filters
- `Esc` - back to task groups

### Notes list (tabs 3-4)

- `↑/↓` or `j/k` - navigate
- `Enter` - edit note
- `n` - new note
- `d` - delete note
- `Esc` - back to project

## Configuration

### Config file location

`~/.config/notion-tui/config.yaml`

### Required settings

```yaml
notion:
  api_token: "secret_xxx"
  projects_database_id: "xxx"

ui:
  theme: "default"
  vim_mode: true

cache:
  enabled: true
  ttl: 300 # seconds
```

## Implementation Plan

### Setup

- initialize Go module
- setup project structure
- install dependencies (bubbletea, lipgloss, bubbles)
- setup Notion API client

### Notion integration

- implement API client wrapper
- create data models for projects, task groups, tasks, notes
- implement CRUD operations for all entities
- add support for reading/writing emoji icons (projects and task groups)
- add error handling and rate limiting

### TUI development

- create base Bubble Tea model
- implement projects list view
- implement project detail view (with tabs)
- implement task groups list component
- implement task list component (within group)
- implement forms/modals for editing
- add keyboard navigation
- implement filtering and search

### Polish

- add styling with lipgloss
- implement error messages and loading states
- add help screen
- write tests
- create documentation
- add configuration file support

### Deployment

- build instructions
- installation guide
- setup GitHub releases

## Potential Challenges

### API limitations

- Notion API rate limits -> implement caching and request throttling
- API response times -> show loading indicators, optimize queries
- nested database queries -> batch requests where possible
- handling task groups (pages vs blocks) -> determine best API approach

### UX considerations

- deep navigation hierarchy (projects > groups > tasks) -> clear breadcrumbs, easy back navigation
- long task/group lists -> pagination or virtualization
- real-time updates -> polling vs manual refresh

### Technical

- API token security -> use env vars or secure config
- error recovery -> graceful degradation, retry logic
- state management -> centralized state, proper updates
- managing nested state (project > group > task) -> clean data flow
- emoji rendering -> terminal compatibility, wide character handling for alignment

## Future Enhancements

- multi-workspace support
- custom views and dashboards
- task dependencies visualization
- time tracking integration
- export functionality (markdown, JSON)
- collaborative features (show who's editing)
- plugin system for custom extensions
- breadcrumb navigation in header
- task group templates
- drag and drop reordering (if feasible in TUI)
