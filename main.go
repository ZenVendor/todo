package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"time"
)

const CMD_COUNT = 0
const CMD_ADD = 1
const CMD_LIST = 2
const CMD_COMPLETE = 3
const CMD_REOPEN = 4
const CMD_HELP = 5

const SW_OPEN = 0
const SW_CLOSED = 1
const SW_OVERDUE = 2
const SW_DUE = 3
const SW_ALL = 4
	
func main () {

    dbLocation := "./"
    
    args := os.Args[1:]
    command := CMD_LIST
    sw := SW_OPEN
    correct := true
    var err error

    if len(args) != 0 {
        correct = false
        if slices.Contains([]string{"count", "c"}, args[0]) {
            command = CMD_COUNT
            correct = true
        }
        if slices.Contains([]string{"add", "a"}, args[0]) && len(args) > 1 {
            command = CMD_ADD
            correct = true
        }
        if slices.Contains([]string{"list", "l"}, args[0]) {
            command = CMD_LIST
            correct = true
        }
        if slices.Contains([]string{"complete", "close", "do", "d"}, args[0]) && len(args) == 2 {
            command = CMD_COMPLETE
            correct = true
        }
        if slices.Contains([]string{"reopen", "open", "undo", "u"}, args[0]) && len(args) == 2 {
            command = CMD_REOPEN
            correct = true
        }
        if slices.Contains([]string{"help", "--help", "-h", "h"}, args[0]) && len(args) == 1 {
            command = CMD_HELP
            correct = true
        }
    }
    
    if command == CMD_HELP {
        PrintHelp()
        return
    }
    if !correct {
        fmt.Printf("Incorrect command: %s\n", args[0])
        PrintHelp()
        return
    }

    sw_arg := -1
    num_arg := -1
    if len(args) > 1 {
        args = os.Args[2:]

        correct = false
        if slices.Contains([]int{CMD_COUNT, CMD_LIST}, command) && len(args) == 1 {
            if slices.Contains([]string{"--all", "-a"}, args[0]) {
                sw = SW_ALL
                correct = true
            }
            if slices.Contains([]string{"--completed", "--closed", "-c"}, args[0]) {
                sw = SW_CLOSED
                correct = true
            }
            if slices.Contains([]string{"--overdue", "-o"}, args[0]) {
                sw = SW_OVERDUE
                correct = true
            }
            if slices.Contains([]string{"--duedate", "--due", "-d"}, args[0]) {
                sw = SW_DUE
                correct = true
            }
        }
        if slices.Contains([]int{CMD_COMPLETE, CMD_REOPEN}, command) && len(args) == 1 {
            correct = true
            num_arg, err = strconv.Atoi(args[0])
            if err != nil {
                correct = false
            }
        }
        if command == CMD_ADD && len(args) == 1 {
                correct = true
        }
        if command == CMD_ADD && len(args) == 3 {
            for i, a := range args {
                if slices.Contains([]string{"--duedate", "--due", "-d"}, a) {
                    sw = SW_DUE
                    sw_arg = i + 1
                }
            }
            if sw_arg > 0 {
                correct = true
            }
        }
    }

    if !correct {
        fmt.Printf("Incorrect arguments: %v\n", args)
        PrintHelp()
        return
    }

    db, err := OpenDB(dbLocation)
    if err != nil {
        log.Fatal(err)
    }

    if command == CMD_COUNT {
        count, err := Count(db, sw)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("%d", count)
        return
    }

    if command == CMD_COMPLETE {
        err := Complete(db, num_arg)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("Task %d has been completed\n", num_arg)
    }

    if command == CMD_REOPEN {
        err := Reopen(db, num_arg)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("Task %d has been reopened\n", num_arg)
    }

    if command == CMD_ADD {
        var t Task

        t.description = args[0]
        t.created = time.Now()
        t.updated = time.Now()

        if len(args) == 3 {
            if sw == SW_DUE {   
                duedate, err := time.Parse("2006-01-02", args[sw_arg])
                if err != nil {
                    log.Fatal(err)
                }
                t.duedate = duedate
                if sw_arg == 1 {
                    t.description = args[2]
                }
            }
        }

        if err = t.AddTask(db); err != nil {
            log.Fatal(err)
        }
        fmt.Printf("Added task: %s\n", t.description)
        if sw == SW_DUE {
            fmt.Printf("Due date: %s\n", t.duedate.Format("2006-01-02"))
        }
    }
   
    if command == CMD_LIST {
        count, err := Count(db, sw)
        if err != nil {
            log.Fatal(err)
        }
        tl, err := List(db, sw)
        if err != nil {
            log.Fatal(err)
        }
        var t_type string 
        switch sw {
            case SW_OPEN:
                t_type = "Open"
            case SW_CLOSED:
                t_type = "Closed"
            case SW_ALL:
                t_type = "All"
            case SW_OVERDUE:
               t_type = "Overdue"
            case SW_DUE:
                t_type = "Open with due date"
        }

        fmt.Printf("%s tasks: %d\n", t_type, count)
        for _, t := range tl {
            t_status := "Open"
            if t.done == 0 && t.duedate.Year() != 1 && t.duedate.Before(time.Now()) {
                t_status = "Overdue"
            }
            if t.done == 1 {
                t_status = "Closed"
            }
            if t.done == 0 {
                if t.duedate.Year() == 1 {
                    fmt.Printf("\t%s %d: %s\n", t_status, t.id, t.description)
                } else {
                    fmt.Printf("\t%s %d: %s, due date: %s\n", t_status, t.id, t.description, t.duedate.Format("2006-01-02"))
                }
            } else {
                fmt.Printf("\t%s %d: %s, completed: %s\n", t_status, t.id, t.description, t.completed.Format("2006-01-02"))
            }
        }
        fmt.Printf("End.\n")
    }

    return
}
