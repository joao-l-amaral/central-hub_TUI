package main

import (
	"central_hub_tui/components"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func setListPanelStyle(width int, height int) lipgloss.Style {
	halfWidth := int(float64(width)*0.3) - 2
	height = int(float64(height) * 0.93)

	return lipgloss.NewStyle().
		Padding(0, 1).
		Width(halfWidth).
		Height(height)
}

func setTabPanelStyle(width int, height int) lipgloss.Style {
	halfWidth := int(float64(width) * 0.7)
	height = int(float64(height) * 0.93)

	return lipgloss.NewStyle().
		Width(halfWidth).
		Height(height)
}

func (m model) View() tea.View {

	projectListStyle := setListPanelStyle(m.width, m.height)
	tabStyle := setTabPanelStyle(m.width, m.height)
	footerStyle := components.GetFooterStyle(m.width, m.footerModel.Height)

	projectListComponentView := m.listModel.List.View()
	projectListPanel := projectListStyle.Render(projectListComponentView)

	tabView := components.TabView(m.tabModel)
	tabPanel := tabStyle.Render(tabView)

	footView := components.GetFooterContent(m.footerModel)
	footer := footerStyle.Render(footView)

	row1 := lipgloss.JoinHorizontal(lipgloss.Top, projectListPanel, tabPanel)

	content := lipgloss.JoinVertical(lipgloss.Left, row1, footer)
	v := tea.NewView(content)
	v.AltScreen = true
	v.WindowTitle = "Central Hub"

	//Actions
	selectProjectInList(m)

	return v
}

func selectProjectInList(m model) {
	if selectedItem, ok := m.listModel.List.SelectedItem().(components.ProjectEntry); ok {
		if selectedItem.IsGit {
			m.tabModel.TabContent[0] = "Project Info - " + selectedItem.Name
			m.tabModel.TabContent[1] = "Git History Tab - " + selectedItem.Name
		}
	} else {
		m.tabModel.TabContent[0] = "Project Info"
		m.tabModel.TabContent[1] = "Git History Tab"
	}
}
