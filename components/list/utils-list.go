package list

import (
	"fmt"
	"strings"

	"central_hub_tui/style"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func RoundedTitleBox(title, content string, width int, height int, isSelected bool, showCounter bool, counter int) string {
	titleLen := len(title) // ╭ title ╮
	fillerLen := max(0, width-titleLen-8)

	focusColor := lipgloss.Color(style.ColorToHex(style.GetNeutralColor()))

	if isSelected {
		focusColor = lipgloss.Color(style.ColorToHex(style.GetPrimaryColor()))
	}

	// Title row: ╭─ title ─╮ (manual top)
	// Box-drawing chars stay in the neutral color; only the title text uses focusColor.
	borderColor := lipgloss.Color(style.ColorToHex(style.GetNeutralColor()))
	borderStyle := lipgloss.NewStyle().Foreground(borderColor)
	styledTitle := lipgloss.NewStyle().Bold(true).Foreground(focusColor).Render(title)

	titlePrefix := "╭──── "
	if showCounter {
		titlePrefix = fmt.Sprintf("╭─[%d] ", counter)
	}

	titleRow := borderStyle.Render(titlePrefix) +
		styledTitle +
		borderStyle.Render(fmt.Sprintf(" %s╮", strings.Repeat("─", fillerLen)))

	// Content box: no top border, rounded sides/bottom
	boxStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderTop(false).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true).
		BorderForeground(lipgloss.Color("#FAF4E9")).
		Padding(1).
		Width(width).
		Height(height - 4) // Minus title row

	// Join
	return lipgloss.JoinVertical(lipgloss.Top, titleRow, boxStyle.Render(content))
}

func SetInnerListHeight(msg tea.WindowSizeMsg) int {
	panelHeight := int(float64(msg.Height) * 0.5)
	return panelHeight - 4
}

// Disable project list help/status bars if there are items to avoid UI clutter.
func ConfigureListOptions(list *list.Model) {
	if len(list.Items()) > 0 {
		list.SetShowTitle(false)
		list.Help.ShowAll = false
		list.SetShowStatusBar(false)
		list.SetShowHelp(false)
	}
}
