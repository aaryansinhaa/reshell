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
   - The Go binary parses `aliases.toml`, `env.toml`, and custom function scripts under `functions/`.
   - It formats environment exports (e.g. `export PATH='/opt/bin:$PATH'`) and aliases using single quotes (`'`) instead of double quotes, preventing parameter expansions. Inner single quotes are escaped (`'\''` for bash/zsh, and `\'` for fish).
   - Custom function scripts are validated to ensure they contain exactly one function block matching the script name, and no executable commands outside that block, preventing malicious startup actions.
   - It outputs the compiled startup hook to `~/.config/reshell/shell/reshell.sh` (or `reshell.fish` for Fish shell).

2. **Hook Injection**:
   - The installer checks `~/.bashrc`, `~/.zshrc`, or `~/.config/fish/config.fish`.
   - If the `reshell` sourcing block is missing, it appends it to the end of the file.
   - It includes a prominent ASCII warning block advising developers not to modify the block manually, as `reshell` rebuilds or replaces it dynamically during `apply` or `clean` operations.

3. **Subshell Execution & Context Limitations**:
   - Sourcing is used to bypass this limitation. When the shell starts up, the generated `reshell.sh` script is sourced directly in the current terminal context, allowing aliases, exports, and functions (such as a custom `mkcd`) to operate natively on the parent shell.

## Auto-Discovery & Import Parser

During setup or migration, `reshell` can automatically discover and import configurations from existing shell configuration files (such as `.bashrc`, `.zshrc`, `.profile`, `.bash_aliases`, and `config.fish`), VS Code user snippets, and Pet snippet manager configurations.

1. **Auto-Discovery**:
   - Parses aliases and environment variables using regular expressions.
   - Parses custom functions using a brace-balancing state machine (or Fish block-nesting counter) to extract the function body.
   - Parses VS Code user snippets (`~/.config/Code/User/snippets/*.json` or workspace `*.code-snippets`) and Pet snippet TOML files (`pet.toml`, `snippet.toml`).

2. **Interactive Conflict Resolution & Deduplication**:
   - Identical items are automatically deduplicated.
   - If a parsed item differs from the current active configuration, `reshell` prompts the user interactively to keep the current value, override it, keep both (by renaming), or skip it.

3. **Secrets Detection**:
   - Identifies potential secrets (e.g., tokens, passwords, keys) using name and value heuristics.
   - Flagged variables are skipped by default. The CLI warns the user that while `reshell`'s Git history is purely local and not pushed to any remote by default, plaintext storage is still discouraged.

4. **Target Profile**:
   - Imports can be directed to a specific isolated configuration profile, creating the profile automatically if it does not exist.
