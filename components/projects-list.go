package components

import (
	"encoding/json"
	"os"
	"path/filepath"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type ListModel struct {
	List list.Model
}

type ProjectEntry struct {
	Name    string            `json:"name"`
	Path    string            `json:"path"`
	IsGit   bool              `json:"isGit"`
	Options map[string]string `json:"options"`
}

type ProjectDTO struct {
	Name    string
	ID      string
	Path    string
	Branch  string
	Changes string
	IsGit   bool
	Options map[string]string
}

var projectPaths []ProjectEntry

func AddProjectToList(p ProjectEntry) ProjectDTO {
	return ProjectDTO{
		Name:    p.Name,
		ID:      p.Name,
		Path:    p.Path,
		Branch:  LoadProjectGitInfo(p),
		IsGit:   p.IsGit,
		Options: p.Options,
	}
}
func (i ProjectDTO) Title() string       { return zone.Mark(i.ID, i.Name) }
func (i ProjectDTO) Description() string { return i.Branch }
func (i ProjectDTO) FilterValue() string { return i.Name }

func SetInnerListHeight(msg tea.WindowSizeMsg) int {
	panelHeight := int(float64(msg.Height) * 0.95)
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

func GetListPanelStyle(width int, height int) lipgloss.Style {
	halfWidth := int(float64(width)*0.3) - 2
	height = int(float64(height) * 0.93)

	return lipgloss.NewStyle().
		Padding(0, 1).
		Width(halfWidth).
		Height(height)
}
