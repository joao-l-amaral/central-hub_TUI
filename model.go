package main

import (
	"central_hub_tui/components"
	"central_hub_tui/components/list"
	"central_hub_tui/components/types"
	"central_hub_tui/utils/git"
	"os/exec"

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

// loadProjectDataCmd runs all git calls for a project in a goroutine.
func loadProjectDataCmd(project types.ProjectDTO) tea.Cmd {
	return func() tea.Msg {
		entry := types.ProjectEntry{
			Name:    project.Name,
			Path:    project.Path,
			IsGit:   project.IsGit,
			Options: project.Options,
		}
		return types.ProjectWorktreeDataMsg{
			Id:        project.ID,
			Info:      "",
			History:   "",
			Worktrees: git.GetGitWorktrees(entry),
		}
	}
}

// Project model
type model struct {
	tabModel            components.TabModel
	projectListModel    types.ProjectListModel
	projectWorktreeList types.ProjectWorktreeListModel
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
		case "enter":
			switch m.focused {
			case FocusProject:
				m.focused = 1
			case FocusWorktree:
				//TODO open the terminal in the path of the selected worktree
				if selected, ok := m.projectWorktreeList.List.SelectedItem().(types.WorktreeItem); ok {
					// if err := os.Chdir(selected.Path); err != nil {
					// 	m.lastAction = "Failed to chdir: " + err.Error()
					// 	return m, nil
					// }

					cmd := exec.Command(
						"pwsh",
						"-NoProfile",
						"-ExecutionPolicy", "Bypass",
						"-File", "scripts/open-wt-tab.ps1",
						"-TargetPath", selected.Path,
					)
					if err := cmd.Start(); err != nil {
						// m.lastAction = "Failed to open terminal: " + err.Error()
						return m, nil
					}

					// m.lastAction = "Navigated to: " + selected.Path
				}
			}
		}
	case types.ProjectWorktreeDataMsg:
		// Only apply if the user hasn't moved to a different project already.
		if msg.Id == m.lastSelectedID {
			m.tabModel.TabContent[0] = msg.Info
			m.tabModel.TabContent[1] = msg.History

			m.projectWorktreeList = list.BuildWorktreeList(m.width, m.height, msg)
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
		listHeight := list.SetInnerListHeight(msg)
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

	list.ConfigureListOptions(&m.projectListModel.List)

	// Detect selection change and fire async load — only for the focused panel.
	var selCmd tea.Cmd
	var selWtCmd tea.Cmd
	switch m.focused {
	case FocusProject:
		m, selCmd = selectProjectInList(m)
	case FocusWorktree:
		m, selWtCmd = selectWorktreeInList(m)
	}

	return m, tea.Batch(cmd, selCmd, selWtCmd)
}

// selectProjectInList detects a selection change and returns an async cmd to load git data.
func selectProjectInList(m model) (model, tea.Cmd) {
	if selectedItem, ok := m.projectListModel.List.SelectedItem().(types.ProjectDTO); ok {
		if selectedItem.ID == m.lastSelectedID {
			return m, nil
		}
		m.lastSelectedID = selectedItem.ID

		if selectedItem.IsGit {
			components.SetTabContent(&m.tabModel, 0, m.spinnerModel.View()+" Loading project info...")
			components.SetTabContent(&m.tabModel, 1, m.spinnerModel.View()+" Loading git history...")
			m.projectWorktreeList.Loading = true
			return m, loadProjectDataCmd(selectedItem)
		}

		components.SetTabContent(&m.tabModel, 0, "Project Info Tab")
		components.SetTabContent(&m.tabModel, 1, "Git History Tab")
		m.projectWorktreeList.Loading = false
	}
	return m, nil
}

func selectWorktreeInList(m model) (model, tea.Cmd) {
	if selectedItem, ok := m.projectWorktreeList.List.SelectedItem().(types.WorktreeItem); ok {
		if selectedItem.ID == m.lastSelectedID {
			return m, nil
		}
		m.lastSelectedID = selectedItem.ID

		components.SetTabContent(&m.tabModel, 0, components.BuildInfoContent(selectedItem))
		components.SetTabContent(&m.tabModel, 1, git.BuildHistoryContent(selectedItem.Path))
	}
	return m, nil
}
