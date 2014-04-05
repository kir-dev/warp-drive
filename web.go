package main

import (
	"github.com/bmizerany/pat"
	"html/template"
	"net/http"
	"strings"
)

const (
	// Maximum number of bytes to be stored in memory when parsing multipart
	// forms. The rest will be stored in temp files.
	MaxMultipartMemory = 100000
)

var (
	templates = template.New("root")
)

func init() {
	mux := pat.New()

	// register routes
	mux.Get("/upload", http.HandlerFunc(uploadPage))

	http.Handle("/", mux)

	// load templates
	templates = template.Must(template.ParseGlob("template/*.html"))
}

func uploadPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "upload", nil)
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	// TODO: re-read templates on every request in dev mode
	if err := templates.ExecuteTemplate(w, makeFullTemplateName(name), data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeFullTemplateName(name string) string {
	if !strings.HasSuffix(name, ".html") {
		return name + ".html"
	}
	return name
}
