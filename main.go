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
	router.HandleFunc("/songs", getSongs)

	log.Fatal(http.ListenAndServe(":8080", handlers.CompressHandler(router)))
}

func index(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	if err != nil {
		log.Fatalf("Something wrong with greeting endppoint: %v", err)
	}
}

func getSongs(w http.ResponseWriter, r *http.Request) {
	songs := scraper.ScrapLetov()

	err := json.NewEncoder(w).Encode(songs)
	if err != nil {
		log.Fatalf("Something wrong with json marshaling: %v", err)
	}
}
