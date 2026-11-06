// Package service loads the embedded content (content + examples),
// parses chapter metadata, and exposes the chapter list and body.
package service

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strconv"
	"strings"
)

// Chapter is the external structure of a single chapter.
// Titles and bodies are bilingual: when en is empty, the frontend falls back to zh.
type Chapter struct {
	ID        string `json:"id"`
	Part      int    `json:"part"`
	Slug      string `json:"slug"`
	TitleZh   string `json:"title_zh"`
	TitleEn   string `json:"title_en"`
	Sort      int    `json:"sort"`
	ContentZh string `json:"content_zh"` // Markdown body in Chinese
	ContentEn string `json:"content_en"` // Markdown body in English (optional)
	Example   string `json:"example"`    // example Go source code
	ExampleID string `json:"example_id"` // example ID (used for frontend caching/hints)
}

// PartMeta is the bilingual title of a single part.
type PartMeta struct {
	Zh string `json:"zh"`
	En string `json:"en"`
}

// Meta is the bilingual book title, subtitle, and part titles from content/book.json.
type Meta struct {
	TitleZh    string              `json:"title_zh"`
	TitleEn    string              `json:"title_en"`
	SubtitleZh string              `json:"subtitle_zh"`
	SubtitleEn string              `json:"subtitle_en"`
	Parts      map[string]PartMeta `json:"parts"`
}

// Service holds the content index.
type Service struct {
	meta     *Meta
	chapters []*Chapter
	byID     map[string]*Chapter
}

// New loads the content from the embedded FS and builds the index.
func New(contentFS, examplesFS fs.FS) (*Service, error) {
	s := &Service{byID: map[string]*Chapter{}}

	if err := s.loadMeta(contentFS); err != nil {
		return nil, err
	}

	err := fs.WalkDir(contentFS, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(p, ".md") || strings.HasSuffix(p, ".en.md") {
			return nil
		}
		data, err := fs.ReadFile(contentFS, p)
		if err != nil {
			return err
		}
		ch, err := parseChapter(p, string(data))
		if err != nil {
			return fmt.Errorf("parse %s: %w", p, err)
		}
		// English body is provided by <name>.en.md in the same directory (optional; left empty and falls back to Chinese when missing)
		if en, err := fs.ReadFile(contentFS, strings.TrimSuffix(p, ".md")+".en.md"); err == nil {
			ch.ContentEn = stripFront(string(en))
		}
		if ex := ch.ExampleID; ex != "" {
			if code, err := loadExample(examplesFS, ex); err == nil {
				ch.Example = code
			}
		}
		s.chapters = append(s.chapters, ch)
		s.byID[ch.ID] = ch
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(s.chapters, func(i, j int) bool {
		if s.chapters[i].Part != s.chapters[j].Part {
			return s.chapters[i].Part < s.chapters[j].Part
		}
		return s.chapters[i].Sort < s.chapters[j].Sort
	})
	return s, nil
}

// List returns the chapter list sorted by (part, sort).
func (s *Service) List() []*Chapter { return s.chapters }

// Get returns the chapter by ID; ok=false when it does not exist.
func (s *Service) Get(id string) (*Chapter, bool) {
	c, ok := s.byID[id]
	return c, ok
}

// Meta returns the book metadata (title, subtitle, part titles).
func (s *Service) Meta() *Meta { return s.meta }

// loadMeta reads the book metadata from content/book.json; uses defaults when the file is missing.
func (s *Service) loadMeta(contentFS fs.FS) error {
	data, err := fs.ReadFile(contentFS, "book.json")
	if err != nil {
		// book.json is not required; it still runs without it (titles empty, handled by the frontend)
		s.meta = &Meta{Parts: map[string]PartMeta{}}
		return nil
	}
	m := &Meta{Parts: map[string]PartMeta{}}
	if err := json.Unmarshal(data, m); err != nil {
		return fmt.Errorf("parse book.json: %w", err)
	}
	s.meta = m
	return nil
}

// parseChapter parses a single chapter file: minimal front matter + Markdown body.
// front-matter format:
//
//	---
//	id: part-0-01
//	title: xxx
//	part: 0
//	sort: 1
//	example: 01-what-is-hash
//	---
//	body...
func parseChapter(file, raw string) (*Chapter, error) {
	ch := &Chapter{}
	rest := raw

	if strings.HasPrefix(raw, "---\n") {
		end := strings.Index(raw, "\n---")
		if end > 0 {
			head := raw[4:end]
			rest = raw[end+4:]
			meta := map[string]string{}
			for _, line := range strings.Split(head, "\n") {
				kv := strings.SplitN(line, ":", 2)
				if len(kv) == 2 {
					meta[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
				}
			}
			ch.ID = meta["id"]
			ch.TitleZh = meta["title"]
			ch.TitleEn = meta["title_en"]
			if ch.TitleEn == "" {
				ch.TitleEn = ch.TitleZh
			}
			ch.ExampleID = meta["example"]
			ch.Part, _ = strconv.Atoi(meta["part"])
			ch.Sort, _ = strconv.Atoi(meta["sort"])
		}
	}

	if ch.ID == "" {
		return nil, fmt.Errorf("missing front matter id")
	}
	if ch.Slug == "" {
		ch.Slug = path.Base(path.Dir(file))
	}

	body := strings.TrimPrefix(rest, "\n")
	ch.ContentZh = body

	// The list page does not need the body, only the title to distinguish entries; here we keep everything.
	return ch, nil
}

// stripFront removes the leading front matter from a markdown snippet and returns the body;
// returns the input unchanged when there is no front matter (English body files may omit front matter).
func stripFront(raw string) string {
	if strings.HasPrefix(raw, "---\n") {
		if end := strings.Index(raw, "\n---"); end > 0 {
			raw = raw[end+4:]
		}
	}
	return strings.TrimPrefix(raw, "\n")
}

// loadExample reads the example Go code from examples/<id>/main.go in the examples FS.
func loadExample(fsys fs.FS, id string) (string, error) {
	data, err := fs.ReadFile(fsys, path.Join("examples", id, "main.go"))
	if err != nil {
		return "", err
	}
	return string(data), nil
}
