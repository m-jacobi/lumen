package index

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuild(t *testing.T) {
	root := t.TempDir()

	paths := []struct {
		rel     string
		content string
	}{
		{rel: "note.md", content: "# Note\nHello"},
		{rel: filepath.Join("sub", "other.md"), content: "# Other\nWorld"},
		{rel: filepath.Join("sub", "ignore.txt"), content: "no"},
		{rel: filepath.Join(".obsidian", "hidden.md"), content: "# Hidden"},
		{rel: filepath.Join("node_modules", "pkg", "file.md"), content: "# Skip"},
	}

	for _, p := range paths {
		abs := filepath.Join(root, p.rel)
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", abs, err)
		}
		if err := os.WriteFile(abs, []byte(p.content), 0o644); err != nil {
			t.Fatalf("write %s: %v", abs, err)
		}
	}

	idx, err := Build(root)
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	if idx.Root != root {
		t.Fatalf("expected root %q, got %q", root, idx.Root)
	}

	if len(idx.Docs) != 2 {
		t.Fatalf("expected 2 docs, got %d", len(idx.Docs))
	}

	if idx.Docs[0].PathRel != "note.md" {
		t.Fatalf("expected first doc path note.md, got %q", idx.Docs[0].PathRel)
	}
	if idx.Docs[1].PathRel != "sub/other.md" {
		t.Fatalf("expected second doc path sub/other.md, got %q", idx.Docs[1].PathRel)
	}

	if idx.Docs[0].Title != "Note" {
		t.Fatalf("expected title %q, got %q", "Note", idx.Docs[0].Title)
	}
	if idx.Docs[0].ContentLower == "" || idx.Docs[1].ContentLower == "" {
		t.Fatalf("expected content to be indexed")
	}
}

func TestBuildWithFileRoot(t *testing.T) {
	tmp := t.TempDir()
	file := filepath.Join(tmp, "file.md")
	if err := os.WriteFile(file, []byte("# Title"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	if _, err := Build(file); err == nil {
		t.Fatalf("expected error when root is a file")
	}
}
