package main

import (
	"central_hub_tui/components"

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

// Project model
type model struct {
	tabModel            components.TabModel
	projectListModel    components.ProjectListModel
	projectWorktreeList components.ProjectWorktreeListModel
	footerModel         components.FooterModel
	width               int
	height              int
	focused             uint
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			if m.focused+1 > FocusTab {
				m.focused = FocusProject
			} else {
				m.focused = m.focused + 1
			}
			return m, nil
		case "shift+tab":
			if m.focused == FocusProject {
				m.focused = FocusTab
			} else {
				m.focused = m.focused - 1
			}
			return m, nil
		case "left", "h":
			if m.focused == FocusTab {
				if m.tabModel.ActiveTab-1 < 0 {
					m.tabModel.ActiveTab = len(m.tabModel.Tabs) - 1
				} else {
					m.tabModel.ActiveTab = max(m.tabModel.ActiveTab-1, 0)
				}
				return m, nil
			}
		case "right", "l":
			if m.focused == FocusTab {
				if m.tabModel.ActiveTab+1 > len(m.tabModel.Tabs)-1 {
					m.tabModel.ActiveTab = 0
				} else {
					m.tabModel.ActiveTab = min(m.tabModel.ActiveTab+1, len(m.tabModel.Tabs)-1)
				}
				return m, nil
			}
		}
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
	case FocusTab:
		// Tab view is static (no interactive bubble component here). Handle
		// key events for tabs above (left/right) so nothing to update here.
		m.tabModel.Focused = true
		cmd = nil
	default:
		m.projectListModel.List, cmd = m.projectListModel.List.Update(msg)
	}

	// Configure common list display options for lists that are initialized
	if len(m.projectListModel.List.Items()) > 0 {
		m.projectListModel.List.Help.ShowAll = false
		m.projectListModel.List.SetShowStatusBar(false)
		m.projectListModel.List.SetShowHelp(false)
		m.projectListModel.List.SetShowTitle(false)
	}
	if len(m.projectWorktreeList.List.Items()) > 0 {
		m.projectWorktreeList.List.Help.ShowAll = false
		m.projectWorktreeList.List.SetShowStatusBar(false)
		m.projectWorktreeList.List.SetShowHelp(false)
		m.projectWorktreeList.List.SetShowTitle(false)
	}

	// Update tab content based on current list selection so view shows immediately.
	m = selectProjectInList(m)
	//TODO load worktrees when project selection changes and update tab
	// content if worktree tab is active. using enter and then focus on worktree panel for now
	// to trigger loading worktrees and updating content.

	return m, cmd
}

// selectProjectInList updates tab content based on the currently selected list item.
func selectProjectInList(m model) model {
	if selectedItem, ok := m.projectListModel.List.SelectedItem().(components.ProjectDTO); ok {
		if selectedItem.IsGit {
			m.tabModel.TabContent[0] = components.BuildInfoContent(selectedItem)
			// Only build the potentially expensive git history when the Git History
			// tab is active. Building it every update caused frequent `git` calls.
			if m.tabModel.ActiveTab == 1 {
				m.tabModel.TabContent[1] = components.BuildHistoryContent(selectedItem)
			}
		} else {
			m.tabModel.TabContent[0] = "Project Info"
			m.tabModel.TabContent[1] = "Git History Tab"
		}
	}
	return m
}
