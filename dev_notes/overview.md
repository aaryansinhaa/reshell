# Welcome to reshell! ⚙️

If you've ever set up a new laptop and spent hours copying aliases, recovering shell scripts from old directories, setting up `$PATH` exports, and installing missing utilities, `reshell` was built for you.

Here is a plain-English, natural-language explanation of how the system works under the hood.

---

## 💡 The Core Idea: "Reproducible Workspaces"

Think of `reshell` as a git repository for your shell profile. Instead of editing your `.bashrc` or `.zshrc` directly (which eventually turns into a messy, unmaintainable graveyard of scripts), you declare your setup in simple configuration files inside `~/.config/reshell/`.

`reshell` then reads these files and compiles them into a single, clean shell script. It injects a hook into your shell profile that loads this compiled script automatically. 

If you get a new computer, you just copy your `~/.config/reshell` folder, install `reshell`, run `setup`, and you are instantly home.

---

## 🛠️ The Moving Parts

Here is how the main components behave and when you should use them:

### 1. Aliases (Shortcuts)
Simple mappings from a short word to a longer command (e.g. `gs` -> `git status`). `reshell` checks these to make sure you don't define the same alias twice or overwrite standard commands.

### 2. Environment Variables
Stored in `env.toml`. Useful for configurations like `$EDITOR` or adding custom directories to your `$PATH`. The cool thing here is that they can be toggled on/off with a single spacebar tap inside the TUI dashboard.

### 3. Custom Functions (RAM-resident helpers)
Written as shell script blocks. When you save a function, it is loaded directly into your active shell's memory. Because they run *in* your current shell, they can change your shell context—like a `mkcd` function that creates a folder and immediately `cd`s into it.

### 4. Scripts (Standalone tasks)
Written as shell files. Unlike functions, scripts run as separate child processes. They are perfect for tasks that take parameters (like a script to resize an image or format files). `reshell` scans `# @param` tags in your scripts and builds interactive forms automatically!

### 5. Workflows (Asynchronous pipelines)
Sometimes you need to run a sequence of tasks (e.g. "Step 1: run clean", "Step 2: build code", "Step 3: upload to server"). Workflows let you chain scripts and terminal commands together. They run in the background, update status signs in real-time, halt if any step fails, and write execution log archives.

---

## 🎨 Design & Interface Aesthetics

We wanted the interface to look like a premium developer utility, not a dry retro command console:
* **Translucency**: The container boundaries have no solid colors, letting your terminal emulator's native background transparency and blur configurations flow right through.
* **Automatic Editor Hook**: When you run `reshell` for the first time, it checks your path for installed text editors (`neovim`, `vscode`, `micro`, etc.) and configures itself, saving you configuration steps.
* **Syntax Highlighting**: All file views, function checkers, and template previews render code syntax colors dynamically so they're easy to read.

---

## 🔍 How Sourcing Works (The Go Limitation)

A common issue with writing terminal managers in languages like Go or Rust is that **child processes cannot alter their parent process**. If the `reshell` Go binary runs a command like `cd ~/projects`, the directory changes for the Go process, but as soon as Go exits, your terminal is still sitting in the exact same folder.

To solve this, `reshell` uses **sourcing**:
1. When you edit things in the TUI, they are saved as static configurations.
2. When you compile them (`Ctrl+a` or `reshell apply`), `reshell` generates a compiled shell script: `~/.config/reshell/shell/reshell.sh`.
3. In your `.bashrc` or `.zshrc`, the injector adds a line that imports (`source`) that file.
4. Because the file is sourced, the active terminal reads and applies the alias bindings and directory commands directly in its own thread!
