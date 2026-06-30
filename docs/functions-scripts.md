# Functions, Scripts & Workflows

Automate and validate terminal actions using custom functions, scripts, and multi-step workflows.

---

## Custom Function Manager

Unlike aliases, custom functions allow conditional logic and parameters directly within shell profiles. They are stored as separate shell scripts inside `~/.config/reshell/functions/`.

### Sourcing & Sudo Dry-Runs

Because functions run inside the shell process, bad syntax (e.g. unclosed brackets or mismatched quotes) can crash your terminal startup.

reshell prevents this by running **Dry-Run Syntax Checks** before updating your configuration:
- In the TUI, select a function and press `v` (or run `reshell function validate <name>`).
- This executes `bash -n` (or `fish -n` for fish scripts) to check for syntax errors without executing the script.
- If there are errors, they are caught and reported immediately.

### Custom Editor:
Press `e` in the TUI to open the selected function body in your default `$EDITOR` (e.g. Vim, Neovim, or Nano) with full system syntax highlighting, then reload the dashboard when you close the editor.

---

## Script Library

For longer actions (such as docker cleanup or project compiling), save them under categories in the script library folder `~/.config/reshell/scripts/<category>/`.

### Parameters Parser

reshell scans your script files to determine if they expect arguments:
1. **Comment Declarations**: Parses `# @param <Name>` comments.
2. **Positional Arguments**: Detects usages of `$1` through `$9`.

In the TUI, executing a script with parameters displays a form asking for your inputs. Once submitted, these values are passed to the script execution block.

---

## Workflow Manager

Workflows represent sequence steps. They are defined in `~/.config/reshell/workflows.toml`.

### Async Execution Steps

Workflows run sequentially:
- Each step runs in the directory defined in the `dir` field.
- If a step exits with a non-zero code (indicating a failure), the workflow halts to prevent subsequent steps from running in a broken environment.
- The TUI displays real-time progress indicators next to each step using spinners and checkmarks.
- Detailed execution reports containing output streams and exit codes are recorded to `~/.config/reshell/logs/workflows/`.
