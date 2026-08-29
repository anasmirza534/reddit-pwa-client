package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/anasmirza534/reddit-pwa-client/internal/reddit"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/home", homeHandler)
	http.HandleFunc("GET /post/{id}", postDetailHandler)

	log.Println("Listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	list, err := reddit.GetHome()
	if err != nil {
		log.Println(err)
	}

	tmpl := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/home.html",
	))
	tmpl.ExecuteTemplate(w, "base", list)
}

func postDetailHandler(w http.ResponseWriter, r *http.Request) {
	postId := r.PathValue("id")
	detail, err := reddit.GetPost(postId)
	if err != nil {
		log.Println(err)
	}

	tmpl := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/comment.html",
	))
	tmpl.ExecuteTemplate(w, "base", detail)
}
