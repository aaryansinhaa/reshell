# TUI Architecture & Bubble Tea State Machine

The reshell dashboard is powered by Charmbracelet's `bubbletea` (Model-View-Update pattern) and styled using `lipgloss`.

## State Representation

The TUI state is tracked in the `model` struct:

- **Navigation**: Controlled by the `activeTab` field. Changing tabs loads appropriate configuration datasets from the disk, resetting the cursor position (`selectedIdx`).
- **Interactive Forms**: A sub-state toggled by `inputMode`. We use `textinput.Model` slices to collect parameters for snippets, aliases, custom functions, or password entries.
- **Subprocess Suspension**: Bubble Tea suspends terminal rendering whenever editing a file or running a interactive script utility. We achieve this by returning `tea.ExecProcess(cmd, callback)` commands. The process gains access to standard input, output, and error streams before returning control to the terminal UI.

## Component Layout Structure

The layout uses a grid block assembly:

```
+-----------------------------------------------------------------------+
|                              Header View                              |
+------------------+----------------------------------------------------+
|                  |                                                    |
|                  |                                                    |
|                  |                  Main Content View                 |
|   Sidebar View   |  (Lists, previews, form fields, or log viewports)  |
|                  |                                                    |
|                  |                                                    |
+------------------+----------------------------------------------------+
|                              Help Footer                              |
+-----------------------------------------------------------------------+
```

1. **Header**: Contains the ASCII brand, project subtitle, and active terminal status.
2. **Sidebar**: Lists navigation options vertically. We compute highlight selections by mapping the current `activeTab` enum value.
3. **Main Content**: Dynamic panel displaying list views with key bindings, interactive creation forms, or background logs viewports.
4. **Help Footer**: A list of active hotkeys based on the open tab.

## Styling System

We maintain consistent layout properties in `styles.go`:
- **Theme**: Base slate background (`#1e1e2e`), slate border lines (`#585b70`), and primary highlights (`#cba6f7` and `#89b4fa`).
- **Margins & Padding**: Standardized border widths and alignment variables prevent window redraw flickering during resizing operations.
