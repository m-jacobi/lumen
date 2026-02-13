package index

import (
	"regexp"
	"strings"
)

var (
	reHeading        = regexp.MustCompile(`(?m)^(#{1,6})\s+(.+?)\s*$`)
	reInlineTag      = regexp.MustCompile(`#[A-Za-z0-9_-]+`)
	reYAMLInlineTags = regexp.MustCompile(`(?m)^\s*tags\s*:\s*\[(.*?)\]\s*$`)
)

func ParseMarkdown(s, filenameBase string) (title string, headings []string, tags []string) {
	title = filenameBase

	matches := reHeading.FindAllStringSubmatch(s, -1)
	for _, m := range matches {
		level := len(m[1])
		text := strings.TrimSpace(m[2])
		if text == "" {
			continue
		}
		headings = append(headings, text)
		if level == 1 && title == filenameBase {
			title = text
		}
	}

	seen := map[string]bool{}

	fm := extractFrontmatter(s)
	if fm != "" {
		collectYAMLBlockTags(fm, seen, &tags)
		collectYAMLInlineTags(fm, seen, &tags)
	}

	for _, t := range reInlineTag.FindAllString(s, -1) {
		if !seen[t] {
			seen[t] = true
			tags = append(tags, t)
		}
	}

	return
}

func extractFrontmatter(s string) string {
	if !strings.HasPrefix(s, "---") {
		return ""
	}
	lines := strings.Split(s, "\n")
	if len(lines) < 2 || strings.TrimSpace(lines[0]) != "---" {
		return ""
	}
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			return strings.Join(lines[1:i], "\n")
		}
	}
	return ""
}

func collectYAMLBlockTags(fm string, seen map[string]bool, out *[]string) {
	lines := strings.Split(fm, "\n")
	inTags := false

	for _, raw := range lines {
		line := strings.TrimRight(raw, "\r")
		trim := strings.TrimSpace(line)

		if !inTags {
			if strings.HasPrefix(strings.ToLower(trim), "tags:") {
				inTags = true
			}
			continue
		}

		if strings.HasPrefix(trim, "- ") {
			val := strings.TrimSpace(strings.TrimPrefix(trim, "- "))
			val = strings.Trim(val, `"'`)
			if val == "" {
				continue
			}
			tag := "#" + val
			if !seen[tag] {
				seen[tag] = true
				*out = append(*out, tag)
			}
			continue
		}

		if trim != "" {
			break
		}
	}
}

func collectYAMLInlineTags(fm string, seen map[string]bool, out *[]string) {
	m := reYAMLInlineTags.FindAllStringSubmatch(fm, -1)
	for _, mm := range m {
		inner := mm[1]
		parts := strings.Split(inner, ",")
		for _, p := range parts {
			val := strings.TrimSpace(p)
			val = strings.Trim(val, `"'`)
			if val == "" {
				continue
			}
			tag := "#" + val
			if !seen[tag] {
				seen[tag] = true
				*out = append(*out, tag)
			}
		}
	}
}
