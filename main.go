package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

var (
	config                 configuration
	db                     *sql.DB
	env                    environment
	imageInsertStmt        *sql.Stmt
	getImagePathStmt       *sql.Stmt
	searchImageByTitleStmt *sql.Stmt
	getImageByHashStmt     *sql.Stmt
	recentImageStmt        *sql.Stmt
)

func main() {
	env0 := flag.String("env", "dev", "possible values: dev, prod, test")
	configPath := flag.String("config", "config/config.json", "path of the confile file")
	port := flag.String("port", ":8080", "port to run on")
	logfile := flag.String("log", "", "path to the logfile; if left empty the app will log to stdout")
	flag.Parse()

	setLogOutput(*logfile)
	env = environment(*env0)
	config = loadConfiguration(*configPath)

	initDb()
	createPreparedStmts()
	createSessionStore()

	log.Printf("Started on %s port in %s mode", *port, env)
	http.ListenAndServe(*port, nil)
}

func initDb() {
	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("dbname=%s port=%d user=%s password=%s", config.Db.Name, config.Db.Port, config.Db.User, config.Db.Pass))
	db.SetMaxIdleConns(config.Db.Pool)
	db.SetMaxOpenConns(config.Db.Pool)
	if err != nil {
		log.Fatal(err)
	}
}

func createPreparedStmts() {
	var err error

	imageInsertStmt, err = db.Prepare(ImageInsertSql)
	if err != nil {
		log.Fatal(err)
	}

	getImagePathStmt, err = db.Prepare(ImageGetPathSql)
	if err != nil {
		log.Fatal(err)
	}

	searchImageByTitleStmt, err = db.Prepare(SearchImageByTitleSql)
	if err != nil {
		log.Fatal(err)
	}

	getImageByHashStmt, err = db.Prepare(ImageByHashSql)
	if err != nil {
		log.Fatal(err)
	}

	recentImageStmt, err = db.Prepare(RecentImageSql)
	if err != nil {
		log.Fatal(err)
	}

}

func setLogOutput(logfile string) {
	if logfile != "" {
		file, err := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic("cannot open file for logging")
		}
		log.SetOutput(file)
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
