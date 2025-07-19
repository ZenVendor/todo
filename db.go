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
	Summary        string
	Priority       int
	DateDue        sql.NullTime
	DateCompleted  sql.NullTime
	Description    string
	ClosingComment string
	Status         *Status
	Group          *Group
	Parent         *Task
	DateCreated    sql.NullTime
	DateUpdated    sql.NullTime
	SysStatus      int
}
type TaskList []Task

type Group struct {
	Id     int
	Name   string
	Counts *Counts
}
type Status struct {
	Id     int
	Name   string
	Counts *Counts
}

type Counts struct {
	All        int
	New        int
	InProgress int
	OnHold     int
	Completed  int
	Open        int
	Overdue    int
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
        create table Tasks (
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

func (t *Task) Add(db *sql.DB) (err error) {
	result, err := db.Exec(`
        INSERT INTO Tasks (
            summary
            , priority
            , date_due
            , description
            , group_id
            , parent_id
        ) VALUES (?, ?, ?, ?, ?, ?);
        `,
		t.Summary,
		t.Priority,
		t.DateDue,
		t.Description,
		t.Group.Id,
		t.Parent.Id,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	t.Id = int(id)
	return err
}

func (g *Group) Add(db *sql.DB) (err error) {
	result, err := db.Exec(
		"INSERT INTO Groups (group_name) VALUES (?);",
		g.Name,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	g.Id = int(id)
	return err
}

func (t *Task) GetById(db *sql.DB) (err error) {
	query := `
        SELECT 
            id 
            , summary
            , priority
            , date_due
            , date_completed
            , description
            , closing_comment
            , status_id
            , status_name
            , group_id
            , group_name
            , parent_id
            , sys_created
            , sys_updated
            , sys_status
        FROM task_list_all;
        WHERE id = ?;
    `
	err = db.QueryRow(query, t.Id).Scan(
		&t.Id,
		&t.Summary,
		&t.Priority,
		&t.DateDue,
		&t.DateCompleted,
		&t.Description,
		&t.ClosingComment,
		&t.Status.Id,
		&t.Status.Name,
		&t.Group.Id,
		&t.Group.Name,
		&t.Parent.Id,
		&t.DateCreated,
		&t.DateUpdated,
		&t.SysStatus,
	)
	return err
}

func (g *Group) GetById(db *sql.DB) (err error) {
	err = db.QueryRow(
		"SELECT id, group_name FROM Groups WHERE id = ?;",
		g.Id,
	).Scan(&g.Id, &g.Name)
	return err
}

func (g *Group) GetByName(db *sql.DB) (err error) {
	err = db.QueryRow(
		"SELECT id, group_name FROM Groups WHERE group_name = ?;",
		g.Name,
	).Scan(&g.Id, &g.Name)
	return err
}

func (t Task) Update(db *sql.DB) (err error) {
	query := `
        UPDATE Tasks SET 
            summary = ?
            , priority = ?
            , date_due = ?
            , date_completed = ?
            , description = ?
            , closing_comment = ?
            , status_id = ?
            , group_id = ?
            , parent_id = ?
            , sys_updated = current_timestamp
        WHERE id = ?;
    `
	_, err = db.Exec(
		query,
		t.Summary,
		t.Priority,
		t.DateDue,
		t.DateCompleted,
		t.Description,
		t.ClosingComment,
		t.Status.Id,
		t.Group.Id,
		t.Parent.Id,
		t.Id,
	)
	return err
}

func (g Group) Update(db *sql.DB) (err error) {
	_, err = db.Exec(
		"UPDATE Groups SET group_name = ?, sys_updated = current_timestamp WHERE id = ?;",
		g.Name,
		g.Id,
	)
	return err
}

func (t Task) Delete(db *sql.DB) (err error) {
	_, err = db.Exec(
		"UPDATE Tasks SET sys_status = 0, sys_updated = current_timestamp WHERE id = ?;",
		t.Id,
	)
	return err
}

func (g Group) Delete(db *sql.DB) (err error) {
	_, err = db.Exec(
		"UPDATE Groups SET sys_status = 0, sys_updated = current_timestamp WHERE id = ?;",
		g.Id,
	)
	return err
}

func ListTasks(db *sql.DB) (tl TaskList, err error) {
	query := `
        SELECT 
            id 
            , summary
            , priority
            , date_due
            , date_completed
            , description
            , closing_comment
            , status_id
            , status_name
            , group_id
            , group_name
            , parent_id
            , sys_created
            , sys_updated
            , sys_status
        FROM task_list_all;
    `
	rows, err := db.Query(query, NullNow())
	if err != nil {
		return tl, err
	}
	defer rows.Close()

	for rows.Next() {
		var t Task
		if err = rows.Scan(
			&t.Id,
			&t.Summary,
			&t.Priority,
			&t.DateDue,
			&t.DateCompleted,
			&t.Description,
			&t.ClosingComment,
			&t.Status.Id,
			&t.Status.Name,
			&t.Group.Id,
			&t.Group.Name,
			&t.Parent.Id,
			&t.DateCreated,
			&t.DateUpdated,
			&t.SysStatus,
		); err != nil {
			return tl, err
		}
		tl = append(tl, t)
	}
	err = rows.Err()
	return tl, err
}

func (c *Counts) GetCounts(db *sql.DB) (err error) {
	query := `
        SELECT 
            count_all
            , count_new
            , count_in_progress
            , count_on_hold
            , count_completed
            , count_open
            , count_overdue
        FROM task_counts
    `
	err = db.QueryRow(query).Scan(
		&c.All,
		&c.New,
		&c.InProgress,
		&c.OnHold,
		&c.Completed,
		&c.Open,
		&c.Overdue,
	)
	return err
}

func (g *Group) GetCounts(db *sql.DB) (err error) {
	query := `
        SELECT 
            count_all
            , count_new
            , count_in_progress
            , count_on_hold
            , count_completed
            , count_open
            , count_overdue
        FROM group_counts
        WHERE id = ?;
    `
	err = db.QueryRow(query, g.Id).Scan(
		&g.Counts.All,
		&g.Counts.New,
		&g.Counts.InProgress,
		&g.Counts.OnHold,
		&g.Counts.Completed,
		&g.Counts.Open,
		&g.Counts.Overdue,
	)
	return err
}

func (s *Status) GetCounts(db *sql.DB) (err error) {
	query := `
        SELECT 
            count_all
            , count_new
            , count_in_progress
            , count_on_hold
            , count_completed
            , count_open
            , count_overdue
        FROM status_counts
        WHERE id = ?;
    `
	err = db.QueryRow(query, s.Id).Scan(
		&s.Counts.All,
		&s.Counts.New,
		&s.Counts.InProgress,
		&s.Counts.OnHold,
		&s.Counts.Completed,
		&s.Counts.Open,
		&s.Counts.Overdue,
	)
	return err
}
