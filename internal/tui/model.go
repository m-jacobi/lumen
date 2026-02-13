package tui

import (
	"lumen/internal/index"
	"lumen/internal/search"
	"lumen/internal/vault"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/glamour"
)

type Config struct {
	VaultName string
	VaultPath string
	Vaults    vault.Vaults

	Mode search.Mode
	Rank bool

	EditorCmd  string
	InitialQ   string
	ShowHidden bool
}

type item struct {
	title string
	desc  string
	path  string
	doc   index.Doc
	score int
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title + " " + i.desc }

type Model struct {
	cfg Config
	idx *index.Index

	in   textinput.Model
	list list.Model
	vp   viewport.Model

	renderer     *glamour.TermRenderer
	renderWrap   int
	previewCache map[string]string

	focus       string
	lastQuery   string
	lastSelPath string

	width  int
	height int
	status string

	debounceID int
	debounceMS int
}

type debouncedSearchMsg struct {
	id int
}

func NewModel(cfg Config, idx *index.Index) Model {
	ti := textinput.New()
	ti.Placeholder = "search… (#tag, words)"
	ti.SetValue(cfg.InitialQ)
	ti.Focus()

	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Notes"
	l.SetShowHelp(false)
	l.DisableQuitKeybindings()

	vp := viewport.New(0, 0)

	m := Model{
		cfg:          cfg,
		idx:          idx,
		in:           ti,
		list:         l,
		vp:           vp,
		previewCache: map[string]string{},
		focus:        "list",
		lastQuery:    cfg.InitialQ,
		lastSelPath:  "",
		renderWrap:   0,
		renderer:     nil,
		debounceMS:   80,
	}

	m.refreshResults()
	return m
}
