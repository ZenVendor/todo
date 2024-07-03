package functions

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"slices"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v2"
)

func NullNow() sql.NullTime {
    return sql.NullTime{Time: time.Now(), Valid: true}
}

type Task struct {
    Id int
    Description string
    Priority int
    Done int
    Due sql.NullTime
    Completed sql.NullTime
    Created sql.NullTime
    Updated sql.NullTime
    Group TaskGroup
}
type TaskList []Task

type TaskGroup struct {
    Id int
    Name string
    Created sql.NullTime
    Updated sql.NullTime
}

type Value struct {
    Name string
    Value interface{}
}
type Values []Value

type Config struct {
    DBLocation string   `yaml:"dblocation"`
    DBName string       `yaml:"dbname"`
    DateFormat string   `yaml:"dateformat"`
}

func (conf *Config) Prepare() (configFile string, err error) {
    configDirs := []string{"."}
    configDir := "."

    home, ok := os.LookupEnv("HOME")
    if ok {
        configDirs = append(configDirs, fmt.Sprintf("%s/.config/todo", home))
        configDir = fmt.Sprintf("%s/.config/todo", home)
    }

    for _, cd := range configDirs {
        if _, err = os.Stat(fmt.Sprintf("%s/todo.yml", cd)); err == nil {
            configFile = fmt.Sprintf("%s/todo.yml", cd)
            break
        }
    }

    if configFile == "" { 
        conf.DBLocation = configDir
        conf.DBName = "todo"
        conf.DateFormat = "2006-01-02"
        
        if _, err = os.Stat(configDir); os.IsNotExist(err) {
            if err = os.MkdirAll(configDir, 0700); err != nil {
                return
            }
        }

        writeConf := fmt.Sprintf("dblocation: %s\ndbname: %s\ndateformat: %s", conf.DBLocation, conf.DBName, conf.DateFormat)
        if err = os.WriteFile(configFile, []byte(writeConf), 0700); err != nil {
            return
        }
    }
    return
}

func (conf *Config) ReadConfig(configFile string) (err error) {
    f, err := os.ReadFile(configFile)
    if err != nil {
        return
    }
    err = yaml.Unmarshal(f, &conf)
    return
}

func (conf *Config) OpenDB() (db *sql.DB, err error) {
    db, err = sql.Open("sqlite3", fmt.Sprintf("%s/%s.db", conf.DBLocation, conf.DBName))
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

func (vs Values) ReadValue(name string) interface{} {
    idx := slices.IndexFunc(vs, func(v Value) bool {
        return v.Name == name
    })
    if idx == -1 {
        return nil
    }
    return vs[idx].Value
}




func TableExists(db *sql.DB) bool {
    var rows int
    if err := db.QueryRow("select count(*) from sqlite_schema where type = 'table' and tbl_name in ('tasklist', 'taskgroup');").Scan(&rows); err != nil {
        log.Fatal(err)
    }
    if rows != 2 {
        return false
    }
    return true
}
    
func CreateTable(db *sql.DB) error {
    query := `
        create table tasklist (
            id integer primary key not null,
            description text not null,
            priority integer null,
            group_id integer not null,
            done integer not null,
            due datetime,
            completed datetime,
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
    query := "insert into tasklist (description, priority, group_id, done, due, completed, created, updated) values (?, ?, ?, ?, ?, ?, ?, ?);"
    r, err := db.Exec(query, t.Description, t.Priority, t.Group.Id, t.Done, t.Due, t.Completed, t.Created, t.Updated)
    if err != nil {
        return
    }
    idr, err := r.LastInsertId()
    t.Id = int(idr)
    return
}

func (g *TaskGroup) Insert(db *sql.DB) (err error) {
    query := "insert into taskgroup (name, created, updated) values (?, ?, ?);"
    r, err := db.Exec(query, g.Name, g.Created, g.Updated)
    if err != nil {
        return
    }
    idr, err := r.LastInsertId()
    g.Id = int(idr)
    return
}
func (t *Task) Select(db *sql.DB) (err error) {
    query := `
        select 
            t.id, t.description, t.priority, t.done, t.due, t.completed, t.created, t.updated
            , g.id, g.name, g.created, g.updated
        from 
            tasklist t
            join taskgroup g on g.id = t.group_id
        where 
            t.id = ?;
    `
    err = db.QueryRow(query, t.Id).Scan(&t.Id, &t.Description, &t.Priority, &t.Done, &t.Due, &t.Completed, &t.Created, &t.Updated, &t.Group.Id, &t.Group.Name, &t.Group.Created, &t.Group.Updated)
    return
}

func (g *TaskGroup) Select(db *sql.DB, byName bool) (err error) {
    if byName {
        query := `
            select g.id, g.name, g.created, g.updated
            from taskgroup g
            where g.Name = ?;
        `
        err = db.QueryRow(query, g.Name).Scan(&g.Id, &g.Name, &g.Created, &g.Updated)
    } else {
        query := `
        select g.id, g.name, g.created, g.updated
        from taskgroup g
        where g.id = ?;
        `
        err = db.QueryRow(query, g.Id).Scan(&g.Id, &g.Name, &g.Created, &g.Updated)
    }
    return
}

func (t Task) Update(db *sql.DB) (err error) {
    query := "update tasklist set description = ?, priority = ?, group_id = ?, done = ?, due = ?, completed = ?, updated = ? where id = ?;"
    _, err = db.Exec(query, t.Description, t.Priority, t.Group.Id, t.Done, t.Due, t.Completed, t.Updated, t.Id)
    return err
}

func (g TaskGroup) Update(db *sql.DB) (err error) {
    query := "update taskgroup set name = ?, updated = ? where id = ?;"
    _, err = db.Exec(query, g.Name, g.Updated, g.Id)
    return err
}

func (t Task) Delete(db *sql.DB) (err error) {
    query := "delete from tasklist where id = ?;"
    _, err = db.Exec(query, t.Id)
    return err
}

func (g TaskGroup) Delete(db *sql.DB) (err error) {
    query := "delete from tasklist where id = ?;"
    _, err = db.Exec(query, g.Id)
    return err
}

    
func CountTask(db *sql.DB, sw int) (count int, err error) {
    query := "select count(*) from tasklist where done = 0;"
    switch sw {
    case SW_ALL:
        query = "select count(*) from tasklist;"
    case SW_CLOSED: 
        query = "select count(*) from tasklist where done = 1;"
    case SW_OVERDUE:
        query = "select count(*) from tasklist where done = 0 and due between '2000-01-01' and ?;"
    }
    err = db.QueryRow(query, time.Now()).Scan(&count)
    return
}

func CountGroup(db *sql.DB, sw int) (vs Values, err error) {
    query := `
        select g.name, count(*)
        from
            taskgroup g
            join tasklist t on t.group_id = g.id
        where t.done = 0
        group by g.name
        order by g.id;

    `
    switch sw {
    case SW_ALL:
        query = `
            select g.name, count(*)
            from
                taskgroup g
                join tasklist t on t.group_id = g.id
            group by g.name;
        `
    case SW_CLOSED:
        query = `
            select g.name, count(*)
            from
                taskgroup g
                join tasklist t on t.group_id = g.id
            where t.done = 1
            group by g.name
            order by g.id;
        `
    case SW_OVERDUE:
        query = `
            select g.name, count(*)
            from
                taskgroup g
                join tasklist t on t.group_id = g.id
            where 
                t.done = 0 
                and t.due between '2001-01-01' and ?
            group by g.name
            order by g.name;
        `
    }

    rows, err := db.Query(query, time.Now()) 
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
    query := "select count(*) from tasklist where done = 0;"
    err = db.QueryRow(query).Scan(&open)
    if err != nil {
        open = -1
    }
    query = "select count(*) from tasklist where done = 0 and due < ?;"
    err = db.QueryRow(query, NullNow()).Scan(&overdue)
    if err != nil {
        overdue = -1
    }
    return 
}

func List(db *sql.DB, sw int) (tl TaskList, err error) {
    query := `
        select 
            t.id, t.description, t.priority, t.done, t.due, t.completed, t.created, t.updated
            , g.id, g.name, g.created, g.updated
        from
            tasklist t
            join taskgroup g on g.id = t.group_id
        where
            t.done = 0
        order by
            t.priority desc, t.due nulls last, g.name 
    `
    switch sw {
    case SW_CLOSED:
        query = `
            select 
                t.id, t.description, t.priority, t.done, t.due, t.completed, t.created, t.updated
                , g.id, g.name, g.created, g.updated
            from
                tasklist t
                join taskgroup g on g.id = t.group_id
            where
                t.done = 1
            order by
                g.name, t.completed desc
        `
    case SW_ALL:
        query = `
            select 
                t.id, t.description, t.priority, t.done, t.due, t.completed, t.created, t.updated
                , g.id, g.name, g.created, g.updated
            from
                tasklist t
                join taskgroup g on g.id = t.group_id
            order by
                t.done, t.priority desc, t.due, g.name 
        `
    case SW_OVERDUE:
        query = `
            select 
                t.id, t.description, t.priority, t.done, t.due, t.completed, t.created, t.updated
                , g.id, g.name, g.created, g.updated
            from
                tasklist t
                join taskgroup g on g.id = t.group_id
            where
                t.done = 0
                and t.due < ?
            order by
                t.priority desc, t.due, g.name 
        `
    }

    rows, err := db.Query(query, NullNow()) 
    if err != nil {
        return 
    }
    defer rows.Close()
    
    for rows.Next() {
        var t Task
        if err = rows.Scan(&t.Id, &t.Description, &t.Priority, &t.Done, &t.Due, &t.Completed, &t.Created, &t.Updated,
                &t.Group.Id, &t.Group.Name, &t.Group.Created, &t.Group.Updated); 
            err != nil {
            return
        }
        tl = append(tl, t)
    }
    err = rows.Err()
    return 
}
