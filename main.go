package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/anasmirza534/reddit-pwa-client/internal/reddit"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	http.Handle("/static/", http.StripPrefix(
		"/static/",
		http.FileServer(http.Dir("static")),
	))
	http.HandleFunc("/", withLogging(handleNotFound))
	http.HandleFunc("/home", withLogging(homeHandler))
	http.HandleFunc("GET /post/{id}", withLogging(postDetailHandler))

	log.Println("Listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func withLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		start := time.Now()
		next(rec, r)
		log.Printf("%s %s %d %v", r.Method, r.URL.Path, rec.status, time.Since(start))
	}
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)

	tmpl := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/not-found.html",
	))
	tmpl.ExecuteTemplate(w, "base", nil)
}

func renderError(w http.ResponseWriter, code int, message string) {
	tmpl := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/error.html",
	))

	w.WriteHeader(code)

	tmpl.ExecuteTemplate(w, "base", struct {
		Code    int
		Message string
	}{code, message})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	after := r.URL.Query().Get("after")

	list, err := reddit.GetHome(after)
	if err != nil {
		log.Println(err)
		renderError(w, http.StatusInternalServerError, "Could not load reddit")
		return
	}

	tmpl := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/home.html",
	))

	if r.Header.Get("HX-Request") == "true" {
		// htmx request (Load more click) — return just the new posts + next button,
		// not the whole page shell.
		tmpl.ExecuteTemplate(w, "postList", list)
		return
	}

	tmpl.ExecuteTemplate(w, "base", list)
}

func postDetailHandler(w http.ResponseWriter, r *http.Request) {
	postId := r.PathValue("id")
	detail, err := reddit.GetPost(postId)
	if err != nil {
		log.Println(err)
		renderError(w, http.StatusInternalServerError, "Could not load reddit")
		return
	}

	tmpl := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/comment.html",
	))
	tmpl.ExecuteTemplate(w, "base", detail)
}
