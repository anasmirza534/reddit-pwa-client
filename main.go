package main

import (
	"log"
	"net/http"
	"text/template"
)

func main() {
	http.HandleFunc("/", homeHandler)
	log.Println("Listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/base.html"))
	tmpl.Execute(w, nil)
}
