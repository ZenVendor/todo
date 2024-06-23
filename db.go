package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
    id int
    description string
    done int
    duedate string
    created string
    completed string
    updated string
}
type TaskList []Task

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
            done integer not null,
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

func OpenDB(location string) (db *sql.DB, err error) {
    db, err = sql.Open("sqlite3", fmt.Sprint(location, "todo.db"))
    if err != nil {
        return
    }
    if !TableExists(db) {
        if err = CreateTable(db); err != nil {
            return
        }
    }
    return 
}

func (t Task) AddTask(db *sql.DB) (err error) {
    query := "insert into tasklist (description, done, duedate, created, updated) values (?, ?, ?, ?, ?);"
    tm := time.Now().Format(time.RFC3339)
    _, err = db.Exec(query, t.description, t.done, t.duedate, tm, tm)
    return err
}
    
func Count(db *sql.DB, sw int) (count int, err error) {
    query := "select count(*) from tasklist where done = 0;"
    switch sw {
    case SW_ALL:
        query = "select count(*) from tasklist;"
    case SW_CLOSED: 
        query = "select count(*) from tasklist where done = 1;"
    case SW_OVERDUE:
        tm := time.Now().Format(time.RFC3339)
        query = fmt.Sprintf("select count(*) from tasklist where done = 0 and duedate < %s;", tm)
    }
    err = db.QueryRow(query).Scan(&count)
    return
}

func List(db *sql.DB, sw int) (tl TaskList, err error) {
    query := "select * from tasklist where done = 0;"

    rows, err := db.Query(query)
    if err != nil {
        return 
    }
    defer rows.Close()
    
    for rows.Next() {
        var t Task
        var due, comp sql.NullString
        if err = rows.Scan(&t.id, &t.description, &t.done, &due, &t.created, &comp, &t.updated); err != nil {
            return
        }
        if due.Valid {
            t.duedate = due.String
        }
        if comp.Valid {
            t.completed = comp.String
        }
        tl = append(tl, t)
    }
    if err = rows.Err(); err != nil {
        return
    }
    return 
}

func (t Task) Complete(db *sql.DB) (err error) {
    query := "update tasklist set done = 1, completed = ?, updated = ?;"
    tm := time.Now().Format(time.RFC3339)
    _, err = db.Exec(query, tm, tm)
    return err
}

func (t Task) Reopen(db *sql.DB) (err error) {
    query := "update tasklist set done = 0, completed = null, updated = ?;"
    tm := time.Now().Format(time.RFC3339)
    _, err = db.Exec(query, tm)
    return err
}


