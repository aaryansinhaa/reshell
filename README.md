# ⚒️ reshell

[![Go Version](https://img.shields.io/github/go-mod/go-version/aaryansinhaa/reshell?color=00ADD8)](https://go.dev/)
[![License](https://img.shields.io/github/license/aaryansinhaa/reshell?color=brightgreen)](LICENSE)
[![Tests](https://img.shields.io/badge/tests-passing-brightgreen)](https://github.com/aaryansinhaa/reshell/actions)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-blueviolet.svg)](https://github.com/aaryansinhaa/reshell/pulls)

**reshell** is a modern, reproducible, and highly portable developer environment and command-line workspace manager written in Go. It provides both a powerful command-line interface (`reshell`) and a premium, fully translucent Terminal User Interface (TUI) dashboard built using the Charmbracelet Bubble Tea framework.

With `reshell`, you can easily configure, track, and synchronize your command aliases, script snippets, custom shell functions, environment variables, system packages, and git configurations in a single version-controlled location.

---

## 🚀 Key Features

* 📊 **Unified Translucent TUI Dashboard**: A premium, blur-through terminal user interface that seamlessly integrates with your terminal emulator's background settings.
* ✍️ **First-Time Setup Wizard**: Automatic welcome onboarding that detects your installed text editors (`nvim`, `vim`, `code`, `micro`, `nano`, etc.) and guides you to configure your workspace preferences.
* 🎨 **Syntax Highlighting**: Real-time syntax highlighting for scripts, functions, snippets, and TOML templates powered by the Chroma rendering engine.
* 🔄 **Workspace Setup Generator**: Build customizable, multi-step workspace setup workflows. Minimize unrelated open windows, load your project directories, and launch default web browser links (Linear, Jira, Spotify, etc.) automatically.
* 📦 **System Package Manager**: Asynchronously installs and uninstalls package requirements across Linux (APT, DNF, Pacman), macOS (Homebrew), and Windows (Winget, Chocolatey) with secure sudo password redirection.
* 🏷️ **Smart Shell Hook Compiler**: Compiles all configurations dynamically and hooks them into your `.bashrc`, `.zshrc`, or `config.fish` profile without modifying files manually.
* 🔑 **Environment & Git Profiles**: Keep track of environment paths and Git profiles globally.
* 📤 **Portable Configurations**: Export or import your entire developer setup as a single TOML manifest or ZIP file for instant machine bootstrap.

---

## 🛠️ Installation & Bootstrapping

### Prerequisites
* Go `1.22` or higher
* Git

### Build from Source
```bash
git clone https://github.com/aaryansinhaa/reshell.git
cd reshell
go build -o reshell
```

### Onboarding Setup
Configure global binary path hooks and profile integrations automatically with the setup command:
```bash
./reshell setup
```
This script copy-installs the `reshell` executable into your local user path (`~/.local/bin/`), registers path variables, and injects shell hook integrations. Simply open a new terminal window to begin using `reshell` globally!

---

## 🕹️ CLI Command Quick Reference

| Command | Action |
| :--- | :--- |
| `reshell` | Launches the interactive TUI Dashboard |
| `reshell apply` | Compiles your active configurations and sources them |
| `reshell clean` | Removes reshell's integration blocks from your shell profile |
| `reshell setup` | Installs reshell binary globally and bootstraps configurations |
| `reshell alias add <name> <value> [desc]` | Registers a command alias |
| `reshell snippet add <name> <code> [desc]` | Stores a code block snippet |
| `reshell snippet copy <name>` | Copies snippet contents to your clipboard |
| `reshell function add <name> <code>` | Creates a custom shell script function |
| `reshell function validate <name>` | Runs a dry-run syntax diagnostic check |
| `reshell script run <cat> <name> [args]` | Runs a library script and writes output logs |
| `reshell workflow run <name>` | Runs a workflow sequence asynchronously |
| `reshell new <template> <name>` | Generates a project skeleton boilerplate |
| `reshell install [repo-url]` | Installs configuration packs or system packages |
| `reshell env add <name> <value>` | Registers environment variables |
| `reshell git apply` | Applies git profiles globally |
| `reshell export <toml-path>` | Exports configurations into a single TOML manifest |
| `reshell import <toml-path>` | Imports configurations from a TOML manifest |

---

## 📂 Configuration Architecture

All configurations are stored in your home directory under `~/.config/reshell/`:

```text
~/.config/reshell/
├── config.toml       # User info, preferred editor, packages, marketplace lists
├── aliases.toml      # Active command aliases
├── snippets.toml     # Script snippets & version history
├── env.toml          # Environment variables
├── workflows.toml    # Workflow definitions
├── functions/        # Raw custom function scripts (.sh, .fish)
├── scripts/          # Library scripts grouped by category
└── logs/             # Workflow and script execution logs
```

For comprehensive tutorials, setup guides, and marketplace documentation, refer to the [docs/](docs/) directory.
