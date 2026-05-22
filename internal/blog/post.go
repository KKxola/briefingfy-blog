package blog

import (
	"html/template"
	"time"
)

type Post struct {
	Title   string
	Slug    string
	Date    time.Time
	Summary string
	HTML    template.HTML
}
