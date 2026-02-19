package gitwt

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type Worktree struct {
	Branch string `json:"branch"`
	Path   string `json:"path"`
}

func Create(ctx context.Context, branch string) (string, error) {
	cmd := exec.CommandContext(ctx, "git-wt", branch)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("git-wt %s failed: %s", branch, strings.TrimSpace(string(exitErr.Stderr)))
		}
		return "", fmt.Errorf("git-wt %s failed: %w", branch, err)
	}
	return strings.TrimSpace(string(out)), nil
}

func Remove(ctx context.Context, branch string, force bool) error {
	flag := "-d"
	if force {
		flag = "-D"
	}
	cmd := exec.CommandContext(ctx, "git-wt", flag, branch)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git-wt %s %s failed: %s", flag, branch, strings.TrimSpace(string(out)))
	}
	return nil
}

func Exists(ctx context.Context, branch string) (bool, error) {
	worktrees, err := List(ctx)
	if err != nil {
		return false, err
	}
	for _, wt := range worktrees {
		if wt.Branch == branch {
			return true, nil
		}
	}
	return false, nil
}

func List(ctx context.Context) ([]Worktree, error) {
	cmd := exec.CommandContext(ctx, "git-wt", "--json")
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("git-wt --json failed: %s", strings.TrimSpace(string(exitErr.Stderr)))
		}
		return nil, fmt.Errorf("git-wt --json failed: %w", err)
	}
	var worktrees []Worktree
	if err := json.Unmarshal(out, &worktrees); err != nil {
		return nil, fmt.Errorf("failed to parse git-wt output: %w", err)
	}
	return worktrees, nil
}
