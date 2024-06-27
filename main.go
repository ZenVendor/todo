package main

import (
	"fmt"
	"log"
	"time"
)

const VERSION = "0.6.2"

func main () {
    var conf Config
    err := conf.Prepare()
    if err != nil {
        log.Fatal(err)
    }

    cmd, sw, vals, valid := ParseArgs(conf.dateFormat)

    if !valid {
        PrintVersion()
        PrintHelp()
        return
    }
    if cmd == CMD_VERSION {
        PrintVersion()
        return
    }

    db, err := conf.OpenDB()
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    if cmd == CMD_COUNT {
        count, err := Count(db, sw)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("%d", count)
        return
    }

    if cmd == CMD_COMPLETE {
        taskId := vals.ReadValue("id").(int)
        err := Complete(db, taskId)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("Task %d has been completed\n", taskId)
    }

    if cmd == CMD_REOPEN {
        taskId := vals.ReadValue("id").(int)
        err := Reopen(db, taskId)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("Task %d has been reopened\n", taskId)
    }

    if cmd == CMD_DELETE {
        taskId := vals.ReadValue("id").(int)
        err := Delete(db, taskId)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("Task %d has been deleted\n", taskId)
    }

    if cmd == CMD_ADD {
        var t Task
        t.description = vals.ReadValue("description").(string)
        if vals.ReadValue("due") != nil {
            t.due = vals.ReadValue("due").(time.Time)
        }      
        t.created = time.Now()
        t.updated = time.Now()

        if err = t.AddTask(db); err != nil {
            log.Fatal(err)
        }

        fmt.Printf("Added task: %s\n", t.description)
        if t.due.Year() != 1 {
            fmt.Printf("Due date: %s\n", t.due.Format(conf.dateFormat))
        }
    }
   
    if cmd == CMD_LIST {
        count, err := Count(db, sw)
        if err != nil {
            log.Fatal(err)
        }
        tl, err := List(db, sw)
        if err != nil {
            log.Fatal(err)
        }
        var tType string 
        switch sw {
            case SW_OPEN:
                tType = "Open"
            case SW_CLOSED:
                tType = "Closed"
            case SW_ALL:
                tType = "All"
            case SW_OVERDUE:
               tType = "Overdue"
        }
        fmt.Printf("%s tasks: %d\n", tType, count)

        for _, t := range tl {
            tStatus := "Open"
            if t.done == 0 && t.due.Year() != 1 && t.due.Before(time.Now()) {
                tStatus = "Overdue"
            }
            if t.done == 1 {
                tStatus = "Closed"
            }
            if t.done == 0 {
                if t.due.Year() == 1 {
                    fmt.Printf("\t%s %d: %s\n", tStatus, t.id, t.description)
                } else {
                    fmt.Printf("\t%s %d: %s, due date: %s\n", tStatus, t.id, t.description, t.due.Format(conf.dateFormat))
                }
            } else {
                fmt.Printf("\t%s %d: %s, completed: %s\n", tStatus, t.id, t.description, t.completed.Format(conf.dateFormat))
            }
        }
        fmt.Printf("End.\n")
    }
    
    if cmd == CMD_UPDATE {
        var t Task

        taskId := vals.ReadValue("id").(int)
        t, err = Select(db, taskId)
        if err != nil {
            log.Fatal(err)
        }

        if vals.ReadValue("description") != nil {
            t.description = vals.ReadValue("description").(string)
        }
        if vals.ReadValue("due") != nil {
            t.due = vals.ReadValue("due").(time.Time)
        }
        t.updated = time.Now()
        if err = t.Update(db); err != nil {
            log.Fatal(err)
        }
        
        var updString string
        for i, v := range vals {
            if i != 0 {
                updString += ", "
            }
            updString += v.name
        }
        fmt.Printf("Updated %s in task %d\n", updString, t.id)
    }
    return
}
