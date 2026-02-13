package index

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Doc struct {
	PathRel  string
	PathAbs  string
	FileName string

	Title    string
	Headings []string
	Tags     []string

	// for fast contains()
	ContentLower string

	ModTime time.Time
}

type Index struct {
	Root string
	Docs []Doc
}

func Build(root string) (*Index, error) {
	st, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !st.IsDir() {
		return nil, fmt.Errorf("vault path is not a directory: %s", root)
	}

	var docs []Doc
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, werr error) error {
		if werr != nil {
			return werr
		}
		if d.IsDir() {
			base := filepath.Base(path)
			if base == ".obsidian" || base == ".trash" || base == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(d.Name()))
		if ext != ".md" {
			return nil
		}

		b, rerr := os.ReadFile(path)
		if rerr != nil {
			return nil
		}
		info, _ := d.Info()

		rel, _ := filepath.Rel(root, path)
		rel = filepath.ToSlash(rel)

		title, heads, tags := ParseMarkdown(string(b), strings.TrimSuffix(d.Name(), ext))

		doc := Doc{
			PathRel:      rel,
			PathAbs:      path,
			FileName:     d.Name(),
			Title:        title,
			Headings:     heads,
			Tags:         tags,
			ContentLower: strings.ToLower(string(b)),
			ModTime:      info.ModTime(),
		}
		docs = append(docs, doc)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &Index{Root: root, Docs: docs}, nil
}
