package note

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Note struct {
	Title    string
	Tags     []string
	Created  time.Time
	Modified time.Time
	Body     string
	Path     string // absolute path to the .md file
}

// frontmatter generates YAML frontmatter for the note.
func (n *Note) frontmatter() string {
	tags := "[]"
	if len(n.Tags) > 0 {
		tags = fmt.Sprintf("[%s]", strings.Join(n.Tags, ", "))
	}
	return fmt.Sprintf(`---
title: %s
tags: %s
created: %s
modified: %s
---`, n.Title, tags, n.Created.Format("2006-01-02"), n.Modified.Format("2006-01-02"))
}

// ToMarkdown returns the full file content.
func (n *Note) ToMarkdown() string {
	return n.frontmatter() + "\n\n" + n.Body
}

// Create creates a new note file in the given directory.
func Create(dir, title string, tags []string) (*Note, error) {
	now := time.Now()
	filename := toFilename(title) + ".md"
	path := filepath.Join(dir, filename)

	n := &Note{
		Title:    title,
		Tags:     tags,
		Created:  now,
		Modified: now,
		Body:     "",
		Path:     path,
	}

	if err := os.WriteFile(path, []byte(n.ToMarkdown()), 0644); err != nil {
		return nil, err
	}
	return n, nil
}

// Load reads a note from a markdown file.
func Load(path string) (*Note, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(path, string(data))
}

// Parse parses note content from a string.
func Parse(path, content string) (*Note, error) {
	n := &Note{Path: path}
	scanner := bufio.NewScanner(strings.NewReader(content))

	// Check for frontmatter
	if scanner.Scan() && strings.TrimSpace(scanner.Text()) == "---" {
		for scanner.Scan() {
			line := scanner.Text()
			if strings.TrimSpace(line) == "---" {
				break
			}
			parseFrontmatterLine(n, line)
		}
	}

	// Rest is body
	var body strings.Builder
	for scanner.Scan() {
		body.WriteString(scanner.Text())
		body.WriteString("\n")
	}
	n.Body = strings.TrimSpace(body.String())

	// Fallback title from filename
	if n.Title == "" {
		base := filepath.Base(path)
		n.Title = strings.TrimSuffix(base, ".md")
	}

	return n, nil
}

func parseFrontmatterLine(n *Note, line string) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return
	}
	key := strings.TrimSpace(parts[0])
	val := strings.TrimSpace(parts[1])

	switch key {
	case "title":
		n.Title = val
	case "tags":
		val = strings.Trim(val, "[]")
		if val != "" {
			for _, t := range strings.Split(val, ",") {
				n.Tags = append(n.Tags, strings.TrimSpace(t))
			}
		}
	case "created":
		if t, err := time.Parse("2006-01-02", val); err == nil {
			n.Created = t
		}
	case "modified":
		if t, err := time.Parse("2006-01-02", val); err == nil {
			n.Modified = t
		}
	}
}

// SetBody updates the note body and modified time.
func (n *Note) SetBody(body string) {
	n.Body = body
	n.Modified = time.Now()
}

// Save writes the note back to disk.
func (n *Note) Save() error {
	n.Modified = time.Now()
	return os.WriteFile(n.Path, []byte(n.ToMarkdown()), 0644)
}

// Delete moves the note to the trash directory.
func (n *Note) Delete(trashDir string) error {
	if err := os.MkdirAll(trashDir, 0755); err != nil {
		return err
	}
	dest := filepath.Join(trashDir, filepath.Base(n.Path))
	return os.Rename(n.Path, dest)
}

// ListAll returns all notes in a directory (recursive).
func ListAll(dir string) ([]*Note, error) {
	var notes []*Note
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip hidden subdirectories (but not the root)
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") && path != dir {
			return filepath.SkipDir
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			n, err := Load(path)
			if err != nil {
				return nil // skip unreadable files
			}
			notes = append(notes, n)
		}
		return nil
	})
	return notes, err
}

// toFilename converts a title to a safe filename.
func toFilename(title string) string {
	s := strings.ToLower(title)
	s = strings.ReplaceAll(s, " ", "-")
	// Remove non-alphanumeric except hyphens
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
