package components

import (
	"encoding/json"
	"os"
	"path/filepath"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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

var projectPaths []ProjectEntry

func AddProjectToList(p ProjectEntry) ProjectEntry {
	return ProjectEntry{Name: p.Name, Path: p.Path, IsGit: p.IsGit, Options: p.Options}
}

func (i ProjectEntry) Title() string       { return i.Name }
func (i ProjectEntry) Description() string { return i.Path }
func (i ProjectEntry) FilterValue() string { return i.Name }

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
