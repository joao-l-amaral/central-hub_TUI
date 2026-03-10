package components

import (
	"strings"

	"central_hub_tui/style"

	"charm.land/lipgloss/v2"
)

type TabModel struct {
	Tabs       []string
	TabContent []string
	Styles     *Styles
	ActiveTab  int
	Height     int
	Width      int
}

type Styles struct {
	doc         lipgloss.Style
	highlight   lipgloss.Style
	inactiveTab lipgloss.Style
	activeTab   lipgloss.Style
	window      lipgloss.Style
	gap         lipgloss.Style
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

func TabStyles(bgIsDark bool) *Styles {
	lightDark := lipgloss.LightDark(bgIsDark)

	inactiveTabBorder := tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder := tabBorderWithBottom("┘", " ", "└")
	highlightColor := lightDark(lipgloss.Color(style.ColorToHex(style.GetPrimaryColor())), lipgloss.Color(style.ColorToHex(style.GetPrimaryColor())))

	s := new(Styles)
	s.doc = lipgloss.NewStyle()
	s.inactiveTab = lipgloss.NewStyle().
		Border(inactiveTabBorder, true).
		BorderForeground(highlightColor).
		Padding(0, 1)
	s.activeTab = s.inactiveTab.
		Border(activeTabBorder, true)
	s.window = lipgloss.NewStyle().
		BorderForeground(highlightColor).
		Padding(2, 0).
		Align(lipgloss.Center).
		Border(lipgloss.NormalBorder()).
		UnsetBorderTop()
	s.gap = lipgloss.NewStyle().Foreground(highlightColor)
	return s
}

func TabView(m TabModel) string {
	if m.Styles == nil {
		return ""
	}

	doc := strings.Builder{}
	s := m.Styles

	var renderedTabs []string

	for i, t := range m.Tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.Tabs)-1, i == m.ActiveTab
		if isActive {
			style = s.activeTab
		} else {
			style = s.inactiveTab
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "└"
		} else if isLast && isActive {
			border.BottomRight = "┘"
		} else if isLast && !isActive {
			border.BottomRight = "┘"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	gapWidth := m.Width - lipgloss.Width(row)
	if gapWidth > 0 {
		filler := s.gap.Render(strings.Repeat("─", gapWidth-1) + "┐")
		row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, filler)
	}
	doc.WriteString(row)
	doc.WriteString("\n")

	// Calculate remaining height for content area
	// Window style overhead: Padding(2,0)=4 rows + border bottom=1 row = 5
	tabRowHeight := lipgloss.Height(row) + 2 // +1 for the newline
	contentHeight := m.Height - tabRowHeight - 1

	doc.WriteString(
		s.window.
			Width(m.Width).        // -2 for left+right border
			Height(contentHeight). // fills remaining terminal height
			Render(m.TabContent[m.ActiveTab]),
	)

	return s.doc.Render(doc.String())
}
