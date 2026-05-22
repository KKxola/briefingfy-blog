package httpserver

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"briefingfy/internal/blog"
)

func TestHomeReturnsOK(t *testing.T) {
	handler := newTestHandler(t)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", response.Code, http.StatusOK)
	}
	if !strings.Contains(response.Body.String(), "Hello World") {
		t.Fatalf("body %q does not contain post title", response.Body.String())
	}
}

func TestPostReturnsOK(t *testing.T) {
	handler := newTestHandler(t)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/posts/hello-world", nil)
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", response.Code, http.StatusOK)
	}
	if !strings.Contains(response.Body.String(), "<p>Hello from a test post.</p>") {
		t.Fatalf("body %q does not contain post HTML", response.Body.String())
	}
}

func TestMissingPostReturnsNotFound(t *testing.T) {
	handler := newTestHandler(t)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/posts/missing", nil)
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("got status %d, want %d", response.Code, http.StatusNotFound)
	}
}

func TestUnexpectedMethodReturnsMethodNotAllowed(t *testing.T) {
	handler := newTestHandler(t)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/", nil)
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("got status %d, want %d", response.Code, http.StatusMethodNotAllowed)
	}
	if response.Header().Get("Allow") != http.MethodGet {
		t.Fatalf("got Allow header %q, want %q", response.Header().Get("Allow"), http.MethodGet)
	}
}

func newTestHandler(t *testing.T) http.Handler {
	t.Helper()

	templatesDirectory := t.TempDir()
	writeTemplate(t, templatesDirectory, "base.html", `{{define "base"}}<!doctype html><html><head><title>{{template "title" .}}</title></head><body>{{template "content" .}}</body></html>{{end}}`)
	writeTemplate(t, templatesDirectory, "home.html", `{{define "title"}}Home{{end}}{{define "content"}}{{range .Posts}}<a href="/posts/{{.Slug}}">{{.Title}}</a>{{end}}{{end}}`)
	writeTemplate(t, templatesDirectory, "post.html", `{{define "title"}}{{.Post.Title}}{{end}}{{define "content"}}<article>{{.Post.HTML}}</article>{{end}}`)

	staticDirectory := filepath.Join(t.TempDir(), "static")
	if err := os.MkdirAll(staticDirectory, 0755); err != nil {
		t.Fatalf("create static directory: %v", err)
	}

	handler, err := New([]blog.Post{
		{
			Title:   "Hello World",
			Slug:    "hello-world",
			Date:    time.Date(2026, 5, 21, 0, 0, 0, 0, time.UTC),
			Summary: "A test post",
			HTML:    template.HTML("<p>Hello from a test post.</p>"),
		},
	}, templatesDirectory, staticDirectory)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	return handler
}

func writeTemplate(t *testing.T, templatesDirectory, name, contents string) {
	t.Helper()

	path := filepath.Join(templatesDirectory, name)
	if err := os.WriteFile(path, []byte(contents), 0644); err != nil {
		t.Fatalf("write template %q: %v", name, err)
	}
}
