package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"briefingfy/internal/blog"
	"briefingfy/internal/httpserver"
)

func main() {
	posts, err := blog.LoadPosts(filepath.Join("content", "posts"))
	if err != nil {
		log.Fatalf("load posts: %v", err)
	}

	handler, err := httpserver.New(posts, filepath.Join("web", "templates"), filepath.Join("web", "static"))
	if err != nil {
		log.Fatalf("create server: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	address := ":" + port
	log.Printf("listening on http://localhost%s", address)
	if err := http.ListenAndServe(address, handler); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
