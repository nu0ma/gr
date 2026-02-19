package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/nu0ma/gr/internal/gitutil"
	"github.com/nu0ma/gr/internal/gitwt"
	"github.com/nu0ma/gr/internal/state"
	"github.com/spf13/cobra"
)

func reviewRun(cmd *cobra.Command, args []string) error {
	branch := args[0]
	ctx := context.Background()

	exists, err := gitwt.Exists(ctx, branch)
	if err != nil {
		return fmt.Errorf("failed to check existing worktrees: %w", err)
	}
	if exists {
		return fmt.Errorf("worktree for branch %q already exists", branch)
	}

	fmt.Printf("Fetching origin/%s...\n", branch)
	if err := gitutil.Fetch(ctx, "origin", branch); err != nil {
		return err
	}

	fmt.Printf("Creating worktree for %s...\n", branch)
	wtPath, err := gitwt.Create(ctx, branch)
	if err != nil {
		return err
	}
	fmt.Printf("Worktree created at: %s\n", wtPath)

	commonDir, err := gitutil.GitCommonDir(ctx)
	if err != nil {
		return err
	}

	st, err := state.Load(commonDir)
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}
	st.Add(state.ReviewSession{
		Branch:       branch,
		WorktreePath: wtPath,
		StartedAt:    time.Now(),
	})
	if err := st.Save(commonDir); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	editor := os.Getenv("GR_EDITOR")
	if editor == "" {
		editor = "nvim"
	}
	fmt.Printf("Opening in %s...\n", editor)
	if err := exec.Command(editor, wtPath).Run(); err != nil {
		fmt.Printf("Warning: failed to open %s: %v\n", editor, err)
	}

	fmt.Printf("\nReview ready! To finish: gr finish %s\n", branch)
	return nil
}
