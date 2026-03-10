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

// Project model
type model struct {
	tabModel    components.TabModel
	listModel   components.ListModel
	footerModel components.FooterModel
	width       int
	height      int
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
			if m.tabModel.ActiveTab+1 > len(m.tabModel.Tabs)-1 {
				m.tabModel.ActiveTab = 0
			} else {
				m.tabModel.ActiveTab = min(m.tabModel.ActiveTab+1, len(m.tabModel.Tabs)-1)
			}

			return m, nil
		case "shift+tab":
			if m.tabModel.ActiveTab-1 < 0 {
				m.tabModel.ActiveTab = len(m.tabModel.Tabs) - 1
			} else {
				m.tabModel.ActiveTab = max(m.tabModel.ActiveTab-1, 0)
			}

			return m, nil
		}
	case tea.MouseClickMsg:
		if msg.Button == tea.MouseLeft {
			// Click at msg.X, msg.Y; check modifiers: msg.Mod
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
		m.listModel.List.SetSize(listWidth, listHeight)

	}

	var cmd tea.Cmd
	m.listModel.List, cmd = m.listModel.List.Update(msg)
	m.listModel.List.Help.ShowAll = false
	m.listModel.List.SetShowStatusBar(false)
	m.listModel.List.SetShowHelp(false)
	m.listModel.List.SetShowTitle(false)

	return m, cmd
}
