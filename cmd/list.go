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
	if _, err := fmt.Fprintln(w, "BRANCH\tPATH\tAGE"); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}
	for _, r := range st.Reviews {
		age := time.Since(r.StartedAt).Truncate(time.Minute)
		if _, err := fmt.Fprintf(w, "%s\t%s\t%s\n", r.Branch, r.WorktreePath, age); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}
	if err := w.Flush(); err != nil {
		return fmt.Errorf("failed to flush output: %w", err)
	}
	return nil
}
