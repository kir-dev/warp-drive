package main

import (
	"fmt"
	"github.com/bmizerany/pat"
	"html/template"
	"log"
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

	http.Handle("/", loggerMiddlerware(mux))

	// load templates
	templates = template.Must(template.ParseGlob("template/*.html"))
}

func uploadPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "upload", nil)
}

func loggerMiddlerware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming request: [%s] %s", r.Method, r.URL.String())
		h.ServeHTTP(w, r)
	})
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	if env.dev() {
		tmpl := template.New(makeFullTemplateName(name))
		tmpl, err := tmpl.ParseFiles("template/" + makeFullTemplateName(name))
		if err != nil {
			fmt.Fprintf(w, "Error in template:\n%s", err.Error())

		} else {
			tmpl.Execute(w, data)
		}
		return
	}

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
