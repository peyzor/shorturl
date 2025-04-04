package main

import (
	"html/template"
	"log"
	"net/http"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/root.html"))
	data := struct {
		Title string
	}{
		Title: "Short URL",
	}
	err := tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handleRoot)

	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Printf("Server is running on %s", srv.Addr)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
