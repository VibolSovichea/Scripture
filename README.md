# etch

**Etch** is a terminal-native note-taking system designed for developers who prioritize speed, keyboard-driven workflows, and local-first data.

Inspired by Neovim and Telescope, Etch brings fast fuzzy search, structured markdown notes, and a focused TUI experience into a single portable binary—built to work seamlessly in local and remote (SSH) environments.

<img width="1710" height="1015" alt="image" src="https://github.com/user-attachments/assets/19870f93-ac3d-44c6-ac13-f70fef784c88" />

---

## Why Etch?

Most note-taking tools are GUI-first and introduce friction for terminal-centric workflows.

Etch is built for developers who:

* live in the terminal
* rely on fast navigation and fuzzy search
* prefer plain text and local ownership of data
* frequently work over SSH or remote servers

Etch removes context switching by keeping your notes where your work already happens.

---

## Features

### Terminal-First Experience

* Fully keyboard-driven interface
* No mouse required
* Optimized for speed and low cognitive overhead

### Neovim-Inspired Workflow

* Familiar keybindings (`j`, `k`, `Enter`, etc.)
* Designed to feel natural for Neovim users

### Telescope-Style Fuzzy Finder

* Real-time fuzzy search across notes
* Keyboard navigation with instant feedback
* Integrated preview pane for fast context switching

### Structured Markdown Notes

* Plain `.md` files with YAML frontmatter
* Fields: `title`, `tags`, `created`, `modified`
* Compatible with tools like Obsidian

### Tag-Based Organization

* Automatic grouping via frontmatter tags
* No vendor lock-in or proprietary format

### Safe Deletion

* Soft delete with trash support
* Prevents accidental data loss

### Single Binary Distribution

* No runtime dependencies
* Works out of the box on macOS and Linux
* Ideal for remote environments and SSH sessions

---

## Installation

### Homebrew (macOS / Linux)

```bash
brew tap vibolsovichea/etch
brew install etch
```

### Go Install

```bash
go install github.com/vibolsovichea/etch/cmd/etch@latest
```

### Prebuilt Binaries

Download from GitHub Releases:
https://github.com/VibolSovichea/etch/releases

---

## Usage

```bash
etch            # Launch the application
etch --version  # Show version
```

On first launch, Etch will prompt for a notes directory.

All notes are stored as plain markdown files, ensuring full portability and interoperability.

---

## Interface

### Dashboard

* ASCII header
* Quick actions
* Recent notes

| Key     | Action     |
| ------- | ---------- |
| `f`     | Find notes |
| `n`     | New note   |
| `q`     | Quit       |
| `j`/`k` | Navigate   |
| `Enter` | Select     |

---

### Finder (Fuzzy Search)

| Key      | Action            |
| -------- | ----------------- |
| Type     | Search            |
| `Enter`  | Open note         |
| `Ctrl+d` | Delete note       |
| `Ctrl+n` | Next result       |
| `Ctrl+p` | Previous result   |
| `Esc`    | Back to dashboard |

---

## Note Format

Etch uses standard markdown with YAML frontmatter:

```markdown
---
title: My Note
tags: [project, idea]
created: 2026-04-21
modified: 2026-04-21
---

Your note content here...
```

This format ensures compatibility with other tools and long-term data ownership.

---

## Architecture

Etch is designed as a modular terminal application with a clear separation of concerns:

* **TUI Framework**: Built using an Elm-inspired architecture (event → update → view)
* **State Management**: Predictable state transitions for UI consistency
* **Views**: Isolated components (dashboard, finder, preview)
* **Storage Layer**: File-based (no database), optimized for portability
* **Search**: Fuzzy matching optimized for interactive performance

This design enables maintainability, extensibility, and responsiveness even with large note collections.

---

## Design Principles

* **Local-first** — your data stays on your machine
* **Composable** — integrates with existing markdown workflows
* **Fast by default** — minimal latency in navigation and search
* **Portable** — single binary, zero dependencies
* **Predictable** — consistent keyboard-driven UX

---

## Roadmap

* Improved fuzzy ranking and scoring
* Plugin or extension system
* Customizable keybindings
* Enhanced tag navigation and filtering
* Performance optimizations for large datasets

---

## License

MIT
