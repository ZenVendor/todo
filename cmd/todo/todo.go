package main

import (
	_ "embed"
	"fmt"
	"log"
	"slices"
	"time"

	f "github.com/ZenVendor/todo/internal/functions"
)

const VERSION = "0.8.0"

//go:embed help.txt
var helpString string

func PrintVersion() {
    fmt.Printf("TODO CLI\tversion: %s\n", VERSION)
}
func PrintHelp() {
    fmt.Println(helpString)
}

func main () {
    var conf f.Config
    configFile, err := conf.Prepare()
    if err != nil {
        log.Fatal(err)
    }
    if err = conf.ReadConfig(configFile); err != nil {
        log.Fatal(err)
    }

    cmd, sw, vals, valid := f.ParseArgs(conf.DateFormat)

    if !valid {
        PrintVersion()
        PrintHelp()
        return
    }
    if cmd == f.CMD_VERSION {
        PrintVersion()
        return
    }

    db, err := conf.OpenDB()
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    if cmd == f.CMD_COUNT {
        count, err := f.CountTask(db, sw)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("%d", count)
        return
    }

    if slices.Contains([]int{f.CMD_COMPLETE, f.CMD_DELETE, f.CMD_REOPEN, f.CMD_UPDATE}, cmd) {
        var t f.Task
        t.Id = vals.ReadValue("id").(int)
        err = t.Select(db)
        if err != nil {
            log.Fatal(err)
        }
        var action string

        switch cmd {
        case f.CMD_COMPLETE:
            t.Done = 1
            t.Completed.Time = time.Now()
            t.Completed.Valid = true
            action = "completed"
        case f.CMD_DELETE:
            action = "deleted"
        case f.CMD_REOPEN:
            t.Done = 0
            t.Completed.Valid = false
            action = "reopened"
        case f.CMD_UPDATE:
            if vals.ReadValue("description") != nil {
                t.Description = vals.ReadValue("description").(string)
            }
            if vals.ReadValue("priority") != nil {
                t.Priority = vals.ReadValue("priority").(int)
            }
            if vals.ReadValue("due") != nil {
                t.Due.Time = vals.ReadValue("due").(time.Time)
                t.Due.Valid = true
            }
            action = "updated"
        }
        t.Updated.Time = time.Now()
        t.Updated.Valid = true

        if cmd == f.CMD_DELETE {
            err = t.Delete(db)
        } else {
            err = t.Update(db)
        }
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("Task %d has been %s.\n", t.Id, action)
    }
        
    if cmd == f.CMD_ADD {
        var t f.Task
        t.Description = vals.ReadValue("description").(string)
        t.Priority = 1
        if vals.ReadValue("priority") != nil {
            t.Priority = vals.ReadValue("priority").(int)
        }      
        t.GroupId = 1
        if vals.ReadValue("group") != nil {
            t.Group.Name = vals.ReadValue("group").(string)
            err = t.Group.Select(db, true)
        }      
        if vals.ReadValue("due") != nil {
            t.Due.Time = vals.ReadValue("due").(time.Time)
            t.Due.Valid = true
        }      
        t.Created.Time = time.Now()
        t.Created.Valid = true
        t.Updated.Time = time.Now()
        t.Updated.Valid = true

        if err = t.Insert(db); err != nil {
            log.Fatal(err)
        }

        fmt.Printf("Added task: %s\n", t.Description)
        if t.Due.Valid {
            fmt.Printf("Due date: %s\n", t.Due.Time.Format(conf.DateFormat))
        }
    }
   
    if cmd == f.CMD_LIST {
        count, err := f.CountTask(db, sw)
        if err != nil {
            log.Fatal(err)
        }
        tl, err := f.List(db, sw)
        if err != nil {
            log.Fatal(err)
        }
        var tType string 
        switch sw {
            case f.SW_OPEN:
                tType = "Open"
            case f.SW_CLOSED:
                tType = "Closed"
            case f.SW_ALL:
                tType = "All"
            case f.SW_OVERDUE:
               tType = "Overdue"
        }
        fmt.Printf("%s tasks: %d\n", tType, count)

        for _, t := range tl {
            tStatus := "Open"
            if t.Done == 0 && t.Due.Valid && t.Due.Time.Before(time.Now()) {
                tStatus = "Overdue"
            }
            if t.Done == 1 {
                tStatus = "Closed"
            }
            if t.Done == 0 {
                if !t.Due.Valid {
                    fmt.Printf("\t%s :: %s %d: %s\n", t.Group.Name, tStatus, t.Id, t.Description)
                } else {
                    fmt.Printf("\t%s :: %s %d: %s, due date: %s\n", t.Group.Name, tStatus, t.Id, t.Description, t.Due.Time.Format(conf.DateFormat))
                }
            } else {
                fmt.Printf("\t%s :: %s %d: %s, completed: %s\n", t.Group.Name, tStatus, t.Id, t.Description, t.Completed.Time.Format(conf.DateFormat))
            }
        }
    }
}
