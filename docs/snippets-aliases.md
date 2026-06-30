# Snippets & Aliases

Learn how to configure bookmarks and shell command mapping overrides.

---

## Snippet Manager

Snippets are multi-line code blocks or shortcuts saved for reuse. They are defined in `~/.config/reshell/snippets.toml`.

### Adding Snippets

Create a snippet using the CLI:

```bash
reshell snippet add mkcd 'mkdir -p "$1" && cd "$1"' "Make and change directory"
```

Or press `n` inside the **Snippets** tab of the TUI.

### Version History

When updating a snippet's code block through the TUI or CLI, the previous content is not discarded. Instead, it is timestamped and appended to the snippet's `history` block inside `snippets.toml`:

```toml
[[snippets]]
name = "mkcd"
code = "mkdir -p \"$1\" && cd \"$1\""
description = "Make and change directory"

[[snippets.history]]
timestamp = "2026-06-28 17:00:00"
code = "mkdir -p \"$1\""
```

You can view the revision history or revert values by inspecting `snippets.toml` directly.

### Actions:
- **Copying**: Press `c` in the TUI to load the selected snippet directly to your system clipboard.
- **Executing**: Press `x` in the TUI to run the snippet code inside a temporary subshell process, showing stdout and stderr before returning to the dashboard.

---

## Alias Manager

Aliases map simple shortcuts to longer shell expressions. They are defined inside `~/.config/reshell/aliases.toml`.

### Duplicate & Conflict Detection

When adding an alias (via the TUI or `reshell alias add`), the manager runs conflict checks to prevent system breakages:

1. **System Override Check**: Queries the path via `exec.LookPath`. If the alias name overrides a binary (e.g. `ls`, `grep`, or `gs`), it prints a warning.
2. **Custom Functions Collision**: Scans the `~/.config/reshell/functions/` directory. If a custom function exists with the same name, it raises a warning.
3. **Duplicate Check**: Checks if the name is already in use by another active alias.

You can still define the alias if you want to force overrides, but the warning flags keep you informed.

### Toggling State

You do not need to delete an alias to disable it. Press `Space` in the **Aliases** tab of the TUI to toggle its `Enabled` state. Disabled aliases are omitted from the compiled output during the next `reshell apply` cycle.
