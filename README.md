
# TODO CLI app
While TODO apps are a dev learning meme, I suddenly found that I'm having trouble keeping track of things and a file doesn't work. Also, I've been playing with go, so it made sense.

Started with flat files, switched to sqlite, maybe will add files later.
1. **Current functions:**
	* Add task with or without due date
	* List tasks (default)
		* open (default)
		* closed
		* all
        * with due date
		* overdue
	* Print task count (to be used for prompt indicator)
		* open (default)
		* closed
		* all
		* overdue
	* Set task completed
	* Reopen task
2. **Planned functions**
	* Task groups (projects)
	* Logging?
	* File support csv/json/something else?

## Use

todo [command] [options] [argument]
	
Program can be executed without any additional argument (defaults to listing open tasks). Other than that a command must follow with optional switches or arguments.
With add command, if there's a due date, the description can come first or last.

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

## Examples
```
todo
todo add "New task" -d "2024-08-13"
todo l --overdue
todo count -a 
todo reopen 3
```

