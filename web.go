package main

import (
	"fmt"
	"github.com/bmizerany/pat"
	"html/template"
	"io/ioutil"
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
	templates map[string]*template.Template
)

func init() {
	mux := pat.New()

	// register routes
	mux.Get("/upload", http.HandlerFunc(uploadPage))
	mux.Post("/upload", http.HandlerFunc(uploadHandler))

	mux.Get("/", http.HandlerFunc(root))
	http.Handle("/", loggerMiddlerware(mux))

	// load templates
	templates = loadTemplates()
}

func root(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", nil)
}

func uploadPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "upload", nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("image")
	if err != nil {
		renderTemplate(w, "upload", map[string]string{"FormError": "Could not upload file."})
		return
	}
	defer file.Close()

	if err = saveFile(r.FormValue("title"), file, header); err != nil {
		log.Printf("Error while saving image: %v", err)
		renderTemplate(w, "upload", map[string]string{
			"FormError": "Could not upload file: " + err.Error(),
		})
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
		tmpl = templates[makeFullTemplateName(name)]
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

func loadTemplates() map[string]*template.Template {
	const layoutTemplate = "template/layout.html"
	files, err := ioutil.ReadDir("template")
	if err != nil {
		log.Fatal(err)
	}

	result := make(map[string]*template.Template)
	for _, file := range files {
		if file.Name() != "layout.html" {
			result[file.Name()] = template.Must(
				template.ParseFiles(layoutTemplate, "template/"+file.Name()))
		}
	}

	return result
}
