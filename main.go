package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	domain "github.com/solairerove/linden-honey-go-scrapper/domain"
)

// TODO  to properties
const (
	dbUsername = "linden-honey-user"
	dbPassword = "linden-honey-pass"
	dbName     = "linden-honey"
	dbPort     = "5430"
)

func main() {

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbUsername, dbPassword, dbName, dbPort)

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	song := domain.Song{
		Title:  "Something",
		Link:   "pff",
		Author: "Letov",
		Album:  "Experiments",
		Verses: []domain.Verse{
			domain.Verse{
				Ordinal: 1,
				Data:    "?",
			},
		},
	}

	song.SaveSong(db)
}
