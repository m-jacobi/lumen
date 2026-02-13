package tui

import (
	"fmt"
	"lumen/internal/index"
	"lumen/internal/search"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) layout() {
	topH := 3
	bodyH := max(1, m.height-topH)

	leftW := max(30, int(float64(m.width)*0.45))
	rightW := max(20, m.width-leftW-1)

	m.in.Width = m.width - 2
	m.list.SetSize(leftW, bodyH)
	m.vp.Width = rightW
	m.vp.Height = bodyH
}

func (m *Model) refreshResults() {
	q := search.Parse(m.in.Value())
	results := search.Search(m.idx, q, m.cfg.Mode, m.cfg.Rank)

	items := make([]list.Item, 0, len(results))
	for _, r := range results {
		title := r.Doc.Title
		desc := r.Doc.PathRel
		if m.cfg.Rank {
			desc = fmt.Sprintf("[%d] %s", r.Score, r.Doc.PathRel)
		}
		items = append(items, item{
			title: title,
			desc:  desc,
			path:  r.Doc.PathAbs,
			doc:   r.Doc,
			score: r.Score,
		})
	}
	m.list.SetItems(items)

	// After a new search, render preview once for the current selection.
	// Reset lastSelPath to force update.
	m.lastSelPath = ""
	m.updatePreview()
}

func (m *Model) ensureRenderer() {
	wrap := m.vp.Width
	if wrap < 20 {
		wrap = 20
	}
	if m.renderer != nil && m.renderWrap == wrap {
		return
	}

	// Important: avoid AutoStyle (can trigger OSC queries on some terminals)
	r, _ := glamour.NewTermRenderer(
		glamour.WithStylePath("dark"), // change to "light" if your terminal is light
		glamour.WithWordWrap(wrap),
	)
	m.renderer = r
	m.renderWrap = wrap
}

func (m *Model) updatePreview() {
	it, ok := m.list.SelectedItem().(item)
	if !ok {
		m.vp.SetContent("")
		return
	}

	m.ensureRenderer()

	cacheKey := fmt.Sprintf("%s|%d", it.path, m.renderWrap)
	if cached, ok := m.previewCache[cacheKey]; ok {
		m.vp.SetContent(cached)
		return
	}

	b, err := os.ReadFile(it.path)
	if err != nil {
		m.vp.SetContent("preview error: " + err.Error())
		return
	}

	header := fmt.Sprintf("# %s\n\n`%s`\n\n", it.doc.Title, it.doc.PathRel)
	md := header + string(b)

	out, rerr := m.renderer.Render(md)
	if rerr != nil {
		// fallback raw
		m.vp.SetContent(md)
		return
	}

	m.previewCache[cacheKey] = out
	m.vp.SetContent(out)
}

func (m Model) View() string {
	mode := mapMode(m.cfg.Mode)
	rank := "off"
	if m.cfg.Rank {
		rank = "on"
	}
	focus := m.focus
	if focus == "" {
		focus = "list"
	}

	top := lipgloss.NewStyle().Bold(true).Render(
		fmt.Sprintf("lumen • vault: %s • mode: %s • rank: %s • focus: %s  (Ctrl+P focus, Tab mode, Ctrl+R rank, Ctrl+V vault, Enter edit, Ctrl+O open)",
			m.cfg.VaultName, mode, rank, focus),
	)
	inputLine := m.in.View()
	status := lipgloss.NewStyle().Faint(true).Render(m.status)

	topBlock := strings.Join([]string{top, inputLine, status}, "\n")

	body := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.list.View(),
		lipgloss.NewStyle().Render(" "),
		m.vp.View(),
	)

	return topBlock + "\n" + body
}

func mapMode(mo search.Mode) string {
	switch mo {
	case search.ModeTags:
		return "tags"
	case search.ModeHeadings:
		return "title"
	case search.ModeContent:
		return "content"
	default:
		return "all"
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// helper for vault switching without import cycles
func indexBuild(root string) (*index.Index, error) {
	return index.Build(root)
}
