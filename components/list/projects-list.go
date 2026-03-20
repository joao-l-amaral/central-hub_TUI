package list

import (
	"encoding/json"
	"os"
	"path/filepath"

	"central_hub_tui/components/types"
	"central_hub_tui/style"

	"charm.land/bubbles/v2/list"
	"charm.land/lipgloss/v2"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

var projectPaths []types.ProjectEntry

func AddProjectToList(p types.ProjectEntry) types.ProjectDTO {
	return types.ProjectDTO{
		Name:    p.Name,
		ID:      p.Name,
		Path:    p.Path,
		IsGit:   p.IsGit,
		Options: p.Options,
	}
}

func init() {
	projectPaths = LoadProjectPaths()
}

func LoadProjectPaths() []types.ProjectEntry {
	candidates := []string{"data/projects.json"}
	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), "projects.json"))
	}
	for _, cfgPath := range candidates {
		data, err := os.ReadFile(cfgPath)
		if err != nil {
			continue
		}
		var entries []types.ProjectEntry
		if err := json.Unmarshal(data, &entries); err != nil {
			continue
		}
		return entries
	}
	return nil
}

// Load from json file
func BuildProjectList() types.ProjectListModel {

	projects := LoadProjectPaths()

	items := make([]list.Item, len(projects))
	for i, p := range projects {
		items[i] = AddProjectToList(p)
	}

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(style.GetPrimaryColor()).
		BorderForeground(style.GetPrimaryColor())

	return types.ProjectListModel{
		List: list.New(items, delegate, 0, 0),
	}
}
