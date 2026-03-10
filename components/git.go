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

func LoadProjectGitInfo(project ProjectEntry) string {
	reBranchClean := regexp.MustCompile(`HEAD detached|fatal`)

	branch, err := gitCmd(project.Path, "rev-parse", "--abbrev-ref", "HEAD")

	if err != nil || reBranchClean.MatchString(branch) {
		branch = ""
	}

	porcelain, _ := gitCmd(project.Path, "status", "--porcelain")
	added, modified, deleted := countChanges(porcelain)

	var changes string
	if added > 0 {
		changes += fmt.Sprintf(" +%d", added)
	}
	if modified > 0 {
		changes += fmt.Sprintf(" ~%d", modified)
	}
	if deleted > 0 {
		changes += fmt.Sprintf(" -%d", deleted)
	}

	var branchDisplay string
	if branch != "" {
		branchDisplay = lipgloss.NewStyle().Foreground(colorForBranch(branch)).Render(branch)
	} else {
		branchDisplay = branch
	}

	return branchDisplay
}
