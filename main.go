package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Task struct {
    id int
    done bool
    desc string
}
type TaskList []Task

func (t Task) String() string {
    return fmt.Sprintf("%d;%t;%s\n", t.id, t.done, t.desc)
}

func printHelp() {
    fmt.Printf("ToDo list\n\n")
    fmt.Printf("Options:\n")
    fmt.Printf("\tadd <description>\tAdd new task.\n")
    fmt.Printf("\tlist [done|all]\tList active (default), done or all tasks.\n")
    fmt.Printf("\tdo <task_id>\tMark task as done.\n")
    fmt.Printf("\tundo <task_id>\tReactivate task.\n")
}


var openTasks TaskList
var closedTasks TaskList

func main () {
    args := os.Args[1:]
    
    if len(args) == 0 || !slices.Contains([]string{"add", "list", "do", "undo"}, args[0]) {
        printHelp()
        return
    }

    file, err := os.OpenFile("tasklist.csv", os.O_RDWR|os.O_CREATE, 0644)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    openCount := 0
    closedCount := 0
    maxId := 0
    var task Task
    s := bufio.NewScanner(file)
    for s.Scan() {
        row := strings.Split(s.Text(), ";")
        task.id, _ = strconv.Atoi(row[0])
        task.done, _ = strconv.ParseBool(row[1])
        task.desc = row[2]

        if task.done {
            closedTasks = append(closedTasks, task)
            closedCount++
        } else {
            openTasks = append(openTasks, task)
            openCount++
        }
        if maxId < task.id {
            maxId = task.id
        }
    }

    if args[0] == "add" {

        if len(strings.TrimSpace(args[1])) == 0 {
            fmt.Printf("Please write description. No task added.\n")
        }
        task = Task{maxId + 1, false, args[1]}
        openCount++
        if _, err = file.WriteString(task.String()); err != nil {
            log.Fatal(err)
        }
        fmt.Printf("Task added with Id %d\n", maxId + 1)
    }
    if args[0] == "list" {
        fmt.Printf("%d open tasks:\n", openCount)
        for _, t := range openTasks {
            fmt.Printf("\t%d: %s\n", t.id, t.desc)
        }
        fmt.Printf("End.\n")
    }

    return
    
}
