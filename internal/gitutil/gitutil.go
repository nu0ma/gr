package gitutil

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

func Fetch(ctx context.Context, remote, branch string) error {
	cmd := exec.CommandContext(ctx, "git", "fetch", remote, branch)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git fetch %s %s failed: %s", remote, branch, strings.TrimSpace(string(out)))
	}
	return nil
}

func GitCommonDir(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--git-common-dir")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git rev-parse --git-common-dir failed: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}
