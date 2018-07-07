package main

import (
	"html/template"
	"log"
	"net/http"
	"regexp"
)

var validPath = regexp.MustCompile("^/(intro|play|sendmove|getmove)/([a-zA-Z0-9]*)$")

var templates = template.Must(template.ParseFiles("server/tmpl/play.html"))

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a := validPath.FindStringSubmatch(r.URL.Path)
		if a == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r)
	}
}

func playHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "play.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/play/", makeHandler(playHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
