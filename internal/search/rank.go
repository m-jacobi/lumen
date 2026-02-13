package search

import (
	"lumen/internal/index"
	"sort"
	"strings"
)

type Result struct {
	Doc   index.Doc
	Score int
}

func Search(idx *index.Index, q Query, mode Mode, rank bool) []Result {
	if strings.TrimSpace(q.Raw) == "" {
		out := make([]Result, 0, len(idx.Docs))
		for _, d := range idx.Docs {
			out = append(out, Result{Doc: d, Score: 0})
		}
		return out
	}

	out := make([]Result, 0, 256)
	for _, d := range idx.Docs {
		score, ok := matchDoc(d, q, mode, rank)
		if ok {
			out = append(out, Result{Doc: d, Score: score})
		}
	}

	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Score != out[j].Score {
			return out[i].Score > out[j].Score
		}
		ti := strings.ToLower(out[i].Doc.Title)
		tj := strings.ToLower(out[j].Doc.Title)
		if ti != tj {
			return ti < tj
		}
		return out[i].Doc.PathRel < out[j].Doc.PathRel
	})
	return out
}

func matchDoc(d index.Doc, q Query, mode Mode, rank bool) (int, bool) {
	score := 0

	titleSpace := strings.ToLower(d.Title + " " + d.PathRel + " " + d.FileName)
	headSpace := strings.ToLower(strings.Join(d.Headings, " "))
	tagSpace := strings.ToLower(strings.Join(d.Tags, " "))
	content := d.ContentLower

	for _, tok := range q.Tokens {
		t := strings.ToLower(tok)
		found := false

		switch mode {
		case ModeTags:
			needle := normalizeTagQuery(t)
			found = strings.Contains(tagSpace, needle)
			if found && rank {
				score += 200
			}
		case ModeHeadings:
			if strings.Contains(titleSpace, t) {
				found = true
				if rank {
					score += 300
				}
			} else if strings.Contains(headSpace, t) {
				found = true
				if rank {
					score += 250
				}
			}
		case ModeContent:
			found = strings.Contains(content, t)
			if found && rank {
				score += 100
			}
		case ModeAll:
			if strings.Contains(titleSpace, t) {
				found = true
				if rank {
					score += 300
				}
			} else if strings.Contains(headSpace, t) {
				found = true
				if rank {
					score += 250
				}
			} else if strings.Contains(tagSpace, normalizeTagQuery(t)) {
				found = true
				if rank {
					score += 200
				}
			} else if strings.Contains(content, t) {
				found = true
				if rank {
					score += 100
				}
			}
		}

		if !found {
			return 0, false
		}
	}

	if !rank {
		score = 0
	}
	return score, true
}

func normalizeTagQuery(t string) string {
	if strings.HasPrefix(t, "#") {
		return t
	}
	return "#" + t
}
