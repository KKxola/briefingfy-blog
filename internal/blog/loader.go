package blog

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark"
)

const dateLayout = "2006-01-02"

func LoadPosts(postsDirectory string) ([]Post, error) {
	pattern := filepath.Join(postsDirectory, "*.md")
	paths, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("find posts: %w", err)
	}

	posts := make([]Post, 0, len(paths))
	for _, path := range paths {
		contents, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read post %q: %w", path, err)
		}

		post, err := parsePost(path, contents)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	sort.Slice(posts, func(left, right int) bool {
		if posts[left].Date.Equal(posts[right].Date) {
			return posts[left].Slug < posts[right].Slug
		}
		return posts[left].Date.After(posts[right].Date)
	})

	return posts, nil
}

func parsePost(path string, contents []byte) (Post, error) {
	frontMatter, body, err := splitFrontMatter(string(contents))
	if err != nil {
		return Post{}, fmt.Errorf("parse post %q: %w", path, err)
	}

	fields := parseFrontMatter(frontMatter)
	if err := validateRequiredFields(path, fields); err != nil {
		return Post{}, err
	}

	date, err := time.Parse(dateLayout, fields["date"])
	if err != nil {
		return Post{}, fmt.Errorf("parse post %q: invalid date %q: %w", path, fields["date"], err)
	}

	htmlBody, err := markdownToHTML(body)
	if err != nil {
		return Post{}, fmt.Errorf("parse post %q: render markdown: %w", path, err)
	}

	return Post{
		Title:   fields["title"],
		Slug:    fields["slug"],
		Date:    date,
		Summary: fields["summary"],
		HTML:    template.HTML(htmlBody),
	}, nil
}

func splitFrontMatter(contents string) (string, string, error) {
	normalized := strings.ReplaceAll(contents, "\r\n", "\n")
	if !strings.HasPrefix(normalized, "---\n") {
		return "", "", fmt.Errorf("missing front matter")
	}

	remaining := normalized[len("---\n"):]
	end := strings.Index(remaining, "\n---\n")
	if end == -1 {
		return "", "", fmt.Errorf("missing front matter closing marker")
	}

	frontMatter := remaining[:end]
	body := remaining[end+len("\n---\n"):]
	return frontMatter, body, nil
}

func parseFrontMatter(frontMatter string) map[string]string {
	fields := make(map[string]string)
	for _, line := range strings.Split(frontMatter, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		key, value, found := strings.Cut(line, ":")
		if !found {
			continue
		}

		fields[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}
	return fields
}

func validateRequiredFields(path string, fields map[string]string) error {
	for _, field := range []string{"title", "slug", "date", "summary"} {
		if fields[field] == "" {
			return fmt.Errorf("parse post %q: missing required front matter field %q", path, field)
		}
	}
	return nil
}

func markdownToHTML(markdown string) (string, error) {
	var html bytes.Buffer
	if err := goldmark.Convert([]byte(markdown), &html); err != nil {
		return "", err
	}
	return html.String(), nil
}
