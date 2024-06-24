package main

import (
    "fmt"
)

func PrintHelp() {
    helpString := `
    ## Usage: 
    todo [command] [options] [argument]

Program can be executed without any additional argument (defaults to listing open tasks). Other than that a command must follow with optional switches or arguments.
With add command, if there's a due date, the description can come first or last.
    
    help | h | --help | -h

    add | a [description]
		--duedate | --due | -d [date]

    list | l
        --completed | --closed | -c
        --overdue | -o
        --duedate | --due | -d
        --all | -a
        
    count | c
        --completed | --closed | -c
        --overdue | -o
        --duedate | --due | -d
        --all | -a

    complete | close | do | d [task_id]

    reopen | open | undo | u [task_id]

Examples:
    todo
    todo add "New task" -d "2024-08-13"
    todo l --overdue
    todo count -a 
    todo reopen 3

`
    fmt.Println(helpString)
}
