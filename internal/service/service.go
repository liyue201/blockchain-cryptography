// Package service 负责加载内嵌教材（content + examples），
// 解析章节元数据并对外提供章节列表与正文。
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

// Chapter 是单个章节的对外结构。
// 标题与正文均支持中英双语：en 为空时前端回退到 zh。
type Chapter struct {
	ID        string `json:"id"`
	Part      int    `json:"part"`
	Slug      string `json:"slug"`
	TitleZh   string `json:"title_zh"`
	TitleEn   string `json:"title_en"`
	Sort      int    `json:"sort"`
	ContentZh string `json:"content_zh"` // Markdown 中文正文
	ContentEn string `json:"content_en"` // Markdown 英文正文（可选）
	Example   string `json:"example"`    // 示例 Go 源码
	ExampleID string `json:"example_id"` // 示例代号（用于前端缓存/提示）
}

// PartMeta 是单个篇（part）的中英标题。
type PartMeta struct {
	Zh string `json:"zh"`
	En string `json:"en"`
}

// Meta 是整本书的中英书名、副标题与篇标题，来自 content/book.json。
type Meta struct {
	TitleZh    string              `json:"title_zh"`
	TitleEn    string              `json:"title_en"`
	SubtitleZh string              `json:"subtitle_zh"`
	SubtitleEn string              `json:"subtitle_en"`
	Parts      map[string]PartMeta `json:"parts"`
}

// Service 持有教材索引。
type Service struct {
	meta     *Meta
	chapters []*Chapter
	byID     map[string]*Chapter
}

// New 从内嵌 FS 加载教材并建立索引。
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
		// 英文正文由同目录的 <name>.en.md 提供（可选，缺省留空、前端回退中文）
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

// List 返回按 (part, sort) 排序的章节列表。
func (s *Service) List() []*Chapter { return s.chapters }

// Get 按 ID 返回章节，不存在时 ok=false。
func (s *Service) Get(id string) (*Chapter, bool) {
	c, ok := s.byID[id]
	return c, ok
}

// Meta 返回整本书的元信息（书名、副标题、篇标题）。
func (s *Service) Meta() *Meta { return s.meta }

// loadMeta 从 content/book.json 读取书的元信息；文件缺失时使用默认元信息。
func (s *Service) loadMeta(contentFS fs.FS) error {
	data, err := fs.ReadFile(contentFS, "book.json")
	if err != nil {
		// book.json 不是必需文件，缺省时仍能运行（标题为空，由前端兜底）
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

// parseChapter 解析单个章节文件：最简 front matter + Markdown 正文。
// front-matter 格式：
//
//	---
//	id: part-0-01
//	title: xxx
//	part: 0
//	sort: 1
//	example: 01-what-is-hash
//	---
//	正文...
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

	// 前端在列表页不需要正文，用 title 区分；这里保留全部。
	return ch, nil
}

// stripFront 去掉一段 markdown 最前面的 front matter，返回正文；
// 无 front matter 时原样返回（英文正文文件可省略 front matter）。
func stripFront(raw string) string {
	if strings.HasPrefix(raw, "---\n") {
		if end := strings.Index(raw, "\n---"); end > 0 {
			raw = raw[end+4:]
		}
	}
	return strings.TrimPrefix(raw, "\n")
}

// loadExample 从 examplesFS 读取约定路径 examples/<id>/main.go。
func loadExample(fsys fs.FS, id string) (string, error) {
	data, err := fs.ReadFile(fsys, path.Join("examples", id, "main.go"))
	if err != nil {
		return "", err
	}
	return string(data), nil
}
