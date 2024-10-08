//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type Page struct {
	Title string
	Body  []byte
}

var templates = template.Must(template.ParseFiles("templates/header.html", "templates/index.html", "templates/edit.html", "templates/view.html", "templates/error.html", "templates/info.html"))
var validPath = regexp.MustCompile("^/(delete|edit|save|view)/([a-zA-Z0-9_-]+)$")
var validStaticPath = regexp.MustCompile("^/(js|css)/([a-zA-Z0-9_./-]+)$")
var pages_path string

func (p *Page) save() error {
	// filename := "pages/" + p.Title + ".txt"
	// Clean the path to avoid Path Traversal attacks - #1
	filename := filepath.Join(pages_path, p.Title+".txt")
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	// filename := "pages/" + title + ".txt"
	// Clean the path to avoid Path Traversal attacks - #1
	filename := filepath.Join(pages_path, title+".txt")
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		renderTemplate(w, "error", title)
	} else {
		renderTemplate(w, "view", p)
	}
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}

	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func deleteHandler(w http.ResponseWriter, r *http.Request, title string) {
	// Clean the path to avoid Path Traversal attacks - #1
	filename := filepath.Join(pages_path, title+".txt")
	err := os.Remove(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func getFileNamesFromPath(files []string) {
	for index, file := range files {
		base := filepath.Base(file)
		ext := filepath.Ext(base)
		files[index] = strings.TrimSuffix(base, ext)
	}
}

func searchFilesInDirpath(dirpath, searchTerm string) ([]string, error) {
	pattern := dirpath + "/*" + searchTerm + "*.txt"
	//log.Println("Search pattern: " + pattern)
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	getFileNamesFromPath(files)
	//log.Println(files)
	sort.Slice(files, func(i, j int) bool { return files[i] < files[j] })

	return files, nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("q")
	pages, err := searchFilesInDirpath(pages_path, searchTerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderTemplate(w, "index", pages)
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	// log.Println("Serve static content: " + r.URL.Path)
	m := validStaticPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, r.URL.Path[1:])
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	host, err := os.Hostname()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderTemplate(w, "info", host)
}

func main() {
	pages_path = os.Getenv("GOWIKI_PAGES_PATH")

	if pages_path == "" {
		pages_path = "./pages"
	}

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/js/", staticHandler)
	http.HandleFunc("/css/", staticHandler)
	http.HandleFunc("/info", infoHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/delete/", makeHandler(deleteHandler))

	port := os.Getenv("GOWIKI_LISTEN_PORT")

	if port == "" {
		port = "8888"
	}

	fmt.Println("Start webserver at *:" + port + ", and pages are stored at " + pages_path)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
