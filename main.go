package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/peyzor/shorturl/db"
	"html/template"
	"log"
	"net/http"
)

func (conf *Config) handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
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

func (conf *Config) handleCreateShort(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	if url == "" {
		fmt.Println("url not set")
		return
	}

	Url, err := conf.Queries.CreateURL(r.Context(), db.CreateURLParams{
		Url:   url,
		Short: "xd",
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%+v\n", Url)
	w.Header().Add("Content-Type", "application/json")
	data, _ := json.Marshal(Url)
	_, _ = w.Write(data)
}

type Config struct {
	Queries *db.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "user=postgres dbname=shorturl sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	queries := db.New(conn)

	mux := http.NewServeMux()

	conf := Config{Queries: queries}

	mux.HandleFunc("GET /home", conf.handleHome)
	mux.HandleFunc("POST /shorts", conf.handleCreateShort)

	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Printf("Server is running on %s", srv.Addr)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
