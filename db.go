package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func TableExists(db *sql.DB) bool {
    var rowId int
    if err := db.QueryRow("select rowid from sqlite_schema where type = 'table' and tbl_name = 'tasklist';").Scan(&rowId); err != nil {
        if err == sql.ErrNoRows {
            return false
        } else {
            log.Fatal(err)
        }
    }
    return true
}
    
func CreateTable(db *sql.DB) error {
    query := `
        create table tasklist (
            id integer not null primary key,
            description text not null,
            done boolean not null,
            duedate date,
            created date not null,
            completed date,
            updated datetime not null
        );
    `
    log.Println("Creating table.")
    _, err := db.Exec(query)
    return err
}



func OpenDB(location string) (*sql.DB, error) {
    db, err := sql.Open("sqlite3", fmt.Sprint(location, "todo.db"))
    if err != nil {
        log.Fatal(err)
    }
    if TableExists(db) {
        return db, err
    }
    if err = CreateTable(db); err != nil {
        log.Fatal(err)
    }
    return db, err
}

func (t Task) AddTask(db *sql.DB) error {
    stmt, err := db.Prepare("insert into tasklist (description, done, duedate, created, updated) values (?, ?, ?, ?, ?);")
    if err != nil {
        return err
    }
    tm := time.Now().Format(time.RFC3339)
    _, err = stmt.Exec(t.description, t.done, t.duedate, tm, tm)
    return err
}
    
func Count(db *sql.DB, sw string) (int, error) {
    var count int
    var stmt *sql.Stmt

    var err error
    switch sw {
    case "closed":
        stmt, err = db.Prepare("select count(*) from tasklist where done = 'true';")
    case "all":
        stmt, err = db.Prepare("select count(*) from tasklist;")
    default:
        stmt, err = db.Prepare("select count(*) from tasklist where done = 'false';")
    }
    if err != nil {
        return count, err     
    }
    err = stmt.QueryRow().Scan(&count)
    return count, err
}

func List(db *sql.DB, sw string) (tl TaskList, err error) {
    var stmt *sql.Stmt

    switch sw {
    case "closed":
    default:
        stmt, err = db.Prepare("select * from tasklist where done = 'false';")
    }
    if err != nil {
        return
    }
    rows, err := stmt.Query()
    if err != nil {
        return 
    }
    defer rows.Close()
    for rows.Next() {
        var t Task
        if err = rows.Scan(&t.id, &t.description, &t.done, &t.duedate, &t.created, &t.completed, &t.updated); err != nil {
            return
        }
        tl = append(tl, t)
    }
    if err = rows.Err(); err != nil {
        return
    }
    return 
}
