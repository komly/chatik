package handlers

import (
	"html/template"
	"net/http"
)

var tmpl = template.Must(template.ParseGlob("templates/*.html"))

type home struct{}

func NewHome() *home {
	s := &home{}
	return s
}

func (h home) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "index.html", nil)
}
