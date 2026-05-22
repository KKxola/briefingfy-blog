package httpserver

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"briefingfy/internal/blog"
)

func New(posts []blog.Post, templatesDirectory string, staticDirectory string) (http.Handler, error) {
	homeTemplate, err := parsePageTemplate(templatesDirectory, "home.html")
	if err != nil {
		return nil, err
	}

	postTemplate, err := parsePageTemplate(templatesDirectory, "post.html")
	if err != nil {
		return nil, err
	}

	postsBySlug := make(map[string]blog.Post, len(posts))
	for _, post := range posts {
		postsBySlug[post.Slug] = post
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler(homeTemplate, posts))
	mux.HandleFunc("/posts/", postHandler(postTemplate, postsBySlug))
	mux.Handle("/static/", getOnly(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDirectory)))))
	return mux, nil
}

func parsePageTemplate(templatesDirectory string, page string) (*template.Template, error) {
	files := []string{
		filepath.Join(templatesDirectory, "base.html"),
		filepath.Join(templatesDirectory, page),
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		return nil, fmt.Errorf("parse templates for %s: %w", page, err)
	}
	return tmpl, nil
}

func homeHandler(homeTemplate *template.Template, posts []blog.Post) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}

		renderTemplate(w, homeTemplate, struct {
			Posts []blog.Post
		}{
			Posts: posts,
		})
	}
}

func postHandler(postTemplate *template.Template, postsBySlug map[string]blog.Post) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}

		slug := strings.TrimPrefix(r.URL.Path, "/posts/")
		if slug == "" || strings.Contains(slug, "/") {
			http.NotFound(w, r)
			return
		}

		post, found := postsBySlug[slug]
		if !found {
			http.NotFound(w, r)
			return
		}

		renderTemplate(w, postTemplate, struct {
			Post blog.Post
		}{
			Post: post,
		})
	}
}

func getOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func methodNotAllowed(w http.ResponseWriter) {
	w.Header().Set("Allow", http.MethodGet)
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

func renderTemplate(w http.ResponseWriter, tmpl *template.Template, data any) {
	var body bytes.Buffer
	if err := tmpl.ExecuteTemplate(&body, "base", data); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body.Bytes())
}
