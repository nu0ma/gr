# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`gr` is a Go CLI tool for PR code review using git worktrees. It wraps [git-wt](https://github.com/k1LoW/git-wt) to create isolated review worktrees from branch names and open them in Cursor.

**Runtime dependencies**: `git-wt` (required), Cursor (optional, warns if missing)

## Build & Development Commands

```bash
make build      # Build binary → ./gr
make test       # Run all tests (go test ./...)
make install    # Install to $GOPATH/bin
make clean      # Remove binary

# Run a single test
go test ./internal/state/ -run TestAddAndList

# Run with verbose output
go test -v ./...
```

## Architecture

The project follows a standard Go CLI layout using Cobra:

- **`main.go`** — Entry point, calls `cmd.Execute()`
- **`cmd/`** — Cobra command definitions
  - `root.go` — Root command with `PersistentPreRunE` that validates `git-wt` is installed. Routes to `review` (with args) or `list` (no args)
  - `review.go` — `gr <branch>`: fetches branch, creates worktree via git-wt, opens in Cursor, records session state
  - `list.go` — `gr list`: displays active review sessions in tabular format
  - `finish.go` — `gr finish [branch]`: removes worktree, cleans up state. Supports `-D` for force removal. Can auto-detect branch when run from within worktree
- **`internal/state/`** — JSON-based session persistence (`gr-state.json` stored in `.git/` common dir). Tracks branch, worktree path, and start time
- **`internal/gitwt/`** — Wrapper around `git-wt` CLI. Executes create/remove/list/exists operations, parses JSON output
- **`internal/gitutil/`** — Git helpers: `Fetch(remote, branch)` and `GitCommonDir()`

### Key Design Decisions

- State file lives in git common directory (`.git/`), so it works correctly across worktrees
- All git/external commands accept `context.Context` for cancellation support
- Stateless command execution: each command loads full state, modifies, and saves (no locking)
- Error wrapping with `fmt.Errorf(...: %w, err)` throughout
