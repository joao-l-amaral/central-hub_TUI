package components

import (
	"fmt"
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

func tabBorderWithBottom(left, right string) lipgloss.Border {
	border := lipgloss.HiddenBorder()

	border.MiddleLeft = left
	border.MiddleRight = right
	// border.Bottom = middle
	// border.BottomRight = right
	return border
}

func TabStyles(bgIsDark bool) *Styles {
	lightDark := lipgloss.LightDark(bgIsDark)

	// inactiveTabBorder := tabBorderWithBottom("─", "─")
	// activeTabBorder := tabBorderWithBottom("─", "─")
	highlightColor := lightDark(lipgloss.Color(style.ColorToHex(style.GetNeutralColor())), lipgloss.Color(style.ColorToHex(style.GetNeutralColor())))

	s := new(Styles)
	s.doc = lipgloss.NewStyle()
	s.inactiveTab = lipgloss.NewStyle().
		// Border(inactiveTabBorder, true).
		BorderForeground(highlightColor).
		Padding(0, 1)
	s.activeTab = s.inactiveTab.
		// Border(activeTabBorder, true).
		Foreground(lipgloss.Color(style.ColorToHex(style.GetPrimaryColor())))
	s.window = lipgloss.NewStyle().
		BorderForeground(highlightColor).
		Align(lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
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

		border.Left = "─"
		border.Right = "─"
		if isFirst && isActive {
			border.Left = "╭"
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.Left = "╭"
			border.BottomLeft = "│"
		} else if isLast && isActive {
			border.Right = "─"
		} else if isLast && !isActive {
			border.Right = "─"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	gapWidth := m.Width - lipgloss.Width(row)
	if gapWidth > 0 {
		filler := s.gap.Render(strings.Repeat("─", gapWidth-1) + "╮")
		row = lipgloss.JoinHorizontal(lipgloss.Center, row, filler)
	}
	doc.WriteString(row)
	doc.WriteString("\n")

	// Calculate remaining height for content area
	// Window style overhead: Padding(2,0)=4 rows + border bottom=1 row = 5
	tabRowHeight := lipgloss.Height(row) + 2 // +1 for the newline
	contentHeight := m.Height - tabRowHeight - 1

	doc.WriteString(
		s.window.
			Width(m.Width).
			Height(contentHeight).
			Align(lipgloss.Left).
			Render(m.TabContent[m.ActiveTab]),
	)

	return s.doc.Render(doc.String())
}

func GetTabPanelStyle(width int, height int) lipgloss.Style {
	halfWidth := int(float64(width) * 0.7)
	height = int(float64(height) * 0.93)

	return lipgloss.NewStyle().
		Width(halfWidth).
		Height(height)
}

// buildDetailContent renders the detail pane content for a git project.
func BuildInfoContent(p ProjectDTO) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Path: %s\n\nBranch: %s%s\n", p.Path, p.Branch, p.Changes))
	b.WriteString("\nEdited files:\n")
	if len(p.EditedFiles) == 0 {
		b.WriteString("  (none)")
	} else {
		for _, fc := range p.EditedFiles {
			var marker string
			var ms lipgloss.Style
			switch fc.Code {
			case "A", "?":
				marker, ms = "+", lipgloss.NewStyle().Foreground(style.GetSuccessColor())
			case "D":
				marker, ms = "-", lipgloss.NewStyle().Foreground(style.GetDangerColor())
			case "M":
				marker, ms = "~", lipgloss.NewStyle().Foreground(style.GetGoldenColor())
			default:
				marker, ms = " ", lipgloss.NewStyle().Foreground(style.GetNeutralColor())
			}
			b.WriteString("  ")
			b.WriteString(ms.Render(marker))
			b.WriteString(" ")
			b.WriteString(fc.Path)
			b.WriteString("\n")
		}
	}
	return b.String()
}
