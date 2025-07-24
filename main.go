package main

import (
	"database/sql"
	"log"
	"os"
)

func main() {
	var conf Config
	var db *sql.DB
	args := os.Args[1:]

	db, err := conf.Prepare()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	parser, err := NewParser(V_LIST, []int{A_OPEN}, map[int]interface{}{}, &conf)
	if err != nil {
		log.Fatal(err)
	}
	if err = parser.Parse(args); err != nil {
		log.Fatal(err)
	}

	err = parser.Verb.Call(&parser, db)
	if err != nil {
		log.Fatal(err)
	}
}
