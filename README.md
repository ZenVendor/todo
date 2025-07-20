# TODO CLI app
This is a command line TODO program.
It does one thing at a time, by design. 

## Major change is in progress

## Config
There is a config file using YAML format. It is created with prepare command.
Prepare command by default creates the config in $XDG_CONFIG_HOME or $HOME/.config/todo if the former is not set.
Using --local argument creates the files in current directory. 
Main program looks for the config file in the current dir, XDG_CONFIG_HOME and then $HOME/.config/todo  

* Default config file: todo_config.yml
* Default db file: todo.db

### Config parameters:
* dblocation: $HOME/.config/todo/
* dbname: "todo.db"
* dateformat "2006-01-02"

## Functionality

1. **Current functions:**
	* Add task
    * Update task 
    * Display task details
    * Set task completed
    * Reopen task
    * Delete task
	* List tasks (default)
		* open (default)
		* closed
		* all
		* overdue
	* Print task count (to be used for prompt indicator)
		* open (default)
		* closed
		* all
		* overdue

## Usage

todo [verb] [required_value] [args] [kwargs]
	
Program can be executed without any additional argument (defaults to listing open tasks). Other than that a command must follow with optional arguments.
Providing invalid date with --due (e.g. --due -) removes due date.




    add, a [short_description]
        --due, -d [date]
        --group, -g [group_name]
        --long, -l [long_description]
        --priority, -p [number]
        --taskid, -t [task_id]

    complete, c [task_id] 
        --comment, -c [closing_comment]

    configure 
        --local
        --reset

    count                     
        --all, -a
        --completed, -c
        --due, -d
        --overdue, -o

    delete [task_id]

    help, h 

    list, l                 
        --all, -a
        --completed, -c
        --due, -d
        --overdue, -o

    reopen, r [task_id]

    show, s [task_id]
        
    update, u [task_id]         
        --due, -d [date]
        --group, -g [group_name]
        --priority, -p [number]
        --short, -s [short_description] 
        --taskid, -t [task_id]

    version, v 

## Examples
```
todo
todo a "New task"
todo add "New task" --due=2024-08-13 --group="Project" --priority=2
todo list --all
todo l -o
todo count -c
todo update 15 --desc="Changed description"
todo u 10 --due= 
todo c 12
todo reopen 3
todo del 5
``` 


