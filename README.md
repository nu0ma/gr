# gr

A CLI tool for PR code review using git worktrees. Wraps [git-wt](https://github.com/k1LoW/git-wt) to create a review worktree from a branch name and open it in Cursor.

## Installation

```bash
go install github.com/nu0ma/gr@latest
```

### Prerequisites

- [git-wt](https://github.com/k1LoW/git-wt)
- [Cursor](https://www.cursor.com/)

## Usage

```bash
# Create a worktree and open in Cursor
gr <branch>

# List active reviews
gr list

# Remove a review worktree
gr finish <branch>

# Remove from within the worktree directory
gr finish

# Force remove
gr finish -D <branch>
```

## Example

```bash
$ gr feature/add-login
Fetching origin/feature/add-login...
Creating worktree for feature/add-login...
Worktree created at: .wt/feature/add-login
Opening in Cursor...

Review ready! To finish: gr finish feature/add-login

$ gr list
BRANCH                PATH                        AGE
feature/add-login     .wt/feature/add-login       5m0s

$ gr finish feature/add-login
Removing worktree for feature/add-login...
Review for feature/add-login finished.
```

## License

MIT
