package blog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadPostsLoadsValidPost(t *testing.T) {
	postsDirectory := t.TempDir()
	writePost(t, postsDirectory, "hello-world.md", `---
title: Hello World
slug: hello-world
date: 2026-05-21
summary: First post
---
# Hello World

Welcome to the blog.
`)

	posts, err := LoadPosts(postsDirectory)
	if err != nil {
		t.Fatalf("LoadPosts returned error: %v", err)
	}

	if len(posts) != 1 {
		t.Fatalf("got %d posts, want 1", len(posts))
	}

	post := posts[0]
	if post.Title != "Hello World" {
		t.Fatalf("got title %q, want %q", post.Title, "Hello World")
	}
	if post.Slug != "hello-world" {
		t.Fatalf("got slug %q, want %q", post.Slug, "hello-world")
	}
	if post.Summary != "First post" {
		t.Fatalf("got summary %q, want %q", post.Summary, "First post")
	}
}

func TestLoadPostsReturnsErrorForMissingRequiredFrontMatter(t *testing.T) {
	postsDirectory := t.TempDir()
	writePost(t, postsDirectory, "missing-summary.md", `---
title: Missing Summary
slug: missing-summary
date: 2026-05-21
---
# Missing Summary
`)

	_, err := LoadPosts(postsDirectory)
	if err == nil {
		t.Fatal("LoadPosts returned nil error, want required field error")
	}
	if !strings.Contains(err.Error(), `missing required front matter field "summary"`) {
		t.Fatalf("got error %q, want missing summary error", err)
	}
}

func TestLoadPostsConvertsMarkdownToHTML(t *testing.T) {
	postsDirectory := t.TempDir()
	writePost(t, postsDirectory, "markdown.md", `---
title: Markdown
slug: markdown
date: 2026-05-21
summary: Markdown conversion
---
This is **strong** text.
`)

	posts, err := LoadPosts(postsDirectory)
	if err != nil {
		t.Fatalf("LoadPosts returned error: %v", err)
	}

	if !strings.Contains(string(posts[0].HTML), "<strong>strong</strong>") {
		t.Fatalf("HTML %q does not contain converted strong text", posts[0].HTML)
	}
}

func TestLoadPostsSortsByDateDescending(t *testing.T) {
	postsDirectory := t.TempDir()
	writePost(t, postsDirectory, "older.md", `---
title: Older
slug: older
date: 2026-05-20
summary: Older post
---
Older post.
`)
	writePost(t, postsDirectory, "newer.md", `---
title: Newer
slug: newer
date: 2026-05-21
summary: Newer post
---
Newer post.
`)

	posts, err := LoadPosts(postsDirectory)
	if err != nil {
		t.Fatalf("LoadPosts returned error: %v", err)
	}

	if got := posts[0].Slug; got != "newer" {
		t.Fatalf("got first slug %q, want %q", got, "newer")
	}
	if got := posts[1].Slug; got != "older" {
		t.Fatalf("got second slug %q, want %q", got, "older")
	}
}

func writePost(t *testing.T, postsDirectory, name, contents string) {
	t.Helper()

	path := filepath.Join(postsDirectory, name)
	if err := os.WriteFile(path, []byte(contents), 0644); err != nil {
		t.Fatalf("write post %q: %v", name, err)
	}
}
