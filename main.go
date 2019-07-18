package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// to properties
const username = "linden-honey-user"
const password = "linden-honey-pass"
const dbname = "linden-honey"

var err error

// to another package nad export plz
type connection struct {
	DB *sql.DB
}

func (db *connection) connect(user, pass, schema string) {

	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, pass, schema)

	// rewrite
	db.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	conn := connection{}
	conn.connect(username, password, dbname)
}
