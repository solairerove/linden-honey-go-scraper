package main

import (
	"database/sql"
	"fmt"

	uuid "github.com/satori/go.uuid"

	_ "github.com/lib/pq"
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

	// song := domain.Song{
	// 	Title:  "Something",
	// 	Link:   "pff",
	// 	Author: "Letov",
	// 	Album:  "Experiments",
	// 	Verses: []domain.Verse{
	// 		domain.Verse{
	// 			Ordinal: 1,
	// 			Data:    "?",
	// 		},
	// 	},
	// }

	var id uuid.UUID
	err = db.QueryRow(`INSERT INTO songs(title, link, author, album) 
											VALUES($1, $2, $3, $4) 
											RETURNING id`,
		"s.Title", "s.Link", "s.Author", "s.Album").Scan(&id)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(id)
}
