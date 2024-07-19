
# TODO CLI app
This is a simple Linux command line TODO program.
It does one thing at a time, by design. 

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

## Use

todo [command] [required-argument] [options] [arguments]
	
Program can be executed without any additional argument (defaults to listing open tasks). Other than that a command must follow with optional switches or arguments.
Providing invalid date with --due (e.g. --due -) removes due date.

    help | h | --help | -h

    version | v | --version | -v

    prepare | prep
        --local

    reset 
        --local

    add | a [description] 
        --due [date]
        --priority [number]
        --group [group_name]

    count                     
        --completed | -c
        --overdue | -o
        --all | -a

    list | l                 
        --completed | -c
        --overdue | -o
        --all | -a

    show | s [task_id]
        
    update | u [task_id]         
        --desc [description] 
        --due [date]
        --priority [number]
        --group [group_name]

    complete | c [task_id] 

    reopen | open [task_id]

    delete | del [task_id]

## Examples
```
todo
todo a "New task"
todo add "New task" --due "2024-08-13" --group "Project" --priority 2
todo list --all
todo l -o
todo count -c
todo update 15 --desc "Changed description"
todo u 10 --due - 
todo c 12
todo reopen 3
todo del 5
```

