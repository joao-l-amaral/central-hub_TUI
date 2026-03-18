package main

import (
	"central_hub_tui/components"
	"central_hub_tui/style"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
)

const (
	listView uint = iota
	titleView
	bodyView
)

const (
	FocusProject uint = iota
	FocusWorktree
	FocusTab
)

// projectDataMsg carries the results of async git data loading for a project.
type projectDataMsg struct {
	id        string
	info      string
	history   string
	worktrees []components.ProjectDTO
}

// loadProjectDataCmd runs all git calls for a project in a goroutine.
func loadProjectDataCmd(worktreeModel components.ProjectWorktreeListModel, project components.ProjectDTO) tea.Cmd {
	return func() tea.Msg {
		entry := components.ProjectEntry{
			Name:    project.Name,
			Path:    project.Path,
			IsGit:   project.IsGit,
			Options: project.Options,
		}
		return projectDataMsg{
			id:        project.ID,
			info:      components.BuildInfoContent(project),
			history:   components.BuildHistoryContent(project),
			worktrees: components.GetGitWorktrees(entry),
		}
	}
}

// Project model
type model struct {
	tabModel            components.TabModel
	projectListModel    components.ProjectListModel
	projectWorktreeList components.ProjectWorktreeListModel
	footerModel         components.FooterModel
	spinnerModel        spinner.Model
	width               int
	height              int
	focused             uint
	lastSelectedID      string
}

func (m model) Init() tea.Cmd {
	return m.spinnerModel.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.focused = m.focused + 1
			if m.focused > 1 {
				m.focused = 0
			}
			return m, nil
		case "left", "h":
			if m.tabModel.ActiveTab-1 < 0 {
				m.tabModel.ActiveTab = len(m.tabModel.Tabs) - 1
			} else {
				m.tabModel.ActiveTab = max(m.tabModel.ActiveTab-1, 0)
			}
			return m, nil
		case "right", "l":
			if m.tabModel.ActiveTab+1 > len(m.tabModel.Tabs)-1 {
				m.tabModel.ActiveTab = 0
			} else {
				m.tabModel.ActiveTab = min(m.tabModel.ActiveTab+1, len(m.tabModel.Tabs)-1)
			}
			return m, nil
		}
	case projectDataMsg:
		// Only apply if the user hasn't moved to a different project already.
		if msg.id == m.lastSelectedID {
			m.tabModel.TabContent[0] = msg.info
			m.tabModel.TabContent[1] = msg.history

			delegate := list.NewDefaultDelegate()
			delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
				Foreground(style.GetPrimaryColor()).
				BorderForeground(style.GetPrimaryColor())
			delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
				Foreground(style.GetNeutralColor()).
				BorderForeground(style.GetPrimaryColor())

			items := make([]list.Item, len(msg.worktrees))
			for i, wt := range msg.worktrees {
				items[i] = wt
			}

			wtWidth := int(float64(m.width)*0.3) - 4
			wtHeight := int(float64(m.height)*0.3) - 4
			m.projectWorktreeList = components.ProjectWorktreeListModel{
				List:              list.New(items, delegate, wtWidth, wtHeight),
				NumberOfWorktrees: len(msg.worktrees),
			}

			// Disable worktrees list help/status bars if there are items to avoid UI clutter.
			if len(m.projectWorktreeList.List.Items()) > 0 {
				m.projectWorktreeList.List.Help.ShowAll = false
				m.projectWorktreeList.List.SetShowStatusBar(false)
				m.projectWorktreeList.List.SetShowHelp(false)
				m.projectWorktreeList.List.SetShowTitle(false)
			}
		}
		return m, nil
	case spinner.TickMsg:
		var spinCmd tea.Cmd
		m.spinnerModel, spinCmd = m.spinnerModel.Update(msg)
		return m, spinCmd
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		panelHeight := int(float64(msg.Height) * 0.93)
		panelWidth := int(float64(msg.Width) * 0.7)
		m.tabModel.Height = panelHeight
		m.tabModel.Width = panelWidth - 2 // -2 for window's left+right border
		listHeight := components.SetInnerListHeight(msg)
		// Calculate the list width (30% of window width - 2, minus padding)
		listWidth := int(float64(msg.Width)*0.3) - 4
		m.projectListModel.List.SetSize(listWidth, listHeight)

		if m.width < 20 {
			m.width = 20
		}
	}

	var cmd tea.Cmd
	switch m.focused {
	case FocusProject:
		m.projectListModel.List, cmd = m.projectListModel.List.Update(msg)
	case FocusWorktree:
		if len(m.projectWorktreeList.List.Items()) > 0 {
			m.projectWorktreeList.List, cmd = m.projectWorktreeList.List.Update(msg)
		} else {
			cmd = nil
		}
	default:
		m.projectListModel.List, cmd = m.projectListModel.List.Update(msg)
	}

	// Disable project list help/status bars if there are items to avoid UI clutter.
	if len(m.projectListModel.List.Items()) > 0 {
		m.projectListModel.List.Help.ShowAll = false
		m.projectListModel.List.SetShowStatusBar(false)
		m.projectListModel.List.SetShowHelp(false)
		m.projectListModel.List.SetShowTitle(false)
	}

	// Detect selection change and fire async git load if needed.
	var selCmd tea.Cmd
	m, selCmd = selectProjectInList(m)

	return m, tea.Batch(cmd, selCmd)
}

// selectProjectInList detects a selection change and returns an async cmd to load git data.
func selectProjectInList(m model) (model, tea.Cmd) {
	if selectedItem, ok := m.projectListModel.List.SelectedItem().(components.ProjectDTO); ok {
		if selectedItem.ID == m.lastSelectedID {
			return m, nil
		}
		m.lastSelectedID = selectedItem.ID

		if selectedItem.IsGit {
			m.tabModel.TabContent[0] = "Loading..."
			m.tabModel.TabContent[1] = "Loading..."
			m.projectWorktreeList.Loading = true
			return m, loadProjectDataCmd(m.projectWorktreeList, selectedItem)
		}

		m.tabModel.TabContent[0] = "Project Info"
		m.tabModel.TabContent[1] = "Git History Tab"
		m.projectWorktreeList.Loading = false
	}
	return m, nil
}
