package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

var (
	config          configuration
	db              *sql.DB
	env             environment
	imageInsertStmt *sql.Stmt
)

func main() {
	env0 := flag.String("env", "dev", "possible values: dev, prod, test")
	configPath := flag.String("config", "config/config.json", "path of the confile file")
	port := flag.String("port", ":8080", "port to run on")
	flag.Parse()

	env = environment(*env0)
	config = loadConfiguration(*configPath)

	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("dbname=%s port=%d user=%s password=%s", config.Db.Name, config.Db.Port, config.Db.User, config.Db.Pass))
	db.SetMaxIdleConns(config.Db.Pool)
	db.SetMaxOpenConns(config.Db.Pool)
	if err != nil {
		log.Fatal(err)
	}

	imageInsertStmt, err = db.Prepare(ImageInsertSql)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Started on %s port", *port)
	http.ListenAndServe(*port, nil)
}
