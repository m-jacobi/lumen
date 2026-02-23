package index

import (
	"reflect"
	"testing"
)

func TestParseMarkdownFrontmatterAndInline(t *testing.T) {
	input := "---\n" +
		"tags:\n" +
		"  - alpha\n" +
		"  - beta\n" +
		"tags: [gamma, \"delta\"]\n" +
		"---\n" +
		"# Title\n" +
		"## Heading 2\n" +
		"Text with #inline and #alpha.\n"

	title, headings, tags := ParseMarkdown(input, "note")

	if title != "Title" {
		t.Fatalf("expected title %q, got %q", "Title", title)
	}

	if !reflect.DeepEqual(headings, []string{"Title", "Heading 2"}) {
		t.Fatalf("unexpected headings: %v", headings)
	}

	expectedTags := []string{"#alpha", "#beta", "#gamma", "#delta", "#inline"}
	if !reflect.DeepEqual(tags, expectedTags) {
		t.Fatalf("expected tags %v, got %v", expectedTags, tags)
	}
}

func TestParseMarkdownUsesFilenameWhenNoH1(t *testing.T) {
	input := "## Heading\n" +
		"Text with #one and #two.\n"

	title, headings, tags := ParseMarkdown(input, "fallback")

	if title != "fallback" {
		t.Fatalf("expected title %q, got %q", "fallback", title)
	}

	if !reflect.DeepEqual(headings, []string{"Heading"}) {
		t.Fatalf("unexpected headings: %v", headings)
	}

	expectedTags := []string{"#one", "#two"}
	if !reflect.DeepEqual(tags, expectedTags) {
		t.Fatalf("expected tags %v, got %v", expectedTags, tags)
	}
}

func TestParseMarkdownBrokenFrontmatterIsIgnored(t *testing.T) {
	input := "---\n" +
		"tags: [alpha]\n" +
		"# Title\n" +
		"Text with #inline.\n"

	_, _, tags := ParseMarkdown(input, "note")

	expectedTags := []string{"#inline"}
	if !reflect.DeepEqual(tags, expectedTags) {
		t.Fatalf("expected tags %v, got %v", expectedTags, tags)
	}
}
