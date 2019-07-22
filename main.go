package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"

	scraper "github.com/solairerove/linden-honey-go-scrapper/scraper"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)
	router.HandleFunc("/songs", songs)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func songs(w http.ResponseWriter, r *http.Request) {
	songs := scraper.ScrapLetov()

	json.NewEncoder(w).Encode(songs)
}
