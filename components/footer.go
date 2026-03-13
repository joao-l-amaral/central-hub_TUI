package components

import (
	"fmt"
	"time"

	"central_hub_tui/style"

	"charm.land/lipgloss/v2"
)

type FooterModel struct {
	Height      int
	ProjectName string
}

func GetFooterStyle(windowWidth int, footerHeight int) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderBottom(false).
		BorderRight(false).
		BorderLeft(false).
		Padding(1, 1).
		Width(windowWidth).
		Align(lipgloss.Center).
		Height(footerHeight)
}

func GetFooterContent(m FooterModel) string {
	projectName := lipgloss.NewStyle().Foreground(style.GetPrimaryColor()).Render(m.ProjectName)
	dataInfo := lipgloss.NewStyle().Foreground(style.GetNeutralColor()).Render(fmt.Sprintf(" TUI - @ %d ", time.Now().Year()))
	authorName := lipgloss.NewStyle().Foreground(style.GetPrimaryColor()).Render("João Amaral")
	leftStr := projectName + dataInfo + authorName
	return lipgloss.NewStyle().Height(1).Render(leftStr)
}
