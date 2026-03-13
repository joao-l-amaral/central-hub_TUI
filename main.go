package main

import (
	"log"

	"central_hub_tui/components"
	"central_hub_tui/style"

	zone "github.com/lrstanley/bubblezone/v2"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

func main() {
	// Initialize the global bubblezone manager early so components can create zones.
	zone.NewGlobal()
	defer zone.Close()
	// Load from json file
	projects := components.LoadProjectPaths()

	items := make([]list.Item, len(projects))
	for i, p := range projects {
		items[i] = components.AddProjectToList(p)
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(style.GetPrimaryColor()).
		BorderForeground(style.GetPrimaryColor())
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(style.GetNeutralColor()).
		BorderForeground(style.GetPrimaryColor())

	tabs := []string{"Info", "Git History", "Worktrees"}
	tabContent := []string{"Info Tab", "Git History Tab", "Worktrees"}

	m := model{
		listModel: components.ListModel{
			List: list.New(items, delegate, 0, 0),
		},
		tabModel: components.TabModel{
			Tabs:       tabs,
			TabContent: tabContent,
			Styles:     components.TabStyles(true),
			ActiveTab:  0,
		},
		footerModel: components.FooterModel{
			ProjectName: "Central Hub",
			Height:      3,
		},
	}

	zone.NewGlobal()

	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		log.Fatalf("unable to run tui: %v", err)
	}
}
