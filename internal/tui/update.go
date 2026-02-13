package tui

import (
	"fmt"
	"os"
	"path/filepath"

	"lumen/internal/open"
	"lumen/internal/search"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.layout()
		// On resize, wrap likely changes -> force preview rerender
		m.lastSelPath = ""
		m.previewCache = map[string]string{}
		m.updatePreview()
		return m, nil

	case tea.KeyMsg:
		key := msg.String()

		// Quit
		if key == "esc" || key == "ctrl+c" || key == "q" {
			return m, tea.Quit
		}

		// Toggle focus
		if key == "ctrl+p" {
			if m.focus == "list" {
				m.focus = "preview"
				m.status = "focus: preview (j/k, ↑/↓, PgUp/PgDn, g/G)"
			} else {
				m.focus = "list"
				m.status = "focus: list"
			}
			return m, nil
		}

		// Actions that work anytime (open)
		switch key {
		case "enter":
			it, ok := m.list.SelectedItem().(item)
			if ok {
				_ = open.InEditor(m.cfg.EditorCmd, it.path)
				m.status = fmt.Sprintf("opened in %s: %s", m.cfg.EditorCmd, filepath.Base(it.path))
			}
			return m, nil

		case "ctrl+o":
			it, ok := m.list.SelectedItem().(item)
			if ok {
				_ = open.InMacOS(it.path)
				m.status = "opened via macOS open"
			}
			return m, nil

		case "ctrl+r":
			m.cfg.Rank = !m.cfg.Rank
			m.refreshResults()
			return m, nil

		case "tab":
			// cycle modes: all -> tags -> headings -> content
			switch m.cfg.Mode {
			case search.ModeAll:
				m.cfg.Mode = search.ModeTags
			case search.ModeTags:
				m.cfg.Mode = search.ModeHeadings
			case search.ModeHeadings:
				m.cfg.Mode = search.ModeContent
			default:
				m.cfg.Mode = search.ModeAll
			}
			m.refreshResults()
			return m, nil

		case "ctrl+v":
			// quick vault cycle
			names := m.cfg.Vaults.Names()
			cur := 0
			for i, n := range names {
				if n == m.cfg.VaultName {
					cur = i
					break
				}
			}
			next := names[(cur+1)%len(names)]
			p, _ := m.cfg.Vaults.Path(next)
			p = expandHome(p)

			idx, err := indexBuild(p)
			if err == nil {
				m.cfg.VaultName = next
				m.cfg.VaultPath = p
				m.idx = idx
				m.previewCache = map[string]string{}
				m.lastSelPath = ""
				m.refreshResults()
				m.status = "vault: " + next
			} else {
				m.status = "vault switch failed: " + err.Error()
			}
			return m, nil
		}

		// Scroll preview only if preview focused
		if m.focus == "preview" {
			switch key {
			case "down", "j":
				m.vp.LineDown(1)
				return m, nil
			case "up", "k":
				m.vp.LineUp(1)
				return m, nil
			case "pgdown", "ctrl+d":
				m.vp.LineDown(m.vp.Height)
				return m, nil
			case "pgup", "ctrl+u":
				m.vp.LineUp(m.vp.Height)
				return m, nil
			case "home", "g":
				m.vp.GotoTop()
				return m, nil
			case "end", "G":
				m.vp.GotoBottom()
				return m, nil
			}
		}
	}

	// Update query input
	var cmdIn tea.Cmd
	m.in, cmdIn = m.in.Update(msg)

	// Only re-search when query actually changed
	newQ := m.in.Value()
	if newQ != m.lastQuery {
		m.lastQuery = newQ
		m.refreshResults() // includes one preview render for top item
	}

	// Update list; selection changes should update preview (but only when selection changed)
	var cmdList tea.Cmd
	m.list, cmdList = m.list.Update(msg)

	// Update preview if selection changed
	it, ok := m.list.SelectedItem().(item)
	if ok {
		if it.path != m.lastSelPath {
			m.lastSelPath = it.path
			m.updatePreview()
		}
	}

	// Viewport update only needed when preview focused (mouse/scroll events)
	var cmdVP tea.Cmd
	if m.focus == "preview" {
		m.vp, cmdVP = m.vp.Update(msg)
	}

	return m, tea.Batch(cmdIn, cmdList, cmdVP)
}

func expandHome(p string) string {
	if len(p) >= 2 && p[:2] == "~/" {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, p[2:])
	}
	return p
}
