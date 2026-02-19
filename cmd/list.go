package cmd

import (
	"context"
	"fmt"
	"text/tabwriter"
	"os"
	"time"

	"github.com/nu0ma/gr/internal/gitutil"
	"github.com/nu0ma/gr/internal/state"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List active review worktrees",
	RunE:  listRun,
}

func listRun(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	commonDir, err := gitutil.GitCommonDir(ctx)
	if err != nil {
		return err
	}

	st, err := state.Load(commonDir)
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	st.CleanStale()
	if err := st.Save(commonDir); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	if len(st.Reviews) == 0 {
		fmt.Println("No active reviews.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "BRANCH\tPATH\tAGE")
	for _, r := range st.Reviews {
		age := time.Since(r.StartedAt).Truncate(time.Minute)
		fmt.Fprintf(w, "%s\t%s\t%s\n", r.Branch, r.WorktreePath, age)
	}
	w.Flush()
	return nil
}
