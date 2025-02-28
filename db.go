package main

import (
	"database/sql"
	"log"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func NullNow() sql.NullTime {
	return sql.NullTime{Time: time.Now(), Valid: true}
}

type Task struct {
	Id             int
	Short          string
	Priority       int
	DateDue        sql.NullTime
	DateCompleted  sql.NullTime
	Long           string
	ClosingComment string
	Status         *Status
	Group          *Group
	Parent         *Task
}
type TaskList []Task

type Group struct {
	Id             int
	Name           string
	SysDateCreated sql.NullTime
	SysDateUpdated sql.NullTime
	SysStatus      int
}
type Status struct {
	Id   int
	Name string
}

type Value struct {
	Name  string
	Value interface{}
}
type Values []Value

func (conf *Config) OpenDB() (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", filepath.Join(conf.DBLocation, conf.DBName))
	return
}

func TableExists(db *sql.DB) bool {
	var rows int
	if err := db.QueryRow("select count(*) from sqlite_schema where type = 'table' and tbl_name in ('Task', 'Group', 'Status');").Scan(&rows); err != nil {
		log.Fatal(err)
	}
	if rows != 3 {
		return false
	}
	return true
}

func CreateTable(db *sql.DB) error {
	query := `
        create table tasklist (
            id integer primary key not null,
            short text not null,
            priority integer null,
            group_id integer not null,
            done integer not null,
            long text,
            comment text,
            due datetime,
            completed datetime,
            task_id integer,
            created datetime not null,
            updated datetime not null
        );
        create table taskgroup (
            id integer primary key not null,
            name text not null,
            created datetime not null,
            updated datetime not null
        );
        insert into taskgroup (name, created, updated) values ('Default', ?, ?);
    `
	log.Println("Creating tables.")
	_, err := db.Exec(query, time.Now(), time.Now())
	return err
}

func (t *Task) Insert(db *sql.DB) (err error) {
	result, err := db.Exec(
		"INSERT INTO Task (short, priority, date_due, status_id, group_id) VALUES (?, ?, ?, ?, ?);",
		t.Short,
		t.Priority,
		t.DateDue,
		t.Status.Id,
		t.Group.Id,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	t.Id = int(id)
	return err
}

func (g *Group) Insert(db *sql.DB) (err error) {
	result, err := db.Exec(
		"INSERT INTO Group (name) VALUES (?);",
		g.Name,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	g.Id = int(id)
	return err
}
func (t *Task) Select(db *sql.DB) (err error) {
	query := `
        SELECT 
            t.id
            , t.short
            , t.priority
            , t.date_due
            , t.date_completed
            , t.long
            , t.closing_comment
            , s.id
            , s.name
            , g.id
            , g.name
        FROM
            Task t
            JOIN Status s ON s.id = t.status_id
            JOIN Group g ON g.id = t.group_id
        WHERE
            t.id = ?;
    `
	err = db.QueryRow(query, t.Id).Scan(
		&t.Id,
		&t.Short,
		&t.Priority,
		&t.DateDue,
		&t.DateCompleted,
		&t.Long,
		&t.ClosingComment,
		&t.Status.Id,
		&t.Status.Name,
		&t.Group.Id,
		&t.Group.Name,
	)
	return err
}

func (g *Group) Select(db *sql.DB, byName bool) (err error) {
	if byName {
		err = db.QueryRow(
			"SELECT id, name FROM Group WHERE g.name = ?;",
			g.Name,
		).Scan(&g.Id, &g.Name)
	} else {
		err = db.QueryRow(
			"SELECT id, name FROM Group WHERE g.id = ?;",
			g.Id,
		).Scan(&g.Id, &g.Name)
	}
	return err
}

func (t Task) Update(db *sql.DB) (err error) {
	query := `
        UPDATE Task SET 
            short = ?
            , priority = ?
            , date_due = ?
            , date_completed = ?
            , long = ?
            , closing_comment = ?
            , status_id = ?
            , group_id = ?
            , sys_date_updated = datetime('now')
        WHERE id = ?;
    `
	_, err = db.Exec(
		query,
		t.Short,
		t.Priority,
		t.DateDue,
		t.DateCompleted,
		t.Long,
		t.ClosingComment,
		t.Status.Id,
		t.Group.Id,
		t.Id,
	)
	return err
}

func (g Group) Update(db *sql.DB) (err error) {
	_, err = db.Exec(
		"UPDATE Group SET name = ?, sys_date_updated = datetime('now') WHERE id = ?;",
		g.Name,
		g.Id,
	)
	return err
}

func (t Task) Delete(db *sql.DB) (err error) {
	_, err = db.Exec(
		"UPDATE Task SET sys_status = 0 WHERE id = ?;",
		t.Id,
	)
	return err
}

func (g Group) Delete(db *sql.DB) (err error) {
	_, err = db.Exec(
		"UPDATE Group SET sys_status = 0 WHERE id = ?;",
		g.Id,
	)
	return err
}

func CountTask(db *sql.DB, sw int) (count int, err error) {
	query := "SELECT COUNT(*) FROM Task WHERE sys_status = 1 AND status_id in (0, 1, 2);"
	switch sw {
	case A_ALL:
		query = "SELECT COUNT(*) FROM Task WHERE sys_status = 1;"
	case A_CLOSED:
		query = "SELECT COUNT(*) FROM Task WHERE sys_status = 1 AND status_id = 3;"
	case A_OVERDUE:
		query = "SELECT COUNT(*) FROM Task WHERE sys_status = 1 AND date('now') > date_due;"
	}
	err = db.QueryRow(query).Scan(&count)
	return
}

func CountGroup(db *sql.DB, sw int) (vs Values, err error) {
	query := `
        SELECT g.id, g.name, COUNT(*)
        FROM
            Group g
            JOIN Task t ON t.group_id = g.id
        WHERE 
            g.sys_status = 1
            AND t.sys_status = 1
            AND t.status_id IN (0, 1, 2)
        GROUP BY g.id, g.name
        ORDER BY g.id;
    `
	switch sw {
	case A_ALL:
		query = `
            SELECT g.id, g.name, COUNT(*)
            FROM
                Group g
                LEFT OUTER JOIN Task t ON t.group_id = g.id
            WHERE
                g.sys_status = 1
                AND t.sys_status = 1
            GROUP BY g.name
            ORDER BY g.id;
        `
	case A_CLOSED:
		query = `
            SELECT g.id, g.name, COUNT(*)
            FROM
                Group g
                JOIN Task t ON t.group_id = g.id
            WHERE 
                g.sys_status = 1
                AND t.sys_status = 1
                AND t.status_id = 3
            GROUP BY g.id, g.name
            ORDER BY g.id;
        `
	case A_OVERDUE:
		query = `
            SELECT g.id, g.name, COUNT(*)
            FROM
                Group g
                JOIN Task t ON t.group_id = g.id
            WHERE 
                g.sys_status = 1
                AND t.sys_status = 1
                AND t.status_id IN (0, 1, 2)
                AND date('now') > t.date_due
            GROUP BY g.id, g.name
            ORDER BY g.id;
        `
	}

	rows, err := db.Query(query)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var v Value
		if err = rows.Scan(&v.Name, &v.Value); err != nil {
			return
		}
		vs = append(vs, v)
	}
	return
}

func CountPrompt(db *sql.DB) (open, overdue int, err error) {
	err = db.QueryRow("SELECT COUNT(*) FROM Task WHERE sys_status = 1 AND status_id IN (0, 1, 2);").Scan(&open)
	if err != nil {
		open = -1
	}
	err = db.QueryRow("SELECT COUNT(*) FROM Task WHERE sys_status = 1 AND status_id IN (0, 1, 2) AND date('now') > date_due;").Scan(&overdue)
	if err != nil {
		overdue = -1
	}
	return open, overdue, err
}

func List(db *sql.DB, sw int) (tl TaskList, err error) {
	query := `
        SELECT 
            t.id
            , t.short
            , t.priority
            , t.date_due
            , t.date_completed
            , t.long
            , t.closing_comment
            , s.id
            , s.name
            , g.id
            , g.name
        FROM
            Task t
            JOIN Status s ON s.id = t.status_id
            JOIN Group g ON g.id = t.group_id
        WHERE
            t.sys_status = 1
            AND status_id IN (0, 1, 2)
        ORDER BY
            t.priority, t.date_due NULLS LAST, g.name 
    `
	switch sw {
	case A_CLOSED:
		query = `
            SELECT 
                t.id
                , t.short
                , t.priority
                , t.date_due
                , t.date_completed
                , t.long
                , t.closing_comment
                , s.id
                , s.name
                , g.id
                , g.name
            FROM
                Task t
                JOIN Status s ON s.id = t.status_id
                JOIN Group g ON g.id = t.group_id
            WHERE
                t.sys_status = 1
                AND status_id = 3
            ORDER BY
                g.name, t.date_completed DESC
        `
	case A_ALL:
		query = `
            SELECT 
                t.id
                , t.short
                , t.priority
                , t.date_due
                , t.date_completed
                , t.long
                , t.closing_comment
                , s.id
                , s.name
                , g.id
                , g.name
            FROM
                Task t
                JOIN Status s ON s.id = t.status_id
                JOIN Group g ON g.id = t.group_id
            WHERE
                t.sys_status = 1
            ORDER BY
                t.date_completed DESC NULLS FIRST, t.priority DESC, t.due, g.name 
        `
	case A_OVERDUE:
		query = `
            SELECT 
                t.id
                , t.short
                , t.priority
                , t.date_due
                , t.date_completed
                , t.long
                , t.closing_comment
                , s.id
                , s.name
                , g.id
                , g.name
            FROM
                Task t
                JOIN Status s ON s.id = t.status_id
                JOIN Group g ON g.id = t.group_id
            WHERE
                t.sys_status = 1
                AND status_id IN (0, 1, 2)
            ORDER BY
                t.priority, t.date_due NULLS LAST, g.name 
        `
	}

	rows, err := db.Query(query, NullNow())
	if err != nil {
		return tl, err
	}
	defer rows.Close()

	for rows.Next() {
		var t Task
		if err = rows.Scan(
			&t.Id,
			&t.Short,
			&t.Priority,
			&t.DateDue,
			&t.DateCompleted,
			&t.Long,
			&t.ClosingComment,
			&t.Status.Id,
			&t.Status.Name,
			&t.Group.Id,
			&t.Group.Name,
		); err != nil {
			return
		}
		tl = append(tl, t)
	}
	err = rows.Err()
	return tl, err
}
