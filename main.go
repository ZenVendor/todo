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

	parser, err := NewParser(V_LIST, []int{A_OPEN}, map[int]interface{}{})
	if err != nil {
		log.Fatal(err)
	}
	if err = parser.Parse(args); err != nil {
		log.Fatal(err)
	}

	db, err = conf.Prepare()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = parser.Verb.Call(&parser, db, &conf)
	if err != nil {
		log.Fatal(err)
	}
}
