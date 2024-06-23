package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

const SW_OPEN = 0
const SW_CLOSED = 1
const SW_OVERDUE = 2
const SW_DUE = 3
const SW_ALL = 4

const CMD_COUNT = 0
const CMD_ADD = 1
const CMD_LIST = 2
const CMD_COMPLETE = 3
const CMD_REOPEN = 4

func printHelp() {
    fmt.Printf("ToDo list\n\n")
    fmt.Printf("Options:\n")
    fmt.Printf("\tadd <description>\tAdd new task.\n")
    fmt.Printf("\tlist [done|all]\tList active (default), done or all tasks.\n")
    fmt.Printf("\tdo <task_id>\tMark task as done.\n")
    fmt.Printf("\tundo <task_id>\tReactivate task.\n\n")
}

func main () {
    
    args := os.Args[1:]

    command := CMD_LIST
    sw := SW_OPEN
    var sw_idx int

    if len(args) != 0 {
        switch args[0] {
        case "count":
            command = CMD_COUNT
        case "add":
            command = CMD_ADD
        case "list":
            command = CMD_LIST
        case "do":
            command = CMD_COMPLETE
        case "undo":
            command = CMD_REOPEN
        default:
            printHelp()
            return
        }
    }

    if len(args) > 1 {
        for i, a := range args[1:] {
            if slices.Contains([]string{"--all", "-a"}, a) {
                sw = SW_ALL
                sw_idx = i
            }
            if slices.Contains([]string{"--closed", "-c"}, a) {
                sw = SW_CLOSED
                sw_idx = i
            }
            if slices.Contains([]string{"--overdue", "-o"}, a) {
                sw = SW_OVERDUE
                sw_idx = i
            }
            if slices.Contains([]string{"--due", "-d"}, a) {
                sw = SW_DUE
                sw_idx = i
            }
        }
    }

    db, err := OpenDB("./")
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

    if command == CMD_ADD {
        var t Task
        if sw == SW_DUE {   
            if len(args) > sw_idx + 2 {
                t.duedate = args[sw_idx + 2]
            }
        }
        desc_idx := 1
        if sw_idx == 0 {
            desc_idx = 3
        }
        if len(strings.TrimSpace(args[desc_idx])) == 0 {
            fmt.Printf("Please write description. No task added.\n")
        }
        fmt.Printf("Adding task: %s\n", args[desc_idx])
        t = Task{0, args[desc_idx], 0, "", "", "", ""}
        if err = t.AddTask(db); err != nil {
            log.Fatal(err)
        }

    }
   
    if command == CMD_LIST {
        count, err := Count(db, SW_OPEN)
        if err != nil {
            log.Fatal(err)
        }
        tl, err := List(db, SW_OPEN)
        if err != nil {
            log.Fatal(err)
        }

        fmt.Printf("Open tasks: %d\n", count)
        for _, t := range tl {
            fmt.Printf("\t%d: %s\n", t.id, t.description)
        }
        fmt.Printf("End.\n")
    }

    return
    
}
