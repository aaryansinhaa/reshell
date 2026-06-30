# reshell Documentation

Welcome to the official documentation for **reshell**, the portable and reproducible Developer Terminal Manager.

reshell helps developers organize, version-control, and share their command-line workspaces. Whether you need to sync shell aliases across multiple systems, run complex workflows, template new projects, or install system dependencies securely, reshell consolidates these workflows under a unified CLI utility and TUI dashboard.

---

## Document Sections

To get started, browse the following topics:

1. **[Getting Started](getting-started.md)**: System requirements, installation, and integration hooks setup.
2. **[Snippets & Aliases](snippets-aliases.md)**: Managing code blocks, system clipboard support, and alias override checks.
3. **[Functions, Scripts & Workflows](functions-scripts.md)**: Dry-running functions, parameterized scripts execution, and workflow automation.
4. **[Packages & Marketplace](package-marketplace.md)**: Synced OS dependency installations, elevated sudo entry piping, and third-party configuration packs.

---

## Design Principles

- **No Hidden State**: All configuration sets are written as human-readable TOML files under `~/.config/reshell`. There are no hidden binary registries or heavy database engines.
- **Dry Sourcing**: Rather than overriding system shells, reshell generates static files sourced by your profile. Sourcing keeps resource overhead at zero and preserves normal terminal startup speeds.
- **Permission Elevation Safety**: Installation actions needing root privileges prompt you securely at runtime rather than requiring the whole reshell process to run as root.
- **100% Portability**: Backing up your `~/.config/reshell` directory or using the ZIP exporter yields an identical terminal workspace on any fresh operating system.
