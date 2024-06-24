package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

    _ "modernc.org/sqlite"
)

type Task struct {
    id int
    description string
    done int
    duedate time.Time
    created time.Time
    completed time.Time
    updated time.Time
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
            duedate datetime,
            created datetime not null,
            completed datetime,
            updated datetime not null
        );
    `
    log.Println("Creating table.")
    _, err := db.Exec(query)
    return err
}

func OpenDB(location string) (db *sql.DB, err error) {
    db, err = sql.Open("sqlite", fmt.Sprint(location, "todo.db"))
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
    _, err = db.Exec(query, t.description, t.done, t.duedate, t.created, t.updated)
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
        query = fmt.Sprintf("select count(*) from tasklist where done = 0 and duedate between '2000-01-01' and '%s';", time.Now())
    }
    err = db.QueryRow(query).Scan(&count)
    return
}

func List(db *sql.DB, sw int) (tl TaskList, err error) {
    var query string
    switch sw {
    case SW_OPEN:
        query = "select * from tasklist where done = 0 order by duedate asc nulls last, created ;"
    case SW_CLOSED:
        query = "select * from tasklist where done = 1 order by completed desc;"
    case SW_ALL:
        query = `
            select * 
            from tasklist 
            order by done, completed desc, duedate asc nulls last, created; 
        `
    case SW_OVERDUE:
        query = fmt.Sprintf("select * from tasklist where done = 0 and duedate between '2000-01-01' and '%s';", time.Now().Format("2006-01-02"))
    case SW_DUE:
        query = "select * from tasklist where done = 0 and duedate > '2000-01-01';"
    }
    
    rows, err := db.Query(query) 
    if err != nil {
        return 
    }
    defer rows.Close()
    
    for next := true; next; next = rows.NextResultSet() {
        for rows.Next() {
            var t Task
            var due, comp sql.NullString
            if err = rows.Scan(&t.id, &t.description, &t.done, &due, &t.created, &comp, &t.updated); err != nil {
                return
            }
            if due.Valid {
                duedate, err := time.Parse(time.RFC3339, due.String)
                if err != nil {
                    return tl, err
                }
                t.duedate = duedate 
            }
            if comp.Valid {
                completed, err := time.Parse(time.RFC3339, comp.String)
                if err != nil {
                    return tl, err
                }
                t.completed = completed
            }
            tl = append(tl, t)
        }
    }
    if err = rows.Err(); err != nil {
        return
    }
    return 
}

func Complete(db *sql.DB, taskId int) (err error) {
    query := "update tasklist set done = 1, completed = ?, updated = ? where id = ?;"
    tm := time.Now().Format(time.RFC3339)
    _, err = db.Exec(query, tm, tm, taskId)
    return err
}

func Reopen(db *sql.DB, taskId int) (err error) {
    query := "update tasklist set done = 0, completed = null, updated = ? where id = ?;"
    tm := time.Now().Format(time.RFC3339)
    _, err = db.Exec(query, tm, taskId)
    return err
}


