package types

import (
	"charm.land/bubbles/v2/list"
	"charm.land/lipgloss/v2"
)

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

func (i ProjectDTO) Title() string       { return i.Name }
func (i ProjectDTO) Description() string { return i.Path }
func (i ProjectDTO) FilterValue() string { return i.Name }

type ProjectWorktreeListModel struct {
	List              list.Model
	window            lipgloss.Style
	gap               lipgloss.Style
	IsSelected        bool
	NumberOfWorktrees int
	Loading           bool
}

// projectDataMsg carries the results of async git data loading for a project.
type ProjectWorktreeDataMsg struct {
	Id        string
	Info      string
	History   string
	Worktrees []WorktreeItem
}

type WorktreeItem ProjectDTO

func (i WorktreeItem) Title() string       { return i.Branch }
func (i WorktreeItem) Description() string { return i.Path }
func (i WorktreeItem) FilterValue() string { return i.Branch }
