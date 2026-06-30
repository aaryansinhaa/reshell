# Component Implementations & Tradeoffs

This document details the mechanics of individual reshell modules.

## Package Installer & Sudo Piping

Installing packages on Linux requires elevation. We chose to prompt the user for their password in the reshell interface rather than running the entire manager application as root.

```
       +------------------+
       |   User Input     | (Password)
       +--------+---------+
                | Stdin Write
                v
  +-------------+-------------+
  |    exec.Command("sudo")   | ---> [apt-get | dnf | pacman] install <pkg>
  +---------------------------+
```

### Flow:
1. `DetectOS` queries `/etc/os-release` or binary paths to find the system package manager.
2. If it is `apt`, `dnf`, or `pacman`, it spawns a subprocess with the `-S` flag (e.g. `sudo -S apt-get install -y <package>`). The `-S` flag directs `sudo` to read the password from standard input.
3. We establish standard input/output/error pipes:
   - The password is written directly to the input pipe (`StdinPipe`).
   - Output channels stream build messages asynchronously to the viewport window.

---

## Script Parameter Scanner

To run custom scripts through the interface, we need to know what parameter inputs are required. We scan files using two methods:

1. **Positionals**: A regex scanner scans for `$1` through `$9` identifiers inside the script code.
2. **Declaratives**: We parse custom metadata tags written in header comments, specifically `# @param <Name>`. This allows developer-friendly naming of arguments (e.g., `# @param ProjectName`) instead of generic positional labels.

---

## Custom Function Validator

Functions are executed directly in the parent terminal context. Loading a function with a syntax error can crash the parent shell startup.

To prevent this, `functions.Validate()` runs a dry-run check before saving:
- It spawns a shell process with the `-n` (noexec) flag: `bash -n ~/.config/reshell/functions/<name>.sh`.
- The shell parses all script lines. If it encounters syntax errors (such as missing brackets or open quotes), it outputs them to stderr and returns a non-zero exit code.
- `reshell` intercepts the output and alerts the developer, preventing broken files from being written.

---

## Configuration Export/Import

The export command bundles user settings into a standard ZIP archive:
- We walk the `~/.config/reshell` directory recursively.
- We skip the `logs/` directory to prevent backing up large log dumps.
- We skip temporary shell files inside `shell/` since they are dynamically rebuilt.
- The standard library `archive/zip` is used to maintain portable archives across Linux, macOS, and Windows.
