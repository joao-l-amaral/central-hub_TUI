package components

import (
	"bytes"
	"fmt"
	"image/color"
	"os/exec"
	"regexp"
	"strings"

	"central_hub_tui/style"

	"charm.land/lipgloss/v2"
)

// gitCmd runs a git command in dir and returns trimmed stdout.
func gitCmd(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

// countChanges parses git porcelain output and counts added/modified/deleted files.
func countChanges(porcelain string) (added, modified, deleted int) {
	lines := strings.Split(strings.TrimSpace(porcelain), "\n")
	for _, l := range lines {
		if l == "" {
			continue
		}
		status := l
		if len(l) >= 2 {
			status = l[0:2]
		}
		if strings.Contains(status, "?") {
			added++
			continue
		}
		if strings.Contains(status, "A") {
			added++
		}
		if strings.Contains(status, "M") {
			modified++
		}
		if strings.Contains(status, "D") {
			deleted++
		}
	}
	return
}

// colorForBranch returns a color based on the branch name convention.
func colorForBranch(branch string) color.Color {
	switch {
	case branch == "main" || branch == "master":
		return style.GetPrimaryColor()
	case strings.HasPrefix(branch, "feature/"):
		return style.GetSuccessColor()
	case strings.HasPrefix(branch, "release/"):
		return style.GetSuccessColor()
	case strings.HasPrefix(branch, "bugfix/"):
		return style.GetDangerColor()
	case strings.HasPrefix(branch, "task/"):
		return style.GetInfoColor()
	default:
		return style.GetNeutralColor()
	}
}

// BuildHistoryContent renders git commit history lazygit-style.
func BuildHistoryContent(project ProjectDTO) string {
	logOut, err := gitCmd(project.Path, "log", "--pretty=format:%h|%an|%cr|%s", "-n", "15") //TODO set the number os commits showing on available space in the tab
	if err != nil || strings.TrimSpace(logOut) == "" {
		return "(no commits)"
	}

	lines := strings.Split(strings.TrimSpace(logOut), "\n")
	var hb strings.Builder

	prStyle := lipgloss.NewStyle().Padding(0, 0, 0, 2).Foreground(style.GetGoldenColor())
	normalStyle := lipgloss.NewStyle().Padding(0, 0, 0, 2).Foreground(style.GetInfoColor())

	for i, ln := range lines {
		parts := strings.SplitN(ln, "|", 4)
		if len(parts) < 4 {
			continue
		}

		author := lipgloss.NewStyle().Padding(0, 0, 0, 8).Foreground(style.GetNeutralColor()).Bold(true).Render(parts[1])
		commitDate := lipgloss.NewStyle().Foreground(style.GetNeutralColor()).Render(parts[2])
		commitMsg := lipgloss.NewStyle().Align(lipgloss.Left).Foreground(style.GetNeutralColor()).Render(parts[3])

		var hash string
		if strings.Contains(parts[3], "Pull request") {
			hash = prStyle.Render(parts[0])
		} else {
			hash = normalStyle.Render(parts[0])
		}

		hb.WriteString(hash)
		hb.WriteString(" ")
		hb.WriteString(commitMsg)
		hb.WriteString("\n")
		hb.WriteString("   ")
		hb.WriteString(author)
		hb.WriteString(" ")
		hb.WriteString(commitDate)
		if i < len(lines)-1 {
			hb.WriteString("\n")
		}
	}
	return hb.String()
}

// get current branch and status for a project, returning a formatted string with colors and change counts.
func LoadProjectGitInfo(project ProjectEntry) string {
	reBranchClean := regexp.MustCompile(`HEAD detached|fatal`)

	branch, err := gitCmd(project.Path, "rev-parse", "--abbrev-ref", "HEAD")

	if err != nil || reBranchClean.MatchString(branch) {
		branch = ""
	}

	porcelain, _ := gitCmd(project.Path, "status", "--porcelain")
	added, modified, deleted := countChanges(porcelain)

	var parts []string
	if added > 0 {
		parts = append(parts, lipgloss.NewStyle().Foreground(style.GetSuccessColor()).Render(fmt.Sprintf("+%d", added)))
	}
	if modified > 0 {
		parts = append(parts, lipgloss.NewStyle().Foreground(style.GetNeutralColor()).Render(fmt.Sprintf("~%d", modified)))
	}
	if deleted > 0 {
		parts = append(parts, lipgloss.NewStyle().Foreground(style.GetDangerColor()).Render(fmt.Sprintf("-%d", deleted)))
	}

	changes := strings.Join(parts, " ")

	var branchDisplay string
	if branch != "" {
		branchDisplay = lipgloss.NewStyle().Foreground(colorForBranch(branch)).Render(branch)
	} else {
		branchDisplay = branch
	}

	return branchDisplay + " " + changes
}

func LoadChangedFiles(project ProjectEntry) []FileChange {
	rePath := regexp.MustCompile(`^(.{2})\s(.*)$`)

	porcelain, _ := gitCmd(project.Path, "status", "--porcelain")

	var editedFiles []FileChange
	if strings.TrimSpace(porcelain) != "" {
		for _, l := range strings.Split(strings.TrimSpace(porcelain), "\n") {
			if l == "" {
				continue
			}
			var path, code string
			if m := rePath.FindStringSubmatch(l); m != nil {
				status := strings.TrimSpace(m[1])
				path = strings.TrimSpace(m[2])
				switch {
				case strings.Contains(status, "?"):
					code = "?"
				case strings.Contains(status, "A"):
					code = "A"
				case strings.Contains(status, "D"):
					code = "D"
				case strings.Contains(status, "M"):
					code = "M"
				default:
					code = status
				}
			} else {
				path = strings.TrimSpace(l)
			}
			if strings.Contains(path, "->") {
				parts := strings.Split(path, "->")
				path = strings.TrimSpace(parts[len(parts)-1])
			}
			editedFiles = append(editedFiles, FileChange{Path: path, Code: code})
		}
	}
	return editedFiles
}

func GetGitWorktrees(project ProjectEntry) []ProjectDTO {
	worktrees, _ := gitCmd(project.Path, "worktree", "list")

	var items []ProjectDTO
	for _, worktree := range strings.Split(strings.TrimSpace(worktrees), "\n") {
		if worktree == "" {
			continue
		}
		// git worktree list output: "<path>  <hash>  [<branch>]"
		fields := strings.Fields(worktree)
		if len(fields) < 1 {
			continue
		}

		path := fields[0]
		name := path
		if idx := strings.LastIndexAny(path, "/\\"); idx >= 0 {
			name = path[idx+1:]
		}
		entry := ProjectEntry{
			Name:    name,
			Path:    path,
			IsGit:   true,
			Options: project.Options,
		}

		items = append(items, AddProjectToList(entry))
	}

	return items
}
