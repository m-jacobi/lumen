<h1 align="center">lumen</h1>

<p align="center">
  <strong>Illuminate your Obsidian vault</strong><br/>
  <sub>Fast · Keyboard-driven · Beautiful Markdown preview</sub>
</p>

<p align="center">
  <a href="https://go.dev/">
    <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white" />
  </a>
  <img src="https://img.shields.io/badge/License-MIT-green" />
  <img src="https://img.shields.io/badge/Platform-macOS%20%7C%20Linux-lightgrey" />
</p>

<p align="center">
  A fast, lightweight terminal interface for searching and previewing your Obsidian vault.
</p>

---

## ✨ Features

- 🔎 Full-text search across `.md` files
- 🏷 YAML frontmatter tags + inline `#tags`
- 🧾 Title & heading search
- 📄 Content-only search
- ⭐ Optional ranking (title > headings > tags > content)
- 🗂 Multiple vaults
- 🎨 Beautiful Markdown preview (powered by Glamour)
- ⚡ Debounced search (120ms by default)
- 🚀 Preview caching for fast navigation
- ⌨ Fully keyboard-driven UI

---

## 🚀 Usage

```bash
lumen
lumen #git
lumen commit amend
```

### Select a vault

```bash
lumen work #meeting
lumen private taxes
```

### Search Modes

| -t | Tags only        |
|----|------------------|
| -H | Titles & heading |
| -c | Content only     |
| -r | Enable ranking   |

#### Examples

```bash
lumen -t #git
lumen -H rebase
lumen -c docker compose
lumen -r commit amend
```

## 🏷 Supported Tag Formats

`lumen` detects tags in both formats:

### YAML frontmatter

```yaml
---
tags:
  - git
  - cheatsheet
---
```

### Inline tags

```md
#git #docker
```

Both become searchable as:

```bash
lumen -t git
lumen -t #git
```

## 📦 Installation

```bash
git clone https://github.com/m-jacobi/lumen
cd lumen
go mod tidy
go build -o lumen
```

Move to your PATH:

```bash
mv lumen /usr/local/bin/
```

## 🛠 Built With

- Cobra
- Bubble Tea
- Bubbles
- Lipgloss
- Glamour

## 🗺 Roadmap Ideas

- Config file (~/.config/lumen/config.yaml)
- Live file watching (fsnotify)
- Tag explorer panel
- Backlink navigation
- Daily note shortcut
- Fuzzy scoring
- Vault picker overlay

## 📜 License

MIT
