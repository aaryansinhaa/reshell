# System Architecture & Core Flow

This document details the layout and lifecycle of the `reshell` developer terminal manager.

## Design Philosophy

The main objective of `reshell` is environment reproducibility. A developer should be able to check out their terminal configuration (aliases, environment variables, custom functions, package list) and re-initialize it instantly on a fresh workstation. 

We avoid stateful databases. Everything is stored in human-readable TOML files under `~/.config/reshell`.

```
              +-----------------------------------------+
              |           ~/.config/reshell             |
              |   (TOML files: config, aliases, etc.)   |
              +-------------------+---------------------+
                                  |
                                  | Load/Save
                                  v
                    +-------------+-------------+
                    |        reshell CLI        | <---+ Sudo Password
                    |     (TUI & Subcommands)   |
                    +-------------+-------------+
                                  |
                                  | Compile
                                  v
              +-------------------+---------------------+
              |   ~/.config/reshell/shell/reshell.sh    |
              |   (Auto-sourced by .bashrc / .zshrc)    |
              +-----------------------------------------+
```

## Shell Sourcing Mechanism

We use a compiler-injector model to bind configurations to the active terminal:

1. **Generation (`reshell apply`)**:
   - The Go binary parses `aliases.toml`, `env.toml`, and files in `functions/`.
   - It formats environment exports (e.g. `export PATH="/opt/bin:$PATH"`) and aliases (e.g. `alias gs="git status"`).
   - Custom functions (written as `.sh` or `.fish` scripts) are sourced sequentially.
   - It outputs a compiled setup script to `~/.config/reshell/shell/reshell.sh` (or `reshell.fish` for Fish shell).

2. **Hook Injection**:
   - The installer checks `~/.bashrc`, `~/.zshrc`, or `~/.config/fish/config.fish`.
   - If the `reshell` sourcing block is missing, it appends it to the end of the file.
   - It includes a prominent ASCII warning block advising developers not to modify the block manually, as `reshell` rebuilds or replaces it dynamically during `apply` or `clean` operations.

3. **Subshell Execution & Context Limitations**:
   - Because commands run as subprocesses in Go cannot modify the environment of the parent shell, commands like changing directories (`cd`) or setting temporary environment variables cannot be executed directly by invoking `reshell`.
   - Sourcing is used to bypass this limitation. When the shell starts up, the generated `reshell.sh` script is sourced directly in the current terminal context, allowing aliases, exports, and functions (such as a custom `mkcd`) to operate natively on the parent shell.
