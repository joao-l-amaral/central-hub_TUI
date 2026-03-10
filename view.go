package main

import (
	"central_hub_tui/components"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func setTabPanelStyle(width int, height int) lipgloss.Style {
	halfWidth := int(float64(width) * 0.7)
	height = int(float64(height) * 0.93)

	return lipgloss.NewStyle().
		Width(halfWidth).
		Height(height)
}

func (m model) View() tea.View {

	projectListStyle := components.GetListPanelStyle(m.width, m.height)
	tabStyle := setTabPanelStyle(m.width, m.height)
	footerStyle := components.GetFooterStyle(m.width, m.footerModel.Height)

	projectListComponentView := m.listModel.List.View()
	projectListPanel := zone.Scan(projectListStyle.Render(projectListComponentView))

	tabView := components.TabView(m.tabModel)
	tabPanel := tabStyle.Render(tabView)

	footView := components.GetFooterContent(m.footerModel)
	footer := footerStyle.Render(footView)

	row1 := lipgloss.JoinHorizontal(lipgloss.Top, projectListPanel, tabPanel)

	content := lipgloss.JoinVertical(lipgloss.Left, row1, footer)
	view := tea.NewView(content)
	view.AltScreen = true
	view.WindowTitle = "Central Hub"
	view.MouseMode = tea.MouseModeAllMotion

	return view
}
