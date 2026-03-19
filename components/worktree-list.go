package components

import (
	"charm.land/bubbles/v2/list"
	"charm.land/lipgloss/v2"
)

type ProjectWorktreeListModel struct {
	List              list.Model
	window            lipgloss.Style
	gap               lipgloss.Style
	IsSelected        bool
	NumberOfWorktrees int
	Loading           bool
}

type WorktreeItem ProjectDTO

func (i WorktreeItem) Title() string       { return i.Branch }
func (i WorktreeItem) Description() string { return i.Path }
func (i WorktreeItem) FilterValue() string { return i.Branch }

func AddWorktreeToList(p ProjectEntry) WorktreeItem {
	return WorktreeItem{
		Name:        p.Name,
		ID:          p.Name,
		Path:        p.Path,
		Branch:      LoadProjectGitInfo(p),
		EditedFiles: LoadChangedFiles(p),
		IsGit:       p.IsGit,
		Options:     p.Options,
	}
}
