package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gr <branch>",
	Short: "Git worktree review tool",
	Long:  "A CLI tool that wraps git-wt for quick PR review using git worktrees.",
	Args:  cobra.MaximumNArgs(1),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if _, err := exec.LookPath("git-wt"); err != nil {
			return fmt.Errorf("git-wt not found in PATH.\nInstall with: brew install k1LoW/tap/git-wt\n  or: go install github.com/k1LoW/git-wt@latest")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return listRun(cmd, args)
		}
		return reviewRun(cmd, args)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(finishCmd)
	rootCmd.AddCommand(listCmd)
}
