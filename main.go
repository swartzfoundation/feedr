package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/swartzfoundation/feedr/frontend"
)

var Time string
var Version = "development"

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/*", http.FileServer(http.FS(frontend.Content)).ServeHTTP)
	log.Printf("Starting server on :3000")
	log.Println("Time: ", Time)
	log.Println("Version: ", Version)
	http.ListenAndServe(":3000", r)
}
