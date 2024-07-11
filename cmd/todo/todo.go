package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"slices"
	"time"

	f "github.com/ZenVendor/todo/internal/functions"
)

const VERSION = "0.9.0"

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

    cmd, sw, vals, err := f.ParseArgs(f.CMD_LIST, f.SW_OPEN)

    if err != nil || cmd == f.CMD_HELP {
        PrintVersion()
        PrintHelp()
        log.Fatal(err)
    }
    if cmd == f.CMD_VERSION {
        PrintVersion()
        return
    }
    if cmd == f.CMD_PREPARE || cmd == f.CMD_RESET {
        local := false
        if sw == f.SW_LOCAL {
            local = true
        }
        reset := false
        if cmd == f.CMD_RESET {
            reset = true
        }
        conf.Prepare(local, reset)
        return
    }

    conf.ReadConfig()
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
        addGroup := false

        t.Id = vals.ReadValue(cmd, f.SW_DEFAULT, conf.DateFormat).(int)
        if err = t.Select(db); err != nil {
            log.Fatal(err)
        }
        switch cmd {
        case f.CMD_COMPLETE:
            t.Done = 1
            t.Completed = f.NullNow()
        case f.CMD_REOPEN:
            t.Done = 0
            t.Completed.Valid = false
        case f.CMD_UPDATE:
            if vals.ValueIsSet(f.SW_DESCRIPTION) {
                t.Description = vals.ReadValue(cmd, f.SW_DESCRIPTION, conf.DateFormat).(string)
            }
            if vals.ValueIsSet(f.SW_PRIORITY) {
                t.Priority = vals.ReadValue(cmd, f.SW_PRIORITY, conf.DateFormat).(int)
            }
            if vals.ValueIsSet(f.SW_DUE) {
                t.Due = vals.ReadValue(cmd, f.SW_DUE, conf.DateFormat).(sql.NullTime)
            }
            if vals.ValueIsSet(f.SW_GROUP) {
                t.Group.Name = vals.ReadValue(cmd, f.SW_GROUP, conf.DateFormat).(string)
                if err = t.Group.Select(db, true); err == sql.ErrNoRows {
                    addGroup = true
                } 
            }
        }
        t.Updated = f.NullNow()

        if cmd == f.CMD_DELETE {
            err = t.Delete(db)
        } else {
            if addGroup {
                t.Group.Created = f.NullNow()
                t.Group.Updated = f.NullNow()
                if err = t.Group.Insert(db); err != nil {
                    log.Fatal(err)
                }
            }
            err = t.Update(db)
        }
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("Task %d has been %sd.\n", t.Id, f.MapCommandDescription()[cmd])
    }
        
    if cmd == f.CMD_ADD {
        var t f.Task
        //t.Id
        t.Description = vals.ReadValue(cmd, f.SW_DEFAULT, conf.DateFormat).(string)
        t.Priority = 1
        if vals.ValueIsSet(f.SW_PRIORITY) {
            t.Priority = vals.ReadValue(cmd, f.SW_PRIORITY, conf.DateFormat).(int)
        } 
        t.Done = 0
        //t.Due
        if vals.ValueIsSet(f.SW_DUE) {
            t.Due = vals.ReadValue(cmd, f.SW_DUE, conf.DateFormat).(sql.NullTime)
        }
        //t.Completed
        t.Created = f.NullNow()
        t.Updated = f.NullNow()
        t.Group.Name = "Default"
        if vals.ValueIsSet(f.SW_GROUP) {
            t.Group.Name = vals.ReadValue(cmd, f.SW_GROUP, conf.DateFormat).(string)
        }
        if err = t.Group.Select(db, true); err == sql.ErrNoRows {
            t.Group.Created = f.NullNow()
            t.Group.Updated = f.NullNow()
            if err = t.Group.Insert(db); err != nil {
                log.Fatal(err)
            }
        }
        if err = t.Insert(db); err != nil {
            log.Fatal(err)
        }

        fmt.Printf("Added task %d: %s\n", t.Id, t.Description)
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
        fmt.Printf("%s tasks: %d\n", f.MapArgumentDescription()[sw], count)

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
