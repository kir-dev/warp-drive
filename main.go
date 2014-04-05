package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

var (
	config configuration
	db     sql.DB
)

func main() {
	println("Hello world!")

	// TODO: make path configrable by command line args
	config = loadConfiguration("config/config.json")

	db, err := sql.Open("postgres", fmt.Sprintf("dbname=%s user=%s password=%s", config.Db.Name, config.Db.User, config.Db.Pass))
	db.SetMaxIdleConns(config.Db.Pool)
	db.SetMaxOpenConns(config.Db.Pool)
	if err != nil {
		log.Fatalln(err)
	}

	// TODO: make port configurable
	http.ListenAndServe(":8080", nil)
}
