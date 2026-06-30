# Getting Started

This guide explains how to install, initialize, and set up reshell.

## Requirements

- **Go**: Version `1.22` or newer.
- **Operating Systems**: Linux (tested on Ubuntu, Debian, Arch, Fedora), macOS, or Windows (via PowerShell/Git Bash).
- **Git**: Required for the marketplace pack cloner.

---

## Installation

Compile the binary from the repository root:

```bash
go build -o reshell
```

To make the command available system-wide, copy the executable into a folder in your `$PATH`, such as `/usr/local/bin`:

```bash
sudo cp reshell /usr/local/bin/
```

Confirm that the installation succeeded:

```bash
reshell --help
```

---

## Active Shell Integration

reshell generates configuration scripts and hooks them into your startup profile. Run:

```bash
reshell apply
```

This updates files depending on your active shell:
- **Zsh**: Appends integration blocks to `~/.zshrc` and compiles imports to `~/.config/reshell/shell/reshell.sh`.
- **Bash**: Appends integration blocks to `~/.bashrc` and compiles imports to `~/.config/reshell/shell/reshell.sh`.
- **Fish**: Appends integration blocks to `~/.config/fish/config.fish` and compiles imports to `~/.config/reshell/shell/reshell.fish`.

### Sourcing Hooks

The added blocks look like this:

```bash
# +---------------------------------------------------------------+
# |                _ __ ___  ___| |__   ___| | |                  |
# |               | '__/ _ \/ __| '_ \ / _ \ | |                  |
# |               | | |  __/\__ \ | | |  __/ | |                  |
# |               |_|  \___||___/_| |_|\___|_|_|                  |
# |                                                               |
# |       WARNING: THIS BLOCK IS AUTOMATICALLY MANAGED BY RESHELL. |
# |       DO NOT MANUALLY UPDATE OR EDIT THE CODE WITHIN IT.      |
# +---------------------------------------------------------------+
# >>> reshell initialize >>>
if [ -f "$HOME/.config/reshell/shell/reshell.sh" ]; then
    . "$HOME/.config/reshell/shell/reshell.sh"
fi
# <<< reshell initialize <<<
```

To remove the integration hooks and restore your original shell files, run:

```bash
reshell clean
```

---

## Accessing the TUI Dashboard

To launch the slate-dark terminal interface, run the binary without subcommands:

```bash
reshell
```

### Hotkey Basics:
- `Tab`: Cycle forward through navigation panels.
- `Shift+Tab`: Cycle backward through navigation panels.
- `Up/Down` or `k/j`: Scroll lists.
- `n`: Register a new item (toggles input forms).
- `e`: Open the active item in your `$EDITOR`.
- `d`: Delete the selected item.
- `q` or `Ctrl+C`: Quit reshell.
