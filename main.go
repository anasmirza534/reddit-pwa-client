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

	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, list)
}

func postDetailHandler(w http.ResponseWriter, r *http.Request) {
	postId := r.PathValue("id")
	list, err := reddit.GetPost(postId)
	if err != nil {
		log.Println(err)
	}

	tmpl := template.Must(template.ParseFiles("templates/comment.html"))
	tmpl.Execute(w, list)
}
