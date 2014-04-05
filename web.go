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
	templates *template.Template
)

func init() {
	mux := pat.New()

	// register routes
	mux.Get("/upload", http.HandlerFunc(uploadPage))
	mux.Post("/upload", http.HandlerFunc(uploadHandler))

	http.Handle("/", loggerMiddlerware(mux))

	// load templates
	templates = template.Must(template.ParseGlob("template/*.html"))
}

func uploadPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "upload", nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if err = saveFile(r.FormValue("title"), file, header); err != nil {
		log.Printf("Error while saving image: %v", err)
		http.Error(w, "Could not save image.", http.StatusBadRequest)
	} else {
		// TODO: redirect to image
		fmt.Fprint(w, "Success!")
	}
}

func loggerMiddlerware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming request: [%s] %s", r.Method, r.URL.String())
		h.ServeHTTP(w, r)
	})
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	var tmpl *template.Template
	if env.dev() {
		var err error
		tmpl, err = template.New("t").ParseFiles("template/layout.html", "template/"+makeFullTemplateName(name))
		if err != nil {
			http.Error(w, fmt.Sprintf("Error in template:\n%s", err), http.StatusInternalServerError)
			return
		}
	} else {
		tmpl = templates
	}

	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeFullTemplateName(name string) string {
	if !strings.HasSuffix(name, ".html") {
		return name + ".html"
	}
	return name
}
