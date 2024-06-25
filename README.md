
# TODO CLI app
While TODO apps are a dev learning meme, I suddenly found that I'm having trouble keeping track of things and a file doesn't work. Also, I've been playing with go, so it made sense.

Started with flat files, switched to sqlite, maybe will add files later.

## Config
* Config dir: $HOME/.config/todo
* Config file: $HOME/.config/todo/todo.yml
* Default db file: $HOME/.config/todo/todo.db

### Config parameters:
* dblocation: $HOME/.config/todo/
* dbname: "todo"
* dateformat "2006-01-02"


1. **Current functions:**
	* Add task with or without due date
    * Update task description and/or due date
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
	* Set task completed
	* Reopen task
    * Delete task

2. **Planned functions**
    * Priority
	* Task groups (projects)
	* Logging?
	* File support csv/json/something else?

## Use

todo [command] [id] [options] [argument]
	
Program can be executed without any additional argument (defaults to listing open tasks). Other than that a command must follow with optional switches or arguments.

    help | h | --help | -h

    add | a [description] [due]

    count                     
    --completed | -c
    --overdue | -o
    --all | -a

    list | l                 
        --completed | -c
        --overdue | -o
        --all | -a
        
    update | u [id]         
        --desc [description] 
        --due [date]

    complete | c [task_id] 

    reopen | open [task_id]

    delete | del [task_id]

## Examples
```
todo
todo a "New task"
todo add "New task" "2024-08-13"
todo list --all
todo l -o
todo count -c
todo update 15 --due "2024-08-13"
todo c 12
todo reopen 3
todo del 5
```

