# gr

A CLI tool for PR code review using git worktrees. Wraps [git-wt](https://github.com/k1LoW/git-wt) to create a review worktree from a branch name and open it in your editor.

## Installation

```bash
go install github.com/nu0ma/gr@latest
```

### Prerequisites

- [git-wt](https://github.com/k1LoW/git-wt)
- [nvim](https://neovim.io/) (default editor)

## Usage

```bash
# Create a worktree and open in nvim
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

## Configuration

### `GR_EDITOR`

Set the `GR_EDITOR` environment variable to change the editor used to open worktrees. Defaults to `nvim`.

| Editor | Value |
|--------|-------|
| nvim (default) | `nvim` |
| Cursor | `cursor` |
| VS Code | `code` |
| Zed | `zed` |

```bash
# Use Cursor
GR_EDITOR=cursor gr feature/add-login

# Use VS Code
GR_EDITOR=code gr feature/add-login

# Persist in your shell config
echo 'export GR_EDITOR=cursor' >> ~/.zshrc
```

## Example

```bash
$ gr feature/add-login
Fetching origin/feature/add-login...
Creating worktree for feature/add-login...
Worktree created at: .wt/feature/add-login
Opening in nvim...

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
