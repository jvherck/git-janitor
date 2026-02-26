# Git Janitor

A fast, interactive Terminal User Interface (TUI) for cleaning up local Git branches.

## The Problem

If you actively review PRs locally or juggle multiple feature branches, your local repository quickly accumulates stale,
merged, or abandoned branches. Cleaning this up manually is tedious. You either have to delete branches one by one in
the JetBrains IDE Git tool window, or manually look up and type out every single branch name you want to remove in the
terminal:

```bash
git branch -D feature/one feature/two bugfix/three stale/four

```

To avoid typing out a massive list, you *can* technically bulk-delete merged branches using native Git alongside bash
utilities, but it requires chaining together a terrifying command pipeline like this:

```bash
git branch --merged | egrep -v "(^\*|main|master|dev)" | xargs git branch -d

```

This native workaround has significant drawbacks:

* **Dangerous:** It runs completely blind. A slight typo in your regular expression could instantly wipe out the wrong
  branches.
* **Hard to remember:** Nobody wants to memorize `egrep -v` syntax and pipe operators just to clean up their workspace.
* **Not cross-platform:** This relies on Unix tools (`grep`, `xargs`) that do not work natively for developers using
  standard Windows Command Prompt or PowerShell.

## The Solution

**Git Janitor** streamlines this entire process by replacing manual typing and dangerous bash scripts with a fast,
keyboard-driven UI directly in your terminal.

Here is why it is a better approach:

* **Visual & Interactive:** You see exactly which branches are queued for deletion before any destructive commands are
  run.
* **Smart & Safe:** It dynamically queries your repository to find the true default branch and hard-locks it, along with
  your currently active branch, so they cannot be accidentally deleted.
* **Zero Memorization:** Just type `git-janitor`, press `m` to select all merged branches, and hit Enter.
* **Cross-Platform:** Compiled as a single Go binary, it works flawlessly across macOS, Linux, and Windows without
  relying on external shell utilities.

## Features

* **Keyboard-driven UI:** Navigate and toggle branches quickly without leaving the terminal.
* **Merged Branch Detection:** Automatically flags branches that have already been safely merged into your default
  branch.
* **Bulk Selection:** Select all unprotected branches or exclusively merged branches with a single keystroke.
* **Branch Protection:** Built-in safeguards prevent the accidental deletion of your currently active branch, as well as
  standard branches like `main`, `master`, and `dev`.

## Installation

Requires [Go](https://go.dev/) 1.20 or later.

**Option 1: Install via Go**

```bash
go install github.com/jvherck/git-janitor@latest

```

*(Ensure your `~/go/bin` directory is in your system's `$PATH`)*

**Option 2: Build from source**

```bash
git clone https://github.com/jvherck/git-janitor.git
cd git-janitor
go build -o git-janitor .
# Move or add the binary to your PATH, e.g., sudo mv git-janitor /usr/local/bin/

```

## Usage

Navigate to any local Git repository in your terminal and run:

```bash
git-janitor

```

### Keybindings

| Key            | Action                                  |
|----------------|-----------------------------------------|
| `↑` / `k`      | Move cursor up                          |
| `↓` / `j`      | Move cursor down                        |
| `Space`        | Toggle selection for the current branch |
| `a`            | Select **all** unprotected branches     |
| `m`            | Select only **merged** branches         |
| `c`            | **Clear** all selections                |
| `Enter`        | Proceed to deletion confirmation        |
| `q` / `Ctrl+C` | Quit                                    |

## Configuration

The UI is built using `bubbletea` and `lipgloss`. Colors, margins, and standard protected branch names are centralized
in the `constants.go` file. You can modify these constants and rebuild the binary to adjust the UI to your preferences.

## Contributing

Contributions, issues, and feature requests are welcome. Feel free to check
the [issues page](https://github.com/jvherck/git-janitor/issues) if you want to contribute.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/NewFeature`)
3. Commit your Changes (`git commit -m 'Add NewFeature'`)
4. Push to the Branch (`git push origin feature/NewFeature`)
5. Open a Pull Request

## License

Distributed under the MIT License. See `LICENSE` for more information.
