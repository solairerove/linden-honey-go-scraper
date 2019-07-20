package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	scrapper "github.com/solairerove/linden-honey-go-scrapper/scrapper"
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

	scrapper.ScrapLetov(db)
}
