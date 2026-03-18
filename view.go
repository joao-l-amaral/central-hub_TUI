package main

import (
	"central_hub_tui/components"
	"central_hub_tui/style"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func (m model) View() tea.View {

	halfWidth := int(float64(m.width)*0.3) - 2

	tabStyle := components.GetTabPanelStyle(m.width+10, m.height).BorderForeground(lipgloss.Color(style.ColorToHex(style.GetPrimaryColor())))
	footerStyle := components.GetFooterStyle(m.width, m.footerModel.Height)

	switch m.focused {
	case FocusProject:
		m.projectListModel.IsSelected = true
		m.projectWorktreeList.IsSelected = false
	case FocusWorktree:
		m.projectListModel.IsSelected = false
		m.projectWorktreeList.IsSelected = true
	}

	projectListComponentView := m.projectListModel.List.View()
	projectListPanel := components.RoundedTitleBox("Projects", projectListComponentView, halfWidth, m.height/2, m.projectListModel.IsSelected, false, 0)

	worktreeView := m.projectWorktreeList.List.View()
	worktreePanel := components.RoundedTitleBox("Worktrees", worktreeView, halfWidth, m.height/2, m.projectWorktreeList.IsSelected, true, m.projectWorktreeList.NumberOfWorktrees)

	// tab panel is always "focused" — left/right always change tabs
	m.tabModel.Focused = true
	tabView := components.TabView(m.tabModel)
	tabPanel := tabStyle.Render(tabView)

	footView := components.GetFooterContent(m.footerModel)
	footer := footerStyle.Render(footView)

	col1 := lipgloss.JoinVertical(lipgloss.Left, projectListPanel, worktreePanel)
	row1 := lipgloss.JoinHorizontal(lipgloss.Top, col1, tabPanel)

	content := lipgloss.JoinVertical(lipgloss.Left, row1, footer)
	view := tea.NewView(content)
	view.AltScreen = true
	view.WindowTitle = "Central Hub"
	view.MouseMode = tea.MouseModeAllMotion

	return view
}
