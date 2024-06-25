package main

import (
	"fmt"
	"log"
    "os"
	"time"
)


func main () {
    
    dateFormat := "2006-01-02"

    cmd, sw, vals, valid := ParseArgs(os.Args[1:], dateFormat)

    if !valid {
        PrintHelp()
        return
    }

    db, err := OpenDB("./")
    if err != nil {
        log.Fatal(err)
    }

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
        fmt.Printf("Task %d has been reopened\n", taskId)
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
            fmt.Printf("Due date: %s\n", t.due.Format(dateFormat))
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
            t_status := "Open"
            if t.done == 0 && t.due.Year() != 1 && t.due.Before(time.Now()) {
                t_status = "Overdue"
            }
            if t.done == 1 {
                t_status = "Closed"
            }
            if t.done == 0 {
                if t.due.Year() == 1 {
                    fmt.Printf("\t%s %d: %s\n", t_status, t.id, t.description)
                } else {
                    fmt.Printf("\t%s %d: %s, due date: %s\n", t_status, t.id, t.description, t.due.Format(dateFormat))
                }
            } else {
                fmt.Printf("\t%s %d: %s, completed: %s\n", t_status, t.id, t.description, t.completed.Format(dateFormat))
            }
        }
        fmt.Printf("End.\n")
    }

    return
}
