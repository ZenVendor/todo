package functions

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v2"
)

const CMD_NONE = 0
const CMD_ADD = 1
const CMD_LIST = 2
const CMD_COUNT = 3
const CMD_UPDATE = 4
const CMD_DELETE = 5
const CMD_REOPEN = 6
const CMD_COMPLETE = 7
const CMD_VERSION = 8

const SW_NONE = 0
const SW_OPEN = 1
const SW_CLOSED = 2
const SW_ALL = 3
const SW_OVERDUE = 4

type Task struct {
    Id int
    Description string
    Done int
    Due time.Time
    Created time.Time
    Completed time.Time
    Updated time.Time
}
type TaskList []Task

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

func ParseArgs(dateFormat string) (cmd, sw int, values Values, valid bool) {
    args := os.Args[1:]
    valid = false

    if len(args) == 0 {
        cmd = CMD_LIST
        sw = SW_OPEN
        valid = true
    } else {
        if slices.Contains([]string{"help", "h", "--help", "-h"}, args[0]) {
            return CMD_NONE, SW_NONE, []Value{}, false 
        }
        if slices.Contains([]string{"version", "v", "--version", "-v"}, args[0]) {
            return CMD_VERSION, SW_NONE, []Value{}, true 
        }
        if slices.Contains([]string{"count"}, args[0]) && slices.Contains([]int{1, 2}, len(args)) {
            cmd = CMD_COUNT
            valid = true
        }
        if slices.Contains([]string{"list", "l"}, args[0]) && slices.Contains([]int{1, 2}, len(args)) {
            cmd = CMD_LIST
            valid = true
        }
        if slices.Contains([]string{"add", "a"}, args[0]) && slices.Contains([]int{2, 3}, len(args)) {
            cmd = CMD_ADD
            valid = true
        }
        if slices.Contains([]string{"complete", "c"}, args[0]) && slices.Contains([]int{2}, len(args)) {
            cmd = CMD_COMPLETE
            valid = true
        }
        if slices.Contains([]string{"reopen", "open"}, args[0]) && slices.Contains([]int{2}, len(args)) {
            cmd = CMD_REOPEN
            valid = true
        }
        if slices.Contains([]string{"delete", "del"}, args[0]) && slices.Contains([]int{2}, len(args)) {
            cmd = CMD_DELETE
            valid = true
        }
        if slices.Contains([]string{"update", "u"}, args[0]) && slices.Contains([]int{4, 6}, len(args)) {
            cmd = CMD_UPDATE
            valid = true
        }
        if !valid {
            return CMD_NONE, SW_NONE, []Value{}, false 
        }

        args = args[1:]
        valid = false

        if len(args) == 0 && slices.Contains([]int{CMD_COUNT, CMD_LIST}, cmd) {
            sw = SW_OPEN
            valid = true
        }
        if len(args) == 1 {
            if slices.Contains([]int{CMD_COUNT, CMD_LIST}, cmd) {
                if slices.Contains([]string{"--all", "-a"}, args[0]) {
                    sw = SW_ALL
                    valid = true
                }
                if slices.Contains([]string{"--completed", "--closed", "-c"}, args[0]) {
                    sw = SW_CLOSED
                    valid = true
                }
                if slices.Contains([]string{"--overdue", "--od", "-o"}, args[0]) {
                    sw = SW_OVERDUE
                    valid = true
                }
            }
            if slices.Contains([]int{CMD_DELETE, CMD_REOPEN, CMD_COMPLETE}, cmd) {
                taskId, err := strconv.Atoi(args[0])
                if err != nil {
                    valid = false
                } else {
                    values = append(values, Value{"id", taskId})
                    valid = true
                }
            }
            if slices.Contains([]int{CMD_ADD}, cmd) {
                values = append(values, Value{"description", args[0]})
                valid = true
            }
        }
        if len(args) == 2 {
            if slices.Contains([]int{CMD_ADD}, cmd) {
                dd, err := time.Parse(dateFormat, args[1])
                if err != nil {
                    valid = false
                } else {
                    values = append(values, Value{"description", args[0]})
                    values = append(values, Value{"due", dd})
                    valid = true
                }
            }
        }
        if len(args) == 3 {
            if slices.Contains([]int{CMD_UPDATE}, cmd) {
                taskId, err := strconv.Atoi(args[0])
                if err != nil {
                    valid = false
                } else {
                    if "--desc" == args[1] {
                        values = append(values, Value{"id", taskId})
                        values = append(values, Value{"description", args[2]})
                        valid = true
                    }
                    if "--due" == args[1] {
                        dd, err := time.Parse(dateFormat, args[2])
                        if err != nil {
                            dd, _ = time.Parse(dateFormat, "0001-01-01")
                        }
                        values = append(values, Value{"id", taskId})
                        values = append(values, Value{"due", dd})
                        valid = true
                    }
                }
            }
        }
        if len(args) == 5 {
            if slices.Contains([]int{CMD_UPDATE}, cmd) {
                taskId, err := strconv.Atoi(args[0])
                if err != nil {
                    valid = false 
                } else {
                    var dd time.Time
                    if "--desc" == args[1] && "--due" == args[3] {
                        ok, _ := regexp.MatchString("[0-9]{4}-[0-9]{2}-[0-9]{2}", args[4])
                        if ok {
                            dd, _ = time.Parse(dateFormat, args[4])
                        }
                    }
                    if "--desc" == args[3] && "--due" == args[1] {
                        ok, _ := regexp.MatchString("[0-9]{4}-[0-9]{2}-[0-9]{2}", args[2])
                        if ok {
                            dd, _ = time.Parse(dateFormat, args[2])
                        }
                    }
                    values = append(values, Value{"id", taskId})
                    values = append(values, Value{"description", args[4]})
                    values = append(values, Value{"due", dd})
                    valid = true
                }
            }
        }
    }
    return 
}

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
            due datetime,
            created datetime not null,
            completed datetime,
            updated datetime not null
        );
    `
    log.Println("Creating table.")
    _, err := db.Exec(query)
    return err
}

func (t Task) AddTask(db *sql.DB) (err error) {
    query := "insert into tasklist (description, done, due, created, updated) values (?, ?, ?, ?, ?);"
    _, err = db.Exec(query, t.Description, t.Done, t.Due, t.Created, t.Updated)
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
        query = fmt.Sprintf("select count(*) from tasklist where done = 0 and due between '2000-01-01' and '%s';", time.Now())
    }
    err = db.QueryRow(query).Scan(&count)
    return
}

func CountPrompt(db *sql.DB) (open, overdue int, err error) {
    query := "select count(*) from tasklist where done = 0;"
    err = db.QueryRow(query).Scan(&open)
    if err != nil {
        open = -1
    }

    query = fmt.Sprintf("select count(*) from tasklist where done = 0 and due between '2000-01-01' and '%s';", time.Now().Format("2006-01-02"))
    err = db.QueryRow(query).Scan(&overdue)
    if err != nil {
        overdue = -1
    }
    return 
}

func List(db *sql.DB, sw int) (tl TaskList, err error) {
    var query string
    switch sw {
    case SW_OPEN:
        query = "select * from tasklist where done = 0 order by due asc nulls last, created ;"
    case SW_CLOSED:
        query = "select * from tasklist where done = 1 order by completed desc;"
    case SW_ALL:
        query = "select * from tasklist order by done, completed desc, due asc nulls last, created;"
    case SW_OVERDUE:
        query = fmt.Sprintf("select * from tasklist where done = 0 and due between '2000-01-01' and '%s';", time.Now().Format("2006-01-02"))
    }
    
    rows, err := db.Query(query) 
    if err != nil {
        return 
    }
    defer rows.Close()
    
    for rows.Next() {
        var t Task
        var due, comp sql.NullString
        if err = rows.Scan(&t.Id, &t.Description, &t.Done, &due, &t.Created, &comp, &t.Updated); err != nil {
            return
        }
        if due.Valid {
            duedate, err := time.Parse(time.RFC3339, due.String)
            if err != nil {
                return tl, err
            }
            t.Due = duedate 
        }
        if comp.Valid {
            completed, err := time.Parse(time.RFC3339, comp.String)
            if err != nil {
                return tl, err
            }
            t.Completed = completed
        }
        tl = append(tl, t)
    }
    if err = rows.Err(); err != nil {
        return
    }
    return 
}

func Complete(db *sql.DB, taskId int) (err error) {
    query := "update tasklist set done = 1, completed = ?, updated = ? where id = ?;"
    now := time.Now().Format(time.RFC3339)
    _, err = db.Exec(query, now, now, taskId)
    return err
}

func Reopen(db *sql.DB, taskId int) (err error) {
    query := "update tasklist set done = 0, completed = null, updated = ? where id = ?;"
    now := time.Now().Format(time.RFC3339)
    _, err = db.Exec(query, now, taskId)
    return err
}

func Delete(db *sql.DB, taskId int) (err error) {
    query := "delete from tasklist where id = ?;"
    _, err = db.Exec(query, taskId)
    return err
}

func Select(db *sql.DB, taskId int) (t Task, err error) {
    query := "select id, description, due from tasklist where id = ?;"
    err = db.QueryRow(query, taskId).Scan(&t.Id, &t.Description, &t.Due)
    return
}
func (t Task) Update(db *sql.DB) (err error) {
    query := "update tasklist set description = ?, due = ?, updated = ? where id = ?;"
    _, err = db.Exec(query, t.Description, t.Due, t.Updated, t.Id)
    return err
}
