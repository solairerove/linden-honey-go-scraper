package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"

	scraper "github.com/solairerove/linden-honey-go-scraper/scraper"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)
	router.HandleFunc("/songs", songs)

	log.Fatal(http.ListenAndServe(":8080", handlers.CompressHandler(router)))
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func songs(w http.ResponseWriter, r *http.Request) {
	songs := scraper.ScrapLetov()

	json.NewEncoder(w).Encode(songs)
}
