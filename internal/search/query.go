package search

import "strings"

type Mode int

const (
	ModeAll Mode = iota
	ModeTags
	ModeHeadings
	ModeContent
)

type Query struct {
	Raw    string
	Tokens []string
}

func Parse(raw string) Query {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return Query{Raw: "", Tokens: nil}
	}
	parts := strings.Fields(raw)
	return Query{Raw: raw, Tokens: parts}
}
