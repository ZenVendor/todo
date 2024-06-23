package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

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
    if len(args) == 0 || !slices.Contains([]string{"add", "list", "do", "undo"}, args[0]) {
        printHelp()
        return
    }

    db, err := OpenDB("./")
    if err != nil {
        log.Fatal(err)
    }

    if args[0] == "add" {
        if len(strings.TrimSpace(args[1])) == 0 {
            fmt.Printf("Please write description. No task added.\n")
        }
        fmt.Printf("Adding task: %s\n", args[1])
        t := Task{0, args[1], 0, "", "", "", ""}
        if err = t.AddTask(db); err != nil {
            log.Fatal(err)
        }

    }
   
    if args[0] == "list" {
        count, err := Count(db, "open")
        if err != nil {
            log.Fatal(err)
        }
        tl, err := List(db, "open")
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
