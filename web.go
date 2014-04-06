package main

import (
	"database/sql"
	"fmt"
	"github.com/bmizerany/pat"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
	mux.Get("/", http.HandlerFunc(root))
	mux.Get("/upload", http.HandlerFunc(uploadPage))
	mux.Post("/upload", http.HandlerFunc(uploadHandler))

	mux.Get("/search", http.HandlerFunc(searchPage))

	mux.Get("/:hash", http.HandlerFunc(getImageHandler))
	mux.Get("/:hash/:width", http.HandlerFunc(getImageHandlerWidth))

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

func getImageHandler(w http.ResponseWriter, r *http.Request) {
	r.URL.Query().Add(":width", "0")
	getImageHandlerWidth(w, r)
}

func getImageHandlerWidth(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get(":hash")
	width, err := strconv.Atoi(r.URL.Query().Get(":width"))

	if err != nil {
		width = 0
	}

	path, err := getImagePath(hash, width)
	switch {
	case err == nil:
		http.ServeFile(w, r, path)
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func searchPage(w http.ResponseWriter, r *http.Request) {
	searchTerm := "%" + r.URL.Query().Get("search-term") + "%"
	sql := "SELECT title, hash FROM images WHERE LOWER(title) LIKE $1"

	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := stmt.Query(searchTerm)
	if err != nil {
		log.Fatal(err)
	}

	data := make(map[string]string)

	for rows.Next() {
		var title string
		var hash string

		if err := rows.Scan(&title, &hash); err != nil {
			log.Fatal(err)
		}

		log.Printf("matched query: " + hash + " " + title)

		data[hash] = title
	}

	renderTemplate(w, "search", data)
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
