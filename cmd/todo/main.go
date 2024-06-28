package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ZenVendor/todo/internal/functions"
)

const VERSION = "0.7.1"

func main () {
    var conf functions.Config
    configFile, err := conf.Prepare()
    if err != nil {
        log.Fatal(err)
    }
    if err = conf.ReadConfig(configFile); err != nil {
        log.Fatal(err)
    }

    cmd, sw, vals, valid := functions.ParseArgs(conf.DateFormat)

    if !valid {
        PrintVersion()
        PrintHelp()
        return
    }
    if cmd == functions.CMD_VERSION {
        PrintVersion()
        return
    }

    db, err := conf.OpenDB()
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    if cmd == functions.CMD_COUNT {
        count, err := functions.Count(db, sw)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("%d", count)
        return
    }

    if cmd == functions.CMD_COMPLETE {
        taskId := vals.ReadValue("id").(int)
        err := functions.Complete(db, taskId)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("Task %d has been completed\n", taskId)
    }

    if cmd == functions.CMD_REOPEN {
        taskId := vals.ReadValue("id").(int)
        err := functions.Reopen(db, taskId)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("Task %d has been reopened\n", taskId)
    }

    if cmd == functions.CMD_DELETE {
        taskId := vals.ReadValue("id").(int)
        err := functions.Delete(db, taskId)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("Task %d has been deleted\n", taskId)
    }

    if cmd == functions.CMD_ADD {
        var t functions.Task
        t.Description = vals.ReadValue("description").(string)
        if vals.ReadValue("due") != nil {
            t.Due = vals.ReadValue("due").(time.Time)
        }      
        t.Created = time.Now()
        t.Updated = time.Now()

        if err = t.AddTask(db); err != nil {
            log.Fatal(err)
        }

        fmt.Printf("Added task: %s\n", t.Description)
        if t.Due.Year() != 1 {
            fmt.Printf("Due date: %s\n", t.Due.Format(conf.DateFormat))
        }
    }
   
    if cmd == functions.CMD_LIST {
        count, err := functions.Count(db, sw)
        if err != nil {
            log.Fatal(err)
        }
        tl, err := functions.List(db, sw)
        if err != nil {
            log.Fatal(err)
        }
        var tType string 
        switch sw {
            case functions.SW_OPEN:
                tType = "Open"
            case functions.SW_CLOSED:
                tType = "Closed"
            case functions.SW_ALL:
                tType = "All"
            case functions.SW_OVERDUE:
               tType = "Overdue"
        }
        fmt.Printf("%s tasks: %d\n", tType, count)

        for _, t := range tl {
            tStatus := "Open"
            if t.Done == 0 && t.Due.Year() != 1 && t.Due.Before(time.Now()) {
                tStatus = "Overdue"
            }
            if t.Done == 1 {
                tStatus = "Closed"
            }
            if t.Done == 0 {
                if t.Due.Year() == 1 {
                    fmt.Printf("\t%s %d: %s\n", tStatus, t.Id, t.Description)
                } else {
                    fmt.Printf("\t%s %d: %s, due date: %s\n", tStatus, t.Id, t.Description, t.Due.Format(conf.DateFormat))
                }
            } else {
                fmt.Printf("\t%s %d: %s, completed: %s\n", tStatus, t.Id, t.Description, t.Completed.Format(conf.DateFormat))
            }
        }
        fmt.Printf("End.\n")
    }
    
    if cmd == functions.CMD_UPDATE {
        var t functions.Task

        taskId := vals.ReadValue("id").(int)
        t, err = functions.Select(db, taskId)
        if err != nil {
            log.Fatal(err)
        }

        if vals.ReadValue("description") != nil {
            t.Description = vals.ReadValue("description").(string)
        }
        if vals.ReadValue("due") != nil {
            t.Due = vals.ReadValue("due").(time.Time)
        }
        t.Updated = time.Now()
        if err = t.Update(db); err != nil {
            log.Fatal(err)
        }
        
        var updString string
        for i, v := range vals {
            if i != 0 {
                updString += ", "
            }
            updString += v.Name
        }
        fmt.Printf("Updated %s in task %d\n", updString, t.Id)
    }
    return
}

func PrintVersion() {
    fmt.Printf("TODO CLI\tversion: %s\n", VERSION)
}

func PrintHelp() {
    helpString := `
Usage: 
    todo [command] [id] [option] [argument]

Without arguments defaults to listing active tasks.
Frequently used commands have single-letter aliases.
In ADD command, description is required and must be provided first, due date can optionally be provided second.
In UPDATE the switches and their following values can be provided in any order.
In commands that require it, task ID must follow the command.
Due date must match the format, default is "YYYY-MM-DD". 

    
    help | h | --help | -h                      display this help

    version | v | --version | -v                display program version

    add | a [description] [due]                 optional due date format 2006-01-02

    count                                       defaults to active tasks
        --completed | -c
        --overdue | -o
        --all | -a

    list | l                                    defaults to active tasks
        --completed | -c
        --overdue | -o
        --all | -a
        
    update | u [id]                             update description, due date, or both. invalid date value removes due date
        --desc [description] 
        --due [date]

    complete | c [task_id]                      set task completed

    reopen | open [task_id]                     reopen completed task

    delete | del [task_id]                      delete task

Examples:
    todo
    todo a "New task"
    todo add "New task" "2024-08-13"
    todo list --all
    todo l -o
    todo count -c
    todo update 15 --due "2024-08-13"
    todo u 10 --due -
    todo c 12
    todo reopen 3
    todo del 5

`
    fmt.Println(helpString)
}
