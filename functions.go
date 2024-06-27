package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"time"

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


type Value struct {
    name string
    value interface{}
}

type Values []Value

type Config struct {
    dbLocation string   `yaml:"dblocation"`
    dbName string       `yaml:"dbname"`
    dateFormat string   `yaml:"dateformat"`
}

func (conf *Config) Prepare() (err error) {
    home, ok := os.LookupEnv("HOME")
    if !ok {
        return fmt.Errorf("HOME environment variable not set.")
    }

    configDir := fmt.Sprintf("%s/.config/todo", home)
    configFile := fmt.Sprintf("%s/todo.yml", configDir)

    conf.dbLocation = configDir
    conf.dbName = "todo"
    conf.dateFormat = "2006-01-02"

    writeConf := fmt.Sprintf("dblocation: %s\ndbname: %s\ndateformat: %s", conf.dbLocation, conf.dbName, conf.dateFormat)

    if _, err = os.Stat(configDir); os.IsNotExist(err) {
        if err = os.MkdirAll(configDir, 0700); err != nil {
            return
        }
    }
    if _, err = os.Stat(configFile); os.IsNotExist(err) {
        if err = os.WriteFile(configFile, []byte(writeConf), 0700); err != nil {
            return
        }
    }
    f, err := os.ReadFile(configFile)
    if err != nil {
        return
    }

    err = yaml.Unmarshal(f, &conf)

    return
}

func (vs Values) ReadValue(name string) interface{} {
    idx := slices.IndexFunc(vs, func(v Value) bool {
        return v.name == name
    })
    if idx == -1 {
        return nil
    }
    return vs[idx].value
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
                    if "--desc" == args[1] && "--due" == args[3] {
                        dd, err := time.Parse(dateFormat, args[4])
                        if err != nil {
                            dd, _ = time.Parse(dateFormat, "0001-01-01")
                        }
                        values = append(values, Value{"id", taskId})
                        values = append(values, Value{"description", args[2]})
                        values = append(values, Value{"due", dd})
                        valid = true
                    }
                    if "--desc" == args[3] && "--due" == args[1] {
                        dd, err := time.Parse(dateFormat, args[2])
                        if err != nil {
                            dd, _ = time.Parse(dateFormat, "0001-01-01")
                        }
                        values = append(values, Value{"id", taskId})
                        values = append(values, Value{"description", args[4]})
                        values = append(values, Value{"due", dd})
                        valid = true
                    }
                }
            }
        }
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
In ADD command, description is required and must be provided first.
In commands that require it, task ID must follow the command.
Values following switches can be provided in any order.

    
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
