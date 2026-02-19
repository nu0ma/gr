package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/nu0ma/gr/internal/gitutil"
	"github.com/nu0ma/gr/internal/gitwt"
	"github.com/nu0ma/gr/internal/state"
	"github.com/spf13/cobra"
)

var forceDelete bool

var finishCmd = &cobra.Command{
	Use:   "finish [branch]",
	Short: "Remove a review worktree",
	Long:  "Remove a review worktree and clean up session state. If no branch is given, detects from current directory.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		commonDir, err := gitutil.GitCommonDir(ctx)
		if err != nil {
			return err
		}

		st, err := state.Load(commonDir)
		if err != nil {
			return fmt.Errorf("failed to load state: %w", err)
		}

		var branch string
		if len(args) > 0 {
			branch = args[0]
		} else {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
			session := st.FindByCwd(cwd)
			if session == nil {
				return fmt.Errorf("no active review found for current directory")
			}
			branch = session.Branch
		}

		fmt.Printf("Removing worktree for %s...\n", branch)
		if err := gitwt.Remove(ctx, branch, forceDelete); err != nil {
			return err
		}

		st.Remove(branch)
		if err := st.Save(commonDir); err != nil {
			return fmt.Errorf("failed to save state: %w", err)
		}

		fmt.Printf("Review for %s finished.\n", branch)
		return nil
	},
}

func init() {
	finishCmd.Flags().BoolVarP(&forceDelete, "force", "D", false, "Force delete worktree")
}
