package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

var (
	config configuration
)

func main() {
	println("Hello world!")

	// TODO: make path configrable by command line args
	config = loadConfiguration("config/config.json")
	fmt.Println(config)

	db, err := sql.Open("postgres", fmt.Sprintf("dbname=%s user=%s password=%s", config.Db.Name, config.Db.User, config.Db.Pass))
	db.SetMaxIdleConns(config.Db.Pool)
	db.SetMaxOpenConns(config.Db.Pool)
	if err != nil {
		log.Fatalln(err)
	}

	var count int
	if err = db.QueryRow("SELECT COUNT(*) FROM images").Scan(&count); err != nil {
		log.Printf("Error: %s", err)
	}
	fmt.Println(count)
}
