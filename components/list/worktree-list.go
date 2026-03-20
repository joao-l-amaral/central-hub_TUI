package list

import (
	"central_hub_tui/components/types"
	"central_hub_tui/style"

	"charm.land/bubbles/v2/list"
)

func BuildWorktreeList(width, height int, msg types.ProjectWorktreeDataMsg) types.ProjectWorktreeListModel {
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(style.GetPrimaryColor()).
		BorderForeground(style.GetPrimaryColor())
	delegate.ShowDescription = false

	items := make([]list.Item, len(msg.Worktrees))
	for i, wt := range msg.Worktrees {
		items[i] = wt
	}

	wtWidth := int(float64(width) * 0.3)
	wtHeight := int(float64(height) * 0.3)
	worktreeList := types.ProjectWorktreeListModel{
		List:              list.New(items, delegate, wtWidth, wtHeight),
		NumberOfWorktrees: len(msg.Worktrees),
	}

	ConfigureListOptions(&worktreeList.List)

	return worktreeList
}
