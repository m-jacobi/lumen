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
	title, headings = extractHeadings(s, filenameBase)
	tags = collectTags(s)
	return title, headings, tags
}

func extractHeadings(s, filenameBase string) (string, []string) {
	title := filenameBase
	var headings []string

	for _, m := range reHeading.FindAllStringSubmatch(s, -1) {
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

	return title, headings
}

func collectTags(s string) []string {
	seen := map[string]bool{}
	var tags []string

	if fm := extractFrontmatter(s); fm != "" {
		collectYAMLTags(fm, seen, &tags)
	}
	collectInlineTags(s, seen, &tags)

	return tags
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

func collectYAMLTags(fm string, seen map[string]bool, out *[]string) {
	lines := strings.Split(fm, "\n")
	inTags := false

	for _, raw := range lines {
		line := strings.TrimRight(raw, "\r")
		trim := strings.TrimSpace(line)

		if !inTags {
			if strings.HasPrefix(strings.ToLower(trim), "tags:") {
				inTags = true
			}
			collectInlineYAMLTags(trim, seen, out)
			continue
		}

		if strings.HasPrefix(strings.ToLower(trim), "tags:") {
			collectInlineYAMLTags(trim, seen, out)
			inTags = false
			continue
		}

		if strings.HasPrefix(trim, "- ") {
			val := strings.TrimSpace(strings.TrimPrefix(trim, "- "))
			val = strings.Trim(val, `"'`)
			if val == "" {
				continue
			}
			addTag("#"+val, seen, out)
			continue
		}

		if trim != "" {
			break
		}
	}
}

func collectInlineYAMLTags(line string, seen map[string]bool, out *[]string) {
	for _, mm := range reYAMLInlineTags.FindAllStringSubmatch(line, -1) {
		inner := mm[1]
		for _, p := range strings.Split(inner, ",") {
			val := strings.TrimSpace(p)
			val = strings.Trim(val, `"'`)
			if val == "" {
				continue
			}
			addTag("#"+val, seen, out)
		}
	}
}

func collectInlineTags(s string, seen map[string]bool, out *[]string) {
	for _, t := range reInlineTag.FindAllString(s, -1) {
		addTag(t, seen, out)
	}
}

func addTag(tag string, seen map[string]bool, out *[]string) {
	if tag == "" {
		return
	}
	if !seen[tag] {
		seen[tag] = true
		*out = append(*out, tag)
	}
}
