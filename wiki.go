//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
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

var templates = template.Must(template.ParseFiles("templates/header.html", "templates/index.html", "templates/edit.html", "templates/view.html", "templates/error.html"))
var validPath = regexp.MustCompile("^/(delete|edit|save|view)/([a-zA-Z0-9_-]+)$")
var validStaticPath = regexp.MustCompile("^/(js|css)/([a-zA-Z0-9_./-]+)$")

func (p *Page) save() error {
	// filename := "pages/" + p.Title + ".txt"
	// Clean the path to avoid Path Traversal attacks - #1
	filename := filepath.Join("pages", p.Title+".txt")
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	// filename := "pages/" + title + ".txt"
	// Clean the path to avoid Path Traversal attacks - #1
	filename := filepath.Join("pages", title+".txt")
	body, err := ioutil.ReadFile(filename)
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
	err := os.Remove("./pages/" + title + ".txt")
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
	pattern := dirpath + "*" + searchTerm + "*.txt"
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
	pages, err := searchFilesInDirpath("./pages/", searchTerm)
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

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/js/", staticHandler)
	http.HandleFunc("/css/", staticHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/delete/", makeHandler(deleteHandler))

	fmt.Println("Start webserver at *:8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
