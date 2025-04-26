package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/peyzor/shorturl/db"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

var baseUrl = "localhost:8080/"
var title = "Short URL"
var shortMaxLen = 4

func randSeq(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	res := make([]rune, length)
	for i := range res {
		res[i] = letters[rand.Intn(len(letters))]
	}
	return string(res)
}

type HomePageData struct {
	Title  string
	Errors []string
}

func renderHomePage(w http.ResponseWriter, data HomePageData) {
	tmpl := template.Must(template.ParseFiles("templates/home.html"))

	err := tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

type ShortsPageData struct {
	Title    string
	Url      string
	Short    string
	ShortUrl string
}

func renderShortsPage(w http.ResponseWriter, data ShortsPageData) {
	tmpl := template.Must(template.ParseFiles("templates/shorts.html"))

	err := tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func (conf *Config) handleHome(w http.ResponseWriter, r *http.Request) {
	data := HomePageData{
		Title: title,
	}

	renderHomePage(w, data)
}

func (conf *Config) handleCreateShort(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	if url == "" {
		log.Println("error occurred: url is not set")

		errors := []string{"URL must be set."}
		data := HomePageData{
			Title:  title,
			Errors: errors,
		}
		renderHomePage(w, data)
		return
	}

	UrlRow, err := conf.Queries.CreateURL(r.Context(), db.CreateURLParams{
		Url:   url,
		Short: randSeq(shortMaxLen),
	})
	if err != nil {
		uniqueConstraintViolationStartsWith := "duplicate key"
		if strings.Contains(err.Error(), uniqueConstraintViolationStartsWith) {
			log.Printf("error occurred: %v\n", err)

			errors := []string{"An error occurred. Please try again."}
			data := HomePageData{
				Title:  title,
				Errors: errors,
			}
			renderHomePage(w, data)
			return
		}
		log.Printf("error occurred: %v\n", err)
		return
	}

	data := ShortsPageData{
		Title:    title,
		Url:      UrlRow.Url,
		Short:    UrlRow.Short,
		ShortUrl: baseUrl + UrlRow.Short,
	}
	renderShortsPage(w, data)
}

func (conf *Config) handleGetShort(w http.ResponseWriter, r *http.Request) {
	short := r.PathValue("short")
	if short == "" {
		log.Println("error occurred: short is not set")
		return
	}

	URLRow, err := conf.Queries.GetShort(r.Context(), short)
	if err != nil {
		log.Printf("error occurred: %v\n", err)
		return
	}

	http.Redirect(w, r, URLRow.Url, http.StatusFound)
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
	mux.HandleFunc("GET /{short}", conf.handleGetShort)

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
