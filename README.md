# etch

A lightweight terminal note-taking app inspired by Khmer wall carvings. Manage your markdown notes with a Neovim-style dashboard and Telescope-style fuzzy finder.

<!-- TODO: Add screenshot/GIF here -->
<!-- ![Screenshot](screenshot.png) -->

## Features

- **Dashboard** — ASCII art header, quick actions, recent notes
- **Telescope-style finder** — fuzzy search with live preview pane
- **Markdown notes** — YAML frontmatter with title, tags, and dates
- **Tag organization** — auto-grouped by frontmatter tags
- **$EDITOR integration** — opens notes in your preferred editor
- **Soft delete** — notes move to trash before permanent removal
- **Single binary** — no runtime dependencies, works over SSH

## Install

### Homebrew (macOS / Linux)

```bash
brew tap vibolsovichea/etch
brew install etch
```

### Go install

```bash
go install github.com/vibolsovichea/etch/cmd/etch@latest
```

### Binary download

Download pre-built binaries from [GitHub Releases](https://github.com/vibolsovichea/etch/releases).

## Usage

```bash
etch            # Launch the app
etch --version  # Show version info
```

On first run, etch will ask where to store your notes. Notes are plain markdown files — compatible with Obsidian and other markdown editors.

## Keybindings

### Dashboard

| Key     | Action         |
|---------|----------------|
| `f`     | Find notes     |
| `n`     | New note       |
| `q`     | Quit           |
| `j`/`k` | Navigate       |
| `Enter` | Select         |

### Finder

| Key      | Action         |
|----------|----------------|
| Type     | Search         |
| `Enter`  | Open in editor |
| `Ctrl+d` | Delete note    |
| `Ctrl+n` | Next result    |
| `Ctrl+p` | Previous result|
| `Esc`    | Back to dashboard |

## Note Format

```markdown
---
title: My Note
tags: [project, idea]
created: 2026-04-21
modified: 2026-04-21
---

Your note content here...
```

## License

MIT
