package main

import (
	"log"

	"central_hub_tui/components"
	componentList "central_hub_tui/components/list"
	componentTabs "central_hub_tui/components/tabs"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
)

func main() {
	buildProjectList := componentList.BuildProjectList()

	tabs := []string{"Info", "Git History"}
	tabContent := []string{componentTabs.CentralHubSplash(), ""}

	m := model{
		projectListModel: buildProjectList,
		spinnerModel:     spinner.New(),
		tabModel: componentTabs.TabModel{
			Tabs:       tabs,
			TabContent: tabContent,
			Styles:     componentTabs.TabStyles(true),
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
