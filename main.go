package main

import (
	"log"

	"central_hub_tui/components"
	componentList "central_hub_tui/components/list"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
)

func main() {
	buildProjectList := componentList.BuildProjectList()

	tabs := []string{"Info", "Git History"}
	tabContent := []string{"Info Tab", "Git History Tab"}

	m := model{
		projectListModel: buildProjectList,
		spinnerModel:     spinner.New(),
		tabModel: components.TabModel{
			Tabs:       tabs,
			TabContent: tabContent,
			Styles:     components.TabStyles(true),
			ActiveTab:  0,
			Focused:    false,
		},
		footerModel: components.FooterModel{
			ProjectName: "Central Hub",
			Height:      3,
		},
	}

	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		log.Fatalf("unable to run tui: %v", err)
	}
}
