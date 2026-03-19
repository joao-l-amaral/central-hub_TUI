package components

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"central_hub_tui/style"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type ProjectListModel struct {
	List       list.Model
	window     lipgloss.Style
	gap        lipgloss.Style
	IsSelected bool
}

type ProjectEntry struct {
	Name    string            `json:"name"`
	Path    string            `json:"path"`
	IsGit   bool              `json:"isGit"`
	Options map[string]string `json:"options"`
}

type FileChange struct {
	Path string
	Code string // A=added, M=modified, D=deleted, ?=untracked
}

type ProjectDTO struct {
	Name        string
	ID          string
	Path        string
	Branch      string
	Changes     string
	IsGit       bool
	EditedFiles []FileChange
	Options     map[string]string
}

var projectPaths []ProjectEntry

func AddProjectToList(p ProjectEntry) ProjectDTO {
	return ProjectDTO{
		Name:    p.Name,
		ID:      p.Name,
		Path:    p.Path,
		IsGit:   p.IsGit,
		Options: p.Options,
	}
}

func (i ProjectDTO) Title() string       { return i.Name }
func (i ProjectDTO) Description() string { return i.Path }
func (i ProjectDTO) FilterValue() string { return i.Name }

func SetInnerListHeight(msg tea.WindowSizeMsg) int {
	panelHeight := int(float64(msg.Height) * 0.5)
	return panelHeight - 4
}

func (p ProjectEntry) GetProjectEntry() ProjectEntry {
	return p
}

func init() {
	projectPaths = LoadProjectPaths()
}

func LoadProjectPaths() []ProjectEntry {
	candidates := []string{"data/projects.json"}
	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), "projects.json"))
	}
	for _, cfgPath := range candidates {
		data, err := os.ReadFile(cfgPath)
		if err != nil {
			continue
		}
		var entries []ProjectEntry
		if err := json.Unmarshal(data, &entries); err != nil {
			continue
		}
		return entries
	}
	return nil
}

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
