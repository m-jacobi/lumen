package index

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Doc struct {
	PathRel      string
	PathAbs      string
	FileName     string
	Title        string
	Headings     []string
	Tags         []string
	ContentLower string
	ModTime      time.Time
}

type Index struct {
	Root string
	Docs []Doc
}

func Build(root string) (*Index, error) {
	if err := validateRoot(root); err != nil {
		return nil, err
	}

	var docs []Doc
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, werr error) error {
		doc, skipDir, err := buildDoc(root, path, d, werr)
		if err != nil {
			return err
		}
		if skipDir {
			return filepath.SkipDir
		}
		if doc == nil {
			return nil
		}
		docs = append(docs, *doc)
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(docs, func(i, j int) bool {
		return docs[i].PathRel < docs[j].PathRel
	})

	return &Index{Root: root, Docs: docs}, nil
}

func validateRoot(root string) error {
	st, err := os.Stat(root)
	if err != nil {
		return err
	}
	if !st.IsDir() {
		return fmt.Errorf("vault path is not a directory: %s", root)
	}
	return nil
}

func buildDoc(root, path string, d fs.DirEntry, werr error) (*Doc, bool, error) {
	if werr != nil {
		return nil, false, werr
	}
	if d.IsDir() {
		return nil, shouldSkipDir(filepath.Base(path)), nil
	}
	if strings.ToLower(filepath.Ext(d.Name())) != ".md" {
		return nil, false, nil
	}

	info, content, rel, err := readDocFields(root, path, d)
	if err != nil {
		return nil, false, err
	}

	base := strings.TrimSuffix(d.Name(), filepath.Ext(d.Name()))
	title, heads, tags := ParseMarkdown(content, base)

	doc := Doc{
		PathRel:      rel,
		PathAbs:      path,
		FileName:     d.Name(),
		Title:        title,
		Headings:     heads,
		Tags:         tags,
		ContentLower: strings.ToLower(content),
		ModTime:      info.ModTime(),
	}
	return &doc, false, nil
}

func readDocFields(root, path string, d fs.DirEntry) (fs.FileInfo, string, string, error) {
	b, rerr := os.ReadFile(path)
	if rerr != nil {
		return nil, "", "", fmt.Errorf("read %s: %w", path, rerr)
	}
	info, ierr := d.Info()
	if ierr != nil {
		return nil, "", "", fmt.Errorf("stat %s: %w", path, ierr)
	}
	content := string(b)

	rel, rerr := filepath.Rel(root, path)
	if rerr != nil {
		return nil, "", "", fmt.Errorf("rel %s: %w", path, rerr)
	}
	rel = filepath.ToSlash(rel)

	return info, content, rel, nil
}

func shouldSkipDir(base string) bool {
	switch base {
	case ".obsidian", ".trash", "node_modules", ".git":
		return true
	default:
		return false
	}
}
